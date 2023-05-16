package collector

import (
	"encoding/xml"
	"io"
	"log"
	"strconv"
	"strings"
	"sync"

	"github.com/antchfx/xmlquery"
	"github.com/prometheus/client_golang/prometheus"
)

const (
	scrapeSubsystem = ""
)

//Metric describes a hpmsa metric
type Metric struct {
	Sources []Source `yaml:"sources"`
}

// Source describes a metric path
type Source struct {
	Path    string               `yaml:"path"`
	XPath   string               `yaml:"xpath"`
	Objects map[string][]*Object `yaml:"objects"`
}

// Object describes a specific xml object
type Object struct {
	Name     string            `yaml:"name"`
	Desc     string            `yaml:"desc"`
	Property string            `yaml:"property"`
	Labels   map[string]string `yaml:"labels"`
	LabelMap map[string]string `yaml:"label_map"`
}

type metricCollector struct {
}

func init() {
	registerCollector("hpmsa_api_collector", defaultEnabled, MetricCollector)
}

//MetricCollector returns a new collector
func MetricCollector() (Collector, error) {
	return &metricCollector{}, nil
}

// Update implements Collector
func (c *metricCollector) Update(ch chan<- prometheus.Metric, target *target) error {
	var wg sync.WaitGroup
	for _, source := range target.Metrics.Sources {
		wg.Add(1)
		go func(source Source) {
			defer wg.Done()
			b, err := APICall(target.Client, target.Hostname, target.SessionKey, source.Path)
			if err != nil {
				log.Print("ERROR: ", err)
				return
			}
			// Query XML response for target objects matching specified XPath
			var r io.Reader
			r = strings.NewReader(string(b))
			root, err := xmlquery.Parse(r)
			if err != nil {
				log.Print("ERROR: ", err)
				return
			}
			var nodes []*xmlquery.Node
			nodes = xmlquery.Find(root, source.XPath)
			// Unmarshal into a Response struct that only contains the objects we care about
			x := &Response{}
			for _, node := range nodes {
				o := &object{}
				err = xml.Unmarshal([]byte(node.OutputXML(true)), o)
				if err != nil {
					log.Print("ERROR: ", err)
					return
				}
				x.Object = append(x.Object, *o)
			}
			if err := ParseXML(ch, x, source); err != nil {
				log.Print("ERROR: ", err)
				return
			}
		}(source)
	}
	wg.Wait()
	return nil
}

// ParseXML will parse the result from the HPMSA endpoint and present it i n prometheus readable format
func ParseXML(ch chan<- prometheus.Metric, x *Response, source Source) error {
	for _, object := range x.Object {
		if objs, ok := source.Objects[object.Name]; ok {
			for _, obj := range objs {
				var value string
				var labelKeys, labelVals []string
				for key, val := range obj.Labels {
					labelKeys = append(labelKeys, key)
					labelVals = append(labelVals, val)
				}
				ParseObject(object, obj, &value, &labelKeys, &labelVals)
				fval, err := strconv.ParseFloat(value, 64)
				if err != nil {
					log.Printf("ERROR: Unable to parse the value for %s which is %v: %v", obj.Name, value, err)
					continue
				}
				ch <- prometheus.MustNewConstMetric(
					prometheus.NewDesc(
						prometheus.BuildFQName(namespace, scrapeSubsystem, obj.Name),
						obj.Desc, labelKeys, nil,
					), prometheus.GaugeValue, fval, labelVals...)
			}
		}
	}
	return nil
}

// ParseObject will parse individual objects for their properties and also loop into resettable-statistics
func ParseObject(object object, obj *Object, value *string, labelKeys *[]string, labelVals *[]string) error {
	for _, property := range object.Property {
		if property.Name == obj.Property {
			*value = property.Value
		}
		if _, ok := obj.LabelMap[property.Name]; ok {
			var alreadyExists bool = false
			for _, k := range *labelKeys {
				if k == obj.LabelMap[property.Name] {
					alreadyExists = true
					break
				}
			}
			if !alreadyExists {
				*labelKeys = append(*labelKeys, obj.LabelMap[property.Name])
				*labelVals = append(*labelVals, property.Value)
			}
		}
	}
	for _, subobject := range object.Object {
		if subobject.Name == "resettable-statistics" {
			ParseObject(subobject, obj, value, labelKeys, labelVals)
		}
	}
	return nil
}

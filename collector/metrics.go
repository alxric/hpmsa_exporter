package collector

import (
	"encoding/xml"
	"log"
	"strconv"
	"sync"

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
			x := &Response{}
			err = xml.Unmarshal(b, x)
			if err != nil {
				log.Print("ERROR: ", err)
				return
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
				for _, property := range object.Property {
					if property.Name == obj.Property {
						value = property.Value
					}
					if _, ok := obj.LabelMap[property.Name]; ok {
						labelKeys = append(labelKeys, obj.LabelMap[property.Name])
						labelVals = append(labelVals, property.Value)
					}
				}
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

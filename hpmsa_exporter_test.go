package main

import (
	"encoding/xml"
	"sync"
	"testing"

	"github.com/alxric/hpmsa_exporter/collector"
	"github.com/prometheus/client_golang/prometheus"
	dto "github.com/prometheus/client_model/go"
	yaml "gopkg.in/yaml.v2"
)

func TestSystemHealth(t *testing.T) {
	var data = `
path: "system"
objects:
  system-information:
    - name: "system_health"
      desc: "System health"
      property: "health-numeric"
`
	b := `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<RESPONSE VERSION="L100" REQUEST="show system">
  <OBJECT basetype="system" name="system-information" oid="1" format="pairs">
    <PROPERTY name="system-name" type="string" size="80" draw="true" sort="string" display-name="System Name">TEST-MSA-01</PROPERTY>
    <PROPERTY name="system-contact" type="string" size="80" draw="true" sort="string" display-name="System Contact">Uninitialized Contact</PROPERTY>
    <PROPERTY name="system-location" type="string" size="80" draw="true" sort="string" display-name="System Location">TEST</PROPERTY>
    <PROPERTY name="system-information" type="string" size="80" draw="true" sort="string" display-name="System Information">Uninitialized Info</PROPERTY>
    <PROPERTY name="midplane-serial-number" type="string" size="33" draw="true" sort="string" display-name="Midplane Serial Number">00C0FF27FA89</PROPERTY>
    <PROPERTY name="vendor-name" type="string" size="17" draw="true" sort="string" display-name="Vendor Name">HP</PROPERTY>
    <PROPERTY name="product-id" type="string" size="17" draw="true" sort="string" display-name="Product ID">MSA 1040 SAS</PROPERTY>
    <PROPERTY name="product-brand" type="string" size="16" draw="true" sort="string" display-name="Product Brand">MSA Storage</PROPERTY>
    <PROPERTY name="scsi-vendor-id" type="string" size="17" draw="true" sort="string" display-name="SCSI Vendor ID">HP</PROPERTY>
    <PROPERTY name="scsi-product-id" type="string" size="17" draw="true" sort="string" display-name="SCSI Product ID">MSA 1040 SAS</PROPERTY>
    <PROPERTY name="enclosure-count" type="uint32" size="17" draw="true" sort="string" display-name="Enclosure Count">1</PROPERTY>
    <PROPERTY name="health" type="string" size="10" draw="true" sort="string" display-name="Health">OK</PROPERTY>
    <PROPERTY name="health-numeric" type="uint32" size="10" draw="true" sort="string" display-name="Health">0</PROPERTY>
    <PROPERTY name="health-reason" type="string" size="80" draw="true" sort="string" display-name="Health Reason"></PROPERTY>
    <PROPERTY name="other-MC-status" type="string" size="40" draw="true" sort="nosort" display-name="Other MC Status">Operational</PROPERTY>
    <PROPERTY name="other-MC-status-numeric" type="uint32" size="40" draw="true" sort="nosort" display-name="Other MC Status">4760</PROPERTY>
    <PROPERTY name="pfuStatus" type="string" size="80" draw="true" sort="string" display-name="PFU Status">Idle</PROPERTY>
    <PROPERTY name="supported-locales" type="string" size="200" draw="true" sort="string" display-name="Supported Locales">English (English), Arabic (العربية), Portuguese (português), Spanish (español), French (français), German (Deutsch), Italian (italiano), Japanese (日本語), Korean (한국어), Dutch (Nederlands), Russian (русский), Chinese-Simplified (简体中文), Chinese-Traditional (繁體中文)</PROPERTY>
    <PROPERTY name="current-node-wwn" type="string" size="17" draw="false" sort="string" display-name="Current Node WWN">500c0ff27fa89000</PROPERTY>
    <PROPERTY name="fde-security-status" type="string" size="20" draw="true" sort="string" display-name="FDE Security Status">Unsecured</PROPERTY>
    <PROPERTY name="fde-security-status-numeric" type="uint32" size="20" draw="true" sort="string" display-name="FDE Security Status">1</PROPERTY>
    <PROPERTY name="platform-type" type="string" size="20" draw="false" sort="string" display-name="Platform Type">Gallium</PROPERTY>
    <PROPERTY name="platform-type-numeric" type="uint32" size="20" draw="false" sort="string" display-name="Platform Type">3</PROPERTY>
    <PROPERTY name="platform-brand" type="string" size="20" draw="false" sort="string" display-name="Platform Brand">HP Cardinals</PROPERTY>
    <PROPERTY name="platform-brand-numeric" type="uint32" size="20" draw="false" sort="string" display-name="Platform Brand">15</PROPERTY>
    <OBJECT basetype="redundancy" name="system-redundancy" oid="2" format="pairs">
      <PROPERTY name="redundancy-mode" type="string" size="140" draw="true" sort="nosort" display-name="Controller Redundancy Mode">Active-Active ULP</PROPERTY>
      <PROPERTY name="redundancy-mode-numeric" type="uint32" size="140" draw="true" sort="nosort" display-name="Controller Redundancy Mode">8</PROPERTY>
      <PROPERTY name="redundancy-status" type="string" size="140" draw="true" sort="nosort" display-name="Controller Redundancy Status">Redundant</PROPERTY>
      <PROPERTY name="redundancy-status-numeric" type="uint32" size="140" draw="true" sort="nosort" display-name="Controller Redundancy Status">2</PROPERTY>
      <PROPERTY name="controller-a-status" type="string" size="40" draw="true" sort="nosort" display-name="Controller A Status">Operational</PROPERTY>
      <PROPERTY name="controller-a-status-numeric" type="uint32" size="40" draw="true" sort="nosort" display-name="Controller A Status">0</PROPERTY>
      <PROPERTY name="controller-a-serial-number" type="string" size="32" draw="true" sort="nosort" display-name="Controller A Serial Number">7CE601M063</PROPERTY>
      <PROPERTY name="controller-b-status" type="string" size="40" draw="true" sort="nosort" display-name="Controller B Status">Operational</PROPERTY>
      <PROPERTY name="controller-b-status-numeric" type="uint32" size="40" draw="true" sort="nosort" display-name="Controller B Status">0</PROPERTY>
      <PROPERTY name="controller-b-serial-number" type="string" size="32" draw="true" sort="nosort" display-name="Controller B Serial Number">7CE603P279</PROPERTY>
      <PROPERTY name="other-MC-status" type="string" size="40" draw="true" sort="nosort" display-name="Other MC Status">Operational</PROPERTY>
      <PROPERTY name="other-MC-status-numeric" type="uint32" size="40" draw="true" sort="nosort" display-name="Other MC Status">4760</PROPERTY>
    </OBJECT>
  </OBJECT>
  <OBJECT basetype="status" name="status" oid="3">
    <PROPERTY name="response-type" type="string" size="12" draw="false" sort="nosort" display-name="Response Type">Success</PROPERTY>
    <PROPERTY name="response-type-numeric" type="uint32" size="12" draw="false" sort="nosort" display-name="Response Type">0</PROPERTY>
    <PROPERTY name="response" type="string" size="180" draw="true" sort="nosort" display-name="Response">Command completed successfully. (2018-05-07 08:25:22)</PROPERTY>
    <PROPERTY name="return-code" type="sint32" size="15" draw="false" sort="nosort" display-name="Return Code">0</PROPERTY>
    <PROPERTY name="component-id" type="string" size="80" draw="false" sort="nosort" display-name="Component ID"></PROPERTY>
    <PROPERTY name="time-stamp" type="string" size="25" draw="false" sort="datetime" display-name="Time">2018-05-07 08:25:22</PROPERTY>
    <PROPERTY name="time-stamp-numeric" type="uint32" size="25" draw="false" sort="datetime" display-name="Time">1525681522</PROPERTY>
  </OBJECT>
</RESPONSE>`
	x := &collector.Response{}
	err := xml.Unmarshal([]byte(b), x)
	if err != nil {
		t.Errorf("Could not unmarshal system xml: %v", err)
	}
	y := collector.Source{}
	err = yaml.Unmarshal([]byte(data), &y)
	if err != nil {
		t.Errorf("Could not unmarshal config: %v", err)
	}
	ch := make(chan prometheus.Metric)
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()

		resp := <-ch
		m := &dto.Metric{}
		resp.Write(m)
		if m.GetGauge().GetValue() != 0 {
			t.Errorf("Invalid system health value: %v", m.GetGauge().GetValue())
		}
	}()
	if err := collector.ParseXML(ch, x, y); err != nil {
		t.Errorf("Could not Parse xml: %v", err)
	}
	wg.Wait()
}

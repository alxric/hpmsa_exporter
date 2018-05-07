# HPMSA Exporter

HPMSA Exporter is a Prometheus exporter used to gather metrics from HP MSA Storage units

### Usage

Start the exporter by simply executing

	./hpmsa_exporter

It will by default look for a configuration file called **config.yml**. This file contains authentication credentials, and
scrape config. You can use the one supplied in this repository or make your own,, as long as the format is the same.
If you want to change the location of the config file, do

	./hpmsa_exporter --config.file config2.yml

When your exporter is up and running, we can test it with curl

	curl localhost:9114/metrics?target=test-msa-01.example.com

If you want to supply a specific config file for a specific host, do

	curl localhost:9114/metrics?target=test-msa-01.example.com&config=config3.yml

### Prometheus configuration

	- job_name: hpmsa
	  scrape_interval: 30s
	  scrape_timeout: 10s
	  metrics_path: /metrics
	  scheme: http
	  static_configs:
	  - targets:
	    - test-msa-01.example.com
	  relabel_configs:
	  - source_labels: [__address__]
		separator: ;
		regex: (.*)
		target_label: __param_target
		replacement: $1
		action: replace
	  - source_labels: [__param_target]
		separator: ;
		regex: (.*)
		target_label: instance
		replacement: $1
		action: replace
	  - separator: ;
		regex: (.*)
		target_label: __address__
		replacement: localhost:9114
		action: replace

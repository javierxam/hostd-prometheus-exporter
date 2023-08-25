# hostd-prometheus-exporter
## Flexible exporter to send Sia's hostd metrics to Prometheus


hostd-prometheus-exporter, a nimble and adaptable metrics exporter that delivers Sia metrics to Prometheus for processing and Grafana for plotting and notifying. hostd exporter simplifies the monitoring of one or multiple Sia instances using Prometheus, Grafana, and other tools.

In short: hostd_prometheus_exporter enables anyone to design custom Sia charts, alerts, and dashboards.


```
$> ./hostd-prometheus-exporter -h
 Usage of ./hostd-prometheus-exporter:
  -address string
        Hostd API address (default "127.0.0.1:9980")
  -passwd string
        Hostd API password (default "Sia is Awesome")
  -port int
        Port to serve Prometheus Metrics on (default 8101)
  -refresh int
        Frequency to get Metrics from Hostd (minutes) (default 1)
```

# Monitoring in Kyma using a custom kube-prometheus-stack

> [!WARNING]
> This guide got reworked and has been moved. Find the updated instructions in [Integrate With Prometheus](https://kyma-project.io/#/telemetry-manager/user/integration/prometheus/README).

## Overview

The Kyma Telemetry module provides collection and integration into observability backends based on OpenTelemetry. OpenTelemetry is the new vendor-neutral player in the cloud-native observability domain and has a growing adoption. However, it still lacks many features and sometimes the typical `kube-prometheus-stack` based on Prometheus, Grafana, and the surrounding helper tools is more appropriate.

The following instructions describe the complete monitoring flow for your service running in Kyma. You get the gist of monitoring applications, such as Prometheus, Grafana, and Alertmanager. You learn how and where you can observe and visualize your service metrics to monitor them for any alerting values.

All the tutorials use the [`monitoring-custom-metrics`](./monitoring-custom-metrics/README.md) example and one of its services called `sample-metrics`. This service exposes the `cpu_temperature_celsius` custom metric on the `/metrics` endpoint. This custom metric is the central element of the whole tutorial set. The metric value simulates the current processor temperature and changes randomly from 60 to 90 degrees Celsius. The alerting threshold in these tutorials is 75 degrees Celsius. If the temperature exceeds this value, the Grafana dashboard, PrometheusRule, and Alertmanager notifications you create inform you about this.

## Sequence of tasks

The instructions cover the following tasks:

 ![Monitoring tutorials](./assets/monitoring-tutorials.svg)

1. [**Deploy a custom Prometheus stack**](./prometheus.md), in which you deploy the [kube-prometheus-stack](https://github.com/prometheus-operator/kube-prometheus) from the upstream Helm chart.

2. [**Observe application metrics**](./monitoring-custom-metrics/README.md), in which you redirect the `cpu_temperature_celsius` metric to the localhost and the Prometheus UI. You later observe how the metric value changes in the predefined 10-second interval in which Prometheus scrapes the metric values from the service's `/metrics` endpoint.

3. [**Create a Grafana dashboard**](./monitoring-grafana-dashboard/README.md), in which you create a Grafana dashboard of a Gauge type for the `cpu_temperature_celsius` metric. This dashboard shows explicitly when the CPU temperature is equal to or higher than the predefined threshold of 75 degrees Celsius, at which point the dashboard turns red.

4. [**Define alerting rules**](./monitoring-alert-rules/README.md), in which you define the `CPUTempHigh` alerting rule by creating a PrometheusRule resource. Prometheus accesses the `/metrics` endpoint every 10 seconds and validates the current value of the `cpu_temperature_celsius` metric. If the value is equal to or higher than 75 degrees Celsius, Prometheus waits for 10 seconds to recheck it. If the value still exceeds the threshold, Prometheus triggers the rule. You can observe both the rule and the alert it generates on the Prometheus dashboard.

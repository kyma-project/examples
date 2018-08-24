# Expose Custom Metrics in Kyma

## Overview

This example shows how to expose custom metrics to Prometheus with a golang service in Kyma.

1. Expose a sample application serving metrics on port 8081
2. Access the exposed metrics in Prometheus

## Prerequisites

- Kyma as the target deployment environment.
- If sidecar injection is not enabled for **default** namespace, then run the following:
    ```bash
    kubectl label namespace default istio-injection=enabled
    ```

## Installation

### Expose a sample metrics application

- Deploy the application, service and servicemonitor
    ```bash
    kubectl apply -f k8s -R
    ```
    
#### Access the exposed metrics in Prometheus

- Run the port-forward on service core-prometheus
    
    ```bash
    kubectl port-forward -n kyma-system svc/core-prometheus 9090:9090
    ```
In [targets](http://localhost:9090/targets#job-sample-metrics-8081) all the **sample-metrics** endpoint will appear.

- Use either one of the registered metrics `cpu_temperature_celsius` or `hd_errors_total` in the `expression` field in [here](http://localhost:9090/graph) and `Execute` to check the values scrapped by Prometheus.

## Cleanup
Run the following commands to completely remove the example and all its resources from the cluster:

- Remove **label** istio-injection from **default** namespace
    ```bash
    kubectl label namespace default istio-injection-
    ```
- Remove **ServiceMonitor** in namespace kyma-system
    ```bash
    kubectl delete servicemonitor -l example=monitoring-custom-metrics -n kyma-system
    ```
- Remove sample-metrics deployments in namespace default
    ```bash
    kubectl delete all -l example=monitoring-custom-metrics
    ```

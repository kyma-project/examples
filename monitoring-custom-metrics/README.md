# Expose Custom Metrics in Kyma

## Overview

This example shows how to expose custom metrics to Prometheus with a Golang service in Kyma. To do so, follow these steps:

1. Expose a sample application serving metrics on the `8081` port.
2. Access the exposed metrics in Prometheus.

## Prerequisites

- Kyma as the target deployment environment.
- If sidecar injection is not enabled for the `default` Namespace, run the following command:
    ```bash
    kubectl label namespace default istio-injection=enabled
    ```

## Installation

### Expose a sample metrics application

- Deploy the application, service, and servicemonitor:
    ```bash
    kubectl apply -f deployment -R
    ```
    
### Access the exposed metrics in Prometheus

- Run the `port-forward` command on the `core-prometheus` service:
    
    ```bash
    kubectl port-forward -n kyma-system svc/core-prometheus 9090:9090
    ```
All the **sample-metrics** endpoints appear as the [`Targets`](http://localhost:9090/targets#job-sample-metrics-8081) list.

- Use either the `cpu_temperature_celsius` or `hd_errors_total` in the **expression** field [here](http://localhost:9090/graph).
- Click the **Execute** button to check the values scrapped by Prometheus.

### Cleanup
Run the following commands to completely remove the example and all its resources from the cluster:

1. Remove the **istio-injection** label from the `default` Namespace.
    ```bash
    kubectl label namespace default istio-injection-
    ```
2. Remove **ServiceMonitor** in the `kyma-system` Namespace.
    ```bash
    kubectl delete servicemonitor -l example=monitoring-custom-metrics -n kyma-system
    ```
3. Remove the `sample-metrics` Deployments in the `default` Namespace.
    ```bash
    kubectl delete all -l example=monitoring-custom-metrics
    ```

# Expose Custom Metrics in Kyma

## Overview

This example shows how to expose custom metrics to Prometheus with a Golang service in Kyma. To do so, follow these steps:

1. Expose a sample application serving metrics on the `8081` port.
2. Access the exposed metrics in Prometheus.

## Prerequisites

- Kyma as the target deployment environment.
- A Namespace to which you deploy the example with the `env: "true"` label. For more information, read the [related documentation](https://github.com/kyma-project/kyma/blob/master/docs/kyma/docs/03-02-namespaces.md).

## Installation

### Expose a sample metrics application

1. Export your Namespace as a variable. Replace the `{namespace}` placeholder in the following command and run it:

    ```bash
    export KYMA_EXAMPLE_NS="{namespace}"
    ```
2. Deploy the service:
    ```bash
    kubectl apply -f deployment/deployment.yaml -n $KYMA_EXAMPLE_NS
    ```
3. Deploy the ServiceMonitor:
    ```bash
    kubectl apply -f deployment/service-monitor.yaml
    ```

### Access the exposed metrics in Prometheus

1. Run the `port-forward` command on the `core-prometheus` service:

    ```bash
    kubectl port-forward -n kyma-system svc/monitoring-prometheus 9090:9090
    ```
All the **sample-metrics** endpoints appear as the [`Targets`](http://localhost:9090/targets#job-sample-metrics-8081) list.

2. Use either `cpu_temperature_celsius` or `hd_errors_total` in the **expression** field [here](http://localhost:9090/graph).
3. Click the **Execute** button to check the values scrapped by Prometheus.

### Cleanup
Run the following commands to completely remove the example and all its resources from the cluster:

1. Remove ServiceMonitor in the `kyma-system` Namespace.
    ```bash
    kubectl delete servicemonitor -l example=monitoring-custom-metrics -n kyma-system
    ```
2. Run the following command to completely remove the example service and all its resources from the cluster:
    ```bash
    kubectl delete all -l example=monitoring-custom-metrics -n $KYMA_EXAMPLE_NS
    ```

# Expose Custom Metrics in Kyma

> [!WARNING]
> This guide has been revised and moved. Find the updated instructions in [Integrate With Prometheus](https://kyma-project.io/#/telemetry-manager/user/integration/prometheus/README).


## Overview

This example shows how to expose custom metrics to Prometheus with a Golang service in Kyma. To do so, follow these steps:

1. Expose a sample application serving metrics on the `8081` port.
2. Access the exposed metrics in Prometheus.

## Prerequisites

- Kyma as the target deployment environment.
- Deployed custom kube-prometheus-stack as described in the [prometheus](../) example.

## Installation

### Expose a sample metrics application

1. Export your Namespace as a variable. Replace the `{namespace}` placeholder in the following command and run it:

    ```bash
    export K8S_NAMESPACE="{namespace}"
    ```

2. Ensure that your Namespace has Istio sidecar injection enabled. This example assumes that the metrics are exposed in a strict mTLS mode:

   ```bash
   kubectl label namespace ${K8S_NAMESPACE} istio-injection=enabled
   ```

3. Deploy the service:

    ```bash
    kubectl apply -f https://raw.githubusercontent.com/kyma-project/examples/main/prometheus/monitoring-custom-metrics/deployment/deployment.yaml -n $K8S_NAMESPACE
    ```

4. Deploy the ServiceMonitor:

    ```bash
    kubectl apply -f https://raw.githubusercontent.com/kyma-project/examples/main/prometheus/monitoring-custom-metrics/deployment/service-monitor.yaml
    ```

### Access the exposed metrics in Prometheus

1. Run the `port-forward` command on the `monitoring-prometheus` service:

    ```bash
    kubectl -n ${K8S_NAMESPACE} port-forward $(kubectl -n ${K8S_NAMESPACE} get service -l app=kube-prometheus-stack-prometheus -oname) 9090
    ```

All the **sample-metrics** endpoints appear as `Targets` under `http://localhost:9090/targets#job-sample-metrics` list.

2. Use either `cpu_temperature_celsius` or `hd_errors_total` in the **expression** field under `http://localhost:9090/graph`.
3. Click the **Execute** button to check the values scraped by Prometheus.

### Cleanup

Run the following commands to completely remove the example and all its resources from the cluster:

1. Remove ServiceMonitor in the `kyma-system` Namespace.

    ```bash
    kubectl delete servicemonitor -l example=monitoring-custom-metrics -n kyma-system
    ```

2. Run the following command to completely remove the example service and all its resources from the cluster:

    ```bash
    kubectl delete all -l example=monitoring-custom-metrics -n $K8S_NAMESPACE
    ```

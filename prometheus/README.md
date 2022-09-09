# Install a custom kube-prometheus-stack in Kyma

## Overview

The Kyma monitoring stack often brings limited configuration options in contrast to the upstream [`kube-prometheus-stack`](https://github.com/prometheus-community/helm-charts/blob/main/charts/kube-prometheus-stack) chart. Modifications might be resetted at next upgrade cycle.

An alternative can be a parallel installation of the upstream chart offering all customization options. This tutorial outlines how to achieve such installation in co-existent to the Kyma monitoring stack.

## Prerequisites

- Kyma as the target deployment environment.
- kubectl > 1.22.x
- helm 3.x

## Installation

### Preparation

1. Export your Namespace as a variable. Replace the `{namespace}` placeholder in the following command and run it:

    ```bash
    export KYMA_EXAMPLE_NS="{namespace}"
    ```

1. Export the Helm release name which you want to use. It can be any name,, be aware that all resources in the cluster will be prefixed with that name. Replace the `{release-name}` placeholder in the following command and run it:
    ```bash
    export HELM_RELEASE_NAME="{release-name}"
    ```

2. Update your helm installation with the required helm repository:

    ```bash
    helm repo add prometheus-community https://prometheus-community.github.io/helm-charts
    helm repo update
    ```

### Install the kube-prometheus-stack

1. Run the Helm upgrade command which will install the chart if not present yet. Change the grafana admin password (at the end of the command) to some value of your choice
    ```bash
    helm upgrade --install -n ${KYMA_EXAMPLE_NS} ${HELM_RELEASE_NAME} prometheus-community/kube-prometheus-stack -f values.yaml --set grafana.adminPassword=myPwd
    ```

Hereby, use the [values.yaml](./values.yaml) provided with this tutorial which contains customized settings derivating from the default settings.

### Verify the installation

1. You should see several pods coming up in the namespace, especially prometheus and alertmanager. Assure that all pods are ending in a "Running" state.
2. Browse the prometheus dashboard and verify that all "Status->Targets" are healthy. Following command will expose the dashboard on `http://localhost:9090`
   ```bash
   kubectl -n ${KYMA_EXAMPLE_NS} port-forward svc/${HELM_RELEASE_NAME}-kube-prometheus-stack-prometheus 9090
   ```
2. Browse the grafana dashboard and verify that the dashboards are showing data. Following command will expose the dashboard on `http://localhost:3000`
   ```bash
   kubectl -n ${KYMA_EXAMPLE_NS} port-forward svc/${HELM_RELEASE_NAME}-grafana 3000:80
   ```

### Deploy a custom workload and scrape it

1. Deploy a custom application which exposes metrics:
  ```bash
  kubectl create namespace myapp
  kubectl label namespace myapp istio-injection=enabled
  kubectl apply -n myapp -f https://raw.githubusercontent.com/kyma-project/examples/main/monitoring-custom-metrics/deployment/deployment.yaml
  ```

2. Configure prometheus to scrape the new target by defining a ServiceMonitor:
  ```bash
  kubectl apply -n myapp -f - <<EOF
  apiVersion: monitoring.coreos.com/v1
  kind: ServiceMonitor
  metadata:
    name: metrics
  spec:
    selector:
      matchLabels:
        app: sample-metrics
    endpoints:
      - port: http
        scheme: https
        tlsConfig: 
        caFile: /etc/prometheus/secrets/istio.default/root-cert.pem
        certFile: /etc/prometheus/secrets/istio.default/cert-chain.pem
        keyFile: /etc/prometheus/secrets/istio.default/key.pem
        insecureSkipVerify: true  # Prometheus does not support Istio security naming, thus skip verifying target Pod certificate
  EOF
  ```

3. Wait a while and check that metric `cpu_temperature_celsius` is explorable in Grafana

### Cleanup

Run the following commands to completely remove the example and all its resources from the cluster:

1. Remove the stack by calling helm:

    ```bash
    helm delete -n ${KYMA_EXAMPLE_NS} ${HELM_RELEASE_NAME}
    ```
2. Remove the sample application:
   
   ```bash
   kubectl delete namespace myapp

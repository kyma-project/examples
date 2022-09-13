# Install a custom kube-prometheus-stack in Kyma

## Overview

The Kyma monitoring stack often brings limited configuration options in contrast to the upstream [`kube-prometheus-stack`](https://github.com/prometheus-community/helm-charts/blob/main/charts/kube-prometheus-stack) chart. Modifications might be resetted at next upgrade cycle.

An alternative can be a parallel installation of the upstream chart offering all customization options. This tutorial outlines how to achieve such installation in co-existent to the Kyma monitoring stack.

Be aware of that this tutorial describes a basic setup which is not to be used in production. Further configuration is usually required like optimizations in regards to the amount of data to scrape and the required resource footprint of the installation. Even a different setup might be required to achieve qualities like high availability, scalability or durable long-term storage.

## Prerequisites

- Kyma as the target deployment environment.
- Kubectl > 1.22.x
- Helm 3.x

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
    helm upgrade --install -n ${KYMA_EXAMPLE_NS} ${HELM_RELEASE_NAME} prometheus-community/kube-prometheus-stack -f https://raw.githubusercontent.com/kyma-project/examples/main/prometheus/values.yaml --set grafana.adminPassword=myPwd
    ```

Hereby, use the [values.yaml](./values.yaml) provided with this tutorial which contains customized settings deviating from the default settings, or create your own one.
The provided values.yaml covers the following adjustments:
- Support parallel operation to a Kyma monitoring stack
- Scraping of Istio strict mTLS workloads

### Verify the installation

1. You should see several pods coming up in the namespace, especially prometheus and alertmanager. Assure that all pods are ending in a "Running" state.
2. Browse the prometheus dashboard and verify that all "Status->Targets" are healthy. Following command will expose the dashboard on `http://localhost:9090`
   ```bash
   kubectl -n ${KYMA_EXAMPLE_NS} port-forward svc/${HELM_RELEASE_NAME}-kube-prometheus-stack-prometheus 9090
   ```
2. Browse the grafana dashboard and verify that the dashboards are showing data. The user `admin` is pre-configured in the helm chart, the password was provided in your helm install command. Following command will expose the dashboard on `http://localhost:3000`:
   ```bash
   kubectl -n ${KYMA_EXAMPLE_NS} port-forward svc/${HELM_RELEASE_NAME}-grafana 3000:80
   ```

### Deploy a custom workload and scrape it

1. Follow the tutorial [monitoring-custom-metrics](./../monitoring-custom-metrics/) but use the steps above for verifying that the metrics are collected.

### Cleanup

Run the following commands to completely remove the example and all its resources from the cluster:

1. Remove the stack by calling helm:

    ```bash
    helm delete -n ${KYMA_EXAMPLE_NS} ${HELM_RELEASE_NAME}
    ```

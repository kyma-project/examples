# Install a custom kube-prometheus-stack in Kyma

## Overview

The Kyma monitoring stack often brings limited configuration options in contrast to the upstream [`kube-prometheus-stack`](https://github.com/prometheus-community/helm-charts/blob/main/charts/kube-prometheus-stack) chart. Modifications might be reset at the next upgrade cycle.

As an alternative, you can install the upstream chart with all customization options parallel. This tutorial outlines how to set up such installation in co-existence to the Kyma monitoring stack.

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

2. Export the Helm release name that you want to use. It can be any name, but be aware that all resources in the cluster will be prefixed with that name. Replace the `{release-name}` placeholder in the following command and run it:
    ```bash
    export HELM_RELEASE_NAME="{release-name}"
    ```

3. Update your Helm installation with the required Helm repository:

    ```bash
    helm repo add prometheus-community https://prometheus-community.github.io/helm-charts
    helm repo update
    ```

### Install the kube-prometheus-stack

1. Run the Helm upgrade command, which installs the chart if it's not present yet. At the end of the command, change the Grafana admin password to some value of your choice.
    ```bash
    helm upgrade --install -n ${KYMA_EXAMPLE_NS} ${HELM_RELEASE_NAME} prometheus-community/kube-prometheus-stack -f https://raw.githubusercontent.com/kyma-project/examples/main/prometheus/values.yaml --set grafana.adminPassword=myPwd
    ```

You can use the [values.yaml](./values.yaml) provided with this tutorial, which contains customized settings deviating from the default settings, or create your own one.
The provided `values.yaml` covers the following adjustments:
- Support parallel operation to a Kyma monitoring stack
- Support the scraping of workload secured with Istio strict mTLS

### Verify the installation

1. You should see several Pods coming up in the Namespace, especially Prometheus and Alertmanager. Assure that all Pods have the "Running" state.
2. Browse the Prometheus dashboard and verify that all "Status->Targets" are healthy. The following command exposes the dashboard on `http://localhost:9090`:
   ```bash
   kubectl -n ${KYMA_EXAMPLE_NS} port-forward svc/${HELM_RELEASE_NAME}-kube-prometheus-stack-prometheus 9090
   ```
3. Browse the Grafana dashboard and verify that the dashboards are showing data. The user `admin` is pre-configured in the Helm chart; the password was provided in your `helm install` command. The following command exposes the dashboard on `http://localhost:3000`:
   ```bash
   kubectl -n ${KYMA_EXAMPLE_NS} port-forward svc/${HELM_RELEASE_NAME}-grafana 3000:80
   ```

### Deploy a custom workload and scrape it

1. Follow the tutorial [monitoring-custom-metrics](./../monitoring-custom-metrics/), but use the steps above to verify that the metrics are collected.

### Cleanup

Run the following command to remove the installation from the cluster:

1. Remove the stack by calling Helm:

    ```bash
    helm delete -n ${KYMA_EXAMPLE_NS} ${HELM_RELEASE_NAME}
    ```

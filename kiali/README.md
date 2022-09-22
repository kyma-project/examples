# Install custom Kiali in Kyma

## Overview

The Kyma Kiali component brings limited configuration options in contrast to the upstream [`kiali-operator`](https://github.com/kiali/helm-charts/tree/master/kiali-operator) chart. Modifications might be reset at next upgrade cycle.

An alternative can be a parallel installation of the upstream chart offering all customization options. This tutorial outlines how to achieve such installation in co-existent to the Kyma stack.

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

1. Export the Helm release name which you want to use. It can be any name, be aware that all resources in the cluster will be prefixed with that name. Replace the `{release-name}` placeholder in the following command and run it:
    ```bash
    export HELM_RELEASE_NAME="{release-name}"
    ```

1. Update your helm installation with the required helm repository:

    ```bash
    helm repo add kiali https://kiali.org/helm-charts
    helm repo update
    ```

### Install the kiali-operator

Kiali recommends to install Kiali via the kiali-operator always. So this tutorial will do this as well by using the kiali-operator helm chart.

1. Run the Helm upgrade command which will install the chart if not present yet.
    ```bash
    helm upgrade --install --create-namespace -n ${KYMA_EXAMPLE_NS} ${HELM_RELEASE_NAME} kiali/kiali-operator --set cr.spec.auth.strategy=anonymous -f https://raw.githubusercontent.com/kyma-project/examples/main/kiali/values.yaml
    ```

Hereby, use the [values.yaml](./values.yaml) provided with this tutorial which contains customized settings deviating from the default settings, or create your own one.

### Verify the installation

1. You should see the kiali-operator and kiali pod coming up in the namespace. Assure that all pods are ending in a "Running" state.
1. Browse the Kiali dashboard. Following command will expose the dashboard on `http://localhost:9090`
   ```bash
   kubectl -n ${KYMA_EXAMPLE_NS} port-forward svc/kiali-server 20001
   ```

### Deploy a custom workload and call it

1. Follow the tutorial [orders-service](./../orders-service/) and see the service communication visualized in Kiali.

## Advanced Topics

### Integrate Jaeger

If you use [Jaeger](https://www.jaegertracing.io/) for distributed tracing, Kiali can utilize your Jaeger instance to [provide traces](https://kiali.io/docs/features/tracing/). Follow the [Jaeger configuration](https://kiali.io/docs/configuration/p8s-jaeger-grafana/jaeger/) guide for integration.

### Integrate Grafana

Kiali can provide links to Istio dashboards in Grafana. The [Grafana configuration](https://kiali.io/docs/configuration/p8s-jaeger-grafana/grafana/) page of the Kiali documentation describes how to integrate Grafana.

### Expose Kiali

Kiali supports different authentication strategies. When exposing the Kiali server, we recommend to use an external OpenID Connect compatible identity provider.

1. Set up Kiali with the [OpenID Connect strategy](https://kiali.io/docs/configuration/authentication/openid/).
1. Follow the [Expose a workload](https://kyma-project.io/docs/kyma/latest/03-tutorials/00-api-exposure/apix-03-expose-workload-apigateway/) Kyma tutorial to expose Kiali.

## Cleanup

Run the following commands to completely remove the example and all its resources from the cluster:

1. Remove the stack by calling helm:

    ```bash
    helm delete -n ${KYMA_EXAMPLE_NS} ${HELM_RELEASE_NAME}
    ```

# Install custom Kiali in Kyma

## Overview

The Kyma Kiali component brings limited configuration options in contrast to the upstream [`kiali-operator`](https://github.com/kiali/helm-charts/tree/master/kiali-operator) chart. Modifications might be reset at next upgrade cycle.

An alternative can be a parallel installation of the upstream chart offering all customization options. The following instructions outline how to achieve such installation in co-existence to the Kyma stack.

## Prerequisites

- Kyma as the target deployment environment
- kubectl > 1.22.x
- Helm 3.x

## Installation

### Preparation

1. Export your Namespace as a variable. Replace the `{namespace}` placeholder in the following command and run it:

    ```bash
    export KYMA_KIALI_NS="{namespace}"
    ```

1. Export the Helm release name that you want to use. The release name must be unique for the chosen Namespace. Be aware that all resources in the cluster will be prefixed with that name. Replace the `{release-name}` placeholder in the following command and run it:
    ```bash
    export HELM_RELEASE_NAME="{release-name}"
    ```

1. Update your Helm installation with the required Helm repository:

    ```bash
    helm repo add kiali https://kiali.org/helm-charts
    helm repo update
    ```

### Install the kiali-operator

> **NOTE:** Kiali recommends to install Kiali always with the Kiali operator; that's why the following step uses the Kiali operator Helm chart.

1. Run the Helm upgrade command, which installs the chart if not present yet.
    ```bash
    helm upgrade --install --create-namespace -n ${KYMA_KIALI_NS} ${HELM_RELEASE_NAME} kiali/kiali-operator --set cr.spec.auth.strategy=anonymous -f https://raw.githubusercontent.com/kyma-project/examples/main/kiali/values.yaml
    ```

You can either use the `[values.yaml](./values.yaml)` provided in this `kiali` folder, which contains customized settings deviating from the default settings, or create your own `values.yaml` file.

### Verify the installation

1. You should see the `kiali-operator` and `kiali` Pod coming up in the Namespace. All Pods must eventually be in a "Running" state.
1. Browse the Kiali dashboard. The following command exposes the dashboard on `http://localhost:20001`:
   ```bash
   kubectl -n ${KYMA_KIALI_NS} port-forward svc/kiali-server 20001
   ```

### Deploy a custom workload and invoke

1. To see the service communication visualized in Kiali, follow the instructions in [orders-service](./../orders-service/).

## Advanced Topics

### Integrate Jaeger

If you use [Jaeger](https://www.jaegertracing.io/) for distributed tracing, Kiali can use your Jaeger instance to [provide traces](https://kiali.io/docs/features/tracing/).

For integration instructions, read [Kiali: Jaeger configuration](https://kiali.io/docs/configuration/p8s-jaeger-grafana/jaeger/).

### Integrate Grafana

Kiali can provide links to Istio dashboards in Grafana. 

For integration instructions, read [Kiali: Grafana configuration](https://kiali.io/docs/configuration/p8s-jaeger-grafana/grafana/).

### Expose Kiali

Kiali supports different authentication strategies. When exposing the Kiali server, we recommend to use an external identity provider compatible with OpenID Connect (OIDC).

1. Set up Kiali with the [OpenID Connect strategy](https://kiali.io/docs/configuration/authentication/openid/).
1. To expose Kiali, follow the [Expose a workload](https://kyma-project.io/docs/kyma/latest/03-tutorials/00-api-exposure/apix-03-expose-workload-apigateway/) instructions.

## Cleanup

When you're done, you can remove the example and all its resources from the cluster.

1. Remove the stack by calling Helm:

    ```bash
    helm delete -n ${KYMA_KIALI_NS} ${HELM_RELEASE_NAME}
    ```

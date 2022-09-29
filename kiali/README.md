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
    helm upgrade --install --create-namespace -n ${KYMA_KIALI_NS} ${HELM_RELEASE_NAME} kiali/kiali-operator -f https://raw.githubusercontent.com/kyma-project/examples/main/kiali/values.yaml
    ```

You can either use the `[values.yaml](./values.yaml)` provided in this `kiali` folder, which contains customized settings deviating from the default settings, or create your own `values.yaml` file.

### Verify the installation

You should see the `kiali-operator` and `kiali-server` Pod coming up in the Namespace. All Pods must eventually be in a "Running" state.

### Access Kiali

To access Kiali, you can either use kubectl port-forwarding or expose it over the Kyma ingress gateway.

* To access Kiali using port-forwarding run:
  ```bash
  kubectl -n ${KYMA_KIALI_NS} port-forward svc/kiali-server 20001
  ```

  Open Kiali in your browser under [http://localhost:20001](http://localhost:20001).

* To expose Kiali over the Kyma ingress gateway, create an APIRule:
  ```bash
  kubectl -n ${KYMA_KIALI_NS} apply -f https://raw.githubusercontent.com/kyma-project/examples/main/kiali/apirule.yaml
  ```
  Get the public URL of your Kiali server:
  ```bash
  kubectl -n ${KYMA_KIALI_NS} get virtualservices.networking.istio.io -ojsonpath='{.items[*].spec.hosts[*]}'
  ```

### Deploy a custom workload and invoke

To see the service communication visualized in Kiali, follow the instructions in [orders-service](./../orders-service/).

## Advanced Topics

### Integrate Jaeger

If you use [Jaeger](https://www.jaegertracing.io/) for distributed tracing, Kiali can use your Jaeger instance to [provide traces](https://kiali.io/docs/features/tracing/).

For integration instructions, read [Kiali: Jaeger configuration](https://kiali.io/docs/configuration/p8s-jaeger-grafana/jaeger/).

To integrate with the Kyma tracing component run:

```bash
export JAEGER_URL=`kubectl -n kyma-system get virtualservices.networking.istio.io tracing -ojsonpath='{.spec.hosts[0]}'`
helm upgrade --install --create-namespace -n ${KYMA_KIALI_NS} ${HELM_RELEASE_NAME} kiali/kiali-operator --set cr.spec.external_services.tracing.enabled=true --set cr.spec.external_services.tracing.url=https://${JAEGER_URL} --set cr.spec.external_services.tracing.in_cluster_url=http://tracing-jaeger-query.kyma-system:16686 --set cr.spec.external_services.tracing.use_grpc=false -f https://raw.githubusercontent.com/kyma-project/examples/main/kiali/values.yaml
```

Ensure that access to the Jaeger query service from the Kiali pod is allowed by the active Istio Authorization Policies. For the Kyma tracing component, use the override values in [tracing-values.yaml](tracing-values.yaml) and adapt the namespace for the used principal.

### Integrate Grafana

Kiali can provide links to Istio dashboards in Grafana.

For integration instructions, read [Kiali: Grafana configuration](https://kiali.io/docs/configuration/p8s-jaeger-grafana/grafana/).

### Authentication

Kiali supports different authentication strategies. The default authentication strategy uses a Kubernetes Service Account Token. If you use a kubeconfig file with a static token, you can use this token to authenticate. Depending on the chosen way to access Kiali, different authentication strategies might be suitable. To learn more about Kiali authentication strategies, read [Kiali: Authentication Strategies](https://kiali.io/docs/configuration/authentication/).

* For Kiali access by port-forwarding, no additional authentication is required and the [anonymous strategy](https://kiali.io/docs/configuration/authentication/anonymous/) can be activated:
  ```bash
  helm upgrade --install --create-namespace -n ${KYMA_KIALI_NS} ${HELM_RELEASE_NAME} kiali/kiali-operator --set cr.spec.auth.strategy=anonymous -f https://raw.githubusercontent.com/kyma-project/examples/main/kiali/values.yaml
  ```
* When exposing the Kiali server over the ingress gateway, we recommend to use an external identity provider compatible with OpenID Connect (OIDC). Find the required settings at [Kiali: OpenID Connect strategy](https://kiali.io/docs/configuration/authentication/openid/).

## Cleanup

When you're done, you can remove the example and all its resources from the cluster.

1. Remove the stack by calling Helm:

    ```bash
    helm delete -n ${KYMA_KIALI_NS} ${HELM_RELEASE_NAME}
    kubectl -n ${KYMA_KIALI_NS} delete -f https://raw.githubusercontent.com/kyma-project/examples/main/kiali/apirule.yaml
    ```

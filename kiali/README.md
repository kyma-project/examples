# Install custom Kiali in Kyma

> [!WARNING]
> This guide is archived and not maintained anymore. Current integration guides around observability are located [here](https://kyma-project.io/#/telemetry-manager/user/README?id=integration-guides).

## Overview

The following instructions outline how to install [`Kiali`](https://github.com/kiali/helm-charts/tree/master/kiali-operator) in Kyma.

## Prerequisites

- Kyma as the target deployment runtime
- A [Prometheus instance preserving Istio metrics](./../prometheus/) deployed to the runtime.
- kubectl version 1.26.x or higher
- Helm 3.x

## Installation

### Preparation

1. Export your Namespace as a variable. Replace the `{namespace}` placeholder in the following command and run it:

    ```bash
    export K8S_NAMESPACE="{namespace}"
    ```

1. Export the Helm release name that you want to use. The release name must be unique for the chosen Namespace. Be aware that all resources in the cluster will be prefixed with that name. Run the following command:
    ```bash
    export HELM_KIALI_RELEASE="kiali"
    ```

1. Update your Helm installation with the required Helm repository:

    ```bash
    helm repo add kiali https://kiali.org/helm-charts
    helm repo update
    ```

### Install the kiali-operator

> **NOTE:** Kiali recommends to install Kiali always with the Kiali operator; that's why the following step uses the Kiali operator Helm chart.

Run the Helm upgrade command, which installs the chart if not present yet.
```bash
export PROM_SERVICE_NAME=$(kubectl -n ${K8S_NAMESPACE} get service -l app=kube-prometheus-stack-prometheus -ojsonpath='{.items[*].metadata.name}')
helm upgrade --install --create-namespace -n $K8S_NAMESPACE $HELM_KIALI_RELEASE kiali/kiali-operator -f https://raw.githubusercontent.com/kyma-project/examples/main/kiali/values.yaml --set cr.spec.external_services.prometheus.url=http://$PROM_SERVICE_NAME.$K8S_NAMESPACE:9090
```

You can either use the [`values.yaml`](./values.yaml) provided in this `kiali` folder, which contains customized settings deviating from the default settings, or create your own `values.yaml` file.

### Verify the installation

Check that the `kiali-operator` and `kiali-server` Pods have been created in the Namespace and are in the `Running` state:
```bash
kubectl -n $K8S_NAMESPACE rollout status deploy $HELM_KIALI_RELEASE-kiali-operator && kubectl -n $K8S_NAMESPACE rollout status deploy kiali-server
```
### Access Kiali

To access Kiali, either use kubectl port forwarding, or expose it using the Kyma Ingress Gateway.

* To access Kiali using port forwarding, run:
  ```bash
  kubectl -n $K8S_NAMESPACE port-forward svc/kiali-server 20001
  ```

  Open Kiali in your browser under `http://localhost:20001` and log in with a [Kubernetes service account token](https://kubernetes.io/docs/reference/access-authn-authz/authentication/#service-account-tokens), for instance, from your kubeconfig file.

* To expose Kiali using the Kyma API Gateway, create an APIRule:
  ```bash
  kubectl -n $K8S_NAMESPACE apply -f https://raw.githubusercontent.com/kyma-project/examples/main/kiali/apirule.yaml
  ```
  Get the public URL of your Kiali server:
  ```bash
  kubectl -n $K8S_NAMESPACE get vs -l apirule.gateway.kyma-project.io/v1beta1=kiali.$K8S_NAMESPACE -ojsonpath='{.items[*].spec.hosts[*]}'
  ```

### Deploy a custom workload and invoke

To see the service communication visualized in Kiali, follow the instructions in [`orders-service`](./../orders-service/).

## Advanced Topics

### Integrate Jaeger

If you use [Jaeger](https://www.jaegertracing.io/) for distributed tracing, Kiali can use your Jaeger instance to [provide traces](https://kiali.io/docs/features/tracing/).

For integration instructions, read [Kiali: Jaeger configuration](https://kiali.io/docs/configuration/p8s-jaeger-grafana/tracing/jaeger/).

### Integrate Grafana

Kiali can provide links to Istio dashboards in Grafana.

For integration instructions, read [Kiali: Grafana configuration](https://kiali.io/docs/configuration/p8s-jaeger-grafana/grafana/).

### Authentication

Kiali supports different authentication strategies. The default authentication strategy uses a Kubernetes Service Account Token. If you use a kubeconfig file with a static token, you can use this token to authenticate. Depending on your preferred way to access Kiali, different authentication strategies might be suitable. To learn more about Kiali authentication strategies, read [Kiali: Authentication Strategies](https://kiali.io/docs/configuration/authentication/).

* For Kiali access by port forwarding, you need no additional authentication, and you can activate the [anonymous strategy](https://kiali.io/docs/configuration/authentication/anonymous/):
  ```bash
  helm upgrade --install --create-namespace -n $K8S_NAMESPACE $HELM_KIALI_RELEASE kiali/kiali-operator --set cr.spec.auth.strategy=anonymous -f https://raw.githubusercontent.com/kyma-project/examples/main/kiali/values.yaml
  ```
* When exposing the Kiali server over the ingress gateway, we recommend to use an external identity provider compatible with OpenID Connect (OIDC). Find the required settings at [Kiali: OpenID Connect strategy](https://kiali.io/docs/configuration/authentication/openid/).

## Cleanup

When you're done, you can remove the example and all its resources from the cluster.

1. Remove the stack by calling Helm:

    ```bash
    helm delete -n $K8S_NAMESPACE $HELM_KIALI_RELEASE
    kubectl -n $K8S_NAMESPACE delete -f https://raw.githubusercontent.com/kyma-project/examples/main/kiali/apirule.yaml
    ```

2. If you created the `$K8S_NAMESPACE` Namespace specifically for this tutorial, remove the Namespace:
    ```bash
    kubectl delete namespace $K8S_NAMESPACE
    ``` 

# Install custom Jaeger in Kyma

## Overview

The Kyma tracing component brings limited configuration options in contrast to the upstream [`jaeger`](https://github.com/jaegertracing/helm-charts/tree/main/charts/jaeger) chart. Modifications might be reset at next upgrade cycle.

An alternative can be a parallel installation of the upstream chart offering all customization options. The following instructions outline how to achieve such installation in co-existence to the Kyma stack.

## Prerequisites

- Kyma as the target deployment environment
- kubectl > 1.22.x
- Helm 3.x

## Installation

### Preparation

1. Export your Namespace as a variable. Replace the `{namespace}` placeholder in the following command and run it:

    ```bash
    export KYMA_NS="{namespace}"
    ```

1. Export the Helm release name that you want to use. The release name must be unique for the chosen Namespace. Be aware that all resources in the cluster will be prefixed with that name. Replace the `{release-name}` placeholder in the following command and run it:
    ```bash
    export HELM_JAEGER_RELEASE="{release-name}"
    ```

1. Update your Helm installation with the required Helm repository:

    ```bash
    helm repo add jaegertracing https://jaegertracing.github.io/helm-charts
    helm repo update
    ```

### Install Jaeger

> **NOTE:** Jaeger recommends to install Jaeger with the [Jaeger operator](https://github.com/jaegertracing/helm-charts/tree/main/charts/jaeger-operator); As the operator requires a cert-manager to be installed, the instructions here will use a plain Jaeger installation. Anyway the installation described is not meant to be used for productive setups.

Run the Helm upgrade command, which installs the chart if not present yet.
```bash
helm upgrade --install --create-namespace -n $KYMA_NS $HELM_JAEGER_RELEASE jaegertracing/jaeger -f https://raw.githubusercontent.com/kyma-project/examples/main/jaeger/values.yaml
```

You can either use the [`values.yaml`](./values.yaml) provided in this `jaeger` folder, which contains customized settings deviating from the default settings, or create your own `values.yaml` file.

### Verify the installation

Check that the `jaeger` Pod has been created in the Namespace and is in the `Running` state:
```bash
kubectl -n $KYMA_NS rollout status deploy $HELM_JAEGER_RELEASE-jaeger
```
### Access Jaeger

To access Jaeger, either use kubectl port forwarding, or expose it using the Kyma Ingress Gateway.

* To access Kiali using port forwarding, run:
  ```bash
  kubectl -n $KYMA_NS port-forward svc/$HELM_JAEGER_RELEASE-jaeger-query 16686
  ```

  Open the Jaeger UI in your browser under [http://localhost:16686](http://localhost:16686).

* To expose Jaeger using the Kyma API Gateway, create an APIRule:
  ```bash
  cat <<EOF | kubectl -n $KYMA_NS apply -f -
  apiVersion: gateway.kyma-project.io/v1beta1
  kind: APIRule
  metadata:
    name: jaeger
  spec:
    host: jaeger-ui
    service:
      name: $HELM_JAEGER_RELEASE-jaeger-query
      port: 16686
    gateway: kyma-system/kyma-gateway
    rules:
      - path: /.*
        methods: ["GET", "POST"]
        accessStrategies:
          - handler: noop
        mutators:
          - handler: noop
  EOF
  ```

  Get the public URL of your Jaeger instance:
  ```bash
  kubectl -n $KYMA_NS get vs -l apirule.gateway.kyma-project.io/v1beta1=jaeger.$KYMA_NS -ojsonpath='{.items[*].spec.hosts[*]}'
  ```

### Deploy a workload and activate Kymas TracePipeline feature

To see distributed traces visualized in Jaeger, follow the instructions in [`trace-demo`](./../trace-demo/).

## Advanced Topics

### Integrate with Grafana

Jaeger can be provided as datasource integrated to Grafana, for example as part of a Grafana installation as described in the [Prometheus tutorial](./../prometheus/README.md). For that, deploy a Grafana datasource like this:

```bash
cat <<EOF | kubectl -n $KYMA_NS apply -f -
apiVersion: v1
kind: ConfigMap
metadata:
  name: jaeger-grafana-datasource
  labels:
    grafana_datasource: "1"
data:
    jaeger-datasource.yaml: |-
      apiVersion: 1
      datasources:
      - name: Jaeger
        type: jaeger
        access: proxy
        url: http://$HELM_JAEGER_RELEASE-jaeger-query.$KYMA_NS:16686
        editable: true
EOF
```
You will need to restart the Grafana instance. Afterwards, the Jaeger datasource will be available in the `Explore` view.

### Authentication

Jaeger does not provide authentication mechanisms by itself. To secure Jaeger please follow the instructions provided at the [Jaeger documentation
](https://www.jaegertracing.io/docs/latest/security/#browser-to-ui)

## Cleanup

When you're done, you can remove the example and all its resources from the cluster.

1. Remove the stack by calling Helm:

    ```bash
    helm delete -n $KYMA_NS $HELM_JAEGER_RELEASE
    kubectl -n $KYMA_NS delete -f https://raw.githubusercontent.com/kyma-project/examples/main/jaeger/apirule.yaml
    ```

2. (Optional) If you created the `$KYMA_NS` Namespace specifically for this tutorial, remove the Namespace:
    ```bash
    kubectl delete namespace $KYMA_NS
    ``` 

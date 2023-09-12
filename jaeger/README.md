# Install custom Jaeger in Kyma

## Overview

The following instructions outline how to use [`Jaeger`](https://github.com/jaegertracing/helm-charts/tree/main/charts/jaeger) as a tracing backend with Kyma's [TracePipeline](https://kyma-project.io/#/telemetry-manager/user/03-traces).

## Prerequisites

- Kyma version 2.10 or higher as the target deployment environment
- kubectl version 1.22.x or higher
- Helm 3.x

## Installation

### Preparation

1. Export your Namespace as a variable. Replace the `{namespace}` placeholder in the following command and run it:

    ```bash
    export KYMA_NS="{namespace}"
    ```

1. Export the Helm release name that you want to use. The release name must be unique for the chosen Namespace. Be aware that all resources in the cluster will be prefixed with that name. Run the following command:
    ```bash
    export HELM_JAEGER_RELEASE="jaeger"
    ```

1. Update your Helm installation with the required Helm repository:

    ```bash
    helm repo add jaegertracing https://jaegertracing.github.io/helm-charts
    helm repo update
    ```

### Install Jaeger

> **NOTE:** It is officially recommended to install Jaeger with the [Jaeger operator](https://github.com/jaegertracing/helm-charts/tree/main/charts/jaeger-operator). Because the operator requires a cert-manager to be installed, the following instructions use a plain Jaeger installation. However, the described installation is not meant to be used for production setups.

Run the Helm upgrade command, which installs the chart if not present yet.
```bash
helm upgrade --install --create-namespace -n $KYMA_NS $HELM_JAEGER_RELEASE jaegertracing/jaeger -f https://raw.githubusercontent.com/kyma-project/examples/main/jaeger/values.yaml
```

You can either use the [`values.yaml`](./values.yaml) provided in this `jaeger` folder, which contains customized settings deviating from the default settings, or create your own `values.yaml` file.

### Verify the installation

Check if the `jaeger` Pod was successfully created in the Namespace and is in the `Running` state:
```bash
kubectl -n $KYMA_NS rollout status deploy $HELM_JAEGER_RELEASE
```

### Activate a TracePipeline

To configure the Kyma trace collector with the deployed Jaeger instance as the backend. To create a new [TracePipeline](https://kyma-project.io/#/telemetry-manager/user/03-traces), 
execute the following command:
   ```bash
   cat <<EOF | kubectl -n $KYMA_NS apply -f -
   apiVersion: telemetry.kyma-project.io/v1alpha1
   kind: TracePipeline
   metadata:
     name: jaeger
   spec:
     output:
       otlp:
         protocol: http
         endpoint:
           value: http://$HELM_JAEGER_RELEASE-collector.$KYMA_NS.svc.cluster.local:4318
   EOF
   ```
  
### Activate Istio Tracing

To [enable Istio](https://kyma-project.io/#/telemetry-manager/user/03-traces?id=step-2-enable-istio-tracing) to report span data, apply an Istio telemetry resource and set the sampling rate to 100%. This approach is not recommended for production.

```bash
cat <<EOF | kubectl apply -f -
apiVersion: telemetry.istio.io/v1alpha1
kind: Telemetry
metadata:
  name: tracing-default
  namespace: istio-system
spec:
  tracing:
  - providers:
    - name: "kyma-traces"
    randomSamplingPercentage: 100.00
EOF
```

### Access Jaeger

To access Jaeger using port forwarding, run:
```bash
kubectl -n $KYMA_NS port-forward svc/$HELM_JAEGER_RELEASE-query 16686
```

Open the Jaeger UI in your browser under `http://localhost:16686`.

### Deploy a workload and activate Kyma's TracePipeline feature

To see distributed traces visualized in Jaeger, follow the instructions in [`trace-demo`](./../trace-demo/).

## Advanced Topics

### Integrate with Grafana

Jaeger can be provided as a data source integrated into Grafana. For example, it can be part of a Grafana installation as described in the [Prometheus tutorial](./../prometheus/README.md). To have a Jaeger data source as part of the Grafana installation, deploy a Grafana data source in the following way:

```bash
cat <<EOF | kubectl -n $KYMA_NS apply -f -
apiVersion: v1
kind: ConfigMap
metadata:
  name: jaeger-grafana-datasource
  labels:
    grafana_datasource: "1"
data:
    jaeger-grafana-datasource.yaml: |-
      apiVersion: 1
      datasources:
      - name: Jaeger-Tracing
        type: jaeger
        access: proxy
        url: http://$HELM_JAEGER_RELEASE-query.$KYMA_NS:16686
        editable: true
EOF
```
Restart the Grafana instance. Afterwards, the Jaeger data source will be available in the `Explore` view.

### Authentication

Jaeger does not provide authentication mechanisms by itself. To secure Jaeger, follow the instructions provided in the [Jaeger documentation
](https://www.jaegertracing.io/docs/latest/security/#browser-to-ui).

### Exposure
>**CAUTION**: The following approach exposes the Jaeger instance as it is without providing any ways of authentication.

To expose Jaeger using Kyma API Gateway, create the following APIRule:
```bash
cat <<EOF | kubectl -n $KYMA_NS apply -f -
apiVersion: gateway.kyma-project.io/v1beta1
kind: APIRule
metadata:
  name: jaeger
spec:
  host: jaeger-ui
  service:
    name: $HELM_JAEGER_RELEASE-query
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

## Cleanup

When you're done, remove the example and all its resources from the cluster.

1. Remove the stack by calling Helm:

    ```bash
    helm delete -n $KYMA_NS $HELM_JAEGER_RELEASE
    ```

2. If you created the `$KYMA_NS` Namespace specifically for this tutorial, remove the Namespace:
    ```bash
    kubectl delete namespace $KYMA_NS
    ``` 

# Install an OTLP-based metrics collector 

## Overview

The following instructions demonstrate how to install [OpenTelemetry Collector](https://github.com/open-telemetry/opentelemetry-collector)s on a Kyma cluster using the official [Helm chart](https://github.com/open-telemetry/opentelemetry-helm-charts/tree/main/charts/opentelemetry-collector) with the goal to collect and ship workload metrics to an OTLP endpoint. For a fully Prometheus-based approach, have a look at the [Prometheus](./../prometheus/README.md) tutorial instead.

The setup brings an OpenTelemetry Collector Deployment acting as gateway, to which all the cluster-wide metrics should be ingested. Then, the gateway enriches the metrics with missing resource attributes and ships them to a target backend.

Furthermore, the setup brings an OpenTelemetry Collector DaemonSet acting as agent running on each node. Here, node-specific metrics relevant for your workload are determined. This is also where an annotation-based workload scraper is running, which scrapes all containers running on the related node.
## Architecture

[Architecture](./assets/overview.drawio.svg)

## Prerequisites

- Kyma Open Source 2.10.x or higher
- kubectl version 1.22.x or higher
- Helm 3.x

## Preparation

1. Export your Namespace as a variable. Replace the `{namespace}` placeholder in the following command and run it:

    ```bash
    export KYMA_NS="{namespace}"
    ```

1. If you haven't created a Namespace yet, do it now:
    ```bash
    kubectl create namespace $KYMA_NS
    ```

1. Update your Helm installation with the required Helm repository:

    ```bash
    helm repo add open-telemetry https://open-telemetry.github.io/opentelemetry-helm-charts
    helm repo update
    ```

## Deploy the gateway

1. Deploy an Otel Collector using the upstream Helm chart with a prepared [values.yaml](./metrics-gateway-values.yaml) file.

   The values file defines a pipeline for receiving OTLP metrics, enriches them with resource attributes that fulfil the Kubernetes semantic conventions, and then exports them to a custom OTLP backend.
   
   The previous instructions don't provide any backend, so you must set the backend configuration. 
   If you don't want to use a Secret, use the following command, and adjust the placeholders `myEndpoint` and `myToken` to your needs:

   ```bash
   helm upgrade metrics-gateway open-telemetry/opentelemetry-collector --version 0.62.2 --install --namespace $KYMA_NS \
     -f https://raw.githubusercontent.com/kyma-project/examples/main/metrics-otlp/metrics-gateway-values.yaml \
     --set config.exporters.otlp.endpoint="{myEndpoint}" \
     --set config.exporters.otlp.headers.Authorization="Bearer {myToken}"
   ```

   
   > **TIP:** It's recommended that you provide tokens using a Secret. To achieve that, you can mount the relevant attributes of your secret via the `extraEnvs` parameter and use placeholders for referencing the actual environment values. Take a look at the provided sample [secret-values.yaml](./secret-values.yaml) file, adjust it to your Secret, and run:
   ```bash
   helm upgrade metrics-gateway open-telemetry/opentelemetry-collector --version 0.62.2 --install --namespace $KYMA_NS \
     -f https://raw.githubusercontent.com/kyma-project/examples/main/metrics-otlp/metrics-gateway-values.yaml \
     -f secret-values.yaml
   ```

1. Verify the deployment.
   Check that the related Pod has been created in the Namespace and is in the `Running` state:
   ```bash
   kubectl -n $KYMA_NS rollout status deploy metrics-gateway
   ```

## Install the agent

1. Deploy the agent.

   The gateway deployed in the previous step can be used to _push_ metrics via OTLP. If you prefer to use the Prometheus _pull_ approach, you need a dedicated instance to scrape your workload. Furthermore, to retrieve node-specific metrics, you need an agent as a DaemonSet.
   
   Deploy an Otel Collector using the upstream Helm chart with a prepared [values](./metrics-agent-values.yaml). The values file defines a metrics pipeline for scraping workload by annotation and pushing them to the gateway. It also defines a second metrics pipeline determining node-specific metrics for your workload from the nodes kubelet and the nodes filesystem itself.

   ```bash
   helm upgrade metrics-agent open-telemetry/opentelemetry-collector --version 0.62.2 --install --namespace $KYMA_NS \
     -f https://raw.githubusercontent.com/kyma-project/examples/main/metrics-otlp/metrics-agent-values.yaml \
     --set config.exporters.otlp.endpoint=metrics-gateway.$KYMA_NS:4317
   ```

1. Verify the deployment.
   Check that the related Pod has been created in the Namespace and is in the `Running` state:
   ```bash
   kubectl -n $KYMA_NS rollout status daemonset metrics-agent
   ```

## Result

Now, you have a setup in place as outlined in the [architecture diagram](#architecture), with a gateway and an agent. By default, the node-specific metrics are automatically ingested by the agent, and pushed via the gateway to the configured backend. Furthermore, metrics for the gateway and agent instances are exported.


## Usage

To add metrics ingestion for your custom workload, you have the following options:

### Prometheus pull-based, with annotations

This approach assumes that you instrumented your application using a library like the [prometheus client library](https://prometheus.io/docs/instrumenting/clientlibs/), having a port in your workload exposed serving a typical Prometheus metrics endpoint.

The agent is configured with a generic scrape configuration, which uses annotations to specify the endpoints to scrape in the cluster. 
Having the annotations in place is everything you need for metrics ingestion to start automatically.

Put the following annotations either to a service that resolves your metrics port, or directly to the pod:

```yaml
prometheus.io/scrape: "true"   # mandatory to enable automatic scraping
prometheus.io/scheme: https    # optional, default is "http" if no Istio sidecar is used. When using a sidecar (Pod has label `security.istio.io/tlsMode=istio`), the default is "https". Use "https" to scrape workloads using Istio client certificates.
prometheus.io/port: "1234"     # optional, configure the port under which the metrics are exposed
prometheus.io/path: /myMetrics # optional, configure the path under which the metrics are exposed
```


> **NOTE:** The agent can scrape endpoints even if the workload uses Istio and accepts only mTLS communication. Because the agent itself should not be part of the Service Mesh in order to observe the Service Mesh, the agent uses a sidecar but has no traffic interception enabled. Instead, it mounts the client certificate and uses the certificate natively for the communication. That is a [recommended approach](https://istio.io/latest/docs/ops/integrations/prometheus/#tls-settings) by Istio.

To try it out, you can install the demo app taken from the tutorial [Expose Custom Metrics in Kyma](https://github.com/kyma-project/examples/tree/main/prometheus/monitoring-custom-metrics) and annotate the workload with the annotations mentioned above.

```bash
kubectl apply -n $KYMA_NS -f https://raw.githubusercontent.com/kyma-project/examples/main/prometheus/monitoring-custom-metrics/deployment/deployment.yaml
kubectl -n $KYMA_NS annotate service sample-metrics prometheus.io/scrape=true
kubectl -n $KYMA_NS annotate service sample-metrics prometheus.io/port=8080
```

The workload exposes the metric `cpu_temperature_celsius` at port `8080`, which is then automatically ingested to your configured backend.

### OpenTelemetry push-based

This approach assumes that you instrumented your workload using the [OpenTelemetry SDK](https://opentelemetry.io/docs/instrumentation/). Here, you only need to tell your workload to which target it should push the metrics. By using the deployed gateway, the application configuration does not need to be aware of any target backend, it only needs to get the static endpoint of the gateway. To achieve this, configure the SDK-specific [environment variables](https://opentelemetry.io/docs/reference/specification/protocol/exporter/). Also, if not done yet as part of your instrumentation, it's strongly recommended that you configure a service name.

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name MyDeployment
spec:
  template:
    spec:
      containers:
      - env:
        - name: OTEL_EXPORTER_OTLP_METRICS_ENDPOINT
          value: http://metrics-gateway.$KYMA_NS:4317
        - name: OTEL_SERVICE_NAME
          valueFrom:
            fieldRef:
              apiVersion: v1
              fieldPath: metadata.labels['app.kubernetes.io/name']
```

## Cleanup

When you're done, you can remove the example and all its resources from the cluster by calling Helm:

```bash
helm delete -n $KYMA_NS metrics-gateway
helm delete -n $KYMA_NS metrics-agent
```

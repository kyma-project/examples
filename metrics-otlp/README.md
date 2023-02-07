# Install an OTLP based metrics collector 

## Overview

The following instructions demonstrates how to install [OpenTelemetry Collector](https://github.com/open-telemetry/otel-collector)s on a Kyma cluster using the official [Helm chart](https://github.com/open-telemetry/opentelemetry-helm-charts/tree/main/charts/opentelemetry-collector) with the goal to collect and ship workload metrics to an OTLP endpoint. For a fully Prometheus-based approach please have a look at the [Prometheus](./../prometheus/README.md) tutorial instead.

The setup will bring an OpenTelemetry Collector Deployment acting as gateway where all the cluster-wide metrics should be ingested to. The gateway will then care about enriching the metrics with missing resource attributes and dies the actual shipment to a target backend.

Furthermore it brings an OpenTelemetry Collector DaemonSet acting as agent running on each node. Here, node specific metrics relevant for your workload are getting determined and an annotation based workload scraper is running there, scraping all containers running on the related node.

[Architecture](./assets/overview.drawio.svg)

## Prerequisites

- Kyma Open Source >= 2.10.x
- kubectl version 1.22.x or higher
- Helm 3.x

## Preparation

1. Export your Namespace as a variable. Replace the `{namespace}` placeholder in the following command and run it:

    ```bash
    export KYMA_NS="{namespace}"
    ```

1. If you don't have created a Namespace yet, do it now and enable Istio injection by label:
    ```bash
    kubectl create namespace $KYMA_NS
    kubectl label namespace $KYMA_NS istio-injection=enabled
    ```

1. Update your Helm installation with the required Helm repository:

    ```bash
    helm repo add open-telemetry https://open-telemetry.github.io/opentelemetry-helm-charts
    helm repo update
    ```

## Install the gateway

1. Deploy the gateway

   Deploy an otel-collector using the upstream helm chart with prepared [values](./metrics-gateway-values.yaml). The values file is defining a pipeline for receiving OTLP metrics, is enriching them with resource attributes fullfilling the kubernetes semantic conventions, and then exports them to a custom OTLP backend.
   
   As this instructions are not providing any backend, the following command sets the backend configuration with placeholders `myEndpoint` and `myToken` which needs to get adjusted to your needs. Be aware of, that a token should be provided using a Secret which gets mounted via the `extraEnvs` parameter. As this cannot be passed via the command itself in an easy way, please consider using a dedicated additional values.yaml file for that.

   ```bash
   helm upgrade metrics-gateway open-telemetry/opentelemetry-collector --version 0.47.0 --install --namespace $KYMA_NS \
     -f https://raw.githubusercontent.com/kyma-project/examples/main/metrics-otlp/metrics-gateway-values.yaml \
     --set config.exporters.otlp.endpoint="{myEndpoint}" \
     --set config.exporters.otlp.headers.Authorization="Bearer {myToken}"
   ```

1. Verify the deployment
   Check that the related Pod has been created in the Namespace and is in the `Running` state:
   ```bash
   kubectl -n $KYMA_NS rollout status deploy metrics-gateway
   ```

## Install the agent

1. Deploy the agent

   The gateway deployed in the previous step can be used to push metrics via OTLP. If you prefer to use the prometheus pull-approach, a dedicated instance is needed for scraping your workload. Furthermore, an agent is required running as a daemonset in order to retrieve node specific metrics.
   
   The following command will deploy an otel-collector using the upstream helm chart with prepared [values](./metrics-agent-values.yaml). The values file is defining a metrics pipeline for scraping workload by annotation and pushing them to the gateway. It defines also a second metrics pipeline determining node-specific metrics for your workload from the nodes kubelet and the nodes filesystem itself.

   ```bash
   helm upgrade metrics-agent open-telemetry/opentelemetry-collector --version 0.47.0 --install --namespace $KYMA_NS \
     -f https://raw.githubusercontent.com/kyma-project/examples/main/metrics-otlp/metrics-agent-values.yaml \
     --set config.exporters.otlp.endpoint=metrics-gateway.$KYMA_NS:4317
   ```

1. Verify the deployment
   Check that the related Pod has been created in the Namespace and is in the `Running` state:
   ```bash
   kubectl -n $KYMA_NS rollout status daemonset metrics-agent
   ```

## Usage

After you have followed the instructions you should have a setup in place as outlined in the diagram on top of the tutorial, with a gateway and an agent. By default, the node specific metrics are getting ingested by the agent automatically and pushed via the gateway to the configured backend. Furthermore, metrics for the gateway and agent instances itself will be exported.

To add metrics ingestion for your custom workload you have two options:

### Prometheus pull-based on base of annotations

That approach assumes that you instrumented your application using a library like the [prometheus client library](https://prometheus.io/docs/instrumenting/clientlibs/), having a port in your workload exposed serving a typical prometheus metrics endpoint.

The agent is configured with a generic scrape configuration which will determine endpoints to scrape in the cluster via annotations. You can either put the annotations to a service resolving your metrics port or directly at the pod itself. The annotations are as followed:

```yaml
prometheus.io/scrape: "true"   # mandatory to enable automatic scraping
prometheus.io/scheme: https    # optional, default is "http" if no Istio sidecar is used. When using a sidecar (Pod has label `security.istio.io/tlsMode=istio`), the default is "https". Use "https" to scrape workloads using Istio client certificates.
prometheus.io/port: "1234"     # optional, configure the port under which the metrics are exposed
prometheus.io/path: /myMetrics # optional, configure the path under which the metrics are exposed
```

Having the annotations in place is everything you need, metrics ingestion should start automatically.

Be aware of, that the agent is capable of scraping endpoints even if the workload is using istio and accepts only mTLS communication. As the agent itself should not be part of the Service Mesh in order to observe the Service Mesh, the agent uses a sidecar but having no traffic interception enabled at all. Instead it mounts the client certificate and uses the certificate natively for the communication. That is a [recommended approach](https://istio.io/latest/docs/ops/integrations/prometheus/#tls-settings) by Istio itself.

### OpenTelemetry push-based

That approach assumes that you instrumented your workload using the [OpenTelemetry SDK](https://opentelemetry.io/docs/instrumentation/). Here, you only need to tell your workload to which target it should push the metrics. By using the deployed gateway, the application configuration does not need to be aware of any target backend, it only needs to get the static endpoint of the gateway. That can be achieved by configuring the SDK-specific [environment variables](https://opentelemetry.io/docs/reference/specification/protocol/exporter/). Also, if not done yet as part of your instrumentation, you strongly should configure a service name.

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

# Install OpenTelemetry Demo Application in Kyma

## Overview

The following instructions install the OpenTelemetry [demo application](https://github.com/open-telemetry/opentelemetry-demo) on a Kyma cluster using a provided [Helm chart](https://github.com/open-telemetry/opentelemetry-helm-charts/tree/main/charts/opentelemetry-demo). The demo application will be configured to push trace data via OTLP to the collector provided by Kyma, so that they will be collected together with the related Istio trace data.

**CAUTION:** The following instructions are under development and rely on a preliminary state of the [Configurable Tracing](https://github.com/kyma-project/kyma/issues/11231) feature. The goal is to have an early E2E scenario for that feature. For now, you must reconfigure the Istio and Telemetry installation, which is only possible using a custom Kyma OpenSource installation. In the final state, this reconfiguration will not be required anymore.

## Prerequisites

- Kyma Open Source version 2.10.x or higher
- kubectl version 1.22.x or higher
- Helm 3.x

### Preparation

1. Export your Namespace as a variable. Replace the `{namespace}` placeholder in the following command and run it:

    ```bash
    export KYMA_NS="{namespace}"
    ```
1. If you don't have created a Namespace yet, do it now:
    ```bash
    kubectl create namespace $KYMA_NS
    ```

1. To enable Istio injection in your Namespace, set the following label:
    ```bash
    kubectl label namespace $KYMA_NS istio-injection=enabled
    ```

1. Export the Helm release name that you want to use. The release name must be unique for the chosen Namespace. Be aware that all resources in the cluster will be prefixed with that name. Run the following command:
    ```bash
    export HELM_OTEL_RELEASE="otel"
    ```

1. Update your Helm installation with the required Helm repository:

    ```bash
    helm repo add open-telemetry https://open-telemetry.github.io/opentelemetry-helm-charts
    helm repo update
    ```

### Activate a Kyma TracePipeline

1. Provide a tracing backend and activate it.
   Install [Jaeger in-cluster](./../jaeger/) or provide a custom backend supporting the OTLP protocol.
2. Activate the Istio tracing feature.
To [enable Istio](https://kyma-project.io/docs/kyma/main/01-overview/main-areas/telemetry/telemetry-03-traces#step-2-enable-istio-tracing) to report span data, apply an Istio telemetry resource and set the sampling rate to 100%. This approach is not recommended for production.
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

### Install the application

Run the Helm upgrade command, which installs the chart if not present yet.
```bash
helm upgrade --version 0.15.4 --install --create-namespace -n $KYMA_NS $HELM_OTEL_RELEASE open-telemetry/opentelemetry-demo -f https://raw.githubusercontent.com/kyma-project/examples/main/trace-demo/values.yaml
```

You can either use the [`values.yaml`](./values.yaml) provided in this `trace-demo` folder, which contains customized settings deviating from the default settings, or create your own `values.yaml` file.

### Verify the application

To verify that the application is running properly, set up port forwarding and call the respective local hosts.

1. To verify the frontend:
   ```bash
   kubectl -n $KYMA_NS port-forward svc/$HELM_OTEL_RELEASE-frontend 8080
   ```
   ```bash
   open http://localhost:8080
   ````

2. To verify that traces arrive in the Jaeger backend:
   ```bash
   kubectl -n kyma-system port-forward svc/tracing-jaeger-query 16686
   ```
   ```bash
   open http://localhost:16686
   ````

3. Enable failures with the feature flag service:
   ```bash
   kubectl -n $KYMA_NS port-forward svc/$HELM_OTEL_RELEASE-featureflagservice 8081
   ```
   ```bash
   open http://localhost:8081
   ````

4. Generate load with the load generator:
   ```bash
   kubectl -n $KYMA_NS port-forward svc/$HELM_OTEL_RELEASE-loadgenerator 8089
   ```
   ```bash
   open http://localhost:8089
   ```

## Advanced

### Browser Instrumentation and missing root spans

The frontend application of the demo uses [browser instrumentation](https://github.com/open-telemetry/opentelemetry-demo/blob/main/docs/services/frontend.md#browser-instrumentation). Because of that, the root span of a trace is created externally to the cluster and will not be captured with the described setup. In Jaeger, you can see a warning at the first span indicating that there is a parent span that is not being captured.

In order to capture the spans reported by the browser, you need to expose the trace endpoint of the collector and configure the frontend to report the spans to that exposed endpoint.

To expose the frontend web app and the relevant trace endpoint, define an Istio Virtualservice as in the following example.
>**CAUTION** The example shows an insecure way of exposing the trace endpoint. Do not use it in production.
```bash
export CLUSTER_DOMAIN={my-domain}

cat <<EOF | kubectl -n $KYMA_NS apply -f -
apiVersion: networking.istio.io/v1beta1
kind: VirtualService
metadata:
  name: frontend
spec:
  gateways:
  - kyma-system/kyma-gateway
  hosts:
  - frontend.$CLUSTER_DOMAIN
  http:
  - route:
    - destination:
        host: telemetry-otlp-traces.kyma-system.svc.cluster.local
        port:
          number: 4318
    match:
    - uri:
        prefix: "/v1/traces"
    corsPolicy:
      allowOrigins:
      - prefix: https://frontend.$CLUSTER_DOMAIN
      allowHeaders:
      - Content-Type
      allowMethods:
      - POST
      allowCredentials: true
  - route:
    - destination:
        host: $HELM_OTEL_RELEASE-frontend.$KYMA_NS.svc.cluster.local
        port:
          number: 8080
EOF
```

Then, update the frontend configuration. To do that, run the following helm command and set the additional configuration:
```bash
helm upgrade --version 0.15.4 --install --create-namespace -n $KYMA_NS $HELM_OTEL_RELEASE open-telemetry/opentelemetry-demo \
-f https://raw.githubusercontent.com/kyma-project/examples/main/trace-demo/values.yaml \
--set 'components.frontend.envOverrides[1].name=PUBLIC_OTEL_EXPORTER_OTLP_TRACES_ENDPOINT' \
--set "components.frontend.envOverrides[1].value=https://frontend.$CLUSTER_DOMAIN/v1/traces"
```

Afterwards, the web application will be accessible in your browser at `https://frontend.$CLUSTER_DOMAIN`. Using the developer tools of the browser you should see that requests to the traces endpoint are successful. Now, every trace should have a proper root span recorded.
## Cleanup

When you're done, you can remove the example and all its resources from the cluster by calling Helm:

```bash
helm delete -n $KYMA_NS $HELM_OTEL_RELEASE
```

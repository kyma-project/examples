# Install OpenTelemetry Demo Application in Kyma

## Overview

The following instructions install the OpenTelemetry [demo application](https://github.com/open-telemetry/opentelemetry-demo) on a Kyma cluster using a provided [Helm chart](https://github.com/open-telemetry/opentelemetry-helm-charts/tree/main/charts/opentelemetry-demo). The demo application will be configured to push trace data via OTLP to the collector provided by Kyma, so that they will be collected together with the related Istio trace data.

**CAUTION:** The following instructions are under development and rely on a preliminary state of the [Configurable Tracing](https://github.com/kyma-project/kyma/issues/11231) feature. The goal is to have an early E2E scenario for that feature. For now, you must reconfigure the Istio and Telemetry installation, which is only possible using a custom Kyma OpenSource installation. In the final state, this reconfiguration will not be required anymore.

## Prerequisites

- Kyma OSS installed from main branch (`kyma deploy -s main`)
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

1. Export the Helm release name that you want to use. The release name must be unique for the chosen Namespace. Be aware that all resources in the cluster will be prefixed with that name. Replace the `{release-name}` placeholder in the following command and run it:
    ```bash
    export HELM_RELEASE="{release-name}"
    ```

1. Update your Helm installation with the required Helm repository:

    ```bash
    helm repo add open-telemetry https://open-telemetry.github.io/opentelemetry-helm-charts
    helm repo update
    ```

### Activate Kyma Tracing preview

1. (Temporary) In the Telemetry operator, enable the tracing feature:

    ```bash
    kubectl apply -f https://raw.githubusercontent.com/kyma-project/kyma/main/components/telemetry-operator/config/crd/bases/telemetry.kyma-project.io_tracepipelines.yaml
    kyma deploy -s main --component telemetry --value telemetry.operator.controllers.tracing.enabled=true
    ```

1. (Temporary) Activate Istio tracing based on w3c-tracecontext and OTLP:
    ```bash
    kyma deploy -s main --component istio -f https://raw.githubusercontent.com/kyma-project/examples/main/trace-demo/istio-values.yaml
    ```
 1. If necessary, restart the relevant workloads.

1. To enable the Jaeger backend, create a new TracePipeline:
   ```bash
   kubectl apply -f https://raw.githubusercontent.com/kyma-project/examples/main/trace-demo/tracepipeline.yaml
   ```

### Install the application

Run the Helm upgrade command, which installs the chart if not present yet.
```bash
helm upgrade --version 0.11.1 --install --create-namespace -n $KYMA_NS $HELM_RELEASE open-telemetry/opentelemetry-demo -f https://raw.githubusercontent.com/kyma-project/examples/main/trace-demo/values.yaml
```

You can either use the [`values.yaml`](./values.yaml) provided in this `trace-demo` folder, which contains customized settings deviating from the default settings, or create your own `values.yaml` file.

### Verify the application

To verify that the application is running properly, set up port forwarding and call the respective local hosts.

1. To verify the frontend:
   ```bash
   kubectl -n $KYMA_NS port-forward svc/$HELM_RELEASE-frontend 8080
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
   kubectl -n $KYMA_NS port-forward svc/$HELM_RELEASE-featureflagservice 8081
   ```
   ```bash
   open http://localhost:8081
   ````

4. Generate load with the load generator:
   ```bash
   kubectl -n $KYMA_NS port-forward svc/$HELM_RELEASE-loadgenerator 8089
   ```
   ```bash
   open http://localhost:8089
   ```

## Cleanup

When you're done, you can remove the example and all its resources from the cluster by calling Helm:

```bash
helm delete -n $KYMA_NS $HELM_RELEASE
```

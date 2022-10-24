# Install OpenTelemtry Demo Application in Kyma

## Overview

This instructions will install the OpenTelemetry [demo app](https://github.com/open-telemetry/opentelemetry-demo) on a Kyma cluster. The demo app is based on the related [helm chart](https://github.com/open-telemetry/opentelemetry-helm-charts/tree/main/charts/opentelemetry-demo)

**CAUTION:** This instructions are under development and are relying on a prelimary state of the [Configurable Tracing](https://github.com/kyma-project/kyma/issues/11231) story. The goal is to have early with an E2E scenario for that feature. The instructions for now require to reconfigure the Istio and telemetry installation which is only possible using a custom Kyma OpenSource installation. In the final state, this reconfiguration will not be required anymore.

## Prerequisites

- Kyma OSS installed from main branch (`kyma deploy -s main`)
- kubectl > 1.22.x
- Helm 3.x

### Preparation

1. Export your Namespace as a variable. Replace the `{namespace}` placeholder in the following command and run it:

    ```bash
    export KYMA_NS="{namespace}"
    ```
1. If you don't have it created yet, now is the time to do sog:
    ```bash
    kubectl create namespace $KYMA_NS
    ```

1. Assure that your namespace has istio-onjection enabled by having the proper label in place
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

### Activate Kyma Tracing Preview feature

1. (Temporary) Enable tracing feature in operator

    ```bash
    kubectl apply -f https://raw.githubusercontent.com/kyma-project/kyma/main/components/telemetry-operator/config/crd/bases/telemetry.kyma-project.io_tracepipelines.yaml
    kyma deploy -s main --component telemetry --value telemetry.operator.controllers.tracing.enabled=true
    ```

1. (Temporary) Activate Istio tracing based on w3c-tracecontext and OTLP
    ```bash
    kyma deploy -s main --component istio -f istio-values.yaml
    ```
    Potentially, you need to restart the relevant workloads

1. Enable Jaeger backend via new tracing feature
   ```bash
   kubectl apply -f tracepipeline.yaml
   ```

### Install the application

Run the Helm upgrade command, which installs the chart if not present yet.
```bash
helm upgrade --version 0.9.6 --install --create-namespace -n $KYMA_NS $HELM_RELEASE open-telemetry/opentelemetry-demo -f values.yaml
```

You can either use the [`values.yaml`](./values.yaml) provided in this `trace-demo` folder, which contains customized settings deviating from the default settings, or create your own `values.yaml` file.

### Verify the application

1. Browse in the frontend by:
   ```bash
   kubectl -n $KYMA_NS port-forward svc/$HELM_RELEASE-frontend 8080
   ````
   and calling `http://localhost:8080`

1. Verify that traces arrive in the Jaeger backend:
   ```bash
   kubectl -n kyma-system port-forward svc/tracing-jaeger-query 16686
   ````
   and calling `http://localhost:16686`

1. Enable failures via the feature flag service:
   ```bash
   kubectl -n $KYMA_NS port-forward svc/$HELM_RELEASE-featureflagservice 8081
   ````
   and calling `http://localhost:8081`

1. Generate load via the load generator:
   ```bash
   kubectl -n $KYMA_NS port-forward svc/$HELM_RELEASE-loadgenerator 8089
   ````
   and calling `http://localhost:8089`


## Cleanup

When you're done, you can remove the example and all its resources from the cluster.

1. Remove the stack by calling Helm:

    ```bash
    helm delete -n $KYMA_NS $HELM_RELEASE
    ```

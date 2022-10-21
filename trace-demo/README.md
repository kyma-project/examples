# Install OpenTelemtry Demo Application in Kyma

## Overview

https://github.com/open-telemetry/opentelemetry-helm-charts/tree/main/charts/opentelemetry-demo
https://github.com/open-telemetry/opentelemetry-demo

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

### Activate Kyma Tracing

1. (Temporary) Enable tracing feature in operator

    ```bash
    kubectl apply -f https://raw.githubusercontent.com/kyma-project/kyma/main/components/telemetry-operator/config/crd/bases/telemetry.kyma-project.io_tracepipelines.yaml
    kyma deploy -s main --component telemetry --value telemetry.operator.controllers.tracing.enabled=true
    ```

1. (Temporary) Activate Istio tracing
    ```bash
    kyma deploy -s main --component istio -f istio-values.yaml
    ```
    Potentially, you need to restart the relevant workloads

1. Enable Jaeger backend
   ```bash
   kubectl apply -f tracepipeline.yaml
   ```

### Install the application

Run the Helm upgrade command, which installs the chart if not present yet.
```bash
helm upgrade --version 0.9.6 --install --create-namespace -n $KYMA_NS $HELM_RELEASE open-telemetry/opentelemetry-demo -f values.yaml
```

You can either use the [`values.yaml`](./values.yaml) provided in this `kiali` folder, which contains customized settings deviating from the default settings, or create your own `values.yaml` file.

### Verify the application

1. Browse in the frontend by:
   ```bash
   kubectl port-forward svc/$KYMA_NS-frontend 8080:8080
   ````
   and calling `http://localhost:8080`

1. Verify that traces arrive in the Jaeger backend:
   ```bash
   kubectl port-forward -n kyma-system svc/tracing-jaeger-query 16686
   ````
   and calling `http://localhost:16686`

1. Enable failures via the feature flag service:
   ```bash
   kubectl port-forward svc/demo2-featureflagservice 8081
   ````
   and calling `http://localhost:8081`

1. Generate load via the load generator:
   ```bash
   kubectl port-forward svc/demo2-loadgenerator 8089
   ````
   and calling `http://localhost:8089`



## Cleanup

When you're done, you can remove the example and all its resources from the cluster.

1. Remove the stack by calling Helm:

    ```bash
    helm delete -n $KYMA_NS $HELM_RELEASE
    ```

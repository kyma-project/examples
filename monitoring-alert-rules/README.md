# Alert Rules Example

## Overview

This example shows how to configure alert rules in Kyma and how to define a new alert rule for AlertManager.

## Prerequisites

- Kyma as the target deployment environment.

## Installation

You need access to the `kyma-system` Namespace to execute the described steps.

### Configure a new alert
1. Create a ConfigMap for the alert-rule.

    ```bash
    kubectl apply -f deployment/alert-rule-configmap.yaml -n kyma-system
    ```

2. Run the `port-forward` command on the `core-prometheus` service to access the Prometheus dashboard.

    ```bash
    kubectl port-forward pod/prometheus-core-0 -n kyma-system 9090:9090
    ```

    Find the **http-db-service-is-not-running** rule [here](http://localhost:9090/rules).

    As the `http-db-service` Deployment does not the exist, the alert is fired [here](http://localhost:9090/alerts).

### Stop the alert from getting fired
1. Export your Environment as a variable. Replace the `{environment}` placeholder in the following command and run it:

    ```bash
    export KYMA_EXAMPLE_ENV="{environment}"
    ```

2. To stop the alert from getting fired, create a Deployment as follows:
```bash
kubectl apply -f ../http-db-service/deployment/deployment.yaml -n $KYMA_EXAMPLE_ENV
```

### Cleanup
Run the following commands to completely remove the example and all its resources from the cluster:

1. Remove the `http-db-service-is-not-running` alert rule from the cluster.

```bash
kubectl delete cm -n kyma-system -l example=monitoring-alert-rules
````

1. Run the following command to completely remove `http-db-service` and all its resources from the cluster:

```bash
kubectl delete all -l example=http-db-service -n $KYMA_EXAMPLE_ENV
```

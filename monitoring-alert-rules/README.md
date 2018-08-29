# Alerting Rules Example

## Overview

This example shows how to configure alert rules in Kyma.

- Define a new alert rule for AlertManager.

## Prerequisites

- Kyma as the target deployment environment.

## Installation
> Note: You need access to the `kyma-system` Namespace to execute the following steps

### Configure a new alert
1. Create a ConfigMap for the alert-rule.

    ```bash
    kubectl apply -f deployment/alert-rule-configmap.yaml -n kyma-system
    ```

2. Run the `port-forward` command on the `core-prometheus` service to access the prometheus dashboard.
    
    ```bash
    kubectl port-forward pod/prometheus-core-0 -n kyma-system 9090:9090
    ```

    In [`http://localhost:9090/rules`](http://localhost:9090/rules), you will find the rule, **http-db-service-is-not-running**.

    As the `http-db-service` Deployment does not the exist, the alert will be fired in [`http://localhost:9090/alerts`](http://localhost:9090/alerts)

### Stop the alert from getting fired 
- Create a Deployment as follows:

    ```bash
    kubectl apply -f ../http-db-service/deployment/deployment.yaml
    ```

### Cleanup
Run the following commands to completely remove the example and all its resources from the cluster:

1. Remove the `http-db-service-is-not-running` alert rule from the cluster.


```bash
kubectl delete cm -n kyma-system -l example=monitoring-alert-rules
````

2. Remove the `http-db-service` Deployment which stopped firing of alerts.

```bash
kubectl delete -f ../http-db-service/deployment/deployment.yaml
```

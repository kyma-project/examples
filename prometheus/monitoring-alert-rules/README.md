# Alert Rules Example

## Overview

This example shows how to deploy and view alerting rules in Kyma.

## Prerequisites

* Kyma as the target deployment environment.
* Deployed custom kube-prometheus-stack as described in the [prometheus](../) example.

## Installation

### Add a new alerting rule

1. Create the PrometheusRule resource holding the configuration of your alerting rule.

    ```bash
    kubectl apply -f https://raw.githubusercontent.com/kyma-project/examples/main/prometheus/monitoring-alert-rules/deployment/alert-rule.yaml
    ```

2. Run the `port-forward` command on the `monitoring-prometheus` service to access the Prometheus dashboard.

    ```bash
    kubectl -n ${KYMA_NS} port-forward $(kubectl -n ${KYMA_NS} get service -l app=kube-prometheus-stack-prometheus -oname) 9090
    ```

3. Go to `http://localhost:9090/rules` and find the **pod-not-running** rule.

    Because the `http-db-service` Deployment does not exist, Alertmanager fires an alert listed at `http://localhost:9090/alerts`.

### Stop the alert from getting fired

1. Export your Namespace as a variable:

    ```bash
    export KYMA_APPLICATION_NS="{namespace}"
    ```

2. To stop the alert from getting fired, create the Deployment:

    ```bash
    kubectl apply -f https://raw.githubusercontent.com/kyma-project/examples/main/http-db-service/deployment/deployment.yaml -n $KYMA_APPLICATION_NS
    ```

### Cleanup

Run the following commands to completely remove the example and all its resources from the cluster:

1. Remove the **pod-not-running** alerting rule from the cluster:

    ```bash
    kubectl delete cm -n kyma-system -l example=monitoring-alert-rules
    ```

2. Remove the **http-db-service** example and all its resources from the cluster:

    ```bash
    kubectl delete all -l example=http-db-service -n $KYMA_APPLICATION_NS
    ```

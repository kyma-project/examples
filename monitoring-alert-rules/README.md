# Alert Rules Example

## Overview

This example shows how to configure alert rules in Kyma and how to define a new alert rule for AlertManager.

## Prerequisites

* Kyma as the target deployment environment.

## Installation

You need access to the `kyma-system` Namespace to execute the described steps.

### Add a new alerting rule

1. Create the PrometheusRule resource holding the configuration of your alerting rule. 

    ```bash
    kubectl apply -f deployment/alert-rule.yaml -n kyma-system
    ```

2. Run the `port-forward` command on the `monitoring-prometheus` service to access the Prometheus dashboard.

    ```bash
    kubectl port-forward pod/prometheus-monitoring-0 -n kyma-system 9090:9090
    ```

3. Go to `http://localhost:9090/rules` and find the **pod-is-not-running** rule.

    As the `http-db-service` Deployment does not the exist, Alertmanager fires an alert listed at `http://localhost:9090/alerts`.

### Stop the alert from getting fired

1. Export your Namespace as a variable. Replace the `{namespace}` placeholder in the following command and run it:

    ```bash
    export KYMA_EXAMPLE_NS="{namespace}"
    ```

2. To stop the alert from getting fired, create a Deployment as follows:

    ```bash
    kubectl apply -f ../http-db-service/deployment/deployment.yaml -n $KYMA_EXAMPLE_NS
    ```

### Cleanup

Run the following commands to completely remove the example and all its resources from the cluster:

1. Remove the `pod-is-not-running` alert rule from the cluster.

    ```bash
    kubectl delete cm -n kyma-system -l example=monitoring-alert-rules
    ```

2. Run the following command to completely remove `http-db-service` and all its resources from the cluster:

    ```bash
    kubectl delete all -l example=http-db-service -n $KYMA_EXAMPLE_NS
    ```
For a complete tutorial on how to expose metrics and define alerting rules, see [these](https://kyma-project.io/docs/components/monitoring/#tutorials-tutorials) tutorials.
# Example Tracing

>**CAUTION** This example is outdated and will be updated soon. Please have a look at the [trace-demo](./../trace-demo/) and [jaeger](./../jaeger/) instead.

## Overview

This example illustrates how to enable tracing for a service deployed in Kyma. For demonstration, it creates a [Go application](src/order-front.go). This application uses [http-db-service](../http-db-service) for CRUD operations on orders.

To understand how traces are propagated, see the [Go application](src/order-front.go). See the [`deployment`](deployment/deployment.yaml) file to learn about port naming and setting the `app` label.

## Prerequisites

- Kyma OS >= 2.10.x
- kubectl >= 1.22.x
- Helm 3.x

>**NOTE:** By default, the sampling rate for Istio is set to `1`, where `100` is the maximum value. This means that only 1 out of 100 requests is sent to Jaeger for trace recording which can affect the number of traces displayed for the service. To change this behavior, adjust the `randomSamplingPercentage` setting in Istio's telemetry resource:
```bash
kubectl -n istio-system edit telemetries.telemetry.istio.io kyma-traces
```

## Installation

1. Export your Namespace as a variable by replacing the `{namespace}` placeholder in the following command and running it:

    ```bash
    export KYMA_EXAMPLE_NS="{namespace}"
    ```

2. Deploy the service. Run this command:

    ```bash
    kubectl apply -f deployment/deployment.yaml -n $KYMA_EXAMPLE_NS
    ```

## Get traces from the example service

1. Call the example service to simulate an incoming order. Run:

    ```bash
    curl -H "Content-Type: application/json" -d '{"orderCode" : "007", "orderPrice" : 12.0}' https://order-front-api.{YOUR_CLUSTER_DOMAIN}/orders
    ```

2. Access the Jaeger UI on the cluster at `http://localhost:16686` using port-forwarding:
```bash
kubectl port-forward -n kyma-system svc/tracing-jaeger-query 16686:16686
```

3. Select **order-front** from the list of available services and click **Find Traces**.

4. The UI displays end-to-end traces for the API call that simulated an incoming order.


## Cleanup

To remove all resources related to this example from your Kyma cluster, run this command:

```bash
kubectl delete all,api -l example=tracing -n $KYMA_EXAMPLE_NS
```

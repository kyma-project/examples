# Example Tracing

## Overview

This example illustrates how to enable tracing for a service deployed in Kyma. For demonstration, it creates a [Go application](src/order-front.go). This application uses [http-db-service](../http-db-service) for CRUD operations on orders.

To understand how traces are propagated, see the [Go application](src/order-front.go). See the [`deployment`](deployment/deployment.yaml) file to learn about port naming and setting the `app` label.

## Prerequisites

- Kyma as the target deployment environment.
- Helm for local installation.
- A Namespace to which you deploy the example with the `env: "true"` label. For more information, read the [related documentation](https://kyma-project.io/docs/root/kyma/#details-namespaces).


## Installation

### Local installation

> **NOTE:** If you use a local Deployment of Kyma on Minikube,  be aware that Jaeger installation is optional, and you cannot install it locally by default. However, you can install it on a Kyma instance and run it locally using Helm.

1. To install Jaeger, go to the [Kyma resources](https://github.com/kyma-project/kyma/tree/master/resources) directory and run the following command:

```bash
helm install -n jaeger -f jaeger/values.yaml --namespace kyma-system --set-string global.domainName=kyma.local --set-string global.isLocalEnv=true jaeger/
```

2. Follow the instructions in the  **Cluster installation** section (skip step 2). You can access the tracing UI locally at `https://jaeger.kyma.local`.

### Cluster installation

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

2. Access the tracing UI on a cluster at `https://jaeger.{YOUR_CLUSTER_DOMAIN}`.

3. Select **order-front** from the list of available services and click **Find Traces**.

4. The UI displays end-to-end traces for the API call that simulated an incoming order.


## Cleanup

To remove all resources related to this example from your Kyma cluster, run this command:

```bash
kubectl delete all,api -l example=tracing -n $KYMA_EXAMPLE_NS
```

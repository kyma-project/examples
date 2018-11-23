# Example Tracing

## Overview

This example illustrates how to enable tracing for a service deployed in Kyma. For demonstration, it creates a [Go application](src/order-front.go). This application uses [http-db-service](../http-db-service) for CRUD operations on orders.

To understand how traces are propagated, see the [Go application](src/order-front.go). See the [`deployment`](deployment/deployment.yaml) file to learn about port naming and setting the `app` label.

## Prerequisites

- Kyma as the target deployment environment.
- An Environment created in Kyma to which you deploy the example application.
- Helm for local installation.


## Installation

### Locall Installation

If you use a local Deployment of Kyma on Minikube, you need to know that Jaeger installation is optional, and you cannot install it locally by default. However, you can install it on a Kyma instance and run it locally using Helm.

To install Jaeger, go to the [Kyma resources](https://github.com/kyma-project/kyma/tree/master/resources) directory and run the following command:

```bash
helm install -n jaeger -f jaeger/values.yaml --namespace kyma-system --set-string global.domainName=kyma.local --set-string global.isLocalEnv=true jaeger/
```

After installing Jaeger, follow the cluster installation steps (skip step 2). The tracing UI locally accesible at `https://jaeger.kyma.local`.

### Cluster installation

1. Export your Environment as variable by replacing the `{environment}` placeholder in the following command and running it:

    ```bash
    export KYMA_EXAMPLE_ENV="{environment}"
    ```

2. Change the value of the **hostname** parameter from `order-front-api.kyma.local` to `order-front-api.{YOUR_CLUSTER_DOMAIN}` in the [`deployment`](deployment/deployment.yaml) file.

3. Deploy the service. Run this command:

    ```bash
    kubectl apply -f deployment/deployment.yaml -n $KYMA_EXAMPLE_ENV
    ```

## Get traces from the example service

1. Call the example service to simulate an incoming order. Run:

    ```bash
    curl -H "Content-Type: application/json" -d '{"orderCode" : "007", "orderPrice" : 12.0}' https://order-front-api.{domain-of-kyma-cluster}/orders
    ```

2. Access the tracing UI on a cluster at `https://jaeger.{YOUR_CLUSTER_DOMAIN}`.

3. Select **order-front** from the list of available services and click **Find Traces**.

4. The UI displays end-to-end traces for the API call that simulated an incoming order.


## Cleanup

To remove all resources related to this example from your Kyma cluster, run this command:

```bash
kubectl delete deployment,svc,api -l example=tracing -n $KYMA_EXAMPLE_ENV
```

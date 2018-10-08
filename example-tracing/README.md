# Example Tracing

## Overview

This example illustrates how to enable tracing for a service deployed in Kyma. For demostration, it creates a [Go application](./order-front.go). This application uses [http-db-service](../http-db-service) for CRUD operations on orders.

To understand how traces are propagated, see the [Go application](./order-front.go).  See the [example-tracing](./example-tracing.yaml) file to learn about port naming and setting the **app** label.

## Prerequisites

- A [Docker](https://docs.docker.com/install)
- An Environment created in Kyma to which you deploy the example application.


## Installation

### Local installation

1. Build the Docker image of the [Go application](./order-front.go). Run:
    
	```bash
    docker build . -t order-front:latest
    ```

2. Run the image:
    
	```bash
    docker run -p 8080:8080 order-front:latest
    ```

### Cluster installation

1. Export the name of the Kyma Environment to which you want to install the example as an environment variable. Run this command:
    
	```bash
    export KYMA_EXAMPLE_ENV="{YOUR_ENV_NAME}"
    ```

2. Change the value of the **hostname** parameter from `order-front-api.kyma.local` to `order-front-api.{YOUR_CLUSTER_DOMAIN}` in the [example-tracing](./example-tracing.yaml) file.

3. Deploy the service. Run:
    
	```bash
    kubectl apply -f ./example-tracing.yaml -n $KYMA_EXAMPLE_ENV
    ```

## Get traces from the example service

1. Call the example service to simulate an incoming order. Run:

	```bash
	curl -H "Content-Type: application/json" -d '{"orderCode" : "007", "orderPrice" : 12.0}' \ 
		https://order-front-api.{domain-of-kyma-cluster}/orders
	```

2. Access the tracing UI either locally at `https://jaeger.kyma.local` or on a cluster at `https://jaeger.{YOUR_CLUSTER_DOMAIN}`.

3. Select **order-front** from the list of available services and click **Find Traces**.

4. The UI displays end-to-end traces for the API call that simulated an incoming order.


## Cleanup

To remove all resources related to this example from your Kyma cluster, run:

```bash
kubectl delete deployment,svc,api -l example=tracing -n $KYMA_EXAMPLE_ENV
```

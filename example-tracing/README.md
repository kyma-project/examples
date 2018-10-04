# Example Tracing

## Overview

This example illustrates how to enable tracing for a service deployed in Kyma. For demostration, it creates a [go application](./order-front.go). This application uses [http-db-service](../http-db-service) for CRUD operations on orders.

Refer to the following to understand configuration related to tracing

1. [go application](./order-front.go) : 
   - How trace headers are propagated?

1. [example-tracing.yaml](./example-tracing.yaml)
   - port naming 
   - `app` label is configured? 

## Prerequisites

- A [Docker](https://docs.docker.com/install) installation
- Kyma as the target deployment environment
- An environment to which to deploy the example

## Installation

### Local installation

1. Build the Docker image
    
	```bash
    docker build . -t order-front:latest
    ```

1. Run the recently built image
    
	```bash
    docker run -p 8080:8080 order-front:latest
    ```

### Cluster installation

1. Export your environment as a variable by replacing the `<environment>` placeholder in the following command and running it:
    
	```bash
    export KYMA_EXAMPLE_ENV="<environment>"
    ```

1. Change the `hostname: order-front-api.kyma.local` to `hostname: order-front-api.{domain-of-kyma-cluster}` in [example-tracing.yaml](./example-tracing.yaml)

1. Deploy the service
    
	```bash
    kubectl apply -f ./example-tracing.yaml -n $KYMA_EXAMPLE_ENV
    ```

### Test the service

1. To test the service, simulate an order via cURL:

	```bash
	curl -H "Content-Type: application/json" -d '{"orderCode" : "007", "orderPrice" : 12.0}' \ 
		https://order-front-api.{domain-of-kyma-cluster}/orders
	```

1. Access the [tracing UI](https://github.com/kyma-project/kyma/blob/master/docs/tracing/docs/001-overview-tracing.md)

1. Select the `order-front` from the service list and `Find Traces`

1. You should see the end-to-end traces for the API call.


### Cleanup

Clean all deployed example resources from Kyma with the following command:

```bash
kubectl delete deployment,svc,api -l example=tracing -n $KYMA_EXAMPLE_ENV
```

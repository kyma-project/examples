# Bind a service to another service

## Overview

This example shows how to bind service to a provided service in Kyma.
The service exposes an HTTP API to perform database transactions and is bound to a provided service giving access to a MSSQL database in Azure.

The example covers the following tasks:

1. Creating a MSSQL service instance from the Azure Broker.
2. Creating a service binding and binding usage to consume the service instance.
3. Deploying a microservice that performs transactions on the MSSQL database through an HTTP API.

## Prerequisites

- Kyma as the target deployment environment.
- [Kubectl](https://kubernetes.io/docs/tasks/tools/install-kubectl/) CLI tool to deploy the example's resources to Kyma.
- A Namespace to which you deploy the example with the `env: "true"` label. For more information, read the [related documentation](https://kyma-project.io/docs/root/kyma/#details-namespaces).

## Installation

Run the following commands to provision a managed MSSQL database, create a binding to it and deploy the service that will use the binding.

1. Export your Namespace as variable by replacing the `<namespace>` placeholder in the following command and running it:
    ```bash
    export KYMA_EXAMPLE_NS="<namespace>"
    ```

2. Create a MSSQL instance.
    ```bash
    kubectl apply -f deployment/mssql-instance.yaml -n $KYMA_EXAMPLE_NS
    ```

3. Ensure that the MSSQL instance is provisioned and running.
    ```bash
    kubectl get serviceinstance/mssql-instance -o jsonpath='{ .status.conditions[0].reason }' -n $KYMA_EXAMPLE_NS
    ```
    > NOTE: Service instances usually take some minutes to be provisioned.

4. Deploy the service binding, the service binding usage and the http-db-service.
    ```bash
    kubectl apply -f deployment/mssql-binding-usage.yaml -n $KYMA_EXAMPLE_NS
    ```

5. Ensure that the MSSQL service binding is provisioned.
    ```bash
    kubectl get ServiceBinding/mssql-instance-binding -o jsonpath='{ .status.conditions[0].reason }' -n $KYMA_EXAMPLE_NS
    ```
6. Verify that the http-db-service is ready.
    ```bash
    kubectl get pods -l example=service-binding -n $KYMA_EXAMPLE_NS
    ```

7. Forward the service port to be able to reach the service.
    ```bash
    kubectl port-forward -n $KYMA_EXAMPLE_NS $(kubectl get pod -n $KYMA_EXAMPLE_NS -l example=service-binding | grep http-db-service | awk '{print $1}') 8017
    ```
8. Check that the service works as expected.

    Create an order:
    ```bash
    curl -d '{"orderId":"66", "total":9000}' -H "Content-Type: application/json" -X POST http://localhost:8017/orders
    ```
    Check the created order
    ```bash
    curl http://localhost:8017/orders
    ```

### Cleanup

Run the following command to completely remove the example and all its resources from the cluster:

```bash
kubectl delete all,sbu,servicebinding,serviceinstance -l example=service-binding-service -n $KYMA_EXAMPLE_NS
```

## Troubleshooting

### Azure broker not available in minikube

Currently this example does not have support for minikube and can't be run locally.

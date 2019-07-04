# Bind a Service to a Lambda Function

## Overview

This example shows how to bind lambda functions to a service.
A lambda function connects to a Redis instance and retrieves statistics for this service.

The example covers these tasks:

1. Create a Redis instance.
2. Create a function using Kubeless CLI.
3. Bind the function to the Redis Service.
4. Call the lambda function.

## Prerequisites

- Install the Kubeless CLI as described in the [Kubeless installation guide](https://kubeless.io/docs/quick-start/).


## Installation

Apply a battery of `yaml` files to run the example.

### Steps

1. Export your Namespace as a variable by replacing the `{namespace}` placeholder in the following command and running it:
    ```bash
    export KYMA_EXAMPLE_NS="{namespace}"
    export KYMA_EXAMPLE_DOMAIN="{kyma domain}"
    ```

2. Create a Redis instance:
    ```bash
    kubectl apply -f deployment/redis-instance.yaml -n $KYMA_EXAMPLE_NS
    ```

3. Ensure that the Redis instance is provisioned:
    ```bash
    kubectl get serviceinstance/redis-instance -o jsonpath='{ .status.conditions[0].reason }' -n $KYMA_EXAMPLE_NS
    ```

4. Create a Redis client through a lambda function, along with the ServiceBinding, and the ServiceBindingUsage custom resource:    
    ```bash
    kubectl apply -f deployment/lambda-function.yaml -n $KYMA_EXAMPLE_NS
    ```

5. Ensure that the Redis ServiceBinding works:
   ```bash
   kubectl get servicebinding/redis-instance-binding -o jsonpath='{ .status.conditions[0].reason }' -n $KYMA_EXAMPLE_NS
   ```

6. Verify that the lambda function is ready:
    ```bash
    kubeless function ls redis-client -n $KYMA_EXAMPLE_NS
    ```

7. Trigger the function.
    The information and statistics about the Redis server appear in the logs of the function Pod.
    ```bash
     curl -ik https://redis-client.$KYMA_EXAMPLE_DOMAIN
    ```

### Cleanup

Use this command to remove the example and all its resources from your Kyma cluster:

```bash
kubectl delete all,api,function,servicebinding,serviceinstance,servicebindingusage -l example=service-binding -n $KYMA_EXAMPLE_NS
```

## Troubleshooting

Make sure the password is injected correctly into the Pod. The password should match the one in the [redis-instance.yaml](./deployment/redis-instance.yaml)

```bash
kubectl exec -n $KYMA_EXAMPLE_NS -it $(kubectl get po -n $KYMA_EXAMPLE_NS -l example=service-binding-lambda --no-headers | awk '{print $1}') bash

env | grep -i redis_password
```

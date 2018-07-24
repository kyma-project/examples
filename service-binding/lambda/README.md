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

- Install the Kubeless CLI as described in the [Kubeless installation guide](https://github.com/kubeless/kubeless#installation).

## Installation

Apply a battery of yaml files to run the example.

### Steps

1. Export your environment as variable by replacing the `<environment>` placeholder in the following command and running it:
    ```bash
    export KYMA_EXAMPLE_ENV="<environment>"
    ```

2. Create a Redis instance and a service binding to the instance.
    ```bash
    kubectl apply -f deployment/redis-instance.yaml,deployment/redis-instance-binding.yaml -n $KYMA_EXAMPLE_ENV
    ```

3. Ensure that the Redis instance and Redis binding service are provisioned.
    ```bash
    kubectl get serviceinstance/redis-instance -o jsonpath='{ .status.conditions[0].reason }' -n $KYMA_EXAMPLE_ENV

    kubectl get servicebinding/redis-instance-binding -o jsonpath='{ .status.conditions[0].reason }' -n $KYMA_EXAMPLE_ENV
    ```

4. Create a lambda function as Redis client.
    ```bash
    kubectl apply -f deployment/lambda-function.yaml -n $KYMA_EXAMPLE_ENV
    ```

5. Create a ServiceBindingUsage resource.
    ```bash
    kubectl apply -f deployment/service-binding-usage.yaml -n $KYMA_EXAMPLE_ENV
    ```

6. Verify that the lambda function is ready.
    ```bash
    kubeless function ls redis-client -n $KYMA_EXAMPLE_ENV
    ```

7. Trigger the function.
    The information and statistics about the Redis server appear in the logs of the function Pod.
    ```bash
    kubeless function call redis-client -n $KYMA_EXAMPLE_ENV
    ```

### Cleanup

Use this command to remove the example and all its resources from your Kyma cluster:

```bash
kubectl delete all,function,servicebinding,serviceinstance,servicebindingusage -l example=service-binding -n $KYMA_EXAMPLE_ENV
```

## Troubleshooting

Make sure the password is injected correctly into the Pod. The password should match the one in the [redis-instance.yaml](./deployment/redis-instance.yaml)

```bash
kubectl exec -n $KYMA_EXAMPLE_ENV -it $(kubectl get po -n $KYMA_EXAMPLE_ENV -l example=service-binding --no-headers | awk '{print $1}') bash

env | grep -i redis_password
```

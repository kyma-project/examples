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

- Install the Kubeless CLI as described in the [Kubeless installation guide](https://github.com/vmware-archive/kubeless/blob/master/docs/quick-start.md).


## Installation

Apply a battery of `yaml` files to run the example.

### Steps

1. Export your Namespace as a variable by replacing the `{namespace}` placeholder in the following command and running it:
    ```bash
    export K8S_NAMESPACE="{namespace}"
    export KYMA_EXAMPLE_DOMAIN="{kyma domain}"
    ```

2. Create a Redis instance:
    ```bash
    kubectl apply -f deployment/redis-instance.yaml -n $K8S_NAMESPACE
    ```

3. Ensure that the Redis instance is provisioned:
    ```bash
    kubectl get serviceinstance/redis-instance -o jsonpath='{ .status.conditions[0].reason }' -n $K8S_NAMESPACE
    ```

4. Create a Redis client through a lambda function, along with the ServiceBinding, and the ServiceBindingUsage custom resource:    
    ```bash
    kubectl apply -f deployment/lambda-function.yaml -n $K8S_NAMESPACE
    ```

5. Ensure that the Redis ServiceBinding works:
   ```bash
   kubectl get servicebinding/redis-instance-binding -o jsonpath='{ .status.conditions[0].reason }' -n $K8S_NAMESPACE
   ```

6. Verify that the lambda function is ready:
    ```bash
    kubeless function ls redis-client -n $K8S_NAMESPACE
    ```

7. Trigger the function.
    The information and statistics about the Redis server appear in the logs of the function Pod.
    ```bash
     curl -ik https://redis-client.$KYMA_EXAMPLE_DOMAIN
    ```

### Cleanup

Use this command to remove the example and all its resources from your Kyma cluster:

```bash
kubectl delete all,api,function,servicebinding,serviceinstance,servicebindingusage -l example=service-binding -n $K8S_NAMESPACE
```

## Troubleshooting

Make sure the password is injected correctly into the Pod. The password should match the one in the [redis-instance.yaml](./deployment/redis-instance.yaml)

```bash
kubectl exec -n $K8S_NAMESPACE -it $(kubectl get po -n $K8S_NAMESPACE -l example=service-binding-lambda --no-headers | awk '{print $1}') bash

env | grep -i redis_password
```

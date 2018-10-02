# Event Bus Lambda Subscription

## Overview

This example shows how to subscribe lambda functions to the **Kyma Event Bus** to receive events.

This example shows how to perform the following tasks:

1. Create a kubeless function to handle published events.
2. Subscribe the created function to a topic, so that events can activate or trigger it.
3. Test triggering of the function by publishing events.

## Prerequisites

- Install the Kubeless CLI as described in the [Kubeless installation guide](https://kubeless.io/docs/quick-start/).

- After installing the CLI, export the name of the `kubeless config` configmap resource name and the `namespace` where it is located into the shell environment.

- An environment to which to deploy the example.

```bash
export KUBELESS_CONFIG=core-kubeless-config
export KUBELESS_NAMESPACE=kyma-system
```

To get information about currently-supported runtimes, use:

```bash
kubeless get-server-config
```

## Installation

1. Export your environment as variable by replacing the `<environment>` placeholder in the following command and running it:
    ```bash
    export KYMA_EXAMPLE_ENV="<environment>"
    ```

2. Create a kubeless function to receive a JSON event
    >**NOTE:** The function must be able to receive `POST` requests.

    You can find the function file [here](js/hello-with-data.js)

    ```bash
    kubeless function deploy hello-with-data --label example=event-bus-lambda-subscription --runtime nodejs8 --handler hello-with-data.main --from-file js/hello-with-data.js -n $KYMA_EXAMPLE_ENV
    ```

3. Set your environment on the subscription endpoint
    ```bash
    #Linux
    sed -i "s/<environment>/$KYMA_EXAMPLE_ENV/g" deployment/subscription.yaml
    #OSX
    sed -i '' "s/<environment>/$KYMA_EXAMPLE_ENV/g" deployment/subscription.yaml
    ```
    To manually edit [subscription.yaml](./deployment/subscription.yaml), replace the `<environment>` placeholder in the endpoint with your environment.


4. Subscribe the function to an Event
    ```bash
    kubectl apply -f deployment/event-activation.yaml,deployment/subscription.yaml -n $KYMA_EXAMPLE_ENV
    ```

5. Test publishing Events to the subscribed function
    - Start a sample publisher.
        The system creates the [`Publisher`](deployment/sample-publisher.yaml) deployment.
        ```bash
        kubectl apply -f deployment/sample-publisher.yaml -n $KYMA_EXAMPLE_ENV
        ```

    - Spawn a shell inside the publisher.
        ```bash
        kubectl exec -n $KYMA_EXAMPLE_ENV $(kubectl get pods -n $KYMA_EXAMPLE_ENV -l app=sample-publisher --output=jsonpath={.items..metadata.name}) -c sample-publisher -i -t -- sh
        ```

    - Publish an Event inside the publisher container.
        >**NOTE:** The **data** field expects a `JSON` object.
        ```bash
        curl -i \
        -H "Content-Type: application/json" \
        -X POST http://core-publish.kyma-system:8080/v1/events \
        -d '{"source-id": "stage.commerce.kyma.local", "event-type": "hello", "event-type-version": "v1", "event-time": "2018-11-02T22:08:41+00:00", "data": { "order-number": 123 }}'
        ```

    - Verify that tailing logs for the function Pod trigger the function.
        ```bash
        kubectl logs -f $(kubectl get po -n $KYMA_EXAMPLE_ENV -l function=hello-with-data --no-headers | grep -i running | awk '{print $1}') -c hello-with-data -n $KYMA_EXAMPLE_ENV
        ```

### Cleanup

Run this command to remove the example and all of its resources:

```bash
kubectl delete all -l example=event-bus-lambda-subscription -n $KYMA_EXAMPLE_ENV
```

# Event Bus Service Subscription

## Overview

This example demonstrates how to publish and consume events using the Event Bus with the event publishing API and push-to-webhook support available in Kyma.

### Details

Kyma comes with the NATS Streaming messaging cluster. Instead of interacting with the NATS Streaming cluster directly, use the publish HTTP API and the push subscription support that Kyma offers to enable basic event publishing and event delivery through push notifications to an HTTP endpoint.

Access the event publishing API from the cluster through port `8080` on the following host names:

* `core-publish.kyma-system`
* `core-publish.kyma-system.svc.cluster.local`

## Prerequisites

* A [Docker](https://docs.docker.com/install) installation if modification of the image is necessary.
* A [kubectl](https://kubernetes.io/docs/tasks/tools/install-kubectl/) installation.
* Kyma as the target deployment environment.
* An environment to which to deploy the example.


## Installation

1. Export your environment as variable by replacing the `<environment>` placeholder in the following command and running it:
    ```bash
    export KYMA_EXAMPLE_ENV="<environment>"
    ```

2. Deploy the example subscriber:
    ```bash
    kubectl apply -f example-subscriber.yaml -n $KYMA_EXAMPLE_ENV
    ```

3. Enable the event activation
    ```bash
    kubectl apply -f example-event-activation.yaml -n $KYMA_EXAMPLE_ENV
    ```
4. Set your environment on the subscription endpoint
    ```bash
    #Linux
    sed -i "s/<environment>/$KYMA_EXAMPLE_ENV/g" example-subscription.yaml
    #OSX
    sed -i '' "s/<environment>/$KYMA_EXAMPLE_ENV/g" example-subscription.yaml
    ```
    To manually edit [example-subscription.yaml](./example-subscription.yaml), replace the `<environment>` placeholder in the endpoint with your environment.

5. Deploy the example subscription:
    ```bash
    kubectl apply -f example-subscription.yaml -n $KYMA_EXAMPLE_ENV
    ```

6. Deploy the example publisher:
    ```bash
    kubectl apply -f example-publisher.yaml -n $KYMA_EXAMPLE_ENV
    ```

7. Verify that the example subscriber and publisher Pods are running:
    ```bash
    kubectl get pods -n $KYMA_EXAMPLE_ENV
    ```
    The system eventually applies EgressRule, which allows the example publisher container to download `cURL`. Then, Istio grants access and the example publisher Pod starts.

8. Once the example subscriber and publisher Pods are running, access the example publisher container through SSH:
    ```bash
    kubectl exec -n $KYMA_EXAMPLE_ENV $(kubectl get pods -n $KYMA_EXAMPLE_ENV -l app=example-publisher --output=jsonpath={.items..metadata.name}) -c example-publisher -i -t -- sh
    ```

9. Publish a message:
    ```bash
    curl -i \
        -H "Content-Type: application/json" \
        -X POST http://core-publish.kyma-system:8080/v1/events \
        -d '{"source": {"source-namespace": "com.github", "source-type": "commerce", "source-environment": "ec-qa"}, "event-type": "test-event-bus", "event-type-version": "v1", "event-time": "2018-11-02T22:08:41+00:00", "data": {"event":{"customer":{"customerID": "1234", "uid": "rick.sanchez@mail.com"}}}}'

    # or use the fully-qualified event publishing API service name
    curl -i \
        -H "Content-Type: application/json" \
        -X POST http://core-publish.kyma-system.svc.cluster.local:3000/v1/events \
        -d '{"source": {"source-namespace": "com.github", "source-type": "commerce", "source-environment": "ec-qa"}, "event-type": "test-event-bus", "event-type-version": "v1", "event-time": "2018-11-02T22:08:41+00:00", "data": {"event":{"customer":{"customerID": "1234", "uid": "rick.sanchez@mail.com"}}}}'
    ```
    > **NOTE:** To send multiple events, change the `orderCode` to have a unique value for each POST request.


### Tracing

To view the traces:

1. [Access the tracing UI](https://github.com/kyma-project/kyma/blob/master/docs/tracing/docs/001-overview-tracing.md#access-jaeger).
2. Select `webhook-service`.
3. Click `Find Traces`.

### Cleanup

Exit the example publisher container and perform a cleanup using these commands:

```bash
exit
kubectl delete deployment,svc,subscription,eventactivation -l example=event-bus -n $KYMA_EXAMPLE_ENV
```

## Troubleshooting

NATS Streaming does not support topic and subject deletion. So, data such as topics and messages from previous runs accumulates if you repeat the installation steps multiple times.

# Event Bus Service Subscription

## Overview

This example demonstrates how to publish and consume Events using the Event Bus with the Event publishing API and push-to-webhook support available in Kyma.

### Details

Kyma comes with the NATS Streaming messaging cluster. Instead of interacting with the NATS Streaming cluster directly, use the publish HTTP API and the push subscription support that Kyma offers, to enable basic Event publishing and Event delivery through push notifications to an HTTP endpoint.

Access the Event publishing API from the cluster through the `8080` port on the following hostnames:

* `event-bus-publish.kyma-system`
* `event-bus-publish.kyma-system.svc.cluster.local`

## Prerequisites

* A [Docker](https://docs.docker.com/install) installation if modification of the image is necessary.
* A [kubectl](https://kubernetes.io/docs/tasks/tools/install-kubectl/) installation.
* Kyma as the target deployment environment.


## Installation

1. Export your Namespace as a variable by replacing the **{namespace}** placeholder in the following command and running it:
    ```bash
    export KYMA_EXAMPLE_NS="{namespace}"
    ```

2. Deploy the example subscriber:
    ```bash
    kubectl apply -f example-subscriber.yaml -n $KYMA_EXAMPLE_NS
    ```

3. Enable the Event activation:
    ```bash
    kubectl apply -f example-event-activation.yaml -n $KYMA_EXAMPLE_NS
    ```
4. Set your Namespace on the subscription endpoint:
    ```bash
    #Linux
    sed -i "s/<namespace>/$KYMA_EXAMPLE_NS/g" example-subscription.yaml
    #OSX
    sed -i '' "s/<namespace>/$KYMA_EXAMPLE_NS/g" example-subscription.yaml
    ```
    To manually edit [example-subscription.yaml](./example-subscription.yaml), replace the **{namespace}** placeholder in the endpoint with your Namespace.

5. Deploy the example subscription:
    ```bash
    kubectl apply -f example-subscription.yaml -n $KYMA_EXAMPLE_NS
    ```

6. Deploy the example publisher:
    ```bash
    kubectl apply -f example-publisher.yaml -n $KYMA_EXAMPLE_NS
    ```

7. Verify that the example subscriber and publisher Pods are running:
    ```bash
    kubectl get pods -n $KYMA_EXAMPLE_NS
    ```
    The system eventually applies EgressRule, which allows the example publisher container to download cURL. Istio grants access and the example publisher Pod starts.

8. Once the example subscriber and publisher Pods are running, access the example publisher container through SSH:
    ```bash
    kubectl exec -n $KYMA_EXAMPLE_NS $(kubectl get pods -n $KYMA_EXAMPLE_NS -l app=example-publisher --output=jsonpath={.items..metadata.name}) -c example-publisher -i -t -- sh
    ```

9. Publish a message:
    ```bash
    curl -i \
        -H "Content-Type: application/json" \
        -X POST http://event-bus-publish.kyma-system:8080/v1/events \
        -d '{"source-id": "external-application", "event-type": "test-event-bus", "event-type-version": "v1", "event-time": "2018-11-02T22:08:41+00:00", "data": {"event":{"customer":{"customerID": "1234", "uid": "rick.sanchez@mail.com"}}}}'

    # or use the fully-qualified Event publishing API service name
    curl -i \
        -H "Content-Type: application/json" \
        -X POST http://event-bus-publish.kyma-system.svc.cluster.local:8080/v1/events \
        -d '{"source-id": "external-application", "event-type": "test-event-bus", "event-type-version": "v1", "event-time": "2018-11-02T22:08:41+00:00", "data": {"event":{"customer":{"customerID": "1234", "uid": "rick.sanchez@mail.com"}}}}'
    ```
    > **NOTE:** To send multiple Events, change the `orderCode` to have a unique value for each POST request.

10. To access attributes of events which are compliant to CloudEvents [specification](https://github.com/cloudevents/spec/blob/v0.3/http-transport-binding.md), you can access event object as follows:

    |Field|JSON path|
    |---|---|
    |ce-specversion | event.extensions.request.headers['ce-specversion'] |
    |ce-type|event.extensions.request.headers['ce-type']|
    |ce-source|event.extensions.request.headers['ce-source']|
    |ce-id|event.extensions.request.headers['ce-id']|
    |ce-time|event.extensions.request.headers['ce-time']|
    |ce-eventtypeversion|event.extensions.request.headers['ce-eventtypeversion']|
    |ce-knativehistory|event.extensions.request.headers['ce-knativehistory']|
    |content-type|event.extensions.request.headers['content-type']|


### Tracing

To view the traces:

1. [Access the tracing UI](https://kyma-project.io/docs/components/tracing/#overview-overview).
2. Select **event-publish-knative-service**.
3. Click **Find Traces**.

### Cleanup

Exit the example publisher container and perform a cleanup using this command:

```bash
exit
kubectl delete all,subscription,eventactivation -l example=event-subscription-service -n $KYMA_EXAMPLE_NS
```

## Troubleshooting

NATS Streaming does not support topic and subject deletion. Thus, data such as topics and messages from previous runs accumulates if you repeat the installation steps multiple times.

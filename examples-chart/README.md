# Example Chart

## Overview

This chart provides an easy way to deploy and test the examples.

## Prerequisites

- Kubernetes 1.10+
- Kyma as the target deployment environment.
- A Namespace to deploy the examples with the `env: "true"` label. For more information, read the [related documentation](https://kyma-project.io/docs/root/kyma/#details-namespaces).

## Details

### Installation

Configure these options in the [`values.yaml`](values.yaml) file:

| Parameter                        | Description |
|--------------------------------- | -----------: |
| examples.image                   | Image for the examples |
| examples.httpDBService.deploy    | Deploy the [HTTP DB Service](../http-db-service) example. |
| examples.httpDBService.deploymentImage | Deployment image for the HTTP DB Service. |
| examples.httpDBService.testImage | Test image for the HTTP DB Service. |
| examples.eventSubscription.lambda.deploy | Deploy the [Event Subscription lambda](../event-subscription/lambda) example. |
| examples.eventEmailService.deploy | Deploy the [Event Email Service](../event-email-service) example. |
| examples.eventEmailService.deploymentImage | Deployment image for the Event Email Service example. |
| rbac.enabled  | Enable RBAC. |

Deploy the examples:

```
helm install -f values.yaml --name examples --namespace <namespace> .
```

### Testing

Deploy the test Pods defined under [`templates/tests`](templates/tests):
```
helm test --cleanup examples
```
The output of this command shows whether the tests passed or failed.

### Cleanup

Clean up the Namespace:
```
helm delete --purge examples
kubectl delete ns <namespace>
```

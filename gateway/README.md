# Basic Gateway Example

## Overview

This basic example demonstrates how to expose a service or a function through an API in a public or secure manner.

## Prerequisites

- Kyma as the target deployment environment.
- [kubectl](https://kubernetes.io/docs/tasks/tools/install-kubectl/) 1.9.0
- An environment to which to deploy the example.

## Installation

First of all, export your environment as variable by replacing the `<environment>` placeholder in the following command and running it.
```bash
export KYMA_EXAMPLE_ENV="<environment>"
```

Run the following commands in either the 'lambda' or the 'service' folder.

```bash
kubectl create -f deployment.yaml,api-without-auth.yaml -n $KYMA_EXAMPLE_ENV
```

#### Test the APIs without authentication

```bash
# for the function:
curl -i https://hello.kyma.local
# for the service
curl -ik https://http-db-service.ykyma.local/orders
```

#### Test the APIs with authentication  

Run the following commands in either the 'lambda' or the 'service' folder.

```bash
kubectl apply -f api-with-auth.yaml -n $KYMA_EXAMPLE_ENV
```

Retrieve a token by accessing the [Dex Web application](https://dex-web.kyma.local/). Sign in with the admin@kyma.cx email address and the generic password from the [Dex ConfigMap](../../../resources/core/charts/dex/templates/pre-install-dex-config-map.yaml) file.

Call the secured APIs:

```bash
# for the function:
curl -i https://hello.kyma.local -H 'Authorization: Bearer <insert copied token here>'
# for the service
curl -ik https://http-db-service.kyma.local/orders -H 'Authorization: Bearer <insert copied token here>'
```

### Cleanup

Run the following command to completely remove the example and all its resources from the cluster:

```bash
# for the function:
kubectl delete function,api -l example=serverless-lambda -n $KYMA_EXAMPLE_ENV

# for the service:
kubectl delete service,deployment,api -l example=http-db-service -n $KYMA_EXAMPLE_ENV
```

### Troubleshooting

- Exposing an API with authentication takes some time to work. Looking through the API resource can help. A good API should have:

```yaml
authenticationstatus
  code: 2
```

```bash
kubectl get api hello -o yaml
apiVersion: gateway.kyma.cx/v1alpha2
kind: api
metadata:
  annotations:
    kubectl.kubernetes.io/last-applied-configuration: |
      {"apiVersion":"gateway.kyma.cx/v1alpha2","kind":"api","metadata":{"annotations":{},"labels":{"example":"serverless-lambda","function":"hello"},"name":"hello","namespace":"default"},"spec":{"authentication":[{"jwt":{"issuer":"https://dex.kyma.local","jwksUri":"http://dex-service.kyma-system.svc.cluster.local:5556/keys"},"type":"JWT"}],"hostname":"hello.kyma.local","service":{"name":"hello","port":8080}}}
  clusterName: ""
  creationTimestamp: 2018-05-29T08:55:04Z
  generation: 0
  labels:
    example: serverless-lambda
    function: hello
  name: hello
  namespace: default
  resourceVersion: "9563"
  selfLink: /apis/gateway.kyma.cx/v1alpha2/namespaces/default/apis/hello
  uid: fa7805fe-631d-11e8-a42f-3e54503ab95f
spec:
  authentication:
  - jwt:
      issuer: https://dex.kyma.local
      jwksUri: http://dex-service.kyma.local
    type: JWT
  hostname: hello.kyma.local
  service:
    name: hello
    port: 8080
status:
  authenticationStatus:
    code: 2
    resource:
      name: hello
      uid: fa807328-631d-11e8-a42f-3e54503ab95f
      version: "9562"
  virtualServiceStatus:
    code: 2
    resource:
      name: hello
      uid: fa7d02e1-631d-11e8-a42f-3e54503ab95f
      version: "9561"
```

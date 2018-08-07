# Basic Gateway Example

## Overview

This basic example demonstrates how to expose a service or a function through an API in a public or secure manner through console UI or manually through kubectl.

## Prerequisites

- Kyma as the target deployment environment.
- An environment to which to deploy the example.

## Installation

### Exposure through console UI

#### Create service

1. Open Kyma console (https://console.kyma.local/) and choose the Environment in which you want to deploy the example (or create one).
2. Click `Deploy new resource to the environment` button, then select deployment.yaml file from the lambda or service directory in this example and click `Upload`.

#### Expose service without authentication

1. Select `Services` button, then click on the name of the service you've created. The name should be the same as the Service or Function name in deployment.yaml file.
2. In the `Exposed APIs` section click `Expose API` button.
3. Fill the `Host` textbox with anything you want, and click `create`. The name you entered will be referenced later with **\<hostname\>**.
4. Click the `Save` button.

#### Test the APIs without authentication

```bash
# for the function:
curl -ik https://<hostname>.kyma.local
# > 200 Hello world

# for the service:
curl -ik https://<hostname>.kyma.local/orders
# > 200 []
```

#### Expose service with authentication

1. Select `Services` button, then click `Expose API`.
2. Fill the `Host` textbox with anything you want, and click `create`. The name you entered will be referenced later with **\<hostname\>**.
3. Check `Secure API` checkbox and fill Issuer and JWKS URI with custom values or leave the default ones.
4. Click the `Save` button.

#### Fetch token

1. In the `Exposed APIs` section click on the API you've recently exposed.
2. In the `Security` section click on `Fetch token` button.
3. Select all text and copy it to clipboard or simply click `Copy to clipboard` button. Make sure that word Bearer is also copied.
4. That token will be referenced later with **\<token\>**.

_NOTE: You don't actually need to expose API first, you can fetch token also while exposing service. After checking `Secure API` checkbox `Fetch token` button appears, you can then fetch token and just return without saving._

#### Test the APIs with authentication

```bash
# without token:
curl -ik https://<hostname>.kyma.local
# > 401 Origin authentication failed.

# with token:
curl -ik https://<hostname>.kyma.local -H 'Authorization: <token>'
# > 200 Hello world
```

### Manual exposure through kubectl

#### Additional prerequisites

- [kubectl](https://kubernetes.io/docs/tasks/tools/install-kubectl/) 1.10.0
- Token fetched from Console UI (look at the NOTE in Fetch token section). Will be referenced later with **\<token\>**.
- Namespace with istio injection enabled. Will be referenced later with **\<namespace\>**. You can enable istio injection in the namespace by labelling that namespace using:

``` bash
kubectl label namespace <namespace> istio-injection=enabled
```

#### Create service

_NOTE: almost all steps in this tutorial are shown for lambda but you can apply them also to service deployment just by changing directory and in some places also names_

Apply one of the deployment.yaml files from the lambda or service directory in this example.

``` bash
kubectl apply -f ./lambda/deployment.yaml -n <namespace>
```

#### Expose service without authentication

``` bash
kubectl apply -f ./lambda/api-without-auth.yaml -n <namespace>
```

##### How that works

To expose service without authentication omit "authentication" field in the repository or set it to an empty list:

```yaml
authentication: []
```

You can also set:

``` yaml
authenticationEnabled: false
```

what forces authentication to be disabled no matter what is in the "authentication" field.

In this sample, both fields are omitted.

#### Test the APIs without authentication

```bash
curl -ik https://hello.kyma.local
# > 200 Hello world
```

#### Expose service with authentication

``` bash
# Create Api with default authentication settings
kubectl apply -f ./lambda/api-with-default-auth.yaml -n <namespace>

# OR

# Create Api with custom authentication settings
kubectl apply -f ./lambda/api-with-auth.yaml -n <namespace>
```

There are two possible ways of exposing secured Api - using default authentication settings and using custom settings. Authentication settings consist of JWKS URI and Issuer. 

To use default settings omit or set empty list in authentication field and set:

``` yaml
authenticationEnabled: true
```

that forces authentication to be enabled, so if authentication settings are not provided it will use default ones.

In api-with-default-auth.yaml "authentication" field is omitted.

To use custom settings fill "authentication" field like it is done in api-with-auth.yaml sample. You can omit "authenticationEnabled" field or set its value to true.

#### Test the APIs with authentication

```bash
# without token:
curl -ik https://<hostname>.kyma.local
# > 401 Origin authentication failed.

# with token:
curl -ik https://<hostname>.kyma.local -H 'Authorization: <token>'
# > 200 Hello world
```

#### Cleanup

Remove the Api by using kubectl delete on the latest Api resource you've applied. E.g. if you created Api from api-with-auth.yaml file just do:

```bash
kubectl delete -f ./lambda/api-with-auth.yaml -n <namespace>
```

Remove service or lambda by simply using kubectl delete on the resource from which the resource was created. E.g. for lambda:

```bash
kubectl delete -f ./lambda/deployment.yaml -n <namespace>
```

## Troubleshooting

Problem: I created an Api resource with authentication enabled but I receive responses with code 200 after making requests without a JWT token. <br><br> 
Solution 1: Wait. If the cluster is under high workload it may take a while for authentication policies to apply. If you still have that problem after a few seconds look at the Solution 2. <br><br>
Solution 2: Probably there is something wrong with JWKS URI you provided. However, if you used the default settings you can be 100% sure that it is proper. If you are using local deployment of Kyma (on minikube) and internal OIDC Identity Provider (like Dex) make sure that JwksUri is provided as FQDN and that it points directly to keys endpoint e.g. http://dex-service.kyma-system.svc.cluster.local:5556/keys as envoy sidecars need to be able to resolve a domain name to proper inside-cluster or outside-cluster IP address. <br><br>
Solution 3: Check if the pod you've created has injected istio-proxy container. At first, do:
``` bash
kubectl get pods -n <namespace>
```
Find the pod created with deployment.yaml, then copy its name. Then do:
``` bash
kc get pod <pod-name> -n <namespace> -o json | jq '.spec.containers[].name'
```
One of the returned strings should be "istio-proxy". If that is not true, the namespace probably doesn't have istio injection enabled. Go to **Additional prerequisites** section in this README to see how to fix that.
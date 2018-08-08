# Basic Gateway Example

## Overview

This basic example demonstrates how to expose a service or a function through an API in a public or secure manner through the console UI or manually using kubectl.

## Prerequisites

- Kyma as the target deployment environment.
- An environment to which you deploy the example.

## Installation

This section contains installation steps on how to expose a service or a function through the console UI, and also manually, using kubectl.

### Exposure through the console UI

#### Create a service

1. Open the Kyma console (https://console.kyma.local/) and choose or create the Environment in which you want to deploy the example.
2. Click the **Deploy new resource to the environment** button, select the `deployment.yaml` file from the `lambda` or `service` directory in this example, and click **Upload**.

#### Expose a service without authentication

1. Select the **Services** button, click on the name of the service you have created. The name should be the same as the Service or Function name in the `deployment.yaml` file.
2. In the **Exposed APIs** section, click the **Expose API** button.
3. Fill the **Host** textbox and click **Create**. The name you entered will be referred to as the **\{hostname\}**.
4. Click the **Save** button.

#### Test the APIs without authentication

```bash
# To perform a test for the function, use the following command:
curl -ik https://{hostname}.kyma.local
# > 200 Hello world

# To perform a test for the service, use the following command:
curl -ik https://{hostname}.kyma.local/orders
# > 200 []
```

#### Expose a service with authentication

1. Select the **Services** button and click **Expose API**.
2. Fill the **Host** textbox and click **Create**. The name you entered will be referred to as the **\{hostname\}**.
3. Check the **Secure API** checkbox and fill the **Issuer** and **JWKS URI** fields with custom values, or leave the default ones.
4. Click the **Save** button.

#### Fetch token

1. In the **Exposed APIs** section, click on the API you have exposed.
2. In the **Security** section, click the **Fetch token** button.
3. Select all text and copy it to the clipboard, or click the **Copy to clipboard** button. Make sure that wordBeareris also copied.
4. The token will be referred to later as **\{token\}**.

>**NOTE:** You do not have to expose API first. You can also fetch token while exposing a service. After checking the **Secure API** checkbox, the **Fetch token** button appears. Then you can fetch token and return without saving.

#### Test the APIs with authentication

```bash
# To perform a test without token, use the following command:
curl -ik https://{hostname}.kyma.local
# > 401 Origin authentication failed.

# To perform a test with token, use the following command:
curl -ik https://{hostname}.kyma.local -H 'Authorization: <token>'
# > 200 Hello world
```

### Manual exposure using kubectl

There are additional prerequisites to exposing a service or a function manually using kubectl:

- [kubectl](https://kubernetes.io/docs/tasks/tools/install-kubectl/) version 1.10.0
- Token fetched from the Console UI which will be referred to as **\{token\}**. For more details, see the **NOTE** in the **Fetch token** section.
- Namespace with Istio injection enabled which will be referred to as **\{namespace\}**. You can enable Istio injection in the Namespace by labeling it using this command:

``` bash
kubectl label namespace {namespace} istio-injection=enabled
```

#### Create a service

>**NOTE:** Almost all steps in this tutorial refer to lambda but you can apply them also to the service Deployment by changing directory and names in some places.

Apply one of the `deployment.yaml` files from the `lambda` or `service` directory in this example.

``` bash
kubectl apply -f ./lambda/deployment.yaml -n {namespace}
```

#### Expose a service without authentication

``` bash
kubectl apply -f ./lambda/api-without-auth.yaml -n {namespace}
```

#### Test the APIs without authentication

To perform a test, use the following command:

```bash
curl -ik https://hello.kyma.local
# > 200 Hello world
```

#### Expose a service with authentication

There are two possible ways of exposing secured Api - using the default authentication settings, and using the custom settings. Authentication settings consist of the JWKS URI and the Issuer.

``` bash
# Create Api with the default authentication settings:
kubectl apply -f ./lambda/api-with-default-auth.yaml -n {namespace}

# OR

# Create Api with the custom authentication settings:
kubectl apply -f ./lambda/api-with-auth.yaml -n {namespace}
```

#### Test the APIs with authentication

```bash
# To perform a test without token, use the following command:
curl -ik https://<hostname>.kyma.local
# > 401 Origin authentication failed.

# To perform a test with token, use the following command:
curl -ik https://<hostname>.kyma.local -H 'Authorization: <token>'
# > 200 Hello world
```

#### Cleanup

Remove the Api by using the `kubectl delete` command on the latest Api resource you have applied. For example, if you have created an Api from the `api-with-auth.yaml` file, run the following command:

```bash
kubectl delete -f ./lambda/api-with-auth.yaml -n {namespace}
```

Remove a service or lambda by using the `kubectl delete` command` on the file from which the resource was created. See the example for lambda:

```bash
kubectl delete -f ./lambda/deployment.yaml -n {namespace}
```

## Troubleshooting

The problem occurs when there is an Api resource with authentication enabled but after making a request without a JWT token, the received responses code is `200`.

Solution 1: Wait. If the cluster is under high workload it may take a while for authentication policies to apply. If you still have that problem after a few seconds look at the Solution 2.


Solution 2: If you did not use the default settings, there might be something wrong with the JWKS URI you provided. If you use a local Deployment of Kyma on Minikube and the internal OIDC Identity Provider such as Dex, make sure that the JWKS URI is provided as FQDN, and that it points directly to the keys endpoint e.g. http://dex-service.kyma-system.svc.cluster.local:5556/keys. Envoy sidecars must be able to resolve a domain name to the proper inside-cluster or outside-cluster IP address.

Solution 3: Check if the Pod you created has the istio-proxy container injected. Run:

``` bash
kubectl get pods -n {namespace}
```

Find the Pod created with the `deployment.yaml` file and copy its name. Run:

``` bash
kc get pod <pod-name> -n {namespace} -o json | jq '.spec.containers[].name'
```

One of the returned strings should be the istio-proxy. If there is no such string, the Namespace probably does not have Istio injection enabled. Read the additional prerequisites at the beginning of the **Manual exposure using kubectl** section in this document to fix that.
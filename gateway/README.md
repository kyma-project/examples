# Basic Gateway Example

## Overview

This basic example demonstrates how to expose a service or a function through an API in a public or secure manner through the console UI, or manually using kubectl.

## Prerequisites

- Kyma as the target deployment environment.
- A Namespace to which you deploy the example with the `env: "true"` label. For more information, read the [related documentation](https://github.com/kyma-project/kyma/blob/master/docs/kyma/docs/011-details-namespaces.md).

## Installation

This section contains installation steps on how to expose a service or a function through the console UI, and manually, using kubectl.

### Exposure through the console UI

#### Create a service

1. Open the [Kyma console](https://console.kyma.local/) and choose or create the Namespace in which you want to deploy the example.
2. Click the **Deploy new resource to the namespace** button, select the `deployment.yaml` file from the `lambda` or `service` directory in this example, and click **Upload**.

#### Expose a service without authentication

1. Select the **Services** button and click on the name of the service you created. The name should be the same as the service or function name in the `deployment.yaml` file.
2. In the **Exposed APIs** section, click the **Expose API** button.
3. Fill the **Host** text box and click **Create**. The name you entered is referred to as the **\{hostname\}**.
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
2. Fill the **Host** text box and click **Create**. The name you entered is referred to as the **\{hostname\}**.
3. Check the **Secure API** checkbox and fill the **Issuer** and **JWKS URI** fields with custom values, or leave the default ones.
4. Click the **Save** button.

#### Fetch token

1. In the **Exposed APIs** section, click on the API you exposed.
2. In the **Security** section, click the **Fetch token** button.
3. Select the whole text and copy it to the clipboard, or click the **Copy to clipboard** button. Make sure that the word `Bearer` is also copied.
4. The token is later referred to as **\{token\}**.

>**NOTE:** You do not have to expose API first. You can also fetch a token while exposing a service. After checking the **Secure API** checkbox, the **Fetch token** button appears. You can fetch the token and return without saving.

#### Test the APIs with authentication

```bash
# To perform a test without the token, use the following command:
curl -ik https://{hostname}.kyma.local
# > 401 Origin authentication failed.

# To perform a test with the token, use the following command:
curl -ik https://{hostname}.kyma.local -H 'Authorization: {token}'
# > 200 Hello world
```

### Manual exposure using kubectl

There are additional prerequisites to exposing a service or a function manually using kubectl:

- [kubectl](https://kubernetes.io/docs/tasks/tools/install-kubectl/) version 1.10.0
- A token fetched from the Console UI which is later referred to as **\{token\}**. For more details, see the **NOTE** in the **Fetch token** section.

#### Create a service

>**NOTE:** Almost all steps in this tutorial refer to a lambda but you can apply them also to the service Deployment by changing the directory and names in some places.

1. Export your Namespace as a variable. Replace the `{namespace}` placeholder in the following command and run it:

    ```bash
    export KYMA_EXAMPLE_NS="{namespace}"
    ```

2. Apply one of the `deployment.yaml` files from the `lambda` or `service` directory in this example.

    ``` bash
    kubectl apply -f ./lambda/deployment.yaml -n $KYMA_EXAMPLE_NS
    ```

#### Expose a service without authentication

Run this command:

``` bash
kubectl apply -f ./lambda/api-without-auth.yaml -n $KYMA_EXAMPLE_NS
```

#### Test the APIs without authentication

To perform a test, use the following command:

```bash
curl -ik https://hello.kyma.local
# > 200 Hello world
```

#### Expose a service with authentication

There are two possible ways of exposing secured Api, either using the default authentication settings or the custom settings. Authentication settings consist of the JWKS URI and the Issuer.

``` bash
# Create Api with the default authentication settings:
kubectl apply -f ./lambda/api-with-default-auth.yaml -n $KYMA_EXAMPLE_NS

# OR

# Create Api with the custom authentication settings:
kubectl apply -f ./lambda/api-with-auth.yaml -n $KYMA_EXAMPLE_NS
```

#### Test the APIs with authentication

```bash
# To perform a test without the token, use the following command:
curl -ik https://{hostname}.kyma.local
# > 401 Origin authentication failed.

# To perform a test with the token, use the following command:
curl -ik https://{hostname}.kyma.local -H 'Authorization: {token}'
# > 200 Hello world
```

#### Cleanup

### Cleanup

Run the following command to completely remove the example and all its resources from the cluster:

```bash
kubectl delete all -l example=gateway -n $KYMA_EXAMPLE_NS
```

## Troubleshooting

The problem occurs when there is an Api resource with authentication enabled, but after making a request without a JWT token, the received response code is `200`.

**Solution 1:** Wait. If the cluster is under high workload, it can take a while for authentication policies to apply. If you still have the problem after a few seconds, look at the Solution 2.

**Solution 2:** If you did not use the default settings, there can be something wrong with the JWKS URI you provided. If you use a local Deployment of Kyma on Minikube and the internal OIDC Identity Provider such as Dex, make sure that the JWKS URI is provided as FQDN, and that it points directly to the keys endpoint, for example, http://dex-service.kyma-system.svc.cluster.local:5556/keys. Envoy sidecars must be able to resolve a domain name to the proper inside-cluster or outside-cluster IP address.

**Solution 3:** Check if the Pod you created has the istio-proxy container injected. Run this command:

``` bash
kubectl get pods -n $KYMA_EXAMPLE_NS
```

Find the Pod created with the `deployment.yaml` file and copy its name. Run this command:

``` bash
kc get pod {pod-name} -n $KYMA_EXAMPLE_NS -o json | jq '.spec.containers[].name'
```

One of the returned strings should be the istio-proxy. If there is no such string, the Namespace probably does not have Istio injection enabled. Read the additional prerequisites at the beginning of the **Manual exposure using kubectl** section in this document to fix that.

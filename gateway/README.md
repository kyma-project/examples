# Basic Gateway Example

## Overview

This basic example demonstrates how to expose a service through an APIRule in a public or secure manner through the console UI, or manually using kubectl.

## Prerequisites

- Kyma as the target deployment environment.


## Installation

This section contains installation steps on how to expose a service through the console UI, and manually, using kubectl.

### Exposure through the console UI

#### Create a service

1. Open the [Kyma console](https://console.kyma.local/) and choose or create the Namespace in which you want to deploy the example.
2. Click the **Deploy new resource** button, select the `deployment.yaml` file from the `service` directory in this example, and click **Upload**.

#### Expose a service without authentication

1. Select the **Services** button and click on the name of the service you created. The name should be the same as the service name in the `deployment.yaml` file - **http-db-service**.
2. In the **API Rules for http-db-service** section, click the **Expose API** button.
3. Fill the **Name** text box.
4. Fill the **Hostname** text box and click **Create**. The name you entered is referred to as the **\{hostname\}**. The domain next to it is referred to as the **\{domain\}**.

>**NOTE:** There are two ways of exposing a service without authentication - `noop` (default) and `allow`. You can switch the method in the dropdown in the **Access strategies** section.

#### Test the APIRule without authentication

```bash
curl -ik https://{hostname}.{domain}/orders
# > 200 []
```

>**NOTE:** If you are using the Kyma deployed locally, add the `{hostname}.kyma.local` to your hosts file.

#### Expose a service with JWT authentication

1. If you **didn't** follow the steps in **Expose a service without authentication** section, go straight to step 2 of this instruction. If you did, you must delete the previously created APIRule. Select the **API Rules** button, click on the trash can icon next to the APIRule and confirm.
2. Select the **Services** button and click on the name of the service you created. The name should be the same as the service name in the `deployment.yaml` file - **http-db-service**.
3. In the **API Rules for http-db-service** section, click the **Expose API** button.
4. Fill the **Name** text box.
5. Fill the **Hostname** text box. The name you entered is referred to as the **\{hostname\}**. The domain next to it is referred to as the **\{domain\}**.
6. In the **Access strategies** section, click on the dropdown and select the **JWT** field. 
7. Click on the **Configure identity provider...** dropdown and select the **Default** configuration.
8. Click the **Create** button.

#### Fetch JWT token

1. On the main Kyma page, click on the **General Settings** button.
2. In the **Kubeconfig** section, click the **Download config** button.
3. Open the downloaded file in a text editor, select the **token** section and copy it to the clipboard.
4. The token is later referred to as **\{jwt-token\}**.

#### Test the APIs with JWT authentication

```bash
# To perform a test without the token, use the following command:
curl https://{hostname}.{domain}/orders
# > {"error":{"code":401,"status":"Unauthorized","request":"1915853b-9780-4751-b26d-903a179e2941","message":"The request could not be authorized"}}

# To perform a test with the token, use the following command:
curl -ik https://{hostname}.{domain}/orders -H 'Authorization: Bearer {jwt-token}'
# > 200 []
```

>**NOTE:** If you are using the Kyma deployed locally, add the `{hostname}.kyma.local` to your hosts file.

#### Expose a service with OAuth2 authentication

1. If you **didn't** follow the steps in **Expose a service without authentication** or **Expose a service with JWT authentication** section, go straight to step 2 of this instruction. If you did, you must delete the previously created APIRule. Select the **API Rules** button, click on the trash can icon next to the APIRule and confirm.
2. Select the **Services** button and click on the name of the service you created. The name should be the same as the service name in the `deployment.yaml` file - **http-db-service**.
3. In the **API Rules for http-db-service** section, click the **Expose API** button.
4. Fill the **Name** text box.
5. Fill the **Hostname** text box. The name you entered is referred to as the **\{hostname\}**. The domain next to it is referred to as the **\{domain\}**.
6. In the **Access strategies** section, click on the dropdown and select the **OAuth2** field. 
7. Fill the **Required scope** text box with `read, write`.
8. Click the **Create** button.

#### Fetch OAuth2 token

1. On the Namespace main page, click **Deploy new resource** button, select the `oauth2client.yaml` file from the `service` directory in this example, and click **Upload**.
2. Fetch the access token with required scopes. The access token in the response is later referred to as **\{oauth2-token\}**. Run:

    ```bash
    curl https://oauth2.{domain}/oauth2/token -H "Authorization: Basic ZXhhbXBsZS1pZDpleGFtcGxlLXNlY3JldA==" -F "grant_type=client_credentials" -F "scope=read write"
    ```

#### Test the APIs with OAuth2 authentication

```bash
# To perform a test without the token, use the following command:
curl https://{hostname}.{domain}/orders
# > {"error":{"code":401,"status":"Unauthorized","request":"1915853b-9780-4751-b26d-903a179e2941","message":"The request could not be authorized"}}

# To perform a test with the token, use the following command:
curl -ik https://{hostname}.{domain}/orders -H 'Authorization: Bearer {oauth2-token}'
# > 200 []
```
 
>**NOTE:** If you are using the Kyma deployed locally, add the `{hostname}.kyma.local` to your hosts file.

### Manual exposure using kubectl

There are additional prerequisites to exposing a service manually using kubectl:

- [kubectl](https://kubernetes.io/docs/tasks/tools/install-kubectl/) in version specified in the [Kyma documentation](https://kyma-project.io/docs/1.12/root/kyma#installation-install-kyma-locally).
- A JWT token fetched from the Console UI which is later referred to as **\{jwt-token\}**. For more details, see the **Fetch JWT token** section in the **Exposure through the console UI**.
- If you run Kyma locally, add the `http-db-service.kyma.local` to your hosts file.

#### Create a service

1. Export your Namespace as a variable. Replace the `{namespace}` placeholder in the following command and run it:

   ```bash
   export KYMA_EXAMPLE_NS="{namespace}"
   ```
   
2. Export your Kyma domain as a variable. Replace the `{domain}` placeholder in the following command and run it::

    ```bash
   export KYMA_EXAMPLE_DOMAIN="{domain}"
   ```

3. Apply the `deployment.yaml` file from the `service` directory in this example.

   ```bash
   kubectl apply -f ./service/deployment.yaml -n $KYMA_EXAMPLE_NS
   ```

#### Expose a service without authentication

Run this command to create an APIRule:

```bash
kubectl apply -f ./service/api-without-auth.yaml -n $KYMA_EXAMPLE_NS
```

#### Test the APIRules without authentication

To perform a test, use the following command:

```bash
curl -ik https://http-db-service.$KYMA_EXAMPLE_DOMAIN/orders
# > 200 []
```

#### Expose a service with JWT authentication

> **NOTE:** If you followed the steps in **Expose a service without authentication** section, the previously created APIRule will be updated after applying the templates.

The JWT authentication settings require to provide a list of trusted issuers. To create an APIRule with the JWT authentication settings, run: 

```bash
cat <<EOF | kubectl apply -f - 
apiVersion: gateway.kyma-project.io/v1alpha1
kind: APIRule
metadata:
  labels:
    example: gateway-service
  name: http-db-service
  namespace: $KYMA_EXAMPLE_NS
spec:
  gateway: kyma-gateway.kyma-system.svc.cluster.local
  rules:
    - accessStrategies:
        - config:
            jwks_urls:
              - http://dex-service.kyma-system.svc.cluster.local:5556/keys
            trusted_issuers:
              - https://dex.$KYMA_EXAMPLE_DOMAIN
          handler: jwt
      methods:
        - GET
        - POST
        - PUT
        - DELETE
      path: /.*
  service:
    host: http-db-service.kyma.local
    name: http-db-service
    port: 8017
```

Or manually adjust the `https://dex.kyma.local` domain in the `trusted_issuers` section of the `api-with-jwt.yaml` file to fit your setup and run:

```bash
kubectl apply -f ./service/api-with-jwt.yaml -n $KYMA_EXAMPLE_NS
```

#### Test the APIs with JWT authentication

```bash
# To perform a test without the token, use the following command:
curl -ik https://http-db-service.$KYMA_EXAMPLE_DOMAIN/orders
# > {"error":{"code":401,"status":"Unauthorized","request":"530f300a-8269-4564-8d0c-9816c692e7c4","message":"The request could not be authorized"}}

# To perform a test with the token, use the following command:
curl -ik https://http-db-service.$KYMA_EXAMPLE_DOMAIN/orders -H 'Authorization: Bearer {jwt-token}'
# > 200 []
```

#### Expose a service with Oauth2 authentication

> **NOTE:** If you followed the steps in **Expose a service without authentication** or **Expose a service with JWT authentication** section, the previously created APIRule will be updated after applying the templates.

Create an APIRule with the OAuth2 authentication settings:
```bash
kubectl apply -f ./service/api-with-oauth2.yaml -n $KYMA_EXAMPLE_NS
```

#### Fetch the OAuth2 token

1. Create an OAuth2Client:

    ```bash
    kubectl apply -f ./service/oauth2client.yaml -n $KYMA_EXAMPLE_NS
    ```

2. Fetch the access token with required scopes. The access token in the response is referred to as the **\{oauth2-token\}**. Run:

    ```bash
    curl https://oauth2.$KYMA_EXAMPLE_DOMAIN/oauth2/token -H "Authorization: Basic ZXhhbXBsZS1pZDpleGFtcGxlLXNlY3JldA==" -F "grant_type=client_credentials" -F "scope=read write"
    ```

#### Test the APIs with OAuth2 authentication

```bash
# To perform a test without the token, use the following command:
curl -ik https://http-db-service.$KYMA_EXAMPLE_DOMAIN/orders
# > {"error":{"code":401,"status":"Unauthorized","request":"530f300a-8269-4564-8d0c-9816c692e7c4","message":"The request could not be authorized"}}

# To perform a test with the token, use the following command:
curl -ik https://http-db-service.$KYMA_EXAMPLE_DOMAIN/orders -H 'Authorization: Bearer {oauth2-token}'
# > 200 []
```


### Cleanup

Run the following command to completely remove the example and all its resources from the cluster:

```bash
kubectl delete all -l example=gateway-service -n $KYMA_EXAMPLE_NS
kubectl delete apirules.gateway.kyma-project.io -l example=gateway-service -n $KYMA_EXAMPLE_NS
kubectl delete oauth2clients.hydra.ory.sh -l example=gateway-service -n $KYMA_EXAMPLE_NS
```

## Troubleshooting

The problem occurs when there is an Api resource with authentication enabled, but after making a request without a JWT token, the received response code is `200`.

**Solution 1:** Wait. If the cluster is under high workload, it can take a while for authentication policies to apply. If you still have the problem after a few seconds, look at the Solution 2.

**Solution 2:** If you did not use the default settings, there can be something wrong with the JWKS URI you provided. If you use a local Deployment of Kyma on Minikube and the internal OIDC Identity Provider such as Dex, make sure that the JWKS URI is provided as FQDN, and that it points directly to the keys endpoint, for example, http://dex-service.kyma-system.svc.cluster.local:5556/keys. Envoy sidecars must be able to resolve a domain name to the proper inside-cluster or outside-cluster IP address.

**Solution 3:** Check if the Pod you created has the istio-proxy container injected. Run this command:

```bash
kubectl get pods -n $KYMA_EXAMPLE_NS
```

Find the Pod created with the `deployment.yaml` file and copy its name. Run this command:

```bash
kc get pod {pod-name} -n $KYMA_EXAMPLE_NS -o json | jq '.spec.containers[].name'
```

One of the returned strings should be the istio-proxy. If there is no such string, the Namespace probably does not have Istio injection enabled. Read the additional prerequisites at the beginning of the **Manual exposure using kubectl** section in this document to fix that.

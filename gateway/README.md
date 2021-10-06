# Basic Gateway Example

## Overview

This basic example demonstrates how to expose a service through an API Rule in a public or secure manner through the console UI, or manually using kubectl.

## Prerequisites

- Kyma as the target deployment environment.


## Installation

This section contains installation steps on how to expose a service through the console UI, and manually, using kubectl.

### Exposure through the console UI

#### Create a service

1. Open the Kyma console and choose or create the Namespace in which you want to deploy the example.
2. Click the **Deploy new resource** button, select the `deployment.yaml` file from the `service` directory in this example, and click **Upload**.

#### Expose a service without authentication

1. Select the **Services** button and click the name of the service you created. The name should be the same as the service name in the `deployment.yaml` file - **http-db-service**.
2. In the **API Rules for http-db-service** section, click the **Expose API** button.
3. Fill the **Name** text box.
4. Fill the **Hostname** text box and click **Create**. The name you entered is referred to as the **\{hostname\}**. The domain next to it is referred to as the **\{domain\}**.

>**NOTE:** There are two ways of exposing a service without authentication - `noop` (default) and `allow`. You can switch the method using the dropdown menu in the **Access strategies** section.

#### Test the API Rule without authentication

```bash
curl -ik https://{hostname}.{domain}/orders
# > 200 []
```

>**NOTE:** If you use Kyma locally, add `{hostname}.{domain}` to your hosts file.

#### Expose a service with JWT authentication

1. If you **didn't** follow the steps in the **Expose a service without authentication** section, go straight to step 2 of this instruction. If you did, you must delete the previously created API Rule. Select the **API Rules** button, click on the trash can icon next to the API Rule and confirm.
2. Select the **Services** button and click the name of the service you created. The name should be the same as the service name in the `deployment.yaml` file - **http-db-service**.
3. In the **API Rules for http-db-service** section, click the **Expose API** button.
4. Fill the **Name** text box.
5. Fill the **Hostname** text box. The name you entered is referred to as the **\{hostname\}**. The domain next to it is referred to as the **\{domain\}**.
6. In the **Access strategies** section, select the **JWT** field from the dropdown list. 
7. Click the **Configure identity provider...** dropdown menu and select the **Default** configuration.
8. Click the **Create** button.

#### Fetch JWT

1. On the main Kyma page, click the **General Settings** button.
2. In the **Kubeconfig** section, click the **Download config** button.
3. Open the downloaded file in a text editor, select the value in the **token** field and copy it to the clipboard.
4. The token is later referred to as **\{jwt\}**.

#### Test the APIs with JWT authentication

```bash
# To perform a test without the token, use the following command:
curl https://{hostname}.{domain}/orders
# > {"error":{"code":401,"status":"Unauthorized","request":"1915853b-9780-4751-b26d-903a179e2941","message":"The request could not be authorized"}}

# To perform a test with the token, use the following command:
curl -ik https://{hostname}.{domain}/orders -H 'Authorization: Bearer {jwt}'
# > 200 []
```

>**NOTE:** If you use Kyma locally, add `{hostname}.{domain}` to your hosts file.

#### Expose a service with OAuth2 authentication

1. If you **didn't** follow the steps in the **Expose a service without authentication** or **Expose a service with JWT authentication** section, go straight to step 2 of this instruction. If you did, you must delete the previously created API Rule. Select the **API Rules** button, click on the trash can icon next to the API Rule and confirm.
2. Select the **Services** button and click the name of the service you created. The name should be the same as the service name in the `deployment.yaml` file - **http-db-service**.
3. In the **API Rules for http-db-service** section, click the **Expose API** button.
4. Fill the **Name** text box.
5. Fill the **Hostname** text box. The name you entered is referred to as the **\{hostname\}**. The domain next to it is referred to as the **\{domain\}**.
6. In the **Access strategies** section, select the **OAuth2** field from the dropdown list. . 
7. Fill the **Required scope** text box with `read, write`.
8. Click the **Create** button.

#### Fetch OAuth2 token

1. On the Namespace main page, click the **Deploy new resource** button, select the `oauth2client.yaml` file from the `service` directory in this example, and click **Upload**.
2. Fetch the access token with the required scopes. The access token in the response is later referred to as **\{oauth2-token\}**. Run:

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
 
>**NOTE:** If you use Kyma locally, add `{hostname}.{domain}` to your hosts file.

### Manual exposure using kubectl

There are additional prerequisites to exposing a service manually using kubectl:

- [kubectl](https://kubernetes.io/docs/tasks/tools/install-kubectl/) in the version specified in the [Kyma documentation](https://kyma-project.io/docs/#installation-install-kyma-locally). It must be configured to point to your Kyma cluster. For more information, see the document about [getting the kubeconfig file](https://kyma-project.io/docs/components/security/#details-iam-kubeconfig-service-get-the-kubeconfig-file-and-configure-the-cli).
- A JWT fetched from the Console UI which is later referred to as **\{jwt\}**. For more details, see the **Fetch JWT** section in the **Exposure through the console UI**.
- If you run Kyma locally, add the `http-db-service.kyma.local` to your hosts file.

#### Create a service

1. Export your Namespace as a variable. Replace the `{namespace}` placeholder in the following command and run it:

   ```bash
   export KYMA_EXAMPLE_NS="{namespace}"
   ```
   
2. Export your Kyma domain as a variable. Replace the `{domain}` placeholder in the following command and run it:

    ```bash
   export KYMA_EXAMPLE_DOMAIN="{domain}"
   ```

3. Apply the `deployment.yaml` file from the `service` directory in this example.

   ```bash
   kubectl apply -f ./service/deployment.yaml -n $KYMA_EXAMPLE_NS
   ```

#### Expose a service without authentication

Run this command to create an API Rule:

```bash
kubectl apply -f ./service/api-without-auth.yaml -n $KYMA_EXAMPLE_NS
```

#### Test the API Rules without authentication

To perform a test, use the following command:

```bash
curl -ik https://http-db-service.$KYMA_EXAMPLE_DOMAIN/orders
# > 200 []
```

#### Expose a service with JWT authentication

> **NOTE:** If you followed the steps in the **Expose a service without authentication** section, the previously created API Rule will be updated after applying the templates.

The JWT authentication settings require to provide a list of trusted issuers. To create an API Rule with the JWT authentication settings, run: 

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
    host: http-db-service
    name: http-db-service
    port: 8017
EOF
```

You can also manually adjust the `https://dex.kyma.local` domain in the `trusted_issuers` section of the `api-with-jwt.yaml` file to fit your setup. After adjusting the domain, run:

```bash
kubectl apply -f ./service/api-with-jwt.yaml -n $KYMA_EXAMPLE_NS
```

#### Test the APIs with JWT authentication

```bash
# To perform a test without the token, use the following command:
curl -ik https://http-db-service.$KYMA_EXAMPLE_DOMAIN/orders
# > {"error":{"code":401,"status":"Unauthorized","request":"530f300a-8269-4564-8d0c-9816c692e7c4","message":"The request could not be authorized"}}

# To perform a test with the token, use the following command:
curl -ik https://http-db-service.$KYMA_EXAMPLE_DOMAIN/orders -H 'Authorization: Bearer {jwt}'
# > 200 []
```

#### Expose a service with Oauth2 authentication

> **NOTE:** If you followed the steps in the **Expose a service without authentication** or **Expose a service with JWT authentication** section, the previously created API Rule will be updated after applying the templates.

Create an API Rule with the OAuth2 authentication settings:
```bash
kubectl apply -f ./service/api-with-oauth2.yaml -n $KYMA_EXAMPLE_NS
```

#### Fetch the OAuth2 token

1. Create an OAuth2 client:

    ```bash
    kubectl apply -f ./service/oauth2client.yaml -n $KYMA_EXAMPLE_NS
    ```

2. Fetch the access token with the required scopes. The access token in the response is referred to as the **\{oauth2-token\}**. Run:

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

**Could not resolve host:** If you run Kyma locally, make sure you have added the hostnames used in this example to your hosts file.

**No healthy upstream:** Check if the Pod you created is running. Run:

```bash
kubectl get pods -n $KYMA_EXAMPLE_NS
```

Wait until all containers of the Pod are running.

**Problems with the JWT authentication:** Make sure you have provided proper domain name in the **Expose a service with JWT authentication** step.

**Upstream connect error or disconnect/reset before headers:** Check if the Pod you created has the istio-proxy container injected. Run this command:

```bash
kubectl get pods -l example=gateway-service -n $KYMA_EXAMPLE_NS -o json | jq '.items[].spec.containers[].name'
```

One of the returned strings should be the istio-proxy. If there is no such string, the Namespace probably does not have Istio injection enabled. For more information, read the document about the [sidecar proxy injection](https://kyma-project.io/docs/kyma/latest/04-operation-guides/operations/smsh-01-istio-disable-sidecar-injection/).

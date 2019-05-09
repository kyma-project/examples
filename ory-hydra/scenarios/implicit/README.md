# Accessing Lambda functions using ORY Hydra with Dex integration example

## Overview

This example illustrates how to secure lambda functions with Hydra JWT tokens using two components of Kyma platform:
- [Hydra](https://www.ory.sh/docs/hydra/) - OAuth 2.0 and OpenID Connect Server.
- [Dex](https://github.com/dexidp/dex) - an identity service that uses OpenID Connect to drive authentication for other apps.

In this scenario, Hydra Server is responsible for issuing JWT tokens and Dex is used as a login provider. Users can authenticate by any method that is configured in Dex.

The scenario consists of following steps:

- Installation of Kyma with `ory`  and `hydra-dex` components.
- Creating a Lambda function protected with Authn policy that points to the Hydra JWKS endpoint.
- Creating a Hydra Oauth2 client with ability to perform Implicit Grant token requests with **openid** scope.
- Fetching a JWT Token from Hydra with Dex as login provider.
- Calling target Lambda with JWT token issued by Hydra.

## Prerequisites
- cURL

## Steps

### Install necessary components
- read instructions from **chart/hydra-dex** chart.
- Install the **chart/hydra-dex** chart.

### Create a lambda function


 - get valid Hydra **Issuer**: `echo $(kubectl get deployments/ory-hydra-oauth2 -n kyma-system -o go-template='{{range (index .spec.template.spec.containers 0).env}}{{if eq .name "OAUTH2_ISSUER_URL"}}{{.value}}{{end}}{{end}}')`
 - the valid Hydra **JWKS URI** for in-cluster calls is:`http://ory-hydra-oauth2.kyma-system.svc.cluster.local/.well-known/jwks.json`
 - Create a lambda function and expose it as HTTPS service. Provide valid **Issuer** and **JWKS URI**. Pay attention to the **Issuer** value, it has to match exactly.


### Create OpenID Connect client



* Replace `<domainName>` with the proper domain name of your ingress gateway for cluster installations or **kyma.local** for local installation.
  Run:  `export DOMAIN_NAME=<domainName>`

* Run  `curl -ik -X POST "https://oauth2-admin.$DOMAIN_NAME/clients" -d '{"grant_types":["implicit"], "response_types":["id_token"], "scope":"openid", "redirect_uris":["http://localhost:8080/callback"], "client_id":"implicit-client", "client_secret":"some-secret"}'`

_Note: The client is using `http://localhost:8080/callback` redirect URI. This doesn't have to be an URL of any real application. Since OpenID Connect Implict flow is browser-based, it's only important to have a valid URL here. The final redirect of the flow contains the token. In case the application does not exist, the browser will report a non-existing address, but the token will be present in the address bar._

### Fetch a JWT token
* create an OpenID Connect token request: `echo "http://oauth2.$DOMAIN_NAME/oauth2/auth?client_id=implicit-client&response_type=id_token&scope=openid&state=8230b269ffa679e9c662cd10e1f1b145&redirect_uri=http://localhost:8080/callback&nonce=$(date | md5)"`
* Copy the URL into your browser
* Authenticate. After successful authentication, you should be redirected to the address that looks like this: `http://localhost:8080/callback#id_token=eyJ...&state=8230b269ffa679e9c662cd10e1f1b145`
* Copy the **id_token** value from the browser address bar. It is long!
* `export JWT=<copied id_token value>`

### Call the lambda with the token

* Export the URL of your Lambda function.  You can copy the URL from the Kyma console. An example: `export LAMBDA_URL="https://demo-hydra-production.kyma.local/"`
* Ensure Lambda host is added to your **/etc/hosts** for local installations
* Call the lambda: `curl -ik -X GET "${LAMBDA_URL}/" -H "Authorization: Bearer ${JWT}"`



## Troubleshooting

- If you see: **Origin authentication failed.** message, check isito-pilot logs for the following entries:
  `warn    Failed to fetch jwt public key from "http://ory-hydra-oauth2.kyma-system.svc.cluster.local/.well-known/jwks.json"`
Istio-pilot sometimes has problems in accessing jwks endpoint to perform JWT verification. Restarting Istio-pilot helps. Note that after restart, not all envoy proxies are updated immediately. Retry for at least two minutes before giving up.

- If **Origin authentication failed.** message is still persistent, decode the Hydra JWT and ensure it's nonce is the same as the one you provided at the token request. Perhaps you're using an invalid token? Re-create the token and try again with a new one.

- If **Origin authentication failed.** message is still persistent, inspect the Authentication Policy for the Lambda. The name of the policy matches the name of the lambda, and it's created in the same namespace as the lambda. Are the issuer and jwksUri fields valid? Try do delete the policy (backup it first). Can you call the lambda now? Re-create the policy and inspect istio-pilot logs.

- If **Origin authentication failed.** message is still persistent, try to call the Hydra JWKS URI from within the cluster using a pod with bash and curl installed. Perhaps there's some networking issue not directly related to token validation.


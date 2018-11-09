# Lambda that calls OCC API in the context of the end user

## Overview

This example shows how to call an OCC API and integrate with the SAP Commerce.

The following diagram illustrates how the web application included in the example interacts with Enterprise Commerce and the lambda function to call the OCC API in the context of a specific user.
![](./diagram.png)

The flow of operations:
1. The Web UI application redirects the user to the Enterprise Commerce to perform authentication.
2. The Enterprise Commerce redirects the user to the Web UI. The access token and the ID token are passed as query parameters.
3. The Web UI application calls the lambda using the access token and the ID token. The access token is passed in the `occ-token` custom header, whereas the ID token is passed in the `Authorization` header.
4. The Application Connector uses the access token to call the OCC API.    

## Prerequisites

To follow this example you need:

- an Environment created in Kyma. For more information, read the [related documentation](https://github.com/kyma-project/kyma/blob/master/docs/kyma/docs/011-details-environments.md).

>**NOTE:** Execute all commands in this example against the Environment you use (-n {Name_Of_The_Environment}).

- an OCC API registered in the Application Connector's Registration API.

## Installation

This section describes how to deploy, expose, and call the lambda function you deployed.

### Deploy the lambda function

1. Before you deploy the lambda, modify the [call-ec-function.yml](call-ec-function.yml) file by setting the `name` container environment variable value to the full URL of the access service that represents the underlying OCC API. For example, for the access service with the `http://ec-default-b515fe5b-7d06-446d-858f-db0128792bb8.kyma-integration` URL:  

```
containers:
- env:
  - name: http://ec-default-b515fe5b-7d06-446d-858f-db0128792bb8.kyma-integration
    value: target_url
```

2. To deploy the lambda function in the Namespace of your choice, run this command:
```
kubectl apply -f ./call-ec-function.yml
```

3. To check if the lambda function is successfully created in the Namespace of your choice, run this command:
```
kubectl describe function/calculate-promotion
```

### Expose the lambda function

1. Edit the [expose-ec-function.yml](expose-ec-function.yaml) file, and change the `$YOUR_DOMAIN` placeholder in the `hostname` parameter value to the domain of your Kyma cluster.
For example, a local Kyma cluster uses the `kyma.local` domain.

2. Expose the lambda function. Specify the Namespace to which you deployed the lambda in the command. Run:
```
kubectl apply -f ./expose-ec-function.yml
```

### Deploy the Web UI application

1. Edit the [web-ui/.env](web-ui/.env) file and set URLs appropriate for used EC instance for:
   **REACT_APP_OAUTH2_ISSUER**, **REACT_APP_OAUTH2_JWKS_URI**, and **REACT_APP_OAUTH2_AUTHORIZE_URL**.
   For most installations only `mycommerce.kyma.cx` domain has to be replaced with EC instance domain.

1. Deploy the Web UI application in the Namespace of your choice and create a Kubernetes service and an Api resource for it:

  ```bash
  kubectl apply -f ./web-ui/deployment.yaml
  ```

1. Expose the Web UI application:

   1. Check the name of the Web UI application Pod (name starting with `call-ec-web-ui-`):

      ```bash
      kubectl get pods
      ```

   1. Forward Web UI application port, so the application will be available on the same port of localhost:

      ```bash
      kubectl port-forward {Pod_Name} 3000
      ```

### Call the lambda

The lambda function expects these query parameters:

- `user-id` - an identifier of an EC user.
- `threshold` - a minimal value of ordered items needed to get a promotion.

After you expose the function, you can invoke by sending this request:
```bash
curl -i -G https://calculate-promotion.{YOUR.CLUSTER.DOMAIN} -d "user-id={customer_id}" -d "threshold=1000" -H "occ-token: {EC_access_token}" -H "Authorization: Bearer {EC_ID_token}"
```

### Use the UI

To access the UI, go to `http://localhost:3000` and authenticate in Enterprise Commerce. After you authenticate, the system redirects you back to the Web UI application.
Find the access token and the ID token issued for your application by Enterprise Commerce in the **Authentication** section. Expand **The Function** section to call the function and check the results of the operation.

### Cleanup
Delete all objects created by the example using this command:
```bash
kubectl delete all -l example=call-ec
```

## Troubleshooting

### Lambda returns Internal Server Error

Internal Server Error may occur in one of the following situations:
- You don't specify the required query parameters when you invoke the lamda.
- You don't set the correct URL of the access service that represents the underlying OCC API in the [call-ec-function.yml](call-ec-function.yml) file.
- The underlying OCC API is not available.
- The underlying OCC API returns a status different than `200`.

Check the logs of the lambda function to troubleshoot the issue.

### Check the log file

Find the Pod that contains the lambda function. List all Pods in a given Namespace using this command:
```bash
kubectl get po
```

Get the logs of the chosen Pod:
```bash
kubectl log {Pod_name}
```      

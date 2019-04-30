# Auth Token Relay Example

## Overview

This example shows how to call an Application API from within Kyma in the conext of an incoming user call.  

## Prerequisites

- Kyma as the target deployment environment.
- A Namespace to which to deploy the example with the `env: "true"` label. For more information, read the [related documentation](https://github.com/kyma-project/kyma/blob/master/docs/kyma/docs/03-02-namespaces.md).

## Installation

Run the following commands to deploy the published service to Kyma:

1. Export your Namespace as variable by replacing the `{namespace}` placeholder in the following command and running it:

    ```bash
    export KYMA_EXAMPLE_NS="{namespace}"
    ```
2. Deploy the setup:
    ```bash
    kubectl -n $KYMA_EXAMPLE_NS apply -f deployment/setup.yaml
    ```
3. Pair and register the API
    ```bash
    kubectl -n $KYMA_EXAMPLE_NS get tokenrequest school-app -o jsonpath='{.status.url}'
    curl -d '{"url":"{TOKEN}"}' -H "Content-Type: application/json" -X POST https://school-mock.{DOMAIN}/connection
    curl -d '{}' -H "Content-Type: application/json" -X POST https://school-mock.{DOMAIN}/local/apis/Schools%20API/register
    ```
4. Bind API to your lambda  
    ```bash
    kubectl -n $KYMA_EXAMPLE_NS get serviceclass

    kubectl -n $KYMA_EXAMPLE_NS apply -f deployment/binding.yaml
    ```

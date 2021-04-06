# Rafter

## Overview

This example illustrates how to use [Rafter](https://kyma-project.io/docs/master/components/rafter/) to store static webpages.

By default, [MinIO](https://min.io/) stores all resources on a cluster, but it also allows you to use different cloud providers. Read Rafter [tutorials](https://kyma-project.io/docs/components/rafter#tutorials-tutorials) for more information.

## Prerequisites

- [kubectl](https://kubernetes.io/docs/tasks/tools/install-kubectl/)
- [Kyma](https://kyma-project.io/docs/) as the target deployment environment

## Installation

1. Export a GitHub webpage URL of a ready-to-use webpage sources.

    Example:

    ```bash
    export GH_WEBPAGE_URL=https://github.com/kyma-project/examples/archive/main.zip
    ```

2. Apply a Bucket custom resource (CR):

    ```bash
    cat <<EOF | kubectl apply -f -
    apiVersion: rafter.kyma-project.io/v1beta1
    kind: Bucket
    metadata:
      name: pages
      namespace: default
    spec:
      region: "us-east-1"
      policy: readonly
    EOF
    ```

3. Apply an Asset CR:

    ```bash
    cat <<EOF | kubectl apply -f -
    apiVersion: rafter.kyma-project.io/v1beta1
    kind: Asset
    metadata:
      name: webpage
      namespace: default
    spec:
      source:
        url: ${GH_WEBPAGE_URL}
        mode: package
        filter: /rafter/webpage/
      bucketRef:
        name: pages
    EOF
    ```

4. Check the value of the **phase** field:

    ```bash
    kubectl get assets.rafter.kyma-project.io webpage -o jsonpath='{.status.phase}'
    ```

    You should get a result exactly like this one:

    ```test
    Ready
    ```

    >**Note:** If the state is `Pending`, wait for a few seconds and try again.

5. Export and merge the values of the **baseUrl** field and the path to the `index.html` file from the Asset CR, and then open it in a default web browser:

    ```bash
    open $(kubectl get assets.rafter.kyma-project.io webpage -o jsonpath='{.status.assetRef.baseUrl}{"/examples-main/rafter/webpage/index.html"}')
    ```

### Cleanup

1. Delete the Asset CR:

    ```bash
    kubectl delete assets.rafter.kyma-project.io webpage
    ```

2. Delete the Bucket CR:

    ```bash
    kubectl delete buckets.rafter.kyma-project.io pages
    ```

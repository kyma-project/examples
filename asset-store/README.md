# Asset Store

## Overview

This example illustrates how to use the [Asset Store](https://kyma-project.io/docs/components/asset-store) to store static webpages.

By default, [Minio](https://min.io/) stores all resources on a cluster, but it also allows you to use different cloud providers. Read the Asset Store [tutorials](https://kyma-project.io/docs/components/asset-store#tutorials-tutorials) for more information.

## Prerequisites

- [kubectl](https://kubernetes.io/docs/tasks/tools/install-kubectl/)
- [Kyma](https://kyma-project.io/docs/) as the target deployment environment

## Installation

1. Export a GitHub webpage URL of a ready-to-use webpage.

    Example:

    ```bash
    export GH_WEBPAGE_URL=https://github.com/kyma-project/examples/archive/master.zip
    ```

2. Apply a Bucket custom resource (CR):

    ```bash
    cat <<EOF | kubectl apply -f -
    apiVersion: assetstore.kyma-project.io/v1alpha2
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
    apiVersion: assetstore.kyma-project.io/v1alpha2
    kind: Asset
    metadata:
      name: webpage
      namespace: default
    spec:
      source:
        url: ${GH_WEBPAGE_URL}
        mode: package
        filter: /asset-store/webpage/
      bucketRef:
        name: pages
    EOF
    ```

4. Check the value of the **phase** field:

    ```bash
    kubectl get assets.assetstore.kyma-project.io webpage -o jsonpath='{.status.phase}'
    ```

    You should get a exacly that result:

    ```test
    Ready
    ```

    >**Note:** if state equals `Pending`, please wait few secounds and try again, but if state equals `Failed`, something went wrong and then you should check reason of filure using extracting the value of the **reason** field from the Asset CR:

    ```bash
    kubectl get assets.assetstore.kyma-project.io webpage -o jsonpath='{.status.reason}'
    ```

5. Export and merge the values of the **baseUrl** field and the path to the `index.html` file from the Asset CR, and then open it in default web browser:

    ```bash
    open $(kubectl get assets.assetstore.kyma-project.io webpage -o jsonpath='{.status.assetRef.baseUrl}{"/"}{.status.assetRef.files[?(@.name=="examples-master/asset-store/webpage/index.html" )].name}')
    ```

### Cleanup

1. Delete the Asset CR:

    ```bash
    kubectl delete assets.assetstore.kyma-project.io webpage
    ```

2. Delete the Bucket CR:

    ```bash
    kubectl delete buckets.assetstore.kyma-project.io pages
    ```

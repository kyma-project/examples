# Asset Store

## Overview

This example illustrates how to use the [Asset Store](https://kyma-project.io/docs/1.5/components/asset-store/) to store static webpages.

By default, [Minio](https://min.io/) stores all resources on a cluster, but it also allows you to use different cloud providers. Read the Asset Store [tutorials](https://kyma-project.io/docs/components/asset-store#tutorials-tutorials) for more information.

## Prerequisites

- [kubectl](https://kubernetes.io/docs/tasks/tools/install-kubectl/)
- [Kyma](https://kyma-project.io/docs/) as the target deployment environment

## Installation

1. Export a GitHub webpage URL of a ready-to-use webpage.

    Example:

    ```bash
    export GH_WEBPAGE_URL=https://github.com/pPrecel/simple-page-for-asset-store/archive/master.zip
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
      policy: readwrite
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
      bucketRef:
        name: pages
    EOF
    ```

4. Describe the Asset CR:

    ```bash
    kubectl describe assets.assetstore.kyma-project.io webpage
    ```

5. Find the **Asset Ref** field and merge **Base URL** with the filename of your `index.html`.

    Example:

    ```yaml
    Status:
      Asset Ref:
        Base URL:  https://minio.kyma.local/pages-1bjc0e7p0qdue/webpage
        Files:
          Name:             simple-page-for-asset-store-master/LICENSE
          Name:             simple-page-for-asset-store-master/README.md
          Name:             simple-page-for-asset-store-master/index.html
          Name:             simple-page-for-asset-store-master/jquery.js
          Name:             simple-page-for-asset-store-master/myscript.js
          Name:             simple-page-for-asset-store-master/style.css
    ```

    This is an example of the **Base URL** merged with the filename of the `index.html`:

    ```url
    https://minio.kyma.local/pages-1bjc0e7p0qdue/webpage/simple-page-for-asset-store-master/index.html
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

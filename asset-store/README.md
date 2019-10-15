# Asset Store

## Overview

This example illustrates how to use the [Asset Store](https://kyma-project.io/docs/1.5/components/asset-store/) to store simple webpages.

## Prerequisites

- [kubectl](https://kubernetes.io/docs/tasks/tools/install-kubectl/)
- [Kyma](https://kyma-project.io/docs/) as the target deployment environment

## Details

By default, [Minio](https://min.io/) stores all resources on a cluster, but it also allows you to use different cloud providers. Read the Asset Store [tutorials](https://kyma-project.io/docs/components/asset-store#tutorials-tutorials) for more information.

### Installation

1. Run the export GH_WEBPAGE_URL={value}, where value contain address to download webpage

    Example:

    ```bash
    export GH_WEBPAGE_URL=https://github.com/pPrecel/simple-page-for-asset-store/archive/master.zip
    ```

2. Apply bucket CR:

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

3. Apply asset CR:

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

### Testing

1. Describe asset CR:

    ```bash
    kubectl describe assets.assetstore.kyma-project.io webpage
    ```

2. Find "Asset Ref" field and merge "Base URL" with file name of your index.html

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

    In this case it should looks like that:

    ```url
    https://minio.kyma.local/pages-1bjc0e7p0qdue/webpage/simple-page-for-asset-store-master/index.html
    ```

### Creanup

1. Delete asset CR:

    ```bash
    kubectl delete assets.assetstore.kyma-project.io webpage
    ```

2. Delete bucket CR:

    ```bash
    kubectl delete buckets.assetstore.kyma-project.io pages
    ```

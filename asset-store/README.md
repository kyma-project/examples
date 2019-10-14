# Asset Store

## Overview

This example illustrates how to use [asset-store](https://kyma-project.io/docs/1.5/components/asset-store/) to store simple websides.

## Prerequisites

- [kubectl](https://kubernetes.io/docs/tasks/tools/install-kubectl/)
- [Kyma](https://kyma-project.io/docs/) as the target deployment environment
-  By default minio store all resources on current cluster, but it works fine also on cloud storage systems. Read [asset-store tutorials](https://kyma-project.io/docs/components/asset-store#tutorials-tutorials) for more informations.


## Details

### Installation

1. Run the export GH_WEBSIDE_URL={value}, where value contain address to download webside

    Example:

    ```bash
    export GH_WEBSIDE_URL=https://github.com/pPrecel/simple-page-for-asset-store/archive/master.zip
    ```

2. Apply bucket CR:

    ```
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

    ```
    cat <<EOF | kubectl apply -f -
    apiVersion: assetstore.kyma-project.io/v1alpha2
    kind: Asset
    metadata:
      name: webside
      namespace: default
    spec:
      bucketRef:
        name: pages
    source:
        url: ${GH_WEBSIDE_URL}
        mode: package
    EOF
    ```

### Creanup

1. Delete bucket CR:

    ```bash
    kubectl delete assetstore.kyma-project.io pages
    ```

2. Delete asset CR:

    ```bash
    kubectl delete assetstore.kyma-project.io webside
    ```

# Event Email Service

## Overview

This example illustrates how to use [asset-store](https://kyma-project.io/docs/1.5/components/asset-store/) to store simple websides on [ABS](https://azure.microsoft.com/en-us/services/storage/blobs/).

## Prerequisites

- [kubectl](https://kubernetes.io/docs/tasks/tools/install-kubectl/)
- [Azure Subscription](https://docs.microsoft.com/en-us/azure/billing/billing-add-change-azure-subscription-administrator)
- [Kyma](https://kyma-project.io/docs/) as the target deployment environment

## Details

### Installation

1. [Set Minio to the Azure Blob Storage Gateway mode](https://kyma-project.io/docs/1.6/components/asset-store/#tutorials-set-minio-to-the-azure-blob-storage-gateway-mode)

2. Run the export GH_WEBSIDE_URL={value}, where value contain address to download webside

    Example:

    ```bash
    export GH_WEBSIDE_URL=https://github.com/pPrecel/simple-page-for-asset-store/archive/master.zip
    ```

3. Apply bucket CR:

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

4. Apply asset CR:

    ```bash
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

### Testing

1. Describe asset CR:

    ```bash
    kubectl describe assets.assetstore.kyma-project.io webside
    ```

2. Find and coppy url addres from x/x/x field

    Example:

3. Find and coppy path to index.html file from section below:

    Example:

4. Merge this two addresses and put it to your favorite browser

### Creanup

1. Delete bucket CR:

    ```bash
    kubectl delete assetstore.kyma-project.io pages
    ```

2. Delete asset CR:

    ```bash
    kubectl delete assetstore.kyma-project.io webside
    ```

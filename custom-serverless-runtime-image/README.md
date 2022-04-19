# Custom serverless runtime image

## Overview

This example shows how to create own custom runtime for a serverless Function based on the Python runtime and `debian:bullseye-slim` base image to provide support for glibc.

## Prerequisites

- Docker as a build tool

## Build an example runtime

1. Export the following environments:

    ```bash
    export IMAGE_NAME=<image_name>
    export IMAGE_TAG=<image_tag>
    ```

2. Build and push the image:

    ```bash
    docker build -t "${IMAGE_NAME}/${IMAGE_TAG}" .
    docker push "${IMAGE_NAME}/${IMAGE_TAG}"
    ```

    > **Note:** You can use it to define your Functions in Kyma. To learn more, read [how to override runtime image](https://kyma-project.io/docs/kyma/main/05-technical-reference/svls-08-function-specification#override-runtime-image).

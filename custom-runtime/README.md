# Custom Runtime

## Overview

This example shows how to create own custom runtime for a serverless function based on python runtime and `debian:bullseye-slim` base image to provide support for glibc.

## Prerequisites

- Docker as a build tool

## Build an example runtime

1. Export following environments:

    ```bash
    export IMAGE_NAME=<image_name>
    export IMAGE_TAG=<image_tag>
    ```

2. Build and push the image:

    ```bash
    docker build -t "${IMAGE_NAME}/${IMAGE_TAG}" .
    docker push "${IMAGE_NAME}/${IMAGE_TAG}"
    ```

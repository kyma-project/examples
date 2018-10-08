#!/usr/bin/env bash
set -e
set -o pipefail

IMAGE_NAME=eu.gcr.io/kyma-project/example/tracing/order-front:0.0.1

echo -e "Start building docker image [ ${IMAGE_NAME} ]"

docker build --no-cache -t ${IMAGE_NAME} .

echo -e "Docker image [ ${IMAGE_NAME} ] has been built successfully ..."

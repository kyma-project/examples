#!/usr/bin/env bash

WEB_UI_DOCKER_HUB=eu.gcr.io/kyma-project/example
APP_NAME=call-ec-web-ui
APP_VERSION=0.0.1

docker build . -t $WEB_UI_DOCKER_HUB/$APP_NAME:$APP_VERSION

docker push $WEB_UI_DOCKER_HUB/$APP_NAME:$APP_VERSION

# Asynchronous communication between Functions

## Overview

This example provides a very simple scenario of asynchronous communication between two Functions, where: 
- The first Function accepts the incoming traffic via HTTP, sanitizes the payload, and publishes the content as an in-cluster event via [Kyma Eventing](https://kyma-project.io/docs/kyma/latest/01-overview/main-areas/eventing/).
- The second Function is a message receiver. It subscribes to the given event type and stores the payload.

This example also provides a template how a git project with Kyma Functions can be structured. Please refer to the Deploy section below.

## Prerequisites

* [Kyma CLI](https://github.com/kyma-project/cli)
* Kyma installed locally or on a cluster

## Deploy
### Deploy via Kyma CLI

Use your favorite IDE to change the logic in the handler.js files.
Run and test changes locally using Kyma CLI.
You can deploy to a test kyma runtime using Kyma CLI 

You can find all installation steps in the [Set asynchronous communication between Functions](https://github.com/kyma-project/kyma/blob/b783d9e6dffc47c0e3c31923aff62371b0a46779/docs/03-tutorials/00-serverless/svls-11-set-asynchronous-connection-of-functions.md) tutorial.


### Auto-deploy code changes
Changes pushed to the `handler.js` files should be automatically pulled by Kyma Serverless as both functions are of `git` type and reference this git repository as the source.

### Deploy via kubectl

Render Kubernetes manifests using the `make render` target. This outputs Kubernetes manifests to the k8s-resources folder.

Deploy to Kyma runtime manually using kubectl or `make deploy`.
There is also a github workflow included which you can use as a template to come up with own automated CI/CD.




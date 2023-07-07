# Asynchronous communication between Functions

## Overview

This example provides a very simple scenario of asynchronous communication between two Functions, where: 
- The first Function accepts the incoming traffic via HTTP, sanitizes the payload, and publishes the content as an in-cluster event via [Kyma Eventing](https://kyma-project.io/docs/kyma/latest/01-overview/eventing/).
- The second Function is a message receiver. It subscribes to the given event type and stores the payload.

This example also provides a template for a git project with Kyma Functions. Please refer to the Deploy section below.

## Prerequisites

* [Kyma CLI](https://github.com/kyma-project/cli)
* Kyma installed locally or on a cluster

## Deploy
### Deploy via Kyma CLI

You can deploy each Function separately using Kyma CLI by running `kyma apply function` in each of the Function's source folders.

You can find all installation steps in the [Set asynchronous communication between Functions](https://kyma-project.io/docs/kyma/latest/03-tutorials/00-serverless/svls-11-set-asynchronous-connection-of-functions/) tutorial.

### Deploy via kubectl

Deploy to Kyma runtime manually using `kubectl apply` or `make deploy` target.
There is also a [github workflow](.github/workflows/deploy.yml) included which you can use as a template to come up with own automated CI/CD.


### Auto-deploy code changes
Changes pushed to the `handler.js` files should be automatically pulled by Kyma Serverless as both Functions are of git type and reference this git repository as the source.



### Test the application

Send an HTTP request to the emitter Function

```bash
curl -H "Content-Type: application/cloudevents+json" -X POST  https://incoming.{your cluster domain} -d '{"foo":"bar"}'
Event sent%
```

Fetch the logs of the receiver Function to observe the incoming message.

```bash

> nodejs16-runtime@0.1.0 start
> node server.js

user code loaded in 0sec 0.649274ms
storing data...
{"foo":"bar"}
```


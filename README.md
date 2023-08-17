# Examples

## Overview

The examples project provides a central repository to showcase and illustrate features and concepts on Kyma.

> **NOTE:** The examples in this repository are provided as is, which means they may not always be up to date with the latest Kyma version. For current examples and tutorials, check our [documentation](https://kyma-project.io) page.

### What an example is

- An example is a small demo that illustrates a particular feature or concept of Kyma.
- An example refers to full, ready-to-use application development requiring no explanation.
- An example must be concise and clear.

### What an example is not

- An example cannot be a lecture or tutorial that guides the user through the topic with steps. Tutorials are part of the product documentation.
- An example is not a production-ready application or service. Do not use examples in a production environment.

### List of examples

The summary of the documentation in the `examples` repository lists all available examples organized by the feature or concept they showcase. This structure provides a quick overview and easy navigation.

| Example | Description | Technology |
|---|---|---|
| [Orders Service](orders-service/README.md) | Test the service and function that expose HTTP APIs to collect simple orders either in the in-memory storage or the Redis database. | Go, NodeJS, Redis |
| [HTTP DB Service](http-db-service/README.md) | Test the service that exposes an HTTP API to access a database on the cluster. | Go, MSSQL |
| [Gateway](gateway/README.md) | Expose APIs for functions or services.  | Kubeless |
| [Service Binding](service-binding/lambda/README.md) | Bind a Redis service to a lambda function. | Kubeless, Redis, NodeJS |
| [Event Email Service](event-email-service/README.md) | Send an automated email upon receiving an Event.  | NodeJS |
| [Tracing](tracing/README.md) | Configure tracing for a service in Kyma. | Go |
| [Custom serverless runtime image](custom-serverless-runtime-image/README.md) | Build custom runtime image. | Python |
| [Kiali](kiali/README.md) | Deploy Kiali to Kyma. | Kiali, Istio, Prometheus |
| [Prometheus](prometheus/README.md) | Deploy a Prometheus stack to monitor custom metrics. | Prometheus, Grafana |
| [Loki](loki/README.md) | Deploy Loki to store Pod logs. | Loki, Grafana |

## Installation

You can run all the examples locally or in a cluster. For instructions on installing and running an example, refer to the `README.md` document of the specific example, either by using the **List of examples** section or by navigating through the project structure.

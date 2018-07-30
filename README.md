# Examples
[![FOSSA Status](https://app.fossa.io/api/projects/git%2Bgithub.com%2Fkyma-project%2Fexamples.svg?type=shield)](https://app.fossa.io/projects/git%2Bgithub.com%2Fkyma-project%2Fexamples?ref=badge_shield)


## Overview

The examples project provides a central repository to showcase and illustrate features and concepts on Kyma.

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
| [HTTP DB service](http-db-service/README.md) | Service exposing an HTTP API to access a database on the cluster. | Go, MSSQL |
| [Event Service Subscription](event-subscription/service/README.md) | The example demonstrating publish and consume features of the Event Bus. | Go |
| [Event Lambda Subscription](event-subscription/lambda/README.md) | Create functions, trigger functions on Events, and bind to services.  | Kubeless |
| [Gateway](gateway/README.md) | Expose APIs for functions or services.  | Kubeless |
| [Service Binding](service-binding/lambda/README.md) | Bind a Redis service to a Lambda Function. | Kubeless, Redis, NodeJS |
| [Calling Enterprise Commerce](call-ec/README.md) | Call Enterprise Commerce in the context of the end user | Kubeless, NodeJS |
| [Alerting Rules](monitoring-alert-rules/README.md) | Configure alert rules in Kyma.  | Prometheus |
| [Custom Metrics in Kyma](monitoring-custom-metrics/README.md) | Expose Custom Metrics in Kyma.  | go, Prometheus |
| [Event Email Service](event-email-service/README.md) | Send an automated email when an event is received.  | NodeJS |

## Installation

You can run all the examples locally or in a cluster. For instructions on installing and running an example, refer to the `README.md` document of the specific example, either by using the [List of examples](#list-of-examples) or by navigating through the project structure.

## License
[![FOSSA Status](https://app.fossa.io/api/projects/git%2Bgithub.com%2Fkyma-project%2Fexamples.svg?type=large)](https://app.fossa.io/projects/git%2Bgithub.com%2Fkyma-project%2Fexamples?ref=badge_large)
# Examples

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

<!-- NOTE: The table with examples is also available in the "kyma" repository. Whenever you update the table, modify this document accordingly: https://kyma-project.io/docs/root/kyma#examples-kyma-features-and-concepts-in-practice . -->

| Example | Description | Technology |
|---|---|---|
| [Orders Service](orders-service/README.md) | Test the service and function that expose HTTP APIs to collect simple orders either in the in-memory storage or the Redis database. | Go, NodeJS, Redis |
| [HTTP DB Service](http-db-service/README.md) | Test the service that exposes an HTTP API to access a database on the cluster. | Go, MSSQL |
| [Gateway](gateway/README.md) | Expose APIs for functions or services.  | Kubeless |
| [Service Binding](service-binding/lambda/README.md) | Bind a Redis service to a lambda function. | Kubeless, Redis, NodeJS |
| [Alert Rules](monitoring-alert-rules/README.md) | Configure alert rules in Kyma.  | Prometheus |
| [Custom Metrics in Kyma](monitoring-custom-metrics/README.md) | Expose custom metrics in Kyma.  | Go, Prometheus |
| [Event Email Service](event-email-service/README.md) | Send an automated email upon receiving an Event.  | NodeJS |
| [Tracing](tracing/README.md) | Configure tracing for a service in Kyma. | Go |
| [Rafter](rafter/README.md) | Store static webpages using Rafter | HTML, CSS |

## Installation

You can run all the examples locally or in a cluster. For instructions on installing and running an example, refer to the `README.md` document of the specific example, either by using the **List of examples** section or by navigating through the project structure.

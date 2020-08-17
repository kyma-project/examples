# Orders Service

## Overview

This example demonstrates Kyma capabilities, such as HTTP endpoints that expose and bind a service to a database. The Orders Service is a sample application (microservice) written in Go. It can expose HTTP endpoints used to CR*D (create, read and delete) basic order JSON entities. The service can run with either an in-memory database that is enabled by default or an external, Redis database.

Additionally, a similar [Serverless](https://kyma-project.io/docs/components/serverless/) Function was created, with the ability to read all records and write a single one. Like the microservice, function can run with either an in-memory database or an Redis instance.

For more information about exposing by API Rule, binding Service Instance and binding event Triggers to microservice/Function, see [official Kyma's started guide](link do guida).

## Prerequisites

- Kyma 1.14 or higher.
- Kubernetes 1.16 or higher.
- [Kubectl](https://kubernetes.io/docs/reference/kubectl/kubectl/) 1.16 or higher.
- [Helm](https://helm.sh/) 3.0 or higher - not required.
- [Docker Compose](https://docs.docker.com/compose/) 3.0 or higher - not required.

## Installation

### By Kubectl

To install Orders Service in Kyma cluster, run:

```bash
kubectl create ns orders-service
kubectl apply -f ./deployment/orders-service.yaml
```

To install Serverless Function in Kyma cluster, run:

```bash
kubectl create ns orders-service
kubectl apply -f ./deployment/function.yaml
```

### By Helm

To install Orders Service in Kyma cluster, run:

```bash
helm install orders-service --namespace orders-service --create-namespace --timeout 60s --wait ./chart
```

Configuration for helm release is in [`values.yaml`](./chart/values.yaml) file.

### By Docker Compose

To build and run the Orders Service locally with Docker Compose, run:

```bash
make docker-compose
```

## Cleanup

### By Kubectl

Run the following command to completely remove the example (Orders Service or Serverless Function) and all its resources from the cluster:

```bash
kubectl delete all -l app=orders-service -n orders-service
kubectl delete ns orders-service
```

### By Helm

Run the following command to completely remove the helm release with example and all its resources from the cluster:

```bash
helm delete orders-service -n orders-service
kubectl delete ns orders-service
```

## Configuration

To configure the service/Function, override the default values of these environment variables:

| Environment variable | Description                                                                   | Required   | Default value |
| ---------------------- | ----------------------------------------------------------------------------- | ------ | ------------- |
| **APP_PORT**       | Specifies the port of running service. Function doesn't use this variable. | NO | `8080`           |
| **APP_REDIS_PREFIX**       | Specifies the prefix for all related to Redis environment variables. See the variables below. | NO | `REDIS_`           |
| **{APP_REDIS_PREFIX}HOST**       | Specifies the host of Redis instance.                       | NO | `nil`            |
| **{APP_REDIS_PREFIX}PORT**       | Specifies the port of Redis instance.                       | NO | `nil`            |
| **{APP_REDIS_PREFIX}REDIS_PASSWORD**       | Specifies the password of Redis instance to authorization.                       | NO | `nil`            |

Example:

```bash
export APP_REDIS_PREFIX="R_"
export R_HOST="abc.com"
export R_PORT="8080"
export R_REDIS_PASSWORD="xyz"
```

## Testing

### Microservice

For ...
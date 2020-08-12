# Orders Service

## Overview

This example demonstrates Kyma capabilities, such as HTTP endpoints that expose and bind a service to database. The service in this example exposes HTTP endpoints simple CRUD (Create, Read, Update, Delete) for basic order JSON entities, as described in the [service's OpenAPI specification](docs/openapi.yaml). The service can run with either an in-memory database or an Redis instance. By default, the in-memory database is enabled. The service in this example uses [Go](http://golang.org).

## Prerequisites

- Kyma 1.14 or higher.
- Kubernetes 1.16 or higher.
- [Helm](https://helm.sh/) 3.0 or higher.
- [Docker Compose](https://docs.docker.com/compose/) 3.0 or higher - not required.

## Installation

To install service in Kyma cluster, run:

```bash

```

Otherwise to build and run the service locally with Docker, run:

```bash
make docker-compose
```

To run the service with an MSSQL database, use the **DbType=mssql** environment variable in the application. To configure the connection to the database, set the environment variables for the values defined in the `config/db.go` file.

The `deployment` folder contains `.yaml` descriptors used for the deployment of the service to Kyma.

Run the following commands to deploy the published service to Kyma:

1. Export your Namespace as variable by replacing the `{namespace}` placeholder in the following command and running it:

    ```bash
    export KYMA_EXAMPLE_NS="{namespace}"
    ```

2. Deploy the service:

    ```bash
    kubectl create namespace $KYMA_EXAMPLE_NS
    kubectl apply -f deployment/deployment.yaml -n $KYMA_EXAMPLE_NS
    ```

### Cleanup

Run the following command to completely remove the example and all its resources from the cluster:

```bash
kubectl delete all -l example=http-db-service -n $KYMA_EXAMPLE_NS
```
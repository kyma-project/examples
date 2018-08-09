# HTTP DB Service

## Overview

This example demonstrates Kyma capabilities, such as HTTP endpoints that expose and bind a service to a database. The service in this example exposes HTTP endpoints used to create and read basic order JSON entities, as described in the [service's API descriptor](docs/api/api.yaml). The service can run with either an in-memory database or an MSSQL instance. By default, the in-memory database is enabled. The service in this example uses the [Go](http://golang.org) language.

## Prerequisites

- A [Docker](https://docs.docker.com/install) installation.
- Kyma as the target deployment environment.
- An MSSQL database for the service's database functionality. You can also use Azure MSSQL that you can provision using the Kyma Open Service Broker API.
- An environment to which to deploy the example.

## Installation

Use these commands to build and run the service with Docker:
```
./build.sh
docker run -it --rm -p 8017:8017 http-db-service:latest
```

To run the service with an MSSQL database, use the **DbType=mssql** environment variable in the application. To configure the connection to the database, set the environment variables for the values defined in the `config/db.go` file.

The `deployment` folder contains `.yaml` descriptors used for the deployment of the service to Kyma.

Run the following commands to deploy the published service to Kyma:

1. Export your environment as variable by replacing the `<environment>` placeholder in the following command and running it:

    ```bash
    export KYMA_EXAMPLE_ENV="<environment>"
    ```
2. Deploy the service
    ```bash
    kubectl apply -f deployment/deployment.yaml -n $KYMA_EXAMPLE_ENV
    ```

### MSSQL Database tests

To run the unit tests on a real MSSQL database, run the following commands on the root of the example:

```bash
docker run -ti -e 'ACCEPT_EULA=Y' -e 'SA_PASSWORD=Password!123' -p 1433:1433 -d microsoft/mssql-server-linux:2017-latest
```

This starts a MSSQL database in a container.

```bash
username=sa password='Password!123' database=master tablename='test_orders' host=localhost port=1433 dbtype=mssql go test ./... -v
```

This runs the specific unit tests for MSSQL databases with all the environment information to connect to the previously started MSSQL database.


### Deployment smoke test

Run the CURL command to verify that the service is running:

```bash
curl --request GET --url http://<host>:<port>/orders
```

### Cleanup

Run the following command to completely remove the example and all its resources from the cluster:

```bash
kubectl delete all -l example=http-db-service -n $KYMA_EXAMPLE_ENV
```
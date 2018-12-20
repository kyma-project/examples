# Event Email Service

## Overview

This example illustrates how to write a service in `NodeJS` that listens for Events and sends an automated email to the address in the Event payload.

## Prerequisites

- A [Docker](https://docs.docker.com/install) installation
- Kyma as the target deployment environment
- An Environment to which to deploy the example

## Installation

### Local installation

1. Build the Docker image:
    ```bash
    docker build . -t event-email-service:latest
    ```

2. Run the recently built image:
    ```bash
    docker run -p 3000:3000 event-email-service:latest
    ```

### Cluster installation

1. Export your Environment as a variable by replacing the **{environment}** placeholder in the following command and running it:
    ```bash
    export KYMA_EXAMPLE_ENV="{environment}"
    ```

2. Deploy the service:
    ```bash
    kubectl apply -f deployment -n $KYMA_EXAMPLE_ENV
    ```

3. Expose the service endpoint:
    ```bash
    kubectl port-forward -n $KYMA_EXAMPLE_ENV $(kubectl get pod -n $KYMA_EXAMPLE_ENV -l example=event-email-service | grep event-email-service | awk '{print $1}') 3000
    ```

### Test the service

To test the service, simulate an Event using cURL:

```bash
curl -H "Content-Type: application/json" -d '{"event":{"customer":{"customerID": "1234", "uid": "rick.sanchez@mail.com"}}}' http://localhost:3000/v1/events/register
```

After sending the Event, you should see a log entry either in your terminal (if running locally) or in the Pod's logs (if running on Kyma) confirming the Event reception.

### Cleanup

Clean all deployed example resources from Kyma with the following command:

```bash
kubectl delete all -l example=event-email-service -n $KYMA_EXAMPLE_ENV
```

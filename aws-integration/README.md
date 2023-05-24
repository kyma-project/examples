# Connect AWS X-Ray and CloudWatch to kyma

## Overview 

Sometimes you may want to transfer your logs, metrics and traces to other systems as well. Kyma provides a possibility to integrate it with some of the AWS components by deploying a custom OTEL Collector which would collect traces, logs, and metrics from kyma system and forward them to the AWS. In order to do this, we recommend you to use [AWS distribution](https://aws-otel.github.io) for OTEL Collector. 
In this tutorial we will show you how to deploy OTEL Collector in order to connect Kyma observability 

## Prerequisistes 

- Kyma as the target deployment environment
- AWS account with permissions to create new users and security policies
- Helm 3.x if you want to deploy open-telemetry sample application

## Installation

### Preparation

1. Export your Namespace as a variable. Replace `{NAMESPACE}` placeholder in the following command and run it:

    ```bash
    export KYMA_NS="{NAMESPACE}"
    ```
1. If you don't have created a Namespace yet, do it now:
    ```bash
    kubectl create namespace $KYMA_NS
    ```
1. Export the AWS region you want to use as variable. Replace `{NAMESPACE}` placeholder in the following command and run it:
    ```bash
    export AWS_REGION="{NAMESPACE}"
    ```

### Create AWS IAM User

As a first step, we need to create an IAM User assign to it specific IAM policies through which our OTEL Collector will be communicating with AWS. 

Firstly, we need to create IAM policy for pushing our metrics:
1. Go to the AWS searchbar, and search for IAM service
1. Go to the `Policies` section. and click `Create policy`
1. Remove the `Popular services` flag, and select `CloudWatch` service. 
1. Now, select `GetMetricData`, `PutMetricData`, `ListMetrics` actions and click `Next`
1. Enter the policy name and click `Create policy`

Now, we should create the policy for CloudWatch Logs:
1. Repeat the first two steps you did while creating the policy for metrics
1. Remove the `Popular services` flag, and select `CloudWatch Logs` service.
1. Now, select `CreateLogGroup`, `CreateLogStream`, `PutLogEvents` actions and click `Next`
1. Specify resource ARN for selected actions
1. Enter the policy name and click `Create policy` 

After creating the IAM Policies, we can finally create an IAM User:
1. Go to the `Users` section and click `Add user`
1. Enter the user name and click next
1. Click to the `Attach policies directly`
1. Select two new policies you added on previous steps as well as the `AWSXrayWriteOnlyAccess`
1. Click `Next` and then `Create User`
1. Open the new user
1. Go to the `Security credentials` tab and click `Create access key`
1. Select `Application running outside AWS` and then click `Next`
1. Describe the purpose of this access key and click `Create access key`
1. Now copy and save `Access key` and `Secret access key`
1. Encode the keys into base64 encoding

### Create a secret with AWS Credentials

In order to connect OTEL Collector to AWS we need to define security credentials in the kyma system. 

1. In the [values.yaml](./aws-secret/values.yaml), replace the `{KEY_ID_BASE64}` and `{KEY_SECRET_BASE64}` to your encoded access keys 
2. Now, create the secret by using 
    ```bash
    kubectl apply -f ./aws-secret/values.yaml
    ```

### Deploy an OTEL Collector

After creating a secret and configuering AWS, we can finally deploy an Otel Collector itself.

1. Deploy an OTEL Collector by calling 
    ```bash
    kubectl apply -f ./aws-otel-collector/values.yaml
    ```

### Create pipelines

After deploying OTEL Collector itself, you should deploy logpipeline, tracepipeline, and metricpipeline. 

1. Deploy a logpipeline by calling 
    ```bash
    kubectl apply -f ./pipelines/logpipeline.yaml
    ```
1. Deploy a tracepipeline by calling 
    ```bash
    kubectl apply -f ./pipelines/tracepipeline.yaml
    ```
1. Deploy a metricpipeline by calling 
    ```bash
    kubectl apply -f ./pipelines/metricpipeline.yaml
    ```

## Verifying the results by deploying sample apps

In order to verify the results of CloudWatch and X-Ray we will deploy sample applications for each service accordingly.

### Verifying X-Ray trace arrival 

In order to deploy sample app which generates traces:
1. Deploy an app using 
    ```bash
    kubectl apply -n ${KYMA_NS} -f ./trace-sample-app/values.yaml
    ```
1. Port-forward an application in order to be able to access it by calling 
    ```bash
    kubectl -n ${KYMA_NS} port-forward svc/sample-app 4567
    ```
1. Make some requests to the application like `localhost:4567` or `localhost:4567/outgoing-http-call`
1. Go to the `AWS X-Ray` tab, and check out the `Traces` section

### Verifying CloudWatch logs and metrics arrival

In order to deploy sample app which generates logs and metrics, we can use OpenTelemetry example. 
1. Execute 
    ```bash
    helm upgrade --version 0.22.2 --install --create-namespace -n ${KYMA_NS} otel-collector open-telemetry/opentelemetry-demo -f ./open-telemetry-sample-app/values.yaml
    ```
1. You can port-forward this application in order to access it via browser by calling 
    ```bash
    kubectl -n ${KYMA_NS} port-forward svc/otel-collector-frontend 8080
    ```
    and going to the `localhost:8080`
1. Now, you can go to `AWS CloudWatch` and look at the main dashboard. There should be logs and metrics arriving

### Creating the dashboard to observe incoming metrics and logs

In order to create a dashboard, you should:
1. Go to the `Dashboards` section
1. Click `Create dashboard` and enter the name
1. Select the widget type and click `Next`
1. Select what you want to observe, either metrics or logs, and click `Next`
1. Decide on which metrics and logs to include and click `Create widget`
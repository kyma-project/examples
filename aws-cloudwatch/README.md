# Integrate Kyma with AWS CloudWatch

## Overview 

The Kyma Telemetry module supports you in integrating with observability backends in a convenient way. This example outlines how to integrate with [AWS CloudWatch](https://aws.amazon.com/cloudwatch) as a backend. As CloudWatch is not supporting OTLP ingestion natively, it will require to deploy the [AWS Distro for OpenTelemetry](https://aws-otel.github.io) additionally. 

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

### Create a secret with AWS Credentials

In order to connect OTEL Collector to AWS we need to define security credentials in the kyma system. 

1. In the [values.yaml](./aws-secret/values.yaml), replace the `{ACCESS_KEY}` and `{SECRET_ACCESS_KEY}` to your access keys, and `{AWS_REGION}` with the AWS region you want to use
2. Now, create the secret by using 
    ```bash
    kubectl apply -f ./aws-secret/values.yaml
    ```

### Deploy an OTEL Collector

After creating a secret and configuering AWS, we can finally deploy an Otel Collector itself.

1. Deploy an OTEL Collector by calling 
    ```bash
    kubectl -n $KYMA_NS apply -f ./aws-otel-collector/values.yaml
    ```

### Create pipelines

After deploying OTEL Collector itself, you should deploy logpipeline, tracepipeline, and metricpipeline. 

1. Deploy a logpipeline by calling 
    ```bash
    kubectl apply -f ./pipelines/logpipeline.yaml
    ```
1. Replace `{NAMESPACE}` and deploy a tracepipeline by calling 
    ```bash
    kubectl apply -f ./pipelines/tracepipeline.yaml
    ```
1. Replace `{NAMESPACE}` and deploy a metricpipeline by calling 
    ```bash
    kubectl apply -f ./pipelines/metricpipeline.yaml
    ```

## Verifying the results by deploying sample apps

In order to verify the results of CloudWatch and X-Ray we will deploy sample applications for each service accordingly.

### Verifying CloudWatch traces, logs, and metrics arrival 

In order to deploy sample app which generates traces that we took from [aws-otel tutorial](https://docs.aws.amazon.com/eks/latest/userguide/sample-app.html):
1. Deploy traffic generator app
    ```bash
    kubectl apply -n ${KYMA_NS} -f ./trace-sample-app/traffic-generator.yaml
    ```
1. Deploy an app using 
    ```bash
    kubectl apply -n ${KYMA_NS} -f ./trace-sample-app/sample-app.yaml
    ```
1. Port-forward an application in order to be able to access it by calling 
    ```bash
    kubectl -n ${KYMA_NS} port-forward svc/sample-app 4567
    ```
1. Make some requests to the application like `localhost:4567` or `localhost:4567/outgoing-http-call`
1. Go to the `AWS X-Ray` tab, and check out the `Traces` section
1. To verify the logs, you can go to `AWS CloudWatch`, then open the `Log groups` and select your cluster. Now, you can open `aws-integration.sample-app-*` and check out the logs of your application.
1. To verify metrics, you can go to the `All metrics`, and open the `aws-integration/otel-collector`

### Creating the dashboard to observe incoming metrics and logs

In order to create a dashboard, you should:
1. Go to the `Dashboards` section
1. Click `Create dashboard` and enter the name
1. Select the widget type and click `Next`
1. Select what you want to observe, either metrics or logs, and click `Next`
1. Decide on which metrics and logs to include and click `Create widget`
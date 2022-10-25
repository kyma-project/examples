# Installing a custom Loki stack in Kyma

## Overview

The Kyma Loki component brings limited configuration options in contrast to the upstream [loki-stack](https://github.com/grafana/helm-charts/tree/main/charts/loki-stack) chart. Modifications might be reset at next upgrade cycle.

An alternative can be a parallel installation of the upstream chart offering all customization options. The following instructions outline how to achieve such installation in co-existence to the Kyma stack.

**CAUTION:** This example uses the Grafana Loki version, which is distributed under AGPL-3.0 only and might not be free of charge for commercial usage.

## Prerequisites

- Kyma as the target deployment environment.
- Kubectl version 1.22.x or higher
- Helm 3.x

## Installation

### Preparation

1. Export your Namespace as a variable. Replace the `{namespace}` placeholder in the following command and run it:

    ```bash
    export KYMA_LOKI_EXAMPLE_NS="{namespace}"
    ```

2. Export the Helm release name that you want to use. It can be any name, but be aware that all resources in the cluster will be prefixed with that name. Replace the `{release-name}` placeholder in the following command and run it:
    ```bash
    export HELM_RELEASE_NAME="{release-name}"
    ```

3. Update your Helm installation with the required Helm repository:

    ```bash
    helm repo add grafana https://grafana.github.io/helm-charts
    helm repo update
    ```

### Install the Loki stack

Run the Helm upgrade command, which installs the chart if not present yet.
 ```bash
helm upgrade --install --create-namespace -n ${KYMA_LOKI_EXAMPLE_NS} ${HELM_RELEASE_NAME} grafana/loki-stack -f https://raw.githubusercontent.com/kyma-project/examples/main/loki/loki-values.yaml --set promtail.enabled=false --set grafana.enabled=false
```

You can either use the [loki-values.yaml](./loki-values.yaml) provided in this `loki` folder, which contains customized settings deviating from the default settings, or create your own `values.yaml` file.


### Verify the installation

Check that the `loki` Pod has been created in the Namespace and is in the `Running` state:

```bash
kubectl -n ${KYMA_LOKI_EXAMPLE_NS} get pod ${HELM_RELEASE_NAME}-0
```

### Activate log shipment using a LogPipeline

1. Download the [logpipeline](https://raw.githubusercontent.com/kyma-project/examples/main/loki/logpipeline-custom.yaml) and replace the `{release-name}` and `{namespace}` placeholder.

2. Apply the modified LogPipeline:


   ```bash
   kubectl apply -f logpipeline-custom.yaml
When the status of the applied LogPipeline resource turned into `Running`, the underlying Fluentbit is reconfigured and log shipment to your Loki instance should be active
> **NOTE:** The used output plugin configuration is using all labels of a Pod to label the Loki log streams. That segregation of the log streams might be not optimal performance wise. Follow [Loki's labelling best practices](https://grafana.com/docs/loki/latest/best-practices/) for a tailormade setup fitting to your workload configuration.

### Verify the setup by accessing logs via the Loki API

1. To access the Loki API, use kubectl port forwarding. Run:
   ```bash
   kubectl -n ${KYMA_LOKI_EXAMPLE_NS} port-forward svc/${HELM_RELEASE_NAME} 3100

### Alternative installation options

The used Helm chart provides configuration for accessing logs with Grafana, which is disabled by default. 
To enable Grafana, run:

```bash
helm upgrade --install --create-namespace -n ${KYMA_LOKI_EXAMPLE_NS} ${HELM_RELEASE_NAME} grafana/loki-stack -f https://raw.githubusercontent.com/kyma-project/examples/main/loki/loki-values.yaml -f https://raw.githubusercontent.com/kyma-project/examples/main/loki/grafana-values.yaml --set grafana.adminPassword=myPwd --set promtail.enabled=false
```

To access Grafana, use kubectl port forwarding. Run:
```bash
kubectl -n ${KYMA_LOKI_EXAMPLE_NS} port-forward svc/${HELM_RELEASE_NAME}-grafana 3000:80
```
Open Grafana in your browser under http://localhost:3000 and log in with `admin` as password `myPwd`. Password can be easly changed by installation parameter `grafana.adminPassword=myPwd` from command above. 

The proposed approach is based Kyma's LogPipeline API which will configure a managed Fluentbit accordingly. Loki itself promotes an own log collector called `promtail` which could be used alternatively. A ready to use setup can be enabled easily by using following Helm command:
```bash
helm upgrade --install --create-namespace -n ${KYMA_LOKI_EXAMPLE_NS} ${HELM_RELEASE_NAME} grafana/loki-stack -f https://raw.githubusercontent.com/kyma-project/examples/main/loki/loki-values.yaml -f https://raw.githubusercontent.com/kyma-project/examples/main/loki/promtail-values.yaml --set grafana.enabled=false
```

> **NOTE:** Following the described instructions will install Loki in a very lightweight setup which will not fullfill production-grade criterias. Please consider to use a scalable setup based on an object storage backend as described in this [guide](https://grafana.com/docs/loki/latest/installation/simple-scalable-helm/).


### Cleanup

1. To remove the installation from cluster, run:

   ```bash
   helm delete -n ${KYMA_LOKI_EXAMPLE_NS} ${HELM_RELEASE_NAME}
   ```

2. To remove deployed LogPipeline instance from cluster, run:
   
   ```bash
   kubectl delete -f https://raw.githubusercontent.com/kyma-project/examples/main/loki/logpipeline-custom.yaml
   ```

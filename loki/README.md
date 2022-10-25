# Installing a custom Loki stack in Kyma

## Overview

The Kyma Loki component brings limited configuration options in contrast to the upstream [loki-stack](https://github.com/grafana/helm-charts/tree/main/charts/loki-stack) chart. Modifications might be reset at next upgrade cycle.

An alternative can be a parallel installation of the upstream chart offering all customization options. The following instructions outline how to achieve such installation in co-existence to the Kyma stack.

**CAUTION:** This example uses the Grafana Loki version, which is distributed under AGPL-3.0 only and might not be free of charge for commercial usage.
> **CAUTION:** These instructions install Loki in a lightweight setup that does not fulfil production-grade qualities. Consider using a scalable setup based on an object storage backend instead (see [Simple scalable deployment of Grafana Loki with Helm](https://grafana.com/docs/loki/latest/installation/simple-scalable-helm/)).

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

## Alternative installation options

## Installation together with Grafana
The used Helm chart supports the deployment of Grafana as well, but is disabled by default. As Grafana provides a very good Loki integration, you might want to install it as well.

To deploy Grafana alongside to Loki and having Loki pre-configured as a datasource, run the following command instead of the original command from the Installation section:

```bash
helm upgrade --install --create-namespace -n ${KYMA_LOKI_EXAMPLE_NS} ${HELM_RELEASE_NAME} grafana/loki-stack -f https://raw.githubusercontent.com/kyma-project/examples/main/loki/loki-values.yaml -f https://raw.githubusercontent.com/kyma-project/examples/main/loki/grafana-values.yaml --set grafana.adminPassword=myPwd

### Installation based on Promtail
Promtail is the recommended log collector to be used for feeding logs to Loki. The instructions provided her are using Kyma's LogPipeline feature based on Fluentbit.
If you prefer to use promtail, you can enable the deployment of it via the Helm chart as well.
The proposed approach is based Kyma's LogPipeline API, which configures a managed Fluent Bit accordingly. Loki itself promotes its own log collector called `promtail`, which you can use alternatively. You can enable a ready-to-use setup with the following Helm command instead of using the one outlined in the installation instructions above:
```bash
helm upgrade --install --create-namespace -n ${KYMA_LOKI_EXAMPLE_NS} ${HELM_RELEASE_NAME} grafana/loki-stack -f https://raw.githubusercontent.com/kyma-project/examples/main/loki/loki-values.yaml -f https://raw.githubusercontent.com/kyma-project/examples/main/loki/promtail-values.yaml
```a



### Cleanup

1. To remove the installation from cluster, run:

   ```bash
   helm delete -n ${KYMA_LOKI_EXAMPLE_NS} ${HELM_RELEASE_NAME}
   ```

2. To remove deployed LogPipeline instance from cluster, run:
   
   ```bash
   kubectl delete -f https://raw.githubusercontent.com/kyma-project/examples/main/loki/logpipeline-custom.yaml
   ```

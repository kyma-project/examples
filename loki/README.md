# Installing a custom Loki stack in Kyma

## Overview

In contrast to the upstream [loki-stack chart](https://github.com/grafana/helm-charts/tree/main/charts/loki-stack), Kyma's Loki component offers limited configuration options. Furthermore, any modifications you make might be reset at next upgrade cycle.

To get all the customization options, follow the instructions to set up a parallel installation of the upstream chart, co-existing with the Kyma stack.

>**CAUTION:** This example uses the Grafana Loki version, which is distributed under AGPL-3.0 only and might not be free of charge for commercial usage.

>**NOTE:** These instructions install Loki in a lightweight setup that does not fulfil production-grade qualities. Consider using a scalable setup based on an object storage backend instead (see [Simple scalable deployment of Grafana Loki with Helm](https://grafana.com/docs/loki/latest/installation/simple-scalable-helm/)).

## Prerequisites

- Kyma as the target deployment environment
- Kubectl version 1.22.x or higher
- Helm 3.x

## Installation

### Preparation

1. Export your Namespace as a variable. Replace the `{NAMESPACE}` placeholder in the following command and run it:

    ```bash
    export KYMA_LOKI_EXAMPLE_NS="{NAMESPACE}"
    ```

2. Export the Helm release name that you want to use. It can be any name, but be aware that all resources in the cluster will be prefixed with that name. Replace the `{HELM_RELEASE_NAME}` placeholder in the following command and run it:

    ```bash
    export HELM_RELEASE_NAME="{HELM_RELEASE_NAME}"
    ```

3. Update your Helm installation with the required Helm repository:

    ```bash
    helm repo add grafana https://grafana.github.io/helm-charts
    helm repo update
    ```

### Install the Loki stack

You install the Loki stack with a Helm upgrade command, which installs the chart if not present yet.

You can either choose an installation based on [Promtail](https://grafana.com/docs/loki/latest/clients/promtail/), which is the log collector recommended by Loki and provides a ready-to-use setup.
Alternatively, you can use Kyma's LogPipeline feature based on Fluent Bit.

The Helm chart supports the deployment of Grafana as well, but it's disabled by default. Because Grafana provides a very good Loki integration, you might want to install it as well.

In any case, you can either use the [loki-values.yaml](./loki-values.yaml) provided in this `loki` folder, which contains customized settings deviating from the default settings, or create your own `values.yaml` file.

<div tabs name="default-settings" group="configuration">
  <details>
  <summary label="promtail-installation">
  Installation with Promtail
  </summary>

To install the Loki stack based on Promtail, run:

```bash
helm upgrade --install --create-namespace -n ${KYMA_LOKI_EXAMPLE_NS} ${HELM_RELEASE_NAME} grafana/loki-stack -f https://raw.githubusercontent.com/kyma-project/examples/main/loki/loki-values.yaml
```
  </details>
  <details>
  <summary label="fluent-bit-installation">
  Installation with Kyma's LogPipeline
  </summary>

>**CAUTION:** This setup uses an unsupported output plugin for the LogPipline. Support for this might be removed in future.

1. To install the Loki stack with Kyma's LogPipeline feature based on Fluent Bit, run:

   ```bash
   helm upgrade --install --create-namespace -n ${KYMA_LOKI_EXAMPLE_NS} ${HELM_RELEASE_NAME} grafana/loki-stack -f https://raw.githubusercontent.com/kyma-project/examples/main/loki/loki-values.yaml --set promtail.enabled=false
   ```

2. Download the [LogPipeline](logpipeline-custom.yaml) and replace the `{HELM_RELEASE_NAME}` and `{NAMESPACE}` placeholder.

3. Apply the modified LogPipeline:

   ```bash
   kubectl apply -f logpipeline-custom.yaml
   ```

When the status of the applied LogPipeline resource turned into `Running`, the underlying Fluent Bit is reconfigured and log shipment to your Loki instance is active.

> **TIP:** The used output plugin configuration uses all labels of a Pod to label the Loki log streams. Performancewise, such segregation of the log streams might be not optimal. Follow [Loki's labelling best practices](https://grafana.com/docs/loki/latest/best-practices/) for a tailormade setup that fits your workload configuration.

  </details>
  <details>
  <summary label="installation-with-grafana">
  Installation with Grafana
  </summary>

  The used Helm chart supports the deployment of Grafana as well, but is disabled by default. Because Grafana provides a very good Loki integration, you might want to install it as well.

  1. To deploy Grafana alongside Loki, with Loki pre-configured as a datasource, run:

     ```bash
     helm upgrade --install --create-namespace -n ${KYMA_LOKI_EXAMPLE_NS} ${HELM_RELEASE_NAME} grafana/loki-stack -f https://raw.githubusercontent.com/kyma-project/examples/main/loki/loki-values.yaml -f https://raw.githubusercontent.com/kyma-project/examples/main/loki/grafana-values.yaml --set grafana.adminPassword=myPwd
     ```

  2. To access the Grafana UI with kubectl port forwarding, run:

     ```bash
     kubectl -n ${KYMA_LOKI_EXAMPLE_NS} port-forward svc/${HELM_RELEASE_NAME}-grafana 3000:80
     ```

  3. Open Grafana in your browser under `http://localhost:3000` and log in with user admin and the password taken from the previous Helm command.
  
  </details>
</div>

### Verify the installation

Check that the `loki` Pod has been created in the Namespace and is in the `Running` state:

```bash
kubectl -n ${KYMA_LOKI_EXAMPLE_NS} get pod -l app=loki,release=${HELM_RELEASE_NAME}
```

### Verify the setup by accessing logs using the Loki API

1. To access the Loki API, use kubectl port forwarding. Run:

   ```bash
   kubectl -n ${KYMA_LOKI_EXAMPLE_NS} port-forward svc/$(kubectl  get svc -n ${KYMA_LOKI_EXAMPLE_NS} -l app=loki,release=${HELM_RELEASE_NAME},variant=headless -ojsonpath='{.items[0].metadata.name}') 3100
   ```

1. Loki queries need a query parameter `time`, provided in nano seconds. To get the current nano seconds in Linux or MacOS, run:

   ```bash
   date +%s
   ```

1. To get the latest logs from Loki, replace the `{NANOSECONDS}` placeholder with the result of the previous command, and run:

   ```bash
   curl -G -s  "http://localhost:3100/loki/api/v1/query" \
     --data-urlencode \
     'query={job="fluentbit"}' \
     --data-urlencode \
     'time={NANOSECONDS}'
   ```

## Cleanup

1. To remove the installation from the cluster, run:

   ```bash
   helm delete -n ${KYMA_LOKI_EXAMPLE_NS} ${HELM_RELEASE_NAME}
   ```

2. To remove the deployed LogPipeline instance from cluster, run:

   ```bash
   kubectl delete -f https://raw.githubusercontent.com/kyma-project/examples/main/loki/logpipeline-custom.yaml
   ```

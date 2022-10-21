# Installing a custom loki-stack in Kyma

## Overview

The Kyma Loki component brings limited configuration options in contrast to the upstream [loki-stack]('https://github.com/grafana/helm-charts/tree/main/charts/loki-stack') chart. Modifications might be reset at next upgrade cycle.

An alternative can be a parallel installation of the upstream chart offering all customization options. The following instructions outline how to achieve such installation in co-existence to the Kyma stack.

**CAUTION:** This example will use the Grafana Loki version which is distributed under AGPL-3.0 only and might not be free of charge for commercial usage.

## Prerequisites

- Kyma as the target deployment environment.
- Kubectl > 1.22.x
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

### Install the loki-stack

Run the Helm upgrade command, which installs the chart if not present yet.
 ```bash
helm upgrade --install -n ${KYMA_LOKI_EXAMPLE_NS} ${HELM_RELEASE_NAME} grafana/loki-stack -f https://raw.githubusercontent.com/kyma-project/examples/main/loki/values.yaml
```

You can either use the [values.yaml]('values.yaml') provided in this loki folder, which contains customized settings deviating from the default settings, or create your own values.yaml file.


### Verify the installation

Check that the `loki` Pod have been created in the `{namespace}` and are in the `Running` state:

```bash
kubectl -n ${KYMA_LOKI_EXAMPLE_NS} get pod ${HELM_RELEASE_NAME}-0
```

### Install LogPipeline with custom plugin

> **NOTE:** Before applying following command, `{release-name}` and `{namespace}` placeholder should be replaced.  

Run the kubectl apply command to deploy custom LogPipeline to emmit logs to the deployed loki instance.

Command below will deploy a LogPipeline with `fluentbit` and loki plugin

 ```bash
kubectl apply -f  https://raw.githubusercontent.com/kyma-project/examples/main/loki/logpipeline-custom.yaml
```

### Accessing Loki and Logs

To access Loki, use kubectl port forwarding run:
```bash
kubectl -n ${KYMA_LOKI_EXAMPLE_NS} portforward svc/${HELM_RELEASE_NAME} 3100
```

Now we can access Loki over `localhost` and query collected logs.

Loki queries need a query parameter `time` in nano second resolution, to get current nano seconds in Linux or MacOS, run following command.

```bash
date +%s
```

Replace `{nanoseconds}` placeholder with result of command above and run following command to get latest logs from Loki

```bash
curl -G -s  "http://localhost:3100/loki/api/v1/query" \
  --data-urlencode \
  'query={job="fluentbit"}'
  --data-urlencode \
  'time={nanoseconds}'
```

### Cleanup

1. To remove the installation from cluster, run:

   ```bash
   helm delete -n ${KYMA_LOKI_EXAMPLE_NS} ${HELM_RELEASE_NAME}
   ```

2. To remove deployed LogPipeline instance from cluster, run:
   
   ```bash
   kubectl delete -f https://raw.githubusercontent.com/kyma-project/examples/main/loki/logpipeline-custom.yaml
   ```


### Additional

The used helm chart provides configuration for accessing logs via Grafana, which is disabled by default. To enable Grafana add following configuration parameters to the used [values.yaml]('values.yaml').

```yaml
grafana:
  enabled: true
```

An alternative to LogPipeline to collect logs and push to the loki, Loki Helm chart deliver log collector `promtail`, promtail disabled in this example, to enable it change [values.yaml]('values.yaml') as
```yaml
promtail:
  enabled: true
```

This example demonstrate basic installation and capabilities of loki with Kyma, for more stable and production grade setup please consult [here]('https://grafana.com/docs/loki/latest/best-practices/')

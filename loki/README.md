# Installing a custom Loki stack in Kyma

## Overview

In contrast to the upstream [loki-stack chart](https://github.com/grafana/helm-charts/tree/main/charts/loki-stack), Kyma's Loki component offers limited configuration options. Furthermore, any modifications you make might be reset at the next upgrade cycle.
To get all the customization options, follow this instructions to set up a parallel installation of the upstream chart, co-existing with the Kyma stack.

>**CAUTION:** This example uses the Grafana Loki version, which is distributed under AGPL-3.0 only and might not be free of charge for commercial usage.

## Prerequisites

- Kyma as the target deployment environment
- Kubectl version 1.22.x or higher
- Helm 3.x

## Installation

### Preparation

1. Export your Namespace as a variable. Replace the `{NAMESPACE}` placeholder in the following command and run it:

    ```bash
    export KYMA_NS="{NAMESPACE}"
    ```

2. Export the Helm release name that you want to use. It can be any name, but be aware that all resources in the cluster will be prefixed with that name. Run the following command:

    ```bash
    export HELM_LOKI_RELEASE="loki"
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
helm upgrade --install --create-namespace -n ${KYMA_NS} ${HELM_LOKI_RELEASE} grafana/loki-stack -f https://raw.githubusercontent.com/kyma-project/examples/main/loki/loki-values.yaml
```
  </details>
  <details>
  <summary label="fluent-bit-installation">
  Installation with Kyma's LogPipeline
  </summary>

>**CAUTION:** This setup uses an unsupported output plugin for the LogPipline. Support for this might be removed in future.

1. To install the Loki stack with Kyma's LogPipeline feature based on Fluent Bit, run:

   ```bash
   helm upgrade --install --create-namespace -n ${KYMA_NS} ${HELM_LOKI_RELEASE} grafana/loki-stack -f https://raw.githubusercontent.com/kyma-project/examples/main/loki/loki-values.yaml --set promtail.enabled=false
   ```

2. Apply the LogPipeline:

   ```bash
   cat <<EOF | kubectl apply -f -
   apiVersion: telemetry.kyma-project.io/v1alpha1
   kind: LogPipeline
   metadata:
     name: custom-loki
   spec:
      input:
         application:
            namespaces:
              system: true
      output:
         custom: |
            name   loki
            host   ${HELM_LOKI_RELEASE}-headless.${KYMA_NS}.svc.cluster.local
            port   3100
            auto_kubernetes_labels off
            labels job=fluentbit, container=$kubernetes['container_name'], namespace=$kubernetes['namespace_name'], pod=$kubernetes['pod_name'], node=$kubernetes['host'], app=$kubernetes['labels']['app'],app=$kubernetes['labels']['app.kubernetes.io/name']
   EOF
   ```

When the status of the applied LogPipeline resource turns into `Running`, the underlying Fluent Bit is reconfigured and log shipment to your Loki instance is active.

  </details>
  <details>
  <summary label="installation-with-grafana">
  Installation with Grafana
  </summary>

  The used Helm chart supports the deployment of Grafana as well, but is disabled by default. Because Grafana provides a very good Loki integration, you might want to install it as well.

  ### Install Loki with Grafana
  1. To deploy Grafana alongside Loki, with Loki pre-configured as a datasource, run:

     ```bash
     helm upgrade --install --create-namespace -n ${KYMA_NS} ${HELM_LOKI_RELEASE} grafana/loki-stack -f https://raw.githubusercontent.com/kyma-project/examples/main/loki/loki-values.yaml -f https://raw.githubusercontent.com/kyma-project/examples/main/loki/grafana-values.yaml --set grafana.adminPassword=myPwd
     ```

  ### Verify the Installation
  1. To access the Grafana UI with kubectl port forwarding, run:

     ```bash
     kubectl -n ${KYMA_NS} port-forward svc/${HELM_LOKI_RELEASE}-grafana 3000:80
     ```
     
     Open Grafana in your browser under `http://localhost:3000` and log in with user admin and the password taken from the previous Helm command.
  
  ### Expose Grafana
  1. To expose Grafana using the Kyma API Gateway, download the APIRule file and replace the `{release-name}` variable with the name of Helm release:
     ```bash
     curl https://raw.githubusercontent.com/kyma-project/examples/main/loki/apirule.yaml -o apirule.yaml
     ```
  1. Create an APIRule:
     ```bash
     kubectl -n ${KYMA_NS} apply -f apirule.yaml 
     ```
  1. Get the public URL of your Loki instance:
     ```bash
     kubectl -n ${KYMA_NS} get vs -l apirule.gateway.kyma-project.io/v1beta1=grafana.${KYMA_NS} -ojsonpath='{.items[*].spec.hosts[*]}'
     ```

  ### Add a Link for Grafana to the Kyma Dashboard
  1. Download the `dashboard-configmap.yaml` file and change `{GRAFANA_LINK}` to the text you retrieved in the previous step.
     ```bash
       curl https://raw.githubusercontent.com/kyma-project/examples/main/loki/dashboard-configmap.yaml -o dashboard-configmap.yaml
     ```
     You can change the label field to change the name of the tab, and the category tab if you wish to move it to another category.

  1. Apply the ConfigMap, and go to Kyma Dashboard. You should see a Link to the newly exposed Grafana under the Observability section. If you already have a busola-config, merge it with the existing one:
     ```bash
     kubectl apply -f dashboard-configmap.yaml 
     ```

  </details>
</div>

>**NOTE:**
>- The used output plugin configuration uses a static label map to assign labels of a Pod to Loki log streams. Activating the `     auto_kubernetes_labels` feature for using all labels of a Pod is not recommended performance-wise. Follow [Loki's labelling best practices](https://grafana.com/docs/loki/latest/best-practices/) for a tailor-made setup that fits your workload configuration.
>- These instructions install Loki in a lightweight setup that does not fulfil production-grade qualities. Consider using a scalable setup based on an object storage backend instead (see [Simple scalable deployment of Grafana Loki with Helm](https://grafana.com/docs/loki/latest/installation/helm/install-scalable/)).


### Verify the installation

Check that the `loki` Pod has been created in the Namespace and is in the `Running` state:

```bash
kubectl -n ${KYMA_NS} get pod -l app=loki,release=${HELM_LOKI_RELEASE}
```

### Verify the setup by accessing logs using the Loki API

1. To access the Loki API, use kubectl port forwarding. Run:

   ```bash
   kubectl -n ${KYMA_NS} port-forward svc/$(kubectl  get svc -n ${KYMA_NS} -l app=loki,release=${HELM_LOKI_RELEASE},variant=headless -ojsonpath='{.items[0].metadata.name}') 3100
   ```

1. Loki queries need a query parameter **time**, provided in nanoseconds. To get the current nanoseconds in Linux or macOS, run:

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
   helm delete -n ${KYMA_NS} ${HELM_LOKI_RELEASE}
   ```

2. To remove the deployed LogPipeline instance from cluster, run:

   ```bash
   kubectl delete -f https://raw.githubusercontent.com/kyma-project/examples/main/loki/logpipeline-custom.yaml
   ```

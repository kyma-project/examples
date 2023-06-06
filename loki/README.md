# Installing a custom Loki stack in Kyma

## Overview

In contrast to the upstream [Loki chart](https://github.com/grafana/loki/tree/main/production/helm/loki), Kyma's Loki component offers limited configuration options. Furthermore, any modifications you make might be reset at the next upgrade cycle.
To get all the customization options, set up a parallel installation of the upstream chart, co-existing with the Kyma stack:

>**CAUTION:** This example uses the Grafana Loki version, which is distributed under AGPL-3.0 only and might not be free of charge for commercial usage.

## Table of Content

- [Prerequisites](#prerequisites)
- [Preparation](#preparation)
- [Loki installation](#loki-installation)
- [Log agent installation](#log-agent-installation)
- [Grafana installation](#grafana-installation)
- [Grafana Exposure](#grafana-exposure)

## Prerequisites

- Kyma as the target deployment environment
- Kubectl version 1.22.x or higher
- Helm 3.x


## Preparation

1. Export your Namespace as a variable. Replace the `{NAMESPACE}` placeholder in the following command and run it:

    ```bash
    export KYMA_NS="{NAMESPACE}"
    ```

2. Export the Helm release names that you want to use. It can be any name, but be aware that all resources in the cluster will be prefixed with that name. Run the following command:

    ```bash
    export HELM_LOKI_RELEASE="loki"
    ```

3. Update your Helm installation with the required Helm repository:

    ```bash
    helm repo add grafana https://grafana.github.io/helm-charts
    helm repo update
    ```

## Loki installation

Depending on your scalability needs and storage requirements, you can install Loki in different [Deployment modes](https://grafana.com/docs/loki/latest/fundamentals/architecture/deployment-modes/). The following instructions install Loki in a lightweight in-cluster solution that does not fulfil production-grade qualities. Consider using a scalable setup based on an object storage backend instead (see [Simple scalable deployment of Grafana Loki with Helm](https://grafana.com/docs/loki/latest/installation/helm/install-scalable/)).

### Install Loki

You install the Loki stack with a Helm upgrade command, which installs the chart if not present yet.

```bash
helm upgrade --install --create-namespace -n ${KYMA_NS} ${HELM_LOKI_RELEASE} grafana/loki -f https://raw.githubusercontent.com/kyma-project/examples/main/loki/loki-values.yaml
```

In any case, you can either create your own `values.yaml` file, or use the [loki-values.yaml](./loki-values.yaml) provided in this `loki` folder, which contains customized settings deviating from the default settings: The provided `values.yaml` file activates the `singleBinary` mode and disables additional components that are usually used when running Loki as a central backend. 

### Verify Loki installation

Check that the `loki` Pod has been created in the Namespace and is in the `Running` state:

```bash
kubectl -n ${KYMA_NS} get pod -l app.kubernetes.io/name=loki
```
## Log agent installation

### Install the log agent

To ingest the application logs from within your cluster to Loki, you can either choose an installation based on [Promtail](https://grafana.com/docs/loki/latest/clients/promtail/), which is the log collector recommended by Loki and provides a ready-to-use setup. Alternatively, you can use Kyma's LogPipeline feature based on Fluent Bit.

<div tabs name="default-settings" group="configuration">
  <details>
  <summary label="promtail-installation">
  Install Promtail
  </summary>

To install Promtail pointing it to the previously installed Loki instance, run:

```bash
helm upgrade --install --create-namespace -n ${KYMA_NS} promtail grafana/promtail -f https://raw.githubusercontent.com/kyma-project/examples/main/loki/promtail-values.yaml --set "config.clients[0].url=https://${HELM_LOKI_RELEASE}.${KYMA_NS}.svc.cluster.local:3100/loki/api/v1/push" 
```
  </details>
  <details>
  <summary label="fluent-bit-installation">
  Install Fluent Bit with Kyma's LogPipeline
  </summary>

>**CAUTION:** This setup uses an unsupported output plugin for the LogPipline.

Apply the LogPipeline:

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
            labels job=fluentbit, container=\$kubernetes['container_name'], namespace=\$kubernetes['namespace_name'], pod=\$kubernetes['pod_name'], node=\$kubernetes['host'], app=\$kubernetes['labels']['app'],app=\$kubernetes['labels']['app.kubernetes.io/name']
   EOF
   ```

When the status of the applied LogPipeline resource turns into `Running`, the underlying Fluent Bit is reconfigured and log shipment to your Loki instance is active.

>**NOTE:** The used output plugin configuration uses a static label map to assign labels of a Pod to Loki log streams. Activating the `     auto_kubernetes_labels` feature for using all labels of a Pod is not recommended performance-wise. Follow [Loki's labelling best practices](https://grafana.com/docs/loki/latest/best-practices/) for a tailor-made setup that fits your workload configuration.

  </details>
</div>

### Verify the setup by accessing logs using the Loki API

1. To access the Loki API, use kubectl port forwarding. Run:

   ```bash
   kubectl -n ${KYMA_NS} port-forward svc/$(kubectl  get svc -n ${KYMA_NS} -l app.kubernetes.io/name=loki -ojsonpath='{.items[0].metadata.name}') 3100
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

## Grafana installation

Because Grafana provides a very good Loki integration, you might want to install it as well.

### Install Grafana

1. To deploy Grafana, run:

   ```bash
   helm upgrade --install --create-namespace -n ${KYMA_NS} grafana grafana/grafana -f https://raw.githubusercontent.com/kyma-project/examples/main/loki/grafana-values.yaml
   ```

1. To enable Loki as Grafana data source, run:
   ```bash
   cat <<EOF | kubectl apply -n ${KYMA_NS} -f -
   apiVersion: v1
   kind: ConfigMap
   metadata:
     name: grafana-loki-datasource
     labels:
       grafana_datasource: "1"
   data:
      loki-datasource.yaml: |-
         apiVersion: 1
         datasources:
         - name: Loki
           type: loki
           access: proxy
           url: "http://${HELM_LOKI_RELEASE}:3100"
           version: 1
           isDefault: false
           jsonData: {}
   EOF
   ```

1. Get the password to access the Grafana UI:

   ```bash
   kubectl get secret --namespace ${KYMA_NS} grafana -o jsonpath="{.data.admin-password}" | base64 --decode ; echo
   ```

1.  To access the Grafana UI with kubectl port forwarding, run:
   ```bash
   kubectl -n ${KYMA_NS} port-forward svc/grafana 3000:80
   ```
   
1. In your browser, open Grafana under `http://localhost:3000` and log in with user admin and the password you retrieved before.
  
## Grafana Exposure

### Expose Grafana
1. To expose Grafana using the Kyma API Gateway, create an APIRule:
   ```bash
   kubectl -n ${KYMA_NS} apply -f apirule.yaml 
   ```
1. Get the public URL of your Loki instance:
   ```bash
   kubectl -n ${KYMA_NS} get virtualservice -l apirule.gateway.kyma-project.io/v1beta1=grafana.${KYMA_NS} -ojsonpath='{.items[*].spec.hosts[*]}'
   ```

### Add a link for Grafana to the Kyma Dashboard
1. Download the `dashboard-configmap.yaml` file and change `{GRAFANA_LINK}` to the public URL of your Grafana instance.
   ```bash
      curl https://raw.githubusercontent.com/kyma-project/examples/main/loki/dashboard-configmap.yaml -o dashboard-configmap.yaml
   ```

1. Optionally, adjust the ConfigMap: You can change the label field to change the name of the tab. If you want to move it to another category, change the category tab.

1. Apply the ConfigMap, and go to Kyma Dashboard. Under the Observability section, you should see a link to the newly exposed Grafana. If you already have a busola-config, merge it with the existing one:
   ```bash
   kubectl apply -f dashboard-configmap.yaml 
   ```
   
## Cleanup

1. To remove the installation from the cluster, run:

   ```bash
   helm delete -n ${KYMA_NS} ${HELM_LOKI_RELEASE}
   helm delete -n ${KYMA_NS} promtail
   helm delete -n ${KYMA_NS} grafana
   ```

2. To remove the deployed LogPipeline instance from cluster, run:

   ```bash
   kubectl delete LogPipeline custom-loki
   ```

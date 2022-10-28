# Install a custom kube-prometheus-stack in Kyma

## Overview

The Kyma monitoring stack often brings limited configuration options in contrast to the upstream [`kube-prometheus-stack`](https://github.com/prometheus-community/helm-charts/blob/main/charts/kube-prometheus-stack) chart. Modifications might be reset at the next upgrade cycle.

As an alternative, you can install the upstream chart with all customization options parallel. This tutorial outlines how to set up such installation in co-existence to the Kyma monitoring stack.

> **CAUTION:** This tutorial describes a basic setup that you should not use in production. Typically, a production setup needs further configuration, like optimizing the amount of data to scrape and the required resource footprint of the installation. To achieve qualities like [high availability](https://prometheus.io/docs/introduction/faq/#can-prometheus-be-made-highly-available), [scalability](https://prometheus.io/docs/introduction/faq/#i-was-told-prometheus-doesnt-scale), or [durable long-term storage](https://prometheus.io/docs/operating/integrations/#remote-endpoints-and-storage), you need a more advanced setup.

> **CAUTION:** This example will use the latest Grafana version which is under AGPL-3.0 and might not be free of charge for commercial usage.

## Prerequisites

- Kyma as the target deployment environment.
- Kubectl > 1.22.x
- Helm 3.x

## Installation

### Preparation
1. Assure that the Kyma monitoring stack running in your cluster is limited to detection of kubernetes resources in the kyma-system namespace only, so that there can be no side-effects with the additional custom stack:
> **NOTE**: That step is only needed for clusters installed manually via the Kyma CLI
    ```bash
    kyma deploy --component monitoring --value monitoring.prometheusOperator.namespaces.releaseNamespace=true
    ```

1. Export your Namespace as a variable. Replace the `{namespace}` placeholder in the following command and run it:

    ```bash
    export KYMA_NS="{namespace}"
    ```
1. If you don't have it created yet, now is the time to do so:
    ```bash
    kubectl create namespace $KYMA_NS
    ```
>**Note**: This Namespace must have **no** Istio sidecar injection enabled so no `istio-injection` label present on the namespace. The Helm chart deploys jobs that will not succeed when sidecar injection is enabled by default.

1. Export the Helm release name that you want to use. It can be any name, but be aware that all resources in the cluster will be prefixed with that name. Replace the `{release-name}` placeholder in the following command and run it:
   ```bash
    export HELM_RELEASE="{release-name}"
    ```

1. Update your Helm installation with the required Helm repository:

    ```bash
    helm repo add prometheus-community https://prometheus-community.github.io/helm-charts
    helm repo update
    ```

### Install the kube-prometheus-stack

1. Run the Helm upgrade command, which installs the chart if it's not present yet. At the end of the command, change the Grafana admin password to some value of your choice.
    ```bash
    helm upgrade --install -n ${KYMA_NS} ${HELM_RELEASE} prometheus-community/kube-prometheus-stack -f https://raw.githubusercontent.com/kyma-project/examples/main/prometheus/values.yaml --set grafana.adminPassword=myPwd
    ```

You can use the [values.yaml](./values.yaml) provided with this tutorial, which contains customized settings deviating from the default settings, or create your own one.
The provided `values.yaml` covers the following adjustments:
- Parallel operation to a Kyma monitoring stack
- Client certificate injection to support scraping of workload secured with Istio strict mTLS
- Active scraping of workload annotated with @prometheus.io/scrape

### Activate scraping of Istio metrics & Grafana dashboards

1. To configure Prometheus for scraping of the Istio-specific metrics from any istio-proxy running in the cluster, deploy a PodMonitor, which scrapes any Pod that has a port with name `.*-envoy-prom` exposed.

    ```bash
    kubectl -n ${KYMA_NS} apply -f https://raw.githubusercontent.com/kyma-project/examples/main/prometheus/istio/podmonitor-istio-proxy.yaml
    ```

2. Deploy a ServiceMonitor definition for the central metrics of the `istiod` deployment:

    ```bash
    kubectl -n ${KYMA_NS} apply -f https://raw.githubusercontent.com/kyma-project/examples/main/prometheus/istio/servicemonitor-istiod.yaml
    ```

3. Get the latest versions of the Istio-specific dashboards.
   Grafana is configured to load dashboards dynamically from ConfigMaps in the cluster, so Istio-specific dashboards can be applied as well. 
   Either follow the [Istio quick start instructions](https://istio.io/latest/docs/ops/integrations/grafana/#option-1-quick-start), or take the prepared ones with the following command:

    ```bash
    kubectl -n ${KYMA_NS} apply -f https://raw.githubusercontent.com/kyma-project/examples/main/prometheus/istio/configmap-istio-grafana-dashboards.yaml
    kubectl -n ${KYMA_NS} apply -f https://raw.githubusercontent.com/kyma-project/examples/main/prometheus/istio/configmap-istio-services-grafana-dashboards.yaml
    ```

   > **NOTE:** This setup collects all Istio metrics on a Pod level, which can lead to cardinality issues. Because  metrics are only needed on service level, for setups having a bigger amount of workloads deployed, it is recommended to use a setup based on federation as described in the [Istio documentation](https://istio.io/latest/docs/ops/best-practices/observability/#using-prometheus-for-production-scale-monitoring).

### Verify the installation

1. You should see several Pods coming up in the Namespace, especially Prometheus and Alertmanager. Assure that all Pods have the "Running" state.
2. Browse the Prometheus dashboard and verify that all "Status->Targets" are healthy. The following command exposes the dashboard on `http://localhost:9090`:
   ```bash
   kubectl -n ${KYMA_NS} port-forward svc/${HELM_RELEASE}-kube-prometh-prometheus 9090
   ```
3. Browse the Grafana dashboard and verify that the dashboards are showing data. The user `admin` is pre-configured in the Helm chart; the password was provided in your `helm install` command. The following command exposes the dashboard on `http://localhost:3000`:
   ```bash
   kubectl -n ${KYMA_NS} port-forward svc/${HELM_RELEASE}-grafana 3000:80
   ```

### Deploy a custom workload and scrape it

1. Follow the tutorial [monitoring-custom-metrics](./../monitoring-custom-metrics/), but use the steps above to verify that the metrics are collected.

### Scrape workload via annotations

Instead of defining a ServiceMonitor per workload for setting up custom metric scraping, you can use a simplified way based on annotations. The used [values.yaml](./values.yaml) defines an `additionalScrapeConfig`, which  scrapes all Pods and services that have these annotations:

```yaml
prometheus.io/scrape: "true"   # mandatory to enable automatic scraping
prometheus.io/scheme: https    # optional, default is http when no istio sidecar is used. When using a sidecar (pod has label security.istio.io/tlsMode=istio) the default will be "https". Use "https" to scrape workloads using istio client certificates. Will only work for annotated services (not pods)
prometheus.io/port: "1234"     # optional, configure the port under which the metrics are exposed
prometheus.io/path: /myMetrics # optional, configure the path under which the metrics are exposed
```

You can try it out by removing the ServiceMonitor from the previous example and instead providing the annotations to the Service manifest.

### Cleanup

1. To remove the installation from the cluster, call Helm:

    ```bash
    helm delete -n ${KYMA_NS} ${HELM_RELEASE}
    ```

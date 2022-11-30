# Create a Grafana dashboard

Follow these sections to create the Gauge dashboard type for the `cpu_temperature_celsius` metric.

## Prerequisites

* Kyma as the target deployment environment.
* Deployed custom kube-prometheus-stack as described in the [prometheus](../) example.
* Deployed [custom-metrics example](../monitoring-custom-metrics).

## Create the dashboard

1. [Access Grafana](../README.md#verify-the-installation).
2. Add a new dashboard with a new panel.
3. For your new query, select **Prometheus** from the data source selector.
4. Pick the `cpu_temperature_celsius` metric.
5. To retrieve the latest metric value on demand, activate the **Instant** switch.
6. From the visualization panels, select the **Gauge** dashboard type.
7. Save your changes and provide a name for the dashboard.

## Configure the dashboard

1. To edit the dashboard settings, go to the **Panel Title** options and select **Edit**.
2. Find the **Field** tab and set the measuring unit to Celsius degrees, indicating the metric data type.
3. Set the minimum metric value to `60` and the maximum value to `90`, indicating the `cpu_temperature_celsius` metric value range.
4. For the dashboard to turn red once the CPU temperature reaches and exceeds 75°C, set a red color threshold to `75`.
5. Go to the **Panel** tab and title the dashboard, for example, `CPU Temperature`.
6. To display this range on the dashboard, make sure that under **Panel > Display**, the threshold labels and threshold markers are activated.
7. Save your changes. We recommend that you add a note to describe the changes made.

## Check the dashboard

Refresh the browser to see how the dashboard changes according to the current value of the `cpu_temperature_celsius` metric.

- If the current metric value ranges from 60 to 74 degrees Celsius, it turns **green**.
- If the current metric value ranges from 75 to 90 degrees Celsius, it turns **red**.

## Add the dashboard as Kubernetes resource

When you create a dashboard to monitor one of your applications (Function, microservice,...), we recommend that you define the dashboard as a Kubernetes ConfigMap resource. In this case, a Grafana sidecar automatically loads the Dashboard on Grafana startup. Following that approach, you can easily keep the dashboard definition together with the Kubernetes resource definitions of your application and port it to different clusters.

1. Create a JSON document with the dashboard definition; for example, by exporting it from Grafana.
2. Create a Kubernetes resource with a unique name for your dashboard and the JSON content, like the following example:

   ```yaml
   apiVersion: v1
   kind: ConfigMap
   metadata:
     name: {UNIQUE_DASHBOARD_NAME}-grafana-dashboard
     labels:
       grafana_dashboard: "1"
   data:
     {UNIQUE_DASHBOARD_NAME}-dashboard.json: |-
       {
         # dashboard JSON content
       }
   ```

3. To apply the Kubernetes resource created in the previous step to your cluster, run:

   ```bash
   kubectl apply -f <{UNIQUE_CONFIGMAP_NAME}.yaml>
   ```

4. Restart the Grafana deployment with the following command:

   ```bash
   kubectl -n ${KYMA_NS} rollout restart deployment ${HELM_RELEASE}-grafana
   ```

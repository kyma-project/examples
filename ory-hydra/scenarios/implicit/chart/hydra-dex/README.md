# Integration of Dex as Hydra Login Provider

## Introduction

This chart bootstraps a [ORY Hydra](https://www.ory.sh/docs/hydra/) "login and consent" application capable of performing Dex-based logins on a [Kyma](https://kyma-project.io) cluster.

ORY Hydra on it's own does not provide any user authentication features. This must be provided externally.
Hydra offers REST-based extension points to integrate with an external login provider.
This chart provides a sample login application that delegates login requests to the Dex running in a Kyma cluster.
In this way, we can use Hydra JWT tokens issued for users authenticated by one of available Dex authentication methods.


## Prerequisites
- Ability to install Kyma from sources, or a running Kyma installation with **ory** component installed.

## Installation

The installation assumes installation from Kyma sources. There is a way to install this chart on already running cluster, assuming that **ory** component is already installed. Look for additional information in _Note_ sections.

1. Add new static client to Dex configuration in the **resources/dex/templates/dex-config-map.yaml** file.
   - Add this entry to the **staticClients** list:
    ```
    - id: hydra-integration
      name: 'Hydra Integration'
      redirectURIs:
      - 'https://oauth2-login-consent.<domainName>/cb'
      secret: <secretValue>
     ```
   - Replace `<domainName>` with the proper domain name of Kyma ingress gateway for cluster installations or **kyma.local** for local installations.
   - Replace `<secretValue>` with a secure random string.

    _Note: Alternatively, you can change this in a running Kyma installation by modifying the respective Config Map and restarting Dex._

3. Change hydra server configuration to point to a valid login-and-consent application.
   - Modify **resources/ory/charts/hydra/values.yaml** file and set the `loginConsent.name` attribute to the same value as defined in **this** chart: `oauth2-login-consent`

   _Note: Alternatively, you can change this in a running Kyma installation by modifying the values of **OAUTH2_CONSENT_URL** and **OAUTH2_LOGIN_URL** environment variables of the Hydra server Deployment (`kubectl edit Deployment ory-hydra-oauth2 -n kyma-system`). Ensure hydra server Pod is redeployed with new values._

4. Install Kyma with **ory** chart enabled.

5. Use Helm to install hydra-dex integration chart.
   - Run: `export DOMAIN_NAME=<domainName>`. Replace `<domainName>` with the proper domain name of Kyma ingress gateway for cluster installations or **kyma.local** for local installation.
   - Run: `helm install . -n hydra-dex --namespace kyma-system --set loginConsent.domainName=${DOMAIN_NAME} --tls`



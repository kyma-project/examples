# Tunnel

## Overview

This example illustrates how to connect application runnig on local machine to Kyma cluster.

## Prerequisites

- [Kyma](https://kyma-project.io/docs/) as the target deployment environment
- [kubectl](https://kubernetes.io/docs/tasks/tools/install-kubectl/) configured to your Kyma cluster
- [Helm](https://helm.sh/docs/intro/install/)
- [inlets](https://github.com/inlets/inlets#install-the-cli) as a reverse proxy and websocket tunnel


## Installation

Create exit node in Kyma cluster:
```
export TOKEN="thesecretpasswordrequiredtoaccessthetunnel" 
export NAMESPACE=test
export PROXY_NAME=proxy
export CHART=https://raw.githubusercontent.com/kyma-project/examples/master/tunnel/inlets-0.1.0.tgz

kubectl create ns $NAMESPACE
helm install $PROXY_NAME $CHART -n $NAMESPACE --set token=$TOKEN
```

Start local application (in separate terminal session):
```
docker run -p 10000:10000 eu.gcr.io/kyma-project/xf-application-mocks/commerce-mock:latest
```

Get remote proxy url, and connect host node to exit node with inlets client:
```
export REMOTE_URL="$(kubectl get virtualservice -n $NAMESPACE -l apirule.gateway.kyma-project.io/v1alpha1=$PROXY_NAME-inlets.$NAMESPACE -o jsonpath='{ .items[0].spec.hosts[0] }')"
echo "Open https://$REMOTE_URL in your browser"

inlets client \
 --remote=wss://$REMOTE_URL \
 --upstream=http://127.0.0.1:10000 \
 --token=$TOKEN
```

Now you can access your local application through remote proxy url printed by the command above, and pair the commerce-mock application with your Kyma runtime.

## Clean up

Stop the inlets client and delete helm release and namespace:
```
helm delete $PROXY_NAME -n $NAMESPACE
kubectl delete ns $NAMESPACE
```
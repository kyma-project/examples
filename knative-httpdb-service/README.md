# HTTP DB Service Example with Knative

## Overview

This is a simple example of how to use Knative's build and serving component with Kyma.
We'll use the source code of the HTTP DB Service Example, but instead of manually building a Docker image, pushing it to a docker repository and referencing it in a Kubernetes deployment we're going to use Knative to build and serve our app.

This example is just meant to illustrate how to build and deploy workloads on Kyma with Knative. Knative will scale your workloads up to multiple or even down to no instances, meaning you will have inconsistencies and data loss when interacting with this example service.

## Prerequisites

- A Kyma installation (locally or on a cluster) with Knative integration installed and with a dedicated namespace you can deploy the example to (See [Installation with Knative](https://kyma-project.io/docs/root/kyma/#installation-installation-with-knative) for how install the Knative integration for Kyma)
- Credentials for [Docker Hub](https://hub.docker.com/) or any other container image registry
- kubectl
- The yaml files in this folder

## Setup

### Set your namespace

To apply everything you do to the namespace you designated for this example, set your namespace with:

```
kubectl config set-context $(kubectl config current-context) --namespace={YOUR_NAMESPACE}
```

Replace `{YOUR_NAMESPACE}` with the name of your namespace. Alternatively you can append `-n={YOUR_NAMESPACE}` to all `kubectl` calls.

### Provide your container registry credentials

Open the file `registry-credentials.yaml` with a text editor. Encode your Docker Hub username as well as your Docker Hub password in base64 and replace the placeholders `{YOUR_ENCODED_USERNAME}` and `{YOUR_ENCODED_PASSWORD}` with the respective values.

To encode your username and your password in base64 you can use tools of your operating system (`{VALUE}` is the placeholder for the username/password you have to encode):

- For macOS and Linux: Execute `echo "{VALUE} | base64` on your shell
  
- For Windows: Open a PowerShell window and execute: `[Convert]::ToBase64String([System.Text.Encoding]::UTF8.GetBytes("{VALUE}"))`

If you want to use another registry than Docker Hub also adjust the registry URL in the file.

Save your changes, then execute this command to push the credentials to Kubernetes:

```
kubectl apply -f registry-credentials.yaml
```

### Configure the service

The file `http-db-service.yaml` contains the description of the target deployment, including the instructions to build the Docker container.
The configuration instructs Knative to pull the sources from the GitHub repository `kyma-project/examples` and let Google's kaniko tool build the docker image and push it to the registry.
As soon as the built image is available, Knative will run the image, make it accessible per https and automatically scale container instances up or down based on the amount of incoming requests.

Open it and replace `{YOUR_USERNAME}` with your Docker Hub username so that your Docker Hub account will be used. Don't base64-encode your username here! 😃
If you are not using Docker Hub adjust the respective parts of the image name as well.

**By default Knative expects images to expose their service on port 8080, whereas the http-db-service example uses the port 8017. In the yaml file the `containerPort` is declared to let Knative know to use this port. However, this is not supported by Knative versions prior to 0.4. If your Kyma installation uses an older version, fork the GitHub repository, change the source code so that port 8080 will be used and replace the GitHub URL in the yaml file with your fork's url.**

Apply the configuration with:

```
kubectl apply -f http-db-service.yaml
```

### Let Knative do the work

Knative should now be starting the build process. You can see that the service is waiting for an image to become available with:

```
kubectl describe ksvc http-db-service
```

You can see that a build has been created by typing:

```
kubectl get builds
```

You can retrieve detailed information on the build's status with (replace `http-db-service-0001` with the build name from the previous step if necessary):

```
kubectl describe build http-db-service-0001
```

Once the service is ready you should be able to access it at `https://http-db-service.{NAMESPACE}.{CLUSTERDOMAIN}`. If you are unsure about the domain, you can look it up with:

```
kubectl get ksvc http-db-service -ojsonpath='{.status.domain}'
```

If you are running Kyma locally in a minikube environment it might be necessary to manually add that domain name to your hosts file (macOS/Linux: `/etc/hosts`, Windows: `C:\Windows\System32\drivers\etc\hosts`) and point it to the IP address of your minikube.


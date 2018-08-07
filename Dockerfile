FROM alpine:3.8

LABEL source="git@github.com:kyma-project/examples.git"

ENV KUBECTL_VERSION=v1.11.1
ENV HELM_VERSION=v2.8.2
ENV KUBELESS_VERSION=v1.0.0-alpha.7

RUN apk add --no-cache curl tar gzip

RUN curl -Lo /usr/bin/kubectl https://storage.googleapis.com/kubernetes-release/release/${KUBECTL_VERSION}/bin/linux/amd64/kubectl && chmod +x /usr/bin/kubectl

RUN curl -Lo /tmp/kubeless.zip https://github.com/kubeless/kubeless/releases/download/${KUBELESS_VERSION}/kubeless_linux-amd64.zip && unzip -q /tmp/kubeless.zip -d /tmp/ && mv /tmp/bundles/kubeless_linux-amd64/kubeless /usr/bin/ && rm -r /tmp/kubeless.zip /tmp/bundles && chmod +x /usr/bin/kubeless

RUN curl -L https://storage.googleapis.com/kubernetes-helm/helm-${HELM_VERSION}-linux-amd64.tar.gz | tar xz && mv linux-amd64/helm /bin/helm && rm -rf linux-amd64

RUN mkdir -p /root/.kube && touch /root/.kube/config

ADD . /
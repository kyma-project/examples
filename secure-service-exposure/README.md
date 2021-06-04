# Secure Service Exposure with OAuth2 Proxy

## This example provides a deployment to expose services in a secure way using [OAuth2 Proxy](https://oauth2-proxy.github.io/oauth2-proxy/).

1. Create a namespace for the proxy, for instance `oauth2-proxy`:

```bash
kubectl create namespace oauth2-proxy
```

2. Create a `Secret` with the OAuth2 Proxy configuration [environment variables](https://oauth2-proxy.github.io/oauth2-proxy/docs/configuration/overview/#environment-variables). See [OAuth2 Proxy docs](https://oauth2-proxy.github.io/oauth2-proxy/docs/configuration/oauth_provider/) for provider specific settings. Set `OAUTH2_PROXY_UPSTREAMS` variable to your in cluster service:

```bash
kubectl -n oauth2-proxy create secret generic oauth2-proxy-secret \
--from-literal="OAUTH2_PROXY_UPSTREAMS=http://my-service.my-namespace.svc:80" \
--from-literal="OAUTH2_PROXY_COOKIE_SECRET=`openssl rand -hex 16`" \
--from-literal="OAUTH2_PROXY_COOKIE_NAME=oauth2-proxy-`openssl rand -hex 16`" \
--from-literal="OAUTH2_PROXY_CLIENT_ID=<my-client-id>" \
--from-literal="OAUTH2_PROXY_CLIENT_SECRET=<my-client-secret>" \
--from-literal="OAUTH2_PROXY_PROVIDER=<my-provider>" \
...
```

3. Deploy OAuth2 Proxy:

```bash
kubectl -n oauth2-proxy apply -f https://raw.githubusercontent.com/kyma-project/examples/main/secure-service-exposure/oauth2_proxy.yaml
```

4. Create `VirtualService`. Adapt the domain in the hosts list to your needs:

```bash
cat <<EOF | kubectl -n oauth2-proxy apply -f -
apiVersion: networking.istio.io/v1alpha3
kind: VirtualService
metadata:
  name: oauth2-proxy
spec:
  hosts:
  - my-service.kyma.example.com
  gateways:
  - kyma-system/kyma-gateway
  http:
  - match:
    - uri:
        regex: /.*
    route:
    - destination:
        port:
          number: 3000
        host: oauth2-proxy
EOF
```

If you want to expose multiple services (for instance Grafana and Kiali), deploy multiple instances of OAuth2 Proxy to different namespaces.

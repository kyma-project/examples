```

  ______                           _              _____ _                _   
 |  ____|                         | |            / ____| |              | |  
 | |__  __  ____ _ _ __ ___  _ __ | | ___  ___  | |    | |__   __ _ _ __| |_ 
 |  __| \ \/ / _` | '_ ` _ \| '_ \| |/ _ \/ __| | |    | '_ \ / _` | '__| __|
 | |____ >  < (_| | | | | | | |_) | |  __/\__ \ | |____| | | | (_| | |  | |_ 
 |______/_/\_\__,_|_| |_| |_| .__/|_|\___||___/  \_____|_| |_|\__,_|_|   \__|
                            | |                                              
                            |_|                                              

```
## Overview

This chart provides an easy way to deploy and test the examples.

## Prerequisites

- Kubernetes 1.10+
- Kyma as the target deployment environment.
- An environment to deploy the examples.

## Details

### Installation

Configure these options on [values.yaml](values.yaml):

| Parameter                        | Description |
|--------------------------------- | -----------: |
| examples.image                   | Image for the examples |
| examples.httpDBService.deploy    | Deploy [HTTP DB Service](../http-db-service) example |
| examples.httpDBService.deploymentImage | Deployment image for HTTP DB Service |
| examples.httpDBService.testImage | Test image for HTTP DB Service |
| examples.eventSubscription.lambda.deploy | Deploy [Event Subscription lambda](../event-subscription/lambda) example |
| examples.eventEmailService.deploy | Deploy [Event Email Service](../event-email-service) example |
| examples.eventEmailService.deploymentImage | Deployment image for Event Email Service example |
| rbac.enabled  | Enable RBAC |

Deploy the examples:
```
Helm install -f values.yaml --name examples --namespace <environment> .
```

#### Testing

Deploy the test pods defined under [templates/tests](templates/tests):
```
helm test --cleanup examples
```
Output of this command will show whether the tests are passed or failed.

#### Cleanup

Cleanup the environment:
```
helm delete --purge examples
kubectl delete ns <environment>
```

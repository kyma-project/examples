#!/usr/bin/env groovy
import groovy.json.JsonSlurperClassic
import groovy.json.JsonOutput

def label = "kyma-${UUID.randomUUID().toString()}"
application = 'examples'
def isMaster = params.GIT_BRANCH == 'master'

dockerPushRoot = isMaster 
    ? "${env.DOCKER_REGISTRY}"
    : "${env.DOCKER_REGISTRY}snapshot/"

dockerImageTag = isMaster
    ? params.APP_VERSION
    : params.GIT_BRANCH

def changes = parseJson("${params.CHANGED_EXAMPLES}") 

//For now, we only have deployment pods for these examples. Once we have for all, we can just eliminate this check.
def deploy = changes.contains("http-db-service") || changes.contains("tests/http-db-service") || changes.contains("event-email-service") || changes.contains("event-subscription/lambda") || changes.contains("service-binding/lambda")

echo """
********************************
Job started with the following parameters:
DOCKER_REGISTRY=${env.DOCKER_REGISTRY}
DOCKER_CREDENTIALS=${env.DOCKER_CREDENTIALS}
GIT_REVISION=${params.GIT_REVISION}
GIT_BRANCH=${params.GIT_BRANCH}
APP_VERSION=${params.APP_VERSION}
CHANGED_EXAMPLES=${JsonOutput.prettyPrint(params.CHANGED_EXAMPLES)}
********************************
"""

podTemplate(label: label) {
    node(label) {
        try {
            timestamps {
                timeout(time:20, unit:"MINUTES") {
                    ansiColor('xterm') {
                        stage("setup") {
                            checkout scm

                            if(dockerImageTag == ""){
                               error("No version for docker tag defined, please set APP_VERSION parameter for master branch or GIT_BRANCH parameter for any branch")
                            }

                            withCredentials([usernamePassword(credentialsId: env.DOCKER_CREDENTIALS, passwordVariable: 'pwd', usernameVariable: 'uname')]) {
                                sh "docker login -u $uname -p '$pwd' $env.DOCKER_REGISTRY"
                            }

                            withCredentials([usernamePassword(credentialsId: 'examples-jenkins-user', passwordVariable: 'pwd', usernameVariable: 'uname')]) {
                                sh "curl -o kubeconfig --user $uname:$pwd https://jenkins.poc.servicefactory.cd.ydev.hybris.com/job/azure/job/ondemand/job/huskiesOnDemand/lastSuccessfulBuild/artifact/kyma/kubeconfig"
                            }
                        }

                        stage("build image $application") {
                            sh "docker build -t $application:latest ."
                        }

                        stage("push image $application") {
                            sh "docker tag ${application}:latest ${dockerPushRoot}${application}:${dockerImageTag}"
                            sh "docker push ${dockerPushRoot}${application}:${dockerImageTag}"
                        }

                        if (deploy) {
                            stage("create namespace for $application") {
                                execute("kubectl create ns ${params.GIT_REVISION}; kubectl label ns ${params.GIT_REVISION} env=true")
                            }

                            stage("Increase resource quota for ns")
                                execute("cat <<EOF | kubectl apply -f- 
apiVersion: v1
kind: List
items:
- apiVersion: v1
  kind: ResourceQuota
  metadata:
    name: kyma-default
    namespace: ${params.GIT_REVISION}
  spec:
    hard:
      limits.memory: 5Gi
EOF")

                            stage("deploy $application") {
                                execute("cd examples-chart && helm install --wait --timeout=600 --name examples -f values.yaml --namespace ${params.GIT_REVISION} . --set examples.image=${dockerPushRoot}${application}:${dockerImageTag} " + configureChart(changes))
                            }

                            stage("test $application") {
                                execute("helm test examples")
                            }
                        }
                    }
                }
            }
        } catch (ex) {
            echo "Got exception: ${ex}"
            currentBuild.result = "FAILURE"
            def body = "${currentBuild.currentResult} ${env.JOB_NAME}${env.BUILD_DISPLAY_NAME}: on branch: ${params.GIT_BRANCH}. See details: ${env.BUILD_URL}"
            emailext body: body, recipientProviders: [[$class: 'DevelopersRecipientProvider'], [$class: 'CulpritsRecipientProvider'], [$class: 'RequesterRecipientProvider']], subject: "${currentBuild.currentResult}: Job '${env.JOB_NAME} [${env.BUILD_NUMBER}]'"
        } finally {
            if (deploy){
                stage("print logs for $application") {
                    execute("kubectl logs -l chart=examples -n ${params.GIT_REVISION}")
                }
                stage("print logs for tests"){
                    execute("kubectl logs -l chart=examples-tests -n ${params.GIT_REVISION}")
                }
                stage("delete $application") {
                    execute("helm delete --purge examples")
                }
                stage("delete namespace for $application") {
                    execute("kubectl delete ns ${params.GIT_REVISION}")
                }
            }
        }
    }
}

def configureChart(changedExamples) {
    def set = ""
    
    def deployHttpDBService = changedExamples.contains("http-db-service")
    def deployHttpDBServiceTests = changedExamples.contains("tests/http-db-service")
    if (deployHttpDBService || deployHttpDBServiceTests) {
        set += "--set examples.httpDBService.deploy=true "
        if (deployHttpDBService) {
            set += "--set examples.httpDBService.deploymentImage=${dockerPushRoot}example/http-db-service:${dockerImageTag} "
        }
        if (deployHttpDBServiceTests) {
            set += "--set examples.httpDBService.testImage=${dockerPushRoot}example/http-db-service-acceptance-tests:${dockerImageTag} "
        }
    }
    if (changedExamples.contains("event-email-service")) {
        set += "--set examples.eventEmailService.deploy=true --set examples.eventEmailService.deploymentImage=${dockerPushRoot}example/event-email-service:${dockerImageTag} "
    }
    if (changedExamples.contains("event-subscription/lambda")) {
        set += "--set examples.eventSubscription.lambda.deploy=true "
    }
    if (changedExamples.contains("service-binding/lambda")) {
        set += "--set examples.serviceBinding.lambda.deploy=true "
    }
    
    return set
}

def execute(command, envs = '') {
    def envText = envs=='' ? '' : "--env $envs"
    def workDir = pwd()
    sh "docker run --rm --env KUBECONFIG=/kubeconfig $envText ${dockerPushRoot}${application}:${dockerImageTag} /bin/sh -c '$command'"
}

@NonCPS
def parseJson(changedExamples) {
    def parser = new JsonSlurperClassic()
    return parser.parseText(changedExamples)
}

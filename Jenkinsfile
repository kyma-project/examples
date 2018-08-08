#!/usr/bin/env groovy
def label = "kyma-${UUID.randomUUID().toString()}"
def application = 'examples'
def isMaster = params.GIT_BRANCH == 'master'

def dockerPushRoot = isMaster 
    ? "${env.DOCKER_REGISTRY}"
    : "${env.DOCKER_REGISTRY}snapshot/"

def dockerImageTag = isMaster
    ? params.APP_VERSION
    : params.GIT_BRANCH

def deploy = false

echo """
********************************
Job started with the following parameters:
DOCKER_REGISTRY=${env.DOCKER_REGISTRY}
DOCKER_CREDENTIALS=${env.DOCKER_CREDENTIALS}
GIT_REVISION=${params.GIT_REVISION}
GIT_BRANCH=${params.GIT_BRANCH}
APP_VERSION=${params.APP_VERSION}
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
                                execute("kubectl create ns ${dockerImageTag}")
                            }

                            stage("deploy $application") {
                                dir("examples-chart") {
                                    execute("helm install --wait --timeout=600 -f values.yaml --namespace ${dockerImageTag} --set examples.image=${dockerPushRoot}${application}:${dockerImageTag} .")
                                }
                            }

                            stage("print logs for $application") {
                                execute("kubectl logs -l chart=examples")
                            }

                            /* stage("test $application") {
                                execute("helm test")
                            }

                            stage("print test logs for $application") {
                                execute("kubectl logs -l chart=examples-tests")
                            } */

                            stage("delete namespace for $application") {
                                execute("kubectl delete ns ${dockerImageTag}")
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
        }
    }
}

def execute(command, envs = '') {
    def envText = envs=='' ? '' : "--env $envs"
    sh "docker run --rm $envText ${dockerPushRoot}${application}:${dockerImageTag} /bin/sh -c '$command'"
}


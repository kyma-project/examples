#!/usr/bin/env groovy
def label = "kyma-${UUID.randomUUID().toString()}"

echo """
********************************
Job started with the following parameters:
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
                        }

                        stage("validate external links") {
                            validateLinks('find . -name "*.md" | grep -v "vendor" | grep -v "./call-ec/web-ui/README.md"')
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

def validateLinks(command) {
    workDir = pwd()
    whiteList = "github.com,localhost,core-publish.kyma-system,kyma-integration,svc.cluster.local,calculate-promotion,kyma.local"
    sh "docker run --rm -v $workDir:/mnt:ro dkhamsing/awesome_bot --allow-dupe --allow-redirect --skip-save-results --allow-ssl --white-list $whiteList `$command`"
}
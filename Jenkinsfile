@Library('whiteblock-dev')_

//General
def APP_NAME = "genesis-cli"
def GCP_REGION = "us-west2"
def DEFAULT_BRANCH = 'master'


//Dev
def DEV_GCP_PROJECT_ID = "infra-dev-249211"
def IMAGE_REPO            = "gcr.io/${DEV_GCP_PROJECT_ID}"
def BINARIES_BUCKET = 'infra-dev-binaries'


pipeline {
  agent any
  environment {
    REV_SHORT = sh(script: "git log --pretty=format:'%h' -n 1", , returnStdout: true).trim()
    DISABLE_AUTH = 'true'
  }
  options {
    buildDiscarder(logRotator(numToKeepStr: '10'))
  }
  stages {
    stage('Static-Analysis') {

      when {
        anyOf {
          changeRequest target: DEFAULT_BRANCH
        }
      }
      steps {
        goFmt()
        goVet()
      }
    }

    stage('push to store') {
      when {
        anyOf {
          branch DEFAULT_BRANCH
        }
      }
      steps {
        script {
          withCredentials([file(credentialsId: 'google-infra-dev-auth', variable: 'GOOGLE_APPLICATION_CREDENTIALS')]) {
            sh """
              gcloud auth activate-service-account --key-file ${GOOGLE_APPLICATION_CREDENTIALS}
              docker build . -t ${IMAGE_REPO}/${APP_NAME}:${BRANCH_NAME}-build-latest
              docker push ${IMAGE_REPO}/${APP_NAME}:${BRANCH_NAME}-build-latest
            """
          }
        }
      }
    }
  }
  post {
    always {
      deleteDir()
      sh "/usr/bin/docker image prune --force --filter 'until=72h'"
      sh "/usr/bin/docker image rm ${IMAGE_REPO}/${APP_NAME}:${BRANCH_NAME}-build-latest || true"
      sh "gcloud auth revoke || true"
    }
    failure {
      script {
        slackNotify(env.BRANCH_NAME == DEFAULT_BRANCH)
      }
    }
  }
}

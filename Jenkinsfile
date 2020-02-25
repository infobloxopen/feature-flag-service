pipeline {
  agent {
    label 'ubuntu_docker_label'
  }
  tools {
    go "Go 1.13"
  }
  options {
    checkoutToSubdirectory('src/github.com/Infoblox-CTO/atlas.feature.flag')
  }
  environment {
    GOPATH = "$WORKSPACE"
    DIRECTORY = "src/github.com/Infoblox-CTO/atlas.feature.flag"
  }
  stages {
    stage('Test') {
      steps {
        sh 'cd $DIRECTORY && make test'
        sh 'helm init --client-only'
        sh 'cd $DIRECTORY && make .helm-lint'
      }
    }
    stage('Build image') {
      steps {
        sh 'cd $DIRECTORY && IMAGE_VERSION=\$(make show-image-version)-j$BUILD_NUMBER make docker'
      }
    }
    stage('Push image') {
      when {
        anyOf { branch 'master'; buildingTag() }
      }
      steps {
        withDockerRegistry([credentialsId: "dockerhub-bloxcicd", url: ""]) {
          sh "cd $DIRECTORY && IMAGE_VERSION=\$(make show-image-version)-j$BUILD_NUMBER make push"
        }
      }
    }
    stage('Push Chart') {
      when {
        anyOf { branch 'master'; buildingTag() }
      }
      steps {
        dir("$WORKSPACE/$DIRECTORY") {
          withDockerRegistry([credentialsId: "dockerhub-bloxcicd", url: ""]) {
            withAWS(region:'us-east-1', credentials:'CICD_HELM') {
              sh "IMAGE_VERSION=\$(make show-image-version)-j$BUILD_NUMBER make push-chart"
            }
          }
          archiveArtifacts artifacts: 'deploy/*.tgz'
          archiveArtifacts artifacts: 'build.properties'
        }
      }
    }
  }
}

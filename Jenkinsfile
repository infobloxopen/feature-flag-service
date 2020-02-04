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
    stage("Lint") {
      steps {
        dir("$DIRECTORY") {
          sh "make fmt && git diff --exit-code"
        }
      }
    }
    stage("Test") {
      steps {
        dir("$DIRECTORY") {
          sh "make test"
        }
      }
    }
    stage("Build") {
      steps {
        withDockerRegistry([credentialsId: "dockerhub-bloxcicd", url: ""]) {
          dir("$DIRECTORY") {
            sh "make docker push"
          }
        }
      }
    }
    stage("Push") {
      when {
        branch "master"
      }
      steps {
        withDockerRegistry([credentialsId: "dockerhub-bloxcicd", url: ""]) {
          dir("$DIRECTORY") {
            sh "make push IMAGE_VERSION=latest"
          }
        }
      }
    }
  }
  post {
    always {
      dir("$DIRECTORY") {
        sh "make clean || true"
      }
      cleanWs()
    }
  }
}

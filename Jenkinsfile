pipeline {
    agent any
    
    environment {
        IMAGE_NAME = 'argyarijal/api-gateway:canary'
        CONTAINER_NAME = ''
        BRANCH_NAME = "main"
        MSG_COMMIT = sh(script: "git log -1 --pretty=%B ${env.GIT_COMMIT}", returnStdout: true).trim()
    }
    
    stages {
       
        stage('Skip pattern') {
            steps {
                script {
                    if (MSG_COMMIT == 'generated by jenkins') {
                        echo 'Build Canceled'
                        currentBuild.result = 'ABORTED'
                        throw new Exception()
                    }
                }
            }
        }

        stage('Docker Build and Push') {
            steps {
                withDockerRegistry([credentialsId: "docker-creds", url: ""]) {
                    retry(3) {
                        timeout(time: 25, unit: 'MINUTES') {
                            sh 'printenv'
                            sh 'DOCKER_BUILDKIT=1 docker build --rm=false -t ${IMAGE_NAME} .'
                            sh 'docker push ${IMAGE_NAME}'
                        }
                    }
                }
            }
        }

        stage('Prune Docker Data') {
            steps {
                sh '''
                if [ $(docker ps -a -q -f name=${CONTAINER_NAME}) ]; then
                    docker stop ${CONTAINER_NAME} || true
                    docker rm ${CONTAINER_NAME} || true
                fi
                '''
            }
        }
        stage('Deploy to portainer') {
            steps {
                script {
                sh "curl -k -X POST https://193.203.167.97:9443/api/stacks/webhooks/936922e1-b1e8-421a-93a0-8df4f580bd5b"
                echo "Deployed to Portainer"
                }
            }
        }
    }
}

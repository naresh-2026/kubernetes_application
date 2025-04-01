pipeline {
    agent any
    environment {
        DOCKER_USERNAME = 'naresh2026'
        DOCKER_REPO = 'repo1'
        DOCKER_CREDENTIALS_ID = 'cred1'
        DOCKER_REGISTRY = 'https://index.docker.io/v1/'
    }
    stages {
        stage('Checkout Code') {
            steps {
                checkout scm
            }
        }
        
        stage('Build and Push Docker Images') {
    steps {
        script {
            docker.withRegistry(${DOCKER_REGISTRY}, DOCKER_CREDENTIALS_ID) {
            def dirs = sh(script: 'find . -type d -mindepth 1 -maxdepth 1', returnStdout: true).trim().split("\n")
            def imageCount = 1  // Start numbering from 1

            dirs.each { dir ->
                def dockerfileExists = sh(script: "test -f ${dir}/Dockerfile && echo 'exists' || echo 'not found'", returnStdout: true).trim()

                if (dockerfileExists == 'exists') {
                    def imageName = "ser${imageCount}"
                    echo "Building Docker Image for ${dir} as ${imageName}"
                    docker.build(imageName, dir)

                    def image_name = "${DOCKER_USERNAME}/${DOCKER_REPO}"
                    echo "Pushing Docker Image ${image_name} to Docker Registry"
                    docker.image(image_name).push("ser${imageCount}")

                    imageCount++
                } else {
                    echo "Skipping ${dir} as no Dockerfile found."
                }
            }
            }
        }
    }
}


        // stage('Push Docker Images') {
        //     steps {
        //         script {
        //             docker.withRegistry("https://${DOCKER_REGISTRY}", DOCKER_CREDENTIALS_ID) {
        //                 def dirs = sh(script: 'find . -type d -mindepth 1 -maxdepth 1', returnStdout: true).trim().split("\n")
        //                 def imageCount = 1

        //                 dirs.each { dir ->
        //                     def imageName = "${DOCKER_USERNAME}/${DOCKER_REPO}:ser${imageCount}"
        //                     echo "Pushing Docker Image ${imageName} to Docker Registry"
        //                     docker.image(imageName).push('latest')
        //                     imageCount++
        //                 }
        //             }
        //         }
        //     }
        // }
    }

    post {
        success {
            echo 'Pipeline executed successfully!'
        }
        failure {
            echo 'Pipeline failed!'
        }
    }
}

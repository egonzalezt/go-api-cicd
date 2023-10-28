def sendDiscordNotification(String result, String latestCommitMessage, String latestCommitAuthor) {
    withCredentials([string(credentialsId: 'discord-webhook-credential-id', variable: 'DISCORD_WEBHOOK_URL')]) {
        def discordSendConfig = [
            footer: JOB_NAME,
            link: env.BUILD_URL,
            result: result,
            title: "Jenkins Pipeline Build: ${JOB_NAME}",
            webhookURL: DISCORD_WEBHOOK_URL,
            thumbnail: "https://cdn.discordapp.com/attachments/678439901544316931/1165804342713000047/icegif-59.gif"
        ]

        if (result == 'SUCCESS') {
            def successDescription = """
                # **Pipeline Execution Successful**

                * Build Number: [${env.BUILD_NUMBER}](${env.BUILD_URL})
                * Branch: ${env.BRANCH_NAME}
                * New Semantic Version: ${env.NEW_SEMANTIC_VERSION}
                * Latest Commit Message: ${latestCommitMessage}
                * Author of the Last Commit: ${latestCommitAuthor}
            """
            discordSendConfig.description = successDescription
            discordSendConfig.image = "https://cdn.discordapp.com/attachments/1081839152942813324/1165799959052951552/undefined_-_Imgur.gif"
        } else {
            def failureDescription = """# **Pipeline Execution Failed**
                * Build Number: [${env.BUILD_NUMBER}](${env.BUILD_URL})
                * Branch: ${env.BRANCH_NAME}
                * Failure Reason: ${result}

            **Additional Details**
                * Latest Commit Message: ${latestCommitMessage}
            """
            try {
                latestTag = sh(script: 'git describe --abbrev=0 --tags', returnStdout: true).trim()
                failureDescription += """
                    * Latest Tag: ${latestTag}
                """
            } catch (Exception e) {
                echo "No tags found. Setting default tag to 0.0.0"
            }
            discordSendConfig.description = failureDescription
            discordSendConfig.image = "https://cdn.discordapp.com/attachments/1082173364552081449/1165807160236716052/kirbo-mad.gif"
        }

        discordSend(discordSendConfig)
    }
}
def getLatestCommitInfo() {
    def commitInfo = sh(
        script: """
        latestCommitMessage=\$(git log -1 --pretty=%B)
        latestCommitAuthor=\$(git log -1 --pretty=%an)
        echo "Latest Commit Message: \$latestCommitMessage"
        echo "Author: \$latestCommitAuthor"
        """,
        returnStdout: true
    ).trim()

    def latestCommitMessage = commitInfo.contains("Latest Commit Message:") ? commitInfo.split("Latest Commit Message:")[1].trim() : "Commit message not found"
    def latestCommitAuthor = commitInfo.contains("Author:") ? commitInfo.split("Author:")[1].trim() : "Commit author not found"

    return [latestCommitMessage, latestCommitAuthor]
}
def determineSemanticVersion() {
    def latestCommitMessage = sh(script: 'git log --pretty=%B -n 1', returnStdout: true).trim()
    def majorKeyword = "release("
    def minorKeyword = "feat("
    def patchKeyword = "fix("
    def latestTag

    try {
        latestTag = sh(script: 'git describe --abbrev=0 --tags', returnStdout: true).trim()
    } catch (Exception e) {
        echo "No tags found. Setting default tag to 0.0.0"
        latestTag = "0.0.0"
    }
    echo "Latest Tag: ${latestTag}"
    def newVersion = latestTag ?: "0.1.0"

    if (latestCommitMessage.contains(majorKeyword)) {
        newVersion = newVersion.tokenize('.').collect { it as Integer }
        newVersion[0]++
        newVersion[1] = 0
    } else if (latestCommitMessage.contains(minorKeyword)) {
        newVersion = newVersion.tokenize('.').collect { it as Integer }
        newVersion[1]++
        newVersion[2] = 0
    } else if (latestCommitMessage.contains(patchKeyword)) {
        newVersion = newVersion.tokenize('.').collect { it as Integer }
        newVersion[2]++
    } else {
        error("Commit does not contain the required format to create the version (feat, fix, release).")
    }

    def newSemanticVersion = newVersion.join('.')
    echo "New Semantic Version: ${newSemanticVersion}"
    currentBuild.description = "Semantic Version: ${newSemanticVersion}"
    env.NEW_SEMANTIC_VERSION = newSemanticVersion
}
pipeline {
    agent {
        label 'agent'
    }
    tools {
        go 'go-1.20'
        dockerTool 'docker'
    }
    environment {
        GO111MODULE = 'on'
        GOPATH = "/home/jenkins/go"
    }
    options {
        disableConcurrentBuilds()
    }
    stages {
        stage('Setup') {
            steps {
                sh 'mkdir -p $HOME/go/bin'
                sh 'export GOPATH=$HOME/go'
                sh 'export PATH=$PATH:$GOPATH/bin'
                sh 'go install github.com/t-yuki/gocover-cobertura@latest'
                sh 'docker ps'
            }
        }
        stage('Get Repo Name') {
           steps {
               script {
                    def repoUrl = sh(returnStdout: true, script: 'git config --get remote.origin.url').trim()
                    def repoName = repoUrl.tokenize('/')[-1].replaceAll('\\.git', '')
                    env.GITHUB_REPOSITORY = repoName
                    echo "Repository name: ${repoName}"
                }
           }
        }
        stage('Install Dependencies') {
            steps {
                sh 'go mod tidy'
            }
        }
        stage('Build') {
            steps {
                sh 'go build'
            }
        }
        stage('Test') {
            steps {
                sh 'ls'
                sh 'pwd'
                sh 'go test -coverprofile=coverage.txt -covermode count ./'
                publishChecks name: 'test', title: 'Test Check', summary: 'Test check through pipeline',
                    text: 'Test check in pipeline script',
                    detailsURL: "${env.BUILD_URL}checks/${currentBuild.number}",
                    actions: [[label:'test-action', description:'Test action', identifier:'test-identifier']]
            }
        }
        stage('Coverage') {
            steps {
                sh '/home/jenkins/go/bin/gocover-cobertura < coverage.txt > coverage.xml'
                cobertura(coberturaReportFile: 'coverage.xml')
            }
        }
        stage('Export Artifacts') {
            steps {
                archiveArtifacts artifacts: 'coverage.xml', fingerprint: true
            }
        }
        stage('Determine Semantic Version') {
            steps {
                script {
                    determineSemanticVersion()
                }
            }
        }
        stage('Create Tag') {
            steps {
                withCredentials([usernamePassword(credentialsId: 'jenkins-github-app',
                                                usernameVariable: 'GITHUB_APP',
                                                passwordVariable: 'GITHUB_ACCESS_TOKEN')]) {
                    sh "git tag ${env.NEW_SEMANTIC_VERSION}"
                    sh "git push https://$GITHUB_APP:$GITHUB_ACCESS_TOKEN@github.com/egonzalezt/${GITHUB_REPOSITORY}.git ${env.NEW_SEMANTIC_VERSION}"
                }
            }
        }
        stage('Build Docker Image and Publish') {
            steps {
                script {
                    def dockerImageName = 'vasitos/go-ci-cd'
                    def dockerImageTag = "${env.NEW_SEMANTIC_VERSION}"
                    sh "docker build -t ${dockerImageName}:${dockerImageTag} ."
                    withCredentials([usernamePassword(credentialsId: 'docker-hub', usernameVariable: 'USERNAME', passwordVariable: 'PASSWORD')]) {
                        sh 'docker login -u $USERNAME -p $PASSWORD'
                    }
                    sh "docker push ${dockerImageName}:${dockerImageTag}"
                }
            }
        }
        stage('Cleanup') {
            steps {
                script {
                    def dockerImageName = 'vasitos/go-ci-cd'
                    def dockerImageTag = "${env.NEW_SEMANTIC_VERSION}"
                    sh "docker rmi ${dockerImageName}:${dockerImageTag}"
                }
            }
        }
    }
    post {
        always {
            script {
                def result = currentBuild.currentResult
                def commitInfo = getLatestCommitInfo()
                def latestCommitMessage = commitInfo[0]
                def latestCommitAuthor = commitInfo[1]
                sendDiscordNotification(result, latestCommitMessage, latestCommitAuthor)
            }
        }
    }
}

pipeline {
    agent {
        label 'agent'
    }
    tools {
        go 'go-1.20'
    }
    environment {
        GO111MODULE = 'on'
        GOPATH = "/home/jenkins/go"
    }
    stages {
        stage('Setup') {
            steps {
                sh 'mkdir -p $HOME/go/bin'
                sh 'export GOPATH=$HOME/go'
                sh 'export PATH=$PATH:$GOPATH/bin'
                sh 'go install github.com/t-yuki/gocover-cobertura@latest'
            }
        }
        stage('Get Repo Name')') {
           steps {
               script {
                  def determineRepoName = {
                      return scm.getUserRemoteConfigs()[0].getUrl().tokenize('/').last().split('.')[0]
                  }
                  env.GITHUB_REPOSITORY = determineRepoName()
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
                    def latestCommitMessage = sh(script: 'git log --pretty=%B -n 1', returnStdout: true).trim()
                    def majorKeyword = "release("
                    def minorKeyword = "feat("
                    def patchKeyword = "fix("
                    def latestTag
                    try {
                        latestTag = sh(script: 'git describe --abbrev=0 --tags', returnStdout: true).trim()
                    } catch (Exception e) {
                        echo "No tags found. Setting default tag to 0.0.0."
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
                    }   else {
                        error("Commit does not contain the required format to create the version (feat, fix, release).")
                    }

                    def newSemanticVersion = newVersion.join('.')
                    echo "New Semantic Version: ${newSemanticVersion}"
                    currentBuild.description = "Semantic Version: ${newSemanticVersion}"
                    env.NEW_SEMANTIC_VERSION = newSemanticVersion
                }
            }
        }
        stage('Create Tag') {
            steps {
                withCredentials([usernamePassword(credentialsId: 'Jenkins-Github-App',
                                                usernameVariable: 'GITHUB_APP',
                                                passwordVariable: 'GITHUB_ACCESS_TOKEN')]) {
                    sh "git tag ${env.NEW_SEMANTIC_VERSION}"
                    sh "git push https://${GITHUB_APP}:${GITHUB_ACCESS_TOKEN}@github.com/egonzalezt/${GITHUB_REPOSITORY}.git ${env.NEW_SEMANTIC_VERSION}"
                }
            }
        }
    }
    post {
        always {
            script {
                withCredentials([string(credentialsId: 'discord-webhook-credential-id', variable: 'DISCORD_WEBHOOK_URL')]) {
                    def discordSendConfig = [
                        description: currentBuild.currentResult == 'SUCCESS' ? "<:LETSFUCKINGOOOOOOOOOOO:809820731134705714> Jenkins Pipeline Build":"<:weynooo:799854983100629022> Jenkins Pipeline Build",
                        footer: JOB_NAME,
                        link: env.BUILD_URL,
                        result: currentBuild.currentResult,
                        title: "Jenkins Pipeline Build: ${JOB_NAME}",
                        webhookURL: DISCORD_WEBHOOK_URL,
                        image: currentBuild.currentResult == 'SUCCESS' ? "https://cdn.discordapp.com/attachments/1081839152942813324/1165799959052951552/undefined_-_Imgur.gif" : "https://cdn.discordapp.com/attachments/1082173364552081449/1165807160236716052/kirbo-mad.gif",
                        thumbnail: "https://cdn.discordapp.com/attachments/678439901544316931/1165804342713000047/icegif-59.gif"
                    ]
                    discordSend(discordSendConfig)
                }
            }
        }
    }
}

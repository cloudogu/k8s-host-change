#!groovy

@Library('github.com/cloudogu/ces-build-lib@4.1.1')
import com.cloudogu.ces.cesbuildlib.*

// Creating necessary git objects
git = new Git(this, "cesmarvin")
git.committerName = 'cesmarvin'
git.committerEmail = 'cesmarvin@cloudogu.com'
gitflow = new GitFlow(this, git)
github = new GitHub(this, git)
changelog = new Changelog(this)
goVersion = "1.24.1"
makefile = new Makefile(this)

// Configuration of repository
repositoryOwner = "cloudogu"
repositoryName = "k8s-host-change"
project = "github.com/${repositoryOwner}/${repositoryName}"
registry = "registry.cloudogu.com"
registry_namespace = "k8s"

// Configuration of branches
productionReleaseBranch = "main"
developmentBranch = "develop"
currentBranch = "${env.BRANCH_NAME}"

helmTargetDir = "target/k8s"
helmChartDir = "${helmTargetDir}/helm"

node('docker') {
    timestamps {
        stage('Checkout') {
            checkout scm
            make 'clean'
        }

        stage('Lint') {
            lintDockerfile()
        }

        stage('Check Markdown Links') {
            Markdown markdown = new Markdown(this)
            markdown.check()
        }

        new Docker(this)
            .image("golang:${goVersion}")
            .mountJenkinsUser()
            .inside("--volume ${WORKSPACE}:/go/src/${project} -w /go/src/${project}") {
                stage('Build') {
                    make 'build-job'
                }

                stage("Unit test") {
                    make 'unit-test'
                    junit allowEmptyResults: true, testResults: 'target/unit-tests/*-tests.xml'
                }

                stage("Review dog analysis") {
                    stageStaticAnalysisReviewDog()
                }

                stage('Generate k8s Resources') {
                    make 'helm-generate'
                    archiveArtifacts "${helmTargetDir}/**/*"
                }

                stage("Lint helm") {
                    make 'helm-lint'
                }
            }

        stage('SonarQube') {
            stageStaticAnalysisSonarQube()
        }

        stage('Trivy scan') {
            Docker docker = new Docker(this)
            def dockerImage = docker.build("cloudogu/${repositoryName}:ci-build")

            Trivy trivy = new Trivy(this)
            trivy.scanImage("cloudogu/${repositoryName}:ci-build", TrivySeverityLevel.CRITICAL, TrivyScanStrategy.UNSTABLE)
            trivy.saveFormattedTrivyReport(TrivyScanFormat.TABLE)
            trivy.saveFormattedTrivyReport(TrivyScanFormat.JSON)
            trivy.saveFormattedTrivyReport(TrivyScanFormat.HTML)
        }

        stageAutomaticRelease()
    }
}

void gitWithCredentials(String command) {
    withCredentials([usernamePassword(credentialsId: 'cesmarvin', usernameVariable: 'GIT_AUTH_USR', passwordVariable: 'GIT_AUTH_PSW')]) {
        sh(
                script: "git -c credential.helper=\"!f() { echo username='\$GIT_AUTH_USR'; echo password='\$GIT_AUTH_PSW'; }; f\" " + command,
                returnStdout: true
        )
    }
}

void stageStaticAnalysisReviewDog() {
    def commitSha = sh(returnStdout: true, script: 'git rev-parse HEAD').trim()

    withCredentials([[$class: 'UsernamePasswordMultiBinding', credentialsId: 'sonarqube-gh', usernameVariable: 'USERNAME', passwordVariable: 'REVIEWDOG_GITHUB_API_TOKEN']]) {
        withEnv(["CI_PULL_REQUEST=${env.CHANGE_ID}", "CI_COMMIT=${commitSha}", "CI_REPO_OWNER=${repositoryOwner}", "CI_REPO_NAME=${repositoryName}"]) {
            make 'static-analysis-ci'
        }
    }
}

void stageStaticAnalysisSonarQube() {
    def scannerHome = tool name: 'sonar-scanner', type: 'hudson.plugins.sonar.SonarRunnerInstallation'
    withSonarQubeEnv {
        sh "git config 'remote.origin.fetch' '+refs/heads/*:refs/remotes/origin/*'"
        gitWithCredentials("fetch --all")

        if (currentBranch == productionReleaseBranch) {
            echo "This branch has been detected as the production branch."
            sh "${scannerHome}/bin/sonar-scanner -Dsonar.branch.name=${env.BRANCH_NAME}"
        } else if (currentBranch == developmentBranch) {
            echo "This branch has been detected as the development branch."
            sh "${scannerHome}/bin/sonar-scanner -Dsonar.branch.name=${env.BRANCH_NAME}"
        } else if (env.CHANGE_TARGET) {
            echo "This branch has been detected as a pull request."
            sh "${scannerHome}/bin/sonar-scanner -Dsonar.pullrequest.key=${env.CHANGE_ID} -Dsonar.pullrequest.branch=${env.CHANGE_BRANCH} -Dsonar.pullrequest.base=${developmentBranch}"
        } else if (currentBranch.startsWith("feature/")) {
            echo "This branch has been detected as a feature branch."
            sh "${scannerHome}/bin/sonar-scanner -Dsonar.branch.name=${env.BRANCH_NAME}"
        } else {
            echo "This branch has been detected as a miscellaneous branch."
            sh "${scannerHome}/bin/sonar-scanner -Dsonar.branch.name=${env.BRANCH_NAME} "
        }
    }
    timeout(time: 2, unit: 'MINUTES') { // Needed when there is no webhook for example
        def qGate = waitForQualityGate()
        if (qGate.status != 'OK') {
            unstable("Pipeline unstable due to SonarQube quality gate failure")
            unstable("Pipeline unstable due to SonarQube quality gate failure")
        }
    }
}

void stageAutomaticRelease() {
    if (gitflow.isReleaseBranch()) {
        String releaseVersion = git.getSimpleBranchName()
        String dockerReleaseVersion = releaseVersion.split("v")[1]
        String controllerVersion = makefile.getVersion()

        stage('Build & Push Image') {
            Docker docker = new Docker(this)
            def dockerImage = docker.build("cloudogu/${repositoryName}:${dockerReleaseVersion}")
            docker.withRegistry('https://registry.hub.docker.com/', 'dockerHubCredentials') {
                dockerImage.push("${dockerReleaseVersion}")
            }
        }

        stage('Finish Release') {
            gitflow.finishRelease(releaseVersion, productionReleaseBranch)
        }

        stage('Push Helm chart to Harbor') {
            new Docker(this)
                .image("golang:${goVersion}")
                .mountJenkinsUser()
                .inside("--volume ${WORKSPACE}:/go/src/${project} -w /go/src/${project}") {
                    // Package operator-chart & crd-chart
                    make 'helm-package'
                    archiveArtifacts "${helmTargetDir}/**/*"

                    // Push charts
                    withCredentials([usernamePassword(credentialsId: 'harborhelmchartpush', usernameVariable: 'HARBOR_USERNAME', passwordVariable: 'HARBOR_PASSWORD')]) {
                        sh ".bin/helm registry login ${registry} --username '${HARBOR_USERNAME}' --password '${HARBOR_PASSWORD}'"

                        sh ".bin/helm push ${helmChartDir}/${repositoryName}-${controllerVersion}.tgz oci://${registry}/${registry_namespace}/"
                    }
                }
        }

        stage('Add Github-Release') {
            releaseId = github.createReleaseWithChangelog(releaseVersion, changelog, productionReleaseBranch)
        }
    }
}

void make(String makeArgs) {
    sh "make ${makeArgs}"
}
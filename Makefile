# Set these to the desired values
ARTIFACT_ID=k8s-host-change
VERSION=0.1.0

GOTAG?=1.20.2
MAKEFILES_VERSION=7.5.0

## Image URL to use all building/pushing image targets
IMAGE_DEV=${K3CES_REGISTRY_URL_PREFIX}/${ARTIFACT_ID}:${VERSION}
IMAGE=cloudogu/${ARTIFACT_ID}:${VERSION}

K8S_RESOURCE_DIR=${WORKDIR}/k8s
K8S_HOST_CHANGE_RESOURCE_YAML=${K8S_RESOURCE_DIR}/k8s-host-change.yaml

include build/make/variables.mk

# make sure to create a statically linked binary otherwise it may quit with
# "exec user process caused: no such file or directory"
GO_BUILD_FLAGS=-mod=vendor -a -tags netgo,osusergo $(LDFLAGS) -o $(BINARY)
# remove DWARF symbol table and strip other symbols to shave ~13 MB from binary
ADDITIONAL_LDFLAGS=-extldflags -static -w -s
LINT_VERSION?=v1.52.1

include build/make/self-update.mk
include build/make/dependencies-gomod.mk
include build/make/build.mk
include build/make/test-common.mk
include build/make/test-integration.mk
include build/make/test-unit.mk
include build/make/static-analysis.mk
include build/make/clean.mk
include build/make/digital-signature.mk
include build/make/mocks.mk

K8S_POST_GENERATE_TARGETS=k8s-generate-job-resource
include build/make/k8s.mk

##@ EcoSystem

.PHONY: build
build: k8s-delete image-import k8s-apply ## Builds a new version of the setup and deploys it into the K8s-EcoSystem.

.PHONY: k8s-generate-job-resource
k8s-generate-job-resource: ${BINARY_YQ} $(K8S_RESOURCE_TEMP_FOLDER) template-dev-only-image-pull-policy template-stage template-log-level ## Generates the final resource yaml.
	@echo "Applying image transformation..."
	@$(BINARY_YQ) -i e "(select(.kind == \"Job\").spec.template.spec.containers[]|select(.image == \"*$(ARTIFACT_ID)*\").image)=\"$(IMAGE_DEV)\"" $(K8S_RESOURCE_TEMP_YAML)
	@echo "Done."

##@ Build

.PHONY: build-job
build-setup: ${SRC} compile ## Builds the setup Go binary.

.PHONY: run
run: ## Run a setup from your host.
	go run ./main.go

.PHONY: k8s-create-temporary-resource
k8s-create-temporary-resource: $(K8S_RESOURCE_TEMP_FOLDER)
	@cp $(K8S_HOST_CHANGE_RESOURCE_YAML) $(K8S_RESOURCE_TEMP_YAML)
	@sed -i "s/'{{ .Version }}'/$(VERSION)/" $(K8S_RESOURCE_TEMP_YAML)

.PHONY: template-dev-only-image-pull-policy
template-dev-only-image-pull-policy: $(BINARY_YQ)
	@if [ ${STAGE}"X" = "development""X" ]; \
		then echo "Setting pull policy to always for development stage!" && $(BINARY_YQ) -i e "(select(.kind == \"Job\").spec.template.spec.containers[]|select(.image == \"*$(ARTIFACT_ID)*\").imagePullPolicy)=\"Always\"" $(K8S_RESOURCE_TEMP_YAML); \
	fi

.PHONY: template-stage
template-stage:
	@echo "Setting STAGE env in deployment to ${STAGE}!"
	@$(BINARY_YQ) -i e "(select(.kind == \"Deployment\").spec.template.spec.containers[]|select(.image == \"*$(ARTIFACT_ID)*\").env[]|select(.name==\"STAGE\").value)=\"${STAGE}\"" $(K8S_RESOURCE_TEMP_YAML)

.PHONY: template-log-level
template-log-level:
	@echo "Setting LOG_LEVEL env in deployment to ${LOG_LEVEL}!"
	@$(BINARY_YQ) -i e "(select(.kind == \"Deployment\").spec.template.spec.containers[]|select(.image == \"*$(ARTIFACT_ID)*\").env[]|select(.name==\"LOG_LEVEL\").value)=\"${LOG_LEVEL}\"" $(K8S_RESOURCE_TEMP_YAML)

##@ Release

.PHONY: job-release
job-release: ## Interactively starts the release workflow.
	@echo "Starting git flow release..."
	@build/make/release.sh host-change

SHELL=/bin/bash

export KUBEVIRT_VERSION = main
export KUBEVIRT_TAG = main

# Use the COMMON_INSTANCETYPES_CRI env variable to control if the following targets are executed within a container.
# Supported runtimes are docker and podman. By default targets run directly on the host.
export COMMON_INSTANCETYPES_IMAGE = quay.io/kubevirtci/common-instancetypes-builder
export COMMON_INSTANCETYPES_IMAGE_TAG = v20241017-30aa54a

# Containerdisk image and tag to use for Validation OS functional tests
export VALIDATION_OS_IMAGE = quay.io/kubevirtci/validation-os-container-disk
export VALIDATION_OS_IMAGE_TAG = 26100.2033.241003-1823

# Packages of golang tools vendored in ./tools
# Version to install is defined in ./tools/go.mod
KUSTOMIZE_PACKAGE ?= sigs.k8s.io/kustomize/kustomize/v5
KUBECONFORM_PACKAGE ?= github.com/yannh/kubeconform/cmd/kubeconform
YQ_PACKAGE ?= github.com/mikefarah/yq/v4

# Version of golangci-lint to install
GOLANGCI_LINT_VERSION ?= v1.55.2

# Kubeconfig default value
KUBECONFIG ?= ~/.kube/config

.PHONY: all
all: lint validate readme test-lint test

.PHONY: build_image
build_image:
	scripts/build_image.sh

.PHONY: push_image
push_image:
	scripts/push_image.sh

.PHONY: lint
lint: generate
	scripts/cri.sh "scripts/lint.sh"

.PHONY: generate
generate: kustomize yq
	scripts/generate.sh

.PHONY: deploy
deploy: lint
	scripts/deploy-kubevirt-and-cdi.sh && KUBECTL=kubectl BASEDIR=$(CURDIR) scripts/sync.sh

.PHONY: functest
functest:
	cd tests && KUBECONFIG=$(KUBECONFIG) go test -v -timeout 0 ./functests/... -ginkgo.v -ginkgo.randomize-all $(FUNCTEST_EXTRA_ARGS)

.PHONY: validate
validate: generate schema kubeconform
	scripts/validate.sh

.PHONY: schema
schema:
	scripts/cri.sh "scripts/schema.sh"

.PHONY: readme
readme: generate
	scripts/readme.sh

.PHONY: check-tree-clean
check-tree-clean: readme test-fmt
	scripts/check-tree-clean.sh

.PHONY: cluster-up
cluster-up:
	scripts/kubevirtci.sh up

.PHONY: cluster-down
cluster-down:
	scripts/kubevirtci.sh down

.PHONY: cluster-sync
cluster-sync: kustomize
	scripts/kubevirtci.sh sync

.PHONY: cluster-sync-containerdisks
cluster-sync-containerdisks:
	scripts/kubevirtci.sh sync-containerdisks

.PHONY: cluster-functest
cluster-functest:
	cd tests && KUBECONFIG=$$(../scripts/kubevirtci.sh kubeconfig) go test -v -timeout 0 ./functests/... -ginkgo.v -ginkgo.randomize-all $(FUNCTEST_EXTRA_ARGS)

.PHONY: kubevirt-up
kubevirt-up:
	scripts/kubevirt.sh up

.PHONY: kubevirt-down
kubevirt-down:
	scripts/kubevirt.sh down

.PHONY: kubevirt-sync
kubevirt-sync: kustomize
	scripts/kubevirt.sh sync

.PHONY: kubevirt-sync-containerdisks
kubevirt-sync-containerdisks:
	scripts/kubevirt.sh sync-containerdisks

.PHONY: kubevirt-functest
kubevirt-functest:
	cd tests && KUBECONFIG=$$(../scripts/kubevirt.sh kubeconfig) go test -v -timeout 0 ./functests/... -ginkgo.v -ginkgo.randomize-all $(FUNCTEST_EXTRA_ARGS)

.PHONY: test
test: generate
	cd tests && go test -v -timeout 0 ./unittests/...

.PHONY: test-fmt
test-fmt:
	cd tests && go fmt ./...

.PHONY: test-vet
test-vet:
	cd tests && go vet ./...

.PHONY: test-lint
test-lint: golangci-lint test-vet
	cd tests && golangci-lint run --timeout 5m

.PHONY: clean
clean:
	rm -rf _bin _build _cluster-up _kubevirt _schemas

# Location to install local binaries to
LOCALBIN ?= $(PWD)/_bin
$(LOCALBIN):
	mkdir -p $(LOCALBIN)
export PATH := $(LOCALBIN):$(PATH)

KUSTOMIZE ?= $(LOCALBIN)/kustomize
kustomize: $(KUSTOMIZE)
$(KUSTOMIZE): $(LOCALBIN)
	cd tools && GOBIN=$(LOCALBIN) go install $(KUSTOMIZE_PACKAGE)

KUBECONFORM ?= $(LOCALBIN)/kubeconform
kubeconform: $(KUBECONFORM)
$(KUBECONFORM): $(LOCALBIN)
	cd tools && GOBIN=$(LOCALBIN) go install $(KUBECONFORM_PACKAGE)

YQ ?= $(LOCALBIN)/yq
yq: $(YQ)
$(YQ): $(LOCALBIN)
	cd tools && GOBIN=$(LOCALBIN) go install $(YQ_PACKAGE)

GOLANGCI_LINT ?= $(LOCALBIN)/golangci-lint
golangci-lint: $(GOLANGCI_LINT)
$(GOLANGCI_LINT): $(LOCALBIN)
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(LOCALBIN) $(GOLANGCI_LINT_VERSION)

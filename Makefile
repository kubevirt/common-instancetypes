SHELL=/bin/bash

# Run `make schema` when updating this version and commit the created files.
# TODO(lyarwood) - Host the expanded JSON schema elsewhere under the kubevirt namespace
export KUBEVIRT_VERSION = main

# Use the COMMON_INSTANCETYPES_CRI env variable to control if the following targets are executed within a container.
# Supported runtimes are docker and podman. By default targets run directly on the host.
export COMMON_INSTANCETYPES_IMAGE = quay.io/kubevirtci/common-instancetypes-builder

.PHONY: all
all: lint validate readme check

.PHONY: build_image
build_image:
	scripts/build_image.sh

.PHONY: push_image
push_image:
	scripts/push_image.sh

.PHONY: lint
lint: generate
	scripts/cri.sh  "scripts/lint.sh"

.PHONY: generate
generate:
	scripts/cri.sh  "scripts/generate.sh"

.PHONY: validate
validate: generate schema
	scripts/cri.sh  "scripts/validate.sh"

.PHONY: schema
schema:
	scripts/cri.sh  "scripts/schema.sh"

.PHONY: readme
readme: generate
	scripts/cri.sh  "scripts/readme.sh"

.PHONY: check
check:
	scripts/cri.sh  "scripts/check.sh"

.PHONY: cluster-up
cluster-up:
	scripts/kubevirtci.sh up

.PHONY: cluster-down
cluster-down:
	scripts/kubevirtci.sh down

.PHONY: cluster-sync
cluster-sync:
	scripts/kubevirtci.sh sync

.PHONY: cluster-functest
cluster-functest:
	scripts/kubevirtci.sh functest

.PHONY: kubevirt-up
kubevirt-up:
	scripts/kubevirt.sh up

.PHONY: kubevirt-down
kubevirt-down:
	scripts/kubevirt.sh down

.PHONY: kubevirt-sync
kubevirt-sync:
	scripts/kubevirt.sh sync

.PHONY: kubevirt-functest
kubevirt-functest:
	scripts/kubevirt.sh functest

.PHONY: clean
clean:
	rm -f common-*-bundle.yaml CHECKSUMS.sha256

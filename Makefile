SHELL=/bin/bash

# Run `make schema` when updating this version and commit the created files.
# TODO(lyarwood) - Host the expanded JSON schema elsewhere under the kubevirt namespace
export KUBEVIRT_VERSION = main

# Use the COMMON_INSTANCETYPES_CRI env variable to control if the following targets are executed within a container.
# Supported runtimes are docker and podman. By default targets run directly on the host.
export COMMON_INSTANCETYPES_IMAGE = quay.io/kubevirtci/common-instancetypes-builder

all: lint schema validate readme check

build_image:
	./scripts/build_image.sh

push_image:
	./scripts/push_image.sh

lint: generate
	./scripts/cri.sh  "./scripts/lint.sh"

generate:
	./scripts/cri.sh  "./scripts/generate.sh"

validate: generate
	./scripts/cri.sh  "./scripts/validate.sh"

schema:
	./scripts/cri.sh  "./scripts/schema.sh"

readme: generate
	./scripts/cri.sh  "./scripts/readme.sh"

check:
	./scripts/cri.sh  "./scripts/check.sh"

.PHONY: cluster-up
cluster-up:
	./scripts/kubevirtci.sh up

.PHONY: cluster-down
cluster-down:
	./scripts/kubevirtci.sh down

.PHONY: cluster-sync
cluster-sync:
	./scripts/kubevirtci.sh sync

.PHONY: cluster-functest
cluster-functest:
	./scripts/kubevirtci.sh functest

.PHONY: kubevirt-up
kubevirt-up:
	./scripts/kubevirt.sh up

.PHONY: kubevirt-down
kubevirt-down:
	./scripts/kubevirt.sh down

.PHONY: kubevirt-sync
kubevirt-sync:
	./scripts/kubevirt.sh sync

.PHONY: kubevirt-functest
kubevirt-functest:
	./scripts/kubevirt.sh functest

clean:
	rm -f common-*-bundle.yaml

.PHONY: all build_image lint generate schema validate readme

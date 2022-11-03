SHELL=/bin/bash

# Run `make schema` when updating this version and commit the created files.
# TODO(lyarwood) - Host the expanded JSON schema elsewhere under the kubevirt namespace
export KUBEVIRT_VERSION = main

# Use the COMMON_INSTANCETYPES_CRI env variable to control if the following targets are executed within a container.
# Supported runtimes are docker and podman. By default targets run directly on the host.
export COMMON_INSTANCETYPES_IMAGE = quay.io/kubevirtci/common-instancetypes-builder

all: lint generate validate readme

build_image:
	./scripts/build_image.sh

lint: 
	./scripts/cri.sh  "./scripts/lint.sh"

generate: 
	./scripts/cri.sh  "./scripts/generate.sh"

validate:
	./scripts/cri.sh  "./scripts/validate.sh"

schema:
	./scripts/cri.sh  "./scripts/schema.sh"

readme:
	./scripts/cri.sh  "./scripts/readme.sh"

clean: 
	rm -f common*.yaml

.PHONY: all build_image lint generate validate readme

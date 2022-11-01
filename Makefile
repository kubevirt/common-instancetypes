SHELL=/bin/bash

# Run `make schema` when updating this version and commit the created files.
# TODO(lyarwood) - Host the expanded JSON schema elsewhere under the kubevirt namespace
export KUBEVIRT_VERSION = main

all: lint generate

lint:
	yamllint .

generate:
	./scripts/generate.sh

schema:
	./scripts/schema.sh

clean: 
	rm -f common*.yaml

.PHONY: all lint generate

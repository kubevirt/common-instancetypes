SHELL=/bin/bash

all: lint generate

lint:
	yamllint .

generate: common-instancetypes.yaml
	echo "---" > common-instancetypes.yaml
	kustomize build >> common-instancetypes.yaml

clean: common-instancetypes.yaml
	rm -f common-instancetypes.yaml && touch common-instancetypes.yaml

.PHONY: all lint generate

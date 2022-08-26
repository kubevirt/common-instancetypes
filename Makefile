SHELL=/bin/bash

all: lint generate

lint:
	yamllint .

generate: 
	echo "---" > common-instancetypes-all-bundle.yaml
	kustomize build >> common-instancetypes-all-bundle.yaml
	echo "---" > common-instancetypes-bundle.yaml
	kustomize build VirtualMachineInstancetypes >> common-instancetypes-bundle.yaml
	echo "---" > common-clusterinstancetypes-bundle.yaml
	kustomize build VirtualMachineClusterInstancetypes >> common-clusterinstancetypes-bundle.yaml
	echo "---" > common-preferences-bundle.yaml
	kustomize build VirtualMachinePreferences >> common-preferences-bundle.yaml
	echo "---" > common-clusterpreferences-bundle.yaml
	kustomize build VirtualMachineClusterPreferences >> common-clusterpreferences-bundle.yaml

clean: 
	rm -f *.yaml

.PHONY: all lint generate

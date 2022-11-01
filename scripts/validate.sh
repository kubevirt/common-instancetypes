#!/bin/bash

# Validate the output of the various kustomization builds available in the repo
for f in ./ VirtualMachineClusterInstancetypes VirtualMachineClusterPreferences VirtualMachineInstancetypes VirtualMachinePreferences; do
    kustomize build ${f} | kubeconform -exit-on-error -strict -schema-location "schemas/${KUBEVIRT_VERSION}/{{ .ResourceKind }}.json"
    if (( $? > 0 )); then
        exit 1
    fi
done

# Validate the generated bundles
for f in $(ls common*.yaml); do
    kubeconform -exit-on-error -strict -schema-location "schemas/${KUBEVIRT_VERSION}/{{ .ResourceKind }}.json"  ${f}
    if (( $? > 0 )); then
        exit 1
    fi
done
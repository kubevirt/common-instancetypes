#!/bin/bash

if ! command -v kubeconform &> /dev/null; then
    echo "kubeconform is not installed, see https://github.com/yannh/kubeconform#installation for more details."
    exit 1
fi

if ! command -v kustomize &> /dev/null; then
    echo "kustomize is not installed, see https://kubectl.docs.kubernetes.io/installation/kustomize/ for more details."
    exit 1
fi

# Validate the freshly generated bundles
for f in $(ls -1 $@); do
    if ! kubeconform -exit-on-error -strict -schema-location "_schemas/${KUBEVIRT_VERSION}/{{ .ResourceKind }}.json"  "${f}" ; then
        exit 1
    fi
done

#!/bin/bash
set -e

function latest_version() {
  local repo="$1"
  curl -s https://api.github.com/repos/kubevirt/"${repo}"/releases/latest | \
    grep '"tag_name":' | \
    sed -E 's/.*"([^"]+)".*/\1/'
}


# Deploying kubevirt
KUBEVIRT_VERSION=$(curl -L https://storage.googleapis.com/kubevirt-prow/devel/release/kubevirt/kubevirt/stable.txt)
KUBEVIRT_NAMESPACE=${1:-kubevirt}

kubectl apply -n "${KUBEVIRT_NAMESPACE}" -f "https://github.com/kubevirt/kubevirt/releases/download/${KUBEVIRT_VERSION}/kubevirt-operator.yaml"
kubectl apply -n "${KUBEVIRT_NAMESPACE}" -f "https://github.com/kubevirt/kubevirt/releases/download/${KUBEVIRT_VERSION}/kubevirt-cr.yaml"

echo "Waiting for Kubevirt to be ready..."
kubectl wait --for=condition=Available --timeout=600s -n "${KUBEVIRT_NAMESPACE}" kv/kubevirt


# Deploying CDI
CDI_VERSION=$(latest_version "containerized-data-importer")
CDI_NAMESPACE=cdi

kubectl apply -n ${CDI_NAMESPACE} -f "https://github.com/kubevirt/containerized-data-importer/releases/download/${CDI_VERSION}/cdi-operator.yaml"
kubectl apply -n ${CDI_NAMESPACE} -f "https://github.com/kubevirt/containerized-data-importer/releases/download/${CDI_VERSION}/cdi-cr.yaml"

echo "Waiting for CDI to be ready..."
kubectl wait --for=condition=Available --timeout=600s -n ${CDI_NAMESPACE} cdi/cdi

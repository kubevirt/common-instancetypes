#!/bin/bash
#
# Copyright 2023 Red Hat, Inc.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

set -ex

# Required by the VMPersistentState feature that currently needs a RWX FS mode SC
# https://kubevirt.io/user-guide/virtual_machines/persistent_tpm_and_uefi_state/
export KUBEVIRT_DEPLOY_NFS_CSI="${KUBEVIRT_DEPLOY_NFS_CSI:-true}"
export KUBEVIRT_MEMORY_SIZE="${KUBEVIRT_MEMORY_SIZE:-16G}"
export KUBEVIRTCI_TAG=${KUBEVIRTCI_TAG:-$(curl -sfL https://storage.googleapis.com/kubevirt-prow/release/kubevirt/kubevirtci/latest)}
export KUBEVIRT_DEPLOY_CDI="true"
export KUBEVIRT_VERSION=${KUBEVIRT_VERSION:-main)}

_base_dir=$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)
_cluster_up_dir="${_base_dir}/_cluster-up"
_virtctl_dir="${_cluster_up_dir}/_out/cmd/virtctl"
_kubectl="${_cluster_up_dir}/cluster-up/kubectl.sh"
_kubessh="${_cluster_up_dir}/cluster-up/ssh.sh"
_kubevirtcicli="${_cluster_up_dir}/cluster-up/cli.sh"
_virtctl="${_cluster_up_dir}/cluster-up/virtctl.sh"
_action=$1
shift

function kubevirtci::fetch_kubevirtci() {
	if [[ ! -d ${_cluster_up_dir} ]]; then
    git clone --depth 1 --branch "${KUBEVIRTCI_TAG}" https://github.com/kubevirt/kubevirtci.git "${_cluster_up_dir}"
  fi
}

function kubevirtci::up() {
  make cluster-up -C "${_cluster_up_dir}"
  KUBECONFIG=$(kubevirtci::kubeconfig)
  export KUBECONFIG

  echo "adding kubevirtci registry to cdi-insecure-registries"
  ${_kubectl} patch cdis/cdi --type merge -p '{"spec": {"config": {"insecureRegistries": ["registry:5000"]}}}'

  # Treat main as stable version
  if [ "$KUBEVIRT_VERSION" = "main" ]; then
    KUBEVIRT_VERSION=$(curl -L https://storage.googleapis.com/kubevirt-prow/devel/release/kubevirt/kubevirt/stable.txt)
  fi

  echo "installing kubevirt..."
  ${_kubectl} apply -f "https://github.com/kubevirt/kubevirt/releases/download/${KUBEVIRT_VERSION}/kubevirt-operator.yaml"
  ${_kubectl} apply -f "https://github.com/kubevirt/kubevirt/releases/download/${KUBEVIRT_VERSION}/kubevirt-cr.yaml"

  # This is needed because kubevirtci does not provide its own virtctl binary
  echo "installing virtctl..."
  mkdir -p "${_virtctl_dir}"
  curl -L "https://github.com/kubevirt/kubevirt/releases/download/${KUBEVIRT_VERSION}/virtctl-${KUBEVIRT_VERSION}-linux-amd64" -o "${_virtctl_dir}/virtctl"
  chmod +x "${_virtctl_dir}/virtctl"

  echo "waiting for kubevirt to become ready, this can take a few minutes..."
  ${_kubectl} -n kubevirt wait kv kubevirt --for condition=Available --timeout=15m

  echo "enabling feature gates to validate instance types and preferences"
  ${_kubectl} patch kv/kubevirt -n kubevirt --type merge -p '{"spec":{"configuration":{"developerConfiguration":{"featureGates": ["GPU", "NUMA", "VMPersistentState"]},"vmStateStorageClass":"nfs-csi"}}}'
}

function kubevirtci::down() {
  make cluster-down -C "${_cluster_up_dir}"
}

function kubevirtci::sync() {
  KUBECTL=${_kubectl} BASEDIR=${_base_dir} "${_base_dir}/scripts/sync.sh"
}

function kubevirtci::sync-containerdisks() {
  podman pull "${VALIDATION_OS_IMAGE}:${VALIDATION_OS_IMAGE_TAG}"
  podman push --tls-verify=false \
    "${VALIDATION_OS_IMAGE}:${VALIDATION_OS_IMAGE_TAG}" \
    "$(kubevirtci::registry)/validation-os-container-disk:latest"
}

function kubevirtci::kubeconfig() {
  "${_cluster_up_dir}/cluster-up/kubeconfig.sh"
}

function kubevirtci::registry() {
  port=$(${_kubevirtcicli} ports registry 2>/dev/null)
  echo "localhost:${port}"
}

function kubevirtci::functest() {
  KUBECTL=${_kubectl} VIRTCTL=${_virtctl} "${_base_dir}/scripts/functest.sh"
}

kubevirtci::fetch_kubevirtci

case ${_action} in
  "up")
    kubevirtci::up
    ;;
  "down")
    kubevirtci::down
    ;;
  "sync")
    kubevirtci::sync
    ;;
  "sync-containerdisks")
    kubevirtci::sync-containerdisks
    ;;
  "kubeconfig")
    kubevirtci::kubeconfig
    ;;
  "registry")
    kubevirtci::registry
    ;;
  "functest")
    kubevirtci::functest
    ;;
  "ssh")
    ${_kubessh} "$@"
    ;;
  "kubectl")
    ${_kubectl} "$@"
    ;;
  "virtctl")
    ${_virtctl} "$@"
    ;;
  *)
    echo "No command provided, known commands are 'up', 'down', 'sync', 'ssh', 'kubeconfig', 'registry', 'kubectl', 'functest'"
    exit 1
    ;;
esac

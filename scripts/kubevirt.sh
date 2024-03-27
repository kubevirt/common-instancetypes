#!/bin/bash
#
# Copyright 2022 Red Hat, Inc.
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
export KUBEVIRT_TAG="${KUBEVIRT_TAG:-main}"

_base_dir=$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)
_kubectl="${_base_dir}/_kubevirt/cluster-up/kubectl.sh"
_kubessh="${_base_dir}/_kubevirt/cluster-up/ssh.sh"
_kubevirtcicli="${_base_dir}/_kubevirt/cluster-up/cli.sh"
_virtctl="${_base_dir}/_kubevirt/cluster-up/virtctl.sh"
_action=$1
shift

function kubevirt::install() {
  if [[ ! -d ${_base_dir}/_kubevirt ]]; then
    git clone --depth 1 --branch "${KUBEVIRT_TAG}" https://github.com/kubevirt/kubevirt.git "${_base_dir}/_kubevirt"
  fi
}

function kubevirt::up() {
  make cluster-up -C "${_base_dir}/_kubevirt" && make cluster-sync -C "${_base_dir}/_kubevirt"

  echo "enabling feature gates to validate instance types and preferences"
  ${_kubectl} patch kv/kubevirt -n kubevirt --type merge -p '{"spec":{"configuration":{"developerConfiguration":{"featureGates": ["GPU", "NUMA", "VMPersistentState"]},"vmStateStorageClass":"nfs-csi"}}}'
}

function kubevirt::down() {
  make cluster-down -C "${_base_dir}/_kubevirt"
}

function kubevirt::sync() {
  KUBECTL=${_kubectl} "${_base_dir}/scripts/sync.sh"
}

function kubevirt::sync-containerdisks() {
  podman pull "${VALIDATION_OS_IMAGE}:${VALIDATION_OS_IMAGE_TAG}"
  podman push --tls-verify=false \
    "${VALIDATION_OS_IMAGE}:${VALIDATION_OS_IMAGE_TAG}" \
    "$(kubevirt::registry)/validation-os-container-disk:latest"
}

function kubevirt::kubeconfig() {
  "${_base_dir}/_kubevirt/cluster-up/kubeconfig.sh"
}

function kubevirt::registry() {
  port=$(${_kubevirtcicli} ports registry 2>/dev/null)
  echo "localhost:${port}"
}

function kubevirt::functest() {
  KUBECTL=${_kubectl} VIRTCTL=${_virtctl} "${_base_dir}/scripts/functest.sh"
}

kubevirt::install

case ${_action} in
  "up")
    kubevirt::up
    ;;
  "down")
    kubevirt::down
    ;;
  "sync")
    kubevirt::sync
    ;;
  "sync-containerdisks")
    kubevirt::sync-containerdisks
    ;;
  "kubeconfig")
    kubevirt::kubeconfig
    ;;
  "registry")
    kubevirt::registry
    ;;
  "functest")
    kubevirt::functest
    ;;
  "ssh")
    ${_kubessh} "$@"
    ;;
  "kubectl")
    ${_kubectl} "$@"
    ;;
  *)
    echo "No command provided, known commands are 'up', 'down', 'sync', 'ssh', 'kubeconfig', 'registry', 'kubectl', 'functest'"
    exit 1
    ;;
esac


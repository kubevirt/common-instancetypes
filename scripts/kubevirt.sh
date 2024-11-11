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
_kubectl="${_base_dir}/_kubevirt/kubevirtci/cluster-up/kubectl.sh"
_kubessh="${_base_dir}/_kubevirt/kubevirtci/cluster-up/ssh.sh"
_kubevirtcicli="${_base_dir}/_kubevirt/kubevirtci/cluster-up/cli.sh"
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
  ${_kubectl} patch kv/kubevirt -n kubevirt --type merge -p '{"spec":{"configuration":{"commonInstancetypesDeployment": {"enabled": false},"developerConfiguration":{"featureGates": ["VMPersistentState"]},"vmStateStorageClass":"nfs-csi"}}}'

  echo "waiting for kubevirt to become ready, this can take a few minutes..."
  ${_kubectl} -n kubevirt wait kv kubevirt --for condition=Available --timeout=15m
}

function kubevirt::down() {
  make cluster-down -C "${_base_dir}/_kubevirt"
}

function kubevirt::sync() {
  KUBECTL=${_kubectl} BASEDIR=${_base_dir} "${_base_dir}/scripts/sync.sh"
}

function kubevirt::sync-containerdisks() {
  podman pull "${VALIDATION_OS_IMAGE}:${VALIDATION_OS_IMAGE_TAG}"
  podman push --tls-verify=false \
    "${VALIDATION_OS_IMAGE}:${VALIDATION_OS_IMAGE_TAG}" \
    "$(kubevirt::registry)/validation-os-container-disk:latest"
}

function kubevirt::kubeconfig() {
  "${_base_dir}/_kubevirt/kubevirtci/cluster-up/kubeconfig.sh"
}

function kubevirt::registry() {
  port=$(${_kubevirtcicli} ports registry 2>/dev/null)
  echo "localhost:${port}"
}

kubevirt::install

cd "${_base_dir}"/_kubevirt

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
  "ssh")
    ${_kubessh} "$@"
    ;;
  "kubectl")
    ${_kubectl} "$@"
    ;;
  *)
    echo "No command provided, known commands are 'up', 'down', 'sync', 'ssh', 'kubeconfig', 'registry', 'kubectl'"
    exit 1
    ;;
esac

cd "${_base_dir}"


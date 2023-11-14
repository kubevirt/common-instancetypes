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

if [ -z "${KUBECTL}" ] || [ -z "${VIRTCTL}" ]; then
    echo "${BASH_SOURCE[0]} expects the following env variables to be provided: KUBECTL, VIRTCTL."
    exit 1
fi

# Fail if no instancetypes or preferences were found
if [ "$(${KUBECTL} get virtualmachineclusterinstancetype -o json | jq '.items | length')" -le 0 ]; then
  echo "functest failed - No VirtualMachineClusterInstancetypes found"
  exit 1
fi
if [ "$(${KUBECTL} get virtualmachineclusterpreference -o json | jq '.items | length')" -le 0 ]; then
  echo "functest failed - No VirtualMachineClusterPreferences found"
  exit 1
fi

# Create a custom tiny instance type for negative tests around preference resource requirements
${VIRTCTL} create instancetype --cpu 1 --memory 64Mi --name tiny | ${KUBECTL} apply -f -

# This func test simply loops over the installed instance types and preferences, assigning each to a VirtualMachine to ensure they are accepted by the webhooks
for preference in $(${KUBECTL} get virtualmachineclusterpreferences --no-headers -o custom-columns=':metadata.name'); do

    # Ensure a VirtualMachine using a preference with resource requirements is rejected if it does not provide enough resources.
    if ${KUBECTL} get virtualmachineclusterpreferences/"${preference}" -o json | jq -er .spec.requirements > /dev/null 2>&1; then
        # TODO(lyarwood): virtctl should be extended with a --cpu switch to allow the non instancetype use case to be tested here
        if ${VIRTCTL} create vm --instancetype tiny --preference "${preference}" --volume-containerdisk name:disk,src:quay.io/containerdisks/fedora:latest --name "vm-${preference}-requirements" | ${KUBECTL} apply -f - ; then
            echo "functest failed - Preference ${preference} should not be able to use virtualmachineclusterinstancetype tiny"
            ${KUBECTL} delete "vm/vm-${preference}-requirements"
            exit 1
        fi
    fi

    # Ensure a VirtualMachine can be created when enough resources are provided using the u1.medium instance type
    if ! ${VIRTCTL} create vm --instancetype u1.medium --preference "${preference}" --volume-containerdisk name:disk,src:quay.io/containerdisks/fedora:latest --name "vm-${preference}" | ${KUBECTL} apply -f - ; then
        echo "functest failed on preference ${preference} using instancetype u1.medium"
        exit 1
    fi

    if ! ${KUBECTL} wait --for=condition=Ready vms/"vm-${preference}" ; then
        echo "functest failed to wait for VirtualMachine to become Ready"
        ${KUBECTL} get vms/"vm-${preference}" -o yaml
        exit 1
    fi

    ${KUBECTL} delete "vm/vm-${preference}"
done

# Cleanup the custom instancetype
${KUBECTL} delete virtualmachineclusterinstancetypes/tiny

for instancetype in $(${KUBECTL} get virtualmachineclusterinstancetypes --no-headers -o custom-columns=':metadata.name'); do
    if ! ${VIRTCTL} create vm --instancetype "${instancetype}" --volume-containerdisk name:disk,src:quay.io/containerdisks/fedora:latest --name "vm-${instancetype}" | ${KUBECTL} apply -f - ; then
        echo "functest failed on instance type ${instancetype}"
        exit 1
    fi
    ${KUBECTL} delete "vm/vm-${instancetype}"
done

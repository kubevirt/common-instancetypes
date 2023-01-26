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

if [ -z "${KUBECTL}" ]; then
    echo "${BASH_SOURCE[0]} expects the following env variables to be provided: KUBECTL."
    exit 1
fi

# This func test simply loops over the installed instance types and preferences, assigning each to a VirtualMachine to ensure they are accepted by the webhooks
for preference in $(${KUBECTL} get virtualmachineclusterpreferences --no-headers -o custom-columns=':metadata.name'); do
    # TODO(lyarwood): Replace with virtctl create vm once 0.59.0 is released
    # ${VIRTCTL} create vm --preference ${preference}
    ${KUBECTL} apply -f - << EOF
---    
apiVersion: kubevirt.io/v1
kind: VirtualMachine
metadata:
  name: vm-${preference}
spec:
  preference:
    name: ${preference}
  running: false
  template:
    spec:
      domain:
        devices: {}
      volumes:
      - containerDisk:
          image: quay.io/containerdisks/fedora:37
        name: containerdisk
EOF
    # We can't inline the above creation call below so stash the return code and check it to keep shellcheck happy
    ret=$?
    if [ $ret -ne 0 ]; then
        echo "functest failed on preference ${preference}"
        exit 1
    fi
    ${KUBECTL} delete "vm/vm-${preference}"
done

for instancetype in $(${KUBECTL} get virtualmachineclusterinstancetypes --no-headers -o custom-columns=':metadata.name'); do
    # TODO(lyarwood): Replace with virtctl create vm once 0.59.0 is released
    # ${VIRTCTL} create vm --instance-type ${instancetype}
    ${KUBECTL} apply -f - << EOF
---    
apiVersion: kubevirt.io/v1
kind: VirtualMachine
metadata:
  name: vm-${instancetype}
spec:
  instancetype:
    name: ${instancetype}
  running: false
  template:
    spec:
      domain:
        devices: {}
      volumes:
      - containerDisk:
          image: quay.io/containerdisks/fedora:37
        name: containerdisk
EOF
    # We can't inline the above creation call below so stash the return code and check it to keep shellcheck happy
    ret=$?
    if [ $ret -ne 0 ]; then
        echo "functest failed on instance type ${instancetype}"
        exit 1
    fi
    ${KUBECTL} delete "vm/vm-${instancetype}"
done

# kubevirt/common-instancetypes

A set of instancetypes and preferences to help create KubeVirt VirtualMachines.

Currently providing the following VirtualMachineClusterInstancetype, VirtualMachineClusterPreference, VirtualMachineInstancetype and VirtualMachinePreference resources:

```
$ ./cluster-up/kubectl.sh kustomize https://github.com/lyarwood/common-instancetypes | ./cluster-up/kubectl.sh apply -f -
selecting docker as container runtime
selecting docker as container runtime
virtualmachineclusterinstancetype.instancetype.kubevirt.io/highperformance.large created
virtualmachineclusterinstancetype.instancetype.kubevirt.io/highperformance.medium created
virtualmachineclusterinstancetype.instancetype.kubevirt.io/highperformance.small created
virtualmachineclusterinstancetype.instancetype.kubevirt.io/highperformance.tiny created
virtualmachineclusterinstancetype.instancetype.kubevirt.io/server.large created
virtualmachineclusterinstancetype.instancetype.kubevirt.io/server.medium created
virtualmachineclusterinstancetype.instancetype.kubevirt.io/server.small created
virtualmachineclusterinstancetype.instancetype.kubevirt.io/server.tiny created
virtualmachineclusterpreference.instancetype.kubevirt.io/centos.7 created
virtualmachineclusterpreference.instancetype.kubevirt.io/centos.7.desktop created
virtualmachineclusterpreference.instancetype.kubevirt.io/centos.8 created
virtualmachineclusterpreference.instancetype.kubevirt.io/centos.8.desktop created
virtualmachineclusterpreference.instancetype.kubevirt.io/centos.9 created
virtualmachineclusterpreference.instancetype.kubevirt.io/centos.9.desktop created
virtualmachineclusterpreference.instancetype.kubevirt.io/fedora.35 created
virtualmachineclusterpreference.instancetype.kubevirt.io/fedora.36 created
virtualmachineclusterpreference.instancetype.kubevirt.io/rhel.7 created
virtualmachineclusterpreference.instancetype.kubevirt.io/rhel.7.desktop created
virtualmachineclusterpreference.instancetype.kubevirt.io/rhel.8 created
virtualmachineclusterpreference.instancetype.kubevirt.io/rhel.8.desktop created
virtualmachineclusterpreference.instancetype.kubevirt.io/rhel.9 created
virtualmachineclusterpreference.instancetype.kubevirt.io/rhel.9.desktop created
virtualmachineclusterpreference.instancetype.kubevirt.io/ubuntu.18.04 created
virtualmachineclusterpreference.instancetype.kubevirt.io/ubuntu.20.04 created
virtualmachineclusterpreference.instancetype.kubevirt.io/ubuntu.22.04 created
virtualmachineclusterpreference.instancetype.kubevirt.io/windows.10 created
virtualmachineclusterpreference.instancetype.kubevirt.io/windows.10.virtio created
virtualmachineclusterpreference.instancetype.kubevirt.io/windows.11 created
virtualmachineclusterpreference.instancetype.kubevirt.io/windows.11.virtio created
virtualmachineclusterpreference.instancetype.kubevirt.io/windows.2k12 created
virtualmachineclusterpreference.instancetype.kubevirt.io/windows.2k12.virtio created
virtualmachineclusterpreference.instancetype.kubevirt.io/windows.2k16 created
virtualmachineclusterpreference.instancetype.kubevirt.io/windows.2k16.virtio created
virtualmachineclusterpreference.instancetype.kubevirt.io/windows.2k19 created
virtualmachineclusterpreference.instancetype.kubevirt.io/windows.2k19.virtio created
virtualmachineinstancetype.instancetype.kubevirt.io/highperformance.large created
virtualmachineinstancetype.instancetype.kubevirt.io/highperformance.medium created
virtualmachineinstancetype.instancetype.kubevirt.io/highperformance.small created
virtualmachineinstancetype.instancetype.kubevirt.io/highperformance.tiny created
virtualmachineinstancetype.instancetype.kubevirt.io/server.large created
virtualmachineinstancetype.instancetype.kubevirt.io/server.medium created
virtualmachineinstancetype.instancetype.kubevirt.io/server.small created
virtualmachineinstancetype.instancetype.kubevirt.io/server.tiny created
virtualmachinepreference.instancetype.kubevirt.io/centos.7 created
virtualmachinepreference.instancetype.kubevirt.io/centos.7.desktop created
virtualmachinepreference.instancetype.kubevirt.io/centos.8 created
virtualmachinepreference.instancetype.kubevirt.io/centos.8.desktop created
virtualmachinepreference.instancetype.kubevirt.io/centos.9 created
virtualmachinepreference.instancetype.kubevirt.io/centos.9.desktop created
virtualmachinepreference.instancetype.kubevirt.io/fedora.35 created
virtualmachinepreference.instancetype.kubevirt.io/fedora.36 created
virtualmachinepreference.instancetype.kubevirt.io/rhel.7 created
virtualmachinepreference.instancetype.kubevirt.io/rhel.7.desktop created
virtualmachinepreference.instancetype.kubevirt.io/rhel.8 created
virtualmachinepreference.instancetype.kubevirt.io/rhel.8.desktop created
virtualmachinepreference.instancetype.kubevirt.io/rhel.9 created
virtualmachinepreference.instancetype.kubevirt.io/rhel.9.desktop created
virtualmachinepreference.instancetype.kubevirt.io/ubuntu.18.04 created
virtualmachinepreference.instancetype.kubevirt.io/ubuntu.20.04 created
virtualmachinepreference.instancetype.kubevirt.io/ubuntu.22.04 created
virtualmachinepreference.instancetype.kubevirt.io/windows.10 created
virtualmachinepreference.instancetype.kubevirt.io/windows.10.virtio created
virtualmachinepreference.instancetype.kubevirt.io/windows.11 created
virtualmachinepreference.instancetype.kubevirt.io/windows.11.virtio created
virtualmachinepreference.instancetype.kubevirt.io/windows.2k12 created
virtualmachinepreference.instancetype.kubevirt.io/windows.2k12.virtio created
virtualmachinepreference.instancetype.kubevirt.io/windows.2k16 created
virtualmachinepreference.instancetype.kubevirt.io/windows.2k16.virtio created
virtualmachinepreference.instancetype.kubevirt.io/windows.2k19 created
virtualmachinepreference.instancetype.kubevirt.io/windows.2k19.virtio created
```
## Requirements

- kustomize
- yamllint

## Build

```
$ make
yamllint .
echo "---" > common-instancetypes.yaml
kustomize build >> common-instancetypes.yaml
```
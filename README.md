# kubevirt/common-instancetypes

A set of [instance types and preferences](https://kubevirt.io/user-guide/virtual_machines/instancetypes/) to help create [`KubeVirt`](https://kubevirt.io)  [`VirtualMachines`](http://kubevirt.io/api-reference/main/definitions.html#_v1alpha1_virtualmachine).

## Installation

[kustomize](https://kustomize.io/) is used to generate these resources before applying them to the cluster with [`kubectl`](https://kubernetes.io/docs/reference/kubectl/):

```
$ kustomize build | kubectl apply -f -
```

A number of resource specific [kustomize](https://kustomize.io/) build configs are also provided and allow limited generation, for example the following command will generate and install only the [`VirtualMachineInstancetype`](http://kubevirt.io/api-reference/main/definitions.html#_v1alpha2_virtualmachineinstancetype) resources:

```
$ kustomize build VirtualMachineInstancetypes | kubectl apply -f -
```

Additional targets are avilable for [`VirtualMachineClusterInstancetypes`](http://kubevirt.io/api-reference/main/definitions.html#_v1alpha2_virtualmachineclusterinstancetype), [`VirtualMachineClusterPreferences`](http://kubevirt.io/api-reference/main/definitions.html#_v1alpha2_virtualmachineclusterpreference) and [`VirtualMachinePreferences`](http://kubevirt.io/api-reference/main/definitions.html#_v1alpha2_virtualmachinepreference).

## Requirements

The following tools are required by each `Makefile` target:

* `make`
  * `make lint`
    * [yamllint](https://github.com/adrienverge/yamllint)
    * [shellcheck](https://github.com/koalaman/shellcheck)
  * `make generate`
    * [kustomize](https://kustomize.io/)
  * `make validate`
    * [kustomize](https://kustomize.io/)
    * [kubeconform](https://github.com/yannh/kubeconform)
  * `make readme`
    * [yq](https://github.com/mikefarah/yq)
* `make schema` (optional)
  * [openapi2jsonschema](https://github.com/instrumenta/openapi2jsonschema)

## Running Makefile targets within a container

A [`Dockerfile`](./Dockerfile) is available that provides all of the above dependencies.

To build and use this simply set the `COMMON_INSTANCETYPES_CRI` env variable to `docker` or `podman` before running `make build_image`:

```
$ COMMON_INSTANCETYPES_CRI=docker make build_image
$ COMMON_INSTANCETYPES_CRI=docker make
```

## Resources

The following instance type and preference resources are provided by this project:

```

VirtualMachineClusterInstancetype
cx1.2xlarge
---
VirtualMachineClusterInstancetype
cx1.4xlarge
---
VirtualMachineClusterInstancetype
cx1.8xlarge
---
VirtualMachineClusterInstancetype
cx1.large
---
VirtualMachineClusterInstancetype
cx1.medium
---
VirtualMachineClusterInstancetype
cx1.xlarge
---
VirtualMachineClusterInstancetype
gn1.2xlarge
---
VirtualMachineClusterInstancetype
gn1.4xlarge
---
VirtualMachineClusterInstancetype
gn1.8xlarge
---
VirtualMachineClusterInstancetype
gn1.xlarge
---
VirtualMachineClusterInstancetype
m1.2xlarge
---
VirtualMachineClusterInstancetype
m1.4xlarge
---
VirtualMachineClusterInstancetype
m1.8xlarge
---
VirtualMachineClusterInstancetype
m1.large
---
VirtualMachineClusterInstancetype
m1.xlarge
---
VirtualMachineClusterInstancetype
o1.2xlarge
---
VirtualMachineClusterInstancetype
o1.4xlarge
---
VirtualMachineClusterInstancetype
o1.8xlarge
---
VirtualMachineClusterInstancetype
o1.large
---
VirtualMachineClusterInstancetype
o1.medium
---
VirtualMachineClusterInstancetype
o1.micro
---
VirtualMachineClusterInstancetype
o1.nano
---
VirtualMachineClusterInstancetype
o1.small
---
VirtualMachineClusterInstancetype
o1.xlarge
---
VirtualMachineClusterInstancetype
u1.2xlarge
---
VirtualMachineClusterInstancetype
u1.4xlarge
---
VirtualMachineClusterInstancetype
u1.8xlarge
---
VirtualMachineClusterInstancetype
u1.large
---
VirtualMachineClusterInstancetype
u1.medium
---
VirtualMachineClusterInstancetype
u1.micro
---
VirtualMachineClusterInstancetype
u1.nano
---
VirtualMachineClusterInstancetype
u1.small
---
VirtualMachineClusterInstancetype
u1.xlarge
---
VirtualMachineClusterPreference
alpine
---
VirtualMachineClusterPreference
centos.7
---
VirtualMachineClusterPreference
centos.7.desktop
---
VirtualMachineClusterPreference
centos.stream8
---
VirtualMachineClusterPreference
centos.stream8.desktop
---
VirtualMachineClusterPreference
centos.stream9
---
VirtualMachineClusterPreference
centos.stream9.desktop
---
VirtualMachineClusterPreference
cirros
---
VirtualMachineClusterPreference
fedora
---
VirtualMachineClusterPreference
rhel.7
---
VirtualMachineClusterPreference
rhel.7.desktop
---
VirtualMachineClusterPreference
rhel.8
---
VirtualMachineClusterPreference
rhel.8.desktop
---
VirtualMachineClusterPreference
rhel.9
---
VirtualMachineClusterPreference
rhel.9.desktop
---
VirtualMachineClusterPreference
ubuntu
---
VirtualMachineClusterPreference
windows.10
---
VirtualMachineClusterPreference
windows.10.virtio
---
VirtualMachineClusterPreference
windows.11
---
VirtualMachineClusterPreference
windows.11.virtio
---
VirtualMachineClusterPreference
windows.2k12
---
VirtualMachineClusterPreference
windows.2k12.virtio
---
VirtualMachineClusterPreference
windows.2k16
---
VirtualMachineClusterPreference
windows.2k16.virtio
---
VirtualMachineClusterPreference
windows.2k19
---
VirtualMachineClusterPreference
windows.2k19.virtio
---
VirtualMachineClusterPreference
windows.2k22
---
VirtualMachineClusterPreference
windows.2k22.virtio
---
VirtualMachineInstancetype
cx1.2xlarge
---
VirtualMachineInstancetype
cx1.4xlarge
---
VirtualMachineInstancetype
cx1.8xlarge
---
VirtualMachineInstancetype
cx1.large
---
VirtualMachineInstancetype
cx1.medium
---
VirtualMachineInstancetype
cx1.xlarge
---
VirtualMachineInstancetype
gn1.2xlarge
---
VirtualMachineInstancetype
gn1.4xlarge
---
VirtualMachineInstancetype
gn1.8xlarge
---
VirtualMachineInstancetype
gn1.xlarge
---
VirtualMachineInstancetype
m1.2xlarge
---
VirtualMachineInstancetype
m1.4xlarge
---
VirtualMachineInstancetype
m1.8xlarge
---
VirtualMachineInstancetype
m1.large
---
VirtualMachineInstancetype
m1.xlarge
---
VirtualMachineInstancetype
o1.2xlarge
---
VirtualMachineInstancetype
o1.4xlarge
---
VirtualMachineInstancetype
o1.8xlarge
---
VirtualMachineInstancetype
o1.large
---
VirtualMachineInstancetype
o1.medium
---
VirtualMachineInstancetype
o1.micro
---
VirtualMachineInstancetype
o1.nano
---
VirtualMachineInstancetype
o1.small
---
VirtualMachineInstancetype
o1.xlarge
---
VirtualMachineInstancetype
u1.2xlarge
---
VirtualMachineInstancetype
u1.4xlarge
---
VirtualMachineInstancetype
u1.8xlarge
---
VirtualMachineInstancetype
u1.large
---
VirtualMachineInstancetype
u1.medium
---
VirtualMachineInstancetype
u1.micro
---
VirtualMachineInstancetype
u1.nano
---
VirtualMachineInstancetype
u1.small
---
VirtualMachineInstancetype
u1.xlarge
---
VirtualMachinePreference
alpine
---
VirtualMachinePreference
centos.7
---
VirtualMachinePreference
centos.7.desktop
---
VirtualMachinePreference
centos.stream8
---
VirtualMachinePreference
centos.stream8.desktop
---
VirtualMachinePreference
centos.stream9
---
VirtualMachinePreference
centos.stream9.desktop
---
VirtualMachinePreference
cirros
---
VirtualMachinePreference
fedora
---
VirtualMachinePreference
rhel.7
---
VirtualMachinePreference
rhel.7.desktop
---
VirtualMachinePreference
rhel.8
---
VirtualMachinePreference
rhel.8.desktop
---
VirtualMachinePreference
rhel.9
---
VirtualMachinePreference
rhel.9.desktop
---
VirtualMachinePreference
ubuntu
---
VirtualMachinePreference
windows.10
---
VirtualMachinePreference
windows.10.virtio
---
VirtualMachinePreference
windows.11
---
VirtualMachinePreference
windows.11.virtio
---
VirtualMachinePreference
windows.2k12
---
VirtualMachinePreference
windows.2k12.virtio
---
VirtualMachinePreference
windows.2k16
---
VirtualMachinePreference
windows.2k16.virtio
---
VirtualMachinePreference
windows.2k19
---
VirtualMachinePreference
windows.2k19.virtio
---
VirtualMachinePreference
windows.2k22
---
VirtualMachinePreference
windows.2k22.virtio

```

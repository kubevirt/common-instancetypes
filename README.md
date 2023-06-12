# kubevirt/common-instancetypes

A set of [instance types and preferences](https://kubevirt.io/user-guide/virtual_machines/instancetypes/) to help create [`KubeVirt`](https://kubevirt.io)  [`VirtualMachines`](http://kubevirt.io/api-reference/main/definitions.html#_v1alpha1_virtualmachine).

## Installation

Bundles of pre-generated instance type and preference resources are available within this repository and can be applied
directly to an existing [`Kubernetes`](https://kubernetes.io/) cluster with [`KubeVirt`](https://kubevirt.io) installed using [`kubectl`](https://kubernetes.io/docs/reference/kubectl/):

```
$ kubectl apply -f ./common-instancetypes-all-bundle.yaml
```

Additionally [kustomize](https://kustomize.io/) can be used to generate these resources before applying them to the cluster with [`kubectl`](https://kubernetes.io/docs/reference/kubectl/)

```
$ kustomize build ./ | kubectl apply -f -
```

A number of resource specific [kustomize](https://kustomize.io/) build configs are also provided and allow limited generation, for example the following command will generate and install only the [`VirtualMachineInstancetype`](http://kubevirt.io/api-reference/main/definitions.html#_v1alpha2_virtualmachineinstancetype) resources:

```
$ kustomize build ./VirtualMachineInstancetypes | kubectl apply -f -
```

Additional targets are avilable for [`./VirtualMachineClusterInstancetypes`](http://kubevirt.io/api-reference/main/definitions.html#_v1alpha2_virtualmachineclusterinstancetype), [`./VirtualMachineClusterPreferences`](http://kubevirt.io/api-reference/main/definitions.html#_v1alpha2_virtualmachineclusterpreference) and [`./VirtualMachinePreferences`](http://kubevirt.io/api-reference/main/definitions.html#_v1alpha2_virtualmachinepreference).

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
highperformance.large
---
VirtualMachineInstancetype
highperformance.medium
---
VirtualMachineInstancetype
highperformance.small
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
n1.2xlarge
---
VirtualMachineInstancetype
n1.4xlarge
---
VirtualMachineInstancetype
n1.8xlarge
---
VirtualMachineInstancetype
n1.large
---
VirtualMachineInstancetype
n1.medium
---
VirtualMachineInstancetype
n1.xlarge
---
VirtualMachineInstancetype
server.large
---
VirtualMachineInstancetype
server.medium
---
VirtualMachineInstancetype
server.micro
---
VirtualMachineInstancetype
server.small
---
VirtualMachineInstancetype
server.tiny
---
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
highperformance.large
---
VirtualMachineClusterInstancetype
highperformance.medium
---
VirtualMachineClusterInstancetype
highperformance.small
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
n1.2xlarge
---
VirtualMachineClusterInstancetype
n1.4xlarge
---
VirtualMachineClusterInstancetype
n1.8xlarge
---
VirtualMachineClusterInstancetype
n1.large
---
VirtualMachineClusterInstancetype
n1.medium
---
VirtualMachineClusterInstancetype
n1.xlarge
---
VirtualMachineClusterInstancetype
server.large
---
VirtualMachineClusterInstancetype
server.medium
---
VirtualMachineClusterInstancetype
server.micro
---
VirtualMachineClusterInstancetype
server.small
---
VirtualMachineClusterInstancetype
server.tiny
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
centos.8.stream
---
VirtualMachineClusterPreference
centos.8.stream.desktop
---
VirtualMachineClusterPreference
centos.9.stream
---
VirtualMachineClusterPreference
centos.9.stream.desktop
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
centos.8.stream
---
VirtualMachinePreference
centos.8.stream.desktop
---
VirtualMachinePreference
centos.9.stream
---
VirtualMachinePreference
centos.9.stream.desktop
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

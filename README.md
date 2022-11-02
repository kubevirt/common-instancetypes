# kubevirt/common-instancetypes

A set of instance types and preferences to help create KubeVirt VirtualMachines.

## Requirements

The following tools are required by each `Makefile` target:

* `make`
  * `make lint`
    * [yamllint](https://github.com/adrienverge/yamllint)
    * [bashate](https://github.com/openstack/bashate)
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

To build and use this simply set the `COMMON_INSTANCETYPES_CRI` env variable to `docker` or `podman` before running `make`:

```
$ COMMON_INSTANCETYPES_CRI=docker make
```

Alternatively you can also build the image and run individual Makefile targets:

```
$ COMMON_INSTANCETYPES_CRI=docker make build_image
[..]
$ COMMON_INSTANCETYPES_CRI=docker make lint
```

## Resources

The following instance type and preference resources are provided by this project:

```

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
server.large
---
VirtualMachineClusterInstancetype
server.medium
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
centos.8
---
VirtualMachineClusterPreference
centos.8.desktop
---
VirtualMachineClusterPreference
centos.9
---
VirtualMachineClusterPreference
centos.9.desktop
---
VirtualMachineClusterPreference
cirros
---
VirtualMachineClusterPreference
fedora.35
---
VirtualMachineClusterPreference
fedora.36
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
ubuntu.18.04
---
VirtualMachineClusterPreference
ubuntu.20.04
---
VirtualMachineClusterPreference
ubuntu.22.04
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
server.large
---
VirtualMachineInstancetype
server.medium
---
VirtualMachineInstancetype
server.small
---
VirtualMachineInstancetype
server.tiny
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
centos.8
---
VirtualMachinePreference
centos.8.desktop
---
VirtualMachinePreference
centos.9
---
VirtualMachinePreference
centos.9.desktop
---
VirtualMachinePreference
cirros
---
VirtualMachinePreference
fedora.35
---
VirtualMachinePreference
fedora.36
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
ubuntu.18.04
---
VirtualMachinePreference
ubuntu.20.04
---
VirtualMachinePreference
ubuntu.22.04
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

```

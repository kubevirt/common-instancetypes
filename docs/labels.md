# Labels

Various labels are used to decorate objects provided by this project to help end users find the right resources for them.

For example, to list all cluster instance types deployed in the current env that provide `4` vCPUs a user can use the following command:

```shell
$ kubectl get virtualmachineclusterinstancetype -linstancetype.kubevirt.io/cpu=4
NAME         AGE
cx1.xlarge   39s
gn1.xlarge   39s
m1.xlarge    39s
n1.large     39s
n1.medium    39s
o1.xlarge    39s
u1.xlarge    39s
```

## Generic

### `instancetype.kubevirt.io/common-instancetypes-version`

The version of the common-instancetypes project used to generate these resources.

### `instancetype.kubevirt.io/vendor`

The vendor of the resource, this is always `kubevirt.io` upstream and should be changed by downstream vendors consuming the project.

### `instancetype.kubevirt.io/icon-pf`

The suggested patternfly icon to use when displaying the resource.

### `instancetype.kubevirt.io/deprecated`

If the resource has been deprecated ahead of removal in a future release of the common-instancetypes project.

## Instance Types

### `instancetype.kubevirt.io/version`

The version of instance type class the resources has been generated from.

### `instancetype.kubevirt.io/class`

The class of the instance type.

### `instancetype.kubevirt.io/cpu`

The number of vCPUs provided by the instance type.

### `instancetype.kubevirt.io/memory`

The amount of memory provided by the instance type.

### `instancetype.kubevirt.io/numa` (optional)

If [NUMA](https://kubevirt.io/user-guide/compute/numa/) [guestmappingpassthrough](https://kubevirt.io/user-guide/compute/numa/#guestmappingpassthrough) is enabled by the instance type.

### `instancetype.kubevirt.io/dedicatedCPUPlacement` (optional)

If [dedicatedCPUPlacement](https://kubevirt.io/user-guide/compute/dedicated_cpu_resources/#requesting-dedicated-cpu-resources) is enabled by the instance type.

### `instancetype.kubevirt.io/isolateEmulatorThread` (optional)

If [isolateEmulatorThread](https://kubevirt.io/user-guide/compute/dedicated_cpu_resources/#requesting-dedicated-cpu-for-qemu-emulator) is enabled by the instance type.

### `instancetype.kubevirt.io/hugepages` (optional)

If [hugepages](https://kubevirt.io/user-guide/compute/hugepages/) are requested by the instance type.

### `instancetype.kubevirt.io/gpus` (optional)

If [GPUs](https://kubevirt.io/user-guide/compute/mediated_devices_configuration/) are requested by the instance type.

## Preferences

### `instancetype.kubevirt.io/os-type`

The underlying type of the workload supported by the preference, current values are `linux` or `windows`.

### `instancetype.kubevirt.io/arch`

The underlying architecture of the workload supported by the preference, current values are `arm64` or `amd64`.

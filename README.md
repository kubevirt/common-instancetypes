# kubevirt/common-instancetypes

A set of [instance types and preferences](https://kubevirt.io/user-guide/virtual_machines/instancetypes/) to help create [`KubeVirt`](https://kubevirt.io)  [`VirtualMachines`](http://kubevirt.io/api-reference/main/definitions.html#_v1alpha1_virtualmachine).

# Overview

## Structure

The available instance types can be divided into two categories:

1. Workload agnostic - or general purpose
2. Workload specific

Workload agnostic instance types are a good starting point to run your workload.
Once you know more about the requirements of your workload, you can start choosing a
workload specific instance type.

The following diagram summarizes the available instance types and their use-cases:

```mermaid
graph TD

classDef grp fill:white,stroke:lightgray,color:gray
classDef series fill:lightyellow,stroke:lightgray
classDef instancetype fill:

wrklda(Workload agnostic)
wrklds(Workload specific)
class wrklda grp

wrklda:::grp --> Universal:::series
Universal([Universal]):::series --> u1:::instancetype

wrklda:::grp --> Overcommitted:::series
Overcommitted([Overcommitted]):::series --> o1:::instancetype

wrklds:::grp --> Computeexclusive:::series
Computeexclusive([Compute Exclusive]):::series --> cx1:::instancetype

wrklds:::grp --> Memory:::series
Memory([Memory]):::series --> m1:::instancetype

wrklds:::grp --> Network:::series
Network([Network]):::series --> n1:::instancetype

wrklds:::grp --> Realtime:::series
Realtime([Realtime]):::series --> rt1:::instancetype
```

### Labels

Various labels are used to decorate objects provided by this project to help end users find the right resources for them. See [docs/labels.md](./docs/labels.md) for more details.

### Schema

<details><summary>Click in order to view the instance type names schema</summary>

```ebnf
instanceTypeName = seriesName , "." , size;

seriesName = ( class | vendorClass ) , version;

class = "u" | "o" | "cx" | "g" | "m" | "n" | "rt";
vendorClass = "g" , vendorHint;
vendorHint = "n" | "i" | "a";
version = "1";

size = "small" | "medium" | "large" | [( "2" | "4" | "8" )] , "xlarge";
```
</details>

# Series

.                           |  U  |  O  |  CX  |  M  |  N  |  RT
----------------------------|-----|-----|------|-----|-----|------
*Hugepages*                 |     |     |  ✓   |  ✓  |  ✓  |  ✓
*Overcommitted Memory*      |     |  ✓  |      |     |     |
*Dedicated CPU*             |     |     |  ✓   |     |  ✓  |  ✓
*Burstable CPU performance* |  ✓  |  ✓  |      |  ✓  |     |
*Isolated emulator threads* |     |     |  ✓   |     |  ✓  |  ✓
*vNUMA*                     |     |     |  ✓   |     |     |  ✓
*vCPU-To-Memory Ratio*      | 1:4 | 1:4 |  1:2 | 1:8 | 1:2 | 1:4

## U Series

The U Series is quite neutral and provides resources for
general purpose applications.

*U* is the abbreviation for "Universal", hinting at the universal
attitude towards workloads.

VMs of instance types will share physical CPU cores on a
time-slice basis with other VMs.

### U Series Characteristics

Specific characteristics of this series are:
- *Burstable CPU performance* - The workload has a baseline compute
  performance but is permitted to burst beyond this baseline, if
  excess compute resources are available.
- *vCPU-To-Memory Ratio (1:4)* - A vCPU-to-Memory ratio of 1:4, for less
  noise per node.

## O Series

The O Series is based on the U Series, with the only difference
being that memory is overcommitted.

*O* is the abbreviation for "Overcommitted".

### UO Series Characteristics

Specific characteristics of this series are:
- *Burstable CPU performance* - The workload has a baseline compute
  performance but is permitted to burst beyond this baseline, if
  excess compute resources are available.
- *Overcommitted Memory* - Memory is over-committed in order to achieve
  a higher workload density.
- *vCPU-To-Memory Ratio (1:4)* - A vCPU-to-Memory ratio of 1:4, for less
  noise per node.

## CX Series

The CX Series provides exclusive compute resources for compute
intensive applications.

*CX* is the abbreviation of "Compute Exclusive".

The exclusive resources are given to the compute threads of the
VM. In order to ensure this, some additional cores (depending
on the number of disks and NICs) will be requested to offload
the IO threading from cores dedicated to the workload.
In addition, in this series, the NUMA topology of the used
cores is provided to the VM.

### CX Series Characteristics

Specific characteristics of this series are:
- *Hugepages* - Hugepages are used in order to improve memory
  performance.
- *Dedicated CPU* - Physical cores are exclusively assigned to every
  vCPU in order to provide fixed and high compute guarantees to the
  workload.
- *Isolated emulator threads* - Hypervisor emulator threads are isolated
  from the vCPUs in order to reduce emaulation related impact on the
  workload.
- *vNUMA* - Physical NUMA topology is reflected in the guest in order to
  optimize guest sided cache utilization.
- *vCPU-To-Memory Ratio (1:2)* - A vCPU-to-Memory ratio of 1:2.

## M Series

The M Series provides resources for memory intensive
applications.

*M* is the abbreviation of "Memory".

### M Series Characteristics

Specific characteristics of this series are:
- *Hugepages* - Hugepages are used in order to improve memory
  performance.
- *Burstable CPU performance* - The workload has a baseline compute
  performance but is permitted to burst beyond this baseline, if
  excess compute resources are available.
- *vCPU-To-Memory Ratio (1:8)* - A vCPU-to-Memory ratio of 1:8, for much
  less noise per node.

## N Series

The N Series provides resources for network intensive DPDK
applications, like VNFs.

*N* is the abbreviation for "Network".

This series of instancetypes requires nodes capable
of running DPDK workloads and being marked with the respective
node-role.kubevirt.io/worker-dpdk label as such.

### N Series Characteristics

Specific characteristics of this series are:
- *Hugepages* - Hugepages are used in order to improve memory
  performance.
- *Dedicated CPU* - Physical cores are exclusively assigned to every
  vCPU in order to provide fixed and high compute guarantees to the
  workload.
- *Isolated emulator threads* - Hypervisor emulator threads are isolated
  from the vCPUs in order to reduce emaulation related impact on the
  workload.
- *vCPU-To-Memory Ratio (1:2)* - A vCPU-to-Memory ratio of 1:2.

## RT Series

The RT Series provides resources for realtime applications, like Oslat.

*RT* is the abbreviation for "realtime".

This series of instance types requires nodes capable of running
realtime applications.

### RT Series Characteristics

Specific characteristics of this series are:
- *Hugepages* - Hugepages are used in order to improve memory
  performance.
- *Dedicated CPU* - Physical cores are exclusively assigned to every
  vCPU in order to provide fixed and high compute guarantees to the
  workload.
- *Isolated emulator threads* - Hypervisor emulator threads are isolated
  from the vCPUs in order to reduce emaulation related impact on the
  workload.
- *vCPU-To-Memory Ratio (1:4)* - A vCPU-to-Memory ratio of 1:4 starting from
  the medium size.

## Development

To get started with customizing or creating your own instancetypes and preferences
see [DEVELOPMENT.md](./DEVELOPMENT.md).

## Resources

The following instancetype resources (cluster-wide and namespaced) are
provided by this project:

Name | vCPUs | Memory
-----|-------|-------
cx1.2xlarge  |  8  |  16Gi
cx1.2xlarge1gi  |  8  |  16Gi
cx1.4xlarge  |  16  |  32Gi
cx1.4xlarge1gi  |  16  |  32Gi
cx1.8xlarge  |  32  |  64Gi
cx1.8xlarge1gi  |  32  |  64Gi
cx1.large  |  2  |  4Gi
cx1.large1gi  |  2  |  4Gi
cx1.medium  |  1  |  2Gi
cx1.medium1gi  |  1  |  2Gi
cx1.xlarge  |  4  |  8Gi
cx1.xlarge1gi  |  4  |  8Gi
m1.2xlarge  |  8  |  64Gi
m1.2xlarge1gi  |  8  |  64Gi
m1.4xlarge  |  16  |  128Gi
m1.4xlarge1gi  |  16  |  128Gi
m1.8xlarge  |  32  |  256Gi
m1.8xlarge1gi  |  32  |  256Gi
m1.large  |  2  |  16Gi
m1.large1gi  |  2  |  16Gi
m1.xlarge  |  4  |  32Gi
m1.xlarge1gi  |  4  |  32Gi
n1.2xlarge  |  16  |  32Gi
n1.4xlarge  |  32  |  64Gi
n1.8xlarge  |  64  |  128Gi
n1.large  |  4  |  8Gi
n1.medium  |  4  |  4Gi
n1.xlarge  |  8  |  16Gi
o1.2xlarge  |  8  |  32Gi
o1.4xlarge  |  16  |  64Gi
o1.8xlarge  |  32  |  128Gi
o1.large  |  2  |  8Gi
o1.medium  |  1  |  4Gi
o1.micro  |  1  |  1Gi
o1.nano  |  1  |  512Mi
o1.small  |  1  |  2Gi
o1.xlarge  |  4  |  16Gi
rt1.2xlarge  |  8  |  32Gi
rt1.4xlarge  |  16  |  64Gi
rt1.8xlarge  |  32  |  128Gi
rt1.large  |  2  |  8Gi
rt1.medium  |  1  |  4Gi
rt1.micro  |  1  |  1Gi
rt1.small  |  1  |  2Gi
rt1.xlarge  |  4  |  16Gi
u1.2xlarge  |  8  |  32Gi
u1.2xmedium  |  2  |  4Gi
u1.4xlarge  |  16  |  64Gi
u1.8xlarge  |  32  |  128Gi
u1.large  |  2  |  8Gi
u1.medium  |  1  |  4Gi
u1.micro  |  1  |  1Gi
u1.nano  |  1  |  512Mi
u1.small  |  1  |  2Gi
u1.xlarge  |  4  |  16Gi

The following preference resources (cluster-wide and namespaced) are
provided by this project:

Name | Guest OS
-----|---------
alpine | Alpine
centos.stream10 | CentOS Stream 10
centos.stream10.desktop | CentOS Stream 10
centos.stream9 | CentOS Stream 9
centos.stream9.desktop | CentOS Stream 9
centos.stream9.dpdk | CentOS Stream 9
cirros | Cirros
debian | Debian
fedora | Fedora (amd64)
fedora.arm64 | Fedora (arm64)
fedora.s390x | Fedora (s390x)
legacy | Legacy Guest
linux | Linux Guest
linux.efi | Linux EFI Guest
linux.virtiotransitional | Linux Virtio Transitional Guest
opensuse.leap | OpenSUSE Leap
opensuse.tumbleweed | OpenSUSE Tumbleweed
oraclelinux | Oracle Linux
rhel.10 | Red Hat Enterprise Linux 10 (amd64)
rhel.10.arm64 | Red Hat Enterprise Linux 10 (arm64)
rhel.10.s390x | Red Hat Enterprise Linux 10 (s390x)
rhel.7 | Red Hat Enterprise Linux 7
rhel.7.desktop | Red Hat Enterprise Linux 7
rhel.8 | Red Hat Enterprise Linux 8
rhel.8.desktop | Red Hat Enterprise Linux 8
rhel.8.dpdk | Red Hat Enterprise Linux 8
rhel.9 | Red Hat Enterprise Linux 9 (amd64)
rhel.9.arm64 | Red Hat Enterprise Linux 9 (arm64)
rhel.9.desktop | Red Hat Enterprise Linux 9 Desktop (amd64)
rhel.9.dpdk | Red Hat Enterprise Linux 9 DPDK (amd64)
rhel.9.realtime | Red Hat Enterprise Linux 9 Realtime (amd64)
rhel.9.s390x | Red Hat Enterprise Linux 9 (s390x)
sles | SUSE Linux Enterprise Server
ubuntu | Ubuntu
windows.10 | Microsoft Windows 10
windows.10.virtio | Microsoft Windows 10 (virtio)
windows.11 | Microsoft Windows 11
windows.11.virtio | Microsoft Windows 11 (virtio)
windows.2k12 | Microsoft Windows Server 2012/2012 R2
windows.2k12.virtio | Microsoft Windows Server 2012/2012 R2 (virtio)
windows.2k16 | Microsoft Windows Server 2016
windows.2k16.virtio | Microsoft Windows Server 2016 (virtio)
windows.2k19 | Microsoft Windows Server 2019
windows.2k19.virtio | Microsoft Windows Server 2019 (virtio)
windows.2k22 | Microsoft Windows Server 2022
windows.2k22.virtio | Microsoft Windows Server 2022 (virtio)
windows.2k25 | Microsoft Windows Server 2025
windows.2k25.virtio | Microsoft Windows Server 2025 (virtio)
windows.2k3 | Microsoft Windows Server 2003
windows.2k8 | Microsoft Windows Server 2008/2008 R2
windows.2k8.virtio | Microsoft Windows Server 2008/2008 R2 (virtio)
windows.7 | Microsoft Windows 7
windows.7.virtio | Microsoft Windows 7 (virtio)
windows.xp | Microsoft Windows XP

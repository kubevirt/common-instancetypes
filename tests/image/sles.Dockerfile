# This builds a containerdisk containing a SLES cloud image.
FROM scratch

ARG disk=./SLES15-SP6-Minimal-VM.x86_64-Cloud-GM.qcow2

COPY --chown=107:107 ${disk} /disk/disk.img

ENV INSTANCETYPE_KUBEVIRT_IO_DEFAULT_INSTANCETYPE u1.medium
ENV INSTANCETYPE_KUBEVIRT_IO_DEFAULT_PREFERENCE sles

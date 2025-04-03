# This builds a containerdisk containing a OL cloud image.
FROM scratch

ARG disk


COPY --chown=107:107 ${disk} /disk/disk.img

ENV INSTANCETYPE_KUBEVIRT_IO_DEFAULT_INSTANCETYPE u1.2xmedium
ENV INSTANCETYPE_KUBEVIRT_IO_DEFAULT_PREFERENCE oraclelinux

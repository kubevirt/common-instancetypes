#!/bin/bash

set -e

if [ -z "$1" ]; then
  echo "Please pass Windows image version to build as first argument."
  exit 1
fi

if [ -z "$2" ]; then
  echo "Please pass path to Windows ISO as second argument."
  exit 1
fi

WINDOWS_VERSION=$1
WINDOWS_ISO=$2

# Default vars
WINDOWS_ISO_EXTRACT_PATH=/tmp/${WINDOWS_VERSION}
WINDOWS_ISO_MODIFIED=/tmp/${WINDOWS_VERSION}.iso
VIRTIO_WIN_ISO_LATEST_URL=https://fedorapeople.org/groups/virt/virtio-win/direct-downloads/latest-virtio/virtio-win.iso
VIRTIO_WIN_ISO=/tmp/virtio-win.iso
# EFI boot by default
WINDOWS_BOOT_IMAGE=(-e efi/microsoft/boot/efisys_noprompt.bin)
DEFAULT_INSTANCETYPE=u1.small

case "$WINDOWS_VERSION" in
  windows10)
    OS_VARIANT=win10
    WINDOWS_ISO_OVERLAYS=(overlays/generic overlays/windows10)
    VIRT_INSTALL_EXTRA_ARGS=(--boot uefi)
    DEFAULT_PREFERENCE=windows.10.virtio
    ;;
  windows11)
    OS_VARIANT=win11
    WINDOWS_ISO_OVERLAYS=(overlays/generic overlays/windows11)
    DEFAULT_INSTANCETYPE=u1.large
    DEFAULT_PREFERENCE=windows.11.virtio
    ;;
  windows2k16)
    OS_VARIANT=win2k16
    WINDOWS_ISO_OVERLAYS=(overlays/windows2k16)
    # Create BIOS bootable ISO
    WINDOWS_BOOT_IMAGE=(-b boot/etfsboot.com)
    DEFAULT_PREFERENCE=windows.2k16.virtio
    ;;
  windows2k19)
    OS_VARIANT=win2k19
    WINDOWS_ISO_OVERLAYS=(overlays/generic overlays/windows2k19)
    # Create BIOS bootable ISO
    WINDOWS_BOOT_IMAGE=(-b boot/etfsboot.com)
    DEFAULT_PREFERENCE=windows.2k19.virtio
    ;;
  windows2k22)
    OS_VARIANT=win2k22
    WINDOWS_ISO_OVERLAYS=(overlays/generic overlays/windows2k22)
    VIRT_INSTALL_EXTRA_ARGS=(--boot uefi)
    DEFAULT_PREFERENCE=windows.2k22.virtio
    ;;
  *)
    echo "Need a valid Windows image version to build: windows10, windows11, windows2k16, windows2k19, windows2k22"
    exit 1
    ;;
esac

# Extract Windows ISO, create a bootable ISO with overlays applied that does not prompt for a key and clean up
# By default the ISO is EFI bootable
7z x "${WINDOWS_ISO}" -o"${WINDOWS_ISO_EXTRACT_PATH}"
mkisofs -udf --allow-limited-size "${WINDOWS_BOOT_IMAGE[@]}" -no-emul-boot \
  -o "${WINDOWS_ISO_MODIFIED}" "${WINDOWS_ISO_EXTRACT_PATH}" "${WINDOWS_ISO_OVERLAYS[@]}"
rm -r "${WINDOWS_ISO_EXTRACT_PATH}"

# Download latest virtio-win.iso if it was not passed as argument, else copy passed ISO to /tmp
if [ -z "$3" ]; then
  curl -L "${VIRTIO_WIN_ISO_LATEST_URL}" -o "${VIRTIO_WIN_ISO}"
else
  cp "$3" "${VIRTIO_WIN_ISO}"
fi

# Set required SELinux contexts
chcon -t virt_content_t "${WINDOWS_ISO_MODIFIED}"
chcon -t virt_content_t "${VIRTIO_WIN_ISO}"

# Start installation process and wait for first reboot of VM, then restart it
virt-install --noautoconsole --wait --name="${WINDOWS_VERSION}" \
  --vcpus=2 --memory=8192 --disk size=15,bus=virtio --video virtio --network user,model=virtio \
  --os-variant="${OS_VARIANT}" --cdrom="${WINDOWS_ISO_MODIFIED}" --disk device=cdrom,path="${VIRTIO_WIN_ISO}" \
  "${VIRT_INSTALL_EXTRA_ARGS[@]}"

# Second loop to wait for VM to shut down after finished install
while LANG=C virsh domstate "${WINDOWS_VERSION}" | grep running > /dev/null; do
  sleep 1
done

# Compress VM disk
qemu-img convert -p -O qcow2 -c "$HOME/.local/share/libvirt/images/${WINDOWS_VERSION}.qcow2" "${WINDOWS_VERSION}.qcow2"

# Remove VM and all images
virsh undefine --nvram --remove-all-storage "${WINDOWS_VERSION}"

# Build containerdisk
podman build -t "${WINDOWS_VERSION}-container-disk:latest" -f - . << EOF
FROM scratch
ADD --chown=107:107 ${WINDOWS_VERSION}.qcow2 /disk/

LABEL instancetype.kubevirt.io/default-instancetype ${DEFAULT_INSTANCETYPE}
LABEL instancetype.kubevirt.io/default-preference ${DEFAULT_PREFERENCE}
LABEL instancetype.kubevirt.io/display-needed true

# Set ENVs for compatibility with crun-vm
ENV INSTANCETYPE_KUBEVIRT_IO_DEFAULT_INSTANCETYPE ${DEFAULT_INSTANCETYPE}
ENV INSTANCETYPE_KUBEVIRT_IO_DEFAULT_PREFERENCE ${DEFAULT_PREFERENCE}
ENV INSTANCETYPE_KUBEVIRT_IO_DISPLAY_NEEDED true
EOF

# Remove VM disk
rm "${WINDOWS_VERSION}.qcow2"

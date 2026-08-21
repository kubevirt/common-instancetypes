# This builds a containerdisk containing a bootable ISO of the Microsoft Validation OS.
# SSH is enabled via a startup script to verify the guest OS is able to start.
#
# bootloader_iso_url only provides the UEFI bootloader, as the Validation OS ISO
# ships no bootable media. Any Windows installation media works, the Windows
# Server evaluation ISO is used as it has a static URL.
#
# The ISOs are cached between builds, keyed by the hash of their URL. Run
# `podman builder prune --build-cache` to refetch, e.g. to pick up a newer
# Validation OS release from the unversioned URL.
#
# For CI purposes a recent build can be found at this private location:
# quay.io/kubevirtci/validation-os-container-disk

FROM quay.io/fedora/fedora-minimal:latest as builder

RUN microdnf install -y p7zip p7zip-plugins wimlib-utils chntpw mkisofs

# Create the working dir, extract ValidationsOS.wim and place it
# as sources/boot.wim into the working dir
ARG validationos_iso_url=https://aka.ms/DownloadValidationOS
RUN --mount=type=cache,id=validationos-isos,sharing=locked,target=/isos \
  iso="/isos/$(echo "$validationos_iso_url" | sha256sum | cut -c1-16).iso" && \
  if [ ! -s "$iso" ]; then curl -fsSL "$validationos_iso_url" -o "$iso.part" && mv "$iso.part" "$iso"; fi && \
  mkdir -p disk/sources && \
  7z x -y '-i!ValidationOS.wim' "$iso" && \
  mv ValidationOS.wim disk/sources/boot.wim

# Extract the bootloader from the Windows installation media into the working dir
ARG bootloader_iso_url=https://software-static.download.prss.microsoft.com/dbazure/888969d5-f34g-4e03-ac9d-1f9786c66749/26100.1742.240906-0331.ge_release_svc_refresh_SERVER_EVAL_x64FRE_en-us.iso
RUN --mount=type=cache,id=validationos-isos,sharing=locked,target=/isos \
  iso="/isos/$(echo "$bootloader_iso_url" | sha256sum | cut -c1-16).iso" && \
  if [ ! -s "$iso" ]; then curl -fsSL "$bootloader_iso_url" -o "$iso.part" && mv "$iso.part" "$iso"; fi && \
  7z x -y '-i!boot' '-i!bootmgr.efi' '-i!efi' "$iso" -o./disk/

# Startup script setting Administrator's password and starting sshd service
# Taken from https://learn.microsoft.com/en-us/windows-hardware/manufacture/desktop/validation-os-run-an-app-on-boot?view=windows-11
COPY <<EOF startnet.valos.cmd
@ECHO OFF
net user Administrator Administrator
net start sshd
netsh advfirewall firewall add rule name="sshd" dir=in action=allow protocol=TCP localport=22 remoteip=any
EOF

# Registry patch enabling the startup script
COPY <<EOF winlogon.reg
Windows Registry Editor Version 5.00

[HKEY_LOCAL_MACHINE\SOFTWARE\Microsoft\Windows NT\CurrentVersion\Winlogon]
"Shell"="cmd /k c:\\windows\\System32\\startnet.valos.cmd"
EOF

# Add startup script to boot.wim and enable running it with registry patch
RUN echo "add startnet.valos.cmd /Windows/System32/startnet.valos.cmd" | wimupdate disk/sources/boot.wim 2 && \
  wimextract disk/sources/boot.wim 2 Windows/System32/config/SOFTWARE && \
  { reged -C -I SOFTWARE HKEY_LOCAL_MACHINE\\SOFTWARE winlogon.reg || true; } && \
  echo "add SOFTWARE Windows/System32/config/SOFTWARE" | wimupdate disk/sources/boot.wim 2

# Create an EFI bootable ISO using the contents in the working dir
RUN mkisofs -V "Microsoft Validation OS" -udf -e efi/microsoft/boot/efisys_noprompt.bin -no-emul-boot -o disk.img disk/

FROM scratch

COPY --from=builder --chown=107:107 /disk.img /disk/

ENV INSTANCETYPE_KUBEVIRT_IO_DEFAULT_INSTANCETYPE u1.2xmedium
ENV INSTANCETYPE_KUBEVIRT_IO_DEFAULT_PREFERENCE windows.11
ENV INSTANCETYPE_KUBEVIRT_IO_DISPLAY_NEEDED true

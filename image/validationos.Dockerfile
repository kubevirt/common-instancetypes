# This builds a containerdisk containing a bootable ISO of the Microsoft Validation OS.
# SSH is enabled via a startup script to verify the guest OS is able to start.
# Call it like this:
# podman build -f validationos.Dockerfile \
#   --build-arg validationos_iso_url="https://aka.ms/DownloadValidationOS"
#   --build-arg windows11_iso_url="$MY_WIN11_DOWNLOAD_URL" \
#
# For CI purposes a recent build can be found at this private location:
# quay.io/kubevirtci/validation-os-container-disk

FROM quay.io/fedora/fedora-minimal:latest as builder

RUN microdnf install -y p7zip p7zip-plugins wimlib-utils chntpw mkisofs

ARG validationos_iso_url
RUN curl -L "$validationos_iso_url" -o validationos.iso

ARG windows11_iso_url
RUN curl -L "$windows11_iso_url" -o windows11.iso

# Create the working dir, extract ValidationsOS.wim and place it
# as sources/boot.wim into the working dir
RUN mkdir -p disk/sources && \
  7z x '-i!ValidationOS.wim' validationos.iso && \
  mv ValidationOS.wim disk/sources/boot.wim

# Startup script setting Administrator's password and starting sshd service
# Taken from https://learn.microsoft.com/en-us/windows-hardware/manufacture/desktop/validation-os-run-an-app-on-boot?view=windows-11
COPY <<EOF startnet.valos.cmd
@ECHO OFF
net user Administrator Administrator
net start sshd
)
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
  reged -C -I SOFTWARE HKEY_LOCAL_MACHINE\\SOFTWARE winlogon.reg || true && \
  echo "add SOFTWARE Windows/System32/config/SOFTWARE" | wimupdate disk/sources/boot.wim 2

# Extract the bootloader from the Windows 11 ISO into the working dir
RUN 7z x '-i!boot' '-i!bootmgr.efi' '-i!efi' windows11.iso -o./disk/

# Create an EFI bootable ISO using the contents in the working dir
RUN mkisofs -V "Microsoft Validation OS" -udf -e efi/microsoft/boot/efisys_noprompt.bin -no-emul-boot -o disk.img disk/

FROM scratch

COPY --from=builder --chown=107:107 /disk.img /disk/

# Install virtio guest drivers
Start-Process msiexec -Wait -ArgumentList "/i E:\virtio-win-gt-x64.msi /qn /passive /norestart"

# Install qemu guest agent
Start-Process msiexec -Wait -ArgumentList "/i E:\guest-agent\qemu-ga-x86_64.msi /qn /passive /norestart"

#
# Taken from https://learn.microsoft.com/de-de/windows-server/administration/openssh/openssh_install_firstuse?tabs=powershell
#

# Install the OpenSSH Server
Add-WindowsCapability -Online -Name OpenSSH.Server~~~~0.0.1.0

# Start the sshd service
Start-Service sshd

# OPTIONAL but recommended:
Set-Service -Name sshd -StartupType 'Automatic'

# Confirm the Firewall rule is configured. It should be created automatically by setup. Run the following to verify
if (!(Get-NetFirewallRule -Name "OpenSSH-Server-In-TCP" -ErrorAction SilentlyContinue | Select-Object Name, Enabled)) {
    Write-Output "Firewall Rule 'OpenSSH-Server-In-TCP' does not exist, creating it..."
    New-NetFirewallRule -Name 'OpenSSH-Server-In-TCP' -DisplayName 'OpenSSH Server (sshd)' -Enabled True -Direction Inbound -Protocol TCP -Action Allow -LocalPort 22
} else {
    Write-Output "Firewall rule 'OpenSSH-Server-In-TCP' has been created and exists."
}

# Download and run SDelete
$url = "https://download.sysinternals.com/files/SDelete.zip"
$file = "C:\SDelete.zip"
$path = "C:\SDelete"
Invoke-WebRequest -Uri $url -OutFile $file
Expand-Archive -LiteralPath $file -DestinationPath $path
Start-Process $path\sdelete64.exe -Wait -ArgumentList "-accepteula -z C:"

# Shut down the machine
Stop-Computer -Force

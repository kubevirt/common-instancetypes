# Install virtio guest drivers
Start-Process msiexec -Wait -ArgumentList "/i E:\virtio-win-gt-x64.msi /qn /passive /norestart"

# Install qemu guest agent
Start-Process msiexec -Wait -ArgumentList "/i E:\guest-agent\qemu-ga-x86_64.msi /qn /passive /norestart"

# Enable Tls12, which is not enabled by default in Windows Server 2016
[Net.ServicePointManager]::SecurityProtocol = [Net.SecurityProtocolType]::Tls12

# Download and run SDelete
$url = "https://download.sysinternals.com/files/SDelete.zip"
$file = "C:\SDelete.zip"
$path = "C:\SDelete"
Invoke-WebRequest -Uri $url -OutFile $file
Expand-Archive -LiteralPath $file -DestinationPath $path
Start-Process $path\sdelete64.exe -Wait -ArgumentList "-accepteula -z C:"

# Shut down the machine
Stop-Computer -Force

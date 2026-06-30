# ssh-wizard installer for Windows (PowerShell)
# Usage: irm https://raw.githubusercontent.com/duong6003/ssh-wizard/master/install.ps1 | iex

$ErrorActionPreference = "Stop"
$Repo = "duong6003/ssh-wizard"
$Binary = "ssh-wizard"

# Get latest version
$Release = Invoke-RestMethod "https://api.github.com/repos/$Repo/releases/latest"
$Version = $Release.tag_name

$Filename = "$Binary-windows-amd64.zip"
$Url = "https://github.com/$Repo/releases/download/$Version/$Filename"

Write-Host "Installing $Binary $Version (windows/amd64)..."

$Tmp = Join-Path $env:TEMP "ssh-wizard-install"
New-Item -ItemType Directory -Force -Path $Tmp | Out-Null

$ZipPath = Join-Path $Tmp $Filename
Invoke-WebRequest -Uri $Url -OutFile $ZipPath
Expand-Archive -Path $ZipPath -DestinationPath $Tmp -Force
Remove-Item $ZipPath

# Install to %LOCALAPPDATA%\Programs\ssh-wizard
$InstallDir = Join-Path $env:LOCALAPPDATA "Programs\ssh-wizard"
New-Item -ItemType Directory -Force -Path $InstallDir | Out-Null
Move-Item -Force (Join-Path $Tmp "$Binary.exe") (Join-Path $InstallDir "$Binary.exe")
Remove-Item -Recurse -Force $Tmp

# Add to user PATH if not already there
$CurrentPath = [Environment]::GetEnvironmentVariable("PATH", "User")
if ($CurrentPath -notlike "*$InstallDir*") {
    [Environment]::SetEnvironmentVariable("PATH", "$CurrentPath;$InstallDir", "User")
    Write-Host "Added $InstallDir to PATH"
    Write-Host "Restart your terminal, then run: $Binary"
} else {
    Write-Host "Installed to $InstallDir"
    Write-Host "Run: $Binary"
}

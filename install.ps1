#!/usr/bin/env pwsh
$ErrorActionPreference = "Stop"

$repo = "hewliyang/spliit-cli"
$binary = "spliit.exe"

# Detect architecture
$arch = if ([Environment]::Is64BitOperatingSystem) {
    if ($env:PROCESSOR_ARCHITECTURE -eq "ARM64") { "arm64" } else { "amd64" }
} else {
    Write-Error "32-bit systems are not supported"
    exit 1
}

# Get latest version
$release = Invoke-RestMethod -Uri "https://api.github.com/repos/$repo/releases/latest"
$version = $release.tag_name
$versionNum = $version.TrimStart('v')

Write-Host "Installing spliit $version (windows/$arch)..."

# Download URL
$url = "https://github.com/$repo/releases/download/$version/spliit_${versionNum}_windows_${arch}.zip"

# Create temp directory
$tmpDir = New-Item -ItemType Directory -Path ([System.IO.Path]::GetTempPath()) -Name ([System.Guid]::NewGuid().ToString())

try {
    # Download
    $zipPath = Join-Path $tmpDir "spliit.zip"
    Write-Host "Downloading from $url"
    Invoke-WebRequest -Uri $url -OutFile $zipPath

    # Extract
    Expand-Archive -Path $zipPath -DestinationPath $tmpDir

    # Install to user's local bin
    $installDir = Join-Path $env:LOCALAPPDATA "Programs\spliit"
    if (!(Test-Path $installDir)) {
        New-Item -ItemType Directory -Path $installDir | Out-Null
    }

    Copy-Item (Join-Path $tmpDir $binary) -Destination $installDir -Force

    # Add to PATH if not already there
    $userPath = [Environment]::GetEnvironmentVariable("Path", "User")
    if ($userPath -notlike "*$installDir*") {
        [Environment]::SetEnvironmentVariable("Path", "$userPath;$installDir", "User")
        Write-Host "Added $installDir to PATH (restart your terminal to use)"
    }

    Write-Host "✓ Installed spliit to $installDir\$binary"
    Write-Host "  Run 'spliit --help' to get started"
}
finally {
    Remove-Item -Recurse -Force $tmpDir
}

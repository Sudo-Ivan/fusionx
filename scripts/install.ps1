# FusionX Installation Script for Windows
# Usage: Invoke-WebRequest -Uri "https://raw.githubusercontent.com/Sudo-Ivan/fusionx/main/scripts/install.ps1" -OutFile "install.ps1"; .\install.ps1

param(
    [string]$InstallDir = "$env:LOCALAPPDATA\Programs\FusionX",
    [switch]$Help,
    [switch]$Version
)

# Configuration
$RepoOwner = "Sudo-Ivan"
$RepoName = "fusionx"
$BinaryName = "fusionx.exe"
$TmpDir = "$env:TEMP\fusionx-install"

# Color codes for console output
$Colors = @{
    Red = [ConsoleColor]::Red
    Green = [ConsoleColor]::Green
    Yellow = [ConsoleColor]::Yellow
    Blue = [ConsoleColor]::Blue
    White = [ConsoleColor]::White
}

# Logging functions
function Write-LogInfo {
    param([string]$Message)
    Write-Host "[INFO] $Message" -ForegroundColor $Colors.Blue
}

function Write-LogSuccess {
    param([string]$Message)
    Write-Host "[SUCCESS] $Message" -ForegroundColor $Colors.Green
}

function Write-LogWarning {
    param([string]$Message)
    Write-Host "[WARNING] $Message" -ForegroundColor $Colors.Yellow
}

function Write-LogError {
    param([string]$Message)
    Write-Host "[ERROR] $Message" -ForegroundColor $Colors.Red
}

# Show help information
function Show-Help {
    Write-Host "FusionX Installation Script for Windows" -ForegroundColor $Colors.Green
    Write-Host ""
    Write-Host "Usage: .\install.ps1 [options]"
    Write-Host ""
    Write-Host "Options:"
    Write-Host "  -InstallDir <path>  Installation directory (default: $env:LOCALAPPDATA\Programs\FusionX)"
    Write-Host "  -Help               Show this help message"
    Write-Host "  -Version            Show version information"
    Write-Host ""
    Write-Host "Examples:"
    Write-Host "  .\install.ps1"
    Write-Host "  .\install.ps1 -InstallDir 'C:\Tools\FusionX'"
    Write-Host ""
    Write-Host "One-liner installation:"
    Write-Host "  Invoke-WebRequest -Uri 'https://raw.githubusercontent.com/$RepoOwner/$RepoName/main/scripts/install.ps1' -OutFile 'install.ps1'; .\install.ps1"
}

# Show version information
function Show-Version {
    Write-Host "FusionX Installation Script v1.0.0" -ForegroundColor $Colors.Green
}

# Detect architecture
function Get-Architecture {
    $arch = $env:PROCESSOR_ARCHITECTURE
    switch ($arch) {
        "AMD64" { return "amd64" }
        "ARM64" { return "arm64" }
        default { 
            Write-LogError "Unsupported architecture: $arch"
            exit 1
        }
    }
}

# Get latest release version
function Get-LatestVersion {
    try {
        $apiUrl = "https://api.github.com/repos/$RepoOwner/$RepoName/releases/latest"
        $response = Invoke-RestMethod -Uri $apiUrl -Method Get
        return $response.tag_name
    }
    catch {
        Write-LogError "Failed to get latest version: $($_.Exception.Message)"
        exit 1
    }
}

# Download file with progress
function Download-File {
    param(
        [string]$Url,
        [string]$OutputPath
    )
    
    try {
        Write-LogInfo "Downloading from: $Url"
        Invoke-WebRequest -Uri $Url -OutFile $OutputPath -UseBasicParsing
        Write-LogSuccess "Downloaded: $OutputPath"
    }
    catch {
        Write-LogError "Failed to download file: $($_.Exception.Message)"
        throw
    }
}

# Verify checksum
function Test-Checksum {
    param(
        [string]$FilePath,
        [string]$ExpectedChecksum
    )
    
    try {
        $actualChecksum = (Get-FileHash -Path $FilePath -Algorithm SHA256).Hash.ToLower()
        $expectedChecksum = $ExpectedChecksum.ToLower()
        
        if ($actualChecksum -eq $expectedChecksum) {
            Write-LogSuccess "Checksum verification passed"
            return $true
        }
        else {
            Write-LogError "Checksum verification failed!"
            Write-LogError "Expected: $expectedChecksum"
            Write-LogError "Actual:   $actualChecksum"
            return $false
        }
    }
    catch {
        Write-LogWarning "Could not verify checksum: $($_.Exception.Message)"
        return $true  # Continue installation if checksum verification fails
    }
}

# Add directory to PATH
function Add-ToPath {
    param([string]$Directory)
    
    $currentPath = [Environment]::GetEnvironmentVariable("PATH", "User")
    if ($currentPath -notlike "*$Directory*") {
        Write-LogInfo "Adding $Directory to user PATH..."
        $newPath = "$currentPath;$Directory"
        [Environment]::SetEnvironmentVariable("PATH", $newPath, "User")
        Write-LogSuccess "Added to PATH. Restart your terminal to use 'fusion' command."
    }
    else {
        Write-LogInfo "Directory already in PATH"
    }
}

# Create desktop shortcut
function New-DesktopShortcut {
    param(
        [string]$TargetPath,
        [string]$ShortcutPath
    )
    
    try {
        $WScriptShell = New-Object -comObject WScript.Shell
        $Shortcut = $WScriptShell.CreateShortcut($ShortcutPath)
        $Shortcut.TargetPath = $TargetPath
        $Shortcut.WorkingDirectory = Split-Path $TargetPath
        $Shortcut.Description = "FusionX RSS Reader"
        $Shortcut.Save()
        Write-LogSuccess "Created desktop shortcut"
    }
    catch {
        Write-LogWarning "Could not create desktop shortcut: $($_.Exception.Message)"
    }
}

# Main installation function
function Install-FusionX {
    Write-LogInfo "Starting FusionX installation for Windows..."
    
    # Detect architecture
    $arch = Get-Architecture
    Write-LogInfo "Detected architecture: $arch"
    
    # Get latest version
    Write-LogInfo "Fetching latest release information..."
    $version = Get-LatestVersion
    Write-LogInfo "Latest version: $version"
    
    # Create temporary directory
    if (Test-Path $TmpDir) {
        Remove-Item $TmpDir -Recurse -Force
    }
    New-Item -ItemType Directory -Path $TmpDir -Force | Out-Null
    
    # Construct download URLs with version in filename (remove 'v' prefix for filename)
    $platform = "windows_$arch"
    $versionClean = $version -replace '^v', ''
    $binaryName = "fusion_$versionClean`_$platform.exe"
    $binaryUrl = "https://github.com/$RepoOwner/$RepoName/releases/download/$version/$binaryName"
    $checksumUrl = "https://github.com/$RepoOwner/$RepoName/releases/download/$version/checksums.txt"
    
    $binaryPath = Join-Path $TmpDir $binaryName
    $checksumPath = Join-Path $TmpDir "checksums.txt"
    
    # Download binary and checksums
    Write-LogInfo "Downloading binary..."
    Download-File -Url $binaryUrl -OutputPath $binaryPath
    
    Write-LogInfo "Downloading checksums..."
    Download-File -Url $checksumUrl -OutputPath $checksumPath
    
    # Verify checksum
    Write-LogInfo "Verifying checksum..."
    $checksumContent = Get-Content $checksumPath
    $expectedChecksum = ($checksumContent | Where-Object { $_ -like "*$binaryName*" }).Split(' ')[0]
    
    if ($expectedChecksum) {
        if (-not (Test-Checksum -FilePath $binaryPath -ExpectedChecksum $expectedChecksum)) {
            Write-LogError "Checksum verification failed. Installation aborted."
            exit 1
        }
    }
    else {
        Write-LogWarning "Could not find checksum for $binaryName in checksums.txt"
    }
    
    # Create installation directory
    Write-LogInfo "Creating installation directory: $InstallDir"
    if (-not (Test-Path $InstallDir)) {
        New-Item -ItemType Directory -Path $InstallDir -Force | Out-Null
    }
    
    # Install binary
    $finalBinaryPath = Join-Path $InstallDir $BinaryName
    Write-LogInfo "Installing binary to: $finalBinaryPath"
    
    # Stop any running instances
    Get-Process -Name "fusionx" -ErrorAction SilentlyContinue | Stop-Process -Force
    
    Copy-Item $binaryPath $finalBinaryPath -Force
    
    # Clean up temporary files
    Remove-Item $TmpDir -Recurse -Force
    
    Write-LogSuccess "FusionX installed successfully!"
    Write-LogInfo "Binary location: $finalBinaryPath"
    
    # Add to PATH
    Add-ToPath -Directory $InstallDir
    
    # Create desktop shortcut
    $desktopPath = [Environment]::GetFolderPath("Desktop")
    $shortcutPath = Join-Path $desktopPath "FusionX.lnk"
    New-DesktopShortcut -TargetPath $finalBinaryPath -ShortcutPath $shortcutPath
    
    # Ask if user wants to start the application
    $response = Read-Host "Do you want to start FusionX now? (y/n)"
    if ($response -match "^[Yy]") {
        Write-LogInfo "Starting FusionX..."
        Write-LogInfo "FusionX will be available at http://localhost:8080"
        Write-LogInfo "Press Ctrl+C to stop"
        
        # Start in new window
        Start-Process -FilePath $finalBinaryPath -WindowStyle Normal
        Write-LogSuccess "FusionX started in a new window"
        Write-LogInfo "Access FusionX at: http://localhost:8080"
    }
    else {
        Write-LogInfo "You can start FusionX later by running: fusionx"
        Write-LogInfo "Or use the desktop shortcut"
        Write-LogInfo "It will be available at: http://localhost:8080"
    }
}

# Handle script parameters
if ($Help) {
    Show-Help
    exit 0
}

if ($Version) {
    Show-Version
    exit 0
}

# Check PowerShell version
if ($PSVersionTable.PSVersion.Major -lt 5) {
    Write-LogError "PowerShell 5.0 or later is required"
    exit 1
}

# Check if running with administrator privileges for system-wide installation
if ($InstallDir -like "C:\Program Files*") {
    if (-NOT ([Security.Principal.WindowsPrincipal] [Security.Principal.WindowsIdentity]::GetCurrent()).IsInRole([Security.Principal.WindowsBuiltInRole] "Administrator")) {
        Write-LogWarning "Administrator privileges required for system-wide installation"
        Write-LogInfo "Either run PowerShell as Administrator or install to user directory"
        $InstallDir = "$env:LOCALAPPDATA\Programs\FusionX"
        Write-LogInfo "Using user installation directory: $InstallDir"
    }
}

# Run main installation
try {
    Install-FusionX
}
catch {
    Write-LogError "Installation failed: $($_.Exception.Message)"
    exit 1
}

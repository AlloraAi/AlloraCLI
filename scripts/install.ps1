# AlloraCLI Windows Installation Script
# This script automates the installation of AlloraCLI on Windows

# Configuration
$RepoOwner = "AlloraAi"
$RepoName = "AlloraCLI"
$BinaryName = "allora.exe"
$InstallDir = "C:\Tools"
$ConfigDir = "$env:APPDATA\alloracli"

# Colors for output
function Write-Success {
    param($Message)
    Write-Host "[SUCCESS] $Message" -ForegroundColor Green
}

function Write-Info {
    param($Message)
    Write-Host "[INFO] $Message" -ForegroundColor Cyan
}

function Write-Warning {
    param($Message)
    Write-Host "[WARNING] $Message" -ForegroundColor Yellow
}

function Write-Error-Message {
    param($Message)
    Write-Host "[ERROR] $Message" -ForegroundColor Red
}

# Check if running as Administrator
function Test-Administrator {
    $currentUser = [Security.Principal.WindowsIdentity]::GetCurrent()
    $principal = New-Object Security.Principal.WindowsPrincipal($currentUser)
    return $principal.IsInRole([Security.Principal.WindowsBuiltInRole]::Administrator)
}

# Get the latest release version
function Get-LatestVersion {
    Write-Info "Fetching latest release information..."
    try {
        $response = Invoke-RestMethod -Uri "https://api.github.com/repos/$RepoOwner/$RepoName/releases/latest"
        $version = $response.tag_name
        Write-Info "Latest version: $version"
        return $version
    }
    catch {
        Write-Error-Message "Failed to fetch latest version: $_"
        exit 1
    }
}

# Download the binary
function Download-Binary {
    param($Version)
    
    $binaryName = "allora-windows-amd64.exe"
    $downloadUrl = "https://github.com/$RepoOwner/$RepoName/releases/download/$Version/$binaryName"
    $tempFile = Join-Path $env:TEMP $binaryName
    
    Write-Info "Downloading AlloraCLI $Version..."
    Write-Info "URL: $downloadUrl"
    
    try {
        Invoke-WebRequest -Uri $downloadUrl -OutFile $tempFile -UseBasicParsing
        Write-Success "Download completed"
        return $tempFile
    }
    catch {
        Write-Error-Message "Failed to download binary: $_"
        exit 1
    }
}

# Install the binary
function Install-Binary {
    param($SourcePath)
    
    Write-Info "Installing binary to $InstallDir..."
    
    # Create install directory if it doesn't exist
    if (-not (Test-Path $InstallDir)) {
        New-Item -ItemType Directory -Path $InstallDir -Force | Out-Null
        Write-Info "Created directory: $InstallDir"
    }
    
    $destPath = Join-Path $InstallDir $BinaryName
    
    # Copy the binary
    try {
        Copy-Item -Path $SourcePath -Destination $destPath -Force
        Write-Success "Binary installed successfully"
        return $destPath
    }
    catch {
        Write-Error-Message "Failed to install binary: $_"
        exit 1
    }
}

# Add to PATH
function Add-ToPath {
    Write-Info "Configuring PATH environment variable..."
    
    # Get current PATH
    $currentPath = [Environment]::GetEnvironmentVariable("PATH", "User")
    
    # Check if already in PATH
    if ($currentPath -like "*$InstallDir*") {
        Write-Info "$InstallDir is already in PATH"
        return
    }
    
    # Add to PATH with proper separator handling
    try {
        if ([string]::IsNullOrEmpty($currentPath)) {
            $newPath = $InstallDir
        }
        elseif ($currentPath.EndsWith(";")) {
            $newPath = "$currentPath$InstallDir"
        }
        else {
            $newPath = "$currentPath;$InstallDir"
        }
        
        [Environment]::SetEnvironmentVariable("PATH", $newPath, "User")
        Write-Success "Added $InstallDir to PATH"
        Write-Warning "Please restart your terminal for PATH changes to take effect"
        
        # Update current session PATH
        $env:PATH = "$env:PATH;$InstallDir"
    }
    catch {
        Write-Error-Message "Failed to update PATH: $_"
        Write-Warning "Please manually add $InstallDir to your PATH"
    }
}

# Setup configuration directory
function Setup-Config {
    Write-Info "Setting up configuration directory..."
    
    if (-not (Test-Path $ConfigDir)) {
        New-Item -ItemType Directory -Path $ConfigDir -Force | Out-Null
        Write-Info "Created configuration directory: $ConfigDir"
    }
    else {
        Write-Info "Configuration directory already exists: $ConfigDir"
    }
    
    # Create plugins directory
    $pluginsDir = Join-Path $ConfigDir "plugins"
    if (-not (Test-Path $pluginsDir)) {
        New-Item -ItemType Directory -Path $pluginsDir -Force | Out-Null
    }
    
    # Create basic config file if it doesn't exist
    $configFile = Join-Path $ConfigDir "config.yaml"
    if (-not (Test-Path $configFile)) {
        $configContent = @"
version: "1.0.0"
agents:
  default:
    type: "general"
    model: "gpt-4"
    max_tokens: 4096
    temperature: 0.7
logging:
  level: "info"
  format: "text"
security:
  encryption: true
  audit_logging: true
"@
        Set-Content -Path $configFile -Value $configContent
        Write-Info "Created basic configuration file: $configFile"
    }
}

# Verify installation
function Test-Installation {
    Write-Info "Verifying installation..."
    
    try {
        $version = & allora --version 2>&1
        Write-Success "Installation successful!"
        Write-Info "Version: $version"
        Write-Info ""
        Write-Info "You can now run 'allora --help' to get started"
        Write-Info "To initialize AlloraCLI, run 'allora init'"
    }
    catch {
        Write-Error-Message "Installation verification failed"
        Write-Warning "Try running 'allora --version' after restarting your terminal"
    }
}

# Main installation function
function Install-AlloraCLI {
    Write-Host ""
    Write-Host "╔═══════════════════════════════════════════════════════════╗" -ForegroundColor Cyan
    Write-Host "║        AlloraCLI Windows Installation Script             ║" -ForegroundColor Cyan
    Write-Host "╚═══════════════════════════════════════════════════════════╝" -ForegroundColor Cyan
    Write-Host ""
    
    # Check for administrator privileges
    if (Test-Administrator) {
        Write-Info "Running with administrator privileges"
    }
    else {
        Write-Warning "Not running as administrator. Installation will proceed to user directory."
    }
    
    # Get latest version
    $version = Get-LatestVersion
    
    # Download binary
    $tempBinary = Download-Binary -Version $version
    
    # Install binary
    $installedPath = Install-Binary -SourcePath $tempBinary
    
    # Setup configuration
    Setup-Config
    
    # Add to PATH
    Add-ToPath
    
    # Clean up temp file
    Remove-Item -Path $tempBinary -Force -ErrorAction SilentlyContinue
    
    # Verify installation
    Test-Installation
    
    Write-Host ""
    Write-Success "Installation completed successfully!"
    Write-Info "Documentation: https://docs.alloraai.com"
    Write-Info "GitHub: https://github.com/$RepoOwner/$RepoName"
    Write-Host ""
}

# Run installation
Install-AlloraCLI

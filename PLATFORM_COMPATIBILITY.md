# Platform Compatibility Guide

This document provides detailed information about AlloraCLI compatibility across different operating systems, architectures, and deployment methods.

## Supported Platforms

### Operating Systems

#### ✅ Linux
- **Distributions**: Ubuntu, Debian, Fedora, CentOS, RHEL, Arch, Alpine, and other modern distributions
- **Architectures**:
  - AMD64 (x86_64) - Primary support
  - ARM64 (aarch64) - Full support
- **Minimum Requirements**:
  - Kernel: 3.10+ (recommended: 4.x+)
  - glibc: 2.17+ (or musl for Alpine)
  - Disk Space: 100 MB
  - RAM: 128 MB minimum, 512 MB recommended

#### ✅ macOS
- **Versions**: macOS 10.15 (Catalina) and later
- **Architectures**:
  - Intel (x86_64) - Full support
  - Apple Silicon (ARM64/M1/M2/M3) - Native support
- **Minimum Requirements**:
  - macOS 10.15+
  - Disk Space: 100 MB
  - RAM: 256 MB minimum, 512 MB recommended

#### ✅ Windows
- **Versions**: Windows 10 (1809+) and Windows 11
- **Architectures**:
  - AMD64 (x86_64) - Primary support
- **Minimum Requirements**:
  - Windows 10 version 1809 or later
  - PowerShell 5.1+ (included by default)
  - Disk Space: 100 MB
  - RAM: 256 MB minimum, 512 MB recommended

### Container Platforms

#### ✅ Docker
- **Base Images**:
  - Alpine Linux (minimal footprint)
  - Ubuntu (full compatibility)
- **Architectures**:
  - linux/amd64
  - linux/arm64
- **Orchestration**: Compatible with Docker Compose, Docker Swarm

#### ✅ Kubernetes
- **Versions**: 1.19+
- **Supported**: Deployments, StatefulSets, DaemonSets
- **Architectures**: linux/amd64, linux/arm64

## Installation Methods by Platform

### Linux

#### Package Managers (Coming Soon)
```bash
# APT (Debian/Ubuntu)
sudo apt-get update
sudo apt-get install allora

# YUM/DNF (RHEL/Fedora/CentOS)
sudo yum install allora

# Pacman (Arch)
sudo pacman -S allora

# Snap (Universal)
sudo snap install allora
```

#### Direct Binary Installation
```bash
# AMD64
curl -L https://github.com/AlloraAi/AlloraCLI/releases/latest/download/allora-linux-amd64 -o allora
chmod +x allora
sudo mv allora /usr/local/bin/

# ARM64
curl -L https://github.com/AlloraAi/AlloraCLI/releases/latest/download/allora-linux-arm64 -o allora
chmod +x allora
sudo mv allora /usr/local/bin/
```

#### Automated Script
```bash
curl -L https://raw.githubusercontent.com/AlloraAi/AlloraCLI/main/scripts/install.sh | bash
```

### macOS

#### Package Managers
```bash
# Homebrew (Coming Soon)
brew tap AlloraAi/tap
brew install allora

# MacPorts (Coming Soon)
sudo port install allora
```

#### Direct Binary Installation
```bash
# Intel Macs
curl -L https://github.com/AlloraAi/AlloraCLI/releases/latest/download/allora-darwin-amd64 -o allora
chmod +x allora
sudo mv allora /usr/local/bin/

# Apple Silicon (M1/M2/M3)
curl -L https://github.com/AlloraAi/AlloraCLI/releases/latest/download/allora-darwin-arm64 -o allora
chmod +x allora
sudo mv allora /usr/local/bin/
```

#### Automated Script
```bash
curl -L https://raw.githubusercontent.com/AlloraAi/AlloraCLI/main/scripts/install.sh | bash
```

### Windows

#### Package Managers
```powershell
# Winget (Coming Soon)
winget install AlloraAi.AlloraCLI

# Chocolatey (Coming Soon)
choco install allora

# Scoop (Coming Soon)
scoop bucket add AlloraAi https://github.com/AlloraAi/scoop-bucket
scoop install allora
```

#### Automated PowerShell Script (Recommended)
```powershell
Invoke-WebRequest -Uri "https://raw.githubusercontent.com/AlloraAi/AlloraCLI/main/scripts/install.ps1" -OutFile "$env:TEMP\install-allora.ps1"
& "$env:TEMP\install-allora.ps1"
```

#### Manual Installation
```powershell
Invoke-WebRequest -Uri "https://github.com/AlloraAi/AlloraCLI/releases/latest/download/allora-windows-amd64.exe" -OutFile "allora.exe"
```

See [WINDOWS_INSTALLATION.md](WINDOWS_INSTALLATION.md) for detailed instructions.

## Configuration Directories

AlloraCLI uses platform-specific directories for configuration:

### Linux
```
$HOME/.config/alloracli/
├── config.yaml          # Main configuration
├── plugins/             # Plugin directory
└── logs/                # Application logs
```

### macOS
```
$HOME/.config/alloracli/
├── config.yaml          # Main configuration
├── plugins/             # Plugin directory
└── logs/                # Application logs
```

### Windows
```
%APPDATA%\alloracli\
├── config.yaml          # Main configuration
├── plugins\             # Plugin directory
└── logs\                # Application logs
```

**Fallback on Windows**: If `APPDATA` is not set, uses `%USERPROFILE%\AppData\Roaming\alloracli\`

## Environment Variables

AlloraCLI supports configuration via environment variables:

### Universal Variables
```bash
# Configuration file location
ALLORA_CONFIG=/path/to/config.yaml

# Logging level (debug, info, warn, error)
ALLORA_LOG_LEVEL=info

# API Keys
ALLORA_API_KEY=your_api_key
ALLORA_OPENAI_API_KEY=your_openai_key
```

### Cloud Provider Variables
```bash
# AWS
AWS_REGION=us-west-2
AWS_PROFILE=default
AWS_ACCESS_KEY_ID=your_key
AWS_SECRET_ACCESS_KEY=your_secret

# Azure
AZURE_SUBSCRIPTION_ID=your_subscription
AZURE_TENANT_ID=your_tenant
AZURE_CLIENT_ID=your_client
AZURE_CLIENT_SECRET=your_secret

# GCP
GOOGLE_CLOUD_PROJECT=your_project
GOOGLE_APPLICATION_CREDENTIALS=/path/to/credentials.json
```

## Platform-Specific Features

### Linux-Specific
- ✅ Native systemd integration
- ✅ Full shell completion (bash, zsh, fish)
- ✅ SELinux compatible
- ✅ AppArmor compatible

### macOS-Specific
- ✅ Keychain integration for secure credential storage
- ✅ Notarized binaries (code signed)
- ✅ Launch Agent support for background services
- ✅ Native zsh completion

### Windows-Specific
- ✅ Windows Credential Manager integration
- ✅ PowerShell completion
- ✅ Windows Service support
- ✅ Event Log integration

## Building from Source

### Prerequisites
- Go 1.21 or higher
- Git
- Make (optional but recommended)

### Universal Build Instructions
```bash
# Clone repository
git clone https://github.com/AlloraAi/AlloraCLI.git
cd AlloraCLI

# Install dependencies
go mod download

# Build for current platform
go build -o allora ./cmd/allora

# Or use Make
make build
```

### Cross-Compilation

#### Build for All Platforms
```bash
make build-all
```

#### Build for Specific Platform
```bash
# Linux AMD64
GOOS=linux GOARCH=amd64 go build -o allora-linux-amd64 ./cmd/allora

# Linux ARM64
GOOS=linux GOARCH=arm64 go build -o allora-linux-arm64 ./cmd/allora

# macOS Intel
GOOS=darwin GOARCH=amd64 go build -o allora-darwin-amd64 ./cmd/allora

# macOS Apple Silicon
GOOS=darwin GOARCH=arm64 go build -o allora-darwin-arm64 ./cmd/allora

# Windows AMD64
GOOS=windows GOARCH=amd64 go build -o allora-windows-amd64.exe ./cmd/allora
```

## Container Deployment

### Docker
```dockerfile
# Using official image
docker pull alloraai/alloracli:latest

# Run container
docker run -it --rm alloraai/alloracli:latest allora --help

# With configuration
docker run -it --rm \
  -v ~/.config/alloracli:/root/.config/alloracli \
  -e ALLORA_API_KEY=your_key \
  alloraai/alloracli:latest allora ask "status"
```

### Docker Compose
```yaml
version: '3.8'
services:
  alloracli:
    image: alloraai/alloracli:latest
    volumes:
      - ./config:/root/.config/alloracli
    environment:
      - ALLORA_API_KEY=${ALLORA_API_KEY}
      - AWS_REGION=us-west-2
    command: allora monitor
```

### Kubernetes
```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: alloracli
spec:
  replicas: 1
  selector:
    matchLabels:
      app: alloracli
  template:
    metadata:
      labels:
        app: alloracli
    spec:
      containers:
      - name: alloracli
        image: alloraai/alloracli:latest
        env:
        - name: ALLORA_API_KEY
          valueFrom:
            secretKeyRef:
              name: allora-secrets
              key: api-key
```

## Troubleshooting

### Platform Detection Issues

#### Linux
```bash
# Check platform
uname -s -m
# Expected: Linux x86_64 or Linux aarch64

# Check binary architecture
file /usr/local/bin/allora
```

#### macOS
```bash
# Check platform
uname -s -m
# Expected: Darwin x86_64 or Darwin arm64

# Check if binary is compatible
file /usr/local/bin/allora

# For Apple Silicon, ensure using ARM64 binary
arch -arm64 allora --version
```

#### Windows
```powershell
# Check platform
[System.Environment]::OSVersion.Platform
[System.Environment]::Is64BitOperatingSystem

# Check PowerShell version
$PSVersionTable.PSVersion
```

### Configuration Directory Issues

#### Check directory location
```bash
# Linux/macOS
allora config show | grep "Config file"

# Windows
allora config show | Select-String "Config file"
```

#### Manually create directory
```bash
# Linux/macOS
mkdir -p ~/.config/alloracli

# Windows PowerShell
New-Item -ItemType Directory -Path "$env:APPDATA\alloracli" -Force
```

### Permission Issues

#### Linux/macOS
```bash
# Fix binary permissions
chmod +x /usr/local/bin/allora

# Fix config directory permissions
chmod 700 ~/.config/alloracli
chmod 600 ~/.config/alloracli/config.yaml
```

#### Windows
```powershell
# Run as Administrator
Start-Process powershell -Verb RunAs

# Check execution policy
Get-ExecutionPolicy
Set-ExecutionPolicy RemoteSigned -Scope CurrentUser
```

## Performance Considerations

### Memory Usage by Platform
- Linux: ~30-50 MB baseline
- macOS: ~40-60 MB baseline
- Windows: ~50-70 MB baseline
- Container: ~35-55 MB baseline (Alpine)

### Startup Time by Platform
- Linux: ~50-80ms
- macOS: ~60-100ms
- Windows: ~80-120ms
- Container: ~60-90ms

## Security Considerations

### Linux
- SELinux policies available
- AppArmor profiles available
- Runs with minimal privileges
- Credential storage in secure keyring

### macOS
- Code signed and notarized
- Keychain integration for credentials
- Sandboxing compatible
- Gatekeeper approved

### Windows
- Authenticode signed
- Credential Manager integration
- Windows Defender compatible
- SmartScreen approved

## Support Matrix

| Platform | Version | Status | Notes |
|----------|---------|--------|-------|
| Ubuntu 20.04+ | AMD64/ARM64 | ✅ Fully Supported | LTS recommended |
| Debian 10+ | AMD64/ARM64 | ✅ Fully Supported | Stable releases |
| CentOS 8+ | AMD64 | ✅ Fully Supported | Use Stream for latest |
| RHEL 8+ | AMD64 | ✅ Fully Supported | Enterprise support |
| Fedora 35+ | AMD64/ARM64 | ✅ Fully Supported | Latest features |
| Arch Linux | AMD64/ARM64 | ✅ Fully Supported | Rolling release |
| Alpine 3.14+ | AMD64/ARM64 | ✅ Fully Supported | Container optimized |
| macOS 10.15+ | Intel | ✅ Fully Supported | Catalina and later |
| macOS 11+ | Apple Silicon | ✅ Fully Supported | Native ARM64 |
| Windows 10 | AMD64 | ✅ Fully Supported | Version 1809+ |
| Windows 11 | AMD64 | ✅ Fully Supported | All versions |
| Docker | linux/amd64 | ✅ Fully Supported | Multi-arch |
| Docker | linux/arm64 | ✅ Fully Supported | Multi-arch |
| Kubernetes | 1.19+ | ✅ Fully Supported | All versions |

## Getting Help

- **Documentation**: [GitHub Repository](https://github.com/AlloraAi/AlloraCLI)
- **Issues**: [GitHub Issues](https://github.com/AlloraAi/AlloraCLI/issues)
- **Discussions**: [GitHub Discussions](https://github.com/AlloraAi/AlloraCLI/discussions)
- **Platform-Specific**:
  - Windows: [WINDOWS_INSTALLATION.md](WINDOWS_INSTALLATION.md)
  - General: [README.md](README.md)
  - FAQ: [docs/faq.md](docs/faq.md)

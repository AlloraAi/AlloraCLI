# 🚀 AlloraCLI Windows Quick Install Reference

## Automated One-Command Installation (Recommended)

```powershell
# Download and run the installation script
Invoke-WebRequest -Uri "https://raw.githubusercontent.com/AlloraAi/AlloraCLI/main/scripts/install.ps1" -OutFile "$env:TEMP\install-allora.ps1"; & "$env:TEMP\install-allora.ps1"
```

This script will:
- ✅ Download the latest version
- ✅ Install to `C:\Tools\`
- ✅ Add to PATH automatically
- ✅ Create configuration directory
- ✅ Verify installation

**After installation, restart your terminal and run:**
```powershell
allora --version
allora init
```

---

## Manual Installation

### One-Command Installation

```powershell
# Download to current directory
Invoke-WebRequest -Uri "https://github.com/AlloraAi/AlloraCLI/releases/latest/download/allora-windows-amd64.exe" -OutFile "allora.exe"

# Test installation
.\allora.exe --version
```

## Setup in 3 Steps

### 1️⃣ Download & Install
```powershell
# Create tools directory
mkdir C:\Tools
cd C:\Tools

# Download AlloraCLI
Invoke-WebRequest -Uri "https://github.com/AlloraAi/AlloraCLI/releases/latest/download/allora-windows-amd64.exe" -OutFile "allora.exe"
```

### 2️⃣ Add to PATH
```powershell
# Add C:\Tools to PATH (run as admin)
$path = [Environment]::GetEnvironmentVariable("PATH", "Machine")
[Environment]::SetEnvironmentVariable("PATH", "$path;C:\Tools", "Machine")

# Restart PowerShell, then test
allora --version
```

### 3️⃣ Initialize
```powershell
# Initialize AlloraCLI
allora init

# Verify setup
allora config show
```

## 🔧 Troubleshooting Quick Fixes

### Command Not Found
```powershell
# Use full path
C:\Tools\allora.exe --version

# Or add to current session PATH
$env:PATH += ";C:\Tools"
```

### Permission Issues
```powershell
# Run PowerShell as Administrator
# Right-click PowerShell → "Run as Administrator"

# Check execution policy
Get-ExecutionPolicy
Set-ExecutionPolicy RemoteSigned -Scope CurrentUser
```

### Configuration Issues
```powershell
# Reset configuration
Remove-Item -Path "$env:APPDATA\alloracli" -Recurse -Force
allora init
```

## 📱 Quick Commands Reference

```powershell
# Basic commands
allora --help              # Show all commands
allora --version           # Show version
allora init               # Initialize configuration
allora config show       # Show current config

# AI features
allora ask "help me"      # Ask AI assistant
allora ask --help         # Ask command options

# Cloud management
allora cloud --help       # Cloud commands
allora monitor --help     # Monitoring commands
```

## 🌐 Configure API Keys

```powershell
# OpenAI (for AI features)
[Environment]::SetEnvironmentVariable("OPENAI_API_KEY", "your-key", "User")

# Or use config
allora config set ai.openai.api_key "your-key"
```

## 📞 Get Help

- **Full Guide:** [WINDOWS_INSTALLATION.md](WINDOWS_INSTALLATION.md)
- **GitHub Issues:** https://github.com/AlloraAi/AlloraCLI/issues
- **Documentation:** https://github.com/AlloraAi/AlloraCLI

## ✅ Success Checklist

- [ ] Downloaded allora.exe (62.9 MB)
- [ ] Renamed to allora.exe
- [ ] Moved to C:\Tools\
- [ ] Added C:\Tools to PATH
- [ ] `allora --version` works
- [ ] `allora init` completed
- [ ] Configuration created in %APPDATA%\alloracli\

**You're ready to use AlloraCLI! 🎉**

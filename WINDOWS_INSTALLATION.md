# 🪟 AlloraCLI Windows Installation Guide

Complete step-by-step installation guide for Windows users.

---

## 📋 Prerequisites

Before installing AlloraCLI, ensure you have:

- **Windows 10** or **Windows 11** (64-bit)
- **Administrator privileges** (for some installation methods)
- **Internet connection** for downloading
- **PowerShell** (pre-installed on Windows 10/11)

---

## 🚀 Installation Methods

### Method 1: Automated Installation (Recommended)

This is the easiest method using PowerShell script.

#### Step 1: Download and Run Installation Script

1. **Open PowerShell** (as Administrator recommended):
   - Press `Windows + X`
   - Select "Windows PowerShell (Admin)" or "Terminal (Admin)"

2. **Run the installation script:**
   ```powershell
   # Download and run the installation script
   Invoke-WebRequest -Uri "https://raw.githubusercontent.com/AlloraAi/AlloraCLI/main/scripts/install.ps1" -OutFile "$env:TEMP\install-allora.ps1"
   
   # Run the script
   & "$env:TEMP\install-allora.ps1"
   ```

3. **The script will automatically:**
   - Fetch the latest release
   - Download the binary
   - Install to `C:\Tools\`
   - Add to PATH
   - Create configuration directory
   - Verify installation

4. **Restart your terminal** to apply PATH changes

#### Step 2: Verify Installation

```powershell
allora --version
allora --help
```

---

### Method 2: Download Pre-built Binary (Manual)

This is the easiest and fastest method for most users.

#### Step 1: Download the Binary

1. **Open your web browser** (Chrome, Edge, Firefox, etc.)

2. **Navigate to the releases page:**
   ```
   https://github.com/AlloraAi/AlloraCLI/releases/latest
   ```

3. **Download the Windows binary:**
   - Look for `allora-windows-amd64.exe` (62.9 MB)
   - Click on the file name to download
   - Save it to your `Downloads` folder

#### Step 2: Rename and Move the Binary

1. **Open File Explorer** (`Windows + E`)

2. **Navigate to your Downloads folder:**
   ```
   C:\Users\[YourUsername]\Downloads\
   ```

3. **Rename the file:**
   - Right-click on `allora-windows-amd64.exe`
   - Select "Rename"
   - Change the name to `allora.exe`

4. **Create a tools directory (recommended):**
   ```
   C:\Tools\
   ```
   - Open File Explorer
   - Navigate to `C:\`
   - Right-click → New → Folder
   - Name it `Tools`

5. **Move the file:**
   - Cut `allora.exe` from Downloads
   - Paste it to `C:\Tools\`

#### Step 3: Add to System PATH (Optional but Recommended)

This allows you to run `allora` from any directory.

1. **Open System Properties:**
   - Press `Windows + R`
   - Type `sysdm.cpl`
   - Press Enter

2. **Open Environment Variables:**
   - Click "Advanced" tab
   - Click "Environment Variables..." button

3. **Edit PATH variable:**
   - In "System variables" section, find "Path"
   - Click "Edit..."
   - Click "New"
   - Add: `C:\Tools`
   - Click "OK" on all dialogs

4. **Restart your terminal** or open a new PowerShell window

#### Step 4: Verify Installation

1. **Open PowerShell:**
   - Press `Windows + X`
   - Select "Windows PowerShell" or "Terminal"

2. **Test the installation:**
   ```powershell
   allora --version
   ```
   
   Expected output:
   ```
   allora version 1.0.0 (commit: unknown, date: unknown)
   ```

3. **View available commands:**
   ```powershell
   allora --help
   ```

---

### Method 2: Manual Installation (Current Directory)

If you prefer to run AlloraCLI from a specific project directory:

#### Step 1: Download to Project Directory

1. **Create or navigate to your project directory:**
   ```powershell
   # Create a new directory
   mkdir C:\AlloraCLI-Project
   cd C:\AlloraCLI-Project
   
   # Or navigate to existing directory
   cd "C:\Path\To\Your\Project"
   ```

2. **Download using PowerShell:**
   ```powershell
   Invoke-WebRequest -Uri "https://github.com/AlloraAi/AlloraCLI/releases/latest/download/allora-windows-amd64.exe" -OutFile "allora.exe"
   ```

3. **Verify download:**
   ```powershell
   dir allora.exe
   ```

#### Step 2: Test Installation

```powershell
.\allora.exe --version
.\allora.exe --help
```

---

### Method 3: Using Windows Package Manager (Winget)

*Note: This method will be available when AlloraCLI is published to winget repository.*

```powershell
# Future installation method
winget install AlloraAi.AlloraCLI
```

---

## 🔧 Initial Setup

After installation, initialize AlloraCLI:

### Step 1: Initialize Configuration

```powershell
allora init
```

This will:
- Create configuration directory: `%APPDATA%\alloracli\`
- Set up your first AI agent
- Configure basic settings

### Step 2: Verify Configuration

```powershell
allora config show
```

### Step 3: Test Basic Functionality

```powershell
# Test help system
allora help

# Test configuration
allora config --help

# Test AI features (basic)
allora ask --help
```

---

## 🌐 Configure Cloud Providers (Optional)

### AWS Configuration

```powershell
# Set AWS region
allora config set aws.region us-west-2

# Set AWS profile (if using AWS CLI profiles)
allora config set aws.profile default
```

### Azure Configuration

```powershell
# Set Azure subscription
allora config set azure.subscription_id "your-subscription-id"

# Set Azure tenant
allora config set azure.tenant_id "your-tenant-id"
```

### Google Cloud Configuration

```powershell
# Set GCP project
allora config set gcp.project_id "your-project-id"

# Set GCP region
allora config set gcp.region us-central1
```

---

## 🔑 API Keys Setup

### OpenAI API Key

```powershell
# Set as environment variable (recommended)
[Environment]::SetEnvironmentVariable("OPENAI_API_KEY", "your-api-key", "User")

# Or configure in AlloraCLI
allora config set ai.openai.api_key "your-api-key"
```

### Anthropic API Key

```powershell
# Set as environment variable
[Environment]::SetEnvironmentVariable("ANTHROPIC_API_KEY", "your-api-key", "User")
```

---

## 🧪 Quick Testing

Run these commands to verify everything works:

```powershell
# Basic commands
allora --version
allora help
allora config show

# AI features (requires API key)
allora ask "What can you help me with?"

# Cloud features (requires cloud credentials)
allora cloud --help
```

---

## 🔧 Troubleshooting

### Common Issues

#### Issue: "allora is not recognized as a command"

**Solution:**
1. Verify the file location: `C:\Tools\allora.exe`
2. Check PATH environment variable includes `C:\Tools`
3. Restart PowerShell/Terminal
4. Use full path: `C:\Tools\allora.exe`

#### Issue: "Execution policy error"

**Solution:**
```powershell
# Check current policy
Get-ExecutionPolicy

# Set policy to allow local scripts
Set-ExecutionPolicy -ExecutionPolicy RemoteSigned -Scope CurrentUser
```

#### Issue: "Access denied" or "Permission error"

**Solution:**
1. Run PowerShell as Administrator
2. Ensure antivirus isn't blocking the executable
3. Check file permissions on the allora.exe file

#### Issue: Configuration file errors

**Solution:**
```powershell
# Reset configuration
Remove-Item -Path "$env:APPDATA\alloracli" -Recurse -Force
allora init
```

### Getting Help

1. **Built-in help:**
   ```powershell
   allora help
   allora [command] --help
   ```

2. **Configuration doctor:**
   ```powershell
   allora config doctor
   ```

3. **Verbose output:**
   ```powershell
   allora --verbose [command]
   ```

4. **GitHub Issues:**
   - Visit: https://github.com/AlloraAi/AlloraCLI/issues
   - Search existing issues or create a new one

---

## 📁 File Locations

### Configuration Files
- **Main config:** `%APPDATA%\alloracli\config.yaml`
- **Logs:** `%APPDATA%\alloracli\logs\`
- **Cache:** `%LOCALAPPDATA%\alloracli\cache\`

### Executable Location
- **Recommended:** `C:\Tools\allora.exe`
- **Alternative:** `C:\Program Files\AlloraCLI\allora.exe`
- **Project-specific:** `[Your-Project]\allora.exe`

---

## 🔄 Updating AlloraCLI

### Manual Update

1. **Download the latest version** from GitHub releases
2. **Replace the existing executable:**
   ```powershell
   # Backup current version (optional)
   copy C:\Tools\allora.exe C:\Tools\allora-backup.exe
   
   # Replace with new version
   copy Downloads\allora-windows-amd64.exe C:\Tools\allora.exe
   ```

3. **Verify the update:**
   ```powershell
   allora --version
   ```

### Future: Automatic Updates

*Note: Auto-update feature will be available in future releases.*

---

## 🚨 Security Considerations

1. **Download only from official sources:**
   - GitHub releases: https://github.com/AlloraAi/AlloraCLI/releases
   - Official website: https://alloraai.com

2. **Verify checksums** (if provided):
   ```powershell
   # Check file hash
   Get-FileHash C:\Tools\allora.exe -Algorithm SHA256
   ```

3. **Run with least privileges** when possible

4. **Keep API keys secure:**
   - Use environment variables
   - Don't commit keys to version control
   - Rotate keys regularly

---

## ✅ Installation Complete!

You now have AlloraCLI installed and ready to use. Next steps:

1. **Read the documentation:** `README.md`
2. **Follow the getting started guide:** `docs/getting-started.md`
3. **Try the examples:** `examples/`
4. **Join the community:** GitHub Discussions

Welcome to AI-powered infrastructure management! 🎉

---

## 📞 Support

- **Documentation:** [GitHub Repository](https://github.com/AlloraAi/AlloraCLI)
- **Issues:** [GitHub Issues](https://github.com/AlloraAi/AlloraCLI/issues)
- **Discussions:** [GitHub Discussions](https://github.com/AlloraAi/AlloraCLI/discussions)
- **Website:** [AlloraAi.com](https://alloraai.com)

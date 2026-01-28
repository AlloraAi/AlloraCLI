#!/bin/bash

# AlloraCLI Installation Script
# This script installs AlloraCLI on Linux and macOS systems

set -e

# Configuration
REPO_OWNER="AlloraAi"
REPO_NAME="AlloraCLI"
BINARY_NAME="allora"
INSTALL_DIR="/usr/local/bin"
CONFIG_DIR="$HOME/.config/alloracli"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Helper functions
log_info() {
    echo -e "${GREEN}[INFO]${NC} $1"
}

log_warn() {
    echo -e "${YELLOW}[WARN]${NC} $1"
}

log_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# Check if running as root for system-wide installation
check_permissions() {
    if [[ "$INSTALL_DIR" == "/usr/local/bin" ]] && [[ $EUID -ne 0 ]]; then
        log_warn "Installing to system directory requires sudo privileges"
        INSTALL_DIR="$HOME/.local/bin"
        log_info "Installing to user directory: $INSTALL_DIR"
        mkdir -p "$INSTALL_DIR"
    fi
}

# Detect OS and architecture
detect_platform() {
    local os
    local arch
    os=$(uname -s | tr '[:upper:]' '[:lower:]')
    arch=$(uname -m)
    
    case $os in
        linux*)
            OS="linux"
            ;;
        darwin*)
            OS="darwin"
            ;;
        *)
            log_error "Unsupported operating system: $os"
            exit 1
            ;;
    esac
    
    case $arch in
        x86_64)
            ARCH="amd64"
            ;;
        arm64|aarch64)
            ARCH="arm64"
            ;;
        *)
            log_error "Unsupported architecture: $arch"
            exit 1
            ;;
    esac
    
    log_info "Detected platform: $OS-$ARCH"
}

# Get the latest release version
get_latest_version() {
    log_info "Fetching latest release information..."
    
    if command -v curl &> /dev/null; then
        VERSION=$(curl -s "https://api.github.com/repos/$REPO_OWNER/$REPO_NAME/releases/latest" | grep '"tag_name":' | sed -E 's/.*"([^"]+)".*/\1/')
    elif command -v wget &> /dev/null; then
        VERSION=$(wget -qO- "https://api.github.com/repos/$REPO_OWNER/$REPO_NAME/releases/latest" | grep '"tag_name":' | sed -E 's/.*"([^"]+)".*/\1/')
    else
        log_error "Neither curl nor wget is available. Please install one of them."
        exit 1
    fi
    
    if [[ -z "$VERSION" ]]; then
        log_error "Failed to fetch latest version"
        exit 1
    fi
    
    log_info "Latest version: $VERSION"
}

# Download the binary
download_binary() {
    local binary_name="$BINARY_NAME-$OS-$ARCH"
    local download_url="https://github.com/$REPO_OWNER/$REPO_NAME/releases/download/$VERSION/$binary_name.tar.gz"
    local temp_dir
    temp_dir=$(mktemp -d)
    local archive_path="$temp_dir/$binary_name.tar.gz"
    
    log_info "Downloading AlloraCLI $VERSION..."
    log_info "URL: $download_url"
    
    if command -v curl &> /dev/null; then
        curl -L -o "$archive_path" "$download_url"
    elif command -v wget &> /dev/null; then
        wget -O "$archive_path" "$download_url"
    else
        log_error "Neither curl nor wget is available"
        exit 1
    fi
    
    if [[ ! -f "$archive_path" ]]; then
        log_error "Failed to download binary"
        exit 1
    fi
    
    log_info "Extracting binary..."
    tar -xzf "$archive_path" -C "$temp_dir"
    
    # Find the binary in the extracted files
    local binary_path="$temp_dir/$binary_name"
    if [[ ! -f "$binary_path" ]]; then
        # Try without extension
        binary_path="$temp_dir/$BINARY_NAME"
        if [[ ! -f "$binary_path" ]]; then
            log_error "Binary not found in archive"
            exit 1
        fi
    fi
    
    log_info "Installing binary to $INSTALL_DIR..."
    if [[ "$INSTALL_DIR" == "/usr/local/bin" ]]; then
        sudo cp "$binary_path" "$INSTALL_DIR/$BINARY_NAME"
        sudo chmod +x "$INSTALL_DIR/$BINARY_NAME"
    else
        cp "$binary_path" "$INSTALL_DIR/$BINARY_NAME"
        chmod +x "$INSTALL_DIR/$BINARY_NAME"
    fi
    
    # Clean up
    rm -rf "$temp_dir"
    
    log_info "Binary installed successfully"
}

# Create configuration directory
setup_config() {
    log_info "Setting up configuration directory..."
    
    if [[ ! -d "$CONFIG_DIR" ]]; then
        mkdir -p "$CONFIG_DIR"
        log_info "Created configuration directory: $CONFIG_DIR"
    else
        log_info "Configuration directory already exists: $CONFIG_DIR"
    fi
    
    # Create plugins directory
    mkdir -p "$CONFIG_DIR/plugins"
    
    # Create a basic config file if it doesn't exist
    local config_file="$CONFIG_DIR/config.yaml"
    if [[ ! -f "$config_file" ]]; then
        cat > "$config_file" << EOF
version: "2.0.0"
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
EOF
        log_info "Created basic configuration file: $config_file"
    fi
}

# Add to PATH if necessary
setup_path() {
    if [[ "$INSTALL_DIR" == "$HOME/.local/bin" ]]; then
        # Check if ~/.local/bin is in PATH
        if [[ ":$PATH:" != *":$HOME/.local/bin:"* ]]; then
            log_warn "$HOME/.local/bin is not in your PATH"
            
            # Add to shell profile
            local shell_profile=""
            if [[ -f "$HOME/.bashrc" ]]; then
                shell_profile="$HOME/.bashrc"
            elif [[ -f "$HOME/.zshrc" ]]; then
                shell_profile="$HOME/.zshrc"
            elif [[ -f "$HOME/.profile" ]]; then
                shell_profile="$HOME/.profile"
            fi
            
            if [[ -n "$shell_profile" ]]; then
                echo "export PATH=\"\$HOME/.local/bin:\$PATH\"" >> "$shell_profile"
                log_info "Added $HOME/.local/bin to PATH in $shell_profile"
                log_info "Please run 'source $shell_profile' or restart your terminal"
            else
                log_warn "Could not find shell profile to update PATH"
                log_warn "Please manually add $HOME/.local/bin to your PATH"
            fi
        fi
    fi
}

# Verify installation
verify_installation() {
    log_info "Verifying installation..."
    
    if command -v "$BINARY_NAME" &> /dev/null; then
        local version_output
        version_output=$("$BINARY_NAME" --version 2>&1)
        log_info "Installation successful!"
        log_info "Version: $version_output"
        
        log_info "You can now run 'allora --help' to get started"
        log_info "To initialize AlloraCLI, run 'allora init'"
    else
        log_error "Installation verification failed"
        log_error "Binary not found in PATH"
        exit 1
    fi
}

# Main installation function
main() {
    log_info "AlloraCLI Installation Script"
    log_info "=============================="
    
    check_permissions
    detect_platform
    get_latest_version
    download_binary
    setup_config
    setup_path
    verify_installation
    
    log_info ""
    log_info "Installation completed successfully!"
    log_info "Documentation: https://docs.alloraai.com"
    log_info "GitHub: https://github.com/$REPO_OWNER/$REPO_NAME"
}

# Handle command line arguments
case "${1:-}" in
    --help|-h)
        echo "AlloraCLI Installation Script"
        echo ""
        echo "Usage: $0 [OPTIONS]"
        echo ""
        echo "Options:"
        echo "  --help, -h    Show this help message"
        echo "  --version     Install specific version"
        echo "  --local       Install to ~/.local/bin (default if not root)"
        echo "  --system      Install to /usr/local/bin (requires sudo)"
        echo ""
        echo "Environment variables:"
        echo "  ALLORA_INSTALL_DIR    Custom installation directory"
        echo "  ALLORA_VERSION        Specific version to install"
        echo ""
        exit 0
        ;;
    --version)
        if [[ -n "${2:-}" ]]; then
            VERSION="$2"
            log_info "Installing specific version: $VERSION"
        else
            log_error "Version not specified"
            exit 1
        fi
        ;;
    --local)
        INSTALL_DIR="$HOME/.local/bin"
        ;;
    --system)
        INSTALL_DIR="/usr/local/bin"
        ;;
esac

# Override with environment variables if set
if [[ -n "${ALLORA_INSTALL_DIR:-}" ]]; then
    INSTALL_DIR="$ALLORA_INSTALL_DIR"
fi

if [[ -n "${ALLORA_VERSION:-}" ]]; then
    VERSION="$ALLORA_VERSION"
fi

# Run main installation
main

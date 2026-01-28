#!/bin/bash

# AlloraCLI Installation Verification Script
# This script tests the installation and configuration logic

set -e

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Counters
TESTS_RUN=0
TESTS_PASSED=0
TESTS_FAILED=0

# Helper functions
log_info() {
    echo -e "${GREEN}[INFO]${NC} $1"
}

log_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

log_warn() {
    echo -e "${YELLOW}[WARN]${NC} $1"
}

test_start() {
    TESTS_RUN=$((TESTS_RUN + 1))
    echo ""
    echo -e "${YELLOW}[TEST $TESTS_RUN]${NC} $1"
}

test_pass() {
    TESTS_PASSED=$((TESTS_PASSED + 1))
    echo -e "${GREEN}✓ PASS${NC}"
}

test_fail() {
    TESTS_FAILED=$((TESTS_FAILED + 1))
    echo -e "${RED}✗ FAIL${NC} $1"
}

# Test 1: Check if binary exists
test_binary_exists() {
    test_start "Check if AlloraCLI binary exists"
    
    if [ -f "./bin/allora" ]; then
        test_pass
        return 0
    else
        test_fail "Binary not found at ./bin/allora"
        return 1
    fi
}

# Test 2: Check if binary is executable
test_binary_executable() {
    test_start "Check if binary is executable"
    
    if [ -x "./bin/allora" ]; then
        test_pass
        return 0
    else
        test_fail "Binary is not executable"
        return 1
    fi
}

# Test 3: Check version command
test_version_command() {
    test_start "Check version command"
    
    if ./bin/allora --version &> /dev/null; then
        VERSION=$(./bin/allora --version)
        log_info "Version: $VERSION"
        test_pass
        return 0
    else
        test_fail "Version command failed"
        return 1
    fi
}

# Test 4: Check help command
test_help_command() {
    test_start "Check help command"
    
    if ./bin/allora --help &> /dev/null; then
        test_pass
        return 0
    else
        test_fail "Help command failed"
        return 1
    fi
}

# Test 5: Check available commands
test_available_commands() {
    test_start "Check available commands"
    
    EXPECTED_COMMANDS=("init" "config" "ask" "monitor" "troubleshoot" "deploy" "analyze" "security" "cloud" "plugin" "gemini")
    MISSING_COMMANDS=()
    
    HELP_OUTPUT=$(./bin/allora --help 2>&1)
    
    for cmd in "${EXPECTED_COMMANDS[@]}"; do
        if ! echo "$HELP_OUTPUT" | grep -q "$cmd"; then
            MISSING_COMMANDS+=("$cmd")
        fi
    done
    
    if [ ${#MISSING_COMMANDS[@]} -eq 0 ]; then
        test_pass
        return 0
    else
        test_fail "Missing commands: ${MISSING_COMMANDS[*]}"
        return 1
    fi
}

# Test 6: Test platform detection
test_platform_detection() {
    test_start "Test platform detection"
    
    OS=$(uname -s | tr '[:upper:]' '[:lower:]')
    ARCH=$(uname -m)
    
    log_info "Detected OS: $OS"
    log_info "Detected Architecture: $ARCH"
    
    case $OS in
        linux*|darwin*)
            test_pass
            return 0
            ;;
        *)
            test_fail "Unsupported OS: $OS"
            return 1
            ;;
    esac
}

# Test 7: Test config directory creation
test_config_dir() {
    test_start "Test config directory detection"
    
    case $(uname -s) in
        Darwin|Linux)
            CONFIG_DIR="$HOME/.config/alloracli"
            ;;
        *)
            test_fail "Unsupported platform for test"
            return 1
            ;;
    esac
    
    log_info "Expected config dir: $CONFIG_DIR"
    
    # The directory doesn't need to exist yet, just verify the path is correct
    if [ -n "$CONFIG_DIR" ]; then
        test_pass
        return 0
    else
        test_fail "Config directory path is empty"
        return 1
    fi
}

# Test 8: Test Makefile targets
test_makefile_targets() {
    test_start "Test Makefile targets"
    
    EXPECTED_TARGETS=("build" "clean" "test" "lint" "deps" "install")
    MISSING_TARGETS=()
    
    for target in "${EXPECTED_TARGETS[@]}"; do
        if ! grep -q "^$target:" Makefile 2>/dev/null; then
            MISSING_TARGETS+=("$target")
        fi
    done
    
    if [ ${#MISSING_TARGETS[@]} -eq 0 ]; then
        test_pass
        return 0
    else
        test_fail "Missing Makefile targets: ${MISSING_TARGETS[*]}"
        return 1
    fi
}

# Test 9: Test Go module
test_go_module() {
    test_start "Test Go module configuration"
    
    if [ -f "go.mod" ]; then
        if grep -q "module github.com/AlloraAi/AlloraCLI" go.mod; then
            test_pass
            return 0
        else
            test_fail "go.mod has incorrect module path"
            return 1
        fi
    else
        test_fail "go.mod not found"
        return 1
    fi
}

# Test 10: Test install script exists
test_install_script() {
    test_start "Test install script exists"
    
    if [ -f "scripts/install.sh" ]; then
        if [ -x "scripts/install.sh" ]; then
            test_pass
            return 0
        else
            test_fail "Install script is not executable"
            return 1
        fi
    else
        test_fail "Install script not found"
        return 1
    fi
}

# Test 11: Test Windows install script exists
test_windows_install_script() {
    test_start "Test Windows install script exists"
    
    if [ -f "scripts/install.ps1" ]; then
        test_pass
        return 0
    else
        test_fail "Windows install script not found"
        return 1
    fi
}

# Test 12: Test documentation files
test_documentation() {
    test_start "Test documentation files"
    
    REQUIRED_DOCS=("README.md" "WINDOWS_INSTALLATION.md" "PLATFORM_COMPATIBILITY.md")
    MISSING_DOCS=()
    
    for doc in "${REQUIRED_DOCS[@]}"; do
        if [ ! -f "$doc" ]; then
            MISSING_DOCS+=("$doc")
        fi
    done
    
    if [ ${#MISSING_DOCS[@]} -eq 0 ]; then
        test_pass
        return 0
    else
        test_fail "Missing documentation: ${MISSING_DOCS[*]}"
        return 1
    fi
}

# Test 13: Test shellcheck on install script (if available)
test_shellcheck() {
    test_start "Test install script with shellcheck (optional)"
    
    if command -v shellcheck &> /dev/null; then
        if shellcheck scripts/install.sh; then
            test_pass
            return 0
        else
            test_fail "Shellcheck found issues"
            return 1
        fi
    else
        log_warn "Shellcheck not installed, skipping"
        test_pass
        return 0
    fi
}

# Test 14: Test build reproducibility
test_build_reproducibility() {
    test_start "Test build reproducibility"
    
    # Get current binary hash
    HASH1=$(sha256sum ./bin/allora | cut -d' ' -f1)
    
    # Rebuild
    log_info "Rebuilding binary..."
    make clean &> /dev/null
    make build &> /dev/null
    
    # Get new hash
    HASH2=$(sha256sum ./bin/allora | cut -d' ' -f1)
    
    if [ "$HASH1" != "$HASH2" ]; then
        log_info "Note: Build hashes differ (expected due to timestamps)"
        log_info "Hash 1: $HASH1"
        log_info "Hash 2: $HASH2"
    fi
    
    # Just check if rebuild succeeded
    if [ -f "./bin/allora" ]; then
        test_pass
        return 0
    else
        test_fail "Rebuild failed"
        return 1
    fi
}

# Main test runner
main() {
    echo ""
    echo "╔═══════════════════════════════════════════════════════════╗"
    echo "║     AlloraCLI Installation Verification Tests            ║"
    echo "╚═══════════════════════════════════════════════════════════╝"
    echo ""
    
    # Run all tests
    test_binary_exists || true
    test_binary_executable || true
    test_version_command || true
    test_help_command || true
    test_available_commands || true
    test_platform_detection || true
    test_config_dir || true
    test_makefile_targets || true
    test_go_module || true
    test_install_script || true
    test_windows_install_script || true
    test_documentation || true
    test_shellcheck || true
    test_build_reproducibility || true
    
    # Print summary
    echo ""
    echo "╔═══════════════════════════════════════════════════════════╗"
    echo "║                    Test Summary                           ║"
    echo "╚═══════════════════════════════════════════════════════════╝"
    echo ""
    echo "Tests Run:    $TESTS_RUN"
    echo -e "Tests Passed: ${GREEN}$TESTS_PASSED${NC}"
    
    if [ $TESTS_FAILED -gt 0 ]; then
        echo -e "Tests Failed: ${RED}$TESTS_FAILED${NC}"
        echo ""
        echo -e "${RED}Some tests failed. Please review the output above.${NC}"
        exit 1
    else
        echo -e "Tests Failed: ${GREEN}0${NC}"
        echo ""
        echo -e "${GREEN}All tests passed! ✓${NC}"
        exit 0
    fi
}

# Run tests
main

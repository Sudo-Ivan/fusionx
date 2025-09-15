#!/bin/bash

set -e

# Configuration
REPO_OWNER="Sudo-Ivan"
REPO_NAME="fusionx"
INSTALL_DIR="/usr/local/bin"
BINARY_NAME="fusionx"
TMP_DIR="/tmp/fusionx-install"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Logging functions
log_info() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

log_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

log_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

log_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# Check if command exists
command_exists() {
    command -v "$1" >/dev/null 2>&1
}

# Detect OS and architecture
detect_platform() {
    local os arch
    
    # Detect OS
    case "$(uname -s)" in
        Linux*)     os="linux";;
        Darwin*)    os="darwin";;
        FreeBSD*)   os="freebsd";;
        OpenBSD*)   os="openbsd";;
        *)          log_error "Unsupported operating system: $(uname -s)"; exit 1;;
    esac
    
    # Detect architecture
    case "$(uname -m)" in
        x86_64|amd64)   arch="amd64";;
        arm64|aarch64)  arch="arm64";;
        armv6l)         arch="armv6";;
        armv7l)         arch="armv7";;
        *)              log_error "Unsupported architecture: $(uname -m)"; exit 1;;
    esac
    
    echo "${os}_${arch}"
}

# Get latest release version
get_latest_version() {
    if command_exists curl; then
        curl -s "https://api.github.com/repos/${REPO_OWNER}/${REPO_NAME}/releases/latest" | grep '"tag_name"' | sed -E 's/.*"([^"]+)".*/\1/'
    elif command_exists wget; then
        wget -qO- "https://api.github.com/repos/${REPO_OWNER}/${REPO_NAME}/releases/latest" | grep '"tag_name"' | sed -E 's/.*"([^"]+)".*/\1/'
    else
        log_error "Neither curl nor wget is available. Please install one of them."
        exit 1
    fi
}

# Download file
download_file() {
    local url="$1"
    local output="$2"
    
    if command_exists curl; then
        curl -L -o "$output" "$url"
    elif command_exists wget; then
        wget -O "$output" "$url"
    else
        log_error "Neither curl nor wget is available."
        exit 1
    fi
}

# Verify checksum
verify_checksum() {
    local file="$1"
    local expected_checksum="$2"
    
    if command_exists sha256sum; then
        local actual_checksum=$(sha256sum "$file" | cut -d' ' -f1)
    elif command_exists shasum; then
        local actual_checksum=$(shasum -a 256 "$file" | cut -d' ' -f1)
    else
        log_warning "No SHA256 checksum tool found. Skipping verification."
        return 0
    fi
    
    if [ "$actual_checksum" = "$expected_checksum" ]; then
        log_success "Checksum verification passed"
        return 0
    else
        log_error "Checksum verification failed!"
        log_error "Expected: $expected_checksum"
        log_error "Actual:   $actual_checksum"
        return 1
    fi
}

# Check if running as root for system installation
check_permissions() {
    if [ "$EUID" -ne 0 ] && [ "$INSTALL_DIR" = "/usr/local/bin" ]; then
        log_warning "Installing to system directory requires root privileges."
        log_info "You can either:"
        log_info "  1. Run with sudo: sudo $0"
        log_info "  2. Install to user directory: INSTALL_DIR=\$HOME/.local/bin $0"
        
        read -p "Install to user directory instead? (y/n): " -n 1 -r
        echo
        if [[ $REPLY =~ ^[Yy]$ ]]; then
            INSTALL_DIR="$HOME/.local/bin"
            mkdir -p "$INSTALL_DIR"
        else
            log_error "Installation cancelled."
            exit 1
        fi
    fi
}

# Main installation function
main() {
    log_info "Starting FusionX installation..."
    
    # Check permissions
    check_permissions
    
    # Detect platform
    local platform=$(detect_platform)
    log_info "Detected platform: $platform"
    
    # Get latest version
    log_info "Fetching latest release information..."
    local version=$(get_latest_version)
    if [ -z "$version" ]; then
        log_error "Failed to get latest version"
        exit 1
    fi
    log_info "Latest version: $version"
    
    # Create temporary directory
    mkdir -p "$TMP_DIR"
    cd "$TMP_DIR"
    
    # Construct download URLs with version in filename (remove 'v' prefix for filename)
    local version_clean="${version#v}"
    local binary_name="fusion_${version_clean}_${platform}"
    local binary_url="https://github.com/${REPO_OWNER}/${REPO_NAME}/releases/download/${version}/${binary_name}"
    local checksum_url="https://github.com/${REPO_OWNER}/${REPO_NAME}/releases/download/${version}/checksums.txt"
    
    # Download binary and checksums
    log_info "Downloading binary..."
    download_file "$binary_url" "$binary_name"
    
    log_info "Downloading checksums..."
    download_file "$checksum_url" "checksums.txt"
    
    # Verify checksum
    log_info "Verifying checksum..."
    local expected_checksum=$(grep "$binary_name" checksums.txt | cut -d' ' -f1)
    if [ -z "$expected_checksum" ]; then
        log_warning "Could not find checksum for $binary_name in checksums.txt"
    else
        verify_checksum "$binary_name" "$expected_checksum"
    fi
    
    # Install binary
    log_info "Installing binary to $INSTALL_DIR..."
    chmod +x "$binary_name"
    
    # Create install directory if it doesn't exist
    mkdir -p "$INSTALL_DIR"
    
    # Move binary to install location with final name
    mv "$binary_name" "$INSTALL_DIR/$BINARY_NAME"
    
    # Clean up
    cd /
    rm -rf "$TMP_DIR"
    
    log_success "FusionX installed successfully!"
    log_info "Binary location: $INSTALL_DIR/$BINARY_NAME"
    
    # Check if install directory is in PATH
    if ! echo "$PATH" | grep -q "$INSTALL_DIR"; then
        log_warning "$INSTALL_DIR is not in your PATH"
        log_info "Add it to your PATH by adding this line to your shell profile:"
        log_info "export PATH=\"\$PATH:$INSTALL_DIR\""
    fi
    
    # Ask if user wants to start the application
    read -p "Do you want to start FusionX now? (y/n): " -n 1 -r
    echo
    if [[ $REPLY =~ ^[Yy]$ ]]; then
        log_info "Starting FusionX..."
        log_info "FusionX will be available at http://localhost:8080"
        log_info "Press Ctrl+C to stop"
        exec "$INSTALL_DIR/$BINARY_NAME"
    else
        log_info "You can start FusionX later by running: $BINARY_NAME"
        log_info "It will be available at http://localhost:8080"
    fi
}

# Handle script arguments
case "${1:-}" in
    --help|-h)
        echo "FusionX Installation Script"
        echo ""
        echo "Usage: $0 [options]"
        echo ""
        echo "Options:"
        echo "  --help, -h     Show this help message"
        echo "  --version, -v  Show version information"
        echo ""
        echo "Environment variables:"
        echo "  INSTALL_DIR    Installation directory (default: /usr/local/bin)"
        echo ""
        echo "Example:"
        echo "  curl -sSL https://raw.githubusercontent.com/${REPO_OWNER}/${REPO_NAME}/main/scripts/install.sh | bash"
        echo "  INSTALL_DIR=\$HOME/.local/bin ./install.sh"
        exit 0
        ;;
    --version|-v)
        echo "Fusion Installation Script v1.0.0"
        exit 0
        ;;
    *)
        main "$@"
        ;;
esac

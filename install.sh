#!/usr/bin/env bash
# Installation script for gohome
# Usage: curl -sSL https://raw.githubusercontent.com/anIcedAntFA/gohome/main/install.sh | bash

set -e

# Configuration
REPO="anIcedAntFA/gohome"
BINARY_NAME="gohome"
INSTALL_DIR="${INSTALL_DIR:-/usr/local/bin}"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Helper functions
info() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

warn() {
    echo -e "${YELLOW}[WARN]${NC} $1"
}

error() {
    echo -e "${RED}[ERROR]${NC} $1"
    exit 1
}

# Detect platform
detect_platform() {
    local os arch

    # Detect OS
    case "$(uname -s)" in
        Linux*)     os="linux" ;;
        Darwin*)    os="darwin" ;;
        MINGW*|MSYS*|CYGWIN*) os="windows" ;;
        *)          error "Unsupported operating system: $(uname -s)" ;;
    esac

    # Detect architecture
    case "$(uname -m)" in
        x86_64|amd64)   arch="x86_64" ;;
        aarch64|arm64)  arch="arm64" ;;
        armv7l)         arch="armv7" ;;
        i386|i686)      arch="i386" ;;
        *)              error "Unsupported architecture: $(uname -m)" ;;
    esac

    echo "${os}_${arch}"
}

# Get latest version from GitHub
get_latest_version() {
    local latest_url="https://api.github.com/repos/${REPO}/releases/latest"
    
    if command -v curl &> /dev/null; then
        curl -sL "${latest_url}" | grep '"tag_name":' | sed -E 's/.*"([^"]+)".*/\1/'
    elif command -v wget &> /dev/null; then
        wget -qO- "${latest_url}" | grep '"tag_name":' | sed -E 's/.*"([^"]+)".*/\1/'
    else
        error "Neither curl nor wget found. Please install one of them."
    fi
}

# Download and install binary
install_binary() {
    local version="$1"
    local platform="$2"
    local os="${platform%%_*}"
    local arch="${platform##*_}"
    
    info "Installing ${BINARY_NAME} ${version} for ${platform}..."

    # Construct download URL
    local archive_ext="tar.gz"
    if [ "$os" = "windows" ]; then
        archive_ext="zip"
    fi
    
    local archive_name="${BINARY_NAME}_${version#v}_${os}_${arch}.${archive_ext}"
    local download_url="https://github.com/${REPO}/releases/download/${version}/${archive_name}"
    
    info "Download URL: ${download_url}"
    
    # Create temporary directory
    local tmp_dir
    tmp_dir=$(mktemp -d)
    trap "rm -rf ${tmp_dir}" EXIT
    
    # Download archive
    info "Downloading ${archive_name}..."
    if command -v curl &> /dev/null; then
        curl -sL "${download_url}" -o "${tmp_dir}/${archive_name}" || error "Failed to download ${archive_name}"
    elif command -v wget &> /dev/null; then
        wget -q "${download_url}" -O "${tmp_dir}/${archive_name}" || error "Failed to download ${archive_name}"
    fi
    
    # Extract archive
    info "Extracting archive..."
    cd "${tmp_dir}"
    if [ "$archive_ext" = "zip" ]; then
        if command -v unzip &> /dev/null; then
            unzip -q "${archive_name}" || error "Failed to extract archive"
        else
            error "unzip command not found. Please install unzip."
        fi
    else
        if command -v tar &> /dev/null; then
            tar -xzf "${archive_name}" || error "Failed to extract archive"
        else
            error "tar command not found. Please install tar."
        fi
    fi
    
    # Find binary
    local binary_path="${tmp_dir}/${BINARY_NAME}"
    if [ "$os" = "windows" ]; then
        binary_path="${binary_path}.exe"
    fi
    
    if [ ! -f "${binary_path}" ]; then
        error "Binary not found in archive: ${binary_path}"
    fi
    
    # Make binary executable
    chmod +x "${binary_path}"
    
    # Install binary
    info "Installing to ${INSTALL_DIR}..."
    
    # Check if we need sudo
    if [ -w "${INSTALL_DIR}" ]; then
        mv "${binary_path}" "${INSTALL_DIR}/${BINARY_NAME}" || error "Failed to install binary"
    else
        warn "Requires sudo to install to ${INSTALL_DIR}"
        sudo mv "${binary_path}" "${INSTALL_DIR}/${BINARY_NAME}" || error "Failed to install binary"
    fi
    
    success "${BINARY_NAME} ${version} installed successfully!"
}

# Verify installation
verify_installation() {
    if command -v "${BINARY_NAME}" &> /dev/null; then
        local installed_version
        installed_version=$("${BINARY_NAME}" --version 2>&1 || true)
        success "Verification successful!"
        info "Installed version: ${installed_version}"
        info "Location: $(command -v ${BINARY_NAME})"
    else
        warn "Binary installed but not found in PATH."
        warn "Please add ${INSTALL_DIR} to your PATH:"
        warn "  export PATH=\"\${PATH}:${INSTALL_DIR}\""
    fi
}

# Main installation flow
main() {
    echo ""
    info "═══════════════════════════════════════════════"
    info "  ${BINARY_NAME} Installation Script"
    info "═══════════════════════════════════════════════"
    echo ""
    
    # Check dependencies
    if ! command -v curl &> /dev/null && ! command -v wget &> /dev/null; then
        error "Neither curl nor wget found. Please install one of them."
    fi
    
    # Detect platform
    local platform
    platform=$(detect_platform)
    info "Detected platform: ${platform}"
    
    # Get latest version
    info "Fetching latest version..."
    local version
    version=$(get_latest_version)
    
    if [ -z "$version" ]; then
        error "Failed to get latest version"
    fi
    
    info "Latest version: ${version}"
    
    # Check if already installed
    if command -v "${BINARY_NAME}" &> /dev/null; then
        local current_version
        current_version=$("${BINARY_NAME}" --version 2>&1 | head -n1 || echo "unknown")
        warn "Found existing installation: ${current_version}"
        
        # Ask for confirmation
        if [ -t 0 ]; then  # Check if running interactively
            read -p "Do you want to reinstall/upgrade? (y/N) " -n 1 -r
            echo
            if [[ ! $REPLY =~ ^[Yy]$ ]]; then
                info "Installation cancelled."
                exit 0
            fi
        else
            info "Proceeding with installation..."
        fi
    fi
    
    # Install binary
    install_binary "${version}" "${platform}"
    
    # Verify installation
    echo ""
    verify_installation
    
    echo ""
    info "═══════════════════════════════════════════════"
    success "Installation complete!"
    info "═══════════════════════════════════════════════"
    echo ""
    info "Get started:"
    info "  ${BINARY_NAME} --help"
    info "  ${BINARY_NAME} --version"
    echo ""
    info "Documentation:"
    info "  https://github.com/${REPO}"
    echo ""
}

# Run main function
main "$@"

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

# Ask user for confirmation
ask() {
    local prompt="$1"
    local default="${2:-N}"
    
    if [ -t 0 ]; then  # Check if running interactively
        if [ "$default" = "Y" ]; then
            read -p "$prompt (Y/n) " -n 1 -r
        else
            read -p "$prompt (y/N) " -n 1 -r
        fi
        echo
        [[ $REPLY =~ ^[Yy]$ ]]
    else
        # Non-interactive: use default
        [ "$default" = "Y" ]
    fi
}

# Check if command exists
has_command() {
    command -v "$1" &> /dev/null
}

# Check and install dependencies
check_dependencies() {
    local os="$1"
    local missing_deps=()
    local installed_deps=()
    
    # Check required tools
    if ! has_command "curl" && ! has_command "wget"; then
        missing_deps+=("curl or wget")
    fi
    
    if [ "$os" = "windows" ]; then
        if ! has_command "unzip"; then
            missing_deps+=("unzip")
        fi
    else
        if ! has_command "tar"; then
            missing_deps+=("tar")
        fi
    fi
    
    # If no missing deps, return success
    if [ ${#missing_deps[@]} -eq 0 ]; then
        return 0
    fi
    
    # Show missing dependencies
    warn "Missing required dependencies:"
    for dep in "${missing_deps[@]}"; do
        echo "  - $dep"
    done
    echo ""
    
    # Try to install dependencies
    info "Attempting to install missing dependencies..."
    
    # Detect package manager and install
    if has_command "pacman"; then
        info "Using pacman (Arch Linux)..."
        if ask "Install dependencies with pacman?" "Y"; then
            for dep in "${missing_deps[@]}"; do
                case "$dep" in
                    "curl or wget") sudo pacman -S --noconfirm curl && installed_deps+=("curl") ;;
                    "tar") sudo pacman -S --noconfirm tar && installed_deps+=("tar") ;;
                    "unzip") sudo pacman -S --noconfirm unzip && installed_deps+=("unzip") ;;
                esac
            done
        else
            error "Cannot proceed without required dependencies."
        fi
    elif has_command "apt-get"; then
        info "Using apt-get (Debian/Ubuntu)..."
        if ask "Install dependencies with apt-get?" "Y"; then
            sudo apt-get update -qq
            for dep in "${missing_deps[@]}"; do
                case "$dep" in
                    "curl or wget") sudo apt-get install -y curl && installed_deps+=("curl") ;;
                    "tar") sudo apt-get install -y tar && installed_deps+=("tar") ;;
                    "unzip") sudo apt-get install -y unzip && installed_deps+=("unzip") ;;
                esac
            done
        else
            error "Cannot proceed without required dependencies."
        fi
    elif has_command "dnf"; then
        info "Using dnf (Fedora/RHEL)..."
        if ask "Install dependencies with dnf?" "Y"; then
            for dep in "${missing_deps[@]}"; do
                case "$dep" in
                    "curl or wget") sudo dnf install -y curl && installed_deps+=("curl") ;;
                    "tar") sudo dnf install -y tar && installed_deps+=("tar") ;;
                    "unzip") sudo dnf install -y unzip && installed_deps+=("unzip") ;;
                esac
            done
        else
            error "Cannot proceed without required dependencies."
        fi
    elif has_command "yum"; then
        info "Using yum (CentOS/RHEL)..."
        if ask "Install dependencies with yum?" "Y"; then
            for dep in "${missing_deps[@]}"; do
                case "$dep" in
                    "curl or wget") sudo yum install -y curl && installed_deps+=("curl") ;;
                    "tar") sudo yum install -y tar && installed_deps+=("tar") ;;
                    "unzip") sudo yum install -y unzip && installed_deps+=("unzip") ;;
                esac
            done
        else
            error "Cannot proceed without required dependencies."
        fi
    elif has_command "zypper"; then
        info "Using zypper (openSUSE)..."
        if ask "Install dependencies with zypper?" "Y"; then
            for dep in "${missing_deps[@]}"; do
                case "$dep" in
                    "curl or wget") sudo zypper install -y curl && installed_deps+=("curl") ;;
                    "tar") sudo zypper install -y tar && installed_deps+=("tar") ;;
                    "unzip") sudo zypper install -y unzip && installed_deps+=("unzip") ;;
                esac
            done
        else
            error "Cannot proceed without required dependencies."
        fi
    elif has_command "apk"; then
        info "Using apk (Alpine Linux)..."
        if ask "Install dependencies with apk?" "Y"; then
            for dep in "${missing_deps[@]}"; do
                case "$dep" in
                    "curl or wget") sudo apk add curl && installed_deps+=("curl") ;;
                    "tar") sudo apk add tar && installed_deps+=("tar") ;;
                    "unzip") sudo apk add unzip && installed_deps+=("unzip") ;;
                esac
            done
        else
            error "Cannot proceed without required dependencies."
        fi
    elif has_command "brew"; then
        info "Using Homebrew (macOS/Linux)..."
        if ask "Install dependencies with brew?" "Y"; then
            for dep in "${missing_deps[@]}"; do
                case "$dep" in
                    "curl or wget") brew install curl && installed_deps+=("curl") ;;
                    "tar") brew install gnu-tar && installed_deps+=("tar") ;;
                    "unzip") brew install unzip && installed_deps+=("unzip") ;;
                esac
            done
        else
            error "Cannot proceed without required dependencies."
        fi
    else
        error "No supported package manager found. Please install dependencies manually: ${missing_deps[*]}"
    fi
    
    # Verify installation
    local install_failed=false
    for dep in "${missing_deps[@]}"; do
        case "$dep" in
            "curl or wget")
                if ! has_command "curl" && ! has_command "wget"; then
                    install_failed=true
                fi
                ;;
            "tar")
                if ! has_command "tar"; then
                    install_failed=true
                fi
                ;;
            "unzip")
                if ! has_command "unzip"; then
                    install_failed=true
                fi
                ;;
        esac
    done
    
    if [ "$install_failed" = true ]; then
        error "Failed to install some dependencies. Please install them manually."
    fi
    
    success "All dependencies installed successfully!"
    
    # Store installed deps for cleanup later
    echo "${installed_deps[@]}"
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
    local fallback_url="https://github.com/${REPO}/releases/latest"
    
    # Try API first
    local version=""
    if has_command "curl"; then
        version=$(curl -sSL "${latest_url}" 2>/dev/null | grep '"tag_name":' | sed -E 's/.*"([^"]+)".*/\1/' | head -1)
    elif has_command "wget"; then
        version=$(wget -qO- "${latest_url}" 2>/dev/null | grep '"tag_name":' | sed -E 's/.*"([^"]+)".*/\1/' | head -1)
    fi
    
    # If API fails, try scraping the releases page
    if [ -z "$version" ]; then
        if has_command "curl"; then
            version=$(curl -sSL "${fallback_url}" 2>/dev/null | grep -o '/releases/tag/v[0-9]\+\.[0-9]\+\.[0-9]\+' | head -1 | sed 's|.*/||')
        elif has_command "wget"; then
            version=$(wget -qO- "${fallback_url}" 2>/dev/null | grep -o '/releases/tag/v[0-9]\+\.[0-9]\+\.[0-9]\+' | head -1 | sed 's|.*/||')
        fi
    fi
    
    echo "$version"
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
    local download_failed=false
    if has_command "curl"; then
        if ! curl -fsSL "${download_url}" -o "${tmp_dir}/${archive_name}"; then
            download_failed=true
        fi
    elif has_command "wget"; then
        if ! wget -q "${download_url}" -O "${tmp_dir}/${archive_name}"; then
            download_failed=true
        fi
    fi
    
    if [ "$download_failed" = true ]; then
        error "Failed to download ${archive_name}. Please check if the release exists: https://github.com/${REPO}/releases/tag/${version}"
    fi
    
    # Verify downloaded file is not empty
    if [ ! -s "${tmp_dir}/${archive_name}" ]; then
        error "Downloaded file is empty. The release asset may not exist."
    fi
    
    # Verify file type
    if [ "$archive_ext" = "zip" ]; then
        if ! file "${tmp_dir}/${archive_name}" | grep -q "Zip archive"; then
            error "Downloaded file is not a valid zip archive. Got: $(file ${tmp_dir}/${archive_name})"
        fi
    else
        if ! file "${tmp_dir}/${archive_name}" | grep -qE "gzip compressed|tar archive"; then
            error "Downloaded file is not a valid tar.gz archive. Got: $(file ${tmp_dir}/${archive_name})"
        fi
    fi
    
    # Extract archive
    info "Extracting archive..."
    cd "${tmp_dir}"
    if [ "$archive_ext" = "zip" ]; then
        unzip -q "${archive_name}" || error "Failed to extract archive"
    else
        tar -xzf "${archive_name}" || error "Failed to extract archive"
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
    if has_command "${BINARY_NAME}"; then
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
    
    # Detect platform first
    local platform
    platform=$(detect_platform)
    local os="${platform%%_*}"
    info "Detected platform: ${platform}"
    
    # Check and install dependencies
    local installed_deps
    installed_deps=$(check_dependencies "$os")
    
    # Get latest version
    info "Fetching latest version..."
    local version
    version=$(get_latest_version)
    
    if [ -z "$version" ]; then
        error "Failed to get latest version"
    fi
    
    info "Latest version: ${version}"
    
    # Check if already installed
    if has_command "${BINARY_NAME}"; then
        local current_version
        current_version=$("${BINARY_NAME}" --version 2>&1 | head -n1 || echo "unknown")
        warn "Found existing installation: ${current_version}"
        
        if ! ask "Do you want to reinstall/upgrade?" "Y"; then
            info "Installation cancelled."
            exit 0
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
    
    # Cleanup dependencies if installed
    if [ -n "$installed_deps" ]; then
        echo ""
        warn "Dependencies installed: ${installed_deps}"
        if ask "Do you want to remove these dependencies?"; then
            info "Removing dependencies..."
            
            if has_command "pacman"; then
                for dep in $installed_deps; do
                    sudo pacman -Rns --noconfirm "$dep"
                done
            elif has_command "apt-get"; then
                for dep in $installed_deps; do
                    sudo apt-get remove -y "$dep"
                done
                sudo apt-get autoremove -y
            elif has_command "dnf"; then
                for dep in $installed_deps; do
                    sudo dnf remove -y "$dep"
                done
            elif has_command "yum"; then
                for dep in $installed_deps; do
                    sudo yum remove -y "$dep"
                done
            elif has_command "zypper"; then
                for dep in $installed_deps; do
                    sudo zypper remove -y "$dep"
                done
            elif has_command "apk"; then
                for dep in $installed_deps; do
                    sudo apk del "$dep"
                done
            elif has_command "brew"; then
                for dep in $installed_deps; do
                    brew uninstall "$dep"
                done
            fi
            
            success "Dependencies removed."
        else
            info "Dependencies kept for future use."
        fi
    fi
}

# Run main function
main "$@"

# Installation script for gohome (Windows PowerShell)
# Usage: irm https://raw.githubusercontent.com/anIcedAntFA/gohome/main/install.ps1 | iex

$ErrorActionPreference = 'Stop'

# Configuration
$REPO = "anIcedAntFA/gohome"
$BINARY_NAME = "gohome"
$INSTALL_DIR = if ($env:GOHOME_INSTALL_DIR) { $env:GOHOME_INSTALL_DIR } else { "$env:LOCALAPPDATA\Programs\gohome" }

# Colors for output
function Write-Info { Write-Host "[INFO] $args" -ForegroundColor Blue }
function Write-Success { Write-Host "[SUCCESS] $args" -ForegroundColor Green }
function Write-Warn { Write-Host "[WARN] $args" -ForegroundColor Yellow }
function Write-Error { Write-Host "[ERROR] $args" -ForegroundColor Red; exit 1 }

# Detect platform
function Get-Platform {
    $arch = switch ($env:PROCESSOR_ARCHITECTURE) {
        "AMD64" { "x86_64" }
        "ARM64" { "arm64" }
        "x86" { "i386" }
        default { Write-Error "Unsupported architecture: $env:PROCESSOR_ARCHITECTURE" }
    }
    return "windows_$arch"
}

# Get latest version from GitHub
function Get-LatestVersion {
    try {
        $response = Invoke-RestMethod -Uri "https://api.github.com/repos/$REPO/releases/latest" -ErrorAction SilentlyContinue
        return $response.tag_name
    }
    catch {
        # Fallback to scraping releases page
        try {
            $html = Invoke-WebRequest -Uri "https://github.com/$REPO/releases/latest" -UseBasicParsing
            if ($html.Content -match '/releases/tag/(v[\d\.]+)') {
                return $matches[1]
            }
        }
        catch {
            Write-Error "Failed to get latest version"
        }
    }
}

# Check if command exists
function Test-Command {
    param($Command)
    return [bool](Get-Command $Command -ErrorAction SilentlyContinue)
}

# Ask user for confirmation
function Confirm-Action {
    param(
        [string]$Prompt,
        [bool]$Default = $false
    )
    
    $defaultText = if ($Default) { "(Y/n)" } else { "(y/N)" }
    $response = Read-Host "$Prompt $defaultText"
    
    if ([string]::IsNullOrWhiteSpace($response)) {
        return $Default
    }
    
    return $response -match '^[Yy]'
}

# Download and install binary
function Install-Binary {
    param(
        [string]$Version,
        [string]$Platform
    )
    
    $archiveName = "${BINARY_NAME}_$($Version.TrimStart('v'))_$Platform.zip"
    $downloadUrl = "https://github.com/$REPO/releases/download/$Version/$archiveName"
    
    Write-Info "Installing $BINARY_NAME $Version for $Platform..."
    Write-Info "Download URL: $downloadUrl"
    
    # Create temporary directory
    $tmpDir = New-Item -ItemType Directory -Path (Join-Path $env:TEMP ([System.IO.Path]::GetRandomFileName()))
    
    try {
        # Download archive
        Write-Info "Downloading $archiveName..."
        $archivePath = Join-Path $tmpDir $archiveName
        
        try {
            Invoke-WebRequest -Uri $downloadUrl -OutFile $archivePath -ErrorAction Stop
        }
        catch {
            Write-Error "Failed to download $archiveName. Please check if the release exists: https://github.com/$REPO/releases/tag/$Version"
        }
        
        # Verify downloaded file
        if ((Get-Item $archivePath).Length -eq 0) {
            Write-Error "Downloaded file is empty. The release asset may not exist."
        }
        
        # Extract archive
        Write-Info "Extracting archive..."
        Expand-Archive -Path $archivePath -DestinationPath $tmpDir -Force
        
        # Find binary
        $binaryPath = Join-Path $tmpDir "$BINARY_NAME.exe"
        if (-not (Test-Path $binaryPath)) {
            Write-Error "Binary not found in archive: $binaryPath"
        }
        
        # Clean up conflicting binaries in GOPATH
        if ($env:GOPATH) {
            $goBin = Join-Path $env:GOPATH "bin\$BINARY_NAME.exe"
            if (Test-Path $goBin) {
                Write-Warn "Found existing binary in GOPATH/bin, removing..."
                Remove-Item $goBin -Force -ErrorAction SilentlyContinue
            }
        }
        
        # Also check common go/bin location
        $homeGoBin = Join-Path $env:USERPROFILE "go\bin\$BINARY_NAME.exe"
        if (Test-Path $homeGoBin) {
            Write-Warn "Found existing binary in ~/go/bin, removing..."
            Remove-Item $homeGoBin -Force -ErrorAction SilentlyContinue
        }
        
        # Create install directory if it doesn't exist
        if (-not (Test-Path $INSTALL_DIR)) {
            New-Item -ItemType Directory -Path $INSTALL_DIR -Force | Out-Null
        }
        
        # Install binary
        Write-Info "Installing to $INSTALL_DIR..."
        $targetPath = Join-Path $INSTALL_DIR "$BINARY_NAME.exe"
        Copy-Item $binaryPath $targetPath -Force
        
        Write-Success "$BINARY_NAME $Version installed successfully!"
        
        # Add to PATH if not already there
        $userPath = [Environment]::GetEnvironmentVariable("Path", "User")
        if ($userPath -notlike "*$INSTALL_DIR*") {
            Write-Warn "Adding $INSTALL_DIR to PATH..."
            [Environment]::SetEnvironmentVariable(
                "Path",
                "$userPath;$INSTALL_DIR",
                "User"
            )
            $env:Path = "$env:Path;$INSTALL_DIR"
            Write-Success "Added to PATH. Restart your terminal to use gohome."
        }
    }
    finally {
        # Cleanup
        Remove-Item $tmpDir -Recurse -Force -ErrorAction SilentlyContinue
    }
}

# Verify installation
function Test-Installation {
    $installedPath = Join-Path $INSTALL_DIR "$BINARY_NAME.exe"
    
    if (Test-Path $installedPath) {
        try {
            $version = & $installedPath --version 2>&1
            Write-Success "Verification successful!"
            Write-Info "Installed version: $version"
            Write-Info "Location: $installedPath"
        }
        catch {
            Write-Warn "Binary installed but could not execute --version"
        }
    }
    else {
        Write-Error "Binary not found at $installedPath"
    }
}

# Main installation flow
function Main {
    Write-Host ""
    Write-Info "═══════════════════════════════════════════════"
    Write-Info "  $BINARY_NAME Installation Script (Windows)"
    Write-Info "═══════════════════════════════════════════════"
    Write-Host ""
    
    # Detect platform
    $platform = Get-Platform
    Write-Info "Detected platform: $platform"
    
    # Get latest version
    Write-Info "Fetching latest version..."
    $version = Get-LatestVersion
    
    if (-not $version) {
        Write-Error "Failed to get latest version"
    }
    
    Write-Info "Latest version: $version"
    
    # Check if already installed
    $installedBinary = Join-Path $INSTALL_DIR "$BINARY_NAME.exe"
    if (Test-Path $installedBinary) {
        try {
            $currentVersion = & $installedBinary --version 2>&1
            Write-Warn "Found existing installation: $currentVersion"
            
            if (-not (Confirm-Action "Do you want to reinstall/upgrade?" $true)) {
                Write-Info "Installation cancelled."
                exit 0
            }
        }
        catch {
            Write-Warn "Found existing installation but could not determine version"
        }
    }
    
    # Install binary
    Install-Binary -Version $version -Platform $platform
    
    # Verify installation
    Write-Host ""
    Test-Installation
    
    Write-Host ""
    Write-Info "═══════════════════════════════════════════════"
    Write-Success "Installation complete!"
    Write-Info "═══════════════════════════════════════════════"
    Write-Host ""
    Write-Info "Get started:"
    Write-Info "  gohome --help"
    Write-Info "  gohome --version"
    Write-Host ""
    Write-Info "Documentation:"
    Write-Info "  https://github.com/$REPO"
    Write-Host ""
    Write-Warn "Note: If 'gohome' command is not found, restart your terminal or run:"
    Write-Warn "  `$env:Path = [System.Environment]::GetEnvironmentVariable('Path','User')"
}

# Run main function
Main

#!/usr/bin/env node

const { spawn } = require('child_process');
const path = require('path');
const os = require('os');
const fs = require('fs');

/**
 * Wrapper script for gohome binary
 * This script locates the platform-specific binary and executes it with all arguments
 */

// Determine the platform-specific binary name and path
function getBinaryPath() {
  const platform = os.platform();
  const arch = os.arch();

  let binaryName = 'gohome';
  if (platform === 'win32') {
    binaryName = 'gohome.exe';
  }

  // Binary is stored in the package's bin directory
  const binaryPath = path.join(__dirname, binaryName);

  if (!fs.existsSync(binaryPath)) {
    console.error(`Error: gohome binary not found at ${binaryPath}`);
    console.error('This indicates that the postinstall script failed to download the binary.');
    console.error(`Platform: ${platform}, Architecture: ${arch}`);
    console.error('\nTroubleshooting:');
    console.error('1. Check your internet connection');
    console.error('2. Verify the release exists: https://github.com/anIcedAntFA/gohome/releases');
    console.error('3. Try manual installation: npm install -g @ngockhoi96/gohome --verbose');
    process.exit(1);
  }

  return binaryPath;
}

// Execute the native binary with all arguments passed through
function main() {
  const binaryPath = getBinaryPath();

  // Spawn the native gohome binary with all command-line arguments
  const child = spawn(binaryPath, process.argv.slice(2), {
    stdio: 'inherit',
    env: process.env
  });

  child.on('error', (err) => {
    console.error(`Error executing gohome binary: ${err.message}`);
    process.exit(1);
  });

  child.on('exit', (code, signal) => {
    if (signal) {
      process.kill(process.pid, signal);
    } else {
      process.exit(code || 0);
    }
  });
}

main();

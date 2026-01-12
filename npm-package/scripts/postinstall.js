#!/usr/bin/env node

const fs = require('fs');
const path = require('path');
const https = require('https');
const { execSync } = require('child_process');

const packageJSON = require('../package.json');
const version = packageJSON.version;

// Strip prerelease suffixes (e.g., 1.0.3-hotfix.1 -> 1.0.3) for GitHub release tag
const releaseVersion = version.split('-')[0];

// Platform and architecture mapping
const PLATFORM_MAPPING = {
  darwin: 'darwin',
  linux: 'linux',
  win32: 'windows'
};

const ARCH_MAPPING = {
  x64: 'x86_64',
  arm64: 'arm64'
};

const platform = PLATFORM_MAPPING[process.platform];
const arch = ARCH_MAPPING[process.arch];

if (!platform || !arch) {
  console.error(`Unsupported platform: ${process.platform} ${process.arch}`);
  console.error('gohome supports:');
  console.error('  - macOS (darwin): x64, arm64');
  console.error('  - Linux: x64, arm64');
  console.error('  - Windows: x64, arm64');
  process.exit(1);
}

// Construct download URL
const ext = 'tar.gz'; // GoReleaser uses tar.gz for all platforms
const filename = `gohome_${releaseVersion}_${platform}_${arch}.${ext}`;
const downloadUrl = `https://github.com/anIcedAntFA/gohome/releases/download/v${releaseVersion}/${filename}`;

console.log(`📦 Installing gohome v${version} for ${platform}/${arch}...`);
console.log(`🔗 Downloading from: ${downloadUrl}`);

// Create bin directory (in parent directory of scripts/)
const binDir = path.join(__dirname, '..', 'bin');
if (!fs.existsSync(binDir)) {
  fs.mkdirSync(binDir, { recursive: true });
}

const archivePath = path.join(__dirname, filename);

// Download function
function download(url, dest) {
  return new Promise((resolve, reject) => {
    const file = fs.createWriteStream(dest);
    https.get(url, (response) => {
      // Handle redirects
      if (response.statusCode === 302 || response.statusCode === 301) {
        return download(response.headers.location, dest)
          .then(resolve)
          .catch(reject);
      }

      if (response.statusCode !== 200) {
        reject(new Error(`Download failed with status ${response.statusCode}`));
        return;
      }

      response.pipe(file);
      file.on('finish', () => {
        file.close();
        resolve();
      });
    }).on('error', (err) => {
      fs.unlinkSync(dest);
      reject(err);
    });
  });
}

// Extract function
function extract(archivePath, binDir, platform) {
  try {
    // All platforms use tar.gz now
    execSync(`tar -xzf "${archivePath}" -C "${binDir}"`, {
      stdio: 'inherit'
    });
    console.log('✅ Extraction completed');
  } catch (error) {
    console.error('❌ Extraction failed:', error.message);
    throw error;
  }
}

// Main installation flow
async function install() {
  try {
    // Download archive
    await download(downloadUrl, archivePath);
    console.log('✅ Download completed');

    // Extract archive
    extract(archivePath, binDir, platform);

    // Clean up archive
    fs.unlinkSync(archivePath);

    // Make binary executable on Unix
    if (platform !== 'windows') {
      const binaryPath = path.join(binDir, 'gohome');
      fs.chmodSync(binaryPath, 0o755);
    }

    console.log('✅ gohome installed successfully!');
    console.log(`📍 Binary location: ${path.join(binDir, platform === 'windows' ? 'gohome.exe' : 'gohome')}`);
    console.log('\n🚀 Get started:');
    console.log('   gohome --help');
  } catch (error) {
    console.error('❌ Installation failed:', error.message);
    console.error('\n💡 Alternative installation methods:');
    console.error('   - Download binary: https://github.com/anIcedAntFA/gohome/releases');
    console.error('   - Go install: go install github.com/anIcedAntFA/gohome/cmd/gohome@latest');
    process.exit(1);
  }
}

install();

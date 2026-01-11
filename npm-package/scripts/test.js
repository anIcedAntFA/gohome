#!/usr/bin/env node

/**
 * Test suite for gohome npm package
 * Verifies that the binary was downloaded and is executable
 */

const { execSync } = require('child_process');
const path = require('path');
const fs = require('fs');
const os = require('os');

console.log('🧪 Running gohome npm package tests...\n');

let passed = 0;
let failed = 0;

function test(name, fn) {
  try {
    fn();
    console.log(`✅ ${name}`);
    passed++;
  } catch (err) {
    console.error(`❌ ${name}`);
    console.error(`   ${err.message}\n`);
    failed++;
  }
}

// Test 1: Check binary exists
test('Binary exists in bin directory', () => {
  const binaryName = os.platform() === 'win32' ? 'gohome.exe' : 'gohome';
  const binaryPath = path.join(__dirname, '..', 'bin', binaryName);
  if (!fs.existsSync(binaryPath)) {
    throw new Error(`Binary not found at ${binaryPath}`);
  }
});

// Test 2: Binary is executable (version check)
test('Binary executes and returns version', () => {
  const binaryName = os.platform() === 'win32' ? 'gohome.exe' : 'gohome';
  const binaryPath = path.join(__dirname, '..', 'bin', binaryName);
  
  try {
    const output = execSync(`"${binaryPath}" --version`, { encoding: 'utf8', timeout: 5000 });
    if (!output.includes('gohome')) {
      throw new Error(`Unexpected version output: ${output}`);
    }
  } catch (err) {
    // Try without version flag (some CLIs don't support it)
    execSync(`"${binaryPath}" --help`, { encoding: 'utf8', timeout: 5000 });
  }
});

// Test 3: Wrapper script exists
test('Wrapper script (gohome.js) exists', () => {
  const wrapperPath = path.join(__dirname, '..', 'bin', 'gohome.js');
  if (!fs.existsSync(wrapperPath)) {
    throw new Error(`Wrapper not found at ${wrapperPath}`);
  }
});

// Test 4: package.json has correct bin configuration
test('package.json has correct bin configuration', () => {
  const packagePath = path.join(__dirname, '..', 'package.json');
  const packageData = JSON.parse(fs.readFileSync(packagePath, 'utf8'));
  
  if (!packageData.bin || !packageData.bin.gohome) {
    throw new Error('package.json missing bin.gohome field');
  }
  
  if (packageData.bin.gohome !== 'bin/gohome.js') {
    throw new Error(`Expected bin.gohome to be 'bin/gohome.js', got '${packageData.bin.gohome}'`);
  }
});

// Test 5: postinstall script is configured
test('postinstall script is configured', () => {
  const packagePath = path.join(__dirname, '..', 'package.json');
  const packageData = JSON.parse(fs.readFileSync(packagePath, 'utf8'));
  
  if (!packageData.scripts || !packageData.scripts.postinstall) {
    throw new Error('package.json missing scripts.postinstall');
  }
  
  if (!packageData.scripts.postinstall.includes('postinstall.js')) {
    throw new Error('postinstall script should reference postinstall.js');
  }
});

// Summary
console.log('\n' + '='.repeat(50));
console.log(`📊 Test Results: ${passed} passed, ${failed} failed`);
console.log('='.repeat(50) + '\n');

process.exit(failed > 0 ? 1 : 0);

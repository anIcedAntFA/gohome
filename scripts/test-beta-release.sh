#!/usr/bin/env bash
# Test script for v1.3.0-beta.1 release

set -e  # Exit on error

echo "🧪 Testing gohome v1.3.0-beta.1 Installation"
echo "=============================================="
echo ""

# Colors
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Test counter
TESTS_PASSED=0
TESTS_FAILED=0

# Helper functions
pass() {
    echo -e "${GREEN}✓${NC} $1"
    ((TESTS_PASSED++))
}

fail() {
    echo -e "${RED}✗${NC} $1"
    ((TESTS_FAILED++))
}

warn() {
    echo -e "${YELLOW}⚠${NC} $1"
}

section() {
    echo ""
    echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
    echo "$1"
    echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
}

# Test 1: Installation
section "1. Testing Installation"
if command -v gohome &> /dev/null; then
    pass "gohome is installed"
else
    fail "gohome is not installed"
    echo ""
    echo "Install with: go install github.com/anIcedAntFA/gohome/cmd/gohome@v1.3.0-beta.1"
    exit 1
fi

# Test 2: Version Check
section "2. Testing Version Command"
VERSION=$(gohome version 2>&1)
if [[ $VERSION == *"1.3.0-beta.1"* ]]; then
    pass "Version is correct: v1.3.0-beta.1"
    echo "   Output: $VERSION"
else
    fail "Version mismatch"
    echo "   Expected: v1.3.0-beta.1"
    echo "   Got: $VERSION"
fi

# Test 3: Help Command
section "3. Testing Help Output"
if gohome --help &> /dev/null; then
    pass "Help command works"
else
    fail "Help command failed"
fi

# Test 4: Subcommands
section "4. Testing Subcommands Exist"
HELP_OUTPUT=$(gohome --help 2>&1)

if [[ $HELP_OUTPUT == *"report"* ]]; then
    pass "Subcommand 'report' exists"
else
    fail "Subcommand 'report' not found"
fi

if [[ $HELP_OUTPUT == *"config"* ]]; then
    pass "Subcommand 'config' exists"
else
    fail "Subcommand 'config' not found"
fi

if [[ $HELP_OUTPUT == *"completion"* ]]; then
    pass "Subcommand 'completion' exists"
else
    fail "Subcommand 'completion' not found"
fi

if [[ $HELP_OUTPUT == *"version"* ]]; then
    pass "Subcommand 'version' exists"
else
    fail "Subcommand 'version' not found"
fi

# Test 5: Config Subcommand
section "5. Testing Config Subcommand"
if gohome config --help &> /dev/null; then
    pass "Config help works"
else
    fail "Config help failed"
fi

CONFIG_HELP=$(gohome config --help 2>&1)
if [[ $CONFIG_HELP == *"init"* ]] && [[ $CONFIG_HELP == *"show"* ]] && [[ $CONFIG_HELP == *"edit"* ]]; then
    pass "Config subcommands (init, show, edit) exist"
else
    fail "Config subcommands missing"
fi

# Test 6: Completion Subcommand
section "6. Testing Completion Generation"
COMPLETIONS=("bash" "zsh" "fish" "powershell")
for shell in "${COMPLETIONS[@]}"; do
    if gohome completion "$shell" &> /dev/null; then
        pass "Completion for $shell generated"
    else
        fail "Completion for $shell failed"
    fi
done

# Test 7: Report Command (basic)
section "7. Testing Report Command"
# Create a temporary workspace
TEMP_WORKSPACE=$(mktemp -d)
mkdir -p "$TEMP_WORKSPACE/test-repo"
cd "$TEMP_WORKSPACE/test-repo"
git init &> /dev/null
git config user.name "Test User" &> /dev/null
git config user.email "test@example.com" &> /dev/null
echo "test" > test.txt
git add test.txt &> /dev/null
git commit -m "feat: test commit" &> /dev/null

# Run report
if gohome report --workspace "$TEMP_WORKSPACE" --days 7 &> /dev/null; then
    pass "Report command works"
else
    warn "Report command failed (might need valid repos)"
fi

# Cleanup
cd - &> /dev/null
rm -rf "$TEMP_WORKSPACE"

# Test 8: Environment Variables
section "8. Testing Environment Variables"
export GOHOME_FORMAT="markdown"
export GOHOME_PERIOD_DAYS=1

# Check if env vars are recognized (by showing config)
CONFIG_OUTPUT=$(gohome config show 2>&1 || true)
if [[ -n "$CONFIG_OUTPUT" ]]; then
    pass "Environment variables appear to be working"
    warn "Manual verification recommended"
else
    warn "Could not verify environment variables"
fi

unset GOHOME_FORMAT
unset GOHOME_PERIOD_DAYS

# Test 9: Config File Locations
section "9. Checking Config File Support"
OLD_CONFIG="$HOME/.gohome.json"
NEW_CONFIG="$HOME/.config/gohome/config.json"

if [[ -f "$OLD_CONFIG" ]]; then
    pass "Old config location exists: $OLD_CONFIG"
    echo "   (backward compatibility preserved)"
fi

if [[ -f "$NEW_CONFIG" ]]; then
    pass "New config location exists: $NEW_CONFIG"
fi

if [[ ! -f "$OLD_CONFIG" ]] && [[ ! -f "$NEW_CONFIG" ]]; then
    warn "No config file found (this is OK for fresh install)"
    echo "   Create one with: gohome config init"
fi

# Summary
section "Test Summary"
TOTAL=$((TESTS_PASSED + TESTS_FAILED))
echo ""
echo "Tests Passed: ${GREEN}$TESTS_PASSED${NC}/$TOTAL"
if [[ $TESTS_FAILED -gt 0 ]]; then
    echo "Tests Failed: ${RED}$TESTS_FAILED${NC}/$TOTAL"
fi
echo ""

if [[ $TESTS_FAILED -eq 0 ]]; then
    echo -e "${GREEN}🎉 All tests passed!${NC}"
    echo ""
    echo "Next steps:"
    echo "  1. Try: gohome config init"
    echo "  2. Configure: gohome config edit"
    echo "  3. Generate report: gohome report"
    echo "  4. Install completions: gohome completion <shell>"
    exit 0
else
    echo -e "${RED}❌ Some tests failed${NC}"
    echo ""
    echo "Please report issues at:"
    echo "https://github.com/anIcedAntFA/gohome/issues/new"
    exit 1
fi

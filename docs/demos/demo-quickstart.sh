#!/usr/bin/env bash
# Demo 1: Quick Start
# Usage: asciinema rec -c "./docs/demos/demo-quickstart.sh" docs/demos/quickstart.cast

set -e

# Typing simulation function
type_command() {
    local cmd="$1"
    local delay="${2:-0.05}"
    
    for ((i=0; i<${#cmd}; i++)); do
        echo -n "${cmd:$i:1}"
        sleep "$delay"
    done
    echo
}

# Clear and show banner
clear
echo "🚀 gohome - Git Activity Aggregator"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo
sleep 1

# 1. Show version
echo "$ gohome --version"
sleep 0.5
gohome --version
sleep 2

# 2. Show help (brief)
echo
echo "$ gohome --help"
sleep 0.5
gohome --help | head -30
sleep 3

# 3. Run basic command
echo
echo "$ gohome"
sleep 0.5
gohome
sleep 3

# 4. Custom time range
echo
echo "$ gohome --days 7 --format table"
sleep 0.5
gohome --days 7 --format table
sleep 3

# End
echo
echo "✨ Try it yourself!"
echo "   curl -fsSL https://raw.githubusercontent.com/anIcedAntFA/gohome/main/scripts/install.sh | bash"
sleep 2

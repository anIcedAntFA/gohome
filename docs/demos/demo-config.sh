#!/usr/bin/env bash
# Demo 2: Advanced Features
# Usage: asciinema rec -c "./docs/demos/demo-config.sh" docs/demos/config.cast

set -e

clear
echo "⚙️  gohome - Configuration & Customization"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo
sleep 1

# 1. Custom tasks
echo "$ gohome -t \"Code review: PR #123\" -t \"Team standup meeting\""
sleep 0.5
gohome -t "Code review: PR #123" -t "Team standup meeting"
sleep 3

# 2. Different formats
echo
echo "$ gohome --format table --style markdown"
sleep 0.5
gohome --format table --style markdown
sleep 3

# 3. Nature style
echo
echo "$ gohome --format table --style nature"
sleep 0.5
gohome --format table --style nature
sleep 3

# 4. Save config
echo
echo "$ gohome --workspaces ~/projects --days 1 --format table --save"
sleep 0.5
gohome --workspaces ~/projects --days 1 --format table --save
sleep 2

echo
echo "✅ Configuration saved to ~/.gohome.json"
sleep 2

# 5. Copy to clipboard
echo
echo "$ gohome --copy"
sleep 0.5
gohome --copy >/dev/null 2>&1 || echo "Report generated!"
sleep 1
echo "📋 Report copied to clipboard!"
sleep 2

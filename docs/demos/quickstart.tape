# VHS Tape: Quick Start Demo
# This tape demonstrates basic gohome usage for new users

# Output settings
Output docs/demos/quickstart.gif

# Terminal settings
Set Shell fish
Set FontSize 20
Set Width 1200
Set Height 800
Set Padding 20
Set Margin 20
Set MarginFill "#282a36"
Set BorderRadius 12
Set WindowBar Colorful
Set TypingSpeed 50ms
Set Theme "Dracula"

# Demo starts
Type "# 🚀 gohome - Git Activity Aggregator"
Enter
Sleep 500ms

Type "# Let's see what we built today!"
Enter
Sleep 1s

# Show version
Type "gohome --version"
Sleep 500ms
Enter
Sleep 2s

# Show help
Type "gohome --help"
Sleep 500ms
Enter
Sleep 3s

# Basic usage - today's commits
Type "# Get today's activity report"
Enter
Type "gohome"
Sleep 500ms
Enter
Sleep 3s

# Custom time range with table format
Type "# Last 5 days in table format"
Enter
Type "gohome --days 5 --format table"
Sleep 500ms
Enter
Sleep 4s

# Save configuration
Type "# Save your preferences"
Enter
Type "gohome --today --icon --scope --save"
Sleep 500ms
Enter
Sleep 3s

# Copy to clipboard
Type "# Copy report to clipboard"
Enter
Type "gohome --copy"
Sleep 500ms
Enter
Sleep 1s
Type "# 📋 Report copied!"
Enter
Sleep 2s

# Final message
Sleep 1s
Type "# ✨ Try it yourself!"
Enter
Type "# curl -sS https://get.ngockhoi96.dev/gohome | sh"
Enter
Sleep 3s

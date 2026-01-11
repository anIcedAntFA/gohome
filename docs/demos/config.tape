# VHS Tape: Configuration & Customization
# This tape demonstrates advanced gohome features

# Output settings
Output docs/demos/config.gif

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

# Setup clean shell (hidden)
Hide
Type "export PS1='$ '"
Enter
Type "clear"
Enter
Show

# Demo starts
Type "# ⚙️  gohome - Configuration & Customization"
Enter
Sleep 500ms

# Custom tasks
Type "# Add custom tasks alongside git commits"
Enter
Type `gohome -t "Code review: PR #123" -t "Team standup meeting"`
Sleep 500ms
Enter
Sleep 3s

# Different output formats
Type "# Try markdown table style"
Enter
Type "gohome --format table --style markdown"
Sleep 500ms
Enter
Sleep 3s

# With icons and scope
Type "# Show emoji icons and commit scopes"
Enter
Type "gohome --icon --scope"
Sleep 500ms
Enter
Sleep 3s

# Save configuration
Type "# Save your preferences"
Enter
Type "gohome --days 5 --format table --save"
Sleep 500ms
Enter
Sleep 2s

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

# Final tip
Type "# 💡 Config saved to ~/.gohome.json"
Enter
Sleep 2s

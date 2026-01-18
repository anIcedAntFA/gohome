# Test Strategy for Phase 4

## Current Coverage Analysis (Baseline: 14.6%)

### ✅ Well-Tested Packages (>80%)
- **internal/scanner**: 87.5% ✅ (comprehensive table-driven tests)
- **internal/config/viper**: 29.9% (validation tests only)
- **internal/spinner**: 31.6% (basic tests only)

### ❌ Zero Coverage Packages (Priority Targets)
1. **cmd/gohome/cmd** (0.0%) - **CRITICAL**
   - config.go: 4 commands (list, get, set, reset)
   - report.go: Main report generation logic
   - version.go: Version display
   - completion.go: Shell completion generation
   - root.go: Root command initialization

2. **internal/parser** (0.0%) - **HIGH PRIORITY**
   - Conventional Commits regex parsing
   - Emoji extraction
   - Critical for correctness

3. **internal/renderer** (0.0%) - **HIGH PRIORITY**  
   - Text/table output formatting
   - Task rendering
   - User-facing output

4. **internal/git** (0.0%) - **HIGH PRIORITY**
   - Git command execution
   - User detection
   - Input sanitization (security)

5. **internal/version** (0.0%) - **MEDIUM PRIORITY**
   - Version string formatting
   - Build info parsing

6. **internal/sys** (0.0%) - **LOW PRIORITY**
   - Clipboard operations (platform-specific)
   - Difficult to test in CI

## Testing Strategy

### Phase 1: Command Tests (Tasks 2-6)
**Target: cmd/gohome/cmd coverage 60-70%**

Each command file needs:
- Command initialization tests
- Flag parsing tests
- Error handling tests
- Happy path integration tests

**Approach:**
- Mock Viper config for deterministic tests
- Use test buffers for output capture
- Test flag validation logic
- Test subcommand routing

**Priority:**
1. config.go (list, get, set, reset subcommands)
2. report.go (main business logic)
3. version.go (simple, quick win)
4. completion.go (shell completion generation)
5. root.go (initialization)

### Phase 2: Parser Tests (Task 7)
**Target: internal/parser coverage >90%**

**Test Cases:**
- Valid conventional commits (feat, fix, chore, etc.)
- Commits with scopes: `feat(api): add endpoint`
- Breaking changes: `feat!: breaking change`
- Emoji extraction: `✨ feat: new feature`
- Invalid/malformed commits
- Edge cases: empty messages, special characters

**Approach:**
- Table-driven tests with comprehensive commit examples
- Test emoji regex patterns
- Benchmark parser performance

### Phase 3: Git Client Tests (Task 8)
**Target: internal/git coverage >80%**

**Test Cases:**
- GetUser() with various git configs
- GetLogs() command construction
- Input sanitization (security tests)
- Branch filtering logic
- Date range calculations
- Error handling (git not installed, invalid repo)

**Approach:**
- Use test git repositories
- Mock os/exec for command testing
- Test injection prevention

### Phase 4: Renderer Tests (Task 9)
**Target: internal/renderer coverage >80%**

**Test Cases:**
- Text format output
- Table format with different styles (normal, markdown)
- Task rendering (text and table)
- Icon and scope toggling
- Empty commit lists
- Special characters in messages

**Approach:**
- Capture output to buffers
- Golden file tests for table layouts
- Test tablewriter integration

### Phase 5: Integration Tests (Task 10)
**Target: End-to-end workflow coverage**

**Test Scenarios:**
1. **Full Report Generation:**
   - Create test git repo with commits
   - Run report command
   - Verify output format
   - Test with different time periods

2. **Config Management:**
   - Set/get/list/reset operations
   - Multi-format config loading (JSON/YAML/TOML)
   - Flag precedence testing

3. **Error Scenarios:**
   - No git repos found
   - Invalid config values
   - Missing git author
   - Invalid time periods

**Approach:**
- Use `testdata/` directory with test repos
- Create temporary config files
- Test complete workflows
- Verify exit codes

### Phase 6: Coverage Verification (Task 11)
**Target: Overall >80% coverage**

**Coverage Goals by Package:**
- cmd/gohome/cmd: 65-70%
- internal/parser: 90%+
- internal/git: 80%+
- internal/renderer: 80%+
- internal/config/viper: 50%+ (improve existing)
- internal/scanner: 90%+ (maintain)
- internal/version: 60%+

**Exclusions (acceptable 0% coverage):**
- internal/sys/clipboard.go (platform-specific, hard to test)
- internal/spinner/example_usage.go (example code)
- cmd/gohome/main.go (entry point, minimal logic)

## Test File Organization

```
cmd/gohome/cmd/
  config_test.go       # Config command tests
  report_test.go       # Report command tests
  version_test.go      # Version command tests
  completion_test.go   # Completion command tests
  root_test.go         # Root command tests

internal/
  parser/
    parser_test.go     # Conventional Commits parsing
  git/
    client_test.go     # Git operations
  renderer/
    printer_test.go    # Output formatting
    
testdata/
  test-repo-1/.git/    # Test git repository
  test-repo-2/.git/    # Multi-branch test repo
  configs/
    test.json          # Test config files
    test.yaml
    test.toml
```

## Testing Tools & Patterns

### Required Packages
```go
import (
    "testing"
    "github.com/stretchr/testify/assert"  // Assertions
    "github.com/stretchr/testify/require" // Fatal assertions
    "bytes"                               // Output capture
    "os"
    "path/filepath"
)
```

### Test Patterns
1. **Table-Driven Tests** (for multiple similar cases)
2. **Test Fixtures** (testdata directory)
3. **Output Capture** (bytes.Buffer for stdout/stderr)
4. **Temp Directories** (t.TempDir() for file operations)
5. **Subtest Groups** (t.Run() for organization)

### Coverage Commands
```bash
# Generate coverage
go test -coverprofile=coverage.out ./...

# View by function
go tool cover -func=coverage.out

# HTML report
go tool cover -html=coverage.out -o coverage.html

# Package-level summary
go test -cover ./...
```

## Timeline Estimate
- **Task 2-6** (CMD tests): 3-4 hours (5 files)
- **Task 7** (Parser tests): 1-2 hours
- **Task 8** (Git tests): 2 hours
- **Task 9** (Renderer tests): 2 hours
- **Task 10** (Integration): 2-3 hours
- **Task 11** (Coverage verification): 1 hour

**Total: 11-14 hours of focused work**

## Success Criteria
- [ ] All packages except sys/clipboard have tests
- [ ] Overall coverage >80%
- [ ] No failing tests in CI
- [ ] All critical paths tested (report generation, config management)
- [ ] Security-sensitive code tested (input sanitization)
- [ ] Documentation updated with testing guidelines

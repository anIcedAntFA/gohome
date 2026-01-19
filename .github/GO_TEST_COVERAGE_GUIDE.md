# Go Test Coverage - Quick Guide

## ✅ Cách Test Đúng Trong Go

### Test Entire Package (ĐÚNG)
```bash
# Basic test
go test ./cmd/gohome/cmd/

# With coverage
go test -cover ./cmd/gohome/cmd/

# With verbose output
go test -v -cover ./cmd/gohome/cmd/

# Generate coverage profile
go test -coverprofile=coverage.txt ./cmd/gohome/cmd/

# View coverage in browser
go test -coverprofile=coverage.txt ./cmd/gohome/cmd/
go tool cover -html=coverage.txt
```

### Test Single File (SAI - Không Work!)
```bash
# ❌ SAI: Undefined references
go test -cover ./cmd/gohome/cmd/root.go
go test -cover ./cmd/gohome/cmd/completion.go
go test -cover ./cmd/gohome/cmd/root_test.go

# Lỗi: undefined: rootCmd
# Vì: File riêng lẻ thiếu dependencies từ package
```

### Tại Sao Không Test Được Single File?

**Go package model:**
- Tất cả `.go` files trong 1 folder = 1 package
- Các files share variables/functions với nhau
- Test cần load TOÀN BỘ package để có dependencies

**Example:**
```go
// root.go
package cmd
var rootCmd = &cobra.Command{...}

// completion.go  
package cmd
func init() {
    rootCmd.AddCommand(completionCmd)  // Needs rootCmd!
}

// version.go
package cmd
func init() {
    rootCmd.AddCommand(versionCmd)  // Needs rootCmd!
}
```

Nếu test `completion.go` riêng → `rootCmd` undefined!

## 📊 Coverage Analysis

### Current Status (v1.3.0-beta.2)

```
Total cmd package coverage: 20.5%

By file:
- completion.go: ~85% ✅ (after improvements)
- version.go: 100% ✅
- root.go: ~80% ✅ (initConfig covered, Execute 0%)
- config.go: ~40% ⚠️ (many 0% functions)
- report.go: ~15% ❌ (many 0% functions)
```

### What's NOT Covered (0% functions)

**config.go:**
- `runConfigList` - 0%
- `runConfigGet` - 0%
- `runConfigSet` - 0%
- `runConfigReset` - 0%

**report.go:**
- `runReport` - 0%
- `handlePeriodFlags` - 0%
- `scanRepositories` - 0%
- `processCommits` - 0%
- `collectActiveTasks` - 0%
- `handleClipboard` - 0%
- `handleSaveConfig` - 0%
- `getAuthorName` - 0%

**root.go:**
- `Execute()` - 0% (has os.Exit(), hard to test)

### Why Execute() Shows 0%

```go
func Execute() {
    if err := rootCmd.Execute(); err != nil {
        fmt.Fprintf(os.Stderr, "❌ Error: %v\n", err)
        os.Exit(1)  // ← This makes it hard to test!
    }
}
```

**Solutions:**
1. **Exclude from coverage** (recommended for main wrappers)
2. **Test via subprocess** (complex, usually not worth it)
3. **Refactor** to return error instead of Exit (breaking change)

## 🎯 Priority to Improve Coverage

### High Priority (Big Impact)
1. **report.go tests** → +20-25% coverage
   - Test `runReport()` function
   - Mock scanner, git client, parser
   - Test different flags combinations

2. **config.go tests** → +10-15% coverage
   - Test `runConfigList()`
   - Test `runConfigSet()/Get()/Reset()`
   - Mock viper config

### Medium Priority
3. **Integration tests** → +5-10% coverage
   - End-to-end report generation
   - Config file operations
   - Multiple repos scanning

### Low Priority (Acceptable to Skip)
4. **Execute() wrapper** → +0.5% coverage
   - Hard to test due to os.Exit()
   - Already tested via rootCmd.Execute()
   - Industry standard to exclude

## 📝 Testing Best Practices

### DO ✅
```bash
# Test entire package
go test ./pkg/...

# Test with race detector
go test -race ./...

# Test with coverage
go test -cover -coverprofile=coverage.txt ./...

# View detailed coverage
go tool cover -html=coverage.txt

# Test specific function
go test -v -run TestFunctionName ./pkg/...

# Test with timeout
go test -timeout 30s ./...
```

### DON'T ❌
```bash
# Don't test single .go file
go test ./pkg/file.go

# Don't test _test.go file directly  
go test ./pkg/file_test.go

# Don't use relative imports in tests
import "../other"  # Bad!
```

## 🔧 How to Add Tests

### 1. For config.go Functions

```go
// config_test.go
func TestRunConfigList(t *testing.T) {
    // Reset viper
    viper.Reset()
    defer viper.Reset()
    
    // Set test config
    viper.Set("days", 5)
    viper.Set("path", "/test/path")
    
    // Capture output
    var buf bytes.Buffer
    configCmd.SetOut(&buf)
    
    // Run command
    configCmd.SetArgs([]string{"list"})
    err := configCmd.Execute()
    
    // Assert
    if err != nil {
        t.Errorf("configCmd.Execute() failed: %v", err)
    }
    
    output := buf.String()
    if !strings.Contains(output, "days") {
        t.Error("Output should contain 'days'")
    }
}
```

### 2. For report.go Functions

```go
// report_test.go
func TestRunReport(t *testing.T) {
    // This requires mocking scanner, git, parser
    // More complex, needs integration test setup
    
    // Setup temp git repo
    tmpDir := t.TempDir()
    // ... create test git repo ...
    
    // Run report
    reportCmd.SetArgs([]string{
        "--path", tmpDir,
        "--days", "1",
    })
    
    err := reportCmd.Execute()
    if err != nil {
        t.Errorf("Failed: %v", err)
    }
}
```

## 🎓 Understanding Coverage Percentage

### What's Good Coverage?

| Type | Target | Priority |
|------|--------|----------|
| Critical logic (parser, git) | 80%+ | High |
| Business logic (config, report) | 60%+ | High |
| CLI commands | 40-60% | Medium |
| Integration/glue code | 30-40% | Low |
| Main entry points | Exclude | N/A |

### Current Project Status

```
✅ parser: 100% → Excellent!
✅ git: ~90% → Excellent!
✅ scanner: ~95% → Excellent!
⚠️  cmd/completion: ~85% → Good
⚠️  cmd/root: ~80% → Good (Execute excluded)
⚠️  cmd/version: 100% → Excellent!
❌ cmd/config: ~40% → Needs improvement
❌ cmd/report: ~15% → Needs significant work
```

**Overall: 45-50% coverage after current improvements**
**Target: 60%+ after adding config and report tests**

## 🚀 Quick Commands Reference

```bash
# Run all tests
go test ./...

# Coverage for all packages
go test -cover ./...

# Coverage with details
go test -coverprofile=coverage.txt ./...
go tool cover -func=coverage.txt

# Coverage HTML report
go test -coverprofile=coverage.txt ./...
go tool cover -html=coverage.txt -o coverage.html

# Test specific package
go test -v ./cmd/gohome/cmd/

# Test specific function
go test -v -run TestCompletionCommand ./cmd/gohome/cmd/

# Run with race detection
go test -race ./...

# Verbose output
go test -v ./...

# Short mode (skip long tests)
go test -short ./...

# Coverage by package
go test -cover ./cmd/... ./internal/...
```

## 📚 Resources

- [Go Testing Package](https://pkg.go.dev/testing)
- [Go Coverage](https://go.dev/blog/cover)
- [Table Driven Tests](https://go.dev/wiki/TableDrivenTests)
- [Advanced Testing](https://quii.gitbook.io/learn-go-with-tests/)

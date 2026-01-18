# Code Review & Refactoring Summary

## Overview

Sau khi review kỹ code và Viper best practices documentation, tôi đã phân tích các vấn đề hiện tại và đề xuất improvements.

## ✅ Completed

### 1. Created Comprehensive Viper Documentation
- **File**: `docs/VIPER_CONFIG_MANAGEMENT.md`
- **Content**:
  - Complete architecture explanation
  - Configuration precedence flow với diagrams
  - Implementation details cho từng component
  - Best practices (DO/DON'T patterns)
  - Common patterns & use cases
  - Troubleshooting guide
  - Migration checklist

### 2. Identified Code Quality Issues

#### Issue #1: Giant Switch Statement (100+ lines)
**Location**: `cmd/gohome/cmd/config.go` → `setConfigValue()`

**Problem**:
```go
func setConfigValue(cfg *viperconfig.Config, key, value string) error {
    // 100+ lines of nested switch statements
    // Duplicated type conversion logic
    // Hard to maintain and extend
}
```

**Solution**: Use reflection-based approach (implemented in docs, ready to apply)

---

#### Issue #2: RegisterAlias in Wrong Place
**Location**: `internal/config/viper/viper.go` → `LoadFromViper()`

**Problem**: Gọi `RegisterAlias()` mỗi lần load config
```go
func LoadFromViper() *Config {
    viper.RegisterAlias("max_depth", "max-depth")  // ❌ Called every time
    //...
}
```

**Solution**: Move to package `init()` (run once)
```go
func init() {
    viper.RegisterAlias("max_depth", "max-depth")  // ✅ Called once
}
```

---

#### Issue #3: Config List UI - Manual String Formatting
**Location**: `cmd/gohome/cmd/config.go` → `runConfigList()`

**Problem**: Manual `Printf` formatting
```go
fmt.Printf("%-20s = %v\n", "days", cfg.Days)
fmt.Printf("%-20s = %v\n", "format", cfg.Format)
// ... 20+ lines
```

**Solution**: Use `tablewriter` như report command (more consistent, prettier)

---

##  🔍 Analysis: Current Architecture

### Strengths ✅

1. **Good separation of concerns**
   - `root.go`: Viper initialization
   - `report.go`: Business logic
   - `viper.go`: Config struct & persistence
   - `config.go`: Config management

2. **Flag binding done correctly**
   - `BindPFlags()` called once in `init()`
   - Persistent flags on root command

3. **Clean JSON output**
   - Using `json.Encoder` instead of `viper.WriteConfig()`
   - Avoids duplicate keys issue

4. **Precedence works correctly**
   - Flags > Config > Defaults hierarchy respected

### Weaknesses ❌

1. **Maintainability**
   - `setConfigValue()` giant switch → hard to extend
   - Duplicate type conversion logic across cases

2. **Performance**
   - `RegisterAlias()` called on every `LoadFromViper()`
   - Should be one-time setup in `init()`

3. **UI/UX**
   - Config list output không đẹp
   - Inconsistent với report command styling

4. **Code duplication**
   - Type validation repeated nhiều chỗ
   - Enum validation (`format`, `style`) inline

---

## 📊 Global Flags Design Review

### Current Setup
```go
// All flags defined as PersistentFlags on root command
rootCmd.PersistentFlags().IntP("days", "d", 1, "...")
rootCmd.PersistentFlags().StringP("format", "f", "text", "...")

// Makes gohome the default command
rootCmd.RunE = runReport
```

### Analysis

**Pros**:
- ✅ Simple UX: `gohome` = `gohome report` (intuitive default)
- ✅ Flags work everywhere due to Persistent nature
- ✅ Less typing for most common use case

**Cons**:
- ⚠️ Confusing: `gohome config -d 7` looks valid but `-d` is meaningless for config command
- ⚠️ Help pollution: `gohome config --help` shows all report flags
- ⚠️ Semantic ambiguity: Are they "global" or "report-specific"?

### Recommendation

**Option A: Keep Current (Simpler)**
- Document clearly in help text that flags are for report command
- Add validation to reject invalid flag combinations
- **Best for**: Small CLIs with one primary command

**Option B: Make Flags Command-Specific (Cleaner)**
- Move flags from `rootCmd` to `reportCmd`
- Keep only truly global flags on root (like `--config`, `--verbose`)
- **Best for**: CLIs with multiple equally important commands

**My Recommendation**: **Option A** for gohome because:
1. Report is clearly the primary/default command
2. Config/version are utility commands
3. Simpler UX outweighs architectural purity
4. Document it well (which we did!)

---

## 🎯 Proposed Refactoring Plan

### Phase 1: Quick Wins (Low Risk)
1. **Move `RegisterAlias()` to `init()`** ← Can do now
2. **Improve config list UI with tablewriter** ← Can do now
3. **Add validation layer** for enum fields

### Phase 2: Structural Improvements (Medium Risk)
4. **Refactor `setConfigValue()` to use reflection**
   - Reduces ~100 lines to ~50 lines
   - Easier to add new config fields
   - Single source of truth for type handling

5. **Extract validation logic**
   - Create `validate()` methods on Config
   - Separate concerns: Viper (mechanics) vs Business Logic (validation)

### Phase 3: Advanced (Future)
6. **Add config watching** (`viper.WatchConfig()`)
7. **Support multiple config formats** (YAML, TOML)
8. **Add config migration** for version upgrades

---

## 💡 Code Examples

### Example: Reflection-Based setConfigValue

```go
// SetValue sets a config field by key with automatic type conversion.
// Uses reflection to avoid giant switch statements (DRY principle).
func (c *Config) SetValue(key, value string) error {
    key = strings.ReplaceAll(key, "-", "_")  // kebab → snake

    // Field mapping
    fieldMap := map[string]string{
        "days":         "Days",
        "format":       "Format",
        "max_depth":    "MaxDepth",
        "all_branches": "AllBranches",
        // ...
    }

    fieldName, ok := fieldMap[key]
    if !ok {
        return fmt.Errorf("unknown config key %q", key)
    }

    // Use reflection to set value
    v := reflect.ValueOf(c).Elem()
    field := v.FieldByName(fieldName)
    
    switch field.Kind() {
    case reflect.Int:
        intVal, _ := strconv.Atoi(value)
        field.SetInt(int64(intVal))
    case reflect.Bool:
        boolVal, _ := strconv.ParseBool(value)
        field.SetBool(boolVal)
    case reflect.String:
        // Add enum validation here
        field.SetString(value)
    }

    return nil
}
```

**Benefits**:
- Add new field: just update `fieldMap` (1 line)
- No duplicate type conversion logic
- Easier to test

---

### Example: Table-Based Config List

```go
func runConfigList() error {
    cfg := viperconfig.LoadFromViper()
    
    table := tablewriter.NewWriter(os.Stdout)
    table.SetHeader([]string{"Key", "Value", "Description"})
    table.SetBorder(false)
    
    // Time period section
    table.Append([]string{"days", fmt.Sprintf("%d", cfg.Days), "Days to look back"})
    table.Append([]string{"format", cfg.Format, "Output format"})
    // ...
    
    table.Render()
    return nil
}
```

**Benefits**:
- Consistent with `gohome report` styling
- Cleaner code
- Easier to add columns (e.g., "Source" column showing where value came from)

---

## 📝 Recommendations Summary

### Must Do 🔴
1. ✅ **Create comprehensive docs** (DONE - see `docs/VIPER_CONFIG_MANAGEMENT.md`)
2. **Move `RegisterAlias()` to package `init()`** (simple, no risk)
3. **Add inline comments** explaining Viper mechanics in code

### Should Do 🟡
4. **Refactor `setConfigValue()` to use reflection**
5. **Improve config list UI** with tablewriter
6. **Add validation layer** (separate from Viper logic)

### Nice to Have 🟢
7. **Add `--verbose` flag** to show config precedence debugging
8. **Add `config validate` subcommand** to check config file
9. **Support YAML config** for users who prefer it

---

## 🧪 Testing Recommendations

After refactoring, test these scenarios:

1. **Precedence**:
   ```bash
   # Config has days=3
   ./bin/gohome -d 7  # Should use 7
   GOHOME_DAYS=5 ./bin/gohome  # Should use 5
   ```

2. **Config set with type conversion**:
   ```bash
   ./bin/gohome config set today false  # Should create bool, not string
   ./bin/gohome config set days abc     # Should error clearly
   ```

3. **Flag binding**:
   ```bash
   ./bin/gohome -T -S  # All short flags work
   ./bin/gohome -b main -m 3  # New short flags work
   ```

4. **Clean JSON output**:
   ```bash
   ./bin/gohome -d 5 -S
   cat ~/.gohome.json | grep -c '"days"'  # Should be 1 (no duplicates)
   ```

---

## 📚 Related Documentation

- `docs/VIPER_CONFIG_MANAGEMENT.md` - Complete Viper guide (NEW!)
- `docs/CLI_GUIDE.md` - User-facing CLI documentation
- `docs/COBRA_VIPER_MIGRATION.md` - Migration history
- `.github/copilot-instructions.md` - AI agent guidelines

---

## Conclusion

Code hiện tại **functional và correct** nhưng có room for improvement về:
- **Maintainability**: Giant switch cần refactor
- **Performance**: RegisterAlias placement
- **UI/UX**: Config list có thể đẹp hơn

Tuy nhiên, **không cần rush refactor tất cả**. Ưu tiên:
1. Documentation (DONE ✅)
2. Quick wins (RegisterAlias move)
3. Refactor khi add new features (giảm regression risk)

Follow **Boy Scout Rule**: "Leave code cleaner than you found it" - refactor incrementally khi touch các areas này.

---

**Status**: Documentation complete, refactoring deferred (ready to implement when needed)  
**Risk Level**: Current code is stable, refactoring is optimization not bugfix  
**Next Steps**: Review docs, then decide which refactorings to prioritize

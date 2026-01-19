# Codecov Guide - Why & How to Use

**Date:** 2026-01-19  
**Tool:** Codecov - Code Coverage Tracking & Visualization  
**Website:** https://codecov.io

---

## 🎯 What is Codecov?

**Codecov** là một **cloud-based service** giúp track, visualize và manage code coverage cho project của bạn. Nó tích hợp với CI/CD để tự động thu thập coverage reports và hiển thị trên UI thân thiện.

### Core Features

1. **Visual Coverage Reports** - Xem coverage qua UI đẹp, có màu sắc
2. **Pull Request Comments** - Tự động comment coverage trên PR/MR
3. **Coverage Trends** - Track coverage theo thời gian
4. **File-level Analysis** - Xem từng file, từng line được cover hay chưa
5. **Status Checks** - Pass/fail PR based on coverage thresholds
6. **Team Dashboard** - Quản lý coverage cho toàn team/organization

---

## 🤔 Tại sao nên dùng Codecov?

### 1. **Visibility & Transparency** 👁️

**Vấn đề không có Codecov:**
```bash
# Local coverage
go test -cover ./...
# Output: 45.1% coverage

# Nhưng:
- Không biết file nào coverage thấp
- Không track được xu hướng (tăng/giảm)
- Team members không thấy được
- Không enforce coverage standards
```

**Với Codecov:**
- ✅ Visual dashboard cho toàn project
- ✅ File-by-file breakdown
- ✅ Line-by-line highlighting (red = not covered)
- ✅ Historical trends graph
- ✅ Compare coverage giữa branches

### 2. **Pull Request Integration** 🔍

**Automated PR Comments:**
```markdown
## Codecov Report
Coverage: 45.1% → 47.3% (+2.2%) ✅
Files Changed: 3
- src/parser.go: 100% ✅
- src/config.go: 65% → 75% (+10%) ✅
- src/main.go: excluded
```

**Benefits:**
- ✅ Reviewers thấy ngay impact lên coverage
- ✅ Catch untested code trước khi merge
- ✅ Encourage devs viết tests
- ✅ Prevent coverage regressions

### 3. **Quality Gates** 🚦

**Status Checks on PRs:**
```yaml
# codecov.yml
status:
  project:
    target: 50%      # Overall project must be ≥50%
    threshold: 5%    # Allow max -5% drop
    
  patch:
    target: 60%      # New code must be ≥60% covered
    threshold: 10%   # Allow max -10% drop on new code
```

**What happens:**
- ❌ PR blocked if coverage drops below threshold
- ✅ PR passes if coverage meets requirements
- 🟡 Warning if approaching threshold

### 4. **Team Accountability** 👥

**For Teams:**
- **Managers:** Track overall code quality
- **Tech Leads:** Identify weak spots
- **Developers:** Know what to test
- **Reviewers:** Objective metric for code reviews

### 5. **Free for Open Source** 💰

- ✅ Unlimited public repos
- ✅ Unlimited users
- ✅ Full features
- ✅ No credit card required

---

## 📊 Understanding Codecov Configuration

### Our Current `.codecov.yml`

```yaml
# Codecov Configuration
# https://docs.codecov.com/docs/codecov-yaml

coverage:
  precision: 2            # Show 2 decimal places (e.g., 45.12%)
  round: down            # Round 45.156% → 45.15% (not 45.16%)
  range: "70...100"      # Color range: red<70%, yellow 70-100%, green≥100%

  status:
    project:
      default:
        target: 50%      # Overall project must maintain ≥50% coverage
        threshold: 5%    # Allow max -5% drop from target

    patch:
      default:
        target: 60%      # New code (in PR) must be ≥60% covered
        threshold: 10%   # Allow max -10% drop for new code

ignore:
  - "cmd/gohome/main.go"    # Exclude entry point (0% OK)
  - "**/*_test.go"          # Exclude test files themselves
  - "**/example_*.go"       # Exclude example code
  - "scripts/**"            # Exclude shell scripts
  - "npm-package/**"        # Exclude NPM wrapper

comment:
  layout: "reach,diff,flags,files"  # What to show in PR comment
  behavior: default                  # Comment on every PR
  require_changes: false             # Comment even if no coverage change
```

### Breaking Down Each Section

#### 1. Coverage Display Settings

```yaml
coverage:
  precision: 2        # Decimals to show
  round: down        # Rounding method
  range: "70...100"  # Color coding
```

**Precision:**
- `0`: 45%
- `1`: 45.1%
- `2`: 45.12% ← **Our choice** (balanced)
- `5`: 45.12345%

**Round:**
- `down`: 45.156% → 45.15% ← **Our choice** (conservative)
- `up`: 45.151% → 45.16% (optimistic)
- `nearest`: 45.155% → 45.16% (standard)

**Range (Color Coding):**
```
"70...100" means:
- Red:    < 70%  (poor coverage)
- Yellow: 70-99% (acceptable)
- Green:  ≥ 100% (impossible, but means >99.99%)

Our range is strict (70% cutoff)
```

**Why we chose "70...100":**
- Critical packages (parser, git) have 90%+ → Green ✅
- CLI packages (cmd/*) have 40-60% → Yellow/Red ⚠️
- Encourages improvement to get green

#### 2. Status Checks (Quality Gates)

```yaml
status:
  project:
    default:
      target: 50%
      threshold: 5%
```

**Project Status (Overall Coverage):**

| Scenario | Current | After PR | Threshold Check | Result |
|----------|---------|----------|----------------|--------|
| Case 1 | 50% | 52% | +2% (within ±5%) | ✅ Pass |
| Case 2 | 50% | 48% | -2% (within ±5%) | ✅ Pass |
| Case 3 | 50% | 44% | -6% (exceeds -5%) | ❌ Fail |
| Case 4 | 50% | 45% | -5% (exactly at limit) | ✅ Pass |

**Formula:**
```
Pass if: current_coverage >= (target - threshold)
Pass if: 44% >= (50% - 5%) = 45%  → Fail!
Pass if: 46% >= (50% - 5%) = 45%  → Pass!
```

**Patch Status (New Code Only):**

```yaml
patch:
  default:
    target: 60%
    threshold: 10%
```

This checks **only the lines changed in PR**.

**Example:**
```diff
PR adds 100 lines:
- 65 lines covered by tests
- 35 lines not covered

Patch coverage: 65/100 = 65% ✅ (target 60%, threshold ±10%)
```

**Why separate patch check?**
- ✅ Enforce higher standards for NEW code
- ✅ Don't penalize for old untested code
- ✅ Gradually improve coverage over time
- ✅ New features must include tests

#### 3. Ignore Patterns

```yaml
ignore:
  - "cmd/gohome/main.go"    # Entry point
  - "**/*_test.go"          # Test files
  - "**/example_*.go"       # Examples
  - "scripts/**"            # Scripts
  - "npm-package/**"        # NPM
```

**Why ignore these?**

**`main.go`:**
```go
// main.go is just 3 lines:
func main() {
    cmd.Execute()  // All logic is in cmd.Execute()
}
// Testing main() requires subprocess, not worth it
// cmd.Execute() is already tested
```

**`*_test.go`:**
- Test files don't need coverage tests
- Would create infinite recursion
- Industry standard to exclude

**`example_*.go`:**
- Example code is for documentation
- Often isolated from main codebase
- Not production code

**`scripts/**` & `npm-package/**`:**
- Shell scripts (not Go code)
- NPM wrapper (JavaScript)
- Different testing methodologies

#### 4. Pull Request Comments

```yaml
comment:
  layout: "reach,diff,flags,files"
  behavior: default
  require_changes: false
```

**Layout Components:**

| Component | Shows |
|-----------|-------|
| `reach` | Overall coverage % |
| `diff` | Coverage change (+2.3%) |
| `flags` | Coverage by flag (e.g., unit vs integration) |
| `files` | File-by-file breakdown |

**Example Comment:**
```markdown
## Codecov Report
Coverage: 45.1% (+2.3%) ✅

Files Changed Coverage Δ:
- parser/parser.go: 100% (+5%) ✅
- cmd/config.go: 55% (+10%) ✅
- scanner/scanner.go: 95% (ø) 

Legend:
✅ Improved
ø No change
❌ Decreased
```

**Behavior:**
- `default`: Comment on every PR
- `once`: Only comment first time
- `new`: Only comment if coverage changed

**require_changes:**
- `false`: Always comment (even if 0% change)
- `true`: Only comment if coverage changed

---

## 🚀 How Codecov Works (Flow)

### End-to-End Process

```
1. Developer pushes code
   ↓
2. GitHub Actions triggers
   ↓
3. Run tests with coverage
   go test -coverprofile=coverage.txt -covermode=atomic ./...
   ↓
4. Upload to Codecov
   uses: codecov/codecov-action@v3
   ↓
5. Codecov processes report
   - Parse coverage.txt
   - Apply .codecov.yml config
   - Calculate metrics
   ↓
6. Codecov posts results
   - Comment on PR
   - Update status check
   - Update dashboard
   ↓
7. Developer sees results
   - ✅ or ❌ on PR
   - Detailed comment
   - Line-by-line view on Codecov.io
```

### GitHub Actions Integration

**Current `.github/workflows/release.yml`:**

```yaml
- name: Run tests
  run: |
    go test -v -coverprofile=coverage.txt ./...
    mkdir -p /tmp/coverage
    mv coverage.txt /tmp/coverage/ || true

- name: Upload coverage to Codecov
  uses: codecov/codecov-action@v3
  with:
    files: /tmp/coverage/coverage.txt
```

**What happens:**
1. Tests run, generate `coverage.txt`
2. `codecov-action` uploads file to Codecov
3. Codecov API processes coverage
4. Results appear on Codecov.io dashboard
5. PR comment auto-posted if in PR context

---

## 🎨 Visual Examples

### 1. Dashboard View

```
Project: gohome
Overall Coverage: 45.1% (+2.3% from last week)

[==========>          ] 45.1%

Top Files:
✅ internal/parser/parser.go:      100% ████████████
✅ internal/git/client.go:         92%  ███████████░
⚠️  cmd/gohome/cmd/config.go:      55%  ██████░░░░░░
❌ cmd/gohome/cmd/report.go:       15%  ██░░░░░░░░░░

Trend (last 30 days):
  45% ┤     ╭╮
  40% ┤    ╭╯╰╮
  35% ┤   ╭╯  ╰╮
  30% ┤  ╭╯    ╰─
  25% ┼──╯
```

### 2. File View (Line-by-Line)

```go
// config.go - Coverage: 55%

 1  package cmd
 2  
 3  ✅ func LoadConfig() {           // Covered
 4  ✅     cfg := viper.New()
 5  ✅     cfg.SetDefault("days", 7)
 6  ✅     cfg.ReadInConfig()
 7  ✅ }
 8  
 9  ❌ func SaveConfig() {           // NOT covered
10  ❌     data := viper.AllSettings()
11  ❌     writeFile(data)
12  ❌ }
```

### 3. PR Status Check

```
✅ codecov/project — Coverage: 45.1% (+2.3%)
   Target: 50%, Threshold: ±5%, Actual: 47.4%
   Status: PASSED (within threshold)
   
✅ codecov/patch — New code: 65% coverage
   Target: 60%, Threshold: ±10%
   Status: PASSED
   
View detailed report: https://codecov.io/gh/.../pull/18
```

---

## ⚙️ Configuration Best Practices

### 1. Set Realistic Targets

**DON'T:**
```yaml
status:
  project:
    target: 90%  # Unrealistic for existing project
    threshold: 0% # Too strict, blocks everything
```

**DO:**
```yaml
status:
  project:
    target: 50%   # Achievable current state
    threshold: 5% # Allow some flexibility
```

**Strategy:**
- Start with **current coverage** as target
- Gradually increase (e.g., +5% every quarter)
- Patch target > project target (new code should be better)

### 2. Use Appropriate Thresholds

**Threshold = How much drop is acceptable**

| Project Type | Recommended Threshold |
|--------------|----------------------|
| New project | 0-2% (strict) |
| Established project | 3-5% (balanced) |
| Legacy project | 5-10% (lenient) |
| Refactoring | 10-20% (temporary) |

**Our choice: 5% project, 10% patch** (balanced for beta)

### 3. Ignore Wisely

**Good ignores:**
- ✅ Entry points (main.go)
- ✅ Generated code
- ✅ Example code
- ✅ Test files
- ✅ Scripts (non-Go)

**Bad ignores:**
- ❌ "Hard to test" code (refactor instead!)
- ❌ Business logic
- ❌ Error handling
- ❌ Entire packages

### 4. Comment Settings

**For active projects:**
```yaml
comment:
  behavior: default      # Comment on every PR
  require_changes: true  # Only if coverage changed
```

**For quiet projects:**
```yaml
comment:
  behavior: once         # Only first time
  require_changes: false # Always show status
```

**Our choice: default + false** (always visible)

---

## 🔧 Advanced Features

### 1. Flags (Separate Coverage by Type)

```yaml
# Upload with flags
- name: Upload unit tests
  uses: codecov/codecov-action@v3
  with:
    files: coverage-unit.txt
    flags: unit

- name: Upload integration tests
  uses: codecov/codecov-action@v3
  with:
    files: coverage-integration.txt
    flags: integration
```

**View separate:**
- Unit test coverage: 75%
- Integration test coverage: 60%
- Combined: 80%

### 2. Multiple Projects

```yaml
coverage:
  status:
    project:
      backend:
        target: 70%
        paths:
          - "internal/**"
      
      cli:
        target: 50%
        paths:
          - "cmd/**"
```

Different targets for different parts!

### 3. Carryforward Flags

```yaml
flags:
  unit:
    carryforward: true  # Use previous if not uploaded
```

**Use case:** Don't run all tests every PR (expensive)

---

## 📈 Measuring Success

### Good Metrics

**1. Trend Direction**
```
Week 1: 40%
Week 2: 42% ↗️
Week 3: 43% ↗️
Week 4: 45% ↗️  ✅ Good trend!
```

**2. Patch Coverage**
```
Last 10 PRs:
- PR #18: +65% new code coverage ✅
- PR #17: +80% new code coverage ✅
- PR #16: +45% new code coverage ⚠️
- PR #15: +90% new code coverage ✅

Average: 70% ✅ Great!
```

**3. Critical Path Coverage**
```
Parser:  100% ✅ (critical)
Git:     92%  ✅ (critical)
Scanner: 95%  ✅ (critical)
CLI:     45%  ⚠️ (less critical, OK)
```

### Bad Metrics to Avoid

❌ **Vanity 100% coverage** (tests everything, even trivial)  
❌ **Gaming the system** (tests that don't assert anything)  
❌ **Coverage obsession** (100% coverage ≠ bug-free)

---

## 🎯 Summary

### Why Use Codecov?

| Benefit | Impact |
|---------|--------|
| **Visibility** | Everyone sees coverage, not just devs |
| **Automation** | No manual coverage checks |
| **Accountability** | Objective metric for code reviews |
| **Trends** | Track improvement over time |
| **Quality Gates** | Prevent regressions automatically |
| **Free** | No cost for open source |

### Our Configuration Explained

```yaml
# Conservative display (round down)
precision: 2, round: down, range: 70-100

# Realistic targets (current state)
project target: 50%, patch target: 60%

# Flexible thresholds (beta project)
project threshold: 5%, patch threshold: 10%

# Smart ignores (entry points, tests, scripts)
ignore: main.go, *_test.go, examples, scripts

# Always visible (transparent)
comment: default behavior, always show
```

### Quick Start

1. **Add `.codecov.yml`** to repo root ✅ (done)
2. **Upload in CI** via codecov-action ✅ (done)
3. **View results** at codecov.io
4. **Iterate** on targets as coverage improves

### Resources

- **Docs:** https://docs.codecov.com
- **Validator:** https://codecov.io/validate
- **Dashboard:** https://codecov.io/gh/anIcedAntFA/gohome
- **Support:** https://community.codecov.com

---

## ✅ Validation

**Our current `.codecov.yml` is VALID:**

```bash
$ curl --data-binary @.codecov.yml https://codecov.io/validate
Valid! ✅
```

**Confirmed settings:**
- ✅ Coverage precision: 2 decimals
- ✅ Rounding: down (conservative)
- ✅ Range: 70-100 (strict yellow cutoff)
- ✅ Project target: 50% (±5%)
- ✅ Patch target: 60% (±10%)
- ✅ Ignores: main.go, tests, examples, scripts
- ✅ PR comments: always enabled

**Ready to use!** 🚀

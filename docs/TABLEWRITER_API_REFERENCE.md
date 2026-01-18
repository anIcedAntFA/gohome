# TableWriter v1.1.3+ API Reference

> **Reference guide for tablewriter v1.1.3+ used in gohome**

## Version Information

- **Current Version**: v1.1.3 (latest)
- **Legacy Version**: v0.0.5 (deprecated, DO NOT USE)
- **Breaking Version**: v1.0.0 (has missing functionality, SKIP IT)

## Installation

```bash
go get github.com/olekukonko/tablewriter@v1.1.3
```

---

## Core API Changes (v1.1.3 vs v0.0.5)

### ❌ OLD API (v0.0.5) - DO NOT USE

```go
// Legacy approach (DEPRECATED)
table := tablewriter.NewWriter(os.Stdout)
table.SetHeader([]string{"Name", "Age"})      // ❌ Old
table.SetBorder(false)                        // ❌ Old
table.SetColumnSeparator("")                  // ❌ Old
table.SetRowSeparator("")                     // ❌ Old
table.SetHeaderLine(false)                    // ❌ Old
table.SetAlignment(tablewriter.ALIGN_LEFT)    // ❌ Old
table.SetHeaderAlignment(tablewriter.ALIGN_LEFT) // ❌ Old
table.Append([]string{"Alice", "25"})
table.Render()
```

### ✅ NEW API (v1.1.3+) - USE THIS

```go
// Modern approach (CURRENT)
table := tablewriter.NewTable(os.Stdout,  // Note: NewTable, not NewWriter
    tablewriter.WithConfig(tablewriter.Config{
        Header: tw.CellConfig{
            Alignment: tw.CellAlignment{Global: tw.AlignCenter},
        },
        Row: tw.CellConfig{
            Alignment: tw.CellAlignment{Global: tw.AlignLeft},
        },
    }),
)
table.Header([]string{"Name", "Age"})  // ✅ New - no "Set" prefix
table.Append([]string{"Alice", "25"})
table.Render()
```

---

## gohome Usage Pattern

### Recommended Pattern for Config List

```go
import (
    "github.com/olekukonko/tablewriter"
    "github.com/olekukonko/tablewriter/tw"
    "os"
)

func printConfigList(cfg *Config) {
    table := tablewriter.NewTable(os.Stdout,
        tablewriter.WithConfig(tablewriter.Config{
            Header: tw.CellConfig{
                Alignment: tw.CellAlignment{Global: tw.AlignCenter},
            },
            Row: tw.CellConfig{
                Alignment: tw.CellAlignment{Global: tw.AlignLeft},
            },
        }),
    )

    // Header without "Set" prefix
    table.Header([]string{"Key", "Value", "Description"})

    // Append rows
    table.Append([]string{"days", fmt.Sprintf("%d", cfg.Days), "Days to look back"})
    table.Append([]string{"format", cfg.Format, "Output format"})
    
    // Render
    table.Render()
}
```

---

## Key Differences Summary

| Feature | v0.0.5 (OLD) | v1.1.3+ (NEW) |
|---------|-------------|---------------|
| **Constructor** | `NewWriter()` | `NewTable()` |
| **Header** | `SetHeader()` | `Header()` |
| **Config** | Individual `Set*()` methods | `WithConfig()` option |
| **Alignment** | `SetAlignment(ALIGN_LEFT)` | `tw.CellAlignment{Global: tw.AlignLeft}` |
| **Borders** | `SetBorder(false)` | `tw.Rendition{Borders: tw.BorderNone}` |
| **Separators** | `SetColumnSeparator("")` | `tw.Separators{BetweenColumns: tw.Off}` |

---

## Configuration Structure

### Config Hierarchy

```
tablewriter.Config
├── Behavior (AutoHide, Compact, TrimWhitespace)
├── Header (CellConfig)
├── Row (CellConfig)
└── Footer (CellConfig)

tw.CellConfig
├── Formatting (AutoFormat, AutoWrap)
├── Alignment (Global, PerColumn)
├── Padding (Global, PerColumn)
├── Merging (Mode: None, Horizontal, Vertical, Both, Hierarchical)
└── Width (ColMaxWidths, ColMinWidths)
```

### Example: Complete Config

```go
config := tablewriter.Config{
    // Table-level behavior
    Behavior: tw.Behavior{
        AutoHideEmptyColumns: tw.Off,
        TrimWhitespace:       tw.On,
    },
    
    // Header styling
    Header: tw.CellConfig{
        Formatting: tw.CellFormatting{AutoFormat: tw.On},
        Alignment:  tw.CellAlignment{Global: tw.AlignCenter},
    },
    
    // Row styling
    Row: tw.CellConfig{
        Alignment:  tw.CellAlignment{Global: tw.AlignLeft},
        ColMaxWidths: tw.CellWidth{Global: 30},
    },
    
    // Footer styling
    Footer: tw.CellConfig{
        Alignment: tw.CellAlignment{Global: tw.AlignRight},
    },
}

table := tablewriter.NewTable(os.Stdout, tablewriter.WithConfig(config))
```

---

## Common Patterns

### Pattern 1: Simple Table (No Borders)

```go
table := tablewriter.NewTable(os.Stdout)  // Default config
table.Header([]string{"Key", "Value"})
table.Append([]string{"name", "gohome"})
table.Render()
```

Output:
```
┌──────┬────────┐
│ KEY  │ VALUE  │
├──────┼────────┤
│ name │ gohome │
└──────┴────────┘
```

### Pattern 2: Clean List (Like Terminal Output)

```go
import "github.com/olekukonko/tablewriter/renderer"

// Minimal borders, no separators
table := tablewriter.NewTable(os.Stdout,
    tablewriter.WithRenderer(renderer.NewBlueprint(tw.Rendition{
        Settings: tw.Settings{
            Borders:    tw.BorderNone,
            Separators: tw.Separators{BetweenColumns: tw.On}, // Only column separator
        },
    })),
)
```

Output:
```
KEY  │ VALUE
─────┼───────
name │ gohome
```

### Pattern 3: Markdown Table

```go
import "github.com/olekukonko/tablewriter/renderer"

table := tablewriter.NewTable(os.Stdout,
    tablewriter.WithRenderer(renderer.NewMarkdown()),
)
table.Header([]string{"Key", "Value"})
table.Append([]string{"name", "gohome"})
table.Render()
```

Output:
```
|  KEY  | VALUE  |
|:-----:|:------:|
| name  | gohome |
```

---

## Alignment Options

```go
// Global alignment (all columns same)
Row: tw.CellConfig{
    Alignment: tw.CellAlignment{Global: tw.AlignLeft},
}

// Per-column alignment
Row: tw.CellConfig{
    Alignment: tw.CellAlignment{
        PerColumn: []tw.Align{
            tw.AlignLeft,   // Column 0
            tw.AlignCenter, // Column 1
            tw.AlignRight,  // Column 2
        },
    },
}

// Available alignments
tw.AlignLeft
tw.AlignCenter
tw.AlignRight
tw.Skip  // Use global alignment
```

---

## Width Control

```go
Row: tw.CellConfig{
    // Max width for all columns
    ColMaxWidths: tw.CellWidth{Global: 30},
    
    // Per-column max widths
    ColMaxWidths: tw.CellWidth{
        PerColumn: []int{10, 20, 30},
    },
    
    // Min widths
    ColMinWidths: tw.CellWidth{Global: 5},
}
```

---

## Wrapping & Formatting

```go
Row: tw.CellConfig{
    Formatting: tw.CellFormatting{
        AutoWrap:   tw.WrapNormal,  // Wrap long text
        AutoFormat: tw.On,           // Auto-capitalize headers
    },
}

// Wrap modes
tw.WrapDisabled  // No wrapping (truncate)
tw.WrapNormal    // Wrap at word boundaries
tw.WrapHard      // Wrap at character boundaries
```

---

## Empty Row Separators

```go
// Add empty row for visual grouping
table.Header([]string{"Section", "Value"})

// Group 1
table.Append([]string{"days", "7"})
table.Append([]string{"weeks", "0"})

// Empty separator
table.Append([]string{"", ""})

// Group 2
table.Append([]string{"format", "text"})
table.Append([]string{"style", "normal"})
```

---

## Migration Checklist

When updating tablewriter code:

- [ ] Replace `tablewriter.NewWriter()` → `tablewriter.NewTable()`
- [ ] Replace `SetHeader()` → `Header()`
- [ ] Replace `SetFooter()` → `Footer()`
- [ ] Replace individual `Set*()` methods with `WithConfig()`
- [ ] Update imports: add `"github.com/olekukonko/tablewriter/tw"`
- [ ] Replace `ALIGN_LEFT` → `tw.AlignLeft`
- [ ] Replace `SetBorder(false)` → `tw.Rendition{Borders: tw.BorderNone}`
- [ ] Test output matches expected format

---

## Common Gotchas

### 1. Wrong Constructor
```go
// ❌ Wrong (v0.0.5 API)
table := tablewriter.NewWriter(os.Stdout)

// ✅ Correct (v1.1.3 API)
table := tablewriter.NewTable(os.Stdout)
```

### 2. Missing "tw" Import
```go
// ❌ Wrong
import "github.com/olekukonko/tablewriter"
// ... tw.AlignLeft not found

// ✅ Correct
import (
    "github.com/olekukonko/tablewriter"
    "github.com/olekukonko/tablewriter/tw"
)
```

### 3. Using Set Methods
```go
// ❌ Wrong (doesn't exist in v1.1.3)
table.SetHeader([]string{"A", "B"})
table.SetBorder(false)

// ✅ Correct
table := tablewriter.NewTable(os.Stdout,
    tablewriter.WithConfig(tablewriter.Config{
        // Configure here
    }),
)
table.Header([]string{"A", "B"})
```

---

## Examples in gohome Codebase

### Current Usage (report.go)

See `internal/renderer/printer.go` for the current implementation:

```go
// Already using v1.1.3+ correctly
table := tablewriter.NewTable(w,
    tablewriter.WithConfig(tablewriter.Config{
        Header: tw.CellConfig{
            Alignment: tw.CellAlignment{Global: tw.AlignCenter},
        },
        Row: tw.CellConfig{
            Alignment: tw.CellAlignment{Global: tw.AlignLeft},
        },
    }),
)
```

### Recommended for Config List

```go
// Similar pattern as printer.go
table := tablewriter.NewTable(os.Stdout,
    tablewriter.WithConfig(tablewriter.Config{
        Header: tw.CellConfig{
            Formatting: tw.CellFormatting{AutoFormat: tw.On},
            Alignment:  tw.CellAlignment{Global: tw.AlignCenter},
        },
        Row: tw.CellConfig{
            Alignment: tw.CellAlignment{Global: tw.AlignLeft},
        },
    }),
)

table.Header([]string{"Key", "Value", "Description"})
// ... append rows
table.Render()
```

---

## References

- **Official Repo**: https://github.com/olekukonko/tablewriter
- **v1.1.3 Release**: https://github.com/olekukonko/tablewriter/releases/tag/v1.1.3
- **gohome Usage**: `internal/renderer/printer.go`
- **Migration Guide**: This document

---

**Last Updated**: January 18, 2026  
**Current gohome tablewriter version**: v1.1.3

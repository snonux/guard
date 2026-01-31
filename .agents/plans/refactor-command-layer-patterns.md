# Feature: Refactor Command Layer - Align All Commands to Standard Pattern

## Feature Description

Refactor all command files in `cmd/guard/commands/*.go` to follow the consistent patterns established in `add.go` and `remove.go`. This involves standardizing error handling, registry operations, warning/error printing, and exit code usage across all 17 command files.

## User Story

As a developer maintaining the guard-tool codebase
I want all command implementations to follow consistent patterns
So that the codebase is maintainable, predictable, and follows uniform error handling practices

## Problem Statement

Currently, the command layer has inconsistent patterns across files:
- 12 commands use `RunE` instead of the standard `Run` pattern
- 11 commands are missing registry save operations
- 8 commands are missing proper warning/error printing
- Error handling varies between `return err` and direct `os.Exit(1)` patterns
- Exit code usage is inconsistent

## Solution Statement

Standardize all commands to follow the pattern established in `add.go` and `remove.go`:
- Use `Run` functions with inline error handling
- Explicitly save registry with error handling
- Print warnings and errors using manager utility functions
- Use `os.Exit(1)` for error conditions

## Feature Metadata

**Feature Type**: Refactor
**Estimated Complexity**: Medium
**Primary Systems Affected**: Command Layer (cmd/guard/commands/)
**Dependencies**: None (internal refactoring)

---

## CONTEXT REFERENCES

### Relevant Codebase Files IMPORTANT: YOU MUST READ THESE FILES BEFORE IMPLEMENTING!

- `cmd/guard/commands/add.go` - Why: Reference pattern for Run function, error handling, registry save, warnings/errors
- `cmd/guard/commands/remove.go` - Why: Reference pattern for Run function, error handling, registry save, warnings/errors
- `internal/manager/warnings.go` (lines 200-220) - Why: PrintWarnings and PrintErrors function signatures
- `cmd/guard/commands/toggle.go` - Why: NEEDS refactoring - remove generic handler, add inline Run handlers
- `cmd/guard/commands/version.go` - Why: NEEDS refactoring - remove version check, update descriptions

### Files to Refactor (13 files need complete conversion)

- `cmd/guard/commands/uninstall.go` - Convert RunE to Run pattern
- `cmd/guard/commands/reset.go` - Convert RunE to Run pattern  
- `cmd/guard/commands/init.go` - Convert RunE to Run pattern
- `cmd/guard/commands/cleanup.go` - Convert RunE to Run pattern
- `cmd/guard/commands/destroy.go` - Convert RunE to Run pattern
- `cmd/guard/commands/clear.go` - Convert RunE to Run pattern
- `cmd/guard/commands/create.go` - Convert RunE to Run pattern
- `cmd/guard/commands/update.go` - Convert RunE to Run pattern
- `cmd/guard/commands/config.go` - Convert RunE to Run pattern
- `cmd/guard/commands/enable.go` - Convert RunE to Run pattern
- `cmd/guard/commands/disable.go` - Convert RunE to Run pattern
- `cmd/guard/commands/toggle.go` - Remove generic handler, add inline Run handlers
- `cmd/guard/commands/version.go` - Remove version check, update descriptions

### Files Needing Partial Updates (3 files)

- `cmd/guard/commands/show.go` - Add auto-detection, Examples, summary counts, collection display
- `cmd/guard/commands/info.go` - Update to exact output format specification

### Relevant Documentation YOU SHOULD READ THESE BEFORE IMPLEMENTING!

- [Cobra Command Documentation](https://pkg.go.dev/github.com/spf13/cobra#Command)
  - Specific section: Run vs RunE function differences
  - Why: Understanding when to use Run vs RunE patterns

### Patterns to Follow

**Standard Command Structure Pattern:**
```go
func NewXxxCmd() *cobra.Command {
    return &cobra.Command{
        Use:   "xxx <target>...",
        Short: "Description",
        Long: `Detailed description with examples:

Examples:
  guard xxx myfile.txt           - Action on file (auto-detected)
  guard xxx myfolder             - Action on folder (auto-detected if directory)
  guard xxx mycollection         - Action on collection (auto-detected)`,
        Run: func(cmd *cobra.Command, args []string) {
            // Implementation with inline error handling
        },
    }
}
```

**Error Handling Pattern:**
```go
if err := someOperation(); err != nil {
    fmt.Fprintf(os.Stderr, "Error: %v\n", err)
    os.Exit(1)
}
```

**Missing Arguments Error Pattern:**
```go
if len(args) == 0 {
    fmt.Fprintln(os.Stderr, "Error: No files specified")
    fmt.Fprintln(os.Stderr, "Usage: guard command <path>...")
    fmt.Fprintln(os.Stderr)
    fmt.Fprintln(os.Stderr, "Use 'guard help command' for more information.")
    os.Exit(1)
}
```

**Registry Save Pattern:**
```go
// Save registry
if err := mgr.SaveRegistry(); err != nil {
    fmt.Fprintf(os.Stderr, "Error: Failed to save registry: %v\n", err)
    os.Exit(1)
}
```

**Warning and Error Printing Pattern:**
```go
// Print warnings
manager.PrintWarnings(mgr.GetWarnings())

// Print errors  
manager.PrintErrors(mgr.GetErrors())

// Exit with error code if there were errors
if mgr.HasErrors() {
    os.Exit(1)
}
```

**Toggle Command Specific Patterns:**

**NO Generic Handler - Each Subcommand Inline:**
```go
// WRONG - Don't use generic toggleHandler function
func toggleHandler(targetType string) func(...) { ... }

// CORRECT - Each subcommand has inline Run handler
func newToggleFileCmd() *cobra.Command {
    return &cobra.Command{
        Use: "file <file>...",
        Run: func(cmd *cobra.Command, args []string) {
            // Inline implementation here
        },
    }
}
```

**Toggle Helper Function Pattern:**
```go
// Helper function that tracks state BEFORE toggling using maps
// Returns true on ERROR (hasError semantics), false on success
func toggleFilesWithOutput(mgr *manager.Manager, files []string) bool {
    // Track registration status BEFORE any toggling
    registrationMap := make(map[string]bool)
    guardStateMap := make(map[string]bool)
    
    for _, file := range files {
        registrationMap[file] = mgr.IsRegisteredFile(file)
        if registrationMap[file] {
            if guard, err := mgr.GetRegistry().GetRegisteredFileGuard(file); err == nil {
                guardStateMap[file] = guard
            }
        }
    }
    
    // Toggle all files
    if err := mgr.ToggleFiles(files); err != nil {
        fmt.Fprintf(os.Stderr, "Error: %v\n", err)
        return true // ERROR
    }
    
    // Count newly registered files AFTER toggling
    newlyRegistered := 0
    for _, file := range files {
        if !registrationMap[file] && mgr.IsRegisteredFile(file) {
            newlyRegistered++
        }
    }
    
    // Print registration message BEFORE individual status messages
    if newlyRegistered > 0 {
        fmt.Printf("Registered %d file(s)\n", newlyRegistered)
    }
    
    // Print results based on previous state
    for _, file := range files {
        wasRegistered := registrationMap[file]
        wasGuarded := guardStateMap[file]
        
        if !wasRegistered {
            fmt.Printf("Guard enabled for %s\n", file)
        } else if wasGuarded {
            fmt.Printf("Guard disabled for %s\n", file)
        } else {
            fmt.Printf("Guard enabled for %s\n", file)
        }
    }
    
    return false // SUCCESS
}
```

**Info Command Exact Output Pattern:**
```go
fmt.Println("Guard - File Permission Management Tool")
fmt.Println()
fmt.Println("Created by Florian Buetow")
fmt.Println("Source code available at github.com/florianbuetow/guard")
fmt.Println()
fmt.Println("Guard helps protect your files from accidental modifications")
fmt.Println("by managing file permissions, ownership, and group settings.")
```

**Version Command Pattern:**
```go
// Remove version check - keep simple
return &cobra.Command{
    Use:   "version",
    Short: "Display version information",
    Long:  "Display the current version of the guard binary.",
    Run: func(cmd *cobra.Command, args []string) {
        fmt.Printf("guard version %s\n", version)
    },
}
```

**Show Command Auto-Detection Pattern:**
```go
// Main show command must use ResolveArguments
files, folders, collections, err := mgr.ResolveArguments(args)
if err != nil {
    fmt.Fprintf(os.Stderr, "Error: %v\n", err)
    os.Exit(1)
}
```

**Show File Info Pattern:**
```go
// printFileInfo must ALWAYS show collections even if empty
func printFileInfo(info manager.FileInfo) {
    prefix := "-"
    if info.Guard {
        prefix = "G"
    }
    
    // Always show collections, even if empty
    fmt.Printf("%s %s (%s)\n", prefix, info.Path, strings.Join(info.Collections, ", "))
}
```

**Show All Files Summary Pattern:**
```go
// When showing all files (no args), print summary at end
func showAllFiles(mgr *manager.Manager) {
    files, err := mgr.ShowFiles(nil)
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error: %v\n", err)
        os.Exit(1)
    }
    
    guarded := 0
    unguarded := 0
    
    for _, file := range files {
        printFileInfo(file)
        if file.Guard {
            guarded++
        } else {
            unguarded++
        }
    }
    
    total := len(files)
    fmt.Printf("%d file(s) total: %d guarded, %d unguarded\n", total, guarded, unguarded)
}
```

**Enable Command Status Tracking Pattern:**
```go
// Track registration and enable status
alreadyRegistered := 0
alreadyEnabled := 0
for _, path := range args {
    if mgr.IsRegisteredFile(path) {
        alreadyRegistered++
        if guard, err := mgr.GetRegistry().GetRegisteredFileGuard(path); err == nil && guard {
            alreadyEnabled++
        }
    }
}

// After enable operation
newlyRegistered := 0
nowEnabled := 0
for _, path := range args {
    if mgr.IsRegisteredFile(path) {
        newlyRegistered++
        if guard, err := mgr.GetRegistry().GetRegisteredFileGuard(path); err == nil && guard {
            nowEnabled++
        }
    }
}

// Print status messages
if newlyRegistered > alreadyRegistered {
    fmt.Printf("Registered %d file(s)\n", newlyRegistered-alreadyRegistered)
}
if nowEnabled > alreadyEnabled {
    fmt.Printf("Enabled %d file(s)\n", nowEnabled-alreadyEnabled)
}
if alreadyEnabled > 0 {
    fmt.Printf("Skipped %d file(s) already enabled\n", alreadyEnabled)
}
```

---

## IMPLEMENTATION PLAN

### Phase 1: Analyze Current State

Understand the current implementation patterns and identify specific changes needed for each command file.

### Phase 2: Convert RunE Commands

Convert all commands using `RunE` pattern to `Run` pattern with inline error handling.

### Phase 3: Add Missing Registry Operations

Add registry save operations with proper error handling to commands that modify state.

### Phase 4: Standardize Warning/Error Handling

Add consistent warning and error printing using manager utility functions.

### Phase 5: Validation

Ensure all commands follow the standard pattern and maintain functionality.

---

## STEP-BY-STEP TASKS

IMPORTANT: Execute every task in order, top to bottom. Each task is atomic and independently testable.

### UPDATE cmd/guard/commands/uninstall.go

- **IMPLEMENT**: Convert RunE to Run pattern with inline error handling, add Examples to Long description
- **PATTERN**: Mirror add.go Run function structure (file:add.go:25-30)
- **IMPORTS**: Ensure fmt and os imports are present
- **GOTCHA**: Maintain existing functionality while changing error handling approach
- **EXAMPLES**: Add Examples section to Long description
- **ERROR_FORMAT**: Use multi-line error format with usage hints
- **VALIDATE**: `go build -o build/guard ./cmd/guard && echo "Build successful"`

### UPDATE cmd/guard/commands/reset.go

- **IMPLEMENT**: Convert RunE to Run pattern, add registry save, add warning/error printing, add Examples
- **PATTERN**: Mirror add.go error handling and registry save (file:add.go:45-55, 65-75)
- **IMPORTS**: Add manager import for PrintWarnings/PrintErrors
- **GOTCHA**: Reset operations may not need registry save - verify business logic
- **EXAMPLES**: Add Examples section to Long description
- **ERROR_FORMAT**: Use multi-line error format with usage hints
- **VALIDATE**: `go build -o build/guard ./cmd/guard && echo "Build successful"`

### UPDATE cmd/guard/commands/init.go

- **IMPLEMENT**: Convert RunE to Run pattern with inline error handling, add Examples
- **PATTERN**: Mirror add.go Run function and error handling (file:add.go:25-40)
- **IMPORTS**: Ensure fmt and os imports are present
- **GOTCHA**: Init creates registry, so save operation is already handled in business logic
- **EXAMPLES**: Add Examples section to Long description
- **ERROR_FORMAT**: Use multi-line error format with usage hints
- **VALIDATE**: `go build -o build/guard ./cmd/guard && echo "Build successful"`

### UPDATE cmd/guard/commands/cleanup.go

- **IMPLEMENT**: Convert RunE to Run pattern, add registry save, add warning/error printing, add Examples
- **PATTERN**: Mirror add.go complete pattern (file:add.go:25-85)
- **IMPORTS**: Add manager import for PrintWarnings/PrintErrors
- **GOTCHA**: Cleanup modifies registry state, needs explicit save
- **EXAMPLES**: Add Examples section to Long description
- **ERROR_FORMAT**: Use multi-line error format with usage hints
- **VALIDATE**: `go build -o build/guard ./cmd/guard && echo "Build successful"`

### UPDATE cmd/guard/commands/destroy.go

- **IMPLEMENT**: Convert RunE to Run pattern, add registry save, add warning/error printing, add Examples
- **PATTERN**: Mirror add.go complete pattern (file:add.go:25-85)
- **IMPORTS**: Add manager import for PrintWarnings/PrintErrors
- **GOTCHA**: Destroy modifies registry state, needs explicit save
- **EXAMPLES**: Add Examples section to Long description
- **ERROR_FORMAT**: Use multi-line error format with usage hints
- **VALIDATE**: `go build -o build/guard ./cmd/guard && echo "Build successful"`

### UPDATE cmd/guard/commands/clear.go

- **IMPLEMENT**: Convert RunE to Run pattern, add registry save, add warning/error printing, add Examples
- **PATTERN**: Mirror add.go complete pattern (file:add.go:25-85)
- **IMPORTS**: Add manager import for PrintWarnings/PrintErrors
- **GOTCHA**: Clear modifies registry state, needs explicit save
- **EXAMPLES**: Add Examples section to Long description
- **ERROR_FORMAT**: Use multi-line error format with usage hints
- **VALIDATE**: `go build -o build/guard ./cmd/guard && echo "Build successful"`

### UPDATE cmd/guard/commands/create.go

- **IMPLEMENT**: Convert RunE to Run pattern, add registry save, add warning/error printing, add Examples
- **PATTERN**: Mirror add.go complete pattern (file:add.go:25-85)
- **IMPORTS**: Add manager import for PrintWarnings/PrintErrors
- **GOTCHA**: Create modifies registry state, needs explicit save
- **EXAMPLES**: Add Examples section to Long description
- **ERROR_FORMAT**: Use multi-line error format with usage hints
- **VALIDATE**: `go build -o build/guard ./cmd/guard && echo "Build successful"`

### UPDATE cmd/guard/commands/update.go

- **IMPLEMENT**: Convert RunE to Run pattern, add registry save, add warning/error printing, add Examples
- **PATTERN**: Mirror add.go complete pattern (file:add.go:25-85)
- **IMPORTS**: Add manager import for PrintWarnings/PrintErrors
- **GOTCHA**: Update modifies registry state, needs explicit save
- **EXAMPLES**: Add Examples section to Long description
- **ERROR_FORMAT**: Use multi-line error format with usage hints
- **VALIDATE**: `go build -o build/guard ./cmd/guard && echo "Build successful"`

### UPDATE cmd/guard/commands/config.go

- **IMPLEMENT**: Convert RunE to Run pattern, add registry save for modification operations, add warning/error printing, add Examples
- **PATTERN**: Mirror add.go complete pattern (file:add.go:25-85)
- **IMPORTS**: Add manager import for PrintWarnings/PrintErrors
- **GOTCHA**: Config may have read-only operations that don't need registry save - check business logic
- **EXAMPLES**: Add Examples section to Long description
- **ERROR_FORMAT**: Use multi-line error format with usage hints
- **VALIDATE**: `go build -o build/guard ./cmd/guard && echo "Build successful"`

### UPDATE cmd/guard/commands/toggle.go

- **IMPLEMENT**: Remove generic toggleHandler function, convert each subcommand to inline Run handlers
- **PATTERN**: Each subcommand uses inline Run function (no shared handler)
- **IMPORTS**: Ensure fmt, os, manager imports are present
- **GOTCHA**: Must track guard state BEFORE toggling using maps, not one-by-one in loop
- **HELPER**: Create toggleFilesWithOutput helper that returns bool for error indication
- **REGISTRATION**: Track registration status in maps before ANY toggling, count newly registered AFTER
- **VALIDATE**: `go build -o build/guard ./cmd/guard && echo "Build successful"`

### UPDATE cmd/guard/commands/version.go

- **IMPLEMENT**: Remove version check logic, update Short and Long descriptions
- **PATTERN**: Keep simple version display without conditional logic
- **IMPORTS**: Ensure fmt import is present
- **GOTCHA**: Remove "if version == "" { version = "dev" }" check completely
- **DESCRIPTIONS**: Short: "Display version information", Long: "Display the current version of the guard binary."
- **VALIDATE**: `go build -o build/guard ./cmd/guard && echo "Build successful"`

### UPDATE cmd/guard/commands/info.go

- **IMPLEMENT**: Convert RunE to Run pattern, update to exact output format specification
- **PATTERN**: Use exact output format with specific spacing and text
- **IMPORTS**: Ensure fmt import is present
- **GOTCHA**: Must match EXACTLY the specified output format including blank lines
- **OUTPUT**: "Guard - File Permission Management Tool" + blank line + author + source + blank line + description
- **VALIDATE**: `go build -o build/guard ./cmd/guard && echo "Build successful"`

### UPDATE cmd/guard/commands/show.go

- **IMPLEMENT**: Add auto-detection using ResolveArguments, add Examples, add summary counts, fix collection display
- **PATTERN**: Main command uses mgr.ResolveArguments for auto-detection
- **IMPORTS**: Add manager import for PrintWarnings/PrintErrors
- **GOTCHA**: printFileInfo must ALWAYS show collections even if empty like "G filename ()"
- **SUMMARY**: newShowFileCmd when no args must print "N file(s) total: X guarded, Y unguarded"
- **EXAMPLES**: Add Examples section to Long description
- **VALIDATE**: `go build -o build/guard ./cmd/guard && echo "Build successful"`

### UPDATE cmd/guard/commands/enable.go

- **IMPLEMENT**: Convert RunE to Run pattern, add status tracking for registration/enable counts
- **PATTERN**: Mirror add.go Run function structure with status message printing
- **IMPORTS**: Add manager import for PrintWarnings/PrintErrors
- **GOTCHA**: Must track "already registered", "already enabled", "newly registered" counts
- **STATUS**: Print "Registered N file(s)", "Enabled N file(s)", "Skipped N file(s) already enabled"
- **VALIDATE**: `go build -o build/guard ./cmd/guard && echo "Build successful"`

### UPDATE cmd/guard/commands/disable.go

- **IMPLEMENT**: Convert RunE to Run pattern, add registry save, add warning/error printing, add Examples
- **PATTERN**: Mirror add.go complete pattern with Long description including Examples section
- **IMPORTS**: Add manager import for PrintWarnings/PrintErrors
- **GOTCHA**: Disable modifies file states, needs explicit registry save
- **EXAMPLES**: Add Examples section to Long description
- **VALIDATE**: `go build -o build/guard ./cmd/guard && echo "Build successful"`

### UPDATE cmd/guard/commands/show.go

- **IMPLEMENT**: Add registry save (if needed), add warning/error printing, add Examples to Long description
- **PATTERN**: Mirror add.go warning/error printing pattern (file:add.go:75-85)
- **IMPORTS**: Add manager import for PrintWarnings/PrintErrors
- **GOTCHA**: Show is read-only, may not need registry save - verify business logic
- **EXAMPLES**: Add Examples section to Long description
- **ERROR_FORMAT**: Use multi-line error format with usage hints if applicable
- **VALIDATE**: `go build -o build/guard ./cmd/guard && echo "Build successful"`

---

## TESTING STRATEGY

### Unit Tests

No new unit tests required - this is a refactoring that maintains existing functionality while standardizing patterns.

### Integration Tests

Use existing shell integration tests to verify functionality is preserved:

### Edge Cases

- Verify error conditions still produce proper exit codes
- Ensure warning aggregation still works correctly
- Confirm registry save operations don't introduce race conditions

---

## VALIDATION COMMANDS

Execute every command to ensure zero regressions and 100% pattern consistency.

### Level 1: Syntax & Style

```bash
go fmt ./...
golangci-lint run
```

### Level 2: Build Validation

```bash
go build -o build/guard ./cmd/guard
```

### Level 3: Integration Tests

```bash
cp build/guard tests/guard
cd tests && ./run-all-tests.sh
```

### Level 4: Manual Validation

```bash
# Test each command maintains functionality
./guard --help
./guard init --help
./guard add --help
# ... test all commands for basic functionality
```

### Level 5: Pattern Consistency Check

```bash
# Verify all commands use Run (not RunE)
grep -r "RunE:" cmd/guard/commands/ || echo "No RunE patterns found - good!"

# Verify error handling patterns
grep -r "fmt.Fprintf(os.Stderr" cmd/guard/commands/ | wc -l

# Verify os.Exit usage
grep -r "os.Exit(1)" cmd/guard/commands/ | wc -l
```

---

## ACCEPTANCE CRITERIA

- [ ] All 17 command files use `Run` function (not `RunE`)
- [ ] All commands use consistent error handling: `fmt.Fprintf(os.Stderr, "Error: %v\n", err)` + `os.Exit(1)`
- [ ] All commands use multi-line error format for missing arguments with usage hints
- [ ] All state-modifying commands call `mgr.SaveRegistry()` with error handling
- [ ] All commands call `manager.PrintWarnings(mgr.GetWarnings())` AND `manager.PrintErrors(mgr.GetErrors())`
- [ ] All commands use `if mgr.HasErrors() { os.Exit(1) }` pattern
- [ ] All commands have Examples section in Long description
- [ ] Toggle command uses inline Run handlers (no generic toggleHandler function)
- [ ] Toggle command has toggleFilesWithOutput helper that returns bool and uses maps for state tracking
- [ ] Toggle command tracks registration in maps BEFORE toggling, counts newly registered AFTER
- [ ] Enable command tracks registration/enable status and prints appropriate counts
- [ ] Info command outputs EXACTLY the specified format with correct spacing
- [ ] Version command removes version check logic and uses specified descriptions
- [ ] Show command uses ResolveArguments for auto-detection in main command
- [ ] Show command printFileInfo ALWAYS shows collections even if empty: "G filename ()"
- [ ] Show file subcommand prints summary when no args: "N file(s) total: X guarded, Y unguarded"
- [ ] All validation commands pass with zero errors
- [ ] Full integration test suite passes
- [ ] No functional regressions introduced
- [ ] Code follows existing project conventions
- [ ] Build succeeds for all intermediate steps

---

## COMPLETION CHECKLIST

- [ ] All 12 RunE commands converted to Run pattern
- [ ] Registry save operations added where needed
- [ ] Warning/error printing added to all commands
- [ ] Consistent error handling implemented
- [ ] All build validations pass
- [ ] Integration tests pass
- [ ] Pattern consistency verified
- [ ] No functional changes to command behavior
- [ ] Exit codes remain consistent

---

## NOTES

**Design Decisions:**
- Maintaining existing command functionality while standardizing implementation patterns
- Using inline error handling instead of error return values for consistency
- Explicit registry save operations for better control and error handling

**Trade-offs:**
- Slightly more verbose error handling in exchange for consistency and explicit control
- Direct os.Exit calls instead of error propagation for clearer failure modes

**Key Risks:**
- Accidentally changing command behavior during refactoring
- Missing registry save operations for state-modifying commands
- Breaking existing shell integration tests

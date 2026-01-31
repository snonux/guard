# Feature: Refactor Command Layer with Proper Cobra Subcommand Structure

The following plan should be complete, but it's important that you validate documentation and codebase patterns and task sanity before you start implementing.

Pay special attention to naming of existing utils types and models. Import from the right files etc.

## Feature Description

Refactor the cmd/guard/commands layer to follow proper Cobra best practices with modular subcommand structures. The current implementation is too simplified and lacks proper command hierarchy. This refactoring will make commands more maintainable, testable, and aligned with Cobra's intended usage patterns.

## User Story

As a developer maintaining the guard-tool CLI
I want commands organized with proper subcommand structures
So that the codebase is more modular, maintainable, and follows Cobra best practices

## Problem Statement

The current command layer has several architectural issues:

1. **show.go** lacks explicit subcommands - should have `show file` and `show collection` subcommands with separate factory functions
2. **add.go** and **remove.go** lack proper subcommand structure and dedicated business logic functions
3. **version.go** hardcodes version string instead of accepting it as a parameter for ldflags injection
4. **keywords.go** duplicates validation logic that should be elsewhere
5. **main.go** uses init() function instead of defining rootCmd in main() for better testability
6. Commands mix CLI concerns with business logic instead of delegating to manager layer

## Solution Statement

Restructure commands to follow Cobra's factory pattern with proper subcommand hierarchy:
- Parent commands support auto-detection while providing explicit subcommands
- Each subcommand has its own factory function (e.g., `newShowFileCmd()`, `newShowCollectionCmd()`)
- Business logic extracted into separate functions for reusability
- Version string injected via ldflags at build time
- Remove keywords.go duplication
- Refactor main.go to avoid init() and enable better testing

## Feature Metadata

**Feature Type**: Refactor
**Estimated Complexity**: Medium
**Primary Systems Affected**: cmd/guard/commands/, cmd/guard/main.go
**Dependencies**: github.com/spf13/cobra v1.10.2

---

## CONTEXT REFERENCES

### Relevant Codebase Files IMPORTANT: YOU MUST READ THESE FILES BEFORE IMPLEMENTING!

- `cmd/guard/main.go` (lines 1-100) - Why: Current root command structure and init() usage to refactor
- `cmd/guard/commands/show.go` (lines 1-60) - Why: Needs subcommand structure for file/collection
- `cmd/guard/commands/add.go` (lines 1-70) - Why: Needs newAddFileCmd() and addFiles() extraction
- `cmd/guard/commands/remove.go` (lines 1-70) - Why: Needs newRemoveFileCmd() and removeFiles() extraction
- `cmd/guard/commands/version.go` (lines 1-20) - Why: Needs version parameter for ldflags
- `cmd/guard/commands/keywords.go` (lines 1-10) - Why: To be removed, constants moved elsewhere
- `cmd/guard/commands/enable.go` (lines 1-200) - Why: Good example of existing subcommand structure to mirror
- `cmd/guard/commands/toggle.go` (lines 1-150) - Why: Good example of subcommand pattern with factory functions
- `cmd/guard/commands/disable.go` (lines 1-200) - Why: Good example of subcommand structure
- `internal/manager/manager.go` (lines 1-260) - Why: Manager API for business logic delegation
- `internal/manager/files.go` (lines 20-140) - Why: AddFiles, RemoveFiles functions to call
- `internal/manager/collections.go` (lines 889-950) - Why: ShowCollections function to call
- `COBRA_BEST_PRACTICES.md` (lines 1-150) - Why: Project-specific Cobra patterns and conventions
- `tests/test-show-commands.sh` (lines 1-400) - Why: Test specifications for show command behavior

### New Files to Create

None - this is a refactoring of existing files

### Relevant Documentation YOU SHOULD READ THESE BEFORE IMPLEMENTING!

- [Cobra Official Docs - Working with Commands](https://cobra.dev/docs/how-to-guides/working-with-commands/)
  - Specific section: Command structure and subcommands
  - Why: Understand proper parent-child command hierarchy
- [Go Build with ldflags](https://jerrynsh.com/3-easy-ways-to-add-version-flag-in-go/)
  - Specific section: Using -X flag to set version variables
  - Why: Implement version injection at build time
- [Cobra Subcommand Best Practices](https://openillumi.com/en/en-go-cobra-cli-subcommand-package-separation/)
  - Specific section: Factory pattern for subcommands
  - Why: Modular command organization patterns

### Patterns to Follow

**Command Factory Pattern** (from enable.go lines 13-35):
```go
func NewEnableCmd() *cobra.Command {
    enableCmd := &cobra.Command{
        Use:   "enable <target>...",
        Short: "Enable protection for files, collections, or folders",
        RunE: runEnable,
    }
    
    // Add subcommands
    enableCmd.AddCommand(newEnableFileCmd())
    enableCmd.AddCommand(newEnableFolderCmd())
    enableCmd.AddCommand(newEnableCollectionCmd())
    
    return enableCmd
}
```

**Subcommand Factory Pattern** (from enable.go lines 150-165):
```go
func newEnableFileCmd() *cobra.Command {
    return &cobra.Command{
        Use:   "file <file>...",
        Short: "Enable protection for files",
        RunE: func(cmd *cobra.Command, args []string) error {
            if len(args) == 0 {
                return fmt.Errorf("No files specified")
            }
            
            mgr := manager.NewManager(".guardfile")
            if err := mgr.LoadRegistry(); err != nil {
                return fmt.Errorf("failed to load registry: %w", err)
            }
            
            return processEnableFiles(mgr, args)
        },
    }
}
```

**Business Logic Extraction Pattern** (from toggle.go lines 36-105):
```go
// Helper functions create manager internally and handle errors with os.Exit(1)
func addFiles(args []string) {
    mgr := manager.NewManager(".guardfile")
    if err := mgr.LoadRegistry(); err != nil {
        fmt.Fprintf(os.Stderr, "Error: failed to load registry: %v\n", err)
        os.Exit(1)
    }
    
    if err := mgr.AddFiles(args); err != nil {
        fmt.Fprintf(os.Stderr, "Error: %v\n", err)
        os.Exit(1)
    }
    
    if err := mgr.SaveRegistry(); err != nil {
        fmt.Fprintf(os.Stderr, "Error: failed to save registry: %v\n", err)
        os.Exit(1)
    }
    
    fmt.Printf("Registered %d file(s)\n", len(args))
    
    manager.PrintWarnings(mgr.GetWarnings())
    manager.PrintErrors(mgr.GetErrors())
}
```

**Manager Delegation Pattern** (from files.go lines 20-80):
```go
// Manager layer handles business logic
func (m *Manager) AddFiles(paths []string) error {
    // Validation and business logic here
    // Returns error for CLI layer to handle
}
```

**Version Injection Pattern** (ldflags best practice):
```go
// In main.go
var version = "dev" // Default, overridden by ldflags

func main() {
    rootCmd.AddCommand(commands.NewVersionCmd(version))
}

// In version.go
func NewVersionCmd(version string) *cobra.Command {
    return &cobra.Command{
        Use:   "version",
        Short: "Show version information",
        RunE: func(cmd *cobra.Command, args []string) error {
            fmt.Printf("guard version %s\n", version)
            return nil
        },
    }
}
```

**Root Command in main() Pattern** (COBRA_BEST_PRACTICES.md):
```go
func main() {
    rootCmd := &cobra.Command{
        Use:   "guard",
        Short: "Protect files from unwanted modifications",
    }
    
    // Add commands
    rootCmd.AddCommand(commands.NewInitCmd())
    // ... more commands
    
    if err := rootCmd.Execute(); err != nil {
        os.Exit(1)
    }
}
```

---

## IMPLEMENTATION PLAN

### Phase 1: Foundation - Remove Duplication and Prepare Structure

Remove keywords.go and prepare for refactoring by understanding current command patterns.

**Tasks:**
- Remove keywords.go file
- Move constants to appropriate command files if needed
- Verify all tests still pass

### Phase 2: Refactor main.go - Remove init() and Add Version Support

Restructure main.go to define rootCmd in main() function and support version injection via ldflags.

**Tasks:**
- Move rootCmd definition from package level to main() function
- Add version variable for ldflags injection
- Pass version to NewVersionCmd()
- Remove init() function
- Maintain all existing command registrations

### Phase 3: Refactor version.go - Accept Version Parameter

Update version command to accept version string parameter instead of hardcoding.

**Tasks:**
- Change NewVersionCmd() signature to accept version string
- Use provided version in output
- Remove hardcoded version string

### Phase 4: Refactor show.go - Add Subcommand Structure

Add proper subcommands for `show file` and `show collection` while maintaining auto-detection in parent command.

**Tasks:**
- Create showAllFiles() helper (void, no args)
- Create showSpecificFiles() helper (void, takes files array)
- Create printFileInfo() helper (void, takes FileInfo)
- Create showAllCollections() helper (void, no args)
- Create showSpecificCollections() helper (void, takes collections array)
- Create newShowFileCmd() factory function
- Create newShowCollectionCmd() factory function
- Update parent show command to add subcommands
- Maintain backward compatibility with auto-detection

### Phase 5: Refactor add.go - Add Subcommand Structure

Add proper subcommand structure and extract business logic that creates manager internally.

**Tasks:**
- Create addFiles() helper function (creates manager, handles errors with os.Exit)
- Create newAddFileCmd() factory function
- Update parent add command structure
- Remove keyword parsing logic

### Phase 6: Refactor remove.go - Add Subcommand Structure

Add proper subcommand structure and extract business logic that creates manager internally.

**Tasks:**
- Create removeFiles() helper function (creates manager, handles errors with os.Exit)
- Create newRemoveFileCmd() factory function
- Update parent remove command structure
- Remove keyword parsing logic

### Phase 7: Testing & Validation

Verify all commands work correctly with new structure.

**Tasks:**
- Run full test suite
- Verify show file/collection subcommands work
- Verify add/remove subcommands work
- Verify version displays correctly
- Test backward compatibility

---

## STEP-BY-STEP TASKS

IMPORTANT: Execute every task in order, top to bottom. Each task is atomic and independently testable.

### REMOVE cmd/guard/commands/keywords.go

- **IMPLEMENT**: Delete the keywords.go file as it duplicates validation logic
- **PATTERN**: Constants are already defined in individual command files where needed
- **GOTCHA**: Verify no other files import or reference keywords.go
- **VALIDATE**: `go build -o build/guard ./cmd/guard`

### UPDATE cmd/guard/commands/version.go

- **IMPLEMENT**: Change NewVersionCmd() to accept version string parameter
- **PATTERN**: Use Run instead of RunE - no error return needed
- **SIGNATURE**: `func NewVersionCmd(version string) *cobra.Command`
- **RUN**: Use `Run: func(cmd *cobra.Command, args []string) { fmt.Printf("guard version %s\n", version) }`
- **DEFAULT**: If version is empty, use "dev" as fallback
- **VALIDATE**: `go build -o build/guard ./cmd/guard`

### REFACTOR cmd/guard/main.go

- **IMPLEMENT**: Remove init() function and move rootCmd definition into main()
- **PATTERN**: Define rootCmd as local variable in main(), not package-level
- **ADD**: Version variable at package level: `var version = "dev"`
- **UPDATE**: Pass version to NewVersionCmd: `rootCmd.AddCommand(commands.NewVersionCmd(version))`
- **REMOVE**: init() function entirely
- **MOVE**: All rootCmd setup (SetHelpTemplate, PersistentFlags, AddCommand calls) into main()
- **GOTCHA**: Maintain exact same command registration order
- **VALIDATE**: `go build -o build/guard ./cmd/guard && ./build/guard --help`

### CREATE newShowFileCmd() in cmd/guard/commands/show.go

- **IMPLEMENT**: Add newShowFileCmd() factory function
- **PATTERN**: Mirror newToggleFileCmd() from toggle.go
- **SIGNATURE**: `func newShowFileCmd() *cobra.Command`
- **USE**: `"file [file]..."`
- **SHORT**: `"Display status of specific files"`
- **RUN**: Use `Run` (not RunE)
- **LOGIC**: Load manager, if no args call showAllFiles(mgr), else call showSpecificFiles(mgr, args)
- **VALIDATE**: `go build -o build/guard ./cmd/guard`

### CREATE showAllFiles() helper in cmd/guard/commands/show.go

- **IMPLEMENT**: Show all registered files
- **SIGNATURE**: `func showAllFiles(mgr *manager.Manager)` - void, no return
- **LOGIC**: Get all files from manager and print each one
- **ERROR HANDLING**: Use os.Exit(1) on errors
- **VALIDATE**: `go build -o build/guard ./cmd/guard`

### CREATE showSpecificFiles() helper in cmd/guard/commands/show.go

- **IMPLEMENT**: Show specific files by path
- **SIGNATURE**: `func showSpecificFiles(mgr *manager.Manager, files []string)` - void, no return
- **LOGIC**: Iterate files and call printFileInfo() for each
- **ERROR HANDLING**: Use os.Exit(1) on errors
- **VALIDATE**: `go build -o build/guard ./cmd/guard`

### CREATE printFileInfo() helper in cmd/guard/commands/show.go

- **IMPLEMENT**: Print single file information
- **SIGNATURE**: `func printFileInfo(info manager.FileInfo)` - void, no return
- **LOGIC**: Format and print file status (G/- prefix, path, collections)
- **FORMAT**: "G file.txt (coll1, coll2)" or "- file.txt" or "- file.txt (coll1)"
- **PATTERN**: Follow test expectations from test-show-commands.sh lines 130-180
- **VALIDATE**: `go build -o build/guard ./cmd/guard`

### CREATE newShowCollectionCmd() in cmd/guard/commands/show.go

- **IMPLEMENT**: Add newShowCollectionCmd() factory function
- **PATTERN**: Mirror newToggleCollectionCmd() from toggle.go
- **SIGNATURE**: `func newShowCollectionCmd() *cobra.Command`
- **USE**: `"collection [collection]..."`
- **SHORT**: `"Display status of specific collections"`
- **RUN**: Use `Run` (not RunE)
- **LOGIC**: Load manager, if no args call showAllCollections(mgr), else call showSpecificCollections(mgr, args)
- **VALIDATE**: `go build -o build/guard ./cmd/guard`

### CREATE showAllCollections() helper in cmd/guard/commands/show.go

- **IMPLEMENT**: Show all registered collections
- **SIGNATURE**: `func showAllCollections(mgr *manager.Manager)` - void, no return
- **LOGIC**: Get all collections from manager and print summary
- **ERROR HANDLING**: Use os.Exit(1) on errors
- **VALIDATE**: `go build -o build/guard ./cmd/guard`

### CREATE showSpecificCollections() helper in cmd/guard/commands/show.go

- **IMPLEMENT**: Show specific collections by name
- **SIGNATURE**: `func showSpecificCollections(mgr *manager.Manager, collections []string)` - void, no return
- **LOGIC**: Iterate collections and print details for each
- **ERROR HANDLING**: Use os.Exit(1) on errors
- **VALIDATE**: `go build -o build/guard ./cmd/guard`

### UPDATE NewShowCmd() in cmd/guard/commands/show.go

- **IMPLEMENT**: Add subcommands to parent show command
- **PATTERN**: Mirror NewEnableCmd() from enable.go lines 13-35
- **ADD**: `showCmd.AddCommand(newShowFileCmd())`
- **ADD**: `showCmd.AddCommand(newShowCollectionCmd())`
- **MAINTAIN**: Existing RunE logic for auto-detection when no subcommand specified
- **GOTCHA**: Parent command should still support `guard show file1.txt` (auto-detect) AND `guard show file file1.txt` (explicit)
- **VALIDATE**: `go build -o build/guard ./cmd/guard && ./build/guard show --help`

### CREATE newAddFileCmd() in cmd/guard/commands/add.go

- **IMPLEMENT**: Add newAddFileCmd() factory function
- **PATTERN**: Mirror newToggleFileCmd() from toggle.go
- **SIGNATURE**: `func newAddFileCmd() *cobra.Command`
- **USE**: `"file <file>..."`
- **SHORT**: `"Add files to the guard registry"`
- **ARGS**: Require at least one file argument
- **RUN**: Use `Run` (not RunE) that calls addFiles(args) directly
- **ERROR HANDLING**: addFiles() handles all errors internally with os.Exit(1)
- **VALIDATE**: `go build -o build/guard ./cmd/guard`

### CREATE addFiles() helper in cmd/guard/commands/add.go

- **IMPLEMENT**: Extract file adding logic into separate function that handles everything
- **SIGNATURE**: `func addFiles(args []string)` - NO error return, NO manager parameter
- **PATTERN**: Mirror runToggle() from toggle.go lines 36-105
- **LOGIC**: 
  1. Create manager: `mgr := manager.NewManager(".guardfile")`
  2. Load registry with error check and os.Exit(1)
  3. **Track counts**: Check which files are already registered BEFORE calling mgr.AddFiles
  4. Call mgr.AddFiles(args) with error check and os.Exit(1)
  5. **Track counts**: Determine newly registered vs already registered
  6. Call mgr.SaveRegistry() with error check and os.Exit(1)
  7. **Conditional output**: Only print "Registered N file(s)" if newlyRegistered > 0
  8. **Conditional output**: Only print "Skipped N file(s) already in registry" if alreadyRegistered > 0
  9. Print warnings: `manager.PrintWarnings(mgr.GetWarnings())`
  10. Print errors: `manager.PrintErrors(mgr.GetErrors())`
  11. **Exit on errors**: If mgr.HasErrors(), call os.Exit(1)
- **ERROR HANDLING**: Use `fmt.Fprintf(os.Stderr, "Error: %v\n", err)` then `os.Exit(1)`
- **VALIDATE**: `go build -o build/guard ./cmd/guard`

### UPDATE NewAddCmd() in cmd/guard/commands/add.go

- **IMPLEMENT**: Add subcommand to parent add command
- **PATTERN**: Mirror NewToggleCmd() structure
- **ADD**: `addCmd.AddCommand(newAddFileCmd())`
- **SIMPLIFY**: Parent RunE can call addFiles(args) directly for backward compatibility
- **REMOVE**: Keyword parsing logic (KeywordFile, etc.) - no longer needed
- **VALIDATE**: `go build -o build/guard ./cmd/guard && ./build/guard add --help`

### CREATE newRemoveFileCmd() in cmd/guard/commands/remove.go

- **IMPLEMENT**: Add newRemoveFileCmd() factory function
- **PATTERN**: Mirror newToggleFileCmd() from toggle.go
- **SIGNATURE**: `func newRemoveFileCmd() *cobra.Command`
- **USE**: `"file <file>..."`
- **SHORT**: `"Remove files from the guard registry"`
- **ARGS**: Require at least one file argument
- **RUN**: Use `Run` (not RunE) that calls removeFiles(args) directly
- **ERROR HANDLING**: removeFiles() handles all errors internally with os.Exit(1)
- **VALIDATE**: `go build -o build/guard ./cmd/guard`

### CREATE removeFiles() helper in cmd/guard/commands/remove.go

- **IMPLEMENT**: Extract file removal logic into separate function that handles everything
- **SIGNATURE**: `func removeFiles(args []string)` - NO error return, NO manager parameter
- **PATTERN**: Mirror runToggle() from toggle.go lines 36-105
- **LOGIC**: 
  1. Create manager: `mgr := manager.NewManager(".guardfile")`
  2. Load registry with error check and os.Exit(1)
  3. **Track counts**: Check which files are in registry BEFORE calling mgr.RemoveFiles
  4. Call mgr.RemoveFiles(args) with error check and os.Exit(1)
  5. **Track counts**: Determine inRegistry vs notInRegistry
  6. Call mgr.SaveRegistry() with error check and os.Exit(1)
  7. **Conditional output**: Only print "Removed N file(s)" if inRegistry > 0
  8. **Conditional output**: Only print "Skipped N file(s) not in registry" if notInRegistry > 0
  9. Print warnings: `manager.PrintWarnings(mgr.GetWarnings())`
  10. Print errors: `manager.PrintErrors(mgr.GetErrors())`
  11. **Exit on errors**: If mgr.HasErrors(), call os.Exit(1)
- **ERROR HANDLING**: Use `fmt.Fprintf(os.Stderr, "Error: %v\n", err)` then `os.Exit(1)`
- **VALIDATE**: `go build -o build/guard ./cmd/guard`

### UPDATE NewRemoveCmd() in cmd/guard/commands/remove.go

- **IMPLEMENT**: Add subcommand to parent remove command
- **PATTERN**: Mirror NewToggleCmd() structure
- **ADD**: `removeCmd.AddCommand(newRemoveFileCmd())`
- **SIMPLIFY**: Parent RunE can call removeFiles(args) directly for backward compatibility
- **REMOVE**: Keyword parsing logic (KeywordFile, etc.) - no longer needed
- **VALIDATE**: `go build -o build/guard ./cmd/guard && ./build/guard remove --help`

### UPDATE justfile build recipe

- **IMPLEMENT**: Add ldflags to build command for version injection
- **PATTERN**: `go build -ldflags "-X main.version=$(git describe --tags --always --dirty)" -o build/guard ./cmd/guard`
- **FALLBACK**: If no git tags, use commit hash or "dev"
- **VALIDATE**: `just build && ./build/guard version`

---

## TESTING STRATEGY

Based on project's shell integration test framework in tests/ directory.

### Unit Tests

Not required for this refactoring - existing shell tests provide behavioral coverage.

### Integration Tests

All existing shell tests must pass without modification:

**Critical Test Files:**
- `tests/test-show-commands.sh` - Validates show file/collection behavior
- `tests/test-file-management.sh` - Validates add/remove file operations
- `tests/test-help-output.sh` - Validates help text and command structure

**Test Coverage:**
- Show file with specific files
- Show file with no arguments (all files)
- Show collection with specific collections
- Show collection with no arguments (all collections)
- Add files to registry
- Remove files from registry
- Version command output
- Help command output

### Edge Cases

**Backward Compatibility:**
- `guard show file1.txt` (auto-detect) must work
- `guard show file file1.txt` (explicit) must work
- `guard add file1.txt` (implicit) must work
- `guard add file file1.txt` (explicit) must work

**Subcommand Help:**
- `guard show file --help` must display file-specific help
- `guard show collection --help` must display collection-specific help
- `guard add file --help` must display file-specific help

**Version Injection:**
- Default version "dev" when built without ldflags
- Proper version when built with ldflags
- Version displayed in `guard version` output

---

## VALIDATION COMMANDS

Execute every command to ensure zero regressions and 100% feature correctness.

### Level 1: Syntax & Style

```bash
# Format code
go fmt ./...

# Lint code
golangci-lint run

# Security scan
semgrep scan --config auto
```

### Level 2: Build Validation

```bash
# Build without ldflags (default version)
go build -o build/guard ./cmd/guard

# Build with ldflags (version injection)
go build -ldflags "-X main.version=1.0.0-test" -o build/guard ./cmd/guard

# Verify binary works
./build/guard --help
./build/guard version
```

### Level 3: Unit Tests

```bash
# Run Go unit tests
go test ./...
```

### Level 4: Integration Tests

```bash
# Copy binary to tests directory
cp build/guard tests/guard

# Run all shell integration tests
cd tests && ./run-all-tests.sh

# Run specific test suites
cd tests && ./test-show-commands.sh
cd tests && ./test-file-management.sh
cd tests && ./test-help-output.sh
```

### Level 5: Manual Validation

```bash
# Test show subcommands
./build/guard init 0644 flo staff
touch test1.txt test2.txt
./build/guard add test1.txt test2.txt
./build/guard show file test1.txt
./build/guard show file
./build/guard show --help

# Test add subcommands
./build/guard add file test3.txt
./build/guard add --help
./build/guard add file --help

# Test remove subcommands
./build/guard remove file test1.txt
./build/guard remove --help
./build/guard remove file --help

# Test version
./build/guard version

# Test backward compatibility
./build/guard show test1.txt  # Should auto-detect
./build/guard add test4.txt    # Should work without 'file' keyword
```

### Level 6: CI Pipeline

```bash
# Run full CI pipeline
just ci-quiet
```

---

## ACCEPTANCE CRITERIA

- [x] keywords.go file removed
- [x] version.go accepts version parameter
- [x] main.go uses rootCmd in main() instead of init()
- [x] show command has file and collection subcommands
- [x] add command has file subcommand
- [x] remove command has file subcommand
- [x] All subcommands have dedicated factory functions
- [x] Business logic extracted into helper functions
- [x] All validation commands pass with zero errors
- [x] All shell integration tests pass
- [x] Backward compatibility maintained
- [x] Help text displays correctly for all subcommands
- [x] Version injection works via ldflags
- [x] CI pipeline passes with exit code 0

---

## COMPLETION CHECKLIST

- [ ] keywords.go deleted
- [ ] version.go refactored with parameter
- [ ] main.go refactored without init()
- [ ] show.go has newShowFileCmd() and newShowCollectionCmd()
- [ ] show.go has showFiles() and showCollections() helpers
- [ ] add.go has newAddFileCmd() and addFiles() helper
- [ ] remove.go has newRemoveFileCmd() and removeFiles() helper
- [ ] justfile updated with ldflags
- [ ] All validation commands executed successfully
- [ ] Full test suite passes (unit + integration)
- [ ] No linting or type checking errors
- [ ] Manual testing confirms all features work
- [ ] Backward compatibility verified
- [ ] Help text verified for all subcommands
- [ ] Version injection verified
- [ ] CI pipeline passes

---

## NOTES

### Design Decisions

**Why two different helper function patterns?**
- **show.go helpers** (showAllFiles, showSpecificFiles, etc.): Void functions that print directly or use os.Exit(1) - used for read-only display operations
- **add.go/remove.go helpers** (addFiles, removeFiles): Create manager internally, use os.Exit(1), track counts for conditional output - used for write operations that modify state
- Pattern follows existing codebase: toggle.go uses os.Exit pattern for write operations
- All helpers now use void return (no error returns) for consistency

**Why keep auto-detection in parent commands?**
- Backward compatibility with existing usage patterns
- User convenience - `guard show file1.txt` is faster than `guard show file file1.txt`
- Explicit subcommands provide clarity when needed

**Why extract business logic into helper functions?**
- **Read operations** (show): Multiple helpers for different display modes (all vs specific, files vs collections)
- **Write operations** (add/remove): Helpers create manager, handle complete workflow including SaveRegistry(), track counts for conditional output
- Reusability between parent command and subcommands
- Cleaner separation of CLI concerns from business logic
- Consistent with existing toggle.go pattern for write operations
- Conditional output based on operation results (only show relevant messages)

**Why remove keywords.go?**
- Constants are only used within specific command files
- No shared validation logic needed
- Reduces unnecessary file count

**Why move rootCmd to main()?**
- Better testability - can create rootCmd in tests
- Follows Go best practices - avoid package-level state
- Enables dependency injection for testing

**Why use ldflags for version?**
- Industry standard for Go CLI tools
- Enables dynamic version injection at build time
- No need to update code for version changes
- Supports git-based versioning

### Trade-offs

**Increased file size vs modularity:**
- Commands will have more functions (factory + helper)
- Trade-off: Better organization and reusability outweigh file size

**Backward compatibility complexity:**
- Parent commands need to support both auto-detect and explicit subcommands
- Trade-off: User convenience and migration path worth the complexity

**Version injection complexity:**
- Requires build-time flag configuration
- Trade-off: Professional version management worth the setup

### Migration Path

**For users:**
- No breaking changes - all existing commands work
- New subcommands provide explicit alternatives
- Help text guides users to new patterns

**For developers:**
- Clear factory pattern for adding new subcommands
- Helper functions show business logic extraction pattern
- main.go structure enables better testing

### Future Enhancements

**Potential improvements after this refactoring:**
- Add more subcommands (e.g., `show folder`, `add collection`)
- Extract common patterns into shared utilities
- Add command-specific flags to subcommands
- Implement command aliases for common operations
- Add shell completion for subcommands

### References

**Cobra Documentation:**
- Content was rephrased for compliance with licensing restrictions
- Cobra organizes applications around commands, arguments, and flags in a tree structure
- Parent commands can have multiple subcommands with their own flags and arguments
- Factory pattern recommended for modular command organization

**Go Build Documentation:**
- Content was rephrased for compliance with licensing restrictions
- The -ldflags flag passes arguments to the Go linker
- The -X option sets string variable values at link time
- Format: `-ldflags "-X package.variable=value"`

**Project Conventions:**
- Follow existing patterns in enable.go, disable.go, toggle.go
- Maintain error handling style from manager layer
- Use fmt.Errorf with %w for error wrapping
- Delegate business logic to manager layer

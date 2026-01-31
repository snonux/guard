# Feature: Refactor Manager Package Architecture

## Feature Description

Comprehensive refactoring of the manager package to eliminate architectural issues including count return values, embedded printing, code duplication between silent/non-silent methods, missing bulk operations for collections, and minimal folder functionality. This refactoring will create a clean separation between business logic (manager layer) and presentation logic (CLI layer).

## User Story

As a developer maintaining the guard-tool codebase
I want a clean manager layer that focuses purely on business logic
So that the CLI commands can handle all presentation concerns and the codebase is more maintainable

## Problem Statement

The current manager package has several architectural issues:
1. **Mixed Responsibilities**: Manager methods return count integers that should be handled by CLI layer
2. **Embedded Output**: Manager layer contains fmt.Printf statements that violate separation of concerns
3. **Code Duplication**: Silent and non-silent method variants create maintenance overhead
4. **Missing Bulk Operations**: Collections lack proper bulk operations and conflict detection
5. **Minimal Folder Support**: Folders.go lacks essential methods and state management

## Solution Statement

Refactor the manager package to:
- Return only errors from manager methods (no counts or output)
- Move all output formatting to CLI commands
- Eliminate silent/non-silent method duplication
- Add comprehensive bulk operations for collections with conflict detection
- Enhance folder functionality with proper state management and path normalization

## Feature Metadata

**Feature Type**: Refactor
**Estimated Complexity**: High
**Primary Systems Affected**: Manager package, CLI commands
**Dependencies**: None (internal refactoring)

---

## CONTEXT REFERENCES

### Relevant Codebase Files IMPORTANT: YOU MUST READ THESE FILES BEFORE IMPLEMENTING!

- `internal/manager/files.go` - Why: Contains current file operations with count returns and printing
- `internal/manager/collections.go` - Why: Shows current collection operations and duplication patterns
- `internal/manager/folders.go` - Why: Minimal implementation that needs enhancement
- `internal/manager/warnings.go` - Why: Warning system that should be preserved and enhanced
- `internal/manager/manager.go` - Why: Core manager structure and orchestration methods
- `internal/registry/registry.go` (lines 1-50) - Why: Registry interface patterns to follow
- `internal/registry/folder_entry.go` - Why: Folder data structures and operations
- `cmd/guard/commands/enable.go` (lines 1-50) - Why: CLI output patterns to follow

### New Files to Create

None - all enums and types will be added to existing files

### Relevant Documentation YOU SHOULD READ THESE BEFORE IMPLEMENTING!

- `COBRA_BEST_PRACTICES.md` - CLI command patterns and output formatting
- `.kiro/steering/tech.md` - Go coding standards and error handling patterns
- `.kiro/steering/structure.md` - Package organization and naming conventions

### Patterns to Follow

**Error Handling Pattern:**
```go
// Manager methods return only error, no counts
func (m *Manager) AddFiles(filePaths []string) error {
    // Business logic only
    return nil
}
```

**CLI Output Pattern:**
```go
// CLI commands handle counting and output
func runAddFiles(cmd *cobra.Command, args []string) error {
    count := 0
    for _, file := range args {
        if err := mgr.AddFile(file); err != nil {
            return err
        }
        count++
    }
    fmt.Printf("Added %d files\n", count)
    return nil
}
```

**Warning Aggregation Pattern:**
```go
// Manager accumulates warnings, CLI prints them
mgr.AddFiles(files)
if mgr.HasWarnings() {
    mgr.PrintWarnings()
}
```

---

## IMPLEMENTATION PLAN

### Phase 1: Foundation - Add Enums to Existing Files

**Tasks:**
- Add EffectiveFolderGuardState enum to folders.go
- Add CollectionConflict type to collections.go
- Establish consistent error patterns

### Phase 2: Refactor Files Operations

**Tasks:**
- Remove count returns from all file methods
- Eliminate printing from manager layer
- Consolidate silent/non-silent method variants
- Add ShowFiles, Cleanup, Reset, Destroy methods
- Update CLI commands to handle counting and output

### Phase 3: Enhance Collections Operations

**Tasks:**
- Add missing bulk operations (ToggleCollections, EnableCollections, DisableCollections)
- Add file management operations (AddFilesToCollections, RemoveFilesFromCollections)
- Add collection management operations (AddCollectionsToCollections, RemoveCollectionsFromCollections)
- Add ShowCollections method
- Implement conflict detection directly in ToggleCollections
- Refactor existing collection methods to remove counts/printing
- Update CLI commands for new collection operations

### Phase 4: Expand Folder Functionality

**Tasks:**
- Add EffectiveFolderGuardState enum with all 5 states
- Add comprehensive folder operations (ToggleFolders, EnableFolders, DisableFolders)
- Implement path normalization helpers
- Create effective folder guard state calculation
- Add folder scanning and validation methods

### Phase 5: Update CLI Integration

**Tasks:**
- Update all CLI commands to handle new manager interface
- Implement proper output formatting in CLI layer
- Add progress reporting for bulk operations
- Ensure warning system integration

---

## STEP-BY-STEP TASKS

### ADD EffectiveFolderGuardState enum to internal/manager/folders.go

- **IMPLEMENT**: Complete folder state enum with all 5 states plus package-level path functions
- **PATTERN**: Follow existing warning types pattern in warnings.go
- **IMPORTS**: os, path/filepath for path normalization
- **GOTCHA**: Use iota for enum values, handle inherited guard state correctly
- **VALIDATE**: `go build ./internal/manager`

```go
// EffectiveFolderGuardState represents the effective guard state of a folder
type EffectiveFolderGuardState int

const (
    FolderNotRegistered EffectiveFolderGuardState = iota
    FolderAllGuarded      // Folder guard=true AND all files guarded
    FolderInheritedGuard  // Folder guard=false BUT all files guarded
    FolderMixedState      // Some files guarded, some unguarded
    FolderAllUnguarded    // All files unguarded
)

// Package-level path normalization functions
func normalizeFolderPath(path string) (string, error)
func folderNameFromPath(path string) string
```

### REFACTOR internal/manager/files.go

- **IMPLEMENT**: Remove all count returns, add ShowFiles, Cleanup, Reset, Destroy methods with result structs
- **PATTERN**: Mirror security layer error-only pattern from registry.go
- **IMPORTS**: Keep existing imports, remove fmt for printing
- **GOTCHA**: Preserve warning accumulation, remove all fmt.Printf calls
- **VALIDATE**: `go build ./internal/manager && go test ./internal/manager`

Key changes:
- `AddFiles([]string) (int, error)` → `AddFiles([]string) error`
- `EnableFiles([]string) (int, error)` → `EnableFiles([]string) error`
- Remove `setFileProtectionSilent` - merge with `setFileProtection`
- Remove all `fmt.Printf` statements
- Add `FileInfo` struct with path, guard status, collections list
- Add `ShowFiles() ([]FileInfo, error)` with collection membership
- Add `CleanupResult` struct with counts
- Add `Cleanup() (CleanupResult, error)` to remove stale entries
- Add `ResetResult` struct with counts
- Add `Reset() (ResetResult, error)` to disable all guards
- Add `Destroy() error` for reset + cleanup + delete guardfile

### REFACTOR internal/manager/collections.go

- **IMPLEMENT**: Remove count returns, add comprehensive collection operations
- **PATTERN**: Follow files.go refactored pattern
- **IMPORTS**: Remove fmt, keep security import
- **GOTCHA**: Implement conflict detection directly in ToggleCollections
- **VALIDATE**: `go build ./internal/manager && go test ./internal/manager`

New methods to add:
```go
func (m *Manager) AddCollections(names []string) error // create empty collections
func (m *Manager) ToggleCollections(names []string) error // with built-in conflict detection
func (m *Manager) EnableCollections(names []string) error
func (m *Manager) DisableCollections(names []string) error
func (m *Manager) AddFilesToCollections(filePaths []string, collectionNames []string) error
func (m *Manager) RemoveFilesFromCollections(filePaths []string, collectionNames []string) error
func (m *Manager) AddCollectionsToCollections(sourceNames []string, targetNames []string) error
func (m *Manager) RemoveCollectionsFromCollections(sourceNames []string, targetNames []string) error
func (m *Manager) ShowCollections() ([]CollectionInfo, error)
```

### REFACTOR internal/manager/folders.go

- **IMPLEMENT**: Add EffectiveFolderGuardState enum and comprehensive folder operations
- **PATTERN**: Mirror collections.go structure for consistency
- **IMPORTS**: Add os, path/filepath for path normalization
- **GOTCHA**: Remove printing, handle inherited guard state correctly
- **VALIDATE**: `go build ./internal/manager && go test ./internal/manager`

New methods to add:
```go
func (m *Manager) ToggleFolders(folderPaths []string) error
func (m *Manager) EnableFolders(folderPaths []string) error  
func (m *Manager) DisableFolders(folderPaths []string) error
func (m *Manager) NormalizeFolderPath(path string) (string, error)
func (m *Manager) GetEffectiveFolderGuardState(folderName string) (EffectiveFolderGuardState, error)
func (m *Manager) ScanFolderFiles(folderPath string, recursive bool) ([]string, error)
```

### UPDATE cmd/guard/commands/add.go

- **IMPLEMENT**: Handle counting and output formatting for file additions
- **PATTERN**: Follow enable.go CLI output pattern
- **IMPORTS**: Add fmt for output formatting
- **GOTCHA**: Count successful operations, handle warnings properly
- **VALIDATE**: `go build ./cmd/guard && ./guard add --help`

```go
func runAdd(cmd *cobra.Command, args []string) error {
    // Resolve arguments using manager
    files, folders, collections, err := mgr.ResolveArguments(args)
    if err != nil {
        return err
    }
    
    // Process each type and count successes
    totalAdded := 0
    if len(files) > 0 {
        if err := mgr.AddFiles(files); err != nil {
            return err
        }
        totalAdded += len(files)
    }
    
    // Print results
    fmt.Printf("Added %d items to registry\n", totalAdded)
    
    // Print warnings if any
    if mgr.HasWarnings() {
        mgr.PrintWarnings()
    }
    
    return nil
}
```

### UPDATE cmd/guard/commands/enable.go

- **IMPLEMENT**: Update to use new manager interface without count returns
- **PATTERN**: Follow existing CLI output patterns in enable.go
- **IMPORTS**: Keep existing imports
- **GOTCHA**: Handle bulk operations properly, maintain auto-detection
- **VALIDATE**: `go build ./cmd/guard && ./guard enable --help`

### UPDATE cmd/guard/commands/disable.go

- **IMPLEMENT**: Update to use new manager interface
- **PATTERN**: Mirror enable.go updates
- **IMPORTS**: Keep existing imports  
- **GOTCHA**: Ensure consistent output formatting with enable command
- **VALIDATE**: `go build ./cmd/guard && ./guard disable --help`

### UPDATE cmd/guard/commands/toggle.go

- **IMPLEMENT**: Update to use new manager interface
- **PATTERN**: Mirror enable.go updates
- **IMPORTS**: Keep existing imports
- **GOTCHA**: Handle mixed states properly for collections/folders
- **VALIDATE**: `go build ./cmd/guard && ./guard toggle --help`

### UPDATE cmd/guard/commands/create.go

- **IMPLEMENT**: Update collection creation to use new bulk operations
- **PATTERN**: Follow add.go updated pattern
- **IMPORTS**: Keep existing imports
- **GOTCHA**: Use new AddCollections method for bulk creation
- **VALIDATE**: `go build ./cmd/guard && ./guard create --help`

### UPDATE cmd/guard/commands/update.go

- **IMPLEMENT**: Add conflict detection for collection updates
- **PATTERN**: Follow existing update command structure
- **IMPORTS**: Keep existing imports
- **GOTCHA**: Check for conflicts before applying updates
- **VALIDATE**: `go build ./cmd/guard && ./guard update --help`

### UPDATE cmd/guard/commands/clear.go

- **IMPLEMENT**: Update to use new bulk collection operations
- **PATTERN**: Follow destroy.go pattern for bulk operations
- **IMPORTS**: Keep existing imports
- **GOTCHA**: Use new ClearCollections method
- **VALIDATE**: `go build ./cmd/guard && ./guard clear --help`

### UPDATE cmd/guard/commands/destroy.go

- **IMPLEMENT**: Update to use new bulk collection operations  
- **PATTERN**: Follow clear.go updated pattern
- **IMPORTS**: Keep existing imports
- **GOTCHA**: Use new DestroyCollections method
- **VALIDATE**: `go build ./cmd/guard && ./guard destroy --help`

---

## TESTING STRATEGY

### Unit Tests

Design unit tests following existing Go testing patterns in the project:
- Test each manager method in isolation
- Mock filesystem operations where needed
- Verify warning accumulation works correctly
- Test error propagation from registry layer

### Integration Tests

Update existing shell integration tests:
- Verify CLI output format remains consistent
- Test bulk operations work correctly
- Validate conflict detection functionality
- Ensure folder operations work as expected

### Edge Cases

Test specific edge cases for this refactoring:
- Empty collections and folders
- Conflicting guard states in collections
- Missing files during bulk operations
- Path normalization edge cases
- Warning aggregation with large numbers of items

---

## VALIDATION COMMANDS

Execute every command to ensure zero regressions and 100% feature correctness.

### Level 1: Syntax & Style

```bash
go fmt ./...
golangci-lint run
```

### Level 2: Unit Tests

```bash
go test ./internal/manager -v
go test ./cmd/guard/commands -v
```

### Level 3: Integration Tests

```bash
go build -o build/guard ./cmd/guard
cp build/guard tests/guard
cd tests && ./run-all-tests.sh
```

### Level 4: Manual Validation

```bash
# Test basic operations still work
./guard init 0644 testuser testgroup
./guard add file1.txt file2.txt
./guard enable file1.txt
./guard show

# Test new bulk operations
./guard create collection1 collection2
./guard update collection1 add file1.txt file2.txt
./guard enable collection1 collection2

# Test folder operations
./guard enable docs/
./guard toggle docs/
./guard show
```

### Level 5: Additional Validation

```bash
# Verify no output changes for existing commands
./guard --help
./guard enable --help
./guard show
```

---

## ACCEPTANCE CRITERIA

- [ ] All manager methods return only errors (no count integers)
- [ ] No fmt.Printf statements remain in manager package
- [ ] Silent/non-silent method duplication eliminated
- [ ] Bulk collection operations implemented (Toggle, Enable, Disable)
- [ ] Collection file management operations implemented (AddFilesToCollections, RemoveFilesFromCollections)
- [ ] Collection-to-collection operations implemented (AddCollectionsToCollections, RemoveCollectionsFromCollections)
- [ ] Collection conflict detection built into ToggleCollections
- [ ] Folder operations enhanced (Toggle, Enable, Disable, path normalization)
- [ ] EffectiveFolderGuardState enum implemented with all 5 states including inherited guard
- [ ] ShowFiles, ShowCollections methods implemented
- [ ] Cleanup, Reset, Destroy methods implemented in files.go
- [ ] All CLI commands updated to handle new manager interface
- [ ] CLI commands provide proper counting and output formatting
- [ ] Warning system preserved and working correctly
- [ ] All existing shell tests pass without modification
- [ ] No regressions in existing functionality
- [ ] Code follows project Go standards and conventions

---

## COMPLETION CHECKLIST

- [ ] All tasks completed in dependency order
- [ ] Each task validation passed immediately
- [ ] All validation commands executed successfully
- [ ] Full test suite passes (unit + integration)
- [ ] No linting or type checking errors
- [ ] Manual testing confirms all operations work
- [ ] Shell integration tests pass unchanged
- [ ] Acceptance criteria all met
- [ ] Code reviewed for quality and maintainability

---

## NOTES

**Key Design Decisions:**
- Manager layer becomes pure business logic with error-only returns
- CLI layer handles all counting, formatting, and user output
- Warning system preserved but enhanced for better aggregation
- Folder state management added for better UX
- Collection conflict detection prevents data inconsistencies

**Migration Strategy:**
- Refactor manager methods first to establish new interface
- Update CLI commands to match new interface
- Preserve existing shell test compatibility
- Maintain backward compatibility for .guardfile format

**Performance Considerations:**
- Bulk operations reduce registry save overhead
- Path normalization cached where possible
- Folder scanning optimized for common use cases
- Warning aggregation prevents output spam

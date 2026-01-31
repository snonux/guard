# Feature: Refactor Warnings System with Typed Warnings

The following plan should be complete, but it's important that you validate documentation and codebase patterns and task sanity before you start implementing.

Pay special attention to naming of existing utils types and models. Import from the right files etc.

## Feature Description

Refactor the manager layer's warning system from simple string-based warnings to a strongly-typed warning system with proper categorization, aggregation, and display formatting. This will enable better warning management, grouping of similar warnings, and cleaner output for users.

## User Story

As a developer maintaining the guard-tool codebase
I want a typed warning system with proper categorization
So that warnings can be aggregated, formatted consistently, and managed more effectively

## Problem Statement

The current warning system in the manager layer uses simple string slices (`[]string`) which:
- Lacks type safety and categorization
- Cannot aggregate similar warnings efficiently
- Provides no structured way to associate warnings with affected items (file paths, collection names)
- Mixes warnings and errors in an ad-hoc manner
- Makes it difficult to format warnings consistently across commands

## Solution Statement

Implement a strongly-typed warning system with:
- `WarningType` enum using iota for type-safe warning categories
- `Warning` struct containing type, message, and associated items
- Aggregation function to group warnings by type and format with bullet points
- Separate tracking for warnings vs errors in the Manager
- Helper methods for adding, retrieving, and displaying warnings/errors

## Feature Metadata

**Feature Type**: Refactor
**Estimated Complexity**: Medium
**Primary Systems Affected**: internal/manager (warnings.go, manager.go, files.go, collections.go, folders.go)
**Dependencies**: None (internal refactor only)

---

## CONTEXT REFERENCES

### Relevant Codebase Files IMPORTANT: YOU MUST READ THESE FILES BEFORE IMPLEMENTING!

- `internal/manager/warnings.go` (all lines) - Why: Current warning system implementation to be refactored
- `internal/manager/manager.go` (lines 26, 34, 348, 364, 434-447) - Why: Manager struct warnings field and GetWarnings/AddWarning methods
- `internal/manager/files.go` (line 26) - Why: Example of AddWarning usage pattern
- `internal/manager/collections.go` (lines 86, 128, 137, 141, 262) - Why: Multiple AddWarning usage patterns
- `internal/manager/folders.go` (lines 37, 41, 50, 54) - Why: AddWarning usage in folder operations
- `cmd/guard/commands/toggle.go` (lines 81, 110, 138-142) - Why: GetWarnings usage and printWarnings helper function
- `internal/manager/manager.go` (lines 14-18) - Why: Const pattern for target types (reference for enum style)
- `internal/filesystem/filesystem_linux.go` (lines 15-19) - Why: Const pattern without iota (reference)

### New Files to Create

None - all changes are modifications to existing files

### Relevant Documentation YOU SHOULD READ THESE BEFORE IMPLEMENTING!

- `.kiro/steering/tech.md` (lines 20-30) - Code Standards section
  - Why: Go naming conventions, error handling, const usage with iota
- `.kiro/steering/tech.md` (lines 31-38) - Error Handling Standards
  - Why: Error wrapping patterns and return value conventions

### Patterns to Follow

**Enum Pattern with iota:**
```go
// From internal/manager/manager.go
const (
	TargetTypeFile       = "file"
	TargetTypeFolder     = "folder"
	TargetTypeCollection = "collection"
)
```

**Struct with Constructor:**
```go
// From internal/manager/manager.go
type Manager struct {
	mu           sync.RWMutex
	registryPath string
	security     *security.Security
	fs           filesystem.Filesystem
	warnings     []string  // This will change to []Warning
}

func NewManager(registryPath string) *Manager {
	return &Manager{
		registryPath: registryPath,
		fs:           filesystem.NewFilesystem(),
		warnings:     make([]string, 0),  // This will change
	}
}
```

**Method Naming Convention:**
```go
// From internal/manager/manager.go
func (m *Manager) GetWarnings() []string { ... }
func (m *Manager) AddWarning(message string) { ... }
```

**Error Output Pattern:**
```go
// From cmd/guard/commands/toggle.go
func printWarnings(warnings []string) {
	for _, warning := range warnings {
		fmt.Fprintf(os.Stderr, "Warning: %s\n", warning)
	}
}
```

**Mutex Usage for Thread Safety:**
```go
// From internal/manager/manager.go
func (m *Manager) AddWarning(message string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.warnings = append(m.warnings, message)
}
```

---

## IMPLEMENTATION PLAN

### Phase 1: Foundation - Define Warning Types and Structures

Create the typed warning system foundation with enums, structs, and type-specific aggregation functions.

**Tasks:**
- Define WarningType enum with iota (including silent WarningFileAlreadyInRegistry)
- Create Warning struct with Type, Message, and Items fields
- Implement NewWarning constructor
- Implement type-specific aggregation functions:
  - aggregateFileMissing (with cleanup suggestion)
  - aggregateCollectionEmpty
  - aggregateFolderEmpty
  - aggregateGeneric (fallback for other types)
- Implement AggregateWarnings dispatcher function
- Implement PrintWarnings (takes []Warning, outputs to stdout)
- Implement PrintErrors (takes []string, outputs to stderr)

### Phase 2: Core Implementation - Update Manager

Update the Manager struct to use the new typed warning system and add error tracking.

**Tasks:**
- Change Manager.warnings field from []string to []Warning
- Add Manager.errors field as []string
- Update AddWarning to accept Warning type
- Add AddError method
- Update GetWarnings to return aggregated string slice
- Add GetErrors, HasWarnings, HasErrors, ClearWarnings, ClearErrors methods
- Add PrintWarnings and PrintErrors functions

### Phase 3: Integration - Update Warning Call Sites

Update all existing AddWarning calls throughout the manager package to use the new typed system.

**Tasks:**
- Update files.go AddWarning calls
- Update collections.go AddWarning calls
- Update folders.go AddWarning calls
- Update manager.go Cleanup method AddWarning calls
- Verify all call sites are migrated

### Phase 4: Testing & Validation

Ensure the refactor maintains existing behavior and passes all tests.

**Tasks:**
- Run unit tests
- Run shell integration tests
- Verify warning output format matches expectations
- Validate thread safety with mutex usage

---

## STEP-BY-STEP TASKS

IMPORTANT: Execute every task in order, top to bottom. Each task is atomic and independently testable.

### UPDATE internal/manager/warnings.go

- **REMOVE**: Entire existing WarningCollector implementation (lines 15-82)
- **REMOVE**: All standalone warning helper functions (WarnFileNotExist, WarnCollectionEmpty, etc.)
- **KEEP**: Warning type constants at top (lines 6-13) - these will be replaced with iota enum
- **IMPLEMENT**: WarningType enum using iota
  ```go
  type WarningType int
  
  const (
  	WarningFileMissing WarningType = iota
  	WarningFileNotInRegistry
  	WarningFileAlreadyInRegistry  // Silent - produces no output
  	WarningCollectionEmpty
  	WarningCollectionNotFound
  	WarningCollectionAlreadyExists
  	WarningFileNotInCollection
  	WarningCollectionHasMissingFiles
  	WarningCollectionCreated
  	WarningFolderEmpty
  	WarningFileAlreadyGuarded
  	WarningGeneric
  )
  ```
- **IMPLEMENT**: Warning struct
  ```go
  type Warning struct {
  	Type    WarningType
  	Message string
  	Items   []string  // Associated file paths, collection names, etc.
  }
  ```
- **IMPLEMENT**: NewWarning constructor
  ```go
  func NewWarning(warnType WarningType, message string, items ...string) Warning {
  	return Warning{
  		Type:    warnType,
  		Message: message,
  		Items:   items,
  	}
  }
  ```
- **IMPLEMENT**: Type-specific aggregation functions
  ```go
  func aggregateFileMissing(warnings []Warning) []string {
  	if len(warnings) == 0 {
  		return nil
  	}
  	
  	var items []string
  	for _, w := range warnings {
  		items = append(items, w.Items...)
  	}
  	
  	if len(items) == 0 {
  		return nil
  	}
  	
  	result := []string{"The following files do not exist on disk:"}
  	for _, item := range items {
  		result = append(result, fmt.Sprintf("  - %s", item))
  	}
  	result = append(result, "Suggestion: Run 'guard cleanup' to remove missing files from registry")
  	return result
  }
  
  func aggregateCollectionEmpty(warnings []Warning) []string {
  	if len(warnings) == 0 {
  		return nil
  	}
  	
  	var items []string
  	for _, w := range warnings {
  		items = append(items, w.Items...)
  	}
  	
  	if len(items) == 0 {
  		return nil
  	}
  	
  	if len(items) == 1 {
  		return []string{fmt.Sprintf("Collection '%s' is empty", items[0])}
  	}
  	
  	result := []string{"The following collections are empty:"}
  	for _, item := range items {
  		result = append(result, fmt.Sprintf("  - %s", item))
  	}
  	return result
  }
  
  func aggregateFolderEmpty(warnings []Warning) []string {
  	if len(warnings) == 0 {
  		return nil
  	}
  	
  	var items []string
  	for _, w := range warnings {
  		items = append(items, w.Items...)
  	}
  	
  	if len(items) == 0 {
  		return nil
  	}
  	
  	if len(items) == 1 {
  		return []string{fmt.Sprintf("Folder '%s' contains no files", items[0])}
  	}
  	
  	result := []string{"The following folders contain no files:"}
  	for _, item := range items {
  		result = append(result, fmt.Sprintf("  - %s", item))
  	}
  	return result
  }
  
  func aggregateGeneric(warnings []Warning) []string {
  	var result []string
  	for _, w := range warnings {
  		if len(w.Items) == 0 {
  			result = append(result, w.Message)
  		} else if len(w.Items) == 1 {
  			result = append(result, fmt.Sprintf("%s: %s", w.Message, w.Items[0]))
  		} else {
  			result = append(result, w.Message+":")
  			for _, item := range w.Items {
  				result = append(result, fmt.Sprintf("  - %s", item))
  			}
  		}
  	}
  	return result
  }
  ```
- **IMPLEMENT**: AggregateWarnings function with type-specific dispatch
  ```go
  func AggregateWarnings(warnings []Warning) []string {
  	if len(warnings) == 0 {
  		return nil
  	}
  	
  	// Group warnings by type
  	grouped := make(map[WarningType][]Warning)
  	for _, w := range warnings {
  		grouped[w.Type] = append(grouped[w.Type], w)
  	}
  	
  	var result []string
  	
  	// Process each warning type with its specific aggregator
  	if warns, ok := grouped[WarningFileMissing]; ok {
  		result = append(result, aggregateFileMissing(warns)...)
  	}
  	
  	if warns, ok := grouped[WarningCollectionEmpty]; ok {
  		result = append(result, aggregateCollectionEmpty(warns)...)
  	}
  	
  	if warns, ok := grouped[WarningFolderEmpty]; ok {
  		result = append(result, aggregateFolderEmpty(warns)...)
  	}
  	
  	// WarningFileAlreadyInRegistry is silent - skip it
  	
  	// Handle all other warning types with generic aggregator
  	for warnType, warns := range grouped {
  		if warnType == WarningFileMissing || 
  		   warnType == WarningCollectionEmpty || 
  		   warnType == WarningFolderEmpty ||
  		   warnType == WarningFileAlreadyInRegistry {
  			continue
  		}
  		result = append(result, aggregateGeneric(warns)...)
  	}
  	
  	return result
  }
  ```
- **PATTERN**: Follow Manager struct pattern from manager.go
- **IMPORTS**: Add "fmt" if not already present
- **GOTCHA**: WarningFileAlreadyInRegistry is intentionally skipped in AggregateWarnings (silent)
- **VALIDATE**: `go build ./internal/manager`

### UPDATE internal/manager/manager.go - Manager struct and constructor

- **UPDATE**: Manager struct (line 26)
  ```go
  type Manager struct {
  	mu           sync.RWMutex
  	registryPath string
  	security     *security.Security
  	fs           filesystem.Filesystem
  	warnings     []Warning  // Changed from []string
  	errors       []string   // New field
  }
  ```
- **UPDATE**: NewManager constructor (line 34)
  ```go
  func NewManager(registryPath string) *Manager {
  	return &Manager{
  		registryPath: registryPath,
  		fs:           filesystem.NewFilesystem(),
  		warnings:     make([]Warning, 0),  // Changed type
  		errors:       make([]string, 0),   // New field
  	}
  }
  ```
- **PATTERN**: Mirror existing Manager patterns
- **GOTCHA**: Don't forget to initialize both warnings and errors slices
- **VALIDATE**: `go build ./internal/manager`

### UPDATE internal/manager/manager.go - Warning/Error methods

- **UPDATE**: AddWarning method (around line 447)
  ```go
  func (m *Manager) AddWarning(warning Warning) {
  	m.mu.Lock()
  	defer m.mu.Unlock()
  	m.warnings = append(m.warnings, warning)
  }
  ```
- **ADD**: AddError method (after AddWarning)
  ```go
  func (m *Manager) AddError(message string) {
  	m.mu.Lock()
  	defer m.mu.Unlock()
  	m.errors = append(m.errors, message)
  }
  ```
- **UPDATE**: GetWarnings method (line 434)
  ```go
  func (m *Manager) GetWarnings() []string {
  	m.mu.Lock()
  	defer m.mu.Unlock()
  	
  	aggregated := AggregateWarnings(m.warnings)
  	m.warnings = m.warnings[:0]  // Clear after retrieving
  	return aggregated
  }
  ```
- **ADD**: GetErrors method
  ```go
  func (m *Manager) GetErrors() []string {
  	m.mu.Lock()
  	defer m.mu.Unlock()
  	
  	errors := make([]string, len(m.errors))
  	copy(errors, m.errors)
  	m.errors = m.errors[:0]  // Clear after retrieving
  	return errors
  }
  ```
- **ADD**: HasWarnings method
  ```go
  func (m *Manager) HasWarnings() bool {
  	m.mu.RLock()
  	defer m.mu.RUnlock()
  	return len(m.warnings) > 0
  }
  ```
- **ADD**: HasErrors method
  ```go
  func (m *Manager) HasErrors() bool {
  	m.mu.RLock()
  	defer m.mu.RUnlock()
  	return len(m.errors) > 0
  }
  ```
- **ADD**: ClearWarnings method
  ```go
  func (m *Manager) ClearWarnings() {
  	m.mu.Lock()
  	defer m.mu.Unlock()
  	m.warnings = m.warnings[:0]
  }
  ```
- **ADD**: ClearErrors method
  ```go
  func (m *Manager) ClearErrors() {
  	m.mu.Lock()
  	defer m.mu.Unlock()
  	m.errors = m.errors[:0]
  }
  ```
- **PATTERN**: Use RLock for read-only operations (Has* methods), Lock for mutations
- **GOTCHA**: GetWarnings now calls AggregateWarnings before returning
- **VALIDATE**: `go build ./internal/manager`

### ADD internal/manager/warnings.go - Print functions

- **ADD**: PrintWarnings function at end of file
  ```go
  func PrintWarnings(warnings []Warning) {
  	aggregated := AggregateWarnings(warnings)
  	for _, warning := range aggregated {
  		fmt.Printf("Warning: %s\n", warning)
  	}
  }
  ```
- **ADD**: PrintErrors function at end of file
  ```go
  func PrintErrors(errors []string) {
  	for _, err := range errors {
  		fmt.Fprintf(os.Stderr, "Error: %s\n", err)
  	}
  }
  ```
- **PATTERN**: PrintWarnings outputs to stdout, PrintErrors to stderr
- **IMPORTS**: Add "os" and "fmt" imports
- **GOTCHA**: PrintWarnings takes []Warning and calls AggregateWarnings internally
- **VALIDATE**: `go build ./internal/manager`

### UPDATE internal/manager/files.go - Migrate AddWarning calls

- **UPDATE**: Line 26 AddWarning call
  ```go
  // Old: m.AddWarning(fmt.Sprintf("file does not exist: %s", cleanPath))
  m.AddWarning(NewWarning(WarningFileMissing, "file does not exist", cleanPath))
  ```
- **PATTERN**: Use NewWarning constructor with appropriate WarningType
- **GOTCHA**: Pass file path as third argument (items parameter)
- **VALIDATE**: `go build ./internal/manager`

### UPDATE internal/manager/collections.go - Migrate AddWarning calls

- **UPDATE**: Line 86
  ```go
  // Old: m.AddWarning(fmt.Sprintf("failed to protect file %s: %v", filePath, err))
  m.AddWarning(NewWarning(WarningGeneric, fmt.Sprintf("failed to protect file: %v", err), filePath))
  ```
- **UPDATE**: Line 128
  ```go
  // Old: m.AddWarning(fmt.Sprintf("collection is empty: %s", collectionName))
  m.AddWarning(NewWarning(WarningCollectionEmpty, "collection is empty", collectionName))
  ```
- **UPDATE**: Line 137
  ```go
  // Old: m.AddWarning(fmt.Sprintf("failed to protect file %s: %v", filePath, err))
  m.AddWarning(NewWarning(WarningGeneric, fmt.Sprintf("failed to protect file: %v", err), filePath))
  ```
- **UPDATE**: Line 141
  ```go
  // Old: m.AddWarning(fmt.Sprintf("failed to update guard flag for %s: %v", filePath, err))
  m.AddWarning(NewWarning(WarningGeneric, fmt.Sprintf("failed to update guard flag: %v", err), filePath))
  ```
- **UPDATE**: Line 262
  ```go
  // Old: m.AddWarning(fmt.Sprintf("failed to %s protection for %s: %v", action, filePath, err))
  m.AddWarning(NewWarning(WarningGeneric, fmt.Sprintf("failed to %s protection: %v", action, err), filePath))
  ```
- **PATTERN**: Use WarningGeneric for error-related warnings, specific types for known conditions
- **GOTCHA**: File paths go in Items parameter, not in Message string
- **VALIDATE**: `go build ./internal/manager`

### UPDATE internal/manager/folders.go - Migrate AddWarning calls

- **UPDATE**: Line 37
  ```go
  // Old: m.AddWarning(fmt.Sprintf("failed to get file info for %s: %v", filePath, err))
  m.AddWarning(NewWarning(WarningGeneric, fmt.Sprintf("failed to get file info: %v", err), filePath))
  ```
- **UPDATE**: Line 41
  ```go
  // Old: m.AddWarning(fmt.Sprintf("failed to register file %s: %v", filePath, err))
  m.AddWarning(NewWarning(WarningGeneric, fmt.Sprintf("failed to register file: %v", err), filePath))
  ```
- **UPDATE**: Line 50
  ```go
  // Old: m.AddWarning(fmt.Sprintf("failed to set protection for %s: %v", filePath, err))
  m.AddWarning(NewWarning(WarningGeneric, fmt.Sprintf("failed to set protection: %v", err), filePath))
  ```
- **UPDATE**: Line 54
  ```go
  // Old: m.AddWarning(fmt.Sprintf("failed to update guard flag for %s: %v", filePath, err))
  m.AddWarning(NewWarning(WarningGeneric, fmt.Sprintf("failed to update guard flag: %v", err), filePath))
  ```
- **PATTERN**: Consistent with collections.go migration
- **VALIDATE**: `go build ./internal/manager`

### UPDATE internal/manager/manager.go - Cleanup method warnings

- **UPDATE**: Line 348
  ```go
  // Old: m.AddWarning(fmt.Sprintf("failed to remove missing file %s: %v", filePath, err))
  m.AddWarning(NewWarning(WarningGeneric, fmt.Sprintf("failed to remove missing file: %v", err), filePath))
  ```
- **UPDATE**: Line 364
  ```go
  // Old: m.AddWarning(fmt.Sprintf("failed to remove empty collection %s: %v", collectionName, err))
  m.AddWarning(NewWarning(WarningGeneric, fmt.Sprintf("failed to remove empty collection: %v", err), collectionName))
  ```
- **PATTERN**: Same as other migrations
- **VALIDATE**: `go build ./internal/manager`

### VERIFY - Build and test entire manager package

- **VERIFY**: All files compile without errors
- **VERIFY**: No remaining string-based AddWarning calls
- **COMMAND**: `go build ./internal/manager`
- **COMMAND**: `go build ./cmd/guard`
- **GOTCHA**: Ensure no compilation errors before proceeding to tests
- **VALIDATE**: `go build ./...`

---

## TESTING STRATEGY

### Unit Tests

No new unit tests required - this is a refactor maintaining existing behavior. Existing tests should continue to pass without modification.

### Integration Tests

Shell integration tests in `tests/` directory validate warning output format. These tests check for "Warning:" prefix and specific warning messages.

**Key test files:**
- `tests/test-folder-empty-warnings.sh` - Validates empty folder warnings
- `tests/test-error-messages.sh` - Validates error message formats
- `tests/test-config.sh` - Validates configuration warnings

### Edge Cases

- **Empty warnings list**: AggregateWarnings should return nil for empty input
- **Single item warning**: Should format as "message: item"
- **Multiple items warning**: Should format with bullet points
- **Mixed warning types**: Should group by type correctly
- **Thread safety**: Multiple goroutines calling AddWarning/GetWarnings
- **Clear operations**: Warnings/errors should be properly cleared

---

## VALIDATION COMMANDS

Execute every command to ensure zero regressions and 100% feature correctness.

### Level 1: Syntax & Style

```bash
go fmt ./...
```

```bash
golangci-lint run
```

```bash
semgrep scan --config auto
```

### Level 2: Build Verification

```bash
go build -o build/guard ./cmd/guard
```

```bash
cp build/guard tests/guard
```

### Level 3: Unit Tests

```bash
go test ./...
```

### Level 4: Integration Tests

```bash
cd tests && ./run-all-tests.sh
```

### Level 5: Manual Validation

Test warning aggregation manually:

```bash
# Initialize guard
./build/guard init 000 $(whoami) staff

# Create empty folder and test warning
mkdir emptyfolder
./build/guard toggle folder emptyfolder 2>&1 | grep -i "warning"

# Should see warning about empty folder
```

Test multiple warnings:

```bash
# Add multiple files that don't exist
./build/guard add file1.txt file2.txt file3.txt 2>&1 | grep -i "warning"

# Should see aggregated warnings with bullet points
```

---

## ACCEPTANCE CRITERIA

- [x] WarningType enum defined with iota and all 12 warning types
- [x] WarningFileAlreadyInRegistry is silent (produces no output)
- [x] Warning struct contains Type, Message, and Items fields
- [x] NewWarning constructor creates Warning instances
- [x] Type-specific aggregation functions implemented:
  - [x] aggregateFileMissing with cleanup suggestion
  - [x] aggregateCollectionEmpty
  - [x] aggregateFolderEmpty
  - [x] aggregateGeneric for fallback
- [x] AggregateWarnings dispatches to type-specific aggregators
- [x] Manager.warnings field changed from []string to []Warning
- [x] Manager.errors field added as []string
- [x] AddWarning method accepts Warning type
- [x] AddError method added
- [x] GetWarnings returns aggregated string slice
- [x] GetErrors, HasWarnings, HasErrors, ClearWarnings, ClearErrors methods added
- [x] PrintWarnings function takes []Warning and outputs to stdout
- [x] PrintErrors function takes []string and outputs to stderr
- [x] All AddWarning call sites migrated to use NewWarning
- [x] All validation commands pass with zero errors
- [x] Shell integration tests pass (especially warning-related tests)
- [x] No regressions in existing functionality
- [x] Thread safety maintained with proper mutex usage

---

## COMPLETION CHECKLIST

- [ ] All tasks completed in order
- [ ] Each task validation passed immediately
- [ ] All validation commands executed successfully
- [ ] Full test suite passes (unit + integration)
- [ ] No linting or type checking errors
- [ ] Manual testing confirms warnings display correctly
- [ ] Acceptance criteria all met
- [ ] Code reviewed for quality and maintainability

---

## NOTES

### Design Decisions

**Why iota for WarningType?**
- Type safety: Prevents invalid warning types at compile time
- Performance: Integer comparison is faster than string comparison
- Maintainability: Easy to add new warning types without breaking existing code

**Why separate warnings and errors?**
- Semantic clarity: Warnings are non-fatal, errors are fatal
- Different handling: Warnings can be aggregated, errors typically stop execution
- Future extensibility: May want different display or logging for errors vs warnings

**Why type-specific aggregation functions?**
- Customization: Each warning type has unique formatting requirements
- User experience: Specific messages like "run guard cleanup" provide actionable guidance
- Maintainability: Easy to modify formatting for one warning type without affecting others
- Extensibility: New warning types can have custom aggregation logic

**Why WarningFileAlreadyInRegistry is silent?**
- Idempotency: File addition is idempotent by design
- User experience: No need to warn about expected behavior
- Noise reduction: Prevents cluttering output with non-issues

**Why PrintWarnings outputs to stdout instead of stderr?**
- Semantic clarity: Warnings are informational, not errors
- User experience: Allows users to pipe warnings separately from errors
- Convention: Many CLI tools output warnings to stdout, errors to stderr

**Why Items as variadic parameter in NewWarning?**
- Convenience: Can create warnings with 0, 1, or many items easily
- Readability: `NewWarning(type, msg, file1, file2)` is cleaner than slice construction
- Common case: Most warnings have 0-1 items, variadic makes this ergonomic

### Trade-offs

**Aggregation complexity vs output quality:**
- Chose type-specific aggregation functions over generic formatting
- Benefit: Each warning type has tailored, actionable messages (e.g., "run guard cleanup")
- Cost: More aggregation functions to maintain (4 functions vs 1)

**Silent warnings for idempotent operations:**
- Chose to make WarningFileAlreadyInRegistry silent
- Benefit: Cleaner output, no noise for expected behavior
- Cost: Less visibility into what operations were skipped (acceptable trade-off)

**stdout vs stderr for warnings:**
- Chose stdout for warnings, stderr for errors
- Benefit: Semantic separation, better pipeline compatibility
- Cost: Different from some CLI conventions (acceptable for better UX)

**Type safety vs flexibility:**
- Chose enum with iota over string constants
- Benefit: Compile-time type checking, better IDE support
- Cost: Cannot create warning types dynamically (acceptable for this use case)

### Future Enhancements

- Add more type-specific aggregation functions as new warning types emerge
- Add warning severity levels (Info, Warning, Error)
- Implement warning filtering/suppression by type
- Add structured logging integration
- Create warning statistics/metrics collection
- Add warning deduplication within same operation
- Consider localization/i18n for warning messages

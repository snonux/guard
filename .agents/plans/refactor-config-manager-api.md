# Feature: Refactor Config Manager API

The following plan should be complete, but its important that you validate documentation and codebase patterns and task sanity before you start implementing.

Pay special attention to naming of existing utils types and models. Import from the right files etc.

## Feature Description

Refactor the config management system in the internal/manager package to use a cleaner API design with pointer parameters for optional updates, dedicated helper functions for common operations, and proper integration with the existing warning system instead of inline print statements.

## User Story

As a developer maintaining the guard-tool codebase
I want a cleaner config management API with pointer-based optional parameters
So that the code is more maintainable, follows Go best practices, and provides better separation of concerns between CLI parsing and business logic

## Problem Statement

The current config.go implementation has several issues:
1. The SetConfig method takes raw string args and parses them in the manager layer, mixing CLI concerns with business logic
2. Inline warning messages are printed directly instead of using the established warning system
3. Duplicate octal parsing logic scattered across methods
4. No dedicated methods for individual config updates
5. Config display formatting is hardcoded in ShowConfig method

## Solution Statement

Refactor the config management to:
1. Use pointer parameters (*string) for optional updates where nil means "don't update"
2. Move CLI argument parsing to the command layer
3. Add dedicated SetConfigMode, SetConfigOwner, SetConfigGroup methods
4. Create helper functions for octal parsing and config formatting
5. Integrate with the existing warning system for guarded items warnings
6. Maintain backward compatibility while improving the API design

## Feature Metadata

**Feature Type**: Refactor
**Estimated Complexity**: Medium
**Primary Systems Affected**: internal/manager/config.go, cmd/guard/commands/config.go
**Dependencies**: Existing warning system, registry security layer

---

## CONTEXT REFERENCES

### Relevant Codebase Files IMPORTANT: YOU MUST READ THESE FILES BEFORE IMPLEMENTING!

- `internal/manager/config.go` (entire file) - Why: Current implementation that needs refactoring
- `internal/manager/warnings.go` (lines 1-50, 200-250) - Why: Warning system patterns to follow
- `internal/manager/manager.go` (lines 300-350) - Why: Shows how warnings are accumulated and retrieved
- `cmd/guard/commands/config.go` (entire file) - Why: CLI layer that will handle argument parsing
- `internal/registry/registry.go` (lines 200-250) - Why: Shows octal parsing patterns already implemented

### New Files to Create

None - this is a refactoring of existing files

### Relevant Documentation YOU SHOULD READ THESE BEFORE IMPLEMENTING!

- [Go Pointer Parameters Best Practices](https://go.dev/doc/effective_go#pointers_vs_values)
  - Specific section: When to use pointers for optional parameters
  - Why: Required for implementing clean optional parameter API
- [Cobra CLI Argument Parsing](https://github.com/spf13/cobra#positional-and-custom-arguments)
  - Specific section: Args validation and parsing
  - Why: Shows proper CLI argument handling patterns

### Patterns to Follow

**Pointer Parameter Pattern:**
```go
// From existing codebase pattern (similar to registry methods)
func (m *Manager) UpdateSomething(required string, optional *string) error {
    if optional != nil {
        // Update the optional field
    }
    // Always update required field
}
```

**Warning System Integration:**
```go
// From internal/manager/warnings.go
m.AddWarning(NewWarning(WarningGeneric, "message", "item"))
```

**Octal Parsing Pattern:**
```go
// From internal/registry/registry.go (lines 200-220)
func octalStringToFileMode(value string) (os.FileMode, error) {
    // Existing implementation to reuse
}
```

**Config Display Pattern:**
```go
// New helper function pattern
func formatConfigValue(value string) string {
    if value == "" {
        return "(empty)"
    }
    return value
}
```

**Output Format Pattern:**
```go
// Config update messages should be indented:
fmt.Printf("Config updated:\n  Mode: %s\n", mode)
// Not inline: fmt.Printf("Config updated: mode set to %s\n", mode)
```

**Warning Integration Pattern:**
```go
// From internal/manager/warnings.go - count guarded items
func (m *Manager) checkAndWarnGuardedFiles() {
    fileCount, collectionCount := m.countGuardedItems()
    if fileCount > 0 || collectionCount > 0 {
        message := fmt.Sprintf("%d file(s) and %d collection(s) are currently guarded.\nThe new config will only apply to future guard operations.\nTo apply the new config to existing guards, disable and re-enable them.", fileCount, collectionCount)
        m.AddWarning(NewWarning(WarningGeneric, message, ""))
    }
}
```

**Error Message Pattern:**
```go
// Registry not found should say:
return fmt.Errorf(".guardfile not found. Run 'guard init' first.")
// Not: return fmt.Errorf("registry not loaded")
```

---

## IMPLEMENTATION PLAN

### Phase 1: Foundation

Create helper functions and establish new API patterns without breaking existing functionality.

**Tasks:**
- Add helper functions for octal parsing and config formatting
- Create internal helper to check for guarded items with warning integration
- Add pointer-based config update methods

### Phase 2: Core Implementation

Implement the new SetConfig method with pointer parameters and refactor existing methods.

**Tasks:**
- Refactor SetConfig to use pointer parameters
- Add dedicated SetConfigMode, SetConfigOwner, SetConfigGroup methods
- Update ShowConfig to use formatting helpers
- Replace inline warnings with warning system integration

### Phase 3: Integration

Update the CLI layer to handle argument parsing and call the new API methods.

**Tasks:**
- Update config command to parse arguments before calling manager
- Add validation and error handling in CLI layer
- Ensure backward compatibility with existing command patterns

### Phase 4: Testing & Validation

Validate the refactoring maintains existing functionality while improving the API.

**Tasks:**
- Run existing tests to ensure no regressions
- Test new pointer parameter behavior
- Validate warning system integration
- Test CLI argument parsing edge cases

---

## STEP-BY-STEP TASKS

IMPORTANT: Execute every task in order, top to bottom. Each task is atomic and independently testable.

### REFACTOR internal/manager/config.go

- **ADD**: Helper function `parseOctalMode(mode string) (os.FileMode, error)`
- **PATTERN**: Mirror octal parsing from `internal/registry/registry.go:octalStringToFileMode`
- **IMPORTS**: No new imports needed
- **GOTCHA**: Ensure 3-digit modes are padded to 4 digits like existing code
- **VALIDATE**: `go build ./internal/manager`

### ADD internal/manager/config.go

- **IMPLEMENT**: Helper function `formatConfigValue(value string) string`
- **PATTERN**: Return "(empty)" for empty strings, otherwise return the value as-is
- **IMPORTS**: No new imports needed
- **GOTCHA**: Handle empty string vs non-empty string display consistently
- **VALIDATE**: `go build ./internal/manager`

### ADD internal/manager/config.go

- **IMPLEMENT**: Helper function `checkAndWarnGuardedFiles(m *Manager)`
- **PATTERN**: Count guarded files and collections, add detailed warning with counts
- **IMPORTS**: No new imports needed
- **GOTCHA**: Use countGuardedItems helper to get accurate counts, format message with proper grammar
- **VALIDATE**: `go build ./internal/manager`

### ADD internal/manager/config.go

- **IMPLEMENT**: Helper function `countGuardedItems() (int, int)`
- **PATTERN**: Count guarded files and collections, return both counts
- **IMPORTS**: No new imports needed
- **GOTCHA**: Check both files and collections for guard status, handle errors gracefully
- **VALIDATE**: `go build ./internal/manager`

### ADD internal/manager/config.go

- **IMPLEMENT**: Method `SetConfigMode(mode string) error`
- **PATTERN**: Take regular string parameter for explicit single value updates
- **IMPORTS**: No new imports needed
- **GOTCHA**: Validate octal mode using parseOctalMode helper, call checkAndWarnGuardedFiles
- **VALIDATE**: `go build ./internal/manager`

### ADD internal/manager/config.go

- **IMPLEMENT**: Method `SetConfigOwner(owner string) error`
- **PATTERN**: Take regular string parameter for explicit single value updates
- **IMPORTS**: No new imports needed
- **GOTCHA**: Call checkAndWarnGuardedFiles after successful update
- **VALIDATE**: `go build ./internal/manager`

### ADD internal/manager/config.go

- **IMPLEMENT**: Method `SetConfigGroup(group string) error`
- **PATTERN**: Take regular string parameter for explicit single value updates
- **IMPORTS**: No new imports needed
- **GOTCHA**: Call checkAndWarnGuardedFiles after successful update
- **VALIDATE**: `go build ./internal/manager`

### REFACTOR internal/manager/config.go

- **IMPLEMENT**: New `SetConfig(mode *string, owner *string, group *string) error` signature
- **PATTERN**: Use pointer parameters, call individual Set methods only for non-nil values
- **IMPORTS**: No new imports needed
- **GOTCHA**: Only call checkAndWarnGuardedFiles once at the end if any updates were made
- **VALIDATE**: `go build ./internal/manager`

### UPDATE internal/manager/config.go

- **IMPLEMENT**: Refactor `ShowConfig()` to use `formatConfigValue` helper
- **PATTERN**: Use formatConfigValue for owner and group, keep mode formatting as-is
- **IMPORTS**: No new imports needed
- **GOTCHA**: Maintain exact same output format for backward compatibility
- **VALIDATE**: `go build ./internal/manager`

### REMOVE internal/manager/config.go

- **IMPLEMENT**: Remove old methods: `processConfigArgs`, `setMultipleConfig`, `setModeAndOwner`, `setSingleConfig`
- **PATTERN**: Clean up deprecated methods after refactoring
- **IMPORTS**: No imports to remove
- **GOTCHA**: Ensure no other files reference these methods
- **VALIDATE**: `go build ./...`

### REMOVE internal/manager/config.go

- **IMPLEMENT**: Remove old `hasGuardedItems()` method and `isOctalMode()` helper
- **PATTERN**: Replace with new helper functions
- **IMPORTS**: No imports to remove
- **GOTCHA**: Ensure functionality is preserved in new helpers
- **VALIDATE**: `go build ./internal/manager`

### UPDATE cmd/guard/commands/config.go

- **IMPLEMENT**: Refactor `runConfigSet` with inline switch case on args[0] for keywords
- **PATTERN**: Switch on "mode", "owner", "group" keywords, call dedicated methods, fallthrough to SetConfig for bulk updates
- **IMPORTS**: No new imports needed
- **GOTCHA**: Handle all existing argument patterns (mode only, mode+owner, mode+owner+group, setting+value)
- **VALIDATE**: `go build ./cmd/guard`

### REMOVE cmd/guard/commands/config.go

- **IMPLEMENT**: Remove the separate parseConfigArgs helper function approach
- **PATTERN**: Use inline switch case parsing instead of separate function
- **IMPORTS**: No imports to remove
- **GOTCHA**: Ensure all argument patterns are handled in the switch case
- **VALIDATE**: `go build ./cmd/guard`

### UPDATE cmd/guard/commands/config.go

- **IMPLEMENT**: Add validation and error handling for parsed arguments
- **PATTERN**: Validate octal modes and argument combinations before calling manager
- **IMPORTS**: Add `strconv` for octal validation
- **GOTCHA**: Provide clear error messages for invalid argument patterns
- **VALIDATE**: `go build ./cmd/guard`

---

## TESTING STRATEGY

### Unit Tests

The project uses shell-based integration tests rather than Go unit tests. Focus on integration testing through the CLI.

### Integration Tests

Test through existing shell test framework in `tests/` directory:

- Test `guard config show` output format unchanged
- Test `guard config set mode 640` single parameter updates
- Test `guard config set 640 owner group` multiple parameter updates
- Test `guard config set owner newowner` and `guard config set group newgroup`
- Test warning messages when guarded items exist
- Test error handling for invalid octal modes

### Edge Cases

- Empty string vs nil parameter handling
- Invalid octal mode validation (non-octal, out of range)
- Registry not loaded error handling
- Argument parsing edge cases (too few, too many args)
- Warning system integration with existing guarded files

---

## VALIDATION COMMANDS

Execute every command to ensure zero regressions and 100% feature correctness.

### Level 1: Syntax & Style

```bash
go fmt ./...
golangci-lint run
```

### Level 2: Build Validation

```bash
go build ./...
go build -o build/guard ./cmd/guard
```

### Level 3: Integration Tests

```bash
cd tests
cp ../build/guard ./guard
./test-config.sh
./run-all-tests.sh
```

### Level 4: Manual Validation

```bash
# Test new API maintains backward compatibility
./guard init 640 testuser testgroup
./guard config show
./guard config set mode 644
./guard config set owner newowner
./guard config set group newgroup
./guard config set 600 finalowner finalgroup
```

### Level 5: Additional Validation

```bash
# Test warning system integration
./guard add file1.txt
./guard enable file1.txt
./guard config set mode 755  # Should show warning about guarded items
```

---

## ACCEPTANCE CRITERIA

- [ ] SetConfig method uses pointer parameters (*string) for optional updates
- [ ] Nil pointer parameters mean "don't update that value"
- [ ] Dedicated SetConfigMode, SetConfigOwner, SetConfigGroup methods exist
- [ ] Helper functions for octal parsing and config formatting implemented
- [ ] Warning system integration replaces inline print statements
- [ ] CLI layer handles argument parsing before calling manager methods
- [ ] All existing command patterns work unchanged (backward compatibility)
- [ ] All validation commands pass with zero errors
- [ ] Integration tests verify end-to-end functionality
- [ ] Code follows existing project patterns and conventions
- [ ] No regressions in existing functionality
- [ ] Warning messages use proper warning system aggregation

---

## COMPLETION CHECKLIST

- [ ] All tasks completed in order
- [ ] Each task validation passed immediately
- [ ] All validation commands executed successfully
- [ ] Full test suite passes (shell integration tests)
- [ ] No linting or build errors
- [ ] Manual testing confirms backward compatibility
- [ ] Acceptance criteria all met
- [ ] Code reviewed for quality and maintainability

---

## NOTES

**Design Decisions:**
- Dedicated methods (SetConfigMode, SetConfigOwner, SetConfigGroup) use regular string parameters for explicit updates
- Only the main SetConfig method uses pointer parameters for optional updates (nil = don't update)
- Helper function formatConfigValue handles single value formatting with "(empty)" for empty strings
- checkAndWarnGuardedFiles adds warnings directly rather than returning boolean
- Config update output uses indented format for better readability
- Error messages are more user-friendly ("guardfile not found, run 'guard init' first")

**Trade-offs:**
- Slightly more complex CLI parsing logic in exchange for cleaner manager API
- Additional helper functions increase code size but improve maintainability
- Pointer parameters require nil checks but provide clear optional semantics

**Future Considerations:**
- This refactoring enables easier testing of config logic in isolation
- The new API makes it easier to add config validation rules
- Helper functions can be reused for future config-related features

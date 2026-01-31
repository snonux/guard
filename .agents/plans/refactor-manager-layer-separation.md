# Feature: Refactor Manager Layer for Proper Separation of Concerns

The following plan should be complete, but its important that you validate documentation and codebase patterns and task sanity before you start implementing.

Pay special attention to naming of existing utils types and models. Import from the right files etc.

## Feature Description

Refactor the Manager layer in `internal/manager/manager.go` to properly separate orchestration logic from command-specific concerns. The manager should only coordinate between the security and filesystem layers without handling UI, printing, or command-specific logic. This refactoring will improve maintainability, testability, and adherence to single responsibility principle.

## User Story

As a developer maintaining the guard-tool codebase
I want the Manager layer to only handle orchestration between other layers
So that command-specific logic is properly separated and the codebase is more maintainable

## Problem Statement

The current Manager layer violates single responsibility principle by:
1. Mixing business logic with UI concerns (printing, formatting)
2. Having side effects in getter methods (GetWarnings clears warnings)
3. Containing command-specific logic that should be in the command layer
4. Lacking proper accessor methods for composed dependencies
5. Missing helper methods for common operations

## Solution Statement

Refactor the Manager to be a pure orchestration layer that:
1. Returns structured data instead of printing
2. Provides clean getter methods without side effects
3. Exposes necessary accessor methods for dependencies
4. Includes helper methods for common operations
5. Moves all UI and command-specific logic to the command layer

## Feature Metadata

**Feature Type**: Refactor
**Estimated Complexity**: Medium
**Primary Systems Affected**: Manager layer, Command layer
**Dependencies**: None (internal refactoring)

---

## CONTEXT REFERENCES

### Relevant Codebase Files IMPORTANT: YOU MUST READ THESE FILES BEFORE IMPLEMENTING!

- `internal/manager/manager.go` (lines 1-440) - Why: Main file being refactored, contains all methods to be modified/removed
- `internal/manager/warnings.go` (lines 39-50) - Why: Warning type definition needed for GetWarnings return type
- `internal/security/security.go` (lines 1-30) - Why: Security layer interface that manager orchestrates
- `internal/filesystem/filesystem.go` - Why: Filesystem layer interface that manager orchestrates
- `internal/registry/registry.go` (lines 1-100) - Why: Registry interface for accessor methods
- `cmd/guard/commands/show.go` (lines 1-50) - Why: Example of current manager usage that needs updating
- `cmd/guard/commands/toggle.go` (lines 1-50) - Why: Example of ResolveArguments usage pattern

### New Files to Create

None - this is a refactoring of existing files

### Relevant Documentation YOU SHOULD READ THESE BEFORE IMPLEMENTING!

- [Go Code Review Comments - Getters](https://github.com/golang/go/wiki/CodeReviewComments#getters)
  - Specific section: Getter naming conventions
  - Why: Proper accessor method patterns
- [Effective Go - Interface Names](https://golang.org/doc/effective_go.html#interface-names)
  - Specific section: Single method interfaces
  - Why: Clean interface design principles

### Patterns to Follow

**Warning Handling Pattern:**
```go
// Current (bad): Side effect in getter
func (m *Manager) GetWarnings() []string {
    aggregated := AggregateWarnings(m.warnings)
    m.warnings = m.warnings[:0] // Side effect!
    return aggregated
}

// New (good): Pure getter
func (m *Manager) GetWarnings() []Warning {
    return m.warnings
}
```

**Accessor Method Pattern:**
```go
// Follow existing pattern from registry.go
func (r *Registry) GetDefaultFileMode() os.FileMode {
    r.mu.RLock()
    defer r.mu.RUnlock()
    // return data
}
```

**Error Handling Pattern:**
```go
// Follow existing pattern from manager methods
if m.security == nil {
    return fmt.Errorf("registry not loaded")
}
```

**Method Naming Pattern:**
- Use existing naming: `IsRegistered*`, `Get*`, `Set*`, `Add*`, `Remove*`
- Private methods: `clearGuardfileImmutableFlag`, `toDisplayPaths`

---

## IMPLEMENTATION PLAN

### Phase 1: Foundation - Clean Up Existing Methods

Remove command-specific logic and UI concerns from manager layer, preparing for clean separation.

**Tasks:**
- Remove methods that belong in command layer
- Fix GetWarnings to return proper type without side effects
- Add new required methods with proper signatures

### Phase 2: Core Implementation - Add New Methods

Implement the new methods required for proper manager functionality.

**Tasks:**
- Add accessor methods for composed dependencies
- Add helper methods for common operations
- Add registry initialization method with proper parameters

### Phase 3: Integration - Update Method Signatures

Update existing methods to follow clean patterns and remove UI concerns.

**Tasks:**
- Update existing methods to return structured data
- Ensure all methods follow single responsibility principle
- Add proper error handling and validation

### Phase 4: Testing & Validation

Ensure refactored manager maintains all functionality while improving separation.

**Tasks:**
- Validate all existing functionality still works
- Ensure no regressions in command behavior
- Verify clean separation of concerns

---

## STEP-BY-STEP TASKS

IMPORTANT: Execute every task in order, top to bottom. Each task is atomic and independently testable.

### REMOVE Methods That Belong in Command Layer

- **REMOVE**: `SetProtection`, `SetProtectionWithType`, `ToggleProtection`, `ToggleProtectionWithType` methods
- **REMOVE**: `Show`, `showFiles`, `showCollections` methods  
- **REMOVE**: `Uninstall` method
- **REMOVE**: `PrintWarnings`, `PrintErrors` methods
- **REMOVE**: `detectTargetType`, `toggleByTypeString` private methods
- **PATTERN**: These methods mix business logic with UI concerns
- **GOTCHA**: Commands currently call these methods - they'll need to be moved to command layer
- **VALIDATE**: `go build ./cmd/guard` (will fail initially, that's expected)

### UPDATE GetWarnings Method Signature

- **REFACTOR**: Change `GetWarnings() []string` to `GetWarnings() []Warning`
- **REMOVE**: Side effect of clearing warnings in GetWarnings
- **PATTERN**: Pure getter method following `internal/registry/registry.go:GetDefaultFileMode` pattern
- **IMPORTS**: Ensure Warning type is properly accessible
- **GOTCHA**: Commands using GetWarnings will need updates
- **VALIDATE**: `go build ./internal/manager`

### ADD InitializeRegistry Method

- **CREATE**: `InitializeRegistry(mode, owner, group string, overwrite bool) error` method
- **REPLACE**: Current `Init` method that has `interactive` parameter
- **PATTERN**: Follow existing `NewRegistry` pattern from `internal/registry/registry.go:NewRegistry`
- **IMPORTS**: Use `registry.RegistryDefaults` struct
- **GOTCHA**: Remove interactive parameter - that's UI concern for command layer
- **VALIDATE**: `go build ./internal/manager`

### ADD clearGuardfileImmutableFlag Private Method

- **CREATE**: `clearGuardfileImmutableFlag() error` private method
- **IMPLEMENT**: Clear immutable flag before any write operation to .guardfile
- **PATTERN**: Follow filesystem operation patterns from `internal/filesystem/filesystem.go`
- **IMPORTS**: Use filesystem layer for actual flag clearing
- **GOTCHA**: Must be called before any registry save operations
- **VALIDATE**: `go build ./internal/manager`

### ADD Accessor Methods

- **CREATE**: `GetRegistry() *security.Security` method
- **CREATE**: `GetFileSystem() *filesystem.FileSystem` method
- **PATTERN**: Simple accessor methods exposing composed dependencies
- **IMPORTS**: Return proper types from internal packages - security wrapper not direct registry
- **GOTCHA**: Manager holds security wrapper, not direct registry reference
- **VALIDATE**: `go build ./internal/manager`

### ADD Collection Helper Methods

- **CREATE**: `IsRegisteredCollection(name string) bool` method
- **CREATE**: `CountFilesInCollection(name string) (int, error)` method
- **PATTERN**: Follow existing `IsRegisteredFile` pattern in manager.go:419
- **IMPORTS**: Delegate to security layer for actual checks
- **GOTCHA**: Ensure proper error handling for non-existent collections
- **VALIDATE**: `go build ./internal/manager`

### ADD toDisplayPaths Helper Method

- **CREATE**: `toDisplayPaths(paths []string) []string` private method
- **IMPLEMENT**: Delegate to `m.security.ToDisplayPath()` for each path in slice
- **PATTERN**: Convert absolute paths to relative display paths by stripping registry root directory prefix
- **IMPORTS**: Delegate to security layer for path conversion logic
- **GOTCHA**: May need to add ToDisplayPath method to security layer if it doesn't exist
- **VALIDATE**: `go build ./internal/manager`

### ADD ResolveArgument Singular Method

- **CREATE**: `ResolveArgument(arg string) (string, error)` method
- **PATTERN**: Extract logic from existing `detectTargetType` method before removing it
- **IMPORTS**: Use same detection logic but return only type as string
- **GOTCHA**: Return type as string ("file", "collection", "folder") - argument is not modified or cleaned
- **VALIDATE**: `go build ./internal/manager`

### UPDATE ResolveArguments Method

- **REFACTOR**: Update `ResolveArguments` to use new `ResolveArgument` method internally
- **PATTERN**: Keep existing signature but use new singular method for implementation
- **IMPORTS**: No new imports needed
- **GOTCHA**: Maintain backward compatibility with existing callers
- **VALIDATE**: `go build ./internal/manager`

### UPDATE LoadRegistry Method Error Handling

- **IMPROVE**: Add specific error messages for common failure cases
- **PATTERN**: Check if file doesn't exist vs corrupted file scenarios  
- **IMPORTS**: Use os.IsNotExist for file existence checks
- **GOTCHA**: Provide helpful error messages: "guardfile not found in current directory. Run 'guard init' to initialize" and "guardfile is corrupted" with recovery suggestions
- **VALIDATE**: `go build ./internal/manager`

### REMOVE Target Type Constants

- **REMOVE**: `TargetTypeFile`, `TargetTypeFolder`, `TargetTypeCollection` constants
- **PATTERN**: These are now handled as strings in command layer
- **GOTCHA**: Any remaining references will cause build failures
- **VALIDATE**: `go build ./cmd/guard` (should pass after command layer updates)

---

## TESTING STRATEGY

### Unit Tests

Design unit tests for new methods following existing test patterns in the project:

- Test `GetWarnings()` returns Warning slice without side effects
- Test `InitializeRegistry()` with various parameter combinations
- Test accessor methods return proper instances
- Test `ResolveArgument()` with different target types
- Test helper methods with edge cases

### Integration Tests

Validate manager orchestration still works correctly:

- Test manager coordinates between security and filesystem layers
- Test warning accumulation and retrieval
- Test registry operations through manager
- Verify no UI concerns remain in manager layer

### Edge Cases

Test specific edge cases for refactored functionality:

- Empty warning lists
- Invalid registry parameters
- Non-existent collections for count operations
- Path formatting edge cases
- Error propagation from underlying layers

---

## VALIDATION COMMANDS

Execute every command to ensure zero regressions and 100% feature correctness.

### Level 1: Syntax & Style

```bash
go fmt ./internal/manager/...
golangci-lint run ./internal/manager/...
```

### Level 2: Unit Tests

```bash
go test ./internal/manager/... -v
```

### Level 3: Integration Tests

```bash
go build ./cmd/guard
go test ./... -v
```

### Level 4: Manual Validation

```bash
# Test manager can be instantiated and basic methods work
cd tests && ./guard init 0644 testuser testgroup
cd tests && ./guard show
```

### Level 5: Additional Validation

```bash
# Run full test suite to ensure no regressions
cd tests && ./run-all-tests.sh
```

---

## ACCEPTANCE CRITERIA

- [ ] Manager layer contains only orchestration logic
- [ ] No printing or UI concerns remain in manager
- [ ] GetWarnings returns Warning slice without side effects
- [ ] All new required methods implemented with proper signatures
- [ ] Accessor methods provide clean access to dependencies
- [ ] Helper methods support common operations
- [ ] All validation commands pass with zero errors
- [ ] No regressions in existing functionality
- [ ] Clean separation between manager and command layers
- [ ] Code follows existing project patterns and conventions

---

## COMPLETION CHECKLIST

- [ ] All methods to be removed are deleted
- [ ] GetWarnings refactored to return Warning slice
- [ ] InitializeRegistry method added with correct signature
- [ ] clearGuardfileImmutableFlag private method added
- [ ] Accessor methods GetRegistry (*security.Security) and GetFileSystem added
- [ ] Collection helper methods added
- [ ] toDisplayPaths helper method added
- [ ] ResolveArgument singular method added (returns string, error)
- [ ] ResolveArguments updated to use new method
- [ ] LoadRegistry error handling improved with specific messages
- [ ] Target type constants removed
- [ ] All validation commands pass
- [ ] No build errors in manager package
- [ ] Manager layer properly separated from UI concerns

---

## NOTES

**Design Decisions:**
- GetWarnings returns Warning slice instead of strings to maintain type safety
- InitializeRegistry replaces Init to remove UI-specific interactive parameter
- Accessor methods expose dependencies for command layer orchestration
- Private helper methods support internal operations without exposing implementation

**Trade-offs:**
- Command layer will need updates to handle UI concerns moved from manager
- Slightly more verbose command implementations but cleaner separation
- Better testability at the cost of some initial refactoring effort

**Future Considerations:**
- Consider interface-based design for better testability
- Evaluate if manager should return result structs instead of using error accumulation
- Consider adding context.Context parameters for cancellation support

# Feature: Refactor Security Layer

The following plan should be complete, but its important that you validate documentation and codebase patterns and task sanity before you start implementing.

Pay special attention to naming of existing utils types and models. Import from the right files etc.

## Feature Description

Refactor the internal/security layer to improve code organization, reduce complexity, ensure proper error handling, and follow Go best practices. The security layer currently wraps the registry and provides path validation, symlink protection, and permission enforcement, but suffers from poor separation of concerns, inconsistent error handling, and violation of Go conventions.

## User Story

As a developer maintaining the guard-tool codebase
I want a well-structured, testable, and maintainable security layer
So that I can easily extend security features, debug issues, and ensure reliable path validation and permission enforcement

## Problem Statement

The current security layer (`internal/security/security.go`) has several architectural and code quality issues:

1. **Poor Separation of Concerns**: The Security struct acts as both a wrapper and contains standalone validation functions
2. **Inconsistent Error Handling**: Some methods return errors, others return false on validation failure
3. **Mixed Responsibilities**: Path validation, symlink checking, and registry wrapping are all in one file
4. **No Unit Tests**: Critical security functionality lacks test coverage
5. **Violation of Go Conventions**: Package-level functions mixed with methods, inconsistent naming
6. **Code Duplication**: Path validation logic is repeated across methods
7. **Poor Interface Design**: The Security struct exposes too many registry methods directly

## Solution Statement

Refactor the security layer into a clean, modular architecture with:
- Separate concerns: path validation, symlink protection, and registry security wrapper
- Consistent error handling patterns following Go conventions
- Comprehensive unit test coverage for all security functions
- Clear interfaces and dependency injection
- Proper abstraction layers that don't leak registry implementation details

## Feature Metadata

**Feature Type**: Refactor
**Estimated Complexity**: Medium
**Primary Systems Affected**: Security layer, Manager layer (interface changes)
**Dependencies**: No external dependencies required

---

## CONTEXT REFERENCES

### Relevant Codebase Files IMPORTANT: YOU MUST READ THESE FILES BEFORE IMPLEMENTING!

- `internal/security/security.go` (entire file) - Why: Current implementation to be refactored
- `internal/manager/manager.go` (lines 16, 39, 61, 92, 98, 156-158) - Why: Shows how manager uses security layer
- `internal/manager/files.go` (lines 30, 45, 58, 91, 97, 102, etc.) - Why: Heavy usage of security methods, interface requirements
- `internal/manager/collections.go` (lines 50, 56, 82, 87, 116) - Why: Collection-related security usage patterns
- `internal/manager/folders.go` (lines 74, 79, 101, 102, 180) - Why: Folder-related security usage patterns
- `internal/manager/config.go` (lines 15-17, 53, 62) - Why: Configuration access patterns through security layer
- `internal/registry/registry.go` (lines 56-95, 99-144) - Why: Registry interface that security wraps
- `internal/registry/file_entry.go` - Why: File-related registry methods that security exposes
- `internal/registry/collection_entry.go` - Why: Collection-related registry methods that security exposes
- `internal/registry/folder_entry.go` - Why: Folder-related registry methods that security exposes

### New Files to Create

- `internal/security/validator.go` - Path validation and symlink protection logic
- `internal/security/wrapper.go` - Registry security wrapper with clean interface
- `internal/security/errors.go` - Security-specific error types and constants
- `internal/security/validator_test.go` - Unit tests for path validation
- `internal/security/wrapper_test.go` - Unit tests for registry wrapper
- `internal/security/security_test.go` - Integration tests for security package

### Relevant Documentation YOU SHOULD READ THESE BEFORE IMPLEMENTING!

- [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments)
  - Specific section: Error handling, naming conventions, package design
  - Why: Ensures refactored code follows Go best practices
- [Effective Go](https://golang.org/doc/effective_go.html)
  - Specific section: Interfaces, error handling, package structure
  - Why: Guide for proper Go package design and error handling patterns

### Patterns to Follow

**Error Handling Pattern** (from existing codebase):
```go
// From internal/manager/files.go:58
if err := m.security.RegisterFile(path, mode, owner, group); err != nil {
    m.AddError(fmt.Sprintf("Error: Failed to register %s: %v", path, err))
    continue
}
```

**Validation Pattern** (from internal/manager/files.go:30):
```go
// Validate all paths first (security check happens regardless of file existence)
if err := m.security.ValidatePaths(paths); err != nil {
    return err
}
```

**Registry Wrapper Pattern** (from current security.go):
```go
// Current pattern - needs improvement
func (s *Security) IsRegisteredFile(filePath string) bool {
    cleanPath, err := ValidateAndCleanPath(filePath)
    if err != nil {
        return false // This loses error information!
    }
    return s.registry.IsRegisteredFile(cleanPath)
}
```

**Go Package Structure Pattern** (from project structure):
- Package-level types and constants at top
- Constructors after type definitions  
- Public methods before private methods
- Helper functions at bottom

**Testing Pattern** (from project justfile):
```go
// Unit tests should use Go's built-in testing framework
// Integration tests in tests/ directory use shell scripts
```

---

## IMPLEMENTATION PLAN

### Phase 1: Foundation - Error Types and Validation

Create security-specific error types and extract path validation logic into a dedicated module with comprehensive error handling.

**Tasks:**
- Define security error types and constants
- Extract path validation logic from security.go
- Create comprehensive path validator with proper error handling
- Add unit tests for path validation

### Phase 2: Registry Wrapper Refactoring

Refactor the Security struct to be a clean wrapper around the registry with proper error propagation and interface design.

**Tasks:**
- Create new registry wrapper with clean interface
- Implement proper error handling for all registry operations
- Remove direct method forwarding, add validation layer
- Maintain backward compatibility with manager layer

### Phase 3: Integration and Testing

Integrate the refactored components, ensure manager layer compatibility, and add comprehensive test coverage.

**Tasks:**
- Update security.go to use new components
- Verify manager layer compatibility
- Add integration tests
- Update imports and ensure no breaking changes

### Phase 4: Cleanup and Documentation

Remove deprecated code, add documentation, and ensure code quality standards.

**Tasks:**
- Remove old validation functions
- Add comprehensive package documentation
- Run linting and formatting
- Verify all tests pass

---

## STEP-BY-STEP TASKS

IMPORTANT: Execute every task in order, top to bottom. Each task is atomic and independently testable.

### CREATE internal/security/errors.go

- **IMPLEMENT**: Security-specific error types and constants
- **PATTERN**: Error handling from internal/manager/files.go:58
- **IMPORTS**: fmt, errors packages
- **GOTCHA**: Use errors.New() for static errors, fmt.Errorf() for dynamic errors
- **VALIDATE**: `go build ./internal/security`

### CREATE internal/security/validator.go

- **IMPLEMENT**: Path validation and symlink protection logic extracted from security.go
- **PATTERN**: Validation pattern from internal/manager/files.go:30
- **IMPORTS**: os, path/filepath, fmt, strings
- **GOTCHA**: Ensure ValidateAndCleanPath returns consistent error types
- **VALIDATE**: `go build ./internal/security && go test ./internal/security -run TestValidator`

### CREATE internal/security/validator_test.go

- **IMPLEMENT**: Comprehensive unit tests for path validation functions
- **PATTERN**: Go testing conventions from project structure
- **IMPORTS**: testing, os, path/filepath, strings
- **GOTCHA**: Test edge cases like empty paths, symlinks, directory traversal
- **VALIDATE**: `go test ./internal/security -run TestValidator -v`

### CREATE internal/security/wrapper.go

- **IMPLEMENT**: Clean registry wrapper with proper error handling
- **PATTERN**: Registry access from internal/manager/files.go:45
- **IMPORTS**: os, internal/registry package
- **GOTCHA**: Don't expose registry methods directly, add validation layer
- **VALIDATE**: `go build ./internal/security`

### CREATE internal/security/wrapper_test.go

- **IMPLEMENT**: Unit tests for registry wrapper functionality
- **PATTERN**: Testing pattern from project conventions
- **IMPORTS**: testing, os, internal/registry
- **GOTCHA**: Mock registry for isolated testing
- **VALIDATE**: `go test ./internal/security -run TestWrapper -v`

### REFACTOR internal/security/security.go

- **IMPLEMENT**: Update Security struct to use new validator and wrapper components
- **PATTERN**: Constructor pattern from internal/filesystem/filesystem.go:17
- **IMPORTS**: Update to use new internal components
- **GOTCHA**: Maintain exact same public interface for manager compatibility
- **VALIDATE**: `go build ./internal/security && go test ./internal/manager -v`

### UPDATE internal/security/security.go - Remove deprecated functions

- **IMPLEMENT**: Remove old ValidatePath, RejectSymlinks, ValidateAndCleanPath functions
- **PATTERN**: Clean removal without breaking imports
- **IMPORTS**: No changes needed
- **GOTCHA**: Ensure no other packages import these functions directly
- **VALIDATE**: `go build ./... && grep -r "ValidatePath\|RejectSymlinks" internal/`

### CREATE internal/security/security_test.go

- **IMPLEMENT**: Integration tests for complete security package
- **PATTERN**: Integration testing approach from project
- **IMPORTS**: testing, os, internal/registry
- **GOTCHA**: Test actual file operations and registry interactions
- **VALIDATE**: `go test ./internal/security -v`

### ADD package documentation to internal/security/security.go

- **IMPLEMENT**: Comprehensive package-level documentation
- **PATTERN**: Godoc conventions from existing packages
- **IMPORTS**: No changes
- **GOTCHA**: Document security guarantees and usage patterns
- **VALIDATE**: `go doc ./internal/security`

---

## TESTING STRATEGY

### Unit Tests

Design unit tests with fixtures and assertions following existing Go testing approaches:

- **Path Validation Tests**: Test all edge cases for ValidatePath, RejectSymlinks, ValidateAndCleanPath
- **Registry Wrapper Tests**: Test error propagation and validation layer
- **Error Type Tests**: Verify error types and messages are consistent

### Integration Tests

- **Security Package Integration**: Test validator + wrapper + registry interactions
- **Manager Compatibility**: Ensure manager layer continues to work without changes
- **File System Operations**: Test actual path validation with real files and symlinks

### Edge Cases

- Empty and whitespace-only paths
- Paths with directory traversal attempts (../, ../../)
- Symlink detection and rejection
- Non-existent files and directories
- Permission denied scenarios
- Registry corruption and recovery

---

## VALIDATION COMMANDS

Execute every command to ensure zero regressions and 100% feature correctness.

### Level 1: Syntax & Style

```bash
go fmt ./internal/security/...
golangci-lint run ./internal/security/...
```

### Level 2: Unit Tests

```bash
go test ./internal/security -v
go test ./internal/security -race
go test ./internal/security -cover
```

### Level 3: Integration Tests

```bash
go test ./internal/manager -v
go test ./... -short
go build ./cmd/guard
```

### Level 4: Manual Validation

```bash
# Test path validation
./guard init 0640 $(whoami) $(id -gn)
echo "test" > testfile.txt
./guard add testfile.txt
./guard show testfile.txt

# Test symlink rejection
ln -s testfile.txt symlink.txt
./guard add symlink.txt  # Should fail with symlink error
rm symlink.txt

# Test directory traversal protection
./guard add ../outside.txt  # Should fail with path validation error
```

### Level 5: Additional Validation

```bash
# Full CI pipeline
just ci-quiet
```

---

## ACCEPTANCE CRITERIA

- [ ] Security layer has clear separation of concerns (validator, wrapper, errors)
- [ ] All security functions have comprehensive unit test coverage (80%+)
- [ ] Manager layer continues to work without any interface changes
- [ ] Path validation is consistent and properly handles all edge cases
- [ ] Error handling follows Go conventions with proper error wrapping
- [ ] No direct registry method exposure - all access goes through validation
- [ ] Symlink protection works correctly and returns proper errors
- [ ] Package documentation clearly explains security guarantees
- [ ] All validation commands pass with zero errors
- [ ] No regressions in existing functionality
- [ ] Code follows Go best practices and project conventions

---

## COMPLETION CHECKLIST

- [ ] All tasks completed in order
- [ ] Each task validation passed immediately
- [ ] All validation commands executed successfully
- [ ] Full test suite passes (unit + integration)
- [ ] No linting or type checking errors
- [ ] Manual testing confirms security features work
- [ ] Manager layer compatibility verified
- [ ] Acceptance criteria all met
- [ ] Code reviewed for quality and maintainability

---

## NOTES

### Design Decisions

1. **Three-Component Architecture**: Separate validator, wrapper, and error handling for clear responsibilities
2. **Backward Compatibility**: Maintain exact same public interface to avoid breaking manager layer
3. **Error Propagation**: Always return errors instead of silent failures (bool returns)
4. **Validation First**: All registry operations go through path validation layer
5. **No Direct Registry Access**: Manager layer cannot bypass security validation

### Trade-offs

- **Performance vs Security**: Added validation layer may have minimal performance impact, but security is prioritized
- **Interface Stability vs Clean Design**: Maintaining backward compatibility limits some design improvements
- **Test Coverage vs Complexity**: Comprehensive testing adds complexity but ensures reliability

### Future Considerations

- Consider adding security audit logging for sensitive operations
- Potential for adding configurable security policies
- Integration with external security scanning tools
- Performance optimization for batch operations

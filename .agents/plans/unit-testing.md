# Feature: Unit Testing

The following plan should be complete, but its important that you validate documentation and codebase patterns and task sanity before you start implementing.

Pay special attention to naming of existing utils types and models. Import from the right files etc.

## Feature Description

Implement comprehensive unit tests for the core functionality files in the guard-tool project. This includes creating unit test files for the filesystem, manager, and registry packages to achieve 60%+ code coverage and ensure reliability of critical file protection operations.

## User Story

As a developer working on guard-tool
I want comprehensive unit tests for core functionality
So that I can confidently refactor code, catch regressions early, and ensure reliable file protection operations across different platforms and scenarios.

## Problem Statement

The guard-tool project currently has minimal unit test coverage for its core business logic layers. While extensive shell integration tests exist (46 test files), the lack of unit tests makes it difficult to:
- Test edge cases and error conditions in isolation
- Validate cross-platform behavior (macOS vs Linux immutable flags)
- Ensure thread safety in concurrent scenarios
- Catch regressions during refactoring
- Achieve reliable CI/CD with fast feedback loops

## Solution Statement

Create comprehensive unit test suites for the three core packages (filesystem, manager, registry) using Go's built-in testing framework with table-driven tests. Focus on testing business logic, error handling, platform-specific behavior, and thread safety while maintaining the existing integration test coverage.

## Feature Metadata

**Feature Type**: Enhancement
**Estimated Complexity**: Medium
**Primary Systems Affected**: internal/filesystem, internal/manager, internal/registry
**Dependencies**: Go testing framework (built-in), temporary file utilities

---

## CONTEXT REFERENCES

### Relevant Codebase Files IMPORTANT: YOU MUST READ THESE FILES BEFORE IMPLEMENTING!

- `internal/filesystem/filesystem.go` (lines 1-330) - Why: Contains all FileSystem methods that need comprehensive testing
- `internal/manager/manager.go` (lines 1-200) - Why: Core Manager struct and business logic orchestration
- `internal/manager/files.go` (lines 1-300) - Why: File operation business logic with complex error handling
- `internal/manager/collections.go` (lines 1-400) - Why: Collection operations with conflict detection logic
- `internal/manager/warnings.go` (lines 1-150) - Why: Warning system that needs validation
- `internal/registry/registry.go` (lines 1-100) - Why: Registry interface and factory functions
- `internal/registry/registry_impl.go` (lines 1-200) - Why: Main registry implementation
- `internal/registry/serializer.go` (lines 1-100) - Why: YAML serialization logic
- `internal/registry/file_repository.go` (lines 1-150) - Why: Thread-safe file operations
- `internal/registry/collection_repository.go` (lines 1-200) - Why: Thread-safe collection operations
- `internal/registry/registry_test.go` (lines 1-100) - Why: Existing test patterns to follow
- `internal/security/security_test.go` (lines 1-50) - Why: Integration test patterns for reference
- `justfile` (lines 53, 68) - Why: Test execution commands in CI pipeline
- `go.mod` - Why: Current dependencies (no new dependencies needed)

### New Files to Create

- `internal/filesystem/filesystem_test.go` - Comprehensive FileSystem method testing
- `internal/manager/manager_test.go` - Manager business logic and orchestration testing  
- `internal/registry/registry_test.go` - Enhanced registry operations testing (extend existing)

### Relevant Documentation YOU SHOULD READ THESE BEFORE IMPLEMENTING!

- [Go Testing Package](https://pkg.go.dev/testing)
  - Specific section: Table-driven tests, t.TempDir(), t.Parallel()
  - Why: Standard patterns for Go unit testing
- [Go Testing Best Practices](https://go.dev/doc/tutorial/add-a-test)
  - Specific section: Test organization and naming conventions
  - Why: Ensures tests follow Go community standards

### Patterns to Follow

**Table-Driven Test Pattern** (from existing registry_test.go):
```go
func TestValidateConfig(t *testing.T) {
    tests := []struct {
        name    string
        config  Config
        wantErr bool
    }{
        {
            name: "valid config",
            config: Config{
                GuardFileMode: "0640",
                GuardOwner:    "user",
                GuardGroup:    "group",
            },
            wantErr: false,
        },
        // ... more test cases
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            err := ValidateConfig(tt.config)
            if (err != nil) != tt.wantErr {
                t.Errorf("ValidateConfig() error = %v, wantErr %v", err, tt.wantErr)
            }
        })
    }
}
```

**Temporary Directory Pattern** (from security_test.go):
```go
func TestSecurity_Integration(t *testing.T) {
    tmpDir := t.TempDir()
    registryPath := tmpDir + "/.guardfile"
    // ... test logic
}
```

**Error Handling Pattern**:
```go
if err != nil {
    t.Fatalf("operation failed: %v", err)
}
```

**Naming Conventions**:
- Test functions: Descriptive names like `TestFileExists`, `TestChmodNonExistent`, `TestGetFileInfoNonExistent`
- Test cases: descriptive names like "valid config", "empty guard mode"
- Variables: `wantErr bool`, `got`, `want` for expected values

---

## IMPLEMENTATION PLAN

### Phase 1: Foundation Setup

Set up testing infrastructure using Go's built-in testing framework for comprehensive unit testing.

**Tasks:**
- Create test helper utilities for common operations
- Set up test fixtures and mock data structures

### Phase 2: Filesystem Package Testing

Implement comprehensive tests for the FileSystem struct covering all public methods, platform-specific behavior, and error conditions.

**Tasks:**
- Test basic file operations (exists, info, permissions)
- Test platform-specific immutable flag operations
- Test directory operations with symlink handling
- Test permission operations with privilege scenarios

### Phase 3: Registry Package Testing

Enhance existing registry tests and add comprehensive coverage for all repository operations, thread safety, and YAML serialization.

**Tasks:**
- Extend existing registry_test.go with comprehensive coverage
- Test thread-safe operations with concurrent access
- Test YAML serialization edge cases and error recovery

### Phase 4: Manager Package Testing

Implement tests for the Manager business logic layer, focusing on complex operations, conflict detection, and error handling.

**Tasks:**
- Test manager initialization and registry lifecycle
- Test file operations with business rule validation
- Test collection operations with conflict detection
- Test warning/error accumulation and aggregation

---

## STEP-BY-STEP TASKS

IMPORTANT: Execute every task in order, top to bottom. Each task is atomic and independently testable.

### ADD go.mod

- **IMPLEMENT**: No new dependencies needed - use Go's built-in testing framework
- **PATTERN**: Standard Go testing without external dependencies
- **IMPORTS**: `testing`, `os`, `path/filepath` only
- **GOTCHA**: Use t.Errorf, t.Fatalf instead of assert functions
- **VALIDATE**: `go test ./... -v`

### CREATE internal/filesystem/filesystem_test.go

- **IMPLEMENT**: Comprehensive FileSystem method testing with table-driven tests
- **PATTERN**: Table-driven tests from `internal/registry/registry_test.go:1-50`
- **IMPORTS**: `testing`, `os`, `path/filepath`, `runtime`
- **GOTCHA**: Platform-specific immutable flag tests need runtime.GOOS checks
- **VALIDATE**: `go test ./internal/filesystem -v`

#### Test Structure for filesystem_test.go:
```go
// Test basic operations
func TestFileExists(t *testing.T)
func TestGetFileInfo(t *testing.T)
func TestGetFileInfoNonExistent(t *testing.T)
func TestHasRootPrivileges(t *testing.T)

// Test permission operations
func TestChmod(t *testing.T)
func TestChmodNonExistent(t *testing.T)
func TestChownCurrentUser(t *testing.T)
func TestChownInvalidUser(t *testing.T)
func TestChownNonExistent(t *testing.T)
func TestChgrpCurrentGroup(t *testing.T)
func TestChgrpInvalidGroup(t *testing.T)
func TestChgrpNonExistent(t *testing.T)

// Test combined operations
func TestApplyPermissions(t *testing.T)
func TestApplyPermissionsEmptyOwnerGroup(t *testing.T)
func TestApplyPermissionsOperationOrder(t *testing.T)
func TestRestorePermissions(t *testing.T)

// Test batch operations
func TestCheckFilesExist(t *testing.T)
func TestCheckFilesExistEmpty(t *testing.T)

// Test immutable flags
func TestImmutableFlagMethodsExist(t *testing.T)
func TestSetImmutableBehavior(t *testing.T)
func TestClearImmutableBehavior(t *testing.T)
func TestIsImmutableBehavior(t *testing.T)
func TestImmutableFlagNonExistentFile(t *testing.T)

// Test error message format
func TestErrorMessageFormat(t *testing.T)
```

### UPDATE internal/registry/registry_test.go

- **IMPLEMENT**: Extend existing tests with comprehensive coverage for all registry operations
- **PATTERN**: Existing table-driven test structure in same file
- **IMPORTS**: Add `sync`, `time` (no external dependencies)
- **GOTCHA**: Thread safety tests need goroutines and race condition detection
- **VALIDATE**: `go test ./internal/registry -v -race`

#### Additional Test Functions to Add:
```go
// Registry Creation
func TestNewRegistryWithValidDefaults(t *testing.T)
func TestNewRegistryWithInvalidMode(t *testing.T)
func TestNewRegistryWithNilDefaults(t *testing.T)
func TestNewRegistryWithValidModeRange(t *testing.T)
func TestNewRegistryWithEmptyOwnerGroup(t *testing.T)
func TestNewRegistryFileAlreadyExists(t *testing.T)
func TestLoadRegistryMissingFile(t *testing.T)

// YAML Persistence
func TestSaveAndLoadRoundTrip(t *testing.T)
func TestLoadCorruptedYAML(t *testing.T)
func TestLoadInvalidModeInYAML(t *testing.T)

// File Operations
func TestRegisterFile(t *testing.T)
func TestRegisterFileDuplicate(t *testing.T)
func TestUnregisterFile(t *testing.T)
func TestUnregisterFileNotRegistered(t *testing.T)
func TestUnregisterFileRemovesFromAllCollections(t *testing.T)
func TestGetRegisteredFiles(t *testing.T)
func TestSetAndGetFileMetadata(t *testing.T)
func TestGetRegisteredFileConfig(t *testing.T)

// Collection Operations
func TestRegisterCollection(t *testing.T)
func TestRegisterCollectionDuplicate(t *testing.T)
func TestUnregisterCollection(t *testing.T)
func TestUnregisterCollectionNotRegistered(t *testing.T)
func TestGetRegisteredCollections(t *testing.T)
func TestCountFilesInCollection(t *testing.T)
func TestAddRegisteredFilesToRegisteredCollections(t *testing.T)
func TestAddRegisteredFilesToRegisteredCollectionsIdempotent(t *testing.T)
func TestAddRegisteredFilesToRegisteredCollectionsNonExistentCollection(t *testing.T)
func TestAddRegisteredFilesToRegisteredCollectionsNonRegisteredFile(t *testing.T)
func TestRemoveRegisteredFilesFromRegisteredCollections(t *testing.T)
func TestSetAndGetCollectionGuard(t *testing.T)

// Thread Safety
func TestConcurrentRegisterFile(t *testing.T)
func TestConcurrentSave(t *testing.T)
func TestConcurrentReadWrite(t *testing.T)

// Default Configuration
func TestGetAndSetDefaultConfig(t *testing.T)
func TestSetDefaultFileModeWithSpecialBits(t *testing.T)

// LastToggle
func TestLastToggle(t *testing.T)
func TestLastTogglePersistence(t *testing.T)
```

### CREATE internal/manager/manager_test.go

- **IMPLEMENT**: Comprehensive Manager business logic testing with mock dependencies
- **PATTERN**: Integration test pattern from `internal/security/security_test.go:1-50`
- **IMPORTS**: `testing`, `os`, `path/filepath` (no external dependencies)
- **GOTCHA**: Manager depends on security and filesystem layers - need mocking or test doubles
- **VALIDATE**: `go test ./internal/manager -v`

#### Test Structure for manager_test.go:
```go
// Helper Functions
func setupTestManager(t *testing.T) (*Manager, string, func())
func createTestFile(t *testing.T, dir string, name string, mode os.FileMode) string

// Test manager lifecycle
func TestNewManager(t *testing.T)
func TestLoadRegistry(t *testing.T)
func TestInitializeRegistry(t *testing.T)
func TestSaveRegistry(t *testing.T)

// Test argument resolution
func TestResolveArgument(t *testing.T)
func TestResolveArguments(t *testing.T)

// Test file operations
func TestAddFiles(t *testing.T)
func TestRemoveFiles(t *testing.T)
func TestToggleFiles(t *testing.T)
func TestEnableFiles(t *testing.T)
func TestDisableFiles(t *testing.T)

// Test collection operations
func TestAddCollections(t *testing.T)
func TestRemoveCollections(t *testing.T)
func TestToggleCollections(t *testing.T)
func TestEnableCollections(t *testing.T)
func TestDisableCollections(t *testing.T)

// Test conflict detection
func TestToggleCollections_ConflictDetection(t *testing.T)

// Test warning/error handling
func TestWarningAccumulation(t *testing.T)
func TestErrorHandling(t *testing.T)
func TestUninstallVerification(t *testing.T)
func TestWarningFileMissingContainsCleanupSuggestion(t *testing.T)
func TestCollectionEnableWithMissingFilesWarning(t *testing.T)
func TestCollectionDisableWithMissingFilesWarning(t *testing.T)
func TestCollectionToggleWithMissingFilesWarning(t *testing.T)
func TestWarningShowsRelativePathsNotAbsolute(t *testing.T)
func TestDestroyCollectionWithMissingFilesWarning(t *testing.T)

// Test cleanup operations
func TestCleanup(t *testing.T)
func TestReset(t *testing.T)
func TestDestroy(t *testing.T)
```

### UPDATE justfile

- **IMPLEMENT**: Add test coverage reporting and race detection to CI pipeline
- **PATTERN**: Existing CI commands in justfile:53-71
- **IMPORTS**: No imports needed
- **GOTCHA**: Coverage reports should exclude vendor and test files
- **VALIDATE**: `just ci-quiet`

#### Add to justfile after line 53:
```bash
# Run tests with coverage
test-coverage:
    @echo ""
    go test ./... -coverprofile=coverage.out -covermode=atomic
    go tool cover -html=coverage.out -o coverage.html
    go tool cover -func=coverage.out | grep total
    @echo ""

# Run tests with race detection
test-race:
    @echo ""
    go test ./... -race -v
    @echo ""
```

### UPDATE .gitignore

- **IMPLEMENT**: Add test artifacts to gitignore
- **PATTERN**: Existing .gitignore patterns
- **IMPORTS**: No imports needed
- **GOTCHA**: Coverage files should not be committed
- **VALIDATE**: `git status` should not show coverage files

#### Add to .gitignore:
```
# Test artifacts
coverage.out
coverage.html
*.test
```

---

## TESTING STRATEGY

### Unit Tests

**Scope**: Test individual methods and functions in isolation with controlled inputs and mocked dependencies.

**Framework**: Go's built-in `testing` package only.

**Coverage Target**: 60%+ for each package, focusing on:
- All public methods and functions
- Error handling paths
- Edge cases and boundary conditions
- Platform-specific behavior

**Test Organization**:
- One test file per source file (`filesystem.go` → `filesystem_test.go`)
- Table-driven tests for multiple input scenarios
- Separate test functions for each public method
- Helper functions for common test setup

### Integration Tests

**Scope**: Test interactions between components using real filesystem operations in temporary directories.

**Existing Coverage**: 46 shell integration tests already provide comprehensive end-to-end coverage.

**Unit Test Integration**: Focus on component boundaries and data flow between layers.

### Edge Cases

**Filesystem Package**:
- Permission operations without root privileges
- Immutable flag operations on unsupported platforms
- Directory operations with symlinks and circular references
- File operations on non-existent files

**Manager Package**:
- Conflict detection with multiple collections sharing files
- Registry corruption and recovery scenarios
- Missing file handling across different operations
- Warning accumulation and aggregation logic

**Registry Package**:
- Concurrent access with multiple goroutines
- YAML serialization with malformed data
- File-collection relationship consistency
- Configuration validation boundary conditions

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
go test ./internal/filesystem -v
go test ./internal/manager -v  
go test ./internal/registry -v
go test ./... -v
```

### Level 3: Race Detection

```bash
go test ./... -race -v
```

### Level 4: Coverage Analysis

```bash
go test ./... -coverprofile=coverage.out -covermode=atomic
go tool cover -func=coverage.out | grep total
```

### Level 5: Integration Tests

```bash
just build
cp build/guard tests/guard
cd tests && ./run-all-tests.sh
```

### Level 6: Full CI Pipeline

```bash
just ci-quiet
```

---

## ACCEPTANCE CRITERIA

- [ ] Unit tests created for all public methods in filesystem, manager, and registry packages
- [ ] Test coverage reaches 60%+ for each target package
- [ ] All table-driven tests follow Go conventions and existing patterns
- [ ] Thread safety tests validate concurrent access scenarios
- [ ] Platform-specific tests handle macOS vs Linux differences
- [ ] Error handling tests cover all failure paths
- [ ] All validation commands pass with zero errors
- [ ] No regressions in existing shell integration tests
- [ ] Test execution time remains under 30 seconds for fast CI feedback
- [ ] Race detection tests pass without data races
- [ ] Coverage reports exclude test files and focus on source code

---

## COMPLETION CHECKLIST

- [ ] filesystem_test.go created with comprehensive method coverage
- [ ] registry_test.go extended with additional test functions
- [ ] manager_test.go created with business logic testing
- [ ] justfile updated with coverage and race detection commands
- [ ] .gitignore updated to exclude test artifacts
- [ ] All unit tests pass individually and collectively
- [ ] Race detection tests pass without warnings
- [ ] Coverage target of 60%+ achieved for each package
- [ ] Integration tests continue to pass
- [ ] CI pipeline executes successfully with new tests

---

## NOTES

**Platform Considerations**: Immutable flag operations behave differently on macOS (chflags) vs Linux (chattr). Tests must use runtime.GOOS to conditionally test platform-specific behavior.

**Thread Safety**: Registry operations use sync.RWMutex for thread safety. Tests should validate concurrent access patterns and race conditions using `go test -race`.

**Mock Strategy**: Manager tests may require mocking filesystem operations to test error conditions without requiring root privileges or creating actual files.

**Coverage Focus**: Prioritize testing business logic and error handling over simple getters/setters. Focus on complex operations like conflict detection, permission management, and YAML serialization.

**Test Performance**: Keep unit tests fast (< 30 seconds total) for rapid CI feedback. Use t.Parallel() for independent tests and t.TempDir() for isolated filesystem operations.

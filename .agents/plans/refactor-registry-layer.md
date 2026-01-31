# Feature: Refactor Registry Layer

The following plan should be complete, but its important that you validate documentation and codebase patterns and task sanity before you start implementing.

Pay special attention to naming of existing utils types and models. Import from the right files etc.

## Feature Description

Refactor the internal/registry layer to improve code organization, reduce complexity, and ensure the registry follows Go best practices. The registry is the YAML-based state persistence system that tracks protected files, collections, and folders. This refactoring will make the code more maintainable while preserving all existing functionality.

## User Story

As a developer maintaining the guard-tool codebase
I want the registry layer to follow Go best practices and have clear separation of concerns
So that the code is easier to understand, test, and extend with new features

## Problem Statement

The current registry implementation has several issues:
1. **Monolithic Registry struct** - Single large struct with 370+ lines handling files, collections, folders, and config
2. **Mixed responsibilities** - Registry handles validation, serialization, business logic, and data access
3. **Large method surface** - 30+ methods on Registry struct violating single responsibility principle
4. **Duplicated validation logic** - Octal string parsing and config validation scattered across files
5. **Inconsistent error handling** - Mix of error patterns and validation approaches
6. **Testing gaps** - No unit tests, only integration tests via shell scripts
7. **Concurrency concerns** - Single RWMutex protecting all operations regardless of data type

## Solution Statement

Refactor the registry layer using Go best practices:
1. **Interface segregation** - Split into focused interfaces (FileRegistry, CollectionRegistry, etc.)
2. **Single responsibility** - Separate validation, serialization, and business logic
3. **Dependency injection** - Use interfaces for testability and modularity
4. **Consistent error handling** - Custom error types and validation patterns
5. **Improved concurrency** - Fine-grained locking and thread-safe operations
6. **Comprehensive testing** - Unit tests for all components

## Feature Metadata

**Feature Type**: Refactor
**Estimated Complexity**: High
**Primary Systems Affected**: Registry layer, Security wrapper, Manager layer
**Dependencies**: No external dependencies, internal refactoring only

---

## CONTEXT REFERENCES

### Relevant Codebase Files IMPORTANT: YOU MUST READ THESE FILES BEFORE IMPLEMENTING!

- `internal/registry/registry.go` (lines 1-370) - Why: Main Registry struct and core functionality to refactor
- `internal/registry/file_entry.go` (lines 1-200) - Why: File operations that need interface extraction
- `internal/registry/collection_entry.go` (lines 1-400) - Why: Collection operations that need interface extraction  
- `internal/registry/folder_entry.go` (lines 1-150) - Why: Folder operations that need interface extraction
- `internal/security/security.go` (lines 1-200) - Why: Security wrapper that depends on Registry interface
- `internal/manager/manager.go` (lines 50-60) - Why: Manager layer that creates Registry instances
- `tests/helpers.sh` (lines 116-200) - Why: Test patterns for registry validation
- `tests/test-guardfile-parsers.sh` - Why: Expected YAML structure and parsing behavior

### New Files to Create

- `internal/registry/interfaces.go` - Registry interfaces and common types
- `internal/registry/config.go` - Configuration management with validation
- `internal/registry/file_repository.go` - File operations implementation
- `internal/registry/collection_repository.go` - Collection operations implementation
- `internal/registry/folder_repository.go` - Folder operations implementation
- `internal/registry/serializer.go` - YAML serialization logic
- `internal/registry/validator.go` - Validation utilities and error types
- `internal/registry/registry_impl.go` - Main registry implementation
- `internal/registry/registry_test.go` - Unit tests for registry components

### Relevant Documentation YOU SHOULD READ THESE BEFORE IMPLEMENTING!

- [Go Interfaces Best Practices](https://go.dev/doc/effective_go#interfaces)
  - Specific section: Interface segregation and composition
  - Why: Required for splitting Registry into focused interfaces
- [Go Error Handling](https://go.dev/blog/error-handling-and-go)
  - Specific section: Custom error types and wrapping
  - Why: Needed for consistent error handling patterns
- [Go Concurrency Patterns](https://go.dev/blog/pipelines)
  - Specific section: sync.RWMutex usage patterns
  - Why: Required for thread-safe registry operations

### Patterns to Follow

**Error Handling Pattern:**
```go
// From existing codebase - internal/registry/registry.go:95
if _, err := os.Stat(registryPath); os.IsNotExist(err) {
    return nil, fmt.Errorf("registry file does not exist: %s", registryPath)
}
```

**Validation Pattern:**
```go
// From existing codebase - internal/registry/registry.go:227
func validateConfig(config Config) error {
    if config.GuardFileMode == "" {
        return fmt.Errorf("guard_mode is required in config")
    }
    // ... more validation
}
```

**Mutex Pattern:**
```go
// From existing codebase - internal/registry/registry.go:147
func (r *Registry) Load() error {
    r.mu.Lock()
    defer r.mu.Unlock()
    // ... implementation
}
```

**YAML Serialization Pattern:**
```go
// From existing codebase - internal/registry/registry.go:193
var registryData RegistryData
registryData.Config = r.config
data, err := yaml.Marshal(&registryData)
```

---

## IMPLEMENTATION PLAN

### Phase 1: Foundation - Interfaces and Types

Extract common interfaces and types that will be used across all registry components.

**Tasks:**
- Define core interfaces for registry operations
- Create common error types and validation utilities
- Establish configuration management patterns

### Phase 2: Component Separation

Split the monolithic Registry into focused components with single responsibilities.

**Tasks:**
- Extract file operations into FileRepository
- Extract collection operations into CollectionRepository  
- Extract folder operations into FolderRepository
- Create serialization component for YAML operations

### Phase 3: Registry Implementation

Create the main registry implementation that composes the separated components.

**Tasks:**
- Implement composite registry that delegates to components
- Ensure thread safety with appropriate locking strategies
- Maintain backward compatibility with existing API

### Phase 4: Testing & Validation

Add comprehensive unit tests and validate the refactoring preserves functionality.

**Tasks:**
- Create unit tests for all components
- Validate against existing shell test expectations
- Performance testing for concurrent operations

---

## STEP-BY-STEP TASKS

IMPORTANT: Execute every task in order, top to bottom. Each task is atomic and independently testable.

### CREATE internal/registry/interfaces.go

- **IMPLEMENT**: Core registry interfaces with single responsibilities covering all 29 methods used by security.go
- **PATTERN**: Interface segregation principle - separate read/write operations
- **IMPORTS**: `os`, `context` for future extensibility
- **GOTCHA**: Must include ALL methods that security.go calls - verified against actual usage
- **VALIDATE**: `go build ./internal/registry`

```go
// FileReader interface for read-only file operations
type FileReader interface {
    IsRegisteredFile(path string) bool
    GetRegisteredFiles() []string
    GetRegisteredFileConfig(path string) (string, string, os.FileMode, bool, error)
    GetRegisteredFileGuard(path string) (bool, error)
}

// FileWriter interface for file modification operations  
type FileWriter interface {
    RegisterFile(path string, mode os.FileMode, owner, group string) error
    UnregisterFile(path string, ignoreMissing bool) error
    SetRegisteredFileGuard(path string, guard bool) error
    SetRegisteredFileConfig(path string, mode os.FileMode, owner, group string, guard bool) error
    RemoveRegisteredFileFromAllRegisteredCollections(path string)
}

// CollectionReader interface for read-only collection operations
type CollectionReader interface {
    IsRegisteredCollection(name string) bool
    GetRegisteredCollections() []string
    GetRegisteredCollectionGuard(name string) (bool, error)
    GetRegisteredCollectionFiles(name string) ([]string, error)
}

// CollectionWriter interface for collection modification operations
type CollectionWriter interface {
    RegisterCollection(name string, files []string) error
    UnregisterCollection(name string, ignoreMissing bool) error
    SetRegisteredCollectionGuard(name string, guard bool) error
    AddRegisteredFilesToRegisteredCollections(collections, files []string) error
    RemoveRegisteredFilesFromRegisteredCollections(collections, files []string) error
}

// FolderReader interface for read-only folder operations
type FolderReader interface {
    IsRegisteredFolder(name string) bool
    GetFolderGuard(name string) (bool, error)
    GetFolder(name string) *Folder
}

// FolderWriter interface for folder modification operations
type FolderWriter interface {
    RegisterFolder(name, path string) error
    SetFolderGuard(name string, guard bool) error
}

// ConfigReader interface for configuration read operations
type ConfigReader interface {
    GetDefaultFileMode() os.FileMode
    GetDefaultFileOwner() string
    GetDefaultFileGroup() string
}

// ConfigWriter interface for configuration write operations
type ConfigWriter interface {
    SetDefaultFileMode(mode os.FileMode) error
    SetDefaultFileOwner(owner string)
    SetDefaultFileGroup(group string)
}

// Persister interface for save operations
type Persister interface {
    Save() error
}

// Registry interface combines all operations - this is what Security.go will use
type Registry interface {
    FileReader
    FileWriter
    CollectionReader
    CollectionWriter
    FolderReader
    FolderWriter
    ConfigReader
    ConfigWriter
    Persister
}

// Folder struct must be preserved for GetFolder return type
type Folder struct {
    Name  string `yaml:"name"`
    Path  string `yaml:"path"`
    Guard bool   `yaml:"guard"`
}
```

### CREATE internal/registry/validator.go

- **IMPLEMENT**: Validation utilities and custom error types
- **PATTERN**: Centralized validation with custom error types - `internal/registry/registry.go:227`
- **IMPORTS**: `fmt`, `errors`, `os`, `strconv`, `strings`
- **GOTCHA**: Preserve exact validation logic from existing code
- **VALIDATE**: `go test ./internal/registry -run TestValidator`

```go
// ValidationError represents validation failures
type ValidationError struct {
    Field   string
    Value   string
    Message string
}

// OctalStringToFileMode validates and converts octal strings
func OctalStringToFileMode(value string) (os.FileMode, error)

// ValidateConfig validates registry configuration
func ValidateConfig(config Config) error
```

### CREATE internal/registry/config.go

- **IMPLEMENT**: Configuration management with validation
- **PATTERN**: Config struct and validation from `internal/registry/registry.go:27-32`
- **IMPORTS**: `fmt`, `os`, `strings`, `sync`
- **GOTCHA**: Maintain exact YAML field names for backward compatibility
- **VALIDATE**: `go test ./internal/registry -run TestConfig`

### CREATE internal/registry/serializer.go

- **IMPLEMENT**: YAML serialization and file I/O operations
- **PATTERN**: Atomic file operations from `internal/registry/registry.go:193-225`
- **IMPORTS**: `fmt`, `os`, `gopkg.in/yaml.v3`
- **GOTCHA**: Preserve exact YAML structure for compatibility with existing .guardfile
- **VALIDATE**: `go test ./internal/registry -run TestSerializer`

### CREATE internal/registry/file_repository.go

- **IMPLEMENT**: File operations extracted from file_entry.go
- **PATTERN**: Method signatures from `internal/registry/file_entry.go`
- **IMPORTS**: `fmt`, `os`, `sync`
- **GOTCHA**: Maintain thread safety with appropriate locking
- **VALIDATE**: `go test ./internal/registry -run TestFileRepository`

### CREATE internal/registry/collection_repository.go

- **IMPLEMENT**: Collection operations extracted from collection_entry.go
- **PATTERN**: Method signatures from `internal/registry/collection_entry.go`
- **IMPORTS**: `fmt`, `os`, `strings`, `sync`
- **GOTCHA**: Preserve collection-to-file relationship management
- **VALIDATE**: `go test ./internal/registry -run TestCollectionRepository`

### CREATE internal/registry/folder_repository.go

- **IMPLEMENT**: Folder operations extracted from folder_entry.go
- **PATTERN**: Method signatures from `internal/registry/folder_entry.go`
- **IMPORTS**: `fmt`, `sync`
- **GOTCHA**: Maintain @ prefix naming convention for folders
- **VALIDATE**: `go test ./internal/registry -run TestFolderRepository`

### CREATE internal/registry/registry_impl.go

- **IMPLEMENT**: Main registry implementation composing all components and implementing Registry interface
- **PATTERN**: Delegation pattern - compose FileRepository, CollectionRepository, etc.
- **IMPORTS**: All repository interfaces and implementations
- **GOTCHA**: Must implement ALL 29 methods from Registry interface for security.go compatibility
- **VALIDATE**: `go build ./internal/registry`

```go
// RegistryImpl composes all repository components
type RegistryImpl struct {
    fileRepo       FileRepository
    collectionRepo CollectionRepository
    folderRepo     FolderRepository
    config         ConfigManager
    serializer     Serializer
}

// Implement all Registry interface methods by delegating to components
// File operations delegate to fileRepo
// Collection operations delegate to collectionRepo
// Folder operations delegate to folderRepo
// Config operations delegate to config
// Save operation delegates to serializer
```

### UPDATE internal/registry/registry.go

- **REFACTOR**: Keep only factory functions and backward compatibility
- **PATTERN**: Preserve `NewRegistry` and `LoadRegistry` function signatures
- **IMPORTS**: Update to use new implementation
- **GOTCHA**: Existing code depends on these exact function signatures
- **VALIDATE**: `go build ./internal/security`

### CREATE internal/registry/registry_test.go

- **IMPLEMENT**: Comprehensive unit tests for all components
- **PATTERN**: Table-driven tests following Go conventions
- **IMPORTS**: `testing`, `os`, `path/filepath`, `gopkg.in/yaml.v3`
- **GOTCHA**: Test concurrent access patterns and error conditions
- **VALIDATE**: `go test ./internal/registry -v -race`

### UPDATE internal/security/security.go

- **REFACTOR**: Update to use Registry interface instead of concrete type
- **PATTERN**: Dependency injection - accept Registry interface instead of *registry.Registry
- **IMPORTS**: Update registry imports to use new interfaces
- **GOTCHA**: All 29 methods must work exactly the same - zero behavior changes allowed
- **VALIDATE**: `go build ./internal/security && go test ./internal/security`

```go
// Security struct now uses Registry interface
type Security struct {
    registry Registry  // Changed from *registry.Registry
}

// All 29 method calls remain identical:
// s.registry.IsRegisteredFile(cleanPath)
// s.registry.GetRegisteredFileGuard(cleanPath)
// ... etc - no changes to method calls
```

### UPDATE internal/manager/manager.go

- **REFACTOR**: Update registry creation to use new implementation
- **PATTERN**: Factory pattern usage from `internal/manager/manager.go:54`
- **IMPORTS**: Update registry imports
- **GOTCHA**: Manager depends on Security wrapper, not Registry directly
- **VALIDATE**: `go build ./internal/manager`

### REMOVE internal/registry/file_entry.go

- **REMOVE**: File after functionality moved to file_repository.go
- **PATTERN**: Clean removal of deprecated files
- **IMPORTS**: N/A
- **GOTCHA**: Ensure no imports reference this file
- **VALIDATE**: `go build ./...`

### REMOVE internal/registry/collection_entry.go

- **REMOVE**: File after functionality moved to collection_repository.go
- **PATTERN**: Clean removal of deprecated files
- **IMPORTS**: N/A
- **GOTCHA**: Ensure no imports reference this file
- **VALIDATE**: `go build ./...`

### REMOVE internal/registry/folder_entry.go

- **REMOVE**: File after functionality moved to folder_repository.go
- **PATTERN**: Clean removal of deprecated files
- **IMPORTS**: N/A
- **GOTCHA**: Ensure no imports reference this file
- **VALIDATE**: `go build ./...`

---

## TESTING STRATEGY

### Unit Tests

Design unit tests with fixtures and assertions following Go testing conventions:

- **File Repository Tests**: Test file registration, unregistration, guard operations
- **Collection Repository Tests**: Test collection creation, file management, guard operations
- **Folder Repository Tests**: Test folder registration and guard operations
- **Serializer Tests**: Test YAML marshaling/unmarshaling with various data structures
- **Validator Tests**: Test all validation functions with valid/invalid inputs
- **Config Tests**: Test configuration management and defaults
- **Integration Tests**: Test composed registry with all components

### Concurrency Tests

- **Race Condition Tests**: Use `go test -race` to detect data races
- **Concurrent Access Tests**: Multiple goroutines accessing registry simultaneously
- **Lock Contention Tests**: Verify appropriate locking granularity

### Compatibility Tests

- **Backward Compatibility**: Ensure existing .guardfile format still works
- **API Compatibility**: Verify all existing public methods still work
- **Shell Test Compatibility**: Existing shell tests should pass unchanged

---

## VALIDATION COMMANDS

Execute every command to ensure zero regressions and 100% feature correctness.

### Level 1: Syntax & Style

```bash
go fmt ./internal/registry/...
golangci-lint run ./internal/registry/...
```

### Level 2: Unit Tests

```bash
go test ./internal/registry -v -race -cover
go test ./internal/security -v
go test ./internal/manager -v
```

### Level 3: Integration Tests

```bash
go build -o build/guard ./cmd/guard
cp build/guard tests/guard
cd tests && ./run-all-tests.sh
```

### Level 4: Manual Validation

```bash
# Test basic registry operations
./guard init --mode 0640 --owner $(whoami) --group $(id -gn)
./guard add test.txt
./guard show
./guard enable test.txt
./guard disable test.txt
```

### Level 5: Performance Validation

```bash
go test ./internal/registry -bench=. -benchmem
go test ./internal/registry -race -count=100
```

---

## ACCEPTANCE CRITERIA

- [ ] Registry interface includes all 29 methods that security.go calls
- [ ] Registry layer follows single responsibility principle with focused interfaces
- [ ] All existing functionality preserved - no behavior changes
- [ ] Comprehensive unit test coverage (80%+) for all registry components
- [ ] Thread-safe operations with appropriate locking granularity
- [ ] Consistent error handling with custom error types
- [ ] All existing shell integration tests pass unchanged
- [ ] No performance regression in registry operations
- [ ] Clean separation between validation, serialization, and business logic
- [ ] Backward compatibility with existing .guardfile format
- [ ] Code follows Go best practices and passes all linting checks
- [ ] Security.go compiles and works with new Registry interface without changes

---

## COMPLETION CHECKLIST

- [ ] Registry interface implements all 29 methods used by security.go
- [ ] All interfaces defined with single responsibilities
- [ ] Validation logic centralized and tested
- [ ] File operations extracted to FileRepository
- [ ] Collection operations extracted to CollectionRepository
- [ ] Folder operations extracted to FolderRepository
- [ ] Configuration operations extracted to ConfigManager
- [ ] YAML serialization isolated in Serializer component
- [ ] Main Registry implementation composes all components
- [ ] Comprehensive unit tests for all components
- [ ] All existing shell tests pass
- [ ] Security and Manager layers updated to use new interfaces
- [ ] Old implementation files removed cleanly
- [ ] Performance benchmarks show no regression
- [ ] Race condition tests pass
- [ ] Documentation updated for new architecture

---

## NOTES

### Design Decisions

1. **Interface Segregation**: Split Registry into focused interfaces (FileReader/Writer, CollectionReader/Writer, etc.) to follow single responsibility principle and improve testability.

2. **Component Composition**: Main Registry implementation composes specialized repositories rather than inheriting, allowing for better separation of concerns and easier testing.

3. **Validation Centralization**: Move all validation logic to dedicated validator package to eliminate duplication and ensure consistency.

4. **Thread Safety**: Use fine-grained locking in individual repositories rather than single coarse-grained lock to improve concurrent performance.

5. **Backward Compatibility**: Maintain exact same public API and YAML format to ensure existing code and .guardfile continue working.

### Trade-offs

- **Complexity vs Maintainability**: Increased number of files and interfaces in exchange for better maintainability and testability
- **Performance vs Safety**: Fine-grained locking may have slight overhead but provides better concurrent access patterns
- **Flexibility vs Simplicity**: More complex architecture enables easier extension and testing

### Migration Strategy

The refactoring preserves all existing public APIs, so no changes are required in calling code. The Security wrapper continues to work unchanged, and all shell tests should pass without modification.

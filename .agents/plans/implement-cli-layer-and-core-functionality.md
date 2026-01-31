# Feature: Implement CLI Layer and Core Functionality

The following plan should be complete, but its important that you validate documentation and codebase patterns and task sanity before you start implementing.

Pay special attention to naming of existing utils types and models. Import from the right files etc.

## Feature Description

Build the complete CLI interface and core business logic layers for guard-tool, a file protection system that prevents AI coding agents from modifying files outside the current work scope. The implementation will create a 6-layer architecture (CLI → Manager → Security + Registry + Filesystem → TUI) with the Registry layer already implemented as the foundation.

## User Story

As a developer using AI coding assistants
I want to quickly protect and unprotect files with simple CLI commands
So that AI tools cannot accidentally modify files outside my current work scope

## Problem Statement

The guard-tool project has a complete Registry layer for YAML-based state management but lacks:
- CLI interface to access registry functionality
- Business logic orchestration (Manager layer)
- Path validation and security (Security layer)  
- Filesystem operations for permissions and immutable flags (Filesystem layer)
- Interactive terminal interface (TUI layer)

## Solution Statement

Implement vertical slices through the architecture, starting with core CLI commands (init, add, enable/disable, show) and their corresponding Manager methods. Each command will integrate Registry operations with new Security and Filesystem layers, following exact behavioral specifications from shell integration tests.

## Feature Metadata

**Feature Type**: New Capability
**Estimated Complexity**: High
**Primary Systems Affected**: CLI, Manager, Security, Filesystem layers
**Dependencies**: github.com/spf13/cobra, platform-specific filesystem operations

---

## CONTEXT REFERENCES

### Relevant Codebase Files IMPORTANT: YOU MUST READ THESE FILES BEFORE IMPLEMENTING!

- `internal/registry/registry.go` (lines 35-95) - Why: Registry struct and NewRegistry/LoadRegistry patterns to follow
- `internal/registry/file_entry.go` (lines 19-37) - Why: RegisterFile method pattern for Manager layer
- `internal/registry/collection_entry.go` - Why: Collection management patterns
- `tests/test-init.sh` (lines 1-30) - Why: Exact CLI behavior specifications for guard init
- `tests/helpers.sh` - Why: .guardfile parsing and assertion patterns
- `justfile` (lines 35-45) - Why: Build and validation command patterns

### New Files to Create

- `cmd/guard/main.go` - Main entry point with -i flag handling and Cobra setup
- `cmd/guard/commands/init.go` - guard init command implementation
- `cmd/guard/commands/add.go` - guard add command implementation  
- `cmd/guard/commands/enable.go` - guard enable command implementation
- `cmd/guard/commands/disable.go` - guard disable command implementation
- `cmd/guard/commands/show.go` - guard show command implementation
- `internal/manager/manager.go` - Business logic orchestration layer
- `internal/security/security.go` - Path validation and symlink rejection
- `internal/filesystem/filesystem.go` - Base filesystem operations interface
- `internal/filesystem/filesystem_darwin.go` - macOS-specific immutable flags
- `internal/filesystem/filesystem_linux.go` - Linux-specific immutable flags

### Relevant Documentation YOU SHOULD READ THESE BEFORE IMPLEMENTING!

- [Cobra User Guide](https://cobra.dev/#getting-started)
  - Specific section: Command organization and flag handling
  - Why: Required for proper CLI structure and flag processing
- [Go Build Constraints](https://pkg.go.dev/go/build#hdr-Build_Constraints)
  - Specific section: Platform-specific build tags
  - Why: Needed for filesystem layer platform abstraction

### Patterns to Follow

**Registry Integration Pattern:**
```go
// From internal/registry/registry.go:56
func NewRegistry(registryPath string, defaults *RegistryDefaults, overwrite bool) (*Registry, error) {
    if err := validateRegistryPath(registryPath); err != nil {
        return nil, err
    }
    // Validation before creation
}
```

**Error Handling Pattern:**
```go
// From internal/registry/registry.go:99
if _, err := os.Stat(registryPath); os.IsNotExist(err) {
    return nil, fmt.Errorf("registry file does not exist: %s", registryPath)
}
```

**CLI Exit Code Pattern (from shell tests):**
- Exit 0: Success and warnings
- Exit 1: Actual errors only
- Idempotent operations (add existing file = skip with message, not error)

**Auto-Detection Priority (enable/disable/toggle commands):**
1. Directory on disk → folder
2. File on disk → file  
3. Registered collection name → collection
4. Registered folder with @ prefix → folder
5. Registered file path → file
6. No match → error

---

## IMPLEMENTATION PLAN

### Phase 1: Foundation

Set up CLI infrastructure and core dependencies before implementing business logic.

**Tasks:**
- Add Cobra dependency to go.mod
- Create main entry point with -i flag handling
- Set up command structure and routing
- Create Manager layer interface

### Phase 2: Core Implementation

Implement vertical slices for essential commands with their Manager methods.

**Tasks:**
- Implement guard init command + Manager.Init()
- Implement guard add command + Manager.AddFile()
- Implement guard show command + Manager.GetStatus()
- Create Security layer for path validation
- Create Filesystem layer interface

### Phase 3: Protection Operations

Implement file protection/unprotection with auto-detection logic.

**Tasks:**
- Implement guard enable command + Manager.SetProtection()
- Implement guard disable command + Manager.SetProtection()
- Add auto-detection logic to Manager layer
- Implement platform-specific filesystem operations

### Phase 4: Testing & Validation

Validate implementation against shell test specifications.

**Tasks:**
- Run shell integration tests for implemented commands
- Fix any behavioral discrepancies
- Add unit tests for Manager layer
- Validate cross-platform filesystem operations

---

## STEP-BY-STEP TASKS

IMPORTANT: Execute every task in order, top to bottom. Each task is atomic and independently testable.

### ADD go.mod

- **IMPLEMENT**: Add Cobra dependency to existing go.mod
- **PATTERN**: Existing go.mod format with gopkg.in/yaml.v3
- **IMPORTS**: `github.com/spf13/cobra v1.8.0`
- **GOTCHA**: Use go mod tidy after adding dependency
- **VALIDATE**: `go mod tidy && go mod verify`

### CREATE cmd/guard/main.go

- **IMPLEMENT**: Main entry point with -i flag handling before Cobra routing
- **PATTERN**: PersistentPreRun pattern for global flag processing
- **IMPORTS**: `github.com/spf13/cobra`, `os`, `fmt`
- **GOTCHA**: Handle -i flag before cobra.Execute() to prevent subcommand routing
- **VALIDATE**: `go build -o build/guard ./cmd/guard && ./build/guard --help`

### CREATE internal/manager/manager.go

- **IMPLEMENT**: Manager struct with Registry dependency and core methods
- **PATTERN**: Registry integration from internal/registry/registry.go:35-42
- **IMPORTS**: `internal/registry`, `os`, `fmt`, `sync`
- **GOTCHA**: Use Registry.Load() and Registry.Save() for persistence
- **VALIDATE**: `go build ./internal/manager`

### CREATE internal/security/security.go

- **IMPLEMENT**: Path validation and symlink rejection functions
- **PATTERN**: validateRegistryPath from internal/registry/registry.go:284
- **IMPORTS**: `os`, `path/filepath`, `fmt`
- **GOTCHA**: Reject paths outside current directory, resolve symlinks
- **VALIDATE**: `go test ./internal/security`

### CREATE cmd/guard/commands/init.go

- **IMPLEMENT**: guard init command with mode/owner/group validation
- **PATTERN**: Shell test specifications from tests/test-init.sh:25-30
- **IMPORTS**: `github.com/spf13/cobra`, `internal/manager`, `fmt`, `os`
- **GOTCHA**: Mode stored as 4-digit octal string (0000-0777), exit code 1 on error
- **VALIDATE**: `./build/guard init 644 testuser testgroup && ls -la .guardfile`

### UPDATE cmd/guard/main.go

- **IMPLEMENT**: Add init command to root command
- **PATTERN**: Cobra command registration pattern
- **IMPORTS**: `cmd/guard/commands`
- **GOTCHA**: Import commands package to register init command
- **VALIDATE**: `./build/guard init --help`

### ADD internal/manager/manager.go

- **IMPLEMENT**: Manager.Init() method using Registry.NewRegistry()
- **PATTERN**: NewRegistry call from internal/registry/registry.go:56-95
- **IMPORTS**: `internal/registry`
- **GOTCHA**: Convert mode string to RegistryDefaults, handle overwrite=false
- **VALIDATE**: `./build/guard init 644 user group && echo $?`

### CREATE cmd/guard/commands/add.go

- **IMPLEMENT**: guard add command with file path validation
- **PATTERN**: Idempotent operations (skip existing with message, not error)
- **IMPORTS**: `github.com/spf13/cobra`, `internal/manager`, `os`
- **GOTCHA**: Exit code 0 for "already registered", not 1
- **VALIDATE**: `touch test.txt && ./build/guard add test.txt && echo $?`

### ADD internal/manager/manager.go

- **IMPLEMENT**: Manager.AddFile() method using Registry.RegisterFile()
- **PATTERN**: RegisterFile from internal/registry/file_entry.go:19-37
- **IMPORTS**: `internal/security`, `os`
- **GOTCHA**: Validate path with Security layer, get current file stats
- **VALIDATE**: `./build/guard add test.txt && ./build/guard add test.txt`

### CREATE cmd/guard/commands/show.go

- **IMPLEMENT**: guard show command displaying registry status
- **PATTERN**: Registry getter methods from file_entry.go:67-76
- **IMPORTS**: `github.com/spf13/cobra`, `internal/manager`, `fmt`
- **GOTCHA**: Show files, collections, folders with guard status
- **VALIDATE**: `./build/guard show`

### ADD internal/manager/manager.go

- **IMPLEMENT**: Manager.GetStatus() method aggregating registry data
- **PATTERN**: Registry.GetRegisteredFiles() and similar getters
- **IMPORTS**: None additional
- **GOTCHA**: Format output consistently with expected CLI format
- **VALIDATE**: `./build/guard show | grep -E "(Files|Collections|Folders)"`

### CREATE internal/filesystem/filesystem.go

- **IMPLEMENT**: Filesystem interface with chmod, chown, immutable flag methods
- **PATTERN**: Interface-based abstraction for platform-specific operations
- **IMPORTS**: `os`, `os/user`
- **GOTCHA**: Define interface that both Darwin and Linux can implement
- **VALIDATE**: `go build ./internal/filesystem`

### CREATE internal/filesystem/filesystem_darwin.go

- **IMPLEMENT**: macOS implementation using chflags for immutable flags
- **PATTERN**: Build constraint //go:build darwin
- **IMPORTS**: `os/exec`, `syscall`
- **GOTCHA**: Use chflags with SF_IMMUTABLE flag, requires sudo
- **VALIDATE**: `go build -tags darwin ./internal/filesystem`

### CREATE internal/filesystem/filesystem_linux.go

- **IMPLEMENT**: Linux implementation using ioctl for immutable flags
- **PATTERN**: Build constraint //go:build linux
- **IMPORTS**: `syscall`, `unsafe`
- **GOTCHA**: Use ioctl with FS_IMMUTABLE_FL flag, requires sudo
- **VALIDATE**: `go build -tags linux ./internal/filesystem`

### CREATE cmd/guard/commands/enable.go

- **IMPLEMENT**: guard enable command with auto-detection logic
- **PATTERN**: Auto-detection priority order from requirements
- **IMPORTS**: `github.com/spf13/cobra`, `internal/manager`
- **GOTCHA**: Check directory → file → collection → folder → registered file
- **VALIDATE**: `touch test.txt && ./build/guard add test.txt && ./build/guard enable test.txt`

### CREATE cmd/guard/commands/disable.go

- **IMPLEMENT**: guard disable command with same auto-detection
- **PATTERN**: Mirror enable.go structure with opposite operation
- **IMPORTS**: `github.com/spf13/cobra`, `internal/manager`
- **GOTCHA**: Same auto-detection logic as enable
- **VALIDATE**: `./build/guard disable test.txt`

### ADD internal/manager/manager.go

- **IMPLEMENT**: Manager.SetProtection() with auto-detection and filesystem calls
- **PATTERN**: Auto-detection priority order, Registry.SetRegisteredFileGuard()
- **IMPORTS**: `internal/filesystem`, `internal/security`
- **GOTCHA**: Detect type first, then call appropriate Registry method + filesystem ops
- **VALIDATE**: `./build/guard enable test.txt && stat test.txt`

### UPDATE cmd/guard/main.go

- **IMPLEMENT**: Register all implemented commands (init, add, show, enable, disable)
- **PATTERN**: Command registration in root command
- **IMPORTS**: All command packages
- **GOTCHA**: Import all command packages to register them
- **VALIDATE**: `./build/guard --help | grep -E "(init|add|show|enable|disable)"`

---

## TESTING STRATEGY

### Shell Integration Tests

Execute existing shell tests to validate exact CLI behavior matches specifications.

**Primary Test**: `tests/test-init.sh`
- Validates guard init command with all argument combinations
- Tests exit codes, .guardfile creation, mode validation
- Verifies error handling for invalid modes and existing files

**Test Execution**: `cd tests && ./run-all-tests.sh`

### Unit Tests

Create unit tests for Manager layer business logic following Go testing conventions.

**Test Files**:
- `internal/manager/manager_test.go` - Manager method testing
- `internal/security/security_test.go` - Path validation testing
- `internal/filesystem/filesystem_test.go` - Filesystem operation testing

### Edge Cases

**CLI Edge Cases**:
- Invalid octal modes (999, 888)
- Non-existent files for add command
- Auto-detection ambiguity (file vs collection name collision)
- Permission denied scenarios

**Filesystem Edge Cases**:
- Files without write permissions
- Symlinks and their targets
- Cross-platform immutable flag behavior

---

## VALIDATION COMMANDS

Execute every command to ensure zero regressions and 100% feature correctness.

### Level 1: Syntax & Style

```bash
go fmt ./...
golangci-lint run
gocyclo -over 15 .
gocognit -over 15 .
semgrep scan --config auto
```

### Level 2: Unit Tests

```bash
go test ./...
go test -race ./...
go test -cover ./...
```

### Level 3: Integration Tests

```bash
cd tests && ./run-all-tests.sh
```

### Level 4: Manual Validation

```bash
# Build and test basic workflow
just build
./build/guard init 644 testuser testgroup
echo "test content" > test.txt
./build/guard add test.txt
./build/guard show
./build/guard enable test.txt
./build/guard disable test.txt
```

### Level 5: Cross-Platform Validation

```bash
# Test platform-specific builds
GOOS=darwin go build ./cmd/guard
GOOS=linux go build ./cmd/guard
```

---

## ACCEPTANCE CRITERIA

- [ ] guard init command creates .guardfile with exact format from shell tests
- [ ] guard add command registers files idempotently (no error on duplicate)
- [ ] guard show command displays registry status in readable format
- [ ] guard enable/disable commands work with auto-detection logic
- [ ] All shell integration tests pass with zero failures
- [ ] Exit codes match specifications (0 for success/warnings, 1 for errors)
- [ ] Platform-specific filesystem operations compile on macOS and Linux
- [ ] Manager layer properly orchestrates Registry, Security, and Filesystem layers
- [ ] Path validation prevents directory traversal and symlink exploits
- [ ] Code follows Go best practices and project conventions
- [ ] All validation commands pass with zero errors

---

## COMPLETION CHECKLIST

- [ ] All tasks completed in dependency order
- [ ] Each task validation passed immediately after implementation
- [ ] Shell integration tests execute successfully
- [ ] Manual workflow testing confirms functionality
- [ ] Cross-platform builds succeed
- [ ] No linting or formatting errors
- [ ] Unit test coverage for Manager layer
- [ ] Auto-detection logic handles all specified cases
- [ ] Filesystem operations work on target platforms
- [ ] CLI behavior exactly matches shell test specifications

---

## NOTES

**Architecture Decision**: Vertical slice implementation ensures each command is fully functional before moving to the next, reducing integration complexity.

**Platform Abstraction**: Filesystem layer uses build tags to handle macOS chflags vs Linux ioctl for immutable flags.

**Security Model**: Security layer validates all paths to prevent traversal attacks and symlink exploits before Registry operations.

**Testing Strategy**: Shell tests serve as executable specifications - implementation must match their exact behavioral requirements.

**Reserved Keywords**: Manager layer must validate collection names against reserved words: to, from, add, remove, file, collection, create, destroy, clear, update, uninstall.

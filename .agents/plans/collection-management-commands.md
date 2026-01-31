# Feature: Collection Management Commands

The following plan should be complete, but its important that you validate documentation and codebase patterns and task sanity before you start implementing.

Pay special attention to naming of existing utils types and models. Import from the right files etc.

## Feature Description

Implement comprehensive collection management commands for the guard-tool CLI, including creation, destruction, updating, clearing, configuration management, maintenance operations, folder operations, and auto-detection capabilities. This feature enables users to group files together for batch protection operations and manage them efficiently through intuitive CLI commands.

## User Story

As a developer using AI coding assistants
I want to group related files into collections and manage them with simple commands
So that I can quickly protect/unprotect sets of files without managing them individually

## Problem Statement

Currently, guard-tool has basic collection support but lacks comprehensive management commands. Users need:
- Easy collection creation and destruction
- File addition/removal from collections
- Batch operations on collections
- Configuration management for protection settings
- Maintenance operations for cleanup
- Dynamic folder tracking
- Auto-detection to simplify command usage

## Solution Statement

Implement a complete set of collection management commands following the existing codebase patterns:
- Collection CRUD operations (create, destroy, update, clear)
- Configuration management with validation and warnings
- Maintenance commands for system cleanup
- Folder operations with dynamic file scanning
- Auto-detection priority system for seamless UX
- Comprehensive error handling and output formatting

## Feature Metadata

**Feature Type**: Enhancement
**Estimated Complexity**: High
**Primary Systems Affected**: CLI Commands, Manager Layer, Registry Layer
**Dependencies**: Existing Cobra CLI framework, YAML persistence, filesystem operations

---

## CONTEXT REFERENCES

### Relevant Codebase Files IMPORTANT: YOU MUST READ THESE FILES BEFORE IMPLEMENTING!

- `cmd/guard/commands/create.go` (lines 13-23) - Why: Existing create command pattern to extend
- `cmd/guard/commands/update.go` (entire file) - Why: Update command structure to enhance
- `internal/manager/manager.go` (lines 427-463) - Why: CreateCollection method pattern
- `internal/manager/manager.go` (lines 466-500) - Why: UpdateCollection method pattern
- `internal/registry/collection_entry.go` (lines 20-27) - Why: Collection struct definition
- `internal/registry/collection_entry.go` (lines 32-68) - Why: RegisterCollection method pattern
- `internal/registry/registry.go` (lines 45-85) - Why: Registry configuration management patterns
- `internal/registry/folder_entry.go` (lines 7-11) - Why: Folder struct for folder operations
- `cmd/guard/commands/constants.go` (entire file) - Why: Remove this file, handle keywords inline in commands
- `tests/test-create.sh` (lines 30-50) - Why: Expected output format patterns
- `tests/test-output-specs.sh` (lines 80-100) - Why: Error message format specifications
- `tests/helpers.sh` (lines 268-302) - Why: Collection validation helper functions

### New Files to Create

- `cmd/guard/commands/destroy.go` - Destroy collections command
- `cmd/guard/commands/clear.go` - Clear collections command  
- `cmd/guard/commands/config.go` - Configuration management commands
- `cmd/guard/commands/cleanup.go` - Cleanup maintenance command
- `cmd/guard/commands/reset.go` - Reset maintenance command
- `cmd/guard/commands/uninstall.go` - Uninstall maintenance command
- `internal/manager/folders.go` - Folder operation methods
- `internal/manager/collections.go` - Collection operation methods
- `internal/manager/config.go` - Configuration management methods
- `internal/manager/files.go` - File operation methods

### Relevant Documentation YOU SHOULD READ THESE BEFORE IMPLEMENTING!

- [Cobra Command Documentation](https://cobra.dev/#commands)
  - Specific section: Subcommands and argument parsing
  - Why: Required for implementing complex command structures like `guard config set mode`
- [Go YAML v3 Documentation](https://pkg.go.dev/gopkg.in/yaml.v3)
  - Specific section: Struct tags and marshaling
  - Why: Needed for proper YAML serialization of new structures

### Patterns to Follow

**Command Structure Pattern:**
```go
func NewXxxCmd() *cobra.Command {
    return &cobra.Command{
        Use:   "command <args>",
        Short: "Brief description",
        Long:  "Detailed description with examples",
        Args:  cobra.MinimumNArgs(1),
        RunE:  runXxx,
    }
}
```

**Manager Method Pattern:**
```go
func (m *Manager) MethodName(args) error {
    m.mu.Lock()
    defer m.mu.Unlock()
    
    if m.registry == nil {
        return fmt.Errorf("registry not loaded")
    }
    
    // Business logic
    
    if err := m.registry.Save(); err != nil {
        return fmt.Errorf("failed to save registry: %w", err)
    }
    
    return nil
}
```

**Error Message Format:**
- Errors: `Error: <specific message>`
- Warnings: `Warning: <specific message>`
- Success: `<Action> <count> <item>(s)` (past tense verbs)

**Registry Validation Pattern:**
```go
if !m.registry.IsRegisteredCollection(name) {
    return fmt.Errorf("collection not found: %s", name)
}
```

**Auto-Detection Priority:**
1. Directory on disk → folder
2. File on disk → file  
3. Registered collection name → collection
4. Registered folder name with @ prefix → folder
5. Registered file path → file
6. None of above → error not found

---

## IMPLEMENTATION PLAN

### Phase 1: Foundation

Extend existing command infrastructure and add missing registry methods for comprehensive collection management.

**Tasks:**
- Add missing registry methods for collection destruction and clearing
- Implement folder operation registry methods
- Add configuration management registry methods
- Create maintenance operation registry methods

### Phase 2: Core Collection Commands

Implement the primary collection management commands following existing patterns.

**Tasks:**
- Enhance create command for multiple collections
- Implement destroy command with file preservation
- Implement update command with add/remove subcommands
- Implement clear command with guard state management

### Phase 3: Configuration Management

Add configuration management commands with validation and warning systems.

**Tasks:**
- Implement config show command
- Implement config set command with subcommands
- Add validation for configuration changes
- Implement warnings for active guard states

### Phase 4: Maintenance & Advanced Features

Add maintenance commands, folder operations, and auto-detection capabilities.

**Tasks:**
- Implement cleanup, reset, and uninstall commands
- Add folder operations with dynamic scanning
- Enhance auto-detection in existing commands
- Add comprehensive error handling and output formatting

---

## STEP-BY-STEP TASKS

IMPORTANT: Execute every task in order, top to bottom. Each task is atomic and independently testable.

### CREATE cmd/guard/commands/destroy.go

- **IMPLEMENT**: Destroy collections command that removes collections but preserves files in registry
- **PATTERN**: Mirror create.go structure with variadic arguments
- **IMPORTS**: `github.com/spf13/cobra`, `github.com/florianbuetow/guard-tool/internal/manager`
- **GOTCHA**: Must disable guard on collection files before destroying collection
- **VALIDATE**: `go build -o build/guard ./cmd/guard && ./build/guard destroy --help`

### UPDATE internal/registry/collection_entry.go

- **IMPLEMENT**: Add `UnregisterCollection` method for destroying collections
- **PATTERN**: Mirror `UnregisterFile` pattern from file_entry.go (lines 180-200)
- **IMPORTS**: No new imports needed
- **GOTCHA**: Must preserve files in registry, only remove collection entry
- **VALIDATE**: `go test ./internal/registry -v -run TestUnregisterCollection`

### UPDATE internal/manager/manager.go

- **IMPLEMENT**: Add `DestroyCollection` method in manager layer
- **PATTERN**: Mirror `CreateCollection` method (lines 427-463)
- **IMPORTS**: No new imports needed
- **GOTCHA**: Must call `setCollectionProtection(name, false)` before destroying
- **VALIDATE**: `go test ./internal/manager -v -run TestDestroyCollection`

### UPDATE cmd/guard/main.go

- **IMPLEMENT**: Add destroy command to root command
- **PATTERN**: Mirror existing `rootCmd.AddCommand(commands.NewCreateCmd())` pattern
- **IMPORTS**: No new imports needed
- **GOTCHA**: Add after existing commands in init() function
- **VALIDATE**: `go build -o build/guard ./cmd/guard && ./build/guard --help | grep destroy`

### CREATE cmd/guard/commands/clear.go

- **IMPLEMENT**: Clear collections command that disables guard and removes files from collection
- **PATTERN**: Mirror destroy.go structure with collection validation
- **IMPORTS**: `github.com/spf13/cobra`, `github.com/florianbuetow/guard-tool/internal/manager`
- **GOTCHA**: Files remain in registry but are removed from collection
- **VALIDATE**: `go build -o build/guard ./cmd/guard && ./build/guard clear --help`

### UPDATE internal/registry/collection_entry.go

- **IMPLEMENT**: Add `ClearCollection` method for removing all files from collection
- **PATTERN**: Modify existing collection to have empty Files slice
- **IMPORTS**: No new imports needed
- **GOTCHA**: Keep collection entry but clear Files slice
- **VALIDATE**: `go test ./internal/registry -v -run TestClearCollection`

### UPDATE internal/manager/manager.go

- **IMPLEMENT**: Add `ClearCollection` method in manager layer
- **PATTERN**: Combine `setCollectionProtection(name, false)` with registry clear
- **IMPORTS**: No new imports needed
- **GOTCHA**: Must disable protection before clearing files
- **VALIDATE**: `go test ./internal/manager -v -run TestClearCollection`

### UPDATE cmd/guard/main.go

- **IMPLEMENT**: Add clear command to root command
- **PATTERN**: Add `rootCmd.AddCommand(commands.NewClearCmd())` in init()
- **IMPORTS**: No new imports needed
- **GOTCHA**: Maintain alphabetical order of commands
- **VALIDATE**: `go build -o build/guard ./cmd/guard && ./build/guard --help | grep clear`

### UPDATE cmd/guard/commands/update.go

- **IMPLEMENT**: Enhance update command with add/remove subcommands
- **PATTERN**: Use cobra subcommands with `cmd.AddCommand()` pattern
- **IMPORTS**: No new imports needed, enhance existing file
- **GOTCHA**: Must validate collection exists before add/remove operations
- **VALIDATE**: `go build -o build/guard ./cmd/guard && ./build/guard update --help`

### UPDATE internal/manager/manager.go

- **IMPLEMENT**: Add `RemoveFilesFromCollection` method
- **PATTERN**: Mirror `UpdateCollection` but remove files instead of add
- **IMPORTS**: No new imports needed
- **GOTCHA**: Files remain in registry, only removed from collection
- **VALIDATE**: `go test ./internal/manager -v -run TestRemoveFilesFromCollection`

### UPDATE internal/registry/collection_entry.go

- **IMPLEMENT**: Add `RemoveRegisteredFilesFromRegisteredCollection` method
- **PATTERN**: Reverse of `AddRegisteredFilesToRegisteredCollections`
- **IMPORTS**: No new imports needed
- **GOTCHA**: Use slice filtering to remove files from Files slice
- **VALIDATE**: `go test ./internal/registry -v -run TestRemoveFilesFromCollection`

### CREATE cmd/guard/commands/config.go

- **IMPLEMENT**: Configuration management command with show and set subcommands
- **PATTERN**: Use cobra subcommands like `guard config show` and `guard config set`
- **IMPORTS**: `github.com/spf13/cobra`, `github.com/florianbuetow/guard-tool/internal/manager`
- **GOTCHA**: Must warn if files are currently guarded when changing config
- **VALIDATE**: `go build -o build/guard ./cmd/guard && ./build/guard config --help`

### UPDATE internal/manager/manager.go

- **IMPLEMENT**: Add `ShowConfig` and `SetConfig` methods
- **PATTERN**: Access registry config directly with validation
- **IMPORTS**: No new imports needed
- **GOTCHA**: Check for active guard states before config changes
- **VALIDATE**: `go test ./internal/manager -v -run TestConfigManagement`

### UPDATE cmd/guard/main.go

- **IMPLEMENT**: Add config command to root command
- **PATTERN**: Add `rootCmd.AddCommand(commands.NewConfigCmd())` in init()
- **IMPORTS**: No new imports needed
- **GOTCHA**: Config is a parent command with subcommands
- **VALIDATE**: `go build -o build/guard ./cmd/guard && ./build/guard config --help`

### CREATE cmd/guard/commands/cleanup.go

- **IMPLEMENT**: Cleanup maintenance command for removing empty collections and missing files
- **PATTERN**: Mirror existing command structure with no arguments
- **IMPORTS**: `github.com/spf13/cobra`, `github.com/florianbuetow/guard-tool/internal/manager`
- **GOTCHA**: Must check file existence with `os.Stat()` before removal
- **VALIDATE**: `go build -o build/guard ./cmd/guard && ./build/guard cleanup --help`

### CREATE cmd/guard/commands/reset.go

- **IMPLEMENT**: Reset maintenance command that disables guard on all files and collections
- **PATTERN**: Mirror cleanup.go structure with no confirmation prompt
- **IMPORTS**: `github.com/spf13/cobra`, `github.com/florianbuetow/guard-tool/internal/manager`
- **GOTCHA**: Must iterate through all registered items and disable protection
- **VALIDATE**: `go build -o build/guard ./cmd/guard && ./build/guard reset --help`

### CREATE cmd/guard/commands/uninstall.go

- **IMPLEMENT**: Uninstall maintenance command that runs reset, cleanup, then deletes guardfile
- **PATTERN**: Chain reset and cleanup operations, then delete registry file
- **IMPORTS**: `github.com/spf13/cobra`, `github.com/florianbuetow/guard-tool/internal/manager`, `os`
- **GOTCHA**: No confirmation prompt - command should execute immediately for scriptability
- **VALIDATE**: `go build -o build/guard ./cmd/guard && ./build/guard uninstall --help`

### UPDATE internal/manager/manager.go

- **IMPLEMENT**: Add `Cleanup`, `Reset`, and `Uninstall` methods or move to domain-specific files
- **PATTERN**: Iterate through registry entries with validation and cleanup
- **IMPORTS**: Add `os` import for file operations
- **GOTCHA**: Reset must disable all protections before cleanup, no confirmation prompts
- **VALIDATE**: `go test ./internal/manager -v -run TestMaintenance`

### UPDATE cmd/guard/main.go

- **IMPLEMENT**: Add maintenance commands to root command
- **PATTERN**: Add cleanup, reset, uninstall commands in init()
- **IMPORTS**: No new imports needed
- **GOTCHA**: Maintain command order and help text consistency
- **VALIDATE**: `go build -o build/guard ./cmd/guard && ./build/guard --help | grep -E "cleanup|reset|uninstall"`

### REMOVE cmd/guard/commands/constants.go

- **IMPLEMENT**: Delete constants.go file and handle keywords inline in commands
- **PATTERN**: Define keywords as const within each command file that needs them
- **IMPORTS**: No imports needed for removal
- **GOTCHA**: Update any imports that reference this file
- **VALIDATE**: `find . -name "*.go" -exec grep -l "commands.Keyword" {} \; | wc -l` should return 0

### UPDATE cmd/guard/commands/enable.go

- **IMPLEMENT**: Enhance enable command with folder support and auto-detection
- **PATTERN**: Add folder keyword detection and auto-detection logic inline
- **IMPORTS**: No new imports needed, enhance existing file
- **GOTCHA**: Define keywords as const within the file, remove constants.go import
- **VALIDATE**: `go build -o build/guard ./cmd/guard && echo "test" > testfile && ./build/guard init 000 flo staff && ./build/guard enable testfile`

### UPDATE cmd/guard/commands/disable.go

- **IMPLEMENT**: Enhance disable command with folder support and auto-detection
- **PATTERN**: Mirror enable.go enhancements
- **IMPORTS**: No new imports needed, enhance existing file
- **GOTCHA**: Same auto-detection priority as enable command
- **VALIDATE**: `go build -o build/guard ./cmd/guard && ./build/guard disable testfile`

### UPDATE cmd/guard/commands/toggle.go

- **IMPLEMENT**: Enhance toggle command with folder support and auto-detection
- **PATTERN**: Mirror enable.go enhancements with toggle logic
- **IMPORTS**: No new imports needed, enhance existing file
- **GOTCHA**: Must handle auto-registration for files not in registry
- **VALIDATE**: `go build -o build/guard ./cmd/guard && ./build/guard toggle testfile`

### UPDATE cmd/guard/commands/show.go

- **IMPLEMENT**: Enhance show command with auto-detection and folder support
- **PATTERN**: Add auto-detection logic to existing show command
- **IMPORTS**: No new imports needed, enhance existing file
- **GOTCHA**: Must display appropriate information based on detected type
- **VALIDATE**: `go build -o build/guard ./cmd/guard && ./build/guard show testfile`

### CREATE internal/manager/folders.go

- **IMPLEMENT**: Folder operation methods for dynamic file scanning
- **PATTERN**: Create methods for folder registration, scanning, and protection
- **IMPORTS**: `os`, `path/filepath`, `fmt`
- **GOTCHA**: Folders scan immediate files non-recursively
- **VALIDATE**: `go test ./internal/manager -v -run TestFolderOperations`

### UPDATE internal/registry/folder_entry.go

- **IMPLEMENT**: Add missing folder registry methods if needed
- **PATTERN**: Ensure complete CRUD operations for folders
- **IMPORTS**: No new imports needed
- **GOTCHA**: Folder names stored with @ prefix in registry
- **VALIDATE**: `go test ./internal/registry -v -run TestFolderEntry`

### UPDATE internal/manager/manager.go

- **IMPLEMENT**: Integrate folder operations into detectTargetType method
- **PATTERN**: Add folder detection logic following priority order
- **IMPORTS**: No new imports needed
- **GOTCHA**: Check directory existence before folder operations
- **VALIDATE**: `go test ./internal/manager -v -run TestDetectTargetType`

### REFACTOR cmd/guard/commands/create.go

- **IMPLEMENT**: Enhance create command to support multiple collections
- **PATTERN**: Use variadic arguments with validation loop
- **IMPORTS**: No new imports needed
- **GOTCHA**: Must validate each collection name for reserved keywords
- **VALIDATE**: `go build -o build/guard ./cmd/guard && ./build/guard create coll1 coll2 coll3`

### UPDATE internal/manager/manager.go

- **IMPLEMENT**: Enhance CreateCollection to handle multiple collections
- **PATTERN**: Add loop for multiple collection creation
- **IMPORTS**: No new imports needed
- **GOTCHA**: Continue on individual failures, report summary
- **VALIDATE**: `go test ./internal/manager -v -run TestCreateMultipleCollections`

---

## TESTING STRATEGY

### Unit Tests

Design unit tests following existing patterns in `internal/` packages:
- Test each new registry method with positive and negative cases
- Test manager methods with mock registry states
- Test error conditions and edge cases
- Use table-driven tests for multiple input scenarios

### Integration Tests

Create shell integration tests following `tests/` directory patterns:
- Test complete command workflows end-to-end
- Validate exact output formats match specifications
- Test error conditions and exit codes
- Test file system interactions and state persistence

### Edge Cases

- Empty collections and missing files
- Reserved keyword validation
- Corrupted guardfile handling
- Permission failures and graceful degradation
- Multiple collections containing same files
- Folder operations with no files or missing directories

---

## VALIDATION COMMANDS

Execute every command to ensure zero regressions and 100% feature correctness.

### Level 1: Syntax & Style

```bash
go fmt ./...
golangci-lint run
semgrep scan --config auto
```

### Level 2: Unit Tests

```bash
go test ./internal/registry -v
go test ./internal/manager -v
go test ./cmd/guard/commands -v
```

### Level 3: Integration Tests

```bash
just build
cp build/guard tests/guard
cd tests && ./run-all-tests.sh
```

### Level 4: Manual Validation

```bash
# Test collection management workflow
./build/guard init 000 flo staff
./build/guard create mygroup
./build/guard update mygroup add file1.txt file2.txt
./build/guard enable mygroup
./build/guard show mygroup
./build/guard disable mygroup
./build/guard clear mygroup
./build/guard destroy mygroup

# Test configuration management
./build/guard config show
./build/guard config set mode 640
./build/guard config set owner newowner
./build/guard config set group newgroup

# Test maintenance commands
./build/guard cleanup
./build/guard reset
```

### Level 5: Additional Validation (Optional)

```bash
# Test build across platforms
GOOS=linux go build -o build/guard-linux ./cmd/guard
GOOS=darwin go build -o build/guard-darwin ./cmd/guard
```

---

## ACCEPTANCE CRITERIA

- [ ] All collection management commands implemented (create, destroy, update, clear)
- [ ] Configuration management commands working (config show, config set)
- [ ] Maintenance commands functional (cleanup, reset, uninstall)
- [ ] Folder operations with dynamic scanning implemented
- [ ] Auto-detection working across all relevant commands
- [ ] All validation commands pass with zero errors
- [ ] Unit test coverage meets requirements (60%+)
- [ ] Integration tests verify end-to-end workflows
- [ ] Code follows existing project conventions and patterns
- [ ] No regressions in existing functionality
- [ ] Error messages follow specified format patterns
- [ ] Output messages use correct past tense verbs and counting format
- [ ] Reserved keyword validation prevents invalid collection names
- [ ] Path security measures implemented (relative paths, symlink rejection)

---

## COMPLETION CHECKLIST

- [ ] All 25+ tasks completed in dependency order
- [ ] Each task validation passed immediately after implementation
- [ ] All validation commands executed successfully
- [ ] Full test suite passes (unit + integration)
- [ ] No linting or type checking errors
- [ ] Manual testing confirms all features work as specified
- [ ] Acceptance criteria all met
- [ ] Code reviewed for quality and maintainability
- [ ] Documentation updated in help text and command descriptions

---

## NOTES

**Design Decisions:**
- Collections and folders are separate concepts - collections are static file lists, folders are dynamic directory scanners
- Auto-detection follows strict priority order to ensure predictable behavior
- Configuration changes warn about active guard states but don't automatically update them
- Maintenance commands execute immediately without confirmation for scriptability
- Error handling follows existing patterns with consistent message formatting
- Keywords handled inline in command files rather than separate constants file
- Manager package organized by domain (folders.go, collections.go, config.go, files.go)

**Implementation Risks:**
- Complex command structure with subcommands requires careful Cobra setup
- Auto-detection logic must be consistent across all commands
- Folder operations need careful file system scanning without recursion
- Registry operations must maintain thread safety and atomic saves

**Performance Considerations:**
- Folder scanning should be efficient for directories with many files
- Registry operations should batch saves to avoid excessive I/O
- Auto-detection should fail fast on missing targets

**Security Considerations:**
- Path validation must prevent traversal attacks
- Symlink rejection prevents security bypasses
- Configuration validation prevents invalid states
- Reserved keyword checking prevents command conflicts

# Feature: TUI (Text User Interface) Implementation

The following plan should be complete, but it's important that you validate documentation and codebase patterns and task sanity before you start implementing.

Pay special attention to naming of existing utils types and models. Import from the right files etc.

## Feature Description

Implement an interactive Text User Interface (TUI) for the guard tool using the Bubble Tea framework. The TUI will provide a user-friendly interface for managing protected files, collections, and folders without requiring command-line arguments. Users can browse registered items, toggle protection status, view file details, and perform quick enable/disable operations through keyboard navigation.

The TUI must mirror the exact output format of existing `guard show` commands and provide sub-second response times as specified in the technical requirements.

## User Story

As a developer using guard-tool
I want to interact with protected files through a visual terminal interface
So that I can quickly browse, toggle, and manage file protection without memorizing CLI commands

## Problem Statement

The current guard-tool requires users to remember specific CLI commands and file paths to manage protection. Users need a more intuitive way to:
- Browse all registered files, collections, and folders
- See protection status at a glance
- Toggle protection with simple keystrokes
- Navigate through items without typing paths

## Solution Statement

Implement a Bubble Tea-based TUI that provides:
- Interactive list navigation with keyboard controls
- Real-time status display (protected/unprotected)
- Quick toggle operations (spacebar/enter)
- Tabbed interface for files/collections/folders
- Status messages and error handling
- Integration with existing manager layer
- Mirror CLI show command output formats

## Feature Metadata

**Feature Type**: New Capability
**Estimated Complexity**: High
**Primary Systems Affected**: CLI entry point, new TUI package, Manager integration
**Dependencies**: Bubble Tea framework, Bubbles components

---

## TUI DESIGN SPECIFICATIONS

### Main Interface Layout

```
╔══════════════════════════════════════════════════════════════════════════════╗
║ Guard Tool - File Protection Manager                                        ║
╠══════════════════════════════════════════════════════════════════════════════╣
║ [1] Files    [2] Collections    [3] Folders                                 ║
╠══════════════════════════════════════════════════════════════════════════════╣
║                                                                              ║
║  > G config.yaml (web, docs)                                                ║
║    - README.md ()                                                           ║
║    G src/main.go (core)                                                     ║
║    - tests/unit.go ()                                                       ║
║                                                                              ║
║                                                                              ║
║                                                                              ║
║                                                                              ║
║                                                                              ║
║                                                                              ║
║                                                                              ║
║                                                                              ║
║                                                                              ║
║                                                                              ║
║                                                                              ║
╠══════════════════════════════════════════════════════════════════════════════╣
║ 4 files total: 2 guarded, 2 unguarded │ t=toggle │ ?=help │ q=quit         ║
╚══════════════════════════════════════════════════════════════════════════════╝
```

### Files Tab Layout

```
╔══════════════════════════════════════════════════════════════════════════════╗
║ Guard Tool - File Protection Manager                                        ║
╠══════════════════════════════════════════════════════════════════════════════╣
║ [1] Files*   [2] Collections    [3] Folders                                 ║
╠══════════════════════════════════════════════════════════════════════════════╣
║                                                                              ║
║  > G config.yaml (web, docs)                                                ║
║    - README.md ()                                                           ║
║    G src/main.go (core)                                                     ║
║    - tests/unit.go ()                                                       ║
║    G .gitignore ()                                                          ║
║    - package.json (web)                                                     ║
║                                                                              ║
║                                                                              ║
║                                                                              ║
║                                                                              ║
║                                                                              ║
║                                                                              ║
║                                                                              ║
║                                                                              ║
║                                                                              ║
║                                                                              ║
╠══════════════════════════════════════════════════════════════════════════════╣
║ 6 files total: 3 guarded, 3 unguarded │ t=toggle │ ?=help │ q=quit         ║
╚══════════════════════════════════════════════════════════════════════════════╝
```

### Collections Tab Layout

```
╔══════════════════════════════════════════════════════════════════════════════╗
║ Guard Tool - File Protection Manager                                        ║
╠══════════════════════════════════════════════════════════════════════════════╣
║ [1] Files    [2] Collections*   [3] Folders                                 ║
╠══════════════════════════════════════════════════════════════════════════════╣
║                                                                              ║
║  > G collection: web (3 files)                                              ║
║    - collection: docs (2 files)                                             ║
║    G collection: core (1 files)                                             ║
║                                                                              ║
║                                                                              ║
║                                                                              ║
║                                                                              ║
║                                                                              ║
║                                                                              ║
║                                                                              ║
║                                                                              ║
║                                                                              ║
║                                                                              ║
║                                                                              ║
║                                                                              ║
║                                                                              ║
╠══════════════════════════════════════════════════════════════════════════════╣
║ 3 collections total: 2 guarded, 1 unguarded │ t=toggle │ ?=help │ q=quit   ║
╚══════════════════════════════════════════════════════════════════════════════╝
```

### Folders Tab Layout

```
╔══════════════════════════════════════════════════════════════════════════════╗
║ Guard Tool - File Protection Manager                                        ║
╠══════════════════════════════════════════════════════════════════════════════╣
║ [1] Files    [2] Collections    [3] Folders*                                ║
╠══════════════════════════════════════════════════════════════════════════════╣
║                                                                              ║
║  > G @src (./src) - All Protected                                           ║
║    - @docs (./docs) - Mixed                                                 ║
║    G @tests (./tests) - All Protected                                       ║
║                                                                              ║
║                                                                              ║
║                                                                              ║
║                                                                              ║
║                                                                              ║
║                                                                              ║
║                                                                              ║
║                                                                              ║
║                                                                              ║
║                                                                              ║
║                                                                              ║
║                                                                              ║
║                                                                              ║
╠══════════════════════════════════════════════════════════════════════════════╣
║ 3 folders total: 2 guarded, 1 unguarded │ t=toggle │ ?=help │ q=quit       ║
╚══════════════════════════════════════════════════════════════════════════════╝
```

### Help Screen Layout

```
╔══════════════════════════════════════════════════════════════════════════════╗
║ Guard Tool - Help                                                           ║
╠══════════════════════════════════════════════════════════════════════════════╣
║                                                                              ║
║  NAVIGATION:                                                                 ║
║    ↑/k, ↓/j     Navigate up/down                                            ║
║    1, 2, 3      Switch tabs (Files, Collections, Folders)                  ║
║                                                                              ║
║  ACTIONS:                                                                    ║
║    t, Space     Toggle protection for selected item                         ║
║    Enter        Toggle protection for selected item                         ║
║                                                                              ║
║  OTHER:                                                                      ║
║    ?            Show/hide this help                                         ║
║    q, Ctrl+C    Quit application                                            ║
║                                                                              ║
║  STATUS INDICATORS:                                                          ║
║    G            Protected (guarded)                                         ║
║    -            Unprotected                                                 ║
║    >            Currently selected item                                     ║
║                                                                              ║
║                                                                              ║
║                                                                              ║
╠══════════════════════════════════════════════════════════════════════════════╣
║ Press ? or Esc to close help                                                ║
╚══════════════════════════════════════════════════════════════════════════════╝
```

### Error/Warning Display

```
╔══════════════════════════════════════════════════════════════════════════════╗
║ Guard Tool - File Protection Manager                                        ║
╠══════════════════════════════════════════════════════════════════════════════╣
║ [1] Files*   [2] Collections    [3] Folders                                 ║
╠══════════════════════════════════════════════════════════════════════════════╣
║                                                                              ║
║  > G config.yaml (web, docs)                                                ║
║    - README.md ()                                                           ║
║    G src/main.go (core)                                                     ║
║                                                                              ║
║                                                                              ║
║                                                                              ║
║                                                                              ║
║                                                                              ║
║                                                                              ║
║                                                                              ║
║                                                                              ║
║                                                                              ║
║                                                                              ║
║                                                                              ║
║                                                                              ║
║                                                                              ║
╠══════════════════════════════════════════════════════════════════════════════╣
║ WARNING: File missing: old.txt │ 3 files total: 2 guarded, 1 unguarded    ║
╚══════════════════════════════════════════════════════════════════════════════╝
```

## KEYBOARD SHORTCUTS SPECIFICATION

### Primary Navigation
- **↑ / k**: Move selection up
- **↓ / j**: Move selection down
- **1**: Switch to Files tab
- **2**: Switch to Collections tab  
- **3**: Switch to Folders tab

### Actions
- **t**: Toggle protection for selected item
- **Space**: Toggle protection for selected item
- **Enter**: Toggle protection for selected item

### System
- **?**: Show/hide help screen
- **Esc**: Close help screen (when open)
- **q**: Quit application
- **Ctrl+C**: Force quit application

### Future Extensions
- **f**: Filter/search (future feature)
- **r**: Refresh display (future feature)

## OUTPUT FORMAT SPECIFICATIONS

### Files Display Format
Mirror `guard show file` output exactly:
```
G filename (collection1, collection2)
- filename ()
```

### Collections Display Format  
Mirror `guard show collection` output exactly:
```
G collection: name (n files)
- collection: name (n files)
```

### Folders Display Format
New format based on folder state:
```
G @foldername (./path) - All Protected
- @foldername (./path) - Mixed
- @foldername (./path) - None Protected
```

### Status Bar Format
```
{count} {type} total: {guarded} guarded, {unguarded} unguarded │ t=toggle │ ?=help │ q=quit
```

## TUI TESTING SPECIFICATIONS

### Unit Tests (Go Testing Framework)
```go
// Test files: internal/tui/*_test.go
func TestTUIModelStateTransitions(t *testing.T)
func TestKeyBindingHandlers(t *testing.T) 
func TestTabNavigation(t *testing.T)
func TestItemListGeneration(t *testing.T)
func TestStatusBarFormatting(t *testing.T)
```

### Integration Tests (tmux + Shell)
```bash
# Test file: tests/test-tui-integration.sh
test_tui_launch_basic() {
    # Test TUI launches without errors
    tmux new-session -d -s test_tui "./guard -i"
    tmux capture-pane -t test_tui -p | grep "Guard Tool"
}

test_tui_navigation() {
    # Test keyboard navigation works
    tmux send-keys -t test_tui "j" "k" "2" "1"
    tmux capture-pane -t test_tui -p | grep "Files\*"
}

test_tui_toggle_operation() {
    # Test toggle functionality
    tmux send-keys -t test_tui "t"
    tmux capture-pane -t test_tui -p | grep "G "
}

test_tui_help_screen() {
    # Test help screen display
    tmux send-keys -t test_tui "?"
    tmux capture-pane -t test_tui -p | grep "NAVIGATION"
}

test_tui_quit() {
    # Test quit functionality
    tmux send-keys -t test_tui "q"
    ! tmux has-session -t test_tui
}
```

### Performance Tests
```bash
# Test response time requirements (sub-second)
test_tui_performance() {
    start_time=$(date +%s%N)
    ./guard -i &
    TUI_PID=$!
    # Wait for TUI to be ready
    sleep 0.1
    kill $TUI_PID
    end_time=$(date +%s%N)
    duration=$(( (end_time - start_time) / 1000000 ))
    [ $duration -lt 1000 ] # Less than 1 second
}
```

### Edge Case Tests
```bash
test_tui_empty_registry()     # No files/collections/folders
test_tui_registry_errors()    # Registry loading failures  
test_tui_permission_errors()  # Filesystem operation failures
test_tui_terminal_resize()    # Terminal size changes
test_tui_invalid_keys()       # Invalid keyboard input
```

---

## CONTEXT REFERENCES

### Relevant Codebase Files IMPORTANT: YOU MUST READ THESE FILES BEFORE IMPLEMENTING!

- `cmd/guard/main.go` (lines 18-30) - Why: Contains TUI mode detection and placeholder implementation
- `internal/manager/manager.go` (lines 14-20, 23-30) - Why: Manager struct and constructor pattern to follow
- `internal/manager/files.go` (lines 10-18) - Why: FileInfo struct for display data pattern
- `internal/security/security.go` (lines 84-110) - Why: GetRegistered* method patterns to follow
- `internal/registry/folder_entry.go` (lines 122-135) - Why: GetRegisteredFolders() method to expose
- `internal/manager/collections.go` (lines 8-15) - Why: Collection validation and reserved keywords
- `internal/manager/warnings.go` - Why: Warning system integration patterns
- `internal/registry/registry.go` (lines 25-35) - Why: Registry data structures and access patterns
- `go.mod` - Why: Current dependencies and Go version
- `justfile` - Why: Build and dependency management patterns

### New Files to Create

- `internal/manager/data_types.go` - CollectionInfo and FolderInfo structs for TUI display
- `internal/manager/data_access.go` - GetCollections() and GetFolders() methods for TUI
- `internal/tui/tui.go` - Main TUI application struct and initialization
- `internal/tui/models.go` - Bubble Tea models for different views
- `internal/tui/views.go` - View rendering functions with ASCII mockups
- `internal/tui/keys.go` - Key binding definitions and handlers
- `internal/tui/items.go` - List item implementations for files/collections/folders
- `internal/tui/tabs.go` - Tab navigation component
- `internal/tui/status.go` - Status bar and message display
- `tests/test-tui-integration.sh` - TUI integration tests using tmux

### Relevant Documentation YOU SHOULD READ THESE BEFORE IMPLEMENTING!

- [Bubble Tea Framework](https://github.com/charmbracelet/bubbletea)
  - Specific section: Model-View-Update architecture
  - Why: Core framework patterns for TUI implementation
- [Bubbles Components](https://github.com/charmbracelet/bubbles)
  - Specific section: List component with filtering
  - Why: Interactive list navigation for files/collections
- [Bubble Tea Tutorial](https://github.com/charmbracelet/bubbletea/tree/master/tutorials)
  - Specific section: Key handling and navigation
  - Why: Keyboard interaction patterns

### Patterns to Follow

**Manager Integration Pattern:**
```go
// From internal/manager/manager.go
type Manager struct {
    registryPath string
    security     *security.Security
    fs           *filesystem.FileSystem
    warnings     []Warning
    errors       []string
}

func NewManager(registryPath string) *Manager {
    return &Manager{
        registryPath: registryPath,
        fs:           filesystem.NewFileSystem(),
        warnings:     make([]Warning, 0),
        errors:       make([]string, 0),
    }
}
```

**Error Handling Pattern:**
```go
// From internal/manager/files.go
if err := m.security.RegisterFile(path, mode, owner, group); err != nil {
    m.AddError(fmt.Sprintf("Error: Failed to register %s: %v", path, err))
    continue
}
```

**Display Data Pattern:**
```go
// From internal/manager/files.go
type FileInfo struct {
    Path        string
    Guard       bool
    Collections []string
}

// New structs needed for TUI:
type CollectionInfo struct {
    Name      string
    Guard     bool
    FileCount int
}

type FolderInfo struct {
    Name  string
    Path  string
    Guard bool
}
```

**Security Layer Pattern:**
```go
// From internal/security/security.go - pattern to follow
func (s *Security) GetRegisteredFiles() []string {
    return s.registry.GetRegisteredFiles()
}

// Missing method to add:
func (s *Security) GetRegisteredFolders() []string {
    return s.registry.GetRegisteredFolders()
}
```

**CLI Integration Pattern:**
```go
// From cmd/guard/main.go
if interactiveMode {
    // TODO: Launch TUI mode
    fmt.Println("Interactive TUI mode not yet implemented")
    os.Exit(0)
}
```

---

## IMPLEMENTATION PLAN

### Phase 1: Prerequisites - Manager Data Access

Add missing methods to Security and Manager layers for TUI data retrieval.

**Tasks:**
- Add GetRegisteredFolders() method to Security layer
- Create CollectionInfo and FolderInfo structs in manager
- Add Manager.GetCollections() and Manager.GetFolders() methods
- Ensure TUI has proper data access patterns

### Phase 2: Foundation Setup

Set up Bubble Tea dependency and basic TUI structure with manager integration.

**Tasks:**
- Add Bubble Tea dependencies to go.mod
- Create TUI package structure
- Implement basic TUI application wrapper
- Integrate with existing manager layer

### Phase 3: Core TUI Components

Implement the main TUI models, views, and navigation system.

**Tasks:**
- Create tab navigation system (Files/Collections/Folders)
- Implement list components for each tab
- Add keyboard navigation and selection
- Create status bar and message display

### Phase 4: Manager Integration

Connect TUI to existing manager operations for data display and manipulation.

**Tasks:**
- Integrate file listing and status display
- Implement collection and folder browsing
- Add toggle operations through manager
- Handle warnings and errors from manager

### Phase 5: Interactive Operations

Add interactive features for protection management.

**Tasks:**
- Implement toggle protection (spacebar/enter)
- Add bulk operations support
- Create confirmation dialogs for destructive operations
- Add help screen and key bindings display

---

## STEP-BY-STEP TASKS

IMPORTANT: Execute every task in order, top to bottom. Each task is atomic and independently testable.

### ADD internal/security/security.go - GetRegisteredFolders method

- **IMPLEMENT**: Add GetRegisteredFolders() method to Security struct
- **PATTERN**: Mirror GetRegisteredFiles() and GetRegisteredCollections() patterns from internal/security/security.go:84-110
- **IMPORTS**: Use existing registry access pattern
- **GOTCHA**: Method must expose registry.GetRegisteredFolders() through security layer
- **VALIDATE**: `go build -o build/guard ./cmd/guard`

### CREATE internal/manager/data_types.go

- **IMPLEMENT**: CollectionInfo and FolderInfo structs for TUI display data
- **PATTERN**: Mirror FileInfo struct pattern from internal/manager/files.go:10-18
- **IMPORTS**: 
  ```go
  import (
      "github.com/florianbuetow/guard-tool/internal/manager"
  )
  ```
- **GOTCHA**: Include all fields needed for TUI display (name, guard status, counts, paths)
- **VALIDATE**: `go build -o build/guard ./cmd/guard`

### ADD internal/manager/data_access.go

- **IMPLEMENT**: GetCollections() and GetFolders() methods for TUI data retrieval
- **PATTERN**: Mirror ShowFiles() method pattern from internal/manager/files.go:383-447
- **IMPORTS**: 
  ```go
  import (
      "fmt"
      "github.com/florianbuetow/guard-tool/internal/manager"
  )
  ```
- **GOTCHA**: Return structured data ([]CollectionInfo, []FolderInfo) not formatted output
- **VALIDATE**: `go build -o build/guard ./cmd/guard`

### UPDATE go.mod

- **IMPLEMENT**: Add Bubble Tea and Bubbles dependencies
- **PATTERN**: Follow existing dependency format in go.mod
- **IMPORTS**: 
  ```
  github.com/charmbracelet/bubbletea v0.27.1
  github.com/charmbracelet/bubbles v0.20.0
  github.com/charmbracelet/lipgloss v0.13.1
  ```
- **GOTCHA**: Use compatible versions that work with Go 1.25.6
- **VALIDATE**: `go mod tidy && go mod verify`

### CREATE internal/tui/tui.go

- **IMPLEMENT**: Main TUI application struct and initialization
- **PATTERN**: Mirror Manager constructor pattern from internal/manager/manager.go:23-30
- **IMPORTS**: 
  ```go
  import (
      "fmt"
      tea "github.com/charmbracelet/bubbletea"
      "github.com/florianbuetow/guard-tool/internal/manager"
  )
  ```
- **GOTCHA**: Handle manager loading errors gracefully
- **VALIDATE**: `go build -o build/guard ./cmd/guard`

### CREATE internal/tui/models.go

- **IMPLEMENT**: Bubble Tea model structs for different views
- **PATTERN**: Follow Bubble Tea Model-View-Update architecture from documentation
- **IMPORTS**: 
  ```go
  import (
      "github.com/charmbracelet/bubbles/list"
      "github.com/charmbracelet/bubbles/key"
      tea "github.com/charmbracelet/bubbletea"
  )
  ```
- **GOTCHA**: Embed list.Model for each tab (files, collections, folders)
- **VALIDATE**: `go build -o build/guard ./cmd/guard`

### CREATE internal/tui/items.go

- **IMPLEMENT**: List item implementations for files, collections, and folders
- **PATTERN**: Follow bubbles list.Item interface (Title(), Description(), FilterValue())
- **IMPORTS**: 
  ```go
  import (
      "fmt"
      "github.com/florianbuetow/guard-tool/internal/manager"
  )
  ```
- **GOTCHA**: Use FileInfo, CollectionInfo, FolderInfo structs from internal/manager/
- **VALIDATE**: `go build -o build/guard ./cmd/guard`

### CREATE internal/tui/keys.go

- **IMPLEMENT**: Key binding definitions and handlers
- **PATTERN**: Follow Bubble Tea key handling from documentation examples
- **IMPORTS**: 
  ```go
  import (
      "github.com/charmbracelet/bubbles/key"
  )
  ```
- **GOTCHA**: Implement exact keyboard shortcuts: ↑/k ↓/j (navigation), 1/2/3 (tabs), t/space/enter (toggle), ?/esc (help), q/ctrl+c (quit)
- **VALIDATE**: `go build -o build/guard ./cmd/guard`

### CREATE internal/tui/tabs.go

- **IMPLEMENT**: Tab navigation component for Files/Collections/Folders
- **PATTERN**: Simple state machine with active tab index
- **IMPORTS**: 
  ```go
  import (
      "github.com/charmbracelet/lipgloss"
  )
  ```
- **GOTCHA**: Use lipgloss for styling active/inactive tabs
- **VALIDATE**: `go build -o build/guard ./cmd/guard`

### CREATE internal/tui/views.go

- **IMPLEMENT**: View rendering functions for each component using DOUBLE-LINE box drawing (╔╗╚╝║═)
- **PATTERN**: Follow Bubble Tea View() method pattern from documentation
- **IMPORTS**: 
  ```go
  import (
      "fmt"
      "strings"
      "github.com/charmbracelet/lipgloss"
  )
  ```
- **GOTCHA**: Must render exact ASCII mockups with proper box drawing, mirror CLI show command output formats
- **VALIDATE**: `go build -o build/guard ./cmd/guard`

### CREATE internal/tui/status.go

- **IMPLEMENT**: Status bar component with exact format: "{count} {type} total: {guarded} guarded, {unguarded} unguarded │ t=toggle │ ?=help │ q=quit"
- **PATTERN**: Follow manager warning/error pattern from internal/manager/warnings.go
- **IMPORTS**: 
  ```go
  import (
      "github.com/charmbracelet/lipgloss"
      "github.com/florianbuetow/guard-tool/internal/manager"
  )
  ```
- **GOTCHA**: Display manager warnings/errors in status bar, use │ separator character
- **VALIDATE**: `go build -o build/guard ./cmd/guard`

### UPDATE cmd/guard/main.go

- **IMPLEMENT**: Replace TUI placeholder with actual TUI launch
- **PATTERN**: Follow existing manager initialization pattern
- **IMPORTS**: 
  ```go
  import (
      "github.com/florianbuetow/guard-tool/internal/tui"
  )
  ```
- **GOTCHA**: Handle registry loading errors and display in TUI
- **VALIDATE**: `go build -o build/guard ./cmd/guard && ./build/guard -i`

### ADD TUI Integration Methods

- **IMPLEMENT**: Methods in TUI for manager operations (toggle, enable, disable)
- **PATTERN**: Use Manager.GetCollections() and Manager.GetFolders() for data retrieval
- **IMPORTS**: Use new manager data access methods
- **GOTCHA**: Handle async operations and update UI state, use structured data not CLI output
- **VALIDATE**: `go build -o build/guard ./cmd/guard && ./build/guard -i`

### ADD Help Screen

- **IMPLEMENT**: Help overlay with exact ASCII mockup layout using DOUBLE-LINE box drawing
- **PATTERN**: Modal overlay pattern from bubbletea examples
- **IMPORTS**: 
  ```go
  import (
      "github.com/charmbracelet/bubbles/viewport"
  )
  ```
- **GOTCHA**: Toggle help with '?' key, dismiss with 'esc' or '?', show all keyboard shortcuts
- **VALIDATE**: `go build -o build/guard ./cmd/guard && ./build/guard -i`

### ADD Error Handling

- **IMPLEMENT**: Error display and recovery in TUI
- **PATTERN**: Follow manager error handling from internal/manager/manager.go
- **IMPORTS**: Use existing manager error types
- **GOTCHA**: Show errors in status bar, allow dismissal with any key
- **VALIDATE**: `go build -o build/guard ./cmd/guard && ./build/guard -i`

### ADD Confirmation Dialogs

- **IMPLEMENT**: Confirmation prompts for destructive operations
- **PATTERN**: Modal dialog pattern with yes/no options
- **IMPORTS**: Use lipgloss for styling
- **GOTCHA**: Only show for operations that modify filesystem (toggle, enable, disable)
- **VALIDATE**: `go build -o build/guard ./cmd/guard && ./build/guard -i`

### CREATE tests/test-tui-integration.sh

- **IMPLEMENT**: TUI integration tests using tmux for terminal simulation
- **PATTERN**: Follow existing shell test patterns from tests/ directory
- **IMPORTS**: Use tmux commands for TUI interaction testing
- **GOTCHA**: Must test all keyboard shortcuts, tab navigation, toggle operations, help screen
- **VALIDATE**: `cd tests && ./test-tui-integration.sh`

---

## TESTING STRATEGY

### Unit Tests

Design unit tests following Go testing conventions:
- Test TUI model state transitions
- Test key binding handlers
- Test item list generation from manager data
- Test tab navigation logic

### Integration Tests

Add TUI integration tests to existing shell test suite:
- Test TUI launch with `guard -i`
- Test basic navigation (requires tmux for terminal simulation)
- Test file listing and status display
- Test toggle operations through TUI

### Edge Cases

Test specific edge cases for TUI:
- Empty registry (no files/collections/folders)
- Registry loading errors
- Manager operation failures
- Terminal resize handling
- Keyboard input edge cases

---

## VALIDATION COMMANDS

Execute every command to ensure zero regressions and 100% feature correctness.

### Level 1: Syntax & Style

```bash
go fmt ./...
golangci-lint run
```

### Level 2: Dependencies

```bash
go mod tidy
go mod verify
go mod download
```

### Level 3: Build Validation

```bash
go build -o build/guard ./cmd/guard
```

### Level 4: Unit Tests

```bash
go test ./internal/tui/...
go test ./...
```

### Level 5: Integration Tests

```bash
cp build/guard tests/guard
cd tests && ./run-all-tests.sh
```

### Level 6: TUI Integration Tests (tmux)

```bash
# Create TUI test file
cd tests && ./test-tui-integration.sh
```

### Level 7: Manual TUI Testing

```bash
# Test TUI launch
./build/guard -i

# Test with existing registry
echo "config:\n  guard_mode: \"0640\"\n  guard_owner: testuser\n  guard_group: testgroup\nfiles: []\ncollections: []\nfolders: []" > .guardfile
./build/guard -i
```

---

## ACCEPTANCE CRITERIA

- [ ] Security.GetRegisteredFolders() method added and functional
- [ ] CollectionInfo and FolderInfo structs created with proper fields
- [ ] Manager.GetCollections() and Manager.GetFolders() methods implemented
- [ ] TUI launches successfully with `guard -i` flag
- [ ] TUI renders exact ASCII mockups with DOUBLE-LINE box drawing (╔╗╚╝║═)
- [ ] Three tabs (Files, Collections, Folders) navigable with 1/2/3 keys
- [ ] List navigation works with ↑/k and ↓/j keys
- [ ] Files display mirrors `guard show file` format: "G filename (collections)"
- [ ] Collections display mirrors `guard show collection` format: "G collection: name (n files)"
- [ ] Folders display shows: "G @name (./path) - State"
- [ ] Toggle works with t/space/enter keys
- [ ] Status bar shows exact format: "{count} {type} total: {guarded} guarded, {unguarded} unguarded │ t=toggle │ ?=help │ q=quit"
- [ ] Help screen accessible with '?' key, dismissible with 'esc'
- [ ] Quit functionality works with 'q' and Ctrl+C
- [ ] Manager integration preserves all existing functionality
- [ ] Error/warning messages display in status bar
- [ ] Terminal resize handled gracefully
- [ ] Sub-second response time requirement met
- [ ] All tmux-based TUI tests pass
- [ ] All validation commands pass with zero errors

---

## COMPLETION CHECKLIST

- [ ] Security.GetRegisteredFolders() method added
- [ ] CollectionInfo and FolderInfo structs created
- [ ] Manager.GetCollections() and Manager.GetFolders() methods implemented
- [ ] All dependencies added to go.mod
- [ ] TUI package structure created
- [ ] Bubble Tea models implemented
- [ ] List components functional
- [ ] Tab navigation working
- [ ] Key bindings implemented
- [ ] Manager integration complete
- [ ] Status display functional
- [ ] Help screen implemented
- [ ] Error handling added
- [ ] CLI integration updated
- [ ] All validation commands pass
- [ ] Manual testing confirms functionality
- [ ] No regressions in existing CLI commands

---

## NOTES

**Design Decisions:**
- Use three separate list models for files/collections/folders to maintain independent state
- Integrate directly with existing manager layer to avoid code duplication
- Follow Bubble Tea best practices for responsive terminal applications
- Maintain consistency with existing CLI command behavior

**Performance Considerations:**
- Lazy load file lists to handle large registries
- Use Bubble Tea's built-in filtering for search functionality
- Minimize manager calls during navigation

**Future Extensibility:**
- TUI structure supports adding new tabs (e.g., configuration, logs)
- Key binding system allows easy addition of new operations
- Status system can be extended for progress indicators

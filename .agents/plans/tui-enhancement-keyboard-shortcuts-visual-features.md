# Feature: TUI Enhancement: Missing Keyboard Shortcuts and Visual Features

The following plan should be complete, but its important that you validate documentation and codebase patterns and task sanity before you start implementing.

Pay special attention to naming of existing utils types and models. Import from the right files etc.

## Feature Description

Enhance the Guard-Tool TUI with missing keyboard shortcuts and visual improvements to provide a more complete and user-friendly terminal interface. This includes recursive folder operations, manual refresh capability, proper symlink handling, and improved text display for long paths.

## User Story

As a developer using Guard-Tool's TUI
I want enhanced keyboard shortcuts and visual feedback
So that I can efficiently manage file protection with recursive operations, manual refresh, proper symlink handling, and readable text display

## Problem Statement

The current TUI implementation lacks several essential features that limit its usability:
1. No recursive toggle operation for folders (only immediate children are affected)
2. No manual refresh capability to reload external changes
3. Symlinks are not visually distinguished and can be incorrectly operated on
4. Long file/collection names are not properly truncated, making the interface hard to read

## Solution Statement

Implement four key enhancements:
1. **Shift+Space Recursive Toggle**: Add recursive folder toggle that affects all files in subdirectories
2. **R Key Refresh**: Add manual refresh to reload .guardfile and rescan file tree
3. **Symlink Visual Handling**: Display symlinks in gray and make them non-interactive
4. **Smart Text Truncation**: Truncate long names in the middle with ellipsis, preserving start and end

## Feature Metadata

**Feature Type**: Enhancement
**Estimated Complexity**: Medium
**Primary Systems Affected**: TUI Layer, Manager Layer, Filesystem Layer
**Dependencies**: Existing Bubble Tea framework, lipgloss styling, filesystem operations

---

## CONTEXT REFERENCES

### Relevant Codebase Files IMPORTANT: YOU MUST READ THESE FILES BEFORE IMPLEMENTING!

- `internal/tui/keys.go` (lines 1-50) - Why: Contains KeyMap structure and key binding definitions to extend
- `internal/tui/keys.go` (lines 60-120) - Why: Contains handleKeyPress function that needs new key handlers
- `internal/tui/models.go` (lines 80-120) - Why: Contains loadData method that needs refresh capability
- `internal/tui/items.go` (lines 1-50) - Why: Contains item rendering logic that needs symlink detection and truncation
- `internal/tui/items.go` (lines 60-140) - Why: Contains delegate Render methods that need visual enhancements
- `internal/filesystem/filesystem.go` (lines 250-280) - Why: Contains CollectFilesRecursive method for recursive operations
- `internal/security/validator.go` (lines 50-65) - Why: Contains symlink detection logic to reuse
- `internal/manager/folders.go` (lines 140-200) - Why: Contains existing folder toggle logic to extend for recursive operations

### New Files to Create

- `internal/tui/text_utils.go` - Text truncation utilities with proper Unicode and ANSI handling
- No new manager files needed - recursive operations handled in TUI layer

### Relevant Documentation YOU SHOULD READ THESE BEFORE IMPLEMENTING!

- [Bubble Tea Key Handling](https://github.com/charmbracelet/bubbletea/blob/master/key.go)
  - Specific section: Key binding and message handling
  - Why: Required for implementing new keyboard shortcuts
- [Lipgloss Styling](https://github.com/charmbracelet/lipgloss#colors)
  - Specific section: Color and styling options
  - Why: Needed for gray symlink styling
- [Go filepath.WalkDir](https://pkg.go.dev/path/filepath#WalkDir)
  - Specific section: Recursive directory traversal
  - Why: Understanding recursive file collection patterns

### Patterns to Follow

**Key Binding Pattern:**
```go
// From internal/tui/keys.go lines 25-30
Toggle: key.NewBinding(
    key.WithKeys("t"),
    key.WithHelp("t", "toggle protection"),
),
```

**Manager Method Pattern:**
```go
// From internal/manager/folders.go lines 144-150
func (m *Manager) ToggleFolders(paths []string) error {
    if len(paths) == 0 {
        return fmt.Errorf("no folders specified")
    }
    // Process logic...
}
```

**Item Rendering Pattern:**
```go
// From internal/tui/items.go lines 70-85
func (d FileItemDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
    i, ok := listItem.(FileItem)
    if !ok {
        return
    }
    str := i.Title()
    if index == m.Index() {
        str = "> " + str
    } else {
        str = "  " + str
    }
    fmt.Fprint(w, str)
}
```

**Symlink Detection Pattern:**
```go
// From internal/security/validator.go lines 51-62
info, err := os.Lstat(path)
if err != nil {
    return nil // Handle error appropriately
}
if info.Mode()&os.ModeSymlink != 0 {
    // Handle symlink case
}
```

---

## IMPLEMENTATION PLAN

### Phase 1: Foundation

Set up utility functions and extend key bindings to support new functionality.

**Tasks:**
- Create text truncation utility with middle ellipsis
- Create symlink detection helper functions
- Extend KeyMap structure with new key bindings
- Add recursive toggle method to manager layer

### Phase 2: Core Implementation

Implement the main keyboard shortcuts and visual enhancements.

**Tasks:**
- Implement Shift+Space recursive toggle handler
- Implement R key refresh handler
- Add symlink visual styling to item delegates
- Integrate text truncation into item title methods

### Phase 3: Integration

Connect new functionality with existing TUI components and ensure proper data flow.

**Tasks:**
- Update handleKeyPress to route new key combinations
- Modify loadData to support refresh scenarios
- Update help screen with new keyboard shortcuts
- Ensure proper error handling and user feedback

### Phase 4: Testing & Validation

Validate all new functionality works correctly and doesn't break existing features.

**Tasks:**
- Test recursive toggle operations on nested folder structures
- Test refresh functionality with external .guardfile changes
- Test symlink visual display and non-interaction
- Test text truncation with various path lengths

---

## STEP-BY-STEP TASKS

IMPORTANT: Execute every task in order, top to bottom. Each task is atomic and independently testable.

### CREATE internal/tui/text_utils.go

- **IMPLEMENT**: Text utilities with proper Unicode and ANSI width calculation
- **IMPLEMENT**: Functions: StringWidth(), TruncateMiddle(), TruncateRight(), PadRight(), PadLeft()
- **IMPLEMENT**: Symlink detection helper that checks if path is symlink
- **PATTERN**: Follow Go utility function patterns with clear error handling
- **IMPORTS**: `os`, `path/filepath`, `github.com/mattn/go-runewidth`, `github.com/muesli/ansi`
- **GOTCHA**: Use three dots `...` not Unicode ellipsis, handle ANSI escape sequences
- **VALIDATE**: `go build ./internal/tui && echo "Text utils compilation successful"`

### REMOVE internal/manager/recursive.go task

- **REASON**: Recursive operations should be in TUI layer, not manager layer
- **ALTERNATIVE**: Use existing `filesystem.CollectFilesRecursive()` directly from TUI

### UPDATE internal/tui/keys.go

- **ADD**: RecursiveToggle key binding with fallback keys: "shift+space", "ctrl+space", "ctrl+@"
- **ADD**: Refresh key binding accepting both "r" and "R"
- **PATTERN**: Follow existing key binding pattern from lines 25-35
- **IMPORTS**: No new imports needed
- **GOTCHA**: Include multiple key combinations for better compatibility
- **VALIDATE**: `go build ./internal/tui && echo "Keys compilation successful"`

### UPDATE internal/tui/keys.go

- **ADD**: Handler cases for RecursiveToggle and Refresh in handleKeyPress function
- **IMPLEMENT**: Call appropriate manager methods and update UI state
- **PATTERN**: Mirror existing toggle handler pattern from lines 85-95
- **IMPORTS**: No new imports needed
- **GOTCHA**: Ensure proper error handling and user feedback messages
- **VALIDATE**: `go build ./internal/tui && echo "Key handlers compilation successful"`

### UPDATE internal/tui/models.go

- **ADD**: RefreshData method that reloads registry and rescans file tree
- **REFACTOR**: Extract common loading logic from loadData for reuse
- **PATTERN**: Follow existing loadData pattern from lines 80-120
- **IMPORTS**: No new imports needed
- **GOTCHA**: Handle registry reload errors gracefully, preserve UI state where possible
- **VALIDATE**: `go build ./internal/tui && echo "Models compilation successful"`

### UPDATE internal/tui/items.go

- **ADD**: FileNode struct with IsSymlink bool field for proper symlink tracking
- **ADD**: Symlink detection to FileItem.Title() method using text_utils
- **ADD**: Text truncation to all Title() methods using text_utils
- **IMPLEMENT**: ColorSymlink = lipgloss.Color("240") and ItemSymlink style
- **PATTERN**: Follow existing Title() method pattern from lines 15-25
- **IMPORTS**: Add lipgloss for styling, text_utils for utilities
- **GOTCHA**: Preserve guard status indicators, handle empty collections list
- **VALIDATE**: `go build ./internal/tui && echo "Items compilation successful"`

### UPDATE internal/tui/items.go

- **MODIFY**: All delegate Render methods to apply symlink styling with Color("240")
- **ADD**: Conditional gray color application for symlink items using ItemSymlink style
- **IMPLEMENT**: Right arrow navigation should NOT traverse into symlinked folders
- **PATTERN**: Follow existing Render method pattern from lines 70-85
- **IMPORTS**: Use lipgloss ColorSymlink and ItemSymlink styles
- **GOTCHA**: Only apply gray to symlinks, preserve selection highlighting
- **VALIDATE**: `go build ./internal/tui && echo "Delegates compilation successful"`

### UPDATE internal/tui/keys.go

- **MODIFY**: handleToggle method to skip symlinks (no-op behavior)
- **ADD**: Symlink detection before toggle operations using text_utils
- **IMPLEMENT**: Recursive toggle using filesystem.CollectFilesRecursive() directly
- **PATTERN**: Add early return for symlinks in toggle logic
- **IMPORTS**: Use text_utils for symlink detection, filesystem for recursive collection
- **GOTCHA**: Provide user feedback when symlink toggle is attempted
- **VALIDATE**: `go build ./internal/tui && echo "Toggle handling compilation successful"`

### UPDATE internal/tui/views.go

- **ADD**: New keyboard shortcuts to help screen
- **MODIFY**: Help text to include Shift+Space and R key descriptions
- **UPDATE**: StatusBarHelp to "↑↓:Navigate  ←→:Expand/Collapse  Space:Toggle  Tab:Switch  R:Refresh  Q:Quit"
- **PATTERN**: Follow existing help text format from renderHelp method
- **IMPORTS**: No new imports needed
- **GOTCHA**: Keep help text concise and aligned with existing format
- **VALIDATE**: `go build ./internal/tui && echo "Views compilation successful"`

### REMOVE internal/manager/manager.go task

- **REASON**: No manager delegation needed since recursive operations are in TUI layer

---

## TESTING STRATEGY

### Unit Tests

Create focused unit tests for new utility functions and core logic:

- Text truncation with various input lengths and Unicode characters
- Symlink detection with different file types
- Recursive toggle operations with nested folder structures
- Refresh functionality with registry state changes

### Integration Tests

Test TUI interactions and manager integration:

- Keyboard shortcut handling in TUI environment
- Visual styling application for symlinks
- Data refresh and UI state preservation
- Error handling and user feedback display

### Edge Cases

- Empty folders for recursive toggle
- Broken symlinks for visual display
- Very long paths for truncation
- Registry corruption for refresh operations
- Permission errors during recursive operations

---

## VALIDATION COMMANDS

Execute every command to ensure zero regressions and 100% feature correctness.

### Level 1: Syntax & Style

```bash
go fmt ./internal/tui/...
golangci-lint run ./internal/tui/...
```

### Level 2: Unit Tests

```bash
go test ./internal/tui/...
go test -race ./internal/tui/...
```

### Level 3: Integration Tests

```bash
go build -o build/guard ./cmd/guard
cp build/guard tests/guard
cd tests && ./test-tui-integration.sh
```

### Level 4: Manual Validation

```bash
# Test TUI with new features
./guard -i
# In TUI:
# 1. Navigate to folder, press Shift+Space (recursive toggle)
# 2. Press R (refresh)
# 3. Navigate to symlink, verify gray color and no-op toggle
# 4. Check long paths are truncated with ellipsis
```

### Level 5: Additional Validation

```bash
# Test recursive operations
mkdir -p test_nested/sub1/sub2
touch test_nested/file1.txt test_nested/sub1/file2.txt test_nested/sub1/sub2/file3.txt
./guard init 644 $(whoami) $(id -gn)
./guard add test_nested/
./guard -i
# Verify recursive toggle affects all nested files
```

---

## ACCEPTANCE CRITERIA

- [ ] Shift+Space performs recursive toggle on selected folder affecting all nested files
- [ ] R key refreshes TUI data by reloading .guardfile and rescanning file tree
- [ ] Symlinks display in gray color (ANSI 7) and are non-interactive
- [ ] Long file/collection names truncate in middle with three dots ...
- [ ] All validation commands pass with zero errors
- [ ] Unit test coverage meets requirements (80%+)
- [ ] Integration tests verify end-to-end TUI workflows
- [ ] Code follows project conventions and patterns
- [ ] No regressions in existing TUI functionality
- [ ] Help screen updated with new keyboard shortcuts
- [ ] Error handling provides clear user feedback
- [ ] Performance remains responsive with large folder structures

---

## COMPLETION CHECKLIST

- [ ] All tasks completed in order
- [ ] Each task validation passed immediately
- [ ] All validation commands executed successfully
- [ ] Full test suite passes (unit + integration)
- [ ] No linting or type checking errors
- [ ] Manual TUI testing confirms all features work
- [ ] Acceptance criteria all met
- [ ] Code reviewed for quality and maintainability
- [ ] Help documentation updated
- [ ] No performance regressions observed

---

## NOTES

**Design Decisions:**
- Middle truncation preserves both path start and filename for better usability
- Gray color (ANSI 7) provides clear visual distinction without being distracting
- Shift+Space combination avoids conflicts with existing Space toggle
- R key follows common refresh conventions in terminal applications

**Performance Considerations:**
- Recursive operations may be slow on large directory trees - consider progress feedback
- Symlink detection uses Lstat which is efficient for single file checks
- Text truncation should be cached if performance becomes an issue

**Security Considerations:**
- Recursive operations respect existing symlink rejection in security layer
- Path validation remains enforced for all new operations
- No new security vulnerabilities introduced by visual enhancements

**Future Extensibility:**
- Text truncation utility can be extended for other UI elements
- Recursive operations pattern can be applied to collections
- Symlink handling can be enhanced with additional visual indicators

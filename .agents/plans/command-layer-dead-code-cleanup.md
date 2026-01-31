# Feature: Command Layer Dead Code Cleanup

The following plan should be complete, but its important that you validate documentation and codebase patterns and task sanity before you start implementing.

Pay special attention to naming of existing utils types and models. Import from the right files etc.

## Feature Description

Clean up the command files in cmd/guard/commands/ by removing redundant helper functions and dead code. The goal is to simplify enable.go, disable.go, config.go, init.go, and fix show.go by consolidating all logic directly into the Run handlers without delegating to helper functions.

## User Story

As a developer maintaining the guard-tool codebase
I want to remove redundant helper functions from command files
So that the code is simpler, more maintainable, and follows Go/Cobra best practices

## Problem Statement

The command layer has accumulated redundant helper functions that add unnecessary complexity:
- enable.go has 8 helper functions that duplicate logic
- disable.go has 8 helper functions that mirror enable.go patterns
- config.go has 2 helper functions that could be inlined
- init.go has prompt functions that are no longer needed
- show.go is missing helper functions that should be restored

This creates maintenance overhead and violates the principle of keeping CLI commands simple.

## Solution Statement

Consolidate all command logic directly into the Run handlers, removing helper function indirection. This follows Go/Cobra best practices of keeping business logic separate from CLI parsing while maintaining simple, readable command handlers.

## Feature Metadata

**Feature Type**: Refactor
**Estimated Complexity**: Medium
**Primary Systems Affected**: CLI command layer (cmd/guard/commands/)
**Dependencies**: None (internal refactoring only)

---

## CONTEXT REFERENCES

### Relevant Codebase Files IMPORTANT: YOU MUST READ THESE FILES BEFORE IMPLEMENTING!

- `cmd/guard/commands/enable.go` (lines 75-270) - Why: Contains 8 helper functions to remove
- `cmd/guard/commands/disable.go` (lines 100-240) - Why: Contains 8 helper functions to remove  
- `cmd/guard/commands/config.go` (lines 162-185) - Why: Contains 2 helper functions to remove
- `cmd/guard/commands/init.go` (lines 60-85) - Why: Contains prompt functions to remove
- `cmd/guard/commands/show.go` (lines 1-200) - Why: Missing helper functions to restore
- `tests/test-show-commands.sh` - Why: Defines expected show command behavior
- `tests/test-init.sh` - Why: Defines expected init command behavior requiring all 3 args

### New Files to Create

None - this is a refactoring task that modifies existing files only.

### Relevant Documentation YOU SHOULD READ THESE BEFORE IMPLEMENTING!

- [Cobra Best Practices](https://cobra.dev/docs/how-to-guides/working-with-commands/)
  - Specific section: Command organization and modular design
  - Why: Confirms that business logic should be separate from CLI parsing
- [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments)
  - Specific section: Function length and complexity guidelines
  - Why: Validates approach of inlining simple helper functions

### Patterns to Follow

**Error Handling Pattern:**
```go
// From existing commands - consistent pattern
if err := mgr.LoadRegistry(); err != nil {
    fmt.Fprintf(os.Stderr, "Error: %v\n", err)
    os.Exit(1)
}
```

**Manager Usage Pattern:**
```go
// From existing commands - standard initialization
mgr := manager.NewManager(".guardfile")
```

**Output Pattern:**
```go
// From existing commands - consistent warning/error handling
manager.PrintWarnings(mgr.GetWarnings())
manager.PrintErrors(mgr.GetErrors())
if mgr.HasErrors() {
    os.Exit(1)
}
```

**Argument Validation Pattern:**
```go
// From existing commands - consistent no-args check
if len(args) == 0 {
    fmt.Fprintln(os.Stderr, "Error: No files specified")
    fmt.Fprintln(os.Stderr, "Usage: guard command <args>...")
    fmt.Fprintln(os.Stderr)
    fmt.Fprintln(os.Stderr, "Use 'guard help command' for more information.")
    os.Exit(1)
}
```

---

## IMPLEMENTATION PLAN

### Phase 1: Foundation

Analyze current helper function usage and prepare for consolidation.

**Tasks:**
- Review all helper functions to understand their current usage
- Identify shared logic that should remain in manager layer
- Prepare inline replacements for simple helper functions

### Phase 2: Core Implementation

Remove helper functions and consolidate logic into Run handlers.

**Tasks:**
- Clean up enable.go by removing 8 helper functions
- Clean up disable.go by removing 8 helper functions  
- Clean up config.go by removing 2 helper functions
- Simplify init.go by removing prompt functions
- Restore missing helper functions in show.go

### Phase 3: Integration

Ensure all commands work correctly after refactoring.

**Tasks:**
- Validate all commands maintain existing behavior
- Run comprehensive test suite
- Verify error handling remains consistent

### Phase 4: Testing & Validation

Comprehensive testing to ensure no regressions.

**Tasks:**
- Run all shell integration tests
- Verify CLI output formats remain unchanged
- Test edge cases and error conditions

---

## STEP-BY-STEP TASKS

IMPORTANT: Execute every task in order, top to bottom. Each task is atomic and independently testable.

### UPDATE cmd/guard/commands/enable.go

- **REMOVE**: Functions runEnable, parseEnableArgs, processEnableTargets, processEnableFiles, processEnableCollections, processEnableGeneric, allTargetsAreFiles, processEnableTargetsIndividually (8 functions total)
- **PATTERN**: Consolidate all logic directly into the main Run handler and subcommand Run handlers
- **IMPORTS**: Keep existing imports unchanged
- **GOTCHA**: Preserve exact output format and error messages for test compatibility
- **VALIDATE**: `go build ./cmd/guard && ./guard enable --help`

### UPDATE cmd/guard/commands/disable.go

- **REMOVE**: Functions runDisable, parseDisableArgs, processDisableTargets, processDisableFiles, processDisableCollections, processDisableGeneric, allDisableTargetsAreFiles, processDisableTargetsIndividually (8 functions total)
- **PATTERN**: Mirror the enable.go consolidation approach for consistency
- **IMPORTS**: Keep existing imports unchanged  
- **GOTCHA**: Maintain identical error handling and output patterns as enable.go
- **VALIDATE**: `go build ./cmd/guard && ./guard disable --help`

### UPDATE cmd/guard/commands/config.go

- **REMOVE**: Functions handleBulkConfigUpdate, isOctalMode (2 functions total)
- **PATTERN**: Inline the bulk config logic directly in the set command Run handler
- **IMPORTS**: Keep existing imports unchanged
- **GOTCHA**: Preserve octal mode validation logic exactly as implemented
- **VALIDATE**: `go build ./cmd/guard && ./guard config --help`

### UPDATE cmd/guard/commands/init.go

- **REMOVE**: Functions promptForOwner, promptForGroup (prompt functions)
- **MODIFY**: Require all 3 arguments (mode, owner, group) or fail with error
- **PATTERN**: Follow existing argument validation patterns from other commands
- **IMPORTS**: Keep existing imports unchanged
- **GOTCHA**: Change behavior to require all args instead of prompting - this is intentional
- **VALIDATE**: `go build ./cmd/guard && ./guard init --help`

### UPDATE cmd/guard/commands/show.go

- **ADD**: Functions showAllFiles, showAllCollections, showSpecificFiles, showSpecificCollections
- **PATTERN**: Extract logic from existing Run handlers into helper functions for better organization
- **IMPORTS**: Keep existing imports unchanged
- **GOTCHA**: This is the opposite of other files - we're adding helpers here for better readability
- **VALIDATE**: `go build ./cmd/guard && ./guard show --help`

---

## TESTING STRATEGY

### Unit Tests

No unit tests exist for command layer - this is CLI integration testing only.

### Integration Tests

Shell-based integration tests in tests/ directory validate all CLI behavior:

- `test-show-commands.sh` - Validates show command output formats
- `test-init.sh` - Validates init command behavior  
- `test-enable-auto-detect.sh` - Validates enable command functionality
- `test-disable-auto-detect.sh` - Validates disable command functionality
- `test-config.sh` - Validates config command behavior

### Edge Cases

- Empty argument lists (should show usage and exit 1)
- Invalid arguments (should show appropriate error messages)
- Missing .guardfile (should show appropriate error)
- Permission failures (should be handled gracefully)

---

## VALIDATION COMMANDS

Execute every command to ensure zero regressions and 100% feature correctness.

### Level 1: Syntax & Style

```bash
# Go formatting and syntax check
go fmt ./cmd/guard/commands/
go vet ./cmd/guard/commands/
```

### Level 2: Build Validation

```bash
# Ensure all commands compile successfully
go build ./cmd/guard
```

### Level 3: Integration Tests

```bash
# Run comprehensive shell test suite
cd tests && ./run-all-tests.sh
```

### Level 4: Manual Validation

```bash
# Test each modified command manually
./guard enable --help
./guard disable --help  
./guard config --help
./guard init --help
./guard show --help

# Test basic functionality
./guard init 640 testuser testgroup
./guard add test.txt
./guard enable test.txt
./guard show
./guard disable test.txt
./guard config show
```

### Level 5: Additional Validation (Optional)

```bash
# Test edge cases
./guard enable  # Should show usage error
./guard init    # Should show usage error (new behavior)
./guard show nonexistent  # Should handle gracefully
```

---

## ACCEPTANCE CRITERIA

- [ ] enable.go has 8 helper functions removed, logic consolidated into Run handlers
- [ ] disable.go has 8 helper functions removed, logic consolidated into Run handlers  
- [ ] config.go has 2 helper functions removed, logic inlined into set command
- [ ] init.go requires all 3 arguments (mode, owner, group) or fails
- [ ] show.go has helper functions added: showAllFiles, showAllCollections, showSpecificFiles, showSpecificCollections
- [ ] All validation commands pass with zero errors
- [ ] Full test suite passes (unit + integration)
- [ ] No changes to CLI output formats or error messages
- [ ] Code follows existing patterns and conventions
- [ ] No regressions in existing functionality

---

## COMPLETION CHECKLIST

- [ ] All tasks completed in order
- [ ] Each task validation passed immediately
- [ ] All validation commands executed successfully
- [ ] Full test suite passes (integration tests)
- [ ] No build or formatting errors
- [ ] Manual testing confirms commands work correctly
- [ ] Acceptance criteria all met
- [ ] Code reviewed for quality and maintainability

---

## NOTES

**Design Decision**: This refactoring follows Go/Cobra best practices by keeping CLI commands simple and focused on argument parsing while delegating business logic to the manager layer. The consolidation removes unnecessary indirection without changing functionality.

**Behavioral Change**: init.go will require all 3 arguments instead of prompting. This is intentional to simplify the command and make it more predictable for scripting.

**Exception**: show.go gets helper functions added (opposite of other files) because the current implementation is too complex for the Run handlers and would benefit from extraction.

**Test Strategy**: Relies entirely on shell integration tests since this is CLI-focused refactoring. The comprehensive test suite in tests/ directory validates all expected behaviors.

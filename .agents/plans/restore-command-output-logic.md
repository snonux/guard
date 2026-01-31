# Feature: Restore Command Layer Output Logic

## Feature Description

Restore the missing status message output logic in toggle.go, enable.go, disable.go, and clear.go files that was lost during command layer simplification. These commands need to output status messages like 'Guard enabled for X', 'Registered N file(s)', etc., following the exact patterns specified in the existing refactor-command-layer-patterns.md plan.

## User Story

As a user of the guard CLI tool
I want to see clear status messages when I enable, disable, toggle, or clear protection
So that I know exactly what actions were performed on which files

## Problem Statement

The toggle.go, enable.go, disable.go, and clear.go command files were simplified but lost their status message output functionality. Users now get no feedback about what operations were performed, making the CLI less user-friendly and breaking expected behavior defined in tests.

## Solution Statement

Restore the output logic by adding status tracking and message printing to each command, following the exact patterns documented in refactor-command-layer-patterns.md. This includes tracking registration status before operations, counting changes, and printing appropriate status messages with per-file messages plus collection/folder summaries.

## Feature Metadata

**Feature Type**: Bug Fix/Enhancement
**Estimated Complexity**: Medium
**Primary Systems Affected**: Command Layer (cmd/guard/commands/)
**Dependencies**: None (internal restoration)

---

## CONTEXT REFERENCES

### Relevant Codebase Files IMPORTANT: YOU MUST READ THESE FILES BEFORE IMPLEMENTING!

- `cmd/guard/commands/add.go` (lines 45-95) - Why: Reference pattern for registration counting and status messages
- `cmd/guard/commands/toggle.go` - Why: NEEDS output logic restoration
- `cmd/guard/commands/enable.go` - Why: NEEDS output logic restoration  
- `cmd/guard/commands/disable.go` - Why: NEEDS output logic restoration
- `cmd/guard/commands/clear.go` - Why: NEEDS output logic restoration
- `tests/test-output-format.sh` - Why: Shows expected output patterns for toggle operations
- `tests/test-output-messages.sh` - Why: Shows expected output patterns for enable/disable operations
- `tests/test-e2e-workflow.sh` - Why: Shows expected count-based output patterns

### New Files to Create

None - this is restoration of existing functionality

### Relevant Documentation YOU SHOULD READ THESE BEFORE IMPLEMENTING!

- [Cobra Command Documentation](https://pkg.go.dev/github.com/spf13/cobra#Command)
  - Specific section: Run function patterns
  - Why: Understanding command structure for output integration

### Patterns to Follow

**Toggle Helper Function Pattern:**
```go
// Helper function that tracks state BEFORE toggling using maps
// Returns true on ERROR (hasError semantics), false on success
func toggleFilesWithOutput(mgr *manager.Manager, files []string) bool {
    // Track registration status BEFORE any toggling
    registrationMap := make(map[string]bool)
    guardStateMap := make(map[string]bool)
    
    for _, file := range files {
        registrationMap[file] = mgr.IsRegisteredFile(file)
        if registrationMap[file] {
            if guard, err := mgr.GetRegistry().GetRegisteredFileGuard(file); err == nil {
                guardStateMap[file] = guard
            }
        }
    }
    
    // Toggle all files
    if err := mgr.ToggleFiles(files); err != nil {
        fmt.Fprintf(os.Stderr, "Error: %v\n", err)
        return true // ERROR
    }
    
    // Count newly registered files AFTER toggling
    newlyRegistered := 0
    for _, file := range files {
        if !registrationMap[file] && mgr.IsRegisteredFile(file) {
            newlyRegistered++
        }
    }
    
    // Print registration message BEFORE individual status messages
    if newlyRegistered > 0 {
        fmt.Printf("Registered %d file(s)\n", newlyRegistered)
    }
    
    // Print results based on previous state
    for _, file := range files {
        wasRegistered := registrationMap[file]
        wasGuarded := guardStateMap[file]
        
        if !wasRegistered {
            fmt.Printf("Guard enabled for %s\n", file)
        } else if wasGuarded {
            fmt.Printf("Guard disabled for %s\n", file)
        } else {
            fmt.Printf("Guard enabled for %s\n", file)
        }
    }
    
    return false // SUCCESS
}
```

**Enable Command Status Tracking Pattern:**
```go
// Track registration and enable status BEFORE operation
alreadyRegistered := 0
alreadyEnabled := 0
for _, path := range args {
    if mgr.IsRegisteredFile(path) {
        alreadyRegistered++
        if guard, err := mgr.GetRegistry().GetRegisteredFileGuard(path); err == nil && guard {
            alreadyEnabled++
        }
    }
}

// After enable operation
newlyRegistered := 0
nowEnabled := 0
for _, path := range args {
    if mgr.IsRegisteredFile(path) {
        newlyRegistered++
        if guard, err := mgr.GetRegistry().GetRegisteredFileGuard(path); err == nil && guard {
            nowEnabled++
        }
    }
}

// Print status messages
if newlyRegistered > alreadyRegistered {
    fmt.Printf("Registered %d file(s)\n", newlyRegistered-alreadyRegistered)
}
if nowEnabled > alreadyEnabled {
    fmt.Printf("Guard enabled for %d file(s)\n", nowEnabled-alreadyEnabled)
}
if alreadyEnabled > 0 {
    fmt.Printf("Skipped %d file(s) already enabled\n", alreadyEnabled)
}
```

**Collection/Folder Output Pattern:**
```go
// For collections, print per-file messages then collection summary
// First: individual file messages for each file in collection
for _, file := range filesInCollection {
    fmt.Printf("Guard enabled for %s\n", file)
}
// Then: collection summary
fmt.Printf("Guard enabled for collection %s\n", collectionName)

// For folders, similar pattern
for _, file := range filesInFolder {
    fmt.Printf("Guard enabled for %s\n", file)
}
fmt.Printf("Guard enabled for folder %s\n", folderName)
```

---

## IMPLEMENTATION PLAN

### Phase 1: Analyze Current State

Review the current simplified command implementations and identify exactly what output logic is missing.

### Phase 2: Restore Toggle Command Output

Add the toggleFilesWithOutput helper function and integrate it into all toggle subcommands.

### Phase 3: Restore Enable Command Output

Add status tracking and output messages to enable commands.

### Phase 4: Restore Disable Command Output

Add status tracking and output messages to disable commands.

### Phase 5: Restore Clear Command Output

Add appropriate output messages to clear command operations.

### Phase 6: Validation

Ensure all restored output matches test expectations and maintains functionality.

---

## STEP-BY-STEP TASKS

IMPORTANT: Execute every task in order, top to bottom. Each task is atomic and independently testable.

### UPDATE cmd/guard/commands/toggle.go

- **IMPLEMENT**: Add toggleFilesWithOutput helper function with state tracking maps
- **PATTERN**: Follow exact helper function pattern from refactor-command-layer-patterns.md
- **IMPORTS**: Ensure fmt import is present for output messages
- **GOTCHA**: Must track registration and guard state BEFORE any toggling operations
- **OUTPUT**: Print "Registered N file(s)" first, then individual "Guard enabled/disabled for X" messages
- **VALIDATE**: `go build -o build/guard ./cmd/guard && echo "Build successful"`

### UPDATE cmd/guard/commands/toggle.go - Integrate Helper

- **IMPLEMENT**: Replace direct mgr.ToggleFiles calls with toggleFilesWithOutput helper in all subcommands
- **PATTERN**: Call helper function and check return value for error handling
- **IMPORTS**: No additional imports needed
- **GOTCHA**: Helper returns bool (true=error, false=success) not error type
- **ERROR_HANDLING**: If helper returns true, skip registry save and exit with error
- **VALIDATE**: `go build -o build/guard ./cmd/guard && echo "Build successful"`

### UPDATE cmd/guard/commands/enable.go - Add File Status Tracking

- **IMPLEMENT**: Add registration and enable status tracking to newEnableFileCmd
- **PATTERN**: Follow enable command status tracking pattern from refactor plan
- **IMPORTS**: No additional imports needed
- **GOTCHA**: Must count "already registered", "already enabled", "newly registered" separately
- **OUTPUT**: Print "Registered N file(s)", "Guard enabled for N file(s)", "Skipped N file(s) already enabled"
- **VALIDATE**: `go build -o build/guard ./cmd/guard && echo "Build successful"`

### UPDATE cmd/guard/commands/enable.go - Add Collection Output

- **IMPLEMENT**: Add per-file and collection summary output messages to newEnableCollectionCmd
- **PATTERN**: Print "Guard enabled for filename" for each file, then "Guard enabled for collection name"
- **IMPORTS**: No additional imports needed
- **GOTCHA**: Must get files in collection and print individual messages before collection summary
- **OUTPUT**: Print per-file messages then "Guard enabled for collection X"
- **VALIDATE**: `go build -o build/guard ./cmd/guard && echo "Build successful"`

### UPDATE cmd/guard/commands/enable.go - Add Folder Output

- **IMPLEMENT**: Add per-file and folder summary output messages to newEnableFolderCmd
- **PATTERN**: Print "Guard enabled for filename" for each file, then "Guard enabled for folder name"
- **IMPORTS**: No additional imports needed
- **GOTCHA**: Must get files in folder and print individual messages before folder summary
- **OUTPUT**: Print per-file messages then "Guard enabled for folder X"
- **VALIDATE**: `go build -o build/guard ./cmd/guard && echo "Build successful"`

### UPDATE cmd/guard/commands/disable.go - Add File Status Tracking

- **IMPLEMENT**: Add status tracking and output messages to newDisableFileCmd
- **PATTERN**: Mirror enable command pattern but for disable operations
- **IMPORTS**: No additional imports needed
- **GOTCHA**: Track "already disabled" vs "newly disabled" counts
- **OUTPUT**: Print "Guard disabled for N file(s)", "Skipped N file(s) already disabled"
- **VALIDATE**: `go build -o build/guard ./cmd/guard && echo "Build successful"`

### UPDATE cmd/guard/commands/disable.go - Add Collection/Folder Output

- **IMPLEMENT**: Add per-file and summary output messages for collection and folder disable operations
- **PATTERN**: Print "Guard disabled for filename" for each file, then "Guard disabled for collection/folder name"
- **IMPORTS**: No additional imports needed
- **GOTCHA**: Must get files in collection/folder and print individual messages before summary
- **OUTPUT**: Print per-file messages then collection/folder summary
- **VALIDATE**: `go build -o build/guard ./cmd/guard && echo "Build successful"`

### UPDATE cmd/guard/commands/clear.go - Add Output Messages

- **IMPLEMENT**: Add status messages for clear operations showing what was cleared
- **PATTERN**: Print messages about files disabled and collections cleared
- **IMPORTS**: No additional imports needed
- **GOTCHA**: Clear operations affect both files and collections
- **OUTPUT**: Print messages about disabled files and cleared collections
- **VALIDATE**: `go build -o build/guard ./cmd/guard && echo "Build successful"`

### UPDATE cmd/guard/commands/enable.go - Main Command Output

- **IMPLEMENT**: Add output logic to main enable command (auto-detection version)
- **PATTERN**: Combine file, folder, and collection output appropriately
- **IMPORTS**: No additional imports needed
- **GOTCHA**: Main command handles mixed types, need to aggregate output properly
- **OUTPUT**: Print appropriate messages for each type of target enabled
- **VALIDATE**: `go build -o build/guard ./cmd/guard && echo "Build successful"`

### UPDATE cmd/guard/commands/disable.go - Main Command Output

- **IMPLEMENT**: Add output logic to main disable command (auto-detection version)
- **PATTERN**: Combine file, folder, and collection output appropriately
- **IMPORTS**: No additional imports needed
- **GOTCHA**: Main command handles mixed types, need to aggregate output properly
- **OUTPUT**: Print appropriate messages for each type of target disabled
- **VALIDATE**: `go build -o build/guard ./cmd/guard && echo "Build successful"`

### UPDATE cmd/guard/commands/toggle.go - Collection/Folder Output

- **IMPLEMENT**: Add per-file and summary output for collection and folder toggle operations
- **PATTERN**: Print per-file messages then collection/folder summary like other commands
- **IMPORTS**: No additional imports needed
- **GOTCHA**: Toggle operations need to track previous state for correct enable/disable messages
- **OUTPUT**: Print per-file "Guard enabled/disabled for X" then "Guard toggled for collection/folder Y"
- **VALIDATE**: `go build -o build/guard ./cmd/guard && echo "Build successful"`

---

## TESTING STRATEGY

### Unit Tests

No new unit tests required - this restores existing functionality that should be covered by integration tests.

### Integration Tests

Use existing shell integration tests to verify output format correctness:

```bash
cd tests && ./test-output-format.sh
cd tests && ./test-output-messages.sh
cd tests && ./test-e2e-workflow.sh
```

### Edge Cases

- Verify toggle operations show correct "enabled" vs "disabled" messages
- Ensure registration counting works correctly for mixed registered/unregistered files
- Confirm collection and folder operations show per-file messages then summary messages

---

## VALIDATION COMMANDS

Execute every command to ensure zero regressions and 100% output correctness.

### Level 1: Syntax & Style

```bash
go fmt ./...
golangci-lint run
```

### Level 2: Build Validation

```bash
go build -o build/guard ./cmd/guard
```

### Level 3: Integration Tests

```bash
cp build/guard tests/guard
cd tests && ./run-all-tests.sh
```

### Level 4: Manual Validation

```bash
# Test toggle output
./guard init 000 flo staff
touch test.txt
./guard toggle test.txt  # Should show "Guard enabled for test.txt"
./guard toggle test.txt  # Should show "Guard disabled for test.txt"

# Test enable output
./guard enable test.txt  # Should show count messages

# Test disable output  
./guard disable test.txt # Should show count messages

# Test clear output
./guard create mycoll
./guard update mycoll add test.txt
./guard clear mycoll     # Should show clear messages
```

### Level 5: Output Format Validation

```bash
# Run specific output format tests
cd tests && ./test-output-format.sh
cd tests && ./test-output-messages.sh
```

---

## ACCEPTANCE CRITERIA

- [ ] Toggle commands print "Guard enabled for X" when enabling individual files
- [ ] Toggle commands print "Guard disabled for X" when disabling individual files
- [ ] Toggle commands print "Registered N file(s)" when registering new files
- [ ] Enable file command prints "Registered N file(s)" for newly registered files
- [ ] Enable file command prints "Guard enabled for N file(s)" for newly enabled files
- [ ] Enable file command prints "Skipped N file(s) already enabled" for already enabled files
- [ ] Enable collection command prints "Guard enabled for filename" for each file in collection
- [ ] Enable collection command prints "Guard enabled for collection name" summary
- [ ] Enable folder command prints "Guard enabled for filename" for each file in folder  
- [ ] Enable folder command prints "Guard enabled for folder name" summary
- [ ] Disable commands print "Guard disabled for N file(s)" for newly disabled files
- [ ] Disable commands print "Skipped N file(s) already disabled" for already disabled files
- [ ] Disable collection command prints "Guard disabled for filename" for each file in collection
- [ ] Disable collection command prints "Guard disabled for collection name" summary
- [ ] Clear commands print appropriate messages about cleared collections
- [ ] All output messages match test expectations exactly
- [ ] All validation commands pass with zero errors
- [ ] Full integration test suite passes
- [ ] No functional regressions introduced
- [ ] Output format matches existing CLI specifications

---

## COMPLETION CHECKLIST

- [ ] toggleFilesWithOutput helper function implemented with state tracking maps
- [ ] All toggle subcommands use helper function for output
- [ ] Enable file subcommand has registration and enable status tracking
- [ ] Enable collection subcommand has per-file and collection summary output
- [ ] Enable folder subcommand has per-file and folder summary output
- [ ] Disable commands have appropriate status tracking and output
- [ ] Disable collection/folder subcommands have per-file and summary output
- [ ] Clear commands have output messages for operations performed
- [ ] Main auto-detection commands aggregate output appropriately
- [ ] All build validations pass
- [ ] Integration tests pass with correct output format
- [ ] Manual testing confirms expected output messages
- [ ] No functional changes to command behavior beyond output

---

## NOTES

**Design Decisions:**
- Following exact patterns from add.go and toggle.go for consistency
- Using per-file messages plus collection/folder summaries to match test expectations
- Maintaining existing command functionality while adding output logic

**Trade-offs:**
- Slightly more complex code in exchange for user-friendly status messages
- Additional state tracking overhead for better user experience

**Key Risks:**
- Breaking existing functionality while adding output logic
- Incorrect state tracking leading to wrong status messages
- Output format mismatches with test expectations

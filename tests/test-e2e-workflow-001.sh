#!/bin/bash

# test-e2e-workflow-001.sh - BASIC WORKFLOW TESTS
# Tests complete workflows through the guard CLI to verify command sequences work correctly

# Source helpers
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
source "$SCRIPT_DIR/helpers-cli.sh"
set -e

# Find guard binary
GUARD_BIN=""
if [ -f "./guard" ]; then
    GUARD_BIN="$(pwd)/guard"
elif command -v guard &> /dev/null; then
    GUARD_BIN="guard"
else
    echo "Error: guard binary not found. Please build it first."
    exit 1
fi

# ============================================================================
# BASIC WORKFLOW TESTS
# ============================================================================
test_basic_file_workflow() {
    log_test "test_basic_file_workflow" \
             "Complete workflow: init -> add -> enable -> disable -> remove"

    # Init
    output=$($GUARD_BIN init 000 flo staff 2>&1)
    assert_contains "$output" "Initialized .guardfile with:" "Init creates guardfile"

    # Add files
    touch file1.txt file2.txt
    output=$($GUARD_BIN add file file1.txt file2.txt 2>&1)
    assert_contains "$output" "Registered 2 file(s)" "Add registers files"

    # Enable files
    output=$($GUARD_BIN enable file file1.txt 2>&1)
    assert_contains "$output" "Guard enabled for 1 file(s)" "Enable shows count"

    # Verify permissions changed
    perms=$(get_file_permissions file1.txt)
    assert_equals "000" "$perms" "File has guard permissions (000)"

    # Disable files (may have permission errors on macOS - that's expected)
    set +e
    output=$($GUARD_BIN disable file file1.txt 2>&1)
    set -e
    assert_contains "$output" "Guard disabled" "Disable shows message"

    # Remove files
    output=$($GUARD_BIN remove file file1.txt 2>&1)
    assert_contains "$output" "Removed 1 file(s)" "Remove shows count"
}

# Run test
run_test test_basic_file_workflow
print_test_summary 1

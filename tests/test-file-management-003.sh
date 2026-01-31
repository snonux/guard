#!/bin/bash

# test-file-management-003.sh - TOGGLE FILE TESTS
# Tests file add, remove, toggle, enable, disable operations

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
# TOGGLE FILE TESTS
# ============================================================================
test_toggle_file_positive() {
    log_test "test_toggle_file_positive" \
             "Toggle file guard flag on and off"

    # Setup
    $GUARD_BIN init 000 flo staff
    touch test.txt
    local initial_perms=$(get_file_permissions "test.txt")
    $GUARD_BIN add file test.txt

    # First toggle (enable)
    $GUARD_BIN toggle file test.txt
    local exit_code1=$?
    assert_exit_code $exit_code1 0 "First toggle should succeed"

    local guard_flag1=$(get_guard_flag "$(pwd)/test.txt")
    assert_equals "true" "$guard_flag1" "Guard flag should be true after first toggle"

    local toggled_perms=$(get_file_permissions "test.txt")
    assert_equals "000" "$toggled_perms" "Permissions should be 000 when guarded"

    # Second toggle (disable)
    $GUARD_BIN toggle file test.txt
    local exit_code2=$?
    assert_exit_code $exit_code2 0 "Second toggle should succeed"

    local guard_flag2=$(get_guard_flag "$(pwd)/test.txt")
    assert_equals "false" "$guard_flag2" "Guard flag should be false after second toggle"

    local restored_perms=$(get_file_permissions "test.txt")
    assert_equals "$initial_perms" "$restored_perms" "Permissions should be restored to original"
}

# Run test
run_test test_toggle_file_positive
print_test_summary 1

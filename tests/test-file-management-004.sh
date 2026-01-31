#!/bin/bash

# test-file-management-004.sh - ENABLE FILE TESTS
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
# ENABLE FILE TESTS
# ============================================================================
test_enable_file_positive() {
    log_test "test_enable_file_positive" \
             "Positive test: Enable guard on file"

    # Setup
    $GUARD_BIN init 000 flo staff
    touch test.txt
    $GUARD_BIN add file test.txt

    # Run
    $GUARD_BIN enable file test.txt
    local exit_code=$?

    # Assert
    assert_exit_code $exit_code 0 "guard enable file should succeed"

    local guard_flag=$(get_guard_flag "$(pwd)/test.txt")
    assert_equals "true" "$guard_flag" "Guard flag should be true"

    local perms=$(get_file_permissions "test.txt")
    assert_equals "000" "$perms" "Permissions should be 000"
}

# Run test
run_test test_enable_file_positive
print_test_summary 1

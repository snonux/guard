#!/bin/bash

# test-output-format-003.sh - DISABLE OUTPUT FORMAT TESTS
# Verifies that CLI output matches the formats specified in CLI-INTERFACE-SPECS.md
# These tests document gaps between spec and implementation - failing tests indicate
# where the implementation needs to be updated to match the spec.

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
# DISABLE OUTPUT FORMAT TESTS
# ============================================================================
test_disable_output_format_single_file() {
    log_test "test_disable_output_format_single_file" \
             "Disable shows count format 'Guard disabled for N file(s)'"

    $GUARD_BIN init 000 flo staff
    touch myfile.txt
    $GUARD_BIN enable file myfile.txt >/dev/null 2>&1

    output=$($GUARD_BIN disable myfile.txt 2>&1)
    local exit_code=$?

    assert_exit_code $exit_code 0 "Disable should succeed"
    assert_contains "$output" "Guard disabled for" "Output should contain 'Guard disabled for'"
    assert_contains "$output" "file(s)" "Output should contain 'file(s)' count format"
}

# Run test
run_test test_disable_output_format_single_file
print_test_summary 1

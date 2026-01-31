#!/bin/bash

# test-output-format-004.sh - PARTIAL/SKIP OUTPUT FORMAT TESTS
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
# PARTIAL/SKIP OUTPUT FORMAT TESTS
# ============================================================================
test_enable_skipped_output_format() {
    log_test "test_enable_skipped_output_format" \
             "Enable already-enabled shows skip message"

    $GUARD_BIN init 000 flo staff
    touch myfile.txt
    $GUARD_BIN enable file myfile.txt >/dev/null 2>&1

    # Enable again - should show skipped
    output=$($GUARD_BIN enable myfile.txt 2>&1)
    local exit_code=$?

    assert_exit_code $exit_code 0 "Enable should succeed (idempotent)"
    assert_contains "$output" "Skipped" "Output should contain 'Skipped' message"
    assert_contains "$output" "already enabled" "Output should indicate already enabled"
}

# Run test
run_test test_enable_skipped_output_format
print_test_summary 1

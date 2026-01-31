#!/bin/bash

# test-update-005.sh - UPDATE ERROR CASES
# Tests modifying collection membership with: guard update <collection> add|remove <files>...
# Replaces: guard add file ... to ... and guard remove file ... from ...

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
# UPDATE ERROR CASES
# ============================================================================
test_update_no_args() {
    log_test "test_update_no_args" \
             "Negative test: guard update without arguments"

    # Setup
    $GUARD_BIN init 000 flo staff

    # Run
    set +e
    $GUARD_BIN update > /dev/null 2>&1
    local exit_code=$?
    set -e

    # Assert
    assert_exit_code $exit_code 1 "guard update without args should fail"
}

# Run test
run_test test_update_no_args
print_test_summary 1

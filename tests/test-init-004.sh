#!/bin/bash

# test-init-004.sh - Test 4: Init with mode out of range (Negative)
# Tests initialization of the guard system with various parameters

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
# Test 4: Init with mode out of range (Negative)
# ============================================================================
test_init_mode_out_of_range() {
    log_test "test_init_mode_out_of_range" \
             "Negative test: guard init with mode > 777"

    # Run guard init with mode 888
    set +e
    $GUARD_BIN init 888 flo staff > /dev/null 2>&1
    local exit_code=$?
    set -e

    # Assert exit code 1
    assert_exit_code $exit_code 1 "guard init with mode 888 should fail"

    # Assert .guardfile not created
    assert_guardfile_not_exists ".guardfile should not be created"
}

# Run test
run_test test_init_mode_out_of_range
print_test_summary 1

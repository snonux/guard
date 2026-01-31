#!/bin/bash

# test-init-001.sh - Test 1: Init with all arguments (Positive)
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
# Test 1: Init with all arguments (Positive)
# ============================================================================
test_init_with_all_args() {
    log_test "test_init_with_all_args" \
             "Positive test: guard init with mode, owner, and group"

    # Run guard init
    $GUARD_BIN init 000 flo staff
    local exit_code=$?

    # Assert exit code
    assert_exit_code $exit_code 0 "guard init should succeed"

    # Assert .guardfile exists
    assert_guardfile_exists ".guardfile should be created"

    # Assert config values
    local mode=$(get_guard_mode_from_config)
    assert_equals "0000" "$mode" "guard_mode should be 0000"

    # Note: We could also check owner/group in config, but focusing on mode for now
}

# Run test
run_test test_init_with_all_args
print_test_summary 1

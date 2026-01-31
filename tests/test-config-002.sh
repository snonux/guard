#!/bin/bash

# test-config-002.sh - Test 2: Config show without .guardfile (Negative)
# Tests configuration management (show and set operations)

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
# Test 2: Config show without .guardfile (Negative)
# ============================================================================
test_config_show_no_guardfile() {
    log_test "test_config_show_no_guardfile" \
             "Negative test: guard config show fails without .guardfile"

    # Ensure no .guardfile exists (setup_test_env ensures clean state)

    # Run guard config show (should fail)
    set +e
    $GUARD_BIN config show > /dev/null 2>&1
    local exit_code=$?
    set -e

    # Assert exit code 1
    assert_exit_code $exit_code 1 "guard config show without .guardfile should fail"
}

# Run test
run_test test_config_show_no_guardfile
print_test_summary 1

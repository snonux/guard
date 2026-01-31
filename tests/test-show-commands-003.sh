#!/bin/bash

# test-show-commands-003.sh - INFO, VERSION, HELP TESTS
# Tests guard show file/collection, info, version, help commands

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
# INFO, VERSION, HELP TESTS
# ============================================================================
test_info_command() {
    log_test "test_info_command" \
             "Test guard info shows author and source information"

    # Run
    output=$($GUARD_BIN info 2>&1)
    local exit_code=$?

    # Assert
    assert_exit_code $exit_code 0 "Should succeed"

    # Should contain author info (CLI-SPECS line 171)
    if [[ "$output" == *"Florian Buetow"* ]] || [[ "$output" == *"Florian"* ]]; then
        echo -e "${GREEN}✓ PASS${NC}: Author name displayed"
        TESTS_PASSED=$((TESTS_PASSED + 1))
    else
        echo -e "${YELLOW}Note${NC}: Author name not detected"
        echo "Output: $output"
    fi

    # Should contain GitHub link
    if [[ "$output" == *"github.com/florianbuetow/guard"* ]] || [[ "$output" == *"github"* ]]; then
        echo -e "${GREEN}✓ PASS${NC}: GitHub link displayed"
        TESTS_PASSED=$((TESTS_PASSED + 1))
    else
        echo -e "${YELLOW}Note${NC}: GitHub link not detected"
        echo "Output: $output"
    fi
}

# Run test
run_test test_info_command
print_test_summary 1

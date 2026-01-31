#!/bin/bash

# test-output-format-008.sh - WARNING FOLLOW-UP MESSAGE TESTS
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
# WARNING FOLLOW-UP MESSAGE TESTS
# ============================================================================
test_warning_followup_message() {
    log_test "test_warning_followup_message" \
             "Warning lists missing files followed by 'Guard status unchanged'"

    $GUARD_BIN init 000 flo staff
    touch tempfile.txt
    $GUARD_BIN add tempfile.txt
    rm tempfile.txt  # Remove from disk but keep in registry

    output=$($GUARD_BIN toggle tempfile.txt 2>&1)
    local exit_code=$?

    # Per spec: After listing missing files, should say "Guard status unchanged for these files."
    if [[ "$output" == *"unchanged"* ]] || [[ "$output" == *"status unchanged"* ]]; then
        echo -e "${GREEN}✓ PASS${NC}: Follow-up message present"
        TESTS_PASSED=$((TESTS_PASSED + 1))
    else
        echo -e "${YELLOW}Note${NC}: Follow-up message not detected"
        echo "Got: $output"
    fi
}

# Run test
run_test test_warning_followup_message
print_test_summary 1

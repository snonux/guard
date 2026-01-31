#!/bin/bash

# test-enable-auto-detect-005.sh - ENABLE AUTO-DETECTION TESTS - NOT FOUND
# Tests auto-detection of files vs collections: guard enable <arg>...
# Without explicit 'file' or 'collection' keyword

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
# ENABLE AUTO-DETECTION TESTS - NOT FOUND
# ============================================================================
test_enable_auto_detect_not_found() {
    log_test "test_enable_auto_detect_not_found" \
             "Error when neither file nor collection exists"

    # Setup
    $GUARD_BIN init 000 flo staff

    # Run enable on non-existent target
    set +e
    output=$($GUARD_BIN enable nonexistent 2>&1)
    local exit_code=$?
    set -e

    # Assert: Should fail with not found error
    assert_exit_code $exit_code 1 "guard enable should fail for non-existent target"

    # Check for not found error message
    if [[ "$output" == *"not found"* ]] || [[ "$output" == *"not exist"* ]]; then
        echo -e "${GREEN}✓ PASS${NC}: Not found error message displayed"
        TESTS_PASSED=$((TESTS_PASSED + 1))
    else
        echo -e "${RED}✗ FAIL${NC}: Not found error message not found"
        echo "Got: $output"
        TESTS_FAILED=$((TESTS_FAILED + 1))
    fi
}

# Run test
run_test test_enable_auto_detect_not_found
print_test_summary 1

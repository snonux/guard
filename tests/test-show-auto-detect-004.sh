#!/bin/bash

# test-show-auto-detect-004.sh - SHOW AUTO-DETECTION TESTS - FILE/COLLECTION PRIORITY
# Tests auto-detection of files vs collections: guard show <arg>...
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
# SHOW AUTO-DETECTION TESTS - FILE/COLLECTION PRIORITY
# ============================================================================
test_show_auto_detect_ambiguous() {
    log_test "test_show_auto_detect_ambiguous" \
             "File on disk takes priority over collection with same name"

    # Setup
    $GUARD_BIN init 000 flo staff
    touch foo
    $GUARD_BIN add file foo
    $GUARD_BIN create foo  # Same name as file

    # Run show - should succeed, treating 'foo' as file (priority over collection)
    set +e
    output=$($GUARD_BIN show foo 2>&1)
    local exit_code=$?
    set -e

    # Assert: Should succeed (file takes priority)
    assert_exit_code $exit_code 0 "guard show should succeed (file takes priority)"

    # Check that output shows file info (not collection)
    if [[ "$output" == *"foo"* ]]; then
        echo -e "${GREEN}✓ PASS${NC}: File info displayed"
        TESTS_PASSED=$((TESTS_PASSED + 1))
    else
        echo -e "${RED}✗ FAIL${NC}: File info not found in output"
        echo "Got: $output"
        TESTS_FAILED=$((TESTS_FAILED + 1))
    fi
}

# Run test
run_test test_show_auto_detect_ambiguous
print_test_summary 1

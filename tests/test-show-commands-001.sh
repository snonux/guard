#!/bin/bash

# test-show-commands-001.sh - SHOW FILE TESTS
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
# SHOW FILE TESTS
# ============================================================================
test_show_file_positive() {
    log_test "test_show_file_positive" \
             "Positive test: Show specific files with guard status"

    # Setup
    $GUARD_BIN init 000 flo staff
    touch file1.txt file2.txt
    $GUARD_BIN add file file1.txt file2.txt
    $GUARD_BIN enable file file1.txt
    # file1 is guarded, file2 is unguarded

    # Run
    output=$($GUARD_BIN show file file1.txt file2.txt 2>&1)
    local exit_code=$?

    # Assert
    assert_exit_code $exit_code 0 "Should succeed"

    # Check for G for file1 (guarded)
    if echo "$output" | grep -qE "^G .*file1\.txt"; then
        echo -e "${GREEN}✓ PASS${NC}: file1.txt shown as guarded G"
        TESTS_PASSED=$((TESTS_PASSED + 1))
    else
        echo -e "${RED}✗ FAIL${NC}: file1.txt not shown as guarded G"
        echo "Output: $output"
        TESTS_FAILED=$((TESTS_FAILED + 1))
    fi

    # Check for - for file2 (unguarded)
    if echo "$output" | grep -qE "^- .*file2\.txt"; then
        echo -e "${GREEN}✓ PASS${NC}: file2.txt shown as unguarded -"
        TESTS_PASSED=$((TESTS_PASSED + 1))
    else
        echo -e "${RED}✗ FAIL${NC}: file2.txt not shown as unguarded -"
        echo "Output: $output"
        TESTS_FAILED=$((TESTS_FAILED + 1))
    fi
}

# Run test
run_test test_show_file_positive
print_test_summary 1

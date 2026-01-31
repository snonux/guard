#!/bin/bash

# test-maintenance-001.sh - CLEANUP TESTS
# Tests guard cleanup, reset, uninstall commands

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
# CLEANUP TESTS
# ============================================================================
test_cleanup_removes_empty_collections() {
    log_test "test_cleanup_removes_empty_collections" \
             "Cleanup should remove empty collections from registry"

    # Setup
    $GUARD_BIN init 000 flo staff
    # OLD: $GUARD_BIN add collection empty1 empty2
    # NEW:
    $GUARD_BIN create empty1 empty2

    # Verify collections exist
    local count_before=$(count_collections_in_registry)

    # Run cleanup
    $GUARD_BIN cleanup
    local exit_code=$?

    # Assert
    assert_exit_code $exit_code 0 "Cleanup should succeed"

    # Empty collections should be removed
    local count_after=$(count_collections_in_registry)
    assert_equals "0" "$count_after" "All empty collections should be removed"
}

# Run test
run_test test_cleanup_removes_empty_collections
print_test_summary 1

#!/bin/bash

# test-tui-milestone-1-007.sh - CATEGORY 7: GUARD STATUS INDICATOR TESTS
# Tests the Text User Interface according to TUI-INTERFACE-SPECS-MILESTONE-1.md
#
# Prerequisites:
# - tmux must be installed (tests will fail if not available)
# - guard binary must be built
#
# Usage:
#   ./tests/test-tui-milestone-1.sh

# Source helpers
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
source "$SCRIPT_DIR/helpers-cli.sh"
source "$SCRIPT_DIR/helpers-tui.sh"
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

# Check for tmux (required for TUI tests)
if ! tui_check_tmux; then
    exit 1
fi

# ============================================================================
# CATEGORY 7: GUARD STATUS INDICATOR TESTS
# ============================================================================
test_indicator_G_explicit() {
    log_test "test_indicator_G_explicit" \
             "[G] indicator shown for explicitly guarded files"

    # Setup
    $GUARD_BIN init 000 flo staff
    touch testfile.txt
    $GUARD_BIN add file testfile.txt
    $GUARD_BIN enable file testfile.txt

    # Launch TUI
    tui_start

    # Assert: Should show [G] indicator
    tui_assert_contains "[G]" "Screen shows [G] for guarded file"

    # Cleanup
    tui_stop
}

# Run test
run_test test_indicator_G_explicit
print_test_summary 1

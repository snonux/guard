#!/bin/bash

# test-tui-milestone-1-009.sh - CATEGORY 9: REFRESH TESTS
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
# CATEGORY 9: REFRESH TESTS
# ============================================================================
test_refresh_new_file() {
    log_test "test_refresh_new_file" \
             "R key refreshes to show newly created file"

    # Setup
    $GUARD_BIN init 000 flo staff
    touch initial.txt

    # Launch TUI
    tui_start

    # Create a new file externally
    touch newfile.txt

    # Refresh
    tui_send_keys "r"
    sleep 0.5

    # Assert: New file should be visible
    tui_assert_contains "newfile.txt" "New file visible after refresh"

    # Cleanup
    tui_stop
}

# Run test
run_test test_refresh_new_file
print_test_summary 1

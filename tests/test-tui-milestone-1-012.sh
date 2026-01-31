#!/bin/bash

# test-tui-milestone-1-012.sh - CATEGORY 12: STATUS BAR TESTS
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
# CATEGORY 12: STATUS BAR TESTS
# ============================================================================
test_status_bar_files_panel() {
    log_test "test_status_bar_files_panel" \
             "Status bar shows correct shortcuts for files panel"

    # Setup
    $GUARD_BIN init 000 flo staff
    touch file.txt

    # Launch TUI
    tui_start

    # Assert: Status bar should show collapse/expand shortcut
    tui_assert_contains "Collapse" "Status bar shows Collapse option"

    # Cleanup
    tui_stop
}

# Run test
run_test test_status_bar_files_panel
print_test_summary 1

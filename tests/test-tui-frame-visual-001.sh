#!/bin/bash

# test-tui-frame-visual-001.sh - 
# Tests verify the double-line frame characters and embedded panel names
# as specified in TUI-INTERFACE-SPECS-MILESTONE-1.md
#
# Prerequisites:
# - tmux must be installed
# - guard binary must be built
#
# Usage:
#   ./tests/test-tui-frame-visual.sh

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
# 
# ============================================================================
test_tui_frame_double_line_corners() {
    log_test "test_tui_frame_double_line_corners" \
             "TUI uses double-line box characters for frame corners"

    # Setup
    $GUARD_BIN init 000 flo staff
    touch file1.txt
    $GUARD_BIN add file file1.txt

    # Launch TUI
    tui_start

    # Assert: Double-line corner characters are present
    tui_assert_contains "╔" "Top-left corner uses double-line (╔)"
    tui_assert_contains "╗" "Top-right corner uses double-line (╗)"
    tui_assert_contains "╚" "Bottom-left corner uses double-line (╚)"
    tui_assert_contains "╝" "Bottom-right corner uses double-line (╝)"

    # Cleanup
    tui_stop
}

# Run test
run_test test_tui_frame_double_line_corners
print_test_summary 1

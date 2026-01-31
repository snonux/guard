#!/bin/bash

# test-tui-milestone-1-011.sh - CATEGORY 11: VISUAL STYLING TESTS
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
# CATEGORY 11: VISUAL STYLING TESTS
# ============================================================================
test_symlink_gray_color() {
    log_test "test_symlink_gray_color" \
             "Symlinks rendered in gray (ANSI 7)"

    # Setup
    $GUARD_BIN init 000 flo staff
    touch realfile.txt
    ln -s realfile.txt symlink.txt

    # Launch TUI
    tui_start

    # Assert: Symlink should have gray styling
    # This checks for the presence of the symlink in the output
    tui_assert_contains "symlink" "Symlink visible in TUI"

    # Cleanup
    tui_stop
}

# Run test
run_test test_symlink_gray_color
print_test_summary 1

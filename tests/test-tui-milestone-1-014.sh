#!/bin/bash

# test-tui-milestone-1-014.sh - CATEGORY 14: NAVIGATION EDGE CASES (NEW)
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
# CATEGORY 14: NAVIGATION EDGE CASES (NEW)
# ============================================================================
test_right_arrow_on_file_noop() {
    log_test "test_right_arrow_on_file_noop" \
             "Right arrow on a file does nothing (Spec line 251)"

    # Setup: Create a file (no folders)
    $GUARD_BIN init 000 flo staff
    touch testfile.txt
    $GUARD_BIN add file testfile.txt

    # Launch TUI
    tui_start

    # Capture initial screen state
    local screen_before=$(tui_capture)

    # Send Right key on the file
    tui_send_keys "Right"

    # Capture new screen state
    local screen_after=$(tui_capture)

    # Assert: TUI still running
    tui_assert_running "TUI still running after Right on file"

    # Assert: Screen should be effectively unchanged (file still selected)
    tui_assert_contains "testfile.txt" "File still visible after Right key"

    # Cleanup
    tui_stop
}

# Run test
run_test test_right_arrow_on_file_noop
print_test_summary 1

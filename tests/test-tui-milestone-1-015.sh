#!/bin/bash

# test-tui-milestone-1-015.sh - CATEGORY 15: FOLDER GUARD INDICATORS (NEW)
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
# CATEGORY 15: FOLDER GUARD INDICATORS (NEW)
# ============================================================================
test_folder_indicator_no_collection() {
    log_test "test_folder_indicator_no_collection" \
             "Folder with no collection shows [ ] indicator (Spec line 233)"

    # Setup
    $GUARD_BIN init 000 flo staff
    mkdir -p emptyfolder
    touch emptyfolder/file.txt
    # Do NOT create collection for this folder

    # Launch TUI
    tui_start

    # Assert: Folder should show [ ] indicator
    tui_assert_contains "[ ]" "Folder shows [ ] indicator for no collection"

    # Cleanup
    tui_stop
}

# Run test
run_test test_folder_indicator_no_collection
print_test_summary 1

#!/bin/bash

# test-tui-guard-states-002.sh - FOLDER GUARD STATE INDICATORS
# Tests all guard state indicators for files, folders, and collections
# according to TUI-EFFECTIVE-GUARD-STATES.md
#
# Guard State Indicators:
#   Files:       [G] guarded, [-] unguarded, [ ] untracked
#   Folders:     [G] all guarded, [-] all unguarded, [~] mixed, [ ] no collection
#   Collections: [G] guarded, [g] implicitly guarded, [~] mixed, [-] unguarded
#
# Prerequisites:
# - tmux must be installed
# - guard binary must be built
#
# Usage:
#   ./tests/test-tui-guard-states.sh

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
# FOLDER GUARD STATE INDICATORS
# ============================================================================
test_folder_indicator_no_collection() {
    log_test "test_folder_indicator_no_collection" \
             "Folder shows [ ] when no collection exists for it"

    # Setup
    $GUARD_BIN init 000 flo staff
    mkdir -p myfolder
    touch myfolder/file.txt
    # Don't create a collection for ./myfolder

    # Launch TUI
    tui_start

    # Assert: Should show [ ] for folder without collection
    tui_assert_contains "[ ]" "Folder without collection shows [ ] indicator"

    # Cleanup
    tui_stop
}

# Run test
run_test test_folder_indicator_no_collection
print_test_summary 1

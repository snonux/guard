package commands

import (
	"fmt"
	"os"

	"github.com/florianbuetow/guard/internal/manager"
	"github.com/spf13/cobra"
)

// NewUninstallCmd creates the uninstall command.
// Per Requirement 8.3: Runs reset, cleanup, verifies, and deletes .guardfile.
func NewUninstallCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "uninstall",
		Short: "Reset, cleanup, verify, and delete the .guardfile",
		Long: `Completely remove guard from the current directory.

This command:
1. Runs reset (disable all guards)
2. Runs cleanup (remove empty collections and missing files)
3. Verifies all existing files have restored permissions
4. Deletes the .guardfile only if verification succeeds

If verification fails, the .guardfile is preserved and an error is returned.`,
		Run: func(cmd *cobra.Command, args []string) {
			mgr := manager.NewManager(".guardfile")

			// Load registry
			if err := mgr.LoadRegistry(); err != nil {
				fmt.Fprintf(os.Stderr, "Error: %v\n", err)
				os.Exit(1)
			}

			// Run uninstall (includes reset, cleanup, verification, and deletion)
			if err := mgr.Destroy(); err != nil {
				fmt.Fprintf(os.Stderr, "Error: %v\n", err)

				// Print warnings and errors
				manager.PrintWarnings(mgr.GetWarnings())
				manager.PrintErrors(mgr.GetErrors())

				os.Exit(1)
			}

			// Print warnings (if any)
			manager.PrintWarnings(mgr.GetWarnings())

			// Print errors (if any)
			manager.PrintErrors(mgr.GetErrors())

			// Exit with error code if there were errors
			if mgr.HasErrors() {
				os.Exit(1)
			}

			// Success message is printed by Destroy()
		},
	}
}

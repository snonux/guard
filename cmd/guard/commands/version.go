package commands

import (
	"fmt"

	"github.com/spf13/cobra"
)

// NewVersionCmd creates the version command.
// Per Requirement 7.2: Displays the current version.
func NewVersionCmd(version string) *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Display version information",
		Long:  "Display the current version of the guard binary.",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("guard version %s\n", version)
		},
	}
}

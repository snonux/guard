package commands

import (
	"fmt"

	"github.com/spf13/cobra"
)

// NewInfoCmd creates the info command.
// Per Requirement 7.1: Displays about text with author and source information.
func NewInfoCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "info",
		Short: "Display information about guard",
		Long:  "Display about text with author and source information.",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Guard - File Permission Management Tool")
			fmt.Println()
			fmt.Println("Created by Florian Buetow")
			fmt.Println("Source code available at github.com/florianbuetow/guard")
			fmt.Println()
			fmt.Println("Guard helps protect your files from accidental modifications")
			fmt.Println("by managing file permissions, ownership, and group settings.")
		},
	}
}

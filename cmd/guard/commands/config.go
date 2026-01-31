package commands

import (
	"fmt"
	"os"

	"github.com/florianbuetow/guard/internal/manager"
	"github.com/spf13/cobra"
)

// NewConfigCmd creates the config command with show and set subcommands.
func NewConfigCmd() *cobra.Command {
	configCmd := &cobra.Command{
		Use:   "config {show|set}",
		Short: "Manage guard configuration",
		Long:  `View or modify guard configuration settings (mode, owner, group).`,
	}

	// Add show subcommand
	configCmd.AddCommand(newConfigShowCmd())

	// Add set subcommand
	configCmd.AddCommand(newConfigSetCmd())

	return configCmd
}

// newConfigShowCmd creates the "config show" subcommand.
func newConfigShowCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "show",
		Short: "Display current configuration",
		Long: `Display the current guard configuration including mode, owner, and group.

Output format:
  Mode:  <octal permission>
  Owner: <username or (empty)>
  Group: <group name or (empty)>`,
		Run: func(cmd *cobra.Command, args []string) {
			mgr := manager.NewManager(".guardfile")

			// Load registry
			if err := mgr.LoadRegistry(); err != nil {
				fmt.Fprintf(os.Stderr, "Error: %v\n", err)
				os.Exit(1)
			}

			// Show config
			if err := mgr.ShowConfig(); err != nil {
				fmt.Fprintf(os.Stderr, "Error: %v\n", err)
				os.Exit(1)
			}
		},
	}
}

// newConfigSetCmd creates the "config set" subcommand.
func newConfigSetCmd() *cobra.Command {
	setCmd := &cobra.Command{
		Use:   "set {mode|owner|group} <value> | <mode> [owner] [group]",
		Short: "Update configuration values",
		Long: `Update guard configuration values.

Single value update:
  guard config set mode <value>   - Set permission mode (octal, 000-777)
  guard config set owner <value>  - Set default owner
  guard config set group <value>  - Set default group

Bulk update (positional):
  guard config set <mode>                 - Update mode only
  guard config set <mode> <owner>         - Update mode and owner
  guard config set <mode> <owner> <group> - Update all three`,
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) == 0 {
				fmt.Fprintln(os.Stderr, "Error: No arguments provided")
				fmt.Fprintln(os.Stderr, "Usage: guard config set {mode|owner|group} <value>")
				fmt.Fprintln(os.Stderr, "   or: guard config set <mode> [owner] [group]")
				fmt.Fprintln(os.Stderr)
				fmt.Fprintln(os.Stderr, "Use 'guard help config set' for more information.")
				os.Exit(1)
			}

			mgr := manager.NewManager(".guardfile")

			// Load registry
			if err := mgr.LoadRegistry(); err != nil {
				fmt.Fprintf(os.Stderr, "Error: %v\n", err)
				os.Exit(1)
			}

			var err error

			// Check if first arg is a keyword (mode/owner/group)
			switch args[0] {
			case "mode":
				if len(args) < 2 {
					fmt.Fprintln(os.Stderr, "Error: mode value required")
					os.Exit(1)
				}
				err = mgr.SetConfigMode(args[1])
			case "owner":
				if len(args) < 2 {
					fmt.Fprintln(os.Stderr, "Error: owner value required")
					os.Exit(1)
				}
				err = mgr.SetConfigOwner(args[1])
			case "group":
				if len(args) < 2 {
					fmt.Fprintln(os.Stderr, "Error: group value required")
					os.Exit(1)
				}
				err = mgr.SetConfigGroup(args[1])
			default:
				// Bulk update: args are positional (mode [owner] [group])
				modeStr := args[0]
				var owner, group *string

				if len(args) > 1 {
					owner = &args[1]
				}
				if len(args) > 2 {
					group = &args[2]
				}

				err = mgr.SetConfig(&modeStr, owner, group)
			}

			if err != nil {
				fmt.Fprintf(os.Stderr, "Error: %v\n", err)
				os.Exit(1)
			}

			// Print warnings
			manager.PrintWarnings(mgr.GetWarnings())

			// Print errors
			manager.PrintErrors(mgr.GetErrors())

			// Exit with error code if there were errors
			if mgr.HasErrors() {
				os.Exit(1)
			}
		},
	}

	return setCmd
}

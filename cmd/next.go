package cmd

import (
	"crumb/helpers"
	"crumb/store"

	"github.com/spf13/cobra"
)

// nextCmd manages the next focus: show, set, or clear.
var nextCmd = &cobra.Command{
	Use:   "next [text]",
	Short: "Get or set the next focus",
	RunE: func(cmd *cobra.Command, args []string) error {
		switch len(args) {
		case 0:
			// Show current focus
			return store.Update(func(data *store.CrumbData) error {
				if data.Next == "" {
					helpers.Info("No current focus set.")
				} else {
					helpers.Info("Current focus: %s", data.Next)
				}
				return nil
			})
		case 1:
			switch args[0] {
			case "clear":
				// Clear focus
				return store.Update(func(data *store.CrumbData) error {
					data.Next = ""
					helpers.Success("Focus cleared.")
					return nil
				})
			default:
				// Set new focus
				return store.Update(func(data *store.CrumbData) error {
					data.Next = args[0]
					helpers.Success("Focus updated.")
					return nil
				})
			}
		default:
			helpers.Error("Usage: crumb next [text]")
			return nil
		}
	},
}

func init() {
	RootCmd.AddCommand(nextCmd)
}
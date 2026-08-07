package cmd

import (
	"crumb/helpers"
	"crumb/store"

	"github.com/spf13/cobra"
)

// ideaCmd manages ideas: add multiple with arguments, list without.
var ideaCmd = &cobra.Command{
	Use:   "idea [text...]",
	Short: "Save or list ideas",
	RunE: func(cmd *cobra.Command, args []string) error {
		switch len(args) {
		case 0:
			// List ideas (newest first)
			return store.Update(func(data *store.CrumbData) error {
				if len(data.Ideas) == 0 {
					helpers.Info("No ideas yet.")
					return nil
				}
				helpers.Info("💡 Ideas (newest first):")
				for i := len(data.Ideas) - 1; i >= 0; i-- {
					helpers.Dim("  %d: %s", len(data.Ideas)-i, data.Ideas[i])
				}
				return nil
			})
		default:
			// Add all provided ideas
			return store.Update(func(data *store.CrumbData) error {
				data.Ideas = append(data.Ideas, args...)
				helpers.Success("Ideas saved.")
				return nil
			})
		}
	},
}

func init() {
	RootCmd.AddCommand(ideaCmd)
}
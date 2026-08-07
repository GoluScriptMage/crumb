package cmd

import (
	"crumb/helpers"
	"crumb/store"

	"github.com/spf13/cobra"
)

// noteCmd manages notes: add with argument, list without.
var noteCmd = &cobra.Command{
	Use:   "note [text]",
	Short: "Save or list notes",
	RunE: func(cmd *cobra.Command, args []string) error {
		switch len(args) {
		case 0:
			// List notes (newest first)
			return store.Update(func(data *store.CrumbData) error {
				if len(data.Notes) == 0 {
					helpers.Info("No notes yet.")
					return nil
				}
				helpers.Info("📝 Notes (newest first):")
				for i := len(data.Notes) - 1; i >= 0; i-- {
					helpers.Dim("  %d: %s", len(data.Notes)-i, data.Notes[i])
				}
				return nil
			})
		case 1:
			// Add single note
			return store.Update(func(data *store.CrumbData) error {
				data.Notes = append(data.Notes, args[0])
				helpers.Success("Note saved.")
				return nil
			})
		default:
			helpers.Error("Usage: crumb note [text]")
			return nil
		}
	},
}

func init() {
	RootCmd.AddCommand(noteCmd)
}
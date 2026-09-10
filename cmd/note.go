package cmd

import (
	"strconv"
	"strings"

	"crumb/helpers"
	"crumb/store"

	"github.com/spf13/cobra"
)

// NoteCmd returns the note command (exported for tests).
func NoteCmd() *cobra.Command { return noteCmd }

// noteCmd manages notes: add, list, delete.
var noteCmd = &cobra.Command{
	Use:   "note [text...|del <id>|del all]",
	Short: "Save, list, or delete notes",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			// Read-only list
			data, err := store.ReadData()
			if err != nil {
				return err
			}
			if len(data.Notes) == 0 {
				helpers.Info("No notes yet.")
				return nil
			}
			helpers.Info("📝 Notes:")
			for i, n := range data.Notes {
				helpers.Dim("  [%d] %s", i+1, n)
			}
			return nil
		}

		// Delete note by index or all
		if args[0] == "del" || args[0] == "delete" {
			if len(args) < 2 {
				helpers.Error("Usage: crumb note del <number> or crumb note del all")
				return nil
			}
			if args[1] == "all" {
				return store.Update(func(data *store.CrumbData) error {
					data.Notes = []string{}
					helpers.Success("All notes deleted.")
					return nil
				})
			}
			idx, err := strconv.Atoi(args[1])
			if err != nil || idx < 1 {
				helpers.Error("Invalid note number: %s", args[1])
				return nil
			}
			return store.Update(func(data *store.CrumbData) error {
				if idx > len(data.Notes) {
					helpers.Error("Note [%d] not found (total %d notes).", idx, len(data.Notes))
					return nil
				}
				deleted := data.Notes[idx-1]
				data.Notes = append(data.Notes[:idx-1], data.Notes[idx:]...)
				helpers.Success("Note [%d] deleted: %s", idx, deleted)
				return nil
			})
		}

		// Add note (joins multiple arguments)
		text := strings.Join(args, " ")
		return store.Update(func(data *store.CrumbData) error {
			data.Notes = append(data.Notes, text)
			helpers.Success("Note saved: %s", text)
			return nil
		})
	},
}

func init() {
	RootCmd.AddCommand(noteCmd)
}

package cmd

import (
	"strconv"
	"strings"

	"crumb/helpers"
	"crumb/store"

	"github.com/spf13/cobra"
)

// IdeaCmd returns the idea command (exported for tests).
func IdeaCmd() *cobra.Command { return ideaCmd }

// ideaCmd manages ideas: add, list, delete.
var ideaCmd = &cobra.Command{
	Use:   "idea [text...|del <id>|del all]",
	Short: "Save, list, or delete ideas",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			// Read-only list
			data, err := store.ReadData()
			if err != nil {
				return err
			}
			if len(data.Ideas) == 0 {
				helpers.Dim("  (no ideas — scratchpad is clear)")
				return nil
			}
			helpers.Info("💡 Ideas:")
			for i, idea := range data.Ideas {
				helpers.Dim("  [%d] %s", i+1, idea)
			}
			return nil
		}

		// Delete idea by index or all
		if args[0] == "del" || args[0] == "delete" {
			if len(args) < 2 {
				helpers.Error("Usage: crumb idea del <number> or crumb idea del all")
				return nil
			}
			if args[1] == "all" {
				return store.Update(func(data *store.CrumbData) error {
					data.Ideas = []string{}
					helpers.Success("All ideas deleted.")
					return nil
				})
			}
			idx, err := strconv.Atoi(args[1])
			if err != nil || idx < 1 {
				helpers.Error("Invalid idea number: %s", args[1])
				return nil
			}
			return store.Update(func(data *store.CrumbData) error {
				if idx > len(data.Ideas) {
					helpers.Error("Idea [%d] not found (total %d ideas).", idx, len(data.Ideas))
					return nil
				}
				deleted := data.Ideas[idx-1]
				data.Ideas = append(data.Ideas[:idx-1], data.Ideas[idx:]...)
				helpers.Success("Idea [%d] deleted: %s", idx, deleted)
				return nil
			})
		}

		// Add idea (joins multiple arguments)
		text := strings.Join(args, " ")
		return store.Update(func(data *store.CrumbData) error {
			data.Ideas = append(data.Ideas, text)
			helpers.Success("Idea saved: %s", text)
			return nil
		})
	},
}

func init() {
	RootCmd.AddCommand(ideaCmd)
}

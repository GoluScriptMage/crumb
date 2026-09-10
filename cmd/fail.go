package cmd

import (
	"crumb/helpers"

	"github.com/spf13/cobra"
)

// FailCmd returns the fail command (exported for tests).
func FailCmd() *cobra.Command { return failCmd }

var failCmd = &cobra.Command{
	Use:   "fail <id>",
	Short: "Mark a task as failed",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			helpers.Error("Usage: crumb fail <id>")
			return nil
		}
		for _, id := range args {
			err := updateTaskStatus(id, "failed", "Task %s marked as failed.")
			if err != nil {
				return err
			}
		}
		return nil
	},
}

func init() {
	RootCmd.AddCommand(failCmd)
}

package cmd

import (
	"crumb/helpers"
	"crumb/store"

	"github.com/spf13/cobra"
)

// RootCmd is the entry point. With no args, shows a dashboard.
var RootCmd = &cobra.Command{
	Use:   "crumb",
	Short: "Crumb — tiny terminal-first developer scratchpad",
	RunE: func(cmd *cobra.Command, args []string) error {
		// Show dashboard: next focus, task count, recent notes/ideas
		return store.Update(func(data *store.CrumbData) error {
			helpers.Info("=== Crumb Dashboard ===")

			// Next focus
			if data.Next != "" {
				helpers.Success("  Focus: %s", data.Next)
			} else {
				helpers.Dim("  Focus: (none)")
			}

			// Active tasks
			pendingCount := 0
			for _, t := range data.Tasks {
				if t.Status == "pending" {
					pendingCount++
				}
			}
			if pendingCount > 0 {
				helpers.Info("  Tasks: %d active", pendingCount)
			} else {
				helpers.Dim("  Tasks: (none)")
			}

			// Recent notes (last 3)
			if len(data.Notes) > 0 {
				helpers.Info("  Recent notes:")
				start := len(data.Notes) - 1
				end := start - 2
				if end < 0 {
					end = 0
				}
				for i := start; i >= end; i-- {
					helpers.Dim("    • %s", data.Notes[i])
				}
			}

			// Recent ideas (last 3)
			if len(data.Ideas) > 0 {
				helpers.Info("  Recent ideas:")
				start := len(data.Ideas) - 1
				end := start - 2
				if end < 0 {
					end = 0
				}
				for i := start; i >= end; i-- {
					helpers.Dim("    • %s", data.Ideas[i])
				}
			}

			return nil
		})
	},
}

// Execute runs the root command.
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		helpers.Error("%v", err)
	}
}
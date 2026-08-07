package cmd

import (
	"fmt"

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
			fmt.Println()
			helpers.Info("=== Crumb Dashboard ===")
			fmt.Println()

			// Next focus
			if data.Next != "" {
				fmt.Printf("  %sFocus:%s %s%s%s\n", helpers.Cyan, helpers.Reset, helpers.Bold+helpers.Green, data.Next, helpers.Reset)
			} else {
				helpers.Dim("  Focus: (none)")
			}
			fmt.Println()

			// Active tasks
			var pendingTasks []store.Task
			for _, t := range data.Tasks {
				if t.Status == "pending" {
					pendingTasks = append(pendingTasks, t)
				}
			}
			if len(pendingTasks) > 0 {
				helpers.Info("  Tasks: %d active", len(pendingTasks))
				for _, t := range pendingTasks {
					fmt.Printf("    %s• [#%s]%s %s%s%s %s(%s)%s\n",
						helpers.Gray, t.ID, helpers.Reset,
						helpers.White, t.Text, helpers.Reset,
						helpers.Yellow, t.Status, helpers.Reset)
				}
			} else {
				helpers.Dim("  Tasks: (none)")
			}
			fmt.Println()

			// Recent notes (last 3)
			if len(data.Notes) > 0 {
				helpers.Info("  Recent notes:")
				start := len(data.Notes) - 1
				end := start - 2
				if end < 0 {
					end = 0
				}
				for i := start; i >= end; i-- {
					fmt.Printf("    %s•%s %s%s%s\n", helpers.Gray, helpers.Reset, helpers.White, data.Notes[i], helpers.Reset)
				}
			}
			fmt.Println()

			// Recent ideas (last 3)
			if len(data.Ideas) > 0 {
				helpers.Info("  Recent ideas:")
				start := len(data.Ideas) - 1
				end := start - 2
				if end < 0 {
					end = 0
				}
				for i := start; i >= end; i-- {
					fmt.Printf("    %s•%s %s%s%s\n", helpers.Gray, helpers.Reset, helpers.White, data.Ideas[i], helpers.Reset)
				}
			}
			fmt.Println()

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

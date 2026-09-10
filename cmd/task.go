package cmd

import (
	"crypto/rand"
	"fmt"

	"crumb/helpers"
	"crumb/store"

	"github.com/spf13/cobra"
)

// GenerateShortId creates a 4-char hex ID for tasks (exported for tests).
func GenerateShortId() string {
	bytes := make([]byte, 4)
	rand.Read(bytes)
	return fmt.Sprintf("%x", bytes)[:4]
}

// TaskCmd returns the task command (exported for tests).
func TaskCmd() *cobra.Command { return taskCmd }

// DoneCmd returns the done command (exported for tests).
func DoneCmd() *cobra.Command { return doneCmd }

// DelCmd returns the del command (exported for tests).
func DelCmd() *cobra.Command { return delCmd }

// DeleteCmd returns the delete alias command (exported for tests).
func DeleteCmd() *cobra.Command { return delCmd }

// taskCmd manages tasks: list or add.
var taskCmd = &cobra.Command{
	Use:   "task [text...]",
	Short: "Manage tasks (list or add)",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			// Read-only list
			data, err := store.ReadData()
			if err != nil {
				return err
			}
			if len(data.Tasks) == 0 {
				helpers.Dim("  (no active tasks — scratchpad is clear)")
				return nil
			}
			helpers.Info("📋 Tasks:")
			helpers.Dim("--------------------------------------------------")
			for i, task := range data.Tasks {
				status := helpers.FormatStatus(task.Status)
				helpers.Dim("  [%d]  %-35s  [#%s]  %s", i+1, task.Text, task.ID, status)
			}
			helpers.Dim("--------------------------------------------------")
			return nil
		}

		// Add new task(s)
		return store.Update(func(data *store.CrumbData) error {
			for _, task := range args {
				id := GenerateShortId()
				newTask := store.Task{
					ID:     id,
					Text:   task,
					Status: "pending",
				}
				data.Tasks = append(data.Tasks, newTask)
				helpers.Success("Task added [#%s]: %s", id, task)
			}
			return nil
		})
	},
}

// doneCmd marks a task as done by ID.
var doneCmd = &cobra.Command{
	Use:   "done [id]",
	Short: "Mark a task as done",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			helpers.Error("Usage: crumb done <id>")
			return nil
		}
		for _, id := range args {
			err := updateTaskStatus(id, "done", "Task %s marked as done.")
			if err != nil {
				return err
			}
		}
		return nil
	},
}

// delCmd deletes a task by ID or clears all with "all".
var delCmd = &cobra.Command{
	Use:     "del [id|all]",
	Aliases: []string{"delete"},
	Short:   "Delete a task by ID or 'all' to delete all",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			helpers.Error("Usage: crumb del <id> or crumb del all")
			return nil
		}
		if args[0] == "all" {
			return store.Update(func(data *store.CrumbData) error {
				data.Tasks = []store.Task{}
				helpers.Success("All tasks deleted.")
				return nil
			})
		}
		for _, id := range args {
			err := updateTaskStatus(id, "delete", "Task %s deleted.")
			if err != nil {
				return err
			}
		}
		return nil
	},
}

// updateTaskStatus finds a task by ID and updates its status in place or deletes it.
func updateTaskStatus(taskID, status, successMsg string) error {
	return store.Update(func(data *store.CrumbData) error {
		found := false
		for i := range data.Tasks {
			if data.Tasks[i].ID == taskID {
				found = true
				if status == "delete" {
					data.Tasks = append(data.Tasks[:i], data.Tasks[i+1:]...)
				} else {
					data.Tasks[i].Status = status
				}
				break
			}
		}

		if !found {
			helpers.Error("Task with ID %s not found.", taskID)
			return nil
		}

		helpers.Success(successMsg, taskID)
		return nil
	})
}

func init() {
	RootCmd.AddCommand(taskCmd)
	RootCmd.AddCommand(doneCmd)
	RootCmd.AddCommand(delCmd)
}

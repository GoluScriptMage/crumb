package cmd

import (
	"crypto/rand"
	"fmt"

	"crumb/helpers"
	"crumb/store"

	"github.com/spf13/cobra"
)

// generateShortId creates a 3-char hex ID for tasks.
func generateShortId() string {
	bytes := make([]byte, 4)
	rand.Read(bytes)
	return fmt.Sprintf("%x", bytes)[:3]
}

// taskCmd manages tasks: add, list, clear.
var taskCmd = &cobra.Command{
	Use:   "task [text|clear]",
	Short: "Manage tasks (add, list, clear)",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			// List tasks
			return store.Update(func(data *store.CrumbData) error {
				if len(data.Tasks) == 0 {
					helpers.Info("📋 No active tasks.")
					return nil
				}
				helpers.Info("📋 Active tasks:")
				helpers.Dim("--------------------------------------------------")
				for i, task := range data.Tasks {
					status := helpers.FormatStatus(task.Status)
					helpers.Dim("  [%d]  %-35s  [#%s]  %s", i+1, task.Text, task.ID, status)
				}
				helpers.Dim("--------------------------------------------------")
				return nil
			})
		}

		// Switch on first argument
		switch args[0] {
		case "clear":
			// Type two times to clear all tasks so we don't accidentally delete them.
			if len(args) > 1 && args[1] == "clear" {
				// Clear all tasks
				return store.Update(func(data *store.CrumbData) error {
					data.Tasks = []store.Task{}
					helpers.Success("All tasks cleared.")
					return nil
				})
			}
		default:
			// Add new task (first arg is the task text)
			return store.Update(func(data *store.CrumbData) error {
				for _, task := range args {
					id := generateShortId()
					newTask := store.Task{
						ID:     id,
						Text:   task,
						Status: "pending",
					}
					data.Tasks = append(data.Tasks, newTask)
					helpers.Success("Task added: %s", task)
				}
				return nil
			})
		}
		return nil
	},
}

// doneCmd marks a task as done by ID (keeps task in list, updates status).
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

// cancelCmd marks a task as canceled by ID (keeps task in list, updates status).
var cancelCmd = &cobra.Command{
	Use:   "cancel [id]",
	Short: "Mark a task as canceled",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			helpers.Error("Usage: crumb cancel <id>")
			return nil
		}
		for _, id := range args {
			err := updateTaskStatus(id, "canceled", "Task %s canceled.")
			if err != nil {
				return err
			}
		}
		return nil
	},
}

// deleteCmd deletes a task by ID (removes it from the list).
var deleteCmd = &cobra.Command{
	Use:   "delete [id]",
	Short: "Delete a task by ID",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			helpers.Error("Usage: crumb delete <id>")
			return nil
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


// Clear all tasks command (for convenience)
var clearCmd = &cobra.Command{
	Use:   "clear [clear]",
	Short: "Clear all tasks",
	RunE: func(cmd *cobra.Command, args []string) error {

		if len(args) == 0 || args[0] != "clear" {
			helpers.Error("Usage: crumb clear clear")
			return nil
		}

		return updateTaskStatus("", "clear", "")
	},
}

// updateTaskStatus finds a task by ID and updates its status in place
func updateTaskStatus(taskID, status, successMsg string) error {
	return store.Update(func(data *store.CrumbData) error {
		found := false
		// For clearing all
		if status == "clear" {
			helpers.Warn("Clearing all tasks...")
			data.Tasks = []store.Task{}
			helpers.Success("All tasks cleared.")
			return nil
		}

		for i := range data.Tasks {
			if data.Tasks[i].ID == taskID {
				data.Tasks[i].Status = status
				found = true
				// Delete the task if the status is "clear" or "delete"
				if status == "delete" {
					data.Tasks = append(data.Tasks[:i], data.Tasks[i+1:]...)
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
	RootCmd.AddCommand(cancelCmd)
	RootCmd.AddCommand(deleteCmd)
	RootCmd.AddCommand(clearCmd)
}

package cmd

import (
	"crumb/store"
	"crypto/rand"
	"fmt"
	"github.com/spf13/cobra"
)

func generateShortId() string {
	bytes := make([]byte, 4)
	rand.Read(bytes)

	return fmt.Sprintf("%x", bytes)[:3]
}

var taskCmd = &cobra.Command{
	Use:   "task [text]",
	Short: "task 'solve 1 problem' is command to save tasks" + "\n" + "task to return all tasks" + "\n" + "task clear to clear all tasks",
	Run: func(cmd *cobra.Command, args []string) {
		// Get the data
		data, err := store.ReadData()
		if err != nil {
			fmt.Println("Cannot get data")
			return
		}

		// If args empty, return all tasks
		if len(args) == 0 {
			// 1. Check if empty
			if len(data.Tasks) == 0 {
				fmt.Println("📋 No active tasks.")
				return
			}

			// 2. Print all tasks
			fmt.Println("📋 Active tasks.")
			fmt.Println("--------------------------------------------------")

			// 3. Print Rows (Looping straight or reverse cleanly)
			for i, task := range data.Tasks {
				// %-2d = 2-space index, %-35s = left-aligned text padded to 35 chars, [%s] = ID tag
				fmt.Printf("[%d]  %-35s  [#%s]\n", i+1, task.Text, task.ID)
			}
			fmt.Println("--------------------------------------------------")
			return
		}

		// If args[0] is clear
		if args[0] == "clear" {
			data.Tasks = []store.Task{}
			err = store.WriteData(data)
			if err != nil {
				fmt.Println("Cannot clear tasks! Try again.")
				return
			}
			return
		}

		// If args[0] is not clear, add new task
		id := generateShortId()
		newTask := store.Task{
			ID:   id,
			Text: args[0],
		}
		data.Tasks = append(data.Tasks, newTask)

		err = store.WriteData(data)
		if err != nil {
			fmt.Println("Cannot write data")
			return
		}
		fmt.Printf("✅ Task added: %s\n", newTask.Text)

	},
}

var doneCmd = &cobra.Command{
	Use:   "done [id]",
	Short: "done [id] is command to mark task as done",
	Run: func(cmd *cobra.Command, args []string) {
		// Get the data
		data, err := store.ReadData()
		if err != nil {
			fmt.Println("Cannot get data")
			return
		}

		if len(args) == 0 {
			fmt.Println("Please provide a task ID to mark as done.")
			return
		}

		taskID := args[0]
		found := false

		for i, task := range data.Tasks {
			if task.ID == taskID {
				data.Done = append(data.Done, task.Text)
				data.Tasks = append(data.Tasks[:i], data.Tasks[i+1:]...)
				found = true
				break
			}
		}

		if !found {
			fmt.Printf("Task with ID %s not found.\n", taskID)
			return
		}

		err = store.WriteData(data)
		if err != nil {
			fmt.Println("Cannot write data")
			return
		}
		fmt.Printf("✅ Task with ID %s marked as done.\n", taskID)
	},
}

func init() {
	RootCmd.AddCommand(taskCmd) // Add the task command to the root command
	RootCmd.AddCommand(doneCmd) // Add the done command to the root command
}

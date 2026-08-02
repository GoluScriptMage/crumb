package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
	"path/filepath"
)

var noteCmd = &cobra.Command{
	Use:   "note [text]",
	Short: "Note is command to save notes",
	Run: func(cmd *cobra.Command, args []string) {
		// Check if any arguments were provided
		if len(args) != 0 {
			fmt.Println(args[0])
		}

		// Get the cwd and read the file
		currentDir, _ := os.Getwd()
		fsPath := filepath.Join(currentDir, "store", "json.go")

		bytes, err := os.ReadFile(fsPath)
		if err != nil{
			fmt.Println("Error reading file:", err)
		}
		fmt.Println(string(bytes))
	},
}

func init() {
	RootCmd.AddCommand(noteCmd)
}

package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var noteCmd = &cobra.Command{
	Use:   "note [text]",
	Short: "Note is command to save notes",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 0 {
			fmt.Println(args[0])
		}
	},
}

func init() {
	RootCmd.AddCommand(noteCmd)
}

package cmd

import (
	"crumb/store"
	"fmt"
	"github.com/spf13/cobra"
)

var nextCmd = &cobra.Command{
	Use:   "next [text]",
	Short: "Next is command to save next actions",
	Run: func(cmd *cobra.Command, args []string) {
		data, err := store.ReadData()

		if err != nil {
			fmt.Println("Error read data")
			return
		}
		// For displaying current focus
		if len(args) == 0 {
			if data.Next == "" {
				fmt.Println("No current focus set.")
			} else {
				fmt.Println("Current focus: ", data.Next)
			}
			return
		}

		// For clearing focus or saving new focus
		if args[0] == "clear" {
			data.Next = ""
			fmt.Println("Focus cleared")
		} else {
			data.Next = args[0]
			fmt.Println("Focus updated")
		}

		// Update the data in the JSON file
		err = store.WriteData(data)

		if err != nil {
			fmt.Println("Cannot write data")
			return
		}
	},
}

func init() {
	RootCmd.AddCommand(nextCmd)
}

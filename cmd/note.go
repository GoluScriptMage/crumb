package cmd

import (
	"crumb/store"
	"fmt"
	"github.com/spf13/cobra"
)

var noteCmd = &cobra.Command{
	Use:   "note [text]",
	Short: "Note is command to save notes",
	Run: func(cmd *cobra.Command, args []string) {

		// get the data
		data, err := store.ReadData()
		if err != nil {
			fmt.Println("Cannot get data")
			return
		}
		// If there are arg to add, append it to notes
		if len(args) > 0 {
			data.Notes = append(data.Notes, args[0])
			err = store.WriteData(data)
			if err != nil {
				fmt.Println("Cannot write data")
				return
			}
		}
		for i := 0; i < len(data.Notes); i++ {
			fmt.Printf("%d: %s\n", i+1, data.Notes[len(data.Notes)-1-i])
		}
	},
}

func init() {
	RootCmd.AddCommand(noteCmd)
}

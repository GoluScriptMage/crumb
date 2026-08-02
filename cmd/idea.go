package cmd

import (
	"crumb/store"
	"fmt"
	"github.com/spf13/cobra"
)

var ideaCmd = &cobra.Command{
	Use:   "idea [text]",
	Short: "Idea is command to save ideas",
	Run: func(cmd *cobra.Command, args []string) {
		data, err := store.ReadData()
		if err != nil {
			fmt.Println("Error reading ideas: ", err)
			return
		}

		if len(args) > 0 {
			// Save all the ideas
			data.Ideas = append(data.Ideas, args...)
			err = store.WriteData(data)

			if err != nil {
				fmt.Println("Cannot write idea")
				return
			}
		}

		for i, idea := range data.Ideas {
			fmt.Printf("%d: %s\n", i+1, idea)
		}
	},
}

func init() {
	RootCmd.AddCommand(ideaCmd)
}

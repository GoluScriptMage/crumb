package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var RootCmd = &cobra.Command{
	Use:   "crumb",
	Short: "Crumb is a tiny terminal-first developer scratchpad",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("=== Crumb Dashboard ===")
	},
}



func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
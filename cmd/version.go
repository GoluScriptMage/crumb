package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var version = "v1.4.0"

// Version returns the current version string (exported for tests).
func Version() string { return version }

// VersionCmd returns the version command (exported for tests).
func VersionCmd() *cobra.Command { return versionCmd }

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the crumb version",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("crumb version", version)
		return nil
	},
}

func init() {
	RootCmd.Version = version
	RootCmd.AddCommand(versionCmd)
}

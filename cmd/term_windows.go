//go:build windows

package cmd

func getTerminalWidth() int {
	return 80
}

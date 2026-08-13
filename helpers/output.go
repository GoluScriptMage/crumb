package helpers

import "fmt"

const (
	Reset  = "\033[0m"
	Bold   = "\033[1m"
	Green  = "\033[32m"
	Yellow = "\033[33m"
	Red    = "\033[31m"
	Cyan   = "\033[36m"
	Gray   = "\033[90m"
	White  = "\033[97m"
)

// Info prints an informational message in cyan.
func Info(format string, args ...any) {
	fmt.Printf(Cyan+format+Reset+"\n", args...)
}

// Success prints a success message in green.
func Success(format string, args ...any) {
	fmt.Printf(Green+format+Reset+"\n", args...)
}

// Error prints an error message in red.
func Error(format string, args ...any) {
	fmt.Printf(Red+format+Reset+"\n", args...)
}

// Warn prints a warning message in yellow.
func Warn(format string, args ...any) {
	fmt.Printf(Yellow+format+Reset+"\n", args...)
}

// Dim prints a dimmed message in gray.
func Dim(format string, args ...any) {
	fmt.Printf(Gray+format+Reset+"\n", args...)
}

// FormatStatus returns a colored status badge (exported for tests).
func FormatStatus(status string) string {
	switch status {
	case "done":
		return Green + "[✔ done]" + Reset
	case "canceled":
		return Red + "[✖ canceled]" + Reset
	case "failed":
		return Red + "[✖ failed]" + Reset
	default: // "pending"
		return Yellow + "[• pending]" + Reset
	}
}
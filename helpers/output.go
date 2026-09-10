package helpers

import "fmt"

const (
	Reset  = "\033[0m"
	Bold   = "\033[1m"
	DimTxt = "\033[2m"
	Green  = "\033[32m"
	Yellow = "\033[33m"
	Red    = "\033[31m"
	Cyan   = "\033[36m"
	Gray   = "\033[90m"
	White  = "\033[97m"
)

// Success prints a subtle confirmation message: dim gray prefix with bold white text.
func Success(format string, args ...any) {
	msg := fmt.Sprintf(format, args...)
	fmt.Printf("%s  ✔  %s%s%s%s\n", Gray, Reset, Bold, msg, Reset)
}

// Error prints a subtle error message: red prefix with muted text.
func Error(format string, args ...any) {
	msg := fmt.Sprintf(format, args...)
	fmt.Printf("%s  ✖  %s%s\n", Red, Reset, msg)
}

// Warn prints a subtle warning message: yellow prefix with muted text.
func Warn(format string, args ...any) {
	msg := fmt.Sprintf(format, args...)
	fmt.Printf("%s  !  %s%s\n", Yellow, Reset, msg)
}

// Info prints a clean informational line.
func Info(format string, args ...any) {
	msg := fmt.Sprintf(format, args...)
	fmt.Println(msg)
}

// Dim prints dimmed/muted gray text.
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

package cmd

import (
	"fmt"
	"math/rand"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"

	"crumb/helpers"
	"crumb/store"

	"github.com/spf13/cobra"
)

// TimerCmd returns the timer command (exported for tests).
func TimerCmd() *cobra.Command { return timerCmd }

type Mascot struct {
	Name  string
	Lines []string
}

var mascots = []Mascot{
	{
		Name: "Sia",
		Lines: []string{
			` /\_/\ `,
			`(˶ᵔ ᵕ ᵔ˶)`,
			` > 🍪 < `,
		},
	},
	{
		Name: "Rabbit",
		Lines: []string{
			`(\__/)`,
			`(•ㅅ•)`,
			`/ づ `,
		},
	},
	{
		Name: "Pikachu",
		Lines: []string{
			`(\__/)`,
			`(o^.^)`,
			`z(_(")(")`,
		},
	},
	{
		Name: "Undead Knight",
		Lines: []string{
			` [•_•] `,
			`<)   )╯ ⚔️`,
			` /   \ `,
		},
	},
}

func visualWidth(s string) int {
	w := 0
	for _, r := range s {
		if r > 0x1F000 || r == '⚔' || r == '🍪' || r == '🎯' || r == '⏳' || r == '✔' || r == '✖' {
			w += 2
		} else {
			w += 1
		}
	}
	return w
}

func printCentered(text string, width int, prefix string, suffix string) {
	vw := visualWidth(text)
	pad := (width - vw) / 2
	if pad < 0 {
		pad = 0
	}
	fmt.Printf("%s%s%s%s\n", strings.Repeat(" ", pad), prefix, text, suffix)
}

func renderTimerFrame(m Mascot, task string, rem int64, total int64, mins int) {
	const boxW = 45
	resetSeq := helpers.Reset

	fmt.Println()
	// 1. Center Mascot
	for _, line := range m.Lines {
		printCentered(line, boxW, helpers.White, resetSeq)
	}
	fmt.Println()

	// 2. Center Task Name
	taskDisplay := fmt.Sprintf("FOCUS: %s", task)
	if len(taskDisplay) > 38 {
		taskDisplay = taskDisplay[:35] + "..."
	}
	printCentered(taskDisplay, boxW, helpers.Bold+helpers.White, resetSeq)
	fmt.Println()

	// 3. Render Aesthetic Box
	innerW := boxW - 2
	fmt.Printf("%s╭%s╮%s\n", helpers.Gray, strings.Repeat("─", innerW), resetSeq)
	fmt.Printf("%s│%s│%s\n", helpers.Gray, strings.Repeat(" ", innerW), resetSeq)

	// Digits
	if rem < 0 {
		rem = 0
	}
	minutes := rem / 60
	seconds := rem % 60
	timeStr := fmt.Sprintf("%02d : %02d", minutes, seconds)
	timePad := (innerW - len(timeStr)) / 2
	timeRemainingSpace := innerW - timePad - len(timeStr)
	fmt.Printf("%s│%s%s%s%s%s%s│%s\n",
		helpers.Gray,
		strings.Repeat(" ", timePad),
		helpers.Bold+helpers.White,
		timeStr,
		resetSeq,
		helpers.Gray,
		strings.Repeat(" ", timeRemainingSpace),
		resetSeq,
	)

	fmt.Printf("%s│%s│%s\n", helpers.Gray, strings.Repeat(" ", innerW), resetSeq)

	// Progress bar
	const barLen = 28
	elapsed := total - rem
	if elapsed < 0 {
		elapsed = 0
	}
	progress := float64(elapsed) / float64(total)
	if progress > 1.0 {
		progress = 1.0
	}
	filled := int(progress * float64(barLen))
	empty := barLen - filled
	if empty < 0 {
		empty = 0
	}

	barStr := strings.Repeat("━", filled)
	if filled > 0 && empty > 0 {
		barStr += "╸"
		empty--
	}
	barStrEmpty := strings.Repeat("─", empty)

	barPad := (innerW - barLen) / 2
	barRightSpace := innerW - barPad - barLen
	fmt.Printf("%s│%s%s%s%s%s%s%s│%s\n",
		helpers.Gray,
		strings.Repeat(" ", barPad),
		helpers.White,
		barStr,
		helpers.Gray,
		barStrEmpty,
		strings.Repeat(" ", barRightSpace),
		resetSeq,
		resetSeq,
	)

	fmt.Printf("%s│%s│%s\n", helpers.Gray, strings.Repeat(" ", innerW), resetSeq)
	fmt.Printf("%s╰%s╯%s\n", helpers.Gray, strings.Repeat("─", innerW), resetSeq)

	// 4. Footer
	footer := fmt.Sprintf("%dm sprint  ·  ctrl+c to background", mins)
	printCentered(footer, boxW, helpers.Gray, resetSeq)
	fmt.Println()
}

var timerCmd = &cobra.Command{
	Use:   "timer [minutes] [task...]",
	Short: "Start, view, or stop the focus timer",
	RunE: func(cmd *cobra.Command, args []string) error {
		// Stop / clear timer
		if len(args) > 0 && (args[0] == "stop" || args[0] == "clear" || args[0] == "del") {
			return store.Update(func(data *store.CrumbData) error {
				data.Timer = nil
				helpers.Success("Timer stopped.")
				return nil
			})
		}

		// Parse minutes if starting a new timer
		if len(args) > 0 {
			mins, err := strconv.Atoi(args[0])
			if err != nil || mins < 1 {
				helpers.Error("Usage: crumb timer <minutes> [task]")
				return nil
			}

			taskName := "Focus Sprint"
			if len(args) > 1 {
				taskName = strings.Join(args[1:], " ")
			}

			err = store.Update(func(data *store.CrumbData) error {
				data.Timer = &store.TimerState{
					Task:      taskName,
					Minutes:   mins,
					StartedAt: time.Now().Unix(),
					Duration:  mins * 60,
				}
				return nil
			})
			if err != nil {
				return err
			}
		}

		// Read current timer
		data, err := store.ReadData()
		if err != nil {
			return err
		}

		if data.Timer == nil {
			helpers.Dim("\n  (no active timer)")
			helpers.Dim("  start a sprint: crumb timer <minutes> <task>")
			helpers.Dim("  example:        crumb timer 50 \"Doing DSA\"\n")
			return nil
		}

		// If running in unit test environment, just render once and return
		if store.IsTestEnv() {
			m := mascots[0]
			rem := int64(data.Timer.Duration) - (time.Now().Unix() - data.Timer.StartedAt)
			renderTimerFrame(m, data.Timer.Task, rem, int64(data.Timer.Duration), data.Timer.Minutes)
			return nil
		}

		// Random mascot selection for this live session
		rand.Seed(time.Now().UnixNano())
		m := mascots[rand.Intn(len(mascots))]

		// Handle Ctrl+C (SIGINT) to detach cleanly
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

		ticker := time.NewTicker(1 * time.Second)
		defer ticker.Stop()

		// Hide cursor
		fmt.Print("\033[?25l")
		defer fmt.Print("\033[?25h")

		for {
			now := time.Now().Unix()
			elapsed := now - data.Timer.StartedAt
			rem := int64(data.Timer.Duration) - elapsed

			// Clear terminal screen and position at top-left
			fmt.Print("\033[H\033[2J")

			renderTimerFrame(m, data.Timer.Task, rem, int64(data.Timer.Duration), data.Timer.Minutes)

			if rem <= 0 {
				fmt.Print("\a") // Terminal bell on completion
				fmt.Print("\033[?25h")
				helpers.Success("Time's up! Great session on %s.", data.Timer.Task)
				fmt.Println()
				return nil
			}

			select {
			case <-sigChan:
				fmt.Print("\033[?25h\n")
				helpers.Dim("  ✔  Timer running in background (re-attach with `crumb timer`).\n")
				return nil
			case <-ticker.C:
			}
		}
	},
}

func init() {
	RootCmd.AddCommand(timerCmd)
}

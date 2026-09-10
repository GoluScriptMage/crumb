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
			` /\_/\\ `,
			`(˶ᵔ ᵕ ᔆ˶)`,
			` /づ🍪 `,
		},
	},
	{
		Name: "Cat",
		Lines: []string{
			` /\_/\\ `,
			`( o.o )`,
			` > ^ < `,
		},
	},
	{
		Name: "Fox",
		Lines: []string{
			` /\   /\ `,
			`(  •ᴗ•  )`,
			` / づ づ `,
		},
	},
	{
		Name: "Bunny",
		Lines: []string{
			`(\   /)`,
			`( •ㅅ• )`,
			`/ づ づ`,
		},
	},
	{
		Name: "Bear",
		Lines: []string{
			` ʕ•ᴥ•ʔ `,
			` /|   |\ `,
			`  |___| `,
		},
	},
	{
		Name: "Frog",
		Lines: []string{
			` @   @ `,
			`( •ᴗ• )`,
			` /|___|\ `,
		},
	},
	{
		Name: "Owl",
		Lines: []string{
			`  ,___,  `,
			` (◉ ◉) `,
			` /|___|\ `,
		},
	},
	{
		Name: "Penguin",
		Lines: []string{
			`  .---.  `,
			` / • • \ `,
			`(   ᴗ   )`,
			` \_____/ `,
		},
	},
	{
		Name: "Ghost",
		Lines: []string{
			` .-"""-. `,
			`/ •   • \`,
			`|   ᴗ   |`,
			`|_/\_/\_|`,
		},
	},
	{
		Name: "Slime",
		Lines: []string{
			`  _____  `,
			` /     \ `,
			`( •ᴗ•  )`,
			` \_____/ `,
		},
	},
	{
		Name: "Robot",
		Lines: []string{
			`┌───────┐`,
			`│ ◉   ◉ │`,
			`│   ▱   │`,
			`└─┬───┬─┘`,
		},
	},
	{
		Name: "Knight",
		Lines: []string{
			`   /‾\   `,
			`  [•_•]  `,
			` <|██|>  `,
			`   / \   `,
		},
	},
	{
		Name: "Wizard",
		Lines: []string{
			`    /\    `,
			`   /__\   `,
			`  (•ᴗ•)  `,
			`  /| |\  `,
		},
	},
	{
		Name: "Ninja",
		Lines: []string{
			`  _____  `,
			` /•   •\ `,
			`(  ᴗ  ) `,
			` /|___|\ `,
		},
	},
	{
		Name: "Samurai",
		Lines: []string{
			`    /‾\   `,
			`   [•_•]  `,
			` <==|██|  `,
			`    / \   `,
		},
	},
	{
		Name: "Pirate",
		Lines: []string{
			`  .---.  `,
			` (•_• )  `,
			` /|___|\ `,
			`  /___\  `,
		},
	},
	{
		Name: "Astronaut",
		Lines: []string{
			` .-------. `,
			`/  •   •  \`,
			`|    ◡    |`,
			` \_______/ `,
		},
	},
	{
		Name: "Demon",
		Lines: []string{
			` /\     /\ `,
			`(  •ᴗ•  )`,
			` \_/\_/  `,
		},
	},
	{
		Name: "Reaper",
		Lines: []string{
			`   .-.    `,
			`  (•_•)   `,
			` /|   |\  `,
			`  |   |___)`,
		},
	},
	{
		Name: "Dragon",
		Lines: []string{
			`    / \___ `,
			` __/  o o \`,
			`/    =ᴥ=  >`,
			`\___/ \___/`,
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

func buildTimerFrame(m Mascot, task string, rem int64, total int64, mins int, termW int) (string, int) {
	const boxW = 45
	resetSeq := helpers.Reset

	indentLen := (termW - boxW) / 2
	if indentLen < 0 {
		indentLen = 0
	}
	indent := strings.Repeat(" ", indentLen)

	var sb strings.Builder
	lineCount := 0

	addLine := func(content string) {
		sb.WriteString(indent)
		sb.WriteString(content)
		sb.WriteString("\033[K\n")
		lineCount++
	}

	addEmptyLine := func() {
		sb.WriteString("\033[K\n")
		lineCount++
	}

	addEmptyLine()

	// 1. Center Mascot
	for _, line := range m.Lines {
		vw := visualWidth(line)
		pad := (boxW - vw) / 2
		if pad < 0 {
			pad = 0
		}
		addLine(fmt.Sprintf("%s%s%s%s", strings.Repeat(" ", pad), helpers.White, line, resetSeq))
	}
	addEmptyLine()

	// 2. Center Task Name
	taskDisplay := fmt.Sprintf("FOCUS: %s", task)
	if len(taskDisplay) > 38 {
		taskDisplay = taskDisplay[:35] + "..."
	}
	vwTask := visualWidth(taskDisplay)
	padTask := (boxW - vwTask) / 2
	if padTask < 0 {
		padTask = 0
	}
	addLine(fmt.Sprintf("%s%s%s%s", strings.Repeat(" ", padTask), helpers.Bold+helpers.White, taskDisplay, resetSeq))
	addEmptyLine()

	// 3. Render Aesthetic Box
	innerW := boxW - 2
	addLine(fmt.Sprintf("%s╭%s╮%s", helpers.Gray, strings.Repeat("─", innerW), resetSeq))
	addLine(fmt.Sprintf("%s│%s│%s", helpers.Gray, strings.Repeat(" ", innerW), resetSeq))

	// Digits
	if rem < 0 {
		rem = 0
	}
	minutes := rem / 60
	seconds := rem % 60
	timeStr := fmt.Sprintf("%02d : %02d", minutes, seconds)
	timePad := (innerW - len(timeStr)) / 2
	timeRemainingSpace := innerW - timePad - len(timeStr)
	addLine(fmt.Sprintf("%s│%s%s%s%s%s%s│%s",
		helpers.Gray,
		strings.Repeat(" ", timePad),
		helpers.Bold+helpers.White,
		timeStr,
		resetSeq,
		helpers.Gray,
		strings.Repeat(" ", timeRemainingSpace),
		resetSeq,
	))

	addLine(fmt.Sprintf("%s│%s│%s", helpers.Gray, strings.Repeat(" ", innerW), resetSeq))

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
	addLine(fmt.Sprintf("%s│%s%s%s%s%s%s%s│%s",
		helpers.Gray,
		strings.Repeat(" ", barPad),
		helpers.White,
		barStr,
		helpers.Gray,
		barStrEmpty,
		strings.Repeat(" ", barRightSpace),
		resetSeq,
		resetSeq,
	))

	addLine(fmt.Sprintf("%s│%s│%s", helpers.Gray, strings.Repeat(" ", innerW), resetSeq))
	addLine(fmt.Sprintf("%s╰%s╯%s", helpers.Gray, strings.Repeat("─", innerW), resetSeq))

	// 4. Footer
	footer := fmt.Sprintf("%dm sprint  ·  ctrl+c to background", mins)
	vwFoot := visualWidth(footer)
	padFoot := (boxW - vwFoot) / 2
	if padFoot < 0 {
		padFoot = 0
	}
	addLine(fmt.Sprintf("%s%s%s%s", strings.Repeat(" ", padFoot), helpers.Gray, footer, resetSeq))
	addEmptyLine()

	return sb.String(), lineCount
}

func renderTimerFrame(m Mascot, task string, rem int64, total int64, mins int) {
	frame, _ := buildTimerFrame(m, task, rem, total, mins, 80)
	fmt.Print(frame)
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

		firstFrame := true
		var lastLineCount int

		for {
			now := time.Now().Unix()
			elapsed := now - data.Timer.StartedAt
			rem := int64(data.Timer.Duration) - elapsed

			termW := getTerminalWidth()
			frame, lineCount := buildTimerFrame(m, data.Timer.Task, rem, int64(data.Timer.Duration), data.Timer.Minutes, termW)

			if !firstFrame && lastLineCount > 0 {
				// Overwrite the last frame in-place!
				fmt.Printf("\033[%dA\r", lastLineCount)
			}
			fmt.Print(frame)
			firstFrame = false
			lastLineCount = lineCount

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

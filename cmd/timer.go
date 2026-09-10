package cmd

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"crumb/helpers"
	"crumb/store"

	"github.com/spf13/cobra"
)

// TimerCmd returns the timer command (exported for tests).
func TimerCmd() *cobra.Command { return timerCmd }

var timerCmd = &cobra.Command{
	Use:   "timer [minutes] [task...]",
	Short: "Start, view, or stop the focus timer",
	RunE: func(cmd *cobra.Command, args []string) error {
		// No args -> show current timer status
		if len(args) == 0 {
			data, err := store.ReadData()
			if err != nil {
				return err
			}

			fmt.Println()
			if data.Timer == nil {
				fmt.Printf("  %s(\\__/)%s\n", helpers.Cyan, helpers.Reset)
				fmt.Printf("  %s(•ㅅ•) / づ%s %s[NO ACTIVE TIMER]%s\n", helpers.Cyan, helpers.Reset, helpers.Gray, helpers.Reset)
				fmt.Println()
				helpers.Dim("  Start a sprint: crumb timer <minutes> <task>")
				helpers.Dim("  Example:        crumb timer 50 \"Doing DSA\"")
				fmt.Println()
				return nil
			}

			now := time.Now().Unix()
			elapsed := now - data.Timer.StartedAt
			rem := int64(data.Timer.Duration) - elapsed

			if rem > 0 {
				m := rem / 60
				s := rem % 60
				fmt.Printf("  %s(\\__/)%s\n", helpers.Cyan, helpers.Reset)
				fmt.Printf("  %s(•ㅅ•) / づ%s %s[IN FOCUS 🎯]%s %s%s%s\n",
					helpers.Cyan, helpers.Reset,
					helpers.Bold+helpers.Green, helpers.Reset,
					helpers.Bold+helpers.White, data.Timer.Task, helpers.Reset)
				fmt.Println()
				fmt.Printf("  %s⏳ %02d:%02d remaining%s %s(%dm total)%s\n",
					helpers.Bold+helpers.Yellow, m, s, helpers.Reset,
					helpers.Gray, data.Timer.Minutes, helpers.Reset)
			} else {
				fmt.Printf("  %s(\\__/)%s\n", helpers.Cyan, helpers.Reset)
				fmt.Printf("  %s(•ㅅ•) / づ%s %s[TIME'S UP 🔔]%s %s%s%s\n",
					helpers.Cyan, helpers.Reset,
					helpers.Bold+helpers.Red, helpers.Reset,
					helpers.Bold+helpers.White, data.Timer.Task, helpers.Reset)
				fmt.Println()
				helpers.Success("  Session completed! Take a break, Architect.")
			}
			fmt.Println()
			return nil
		}

		// Stop / clear timer
		if args[0] == "stop" || args[0] == "clear" || args[0] == "del" {
			return store.Update(func(data *store.CrumbData) error {
				data.Timer = nil
				fmt.Println()
				fmt.Printf("  %s(\\__/)%s\n", helpers.Cyan, helpers.Reset)
				fmt.Printf("  %s(•ㅅ•) / づ%s %s[TIMER STOPPED]%s\n", helpers.Cyan, helpers.Reset, helpers.Yellow, helpers.Reset)
				fmt.Println()
				return nil
			})
		}

		// Parse minutes
		mins, err := strconv.Atoi(args[0])
		if err != nil || mins < 1 {
			helpers.Error("Invalid duration: %s. Example: crumb timer 50 \"Doing DSA\"", args[0])
			return nil
		}

		taskName := "Focus Sprint"
		if len(args) > 1 {
			taskName = strings.Join(args[1:], " ")
		}

		return store.Update(func(data *store.CrumbData) error {
			data.Timer = &store.TimerState{
				Task:      taskName,
				Minutes:   mins,
				StartedAt: time.Now().Unix(),
				Duration:  mins * 60,
			}

			fmt.Println()
			fmt.Printf("  %s(\\__/)%s\n", helpers.Cyan, helpers.Reset)
			fmt.Printf("  %s(•ㅅ•) / づ%s %s[TIMER STARTED ⏱️]%s\n", helpers.Cyan, helpers.Reset, helpers.Bold+helpers.Green, helpers.Reset)
			fmt.Println()
			fmt.Printf("  %sTask:%s     %s%s%s\n", helpers.Cyan, helpers.Reset, helpers.Bold+helpers.White, taskName, helpers.Reset)
			fmt.Printf("  %sDuration:%s %s%dm%s\n", helpers.Cyan, helpers.Reset, helpers.Green, mins, helpers.Reset)
			fmt.Println()
			return nil
		})
	},
}

func init() {
	RootCmd.AddCommand(timerCmd)
}

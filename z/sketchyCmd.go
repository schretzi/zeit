package z

import (
	"fmt"
	"log"
	"os/exec"
	"time"

	"github.com/spf13/cobra"
)

var sketchyCmd = &cobra.Command{
	Use:   "sketchy",
	Short: "Currently tracking in sketchybar",
	Long:  "Show currently tracking activity in sketchybar.",
	Run: func(cmd *cobra.Command, args []string) {
		user := GetCurrentUser()
		var shellCmd string
		runningEntryId, err := database.GetRunningEntryId(user)
		if err != nil {
			log.Fatalf(ErrorString, CharError, err)
		}

		if runningEntryId == "" {
			shellCmd = "label=No active Tracking"
		} else {
			runningEntry, err := database.GetEntry(user, runningEntryId)
			if err != nil {
				log.Fatalf(ErrorString, CharError, err)
			}

			dur := time.Since(runningEntry.Begin)
			shellCmd = "label=" + runningEntry.Project + " / " + runningEntry.Task + ":" + fmtDuration(dur)
		}
		runCmd := exec.Command("/opt/homebrew/bin/sketchybar", "-m", "--set", "zeit", shellCmd)
		_, err = runCmd.Output()

		if err != nil {
			fmt.Println(err.Error())
			return
		}
	},
}

func init() {
	rootCmd.AddCommand(sketchyCmd)
}

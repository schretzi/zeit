package z

import (
	"fmt"
	"log"

	"github.com/spf13/cobra"
)

var trackingCmd = &cobra.Command{
	Use:   "tracking",
	Short: "Currently tracking activity",
	Long:  "Show currently tracking activity.",
	Run: func(cmd *cobra.Command, args []string) {
		user := GetCurrentUser()

		runningEntryId, err := database.GetRunningEntryId(user)
		if err != nil {
			log.Fatalf(ErrorString, CharError, err)
		}

		if runningEntryId == "" {
			log.Fatalf("%s not running\n", CharFinish)
		}

		runningEntry, err := database.GetEntry(user, runningEntryId)
		if err != nil {
			log.Fatalf(ErrorString, CharError, err)
		}

		fmt.Print(runningEntry.GetOutputForTrack(true, true))
	},
}

func init() {
	rootCmd.AddCommand(trackingCmd)
}

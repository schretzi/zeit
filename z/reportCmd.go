package z

import (
	"fmt"
	"time"

	"github.com/shopspring/decimal"
	"github.com/spf13/cobra"
)
type taskPerDay struct {
	name string
	dur time.Duration
	notes []string
}

type project per Day {
	name     
}


}
var reportCmd = &cobra.Command{
	Use:   "report",
	Short: "report times an day / project / task level",
	Long:  "Reporting summaries on daily, project, task level for a given range",
	Run: func(cmd *cobra.Command, args []string) {
		var projectPerDay = make(map[string]map[string]float32)
		var taskPerDay = make
		filteredEntries := listEntries()
		

		totalHours := decimal.NewFromInt(0)
		for _, entry := range filteredEntries {
			totalHours = totalHours.Add(entry.GetDuration())
			fmt.Printf("%s\n", entry.GetOutput(false))
		}

		if listTotalTime {
			fmt.Printf("\nTOTAL: %s H\n\n", fmtHours(totalHours))
		}
	},
}

func init() {
	rootCmd.AddCommand(reportCmd)

	reportCmd.Flags().StringVar(&since, "since", "", "Date/time to start the list from")
	reportCmd.Flags().StringVar(&until, "until", "", "Date/time to list until")
	reportCmd.Flags().StringVar(&listRange, "range", "", "shortcut to set since/until for a given range (today, yesterday, thisWeek, lastWeek, thisMonth, lastMonth)")
	reportCmd.Flags().StringVarP(&project, "project", "p", "", "Project to be listed")
	reportCmd.Flags().StringVarP(&task, "task", "t", "", "Task to be listed")

	flagName := "task"
	reportCmd.RegisterFlagCompletionFunc(flagName, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		user := GetCurrentUser()
		entries, _ := database.ListEntries(user)
		_, tasks := listProjectsAndTasks(entries)
		return tasks, cobra.ShellCompDirectiveDefault
	})
}

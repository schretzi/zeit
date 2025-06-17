package z

import (
	"fmt"
	"time"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	listTotalTime            bool
	listOnlyProjectsAndTasks bool
	listOnlyTasks            bool
	appendProjectIDToTask    bool
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List activities",
	Long:  "List all tracked activities.",
	Run: func(cmd *cobra.Command, args []string) {
		filteredEntries := listEntries()

		var totalHours time.Duration = 0
		for _, entry := range filteredEntries {
			totalHours += entry.GetDuration()
			if showNotesFlag {
				fmt.Printf("%s\n", entry.GetOutput(false, true))
			} else {
				fmt.Printf("%s\n", entry.GetOutput(false, false))
		  }
		}

		if listTotalTime == true {
			fmt.Printf("\nTOTAL: %s H\n\n", fmtDuration(totalHours))
		}
		return
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
	listCmd.Flags().StringVar(&since, "since", "", "Date/time to start the list from")
	listCmd.Flags().StringVar(&until, "until", "", "Date/time to list until")
	listCmd.Flags().StringVar(&listRange, "range", "", "Shortcut for --since and --until that accepts: "+strings.Join(Ranges(), ", "))
	listCmd.Flags().StringVarP(&project, "project", "p", "", "Project to be listed")
	listCmd.Flags().StringVarP(&task, "task", "t", "", "Task to be listed")
	listCmd.Flags().BoolVar(&fractional, "decimal", false, "Show fractional hours in decimal format instead of minutes")
	listCmd.Flags().BoolVar(&listTotalTime, "total", false, "Show total time of hours for listed activities")
	listCmd.Flags().BoolVar(&listOnlyProjectsAndTasks, "only-projects-and-tasks", false, "Only list projects and their tasks, no entries")
	listCmd.Flags().BoolVar(&listOnlyTasks, "only-tasks", false, "Only list tasks, no projects nor entries")
	listCmd.Flags().BoolVar(&appendProjectIDToTask, "append-project-id-to-task", false, "Append project ID to tasks in the list")
	listCmd.Flags().BoolVar(&showNotesFlag, "notes", false, "Show notes from task")
	viper.BindPFlag("list.notes", listCmd.Flags().Lookup("showNotesFlag"))

	flagName := "task"
	listCmd.RegisterFlagCompletionFunc(flagName, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		user := GetCurrentUser()
		entries, _ := database.ListEntries(user)
		_, tasks := listProjectsAndTasks(entries)
		return tasks, cobra.ShellCompDirectiveDefault
	})

	flagName = "range"
	listCmd.RegisterFlagCompletionFunc(flagName, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		ranges := Ranges()
		return ranges, cobra.ShellCompDirectiveDefault
	})
}

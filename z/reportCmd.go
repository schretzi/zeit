package z

import (
  "fmt"
  "sort"
  "time"

  "github.com/go-gota/gota/dataframe"
  "github.com/go-gota/gota/series"
  "github.com/spf13/cobra"
  "github.com/spf13/viper"
  "golang.org/x/exp/maps"
)

type reportEntry struct {
  Day      string
  Project  string
  Task     string
  Duration float64
  Notes    string
}

var reportCmd = &cobra.Command{
  Use:   "report",
  Short: "report times an day / project / task level",
  Long:  "Reporting summaries on daily, project, task level for a given range",
  Run: func(cmd *cobra.Command, args []string) {

    // Set Aggregation function on dataframe
    var aggregate []dataframe.AggregationType
    aggregate = append(aggregate, dataframe.Aggregation_SUM, dataframe.Aggregation_MIN)

    if since == "" && until == "" && listRange == "" {
      // For me report without any time limit makes no sense, so I use default from Config if set
      listRange = viper.GetString("report.default")
    }

    // Filter entries and load into dataframe for further analysis
    filteredEntries := listEntries()
    if listRange != "" {
      fmt.Println("Reporting for Timerange:", listRange, "/", sinceTime.Format(DateFormat), "-", untilTime.Format(DateFormat))
    }
    var reportEntries []reportEntry
    for _, fe := range filteredEntries {
      var entryDuration float64
      if fe.Finish.IsZero() {
        entryDuration = -1
      } else {
        entryDuration = time.Duration(fe.Finish.Sub(fe.Begin)).Seconds()
      }
      reportEntries = append(reportEntries, reportEntry{fe.Begin.Format(DateFormat), fe.Project, fe.Task, entryDuration, fe.Notes})
    }
    df := dataframe.LoadStructs(reportEntries)

    // Group and Order by Day
    groupedDay := df.GroupBy("Day")

    daySum := groupedDay.Aggregation(aggregate, []string{ColDuration, ColDuration})

    dayKeys := maps.Keys(groupedDay.GetGroups())
    sort.Strings(dayKeys)

    for _, dayKey := range dayKeys {
      entriesForDay := daySum.Filter(dataframe.F{Colname: "Day", Comparator: series.Eq, Comparando: dayKey})

      durFloat := entriesForDay.Col("Duration_SUM").Float()
      durMin := entriesForDay.Col("Duration_MIN").Float()
      durDay := time.Duration(durFloat[0] * float64(time.Second))
      durDay = durDay.Round(time.Minute)
      if durMin[0] < 0 {
        durDay = durDay + time.Second
        fmt.Println("\n", dayKey, ": ", fmtDuration(durDay), RunningFlag)
      } else {
        fmt.Println("\n", dayKey, ": ", fmtDuration(durDay))
      }
      // Group and Sum on Project for this Day
      filteredProject := df.Filter(dataframe.F{Colname: "Day", Comparator: series.Eq, Comparando: dayKey})
      groupedProject := filteredProject.GroupBy("Project")
      projectSum := groupedProject.Aggregation(aggregate, []string{ColDuration, ColDuration})
      projectKeys := maps.Keys(groupedProject.GetGroups())
      sort.Strings(projectKeys)
      for _, projectKey := range projectKeys {
        entriesForProject := projectSum.Filter(dataframe.F{Colname: "Project", Comparator: series.Eq, Comparando: projectKey})
        projectDurFloat := entriesForProject.Col("Duration_SUM").Float()
        projectDurMin := entriesForProject.Col("Duration_MIN").Float()
        projectDur := time.Duration(projectDurFloat[0] * float64(time.Second))
        projectDur = projectDur.Round(time.Minute)
        if projectDurMin[0] < 0 {
          projectDur = projectDur + time.Second
          fmt.Println("    ", projectKey, ": ", fmtDuration(projectDur), RunningFlag)
        } else {
          fmt.Println("    ", projectKey, ": ", fmtDuration(projectDur))
        }

        // Group and Sum on Tasks for that Project on this Day
        filteredTask := filteredProject.Filter(dataframe.F{Colname: "Project", Comparator: series.Eq, Comparando: projectKey})
        groupedTask := filteredTask.GroupBy("Task")
        taskSum := groupedTask.Aggregation(aggregate, []string{ColDuration, ColDuration})
        taskKeys := maps.Keys(groupedTask.GetGroups())
        sort.Strings(taskKeys)
        for _, taskKey := range taskKeys {
          entriesForTask := taskSum.Filter(dataframe.F{Colname: "Task", Comparator: series.Eq, Comparando: taskKey})
          taskDurFloat := entriesForTask.Col("Duration_SUM").Float()
          taskDurMin := entriesForTask.Col("Duration_MIN").Float()
          taskDur := time.Duration(taskDurFloat[0] * float64(time.Second))
          taskDur = taskDur.Round(time.Minute)
          if taskDurMin[0] < 0 {
            taskDur = taskDur + time.Second
            fmt.Println("        ", taskKey, ": ", fmtDuration(taskDur), RunningFlag)
          } else {
            fmt.Println("        ", taskKey, ": ", fmtDuration(taskDur))
          }
        }
      }
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

package z

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/jinzhu/now"
	"github.com/spf13/cobra"
)

func exportZeitJson(entries []Entry) (string, error) {
	stringified, err := json.Marshal(entries)
	if err != nil {
		return "", err
	}

	return string(stringified), nil
}

func exportTymeJson(entries []Entry) (string, error) {
	tyme := Tyme{}
	err := tyme.FromEntries(entries)
	if err != nil {
		return "", err
	}

	return tyme.Stringify(), nil
}

var exportCmd = &cobra.Command{
	Use:   "export ([flags])",
	Short: "Export tracked activities",
	Long:  "Export tracked activities to various formats.",
	// Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		var entries []Entry
		var err error

		user := GetCurrentUser()

		entries, err = database.ListEntries(user)
		if err != nil {
			log.Fatalf(ErrorString, CharError, err)
		}

		var sinceTime time.Time
		var untilTime time.Time

		if since != "" {
			sinceTime, err = now.Parse(since)
			if err != nil {
				log.Fatalf(ErrorString, CharError, err)
			}
		}

		if until != "" {
			untilTime, err = now.Parse(until)
			if err != nil {
				log.Fatalf(ErrorString, CharError, err)
			}
		}

		var filteredEntries []Entry
		filteredEntries, err = GetFilteredEntries(entries, project, task, sinceTime, untilTime)
		if err != nil {
			log.Fatalf(ErrorString, CharError, err)
		}

		var output string = ""
		switch format {
		case "zeit":
			output, err = exportZeitJson(filteredEntries)
			if err != nil {
				log.Fatalf(ErrorString, CharError, err)
			}
		case "tyme":
			output, err = exportTymeJson(filteredEntries)
			if err != nil {
				log.Fatalf(ErrorString, CharError, err)
			}
		default:
			log.Fatalf("%s specify an export format; see `zeit export --help` for more info\n", CharError)
		}

		fmt.Printf("%s\n", output)
	},
}

func init() {
	rootCmd.AddCommand(exportCmd)

	exportCmd.Flags().StringVar(&format, "format", "zeit", "Format to export, possible values: zeit, tyme")
	exportCmd.Flags().StringVar(&since, "since", "", "Date/time to start the export from")
	exportCmd.Flags().StringVar(&until, "until", "", "Date/time to export until")
	exportCmd.Flags().StringVarP(&project, "project", "p", "", "Project to be exported")
	exportCmd.Flags().StringVarP(&task, "task", "t", "", "Task to be exported")
}

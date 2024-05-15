package z

import (
	"fmt"
	"log"

	"github.com/gookit/color"
	"github.com/spf13/cobra"
)

var eraseCmd = &cobra.Command{
	Use:   "erase ([flags]) [id]",
	Short: "Erase activity",
	Long:  "Erase tracked activity.",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		user := GetCurrentUser()
		id := args[0]

		err := database.EraseEntry(user, id)
		if err != nil {
			log.Fatalf(ErrorString, CharError, err)
		}

		fmt.Printf("%s erased %s\n", CharInfo, color.FgLightWhite.Render(id))
	},
}

func init() {
	rootCmd.AddCommand(eraseCmd)
}

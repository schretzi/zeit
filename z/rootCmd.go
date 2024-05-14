package z

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"os"

	"github.com/gookit/color"
)

var database *Database

var begin string
var finish string
var project string
var task string
var notes string

var since string
var until string

var format string
var force bool

var noColors bool
var Debug bool
var cfgFile string

const (
	CharTrack  = " ▶"
	CharFinish = " ■"
	CharErase  = " ◀"
	CharError  = " ▲"
	CharInfo   = " ●"
	CharMore   = " ◆"
)

var rootCmd = &cobra.Command{
	Use:   "zeit",
	Short: "Command line Zeiterfassung",
	Long:  `A command line time tracker.`,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Printf("%s %+v\n", CharError, err)
		os.Exit(-1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().BoolVar(&noColors, "no-colors", false, "Do not use colors in output")
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $XDG_CONFIG_HOME/zeit/zeit.yaml)")
	rootCmd.PersistentFlags().BoolVarP(&Debug, "debug", "d", false, "Display debugging output in the console. (default: false)")
	viper.BindPFlag("debug", rootCmd.PersistentFlags().Lookup("debug"))

}

func initConfig() {

	if noColors {
		color.Disable()
	}
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		viper.AddConfigPath("$XDG_CONFIG_HOME/zeit")
		viper.AddConfigPath(home + "/.config/zeit")
		viper.SetConfigType("yaml")
		viper.SetConfigName("zeit")
	}

	viper.SetEnvPrefix("zeit")
	viper.BindEnv("db")

	if err := viper.ReadInConfig(); err != nil {
		fmt.Printf("Error in Reading Configuration: %v", err)
		os.Exit(1)
	}

	if viper.GetBool("debug") {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
		fmt.Fprintln(os.Stderr, "Using Database file:", viper.GetString("db"))
	}

	var err error
	database, err = InitDatabase()
	if err != nil {
		fmt.Printf("%s %+v\n", CharError, err)
		os.Exit(1)
	}
}

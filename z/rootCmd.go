package z

import (
	"fmt"
	"log"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"os"

	"github.com/gookit/color"
)

var database *Database

var begin string
var finish string
var switchString string
var project string
var task string
var notes string

var since string
var until string
var listRange string
var sinceTime time.Time
var untilTime time.Time

var format string
var force bool

var noColors bool
var debug bool
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
		log.Fatalf(ErrorString, CharError, err)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $XDG_CONFIG_HOME/zeit/zeit.yaml)")

	rootCmd.PersistentFlags().BoolVar(&noColors, FlagNoColors, false, "Do not use colors in output")
	viper.BindPFlag(FlagNoColors, rootCmd.PersistentFlags().Lookup(FlagNoColors))

	rootCmd.PersistentFlags().BoolVarP(&debug, FlagDebug, "d", false, "Display debugging output in the console. (default: false)")
	viper.BindPFlag(FlagDebug, rootCmd.PersistentFlags().Lookup(FlagDebug))
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
		// Set default values for parameters
		viper.Set("debug", false)
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

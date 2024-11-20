/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"

	"datadog_import/common"
	"datadog_import/logctx"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	cfgFile, logLevel string
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "datadog_import",
	Short: "Converts Dynatrace synthetic monitors to Datadog synthetic monitors",
	Long:  `datadog_import is a tool that converts Dynatrace synthetic monitors to Datadog synthetic monitors.`,
	Run:   run,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	rootCmd.PersistentFlags().StringVar(&logLevel, "log", "info", "log level (default is info))")
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is ./config.yaml)")

	viper.BindPFlag("log", rootCmd.PersistentFlags().Lookup("log"))
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		// Search config in home directory with name "config" (without extension).
		viper.AddConfigPath(home)
		viper.AddConfigPath(".")
		viper.SetConfigType("yaml")
		viper.SetConfigName("config")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	}
}

func run(cmd *cobra.Command, args []string) {
	// logrus.SetFormatter(&logrus.JSONFormatter{})
	logrus.SetReportCaller(true)
	logctx.Default = logrus.NewEntry(logrus.StandardLogger())
	var config ViperConfig
	err := viper.Unmarshal(&config)
	common.Check(cmd.Context(), err)
	config.SetLogLevel()
	ctx := logctx.New(cmd.Context(), logctx.Default)
	conv := config.BuildConverter(ctx)
	conv.Convert()
}

/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"

	"dynatrace_to_datadog/common"
	"dynatrace_to_datadog/converter"
	"dynatrace_to_datadog/datadog/api"
	"dynatrace_to_datadog/dynatrace"
	"dynatrace_to_datadog/logctx"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	cfgFile, logLevel string
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "dynatrace_to_datadog",
	Short: "Converts Dynatrace synthetic monitors to Datadog synthetic monitors",
	Long:  `dynatrace_to_datadog is a tool that converts Dynatrace synthetic monitors to Datadog synthetic monitors.`,
	Run:   makeConverter,
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

func makeConverter(cmd *cobra.Command, args []string) {
	level, err := logrus.ParseLevel(viper.GetString("log"))
	if err != nil {
		logctx.Default.Panic(err)
	}
	logrus.SetLevel(level)
	// logrus.SetFormatter(&logrus.JSONFormatter{})
	// logrus.SetReportCaller(true)
	logctx.Default = logrus.NewEntry(logrus.StandardLogger())
	ctx := logctx.New(cmd.Context(), logctx.Default)
	conv := converter.Converter{}
	dynatraceConfig := viper.Sub("dynatrace")
	if dynatraceConfig != nil {
		var dynaConf dynatrace.Config
		err := dynatraceConfig.Unmarshal(&dynaConf)
		common.Check(ctx, err)
		conv.Reader, err = dynaConf.GetReader()
		common.Check(ctx, err)
		conv.Transform = dynaConf.GetTransformer()
	} else {
		panic(fmt.Errorf("dynatrace config is nil"))
	}
	ddConfig := viper.Sub("datadog")
	if ddConfig != nil {
		var datadogConf api.Config
		err := ddConfig.Unmarshal(&datadogConf)
		common.Check(ctx, err)
		conv.Writers = append(conv.Writers, datadogConf.NewDatadogWriter())
	}

	outputPath := viper.GetString("output")
	if outputPath != "" {
		writer, err := common.NewFileWriter(outputPath)
		common.Check(ctx, err)
		conv.Writers = append(conv.Writers, writer)
	}

	if len(conv.Writers) == 0 {
		panic(fmt.Errorf("no output found"))
	}
	conv.Convert(ctx)
}

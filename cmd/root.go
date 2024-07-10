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
	Short: "A brief description of your application",
	Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	Run: makeConverter,
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

func makeConfig() *common.Config {
	log := logrus.New()
	level, err := logrus.ParseLevel(viper.GetString("log"))
	common.Check(err)
	log.SetLevel(level)
	return &common.Config{
		Log: log,
	}
}

func makeConverter(cmd *cobra.Command, args []string) {
	conf := makeConfig()
	conv := converter.Converter{Config: conf}
	dynatraceConfig := viper.Sub("dynatrace")
	if dynatraceConfig != nil {
		var dynaConf dynatrace.Config
		err := dynatraceConfig.Unmarshal(&dynaConf)
		common.Check(err)
		conv.Reader, err = dynaConf.GetReader()
		common.Check(err)
		conv.Transform = dynaConf.GetTransformer()
	} else {
		panic(fmt.Errorf("dynatrace config is nil"))
	}
	ddConfig := viper.Sub("datadog")
	if ddConfig != nil {
		var datadogConf api.Config
		err := ddConfig.Unmarshal(&datadogConf)
		common.Check(err)
		conv.Writers = append(conv.Writers, datadogConf.NewSyntheticsBrowserWriter())
	}

	outputPath := viper.GetString("output")
	if outputPath != "" {
		writer, err := common.NewFileWriter(outputPath)
		common.Check(err)
		conv.Writers = append(conv.Writers, writer)
	}

	if len(conv.Writers) == 0 {
		panic(fmt.Errorf("no output found"))
	}
	conv.Convert()
}

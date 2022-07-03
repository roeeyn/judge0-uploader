/*
Copyright © 2022 Rodrigo Medina rodrigo.medina.neri@gmail.com

*/
package cmd

import (
	"fmt"
	"os"

	logger "github.com/roeeyn/judge0-uploader/pkg/logger"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string
var verbose bool

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "judge0-uploader",
	Short: "Utility to upload local files to Judge0",
	Long: `This CLI utility is used to upload local files to Judge0.

This is useful to load local coding challenges to Judge0, execute them and get the results.`,
	// Uncomment th following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
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

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.judge0-uploader.yaml)")
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "verbose output")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	// rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
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

		// Search config in home directory with name ".judge0-uploader" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigType("yaml")
		viper.SetConfigName(".judge0-uploader")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	err := viper.ReadInConfig()
	configFileUsed := viper.ConfigFileUsed()
	isTokenSet := viper.IsSet("judge0_auth_token")

	logger.IsVerbose = verbose

	if err == nil || isTokenSet {
		logger.LogInfo(fmt.Sprintf("Using config file: %s", configFileUsed))
		logger.LogInfo(fmt.Sprintf("Is Token Set?: %t", isTokenSet))
	} else {
		logger.LogError(fmt.Errorf("Searched for token in: %s", configFileUsed))
		logger.LogError(fmt.Errorf("No valid config file found in: %s", viper.ConfigFileUsed()))
		logger.LogFatal(fmt.Errorf("A Judge0 Auth Token is needed to use this utility."))
		os.Exit(1)
	}
}

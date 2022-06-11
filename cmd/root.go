/*
Copyright © 2022 Rodrigo Medina rodrigo.medina.neri@gmail.com

*/
package cmd

import (
	"fmt"
	"os"

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
	// Uncomment the following line if your bare application
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
	if err := viper.ReadInConfig(); err == nil {
		configFileUsed := viper.ConfigFileUsed()
		isTokenSet := viper.IsSet("judge0_auth_token")

		fmt.Println("Using config file:", configFileUsed)
		fmt.Println("Is Token Set?", isTokenSet)

		if !isTokenSet {
			fmt.Fprintln(os.Stderr, "A Judge0 Auth Token is needed to use this utility.")
			fmt.Fprintln(os.Stderr, "Searched for token in:", configFileUsed)
			os.Exit(1)
		}

	}
}

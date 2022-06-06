/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// submitCmd represents the submit command
var submitCmd = &cobra.Command{
	Use:   "submit",
	Short: "Submit local files from specified directory to Judge0",
	Long: `Submit local coding challenge files from specified directory to the Judge0 server.

We're expecting that the directory contains the following files:
- run (bash script)
- main*
- test*
- testframework*

* This files should have the same extension.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("submit called")
	},
}

func init() {
	rootCmd.AddCommand(submitCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// submitCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// submitCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

/*
Copyright © 2022 Rodrigo Medina rodrigo.medina.neri@gmail.com

*/
package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/viper"

	submitter "github.com/roeeyn/judge0-uploader/pkg/j0_submitter"
	"github.com/spf13/cobra"
)

// submitCmd represents the submit command
var submitCmd = &cobra.Command{
	Use:   "submit YOUR_CHALLENGE_FOLDER",
	Short: "Submit local files from specified directory to Judge0",
	Long: `Submit local coding challenge files from specified directory to the Judge0 server.

We're expecting that the directory contains the following files:
- run (bash script)
- index*
- test*
- testframework*

* This files should have the same extension.`,
	Run:  run,
	Args: cobra.ExactArgs(1),
}

func run(cmd *cobra.Command, args []string) {
	fmt.Println("submit command called")
	fmt.Println("Challenge Relative Path:", args[0])

	j0AuthToken := viper.GetString("judge0_auth_token")
	isVerbose := viper.GetBool("verbose")
	j0Submitter := submitter.NewJ0Submitter(args[0], j0AuthToken, isVerbose)

	submissionId, err := j0Submitter.Run()

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Println("Result Submission ID:", submissionId)
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

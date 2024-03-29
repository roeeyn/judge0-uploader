/*
Copyright © 2022 Rodrigo Medina rodrigo.medina.neri@gmail.com

*/
package cmd

import (
	"fmt"

	"github.com/spf13/viper"

	logger "github.com/roeeyn/judge0-uploader/pkg/logger"
	submitter "github.com/roeeyn/judge0-uploader/pkg/submitter"
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
	Run:  runSubmit,
	Args: cobra.ExactArgs(1),
}

func runSubmit(cmd *cobra.Command, args []string) {
	j0AuthToken := viper.GetString("judge0_auth_token")
	j0ServerUrl := viper.GetString("judge0_server_url")
	submitter := submitter.NewSubmitter(args[0], j0AuthToken, j0ServerUrl)

	logger.LogInfo("Submit command called")
	logger.LogInfo(fmt.Sprintf("Challenge Relative Path: %s", args[0]))

	submissionId, err := submitter.Run()
	if err != nil {
		logger.LogFatal(err)
	}

	logger.LogInfo(fmt.Sprintf("Result Submission ID: %s", submissionId))
	fmt.Print(submissionId)
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

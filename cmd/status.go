/*
Copyright © 2022 Rodrigo Medina rodrigo.medina.neri@gmail.com

*/
package cmd

import (
	"fmt"

	logger "github.com/roeeyn/judge0-uploader/pkg/logger"
	statusFetcher "github.com/roeeyn/judge0-uploader/pkg/status_fetcher"
	utils "github.com/roeeyn/judge0-uploader/pkg/utils"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// statusCmd represents the status command
var statusCmd = &cobra.Command{
	Use:   "status YOUR_SUBMISSION_ID",
	Short: "Get the status of a challenge submission",
	Long: `Get the status of a challenge submission by its submission ID.

Based on the submission ID of the submission we can fetch the status from it,
and you can choose to wait until the execution has finished, or just get the current
status as it is.

A execution is considered finished when the status is neither "In Queue" or "Processing".`,
	Run:  runStatus,
	Args: cobra.ExactArgs(1),
}

func runStatus(cmd *cobra.Command, args []string) {
	j0AuthToken := viper.GetString("judge0_auth_token")
	j0ServerUrl := viper.GetString("judge0_server_url")
	submissionId := args[0]
	statusFetcher := statusFetcher.NewStatusFetcher(j0AuthToken, j0ServerUrl, submissionId)

	logger.LogInfo("Status command called")
	logger.LogInfo(fmt.Sprintf("Requesting status for Submission ID: %s", submissionId))

	submissionStatusResponse, err := statusFetcher.Run()
	if err != nil {
		logger.LogFatal(err)
	}

	logger.LogInfo(fmt.Sprintf("Status ID: %d", submissionStatusResponse.Status.Id))
	logger.LogInfo(fmt.Sprintf("Status description: %s", submissionStatusResponse.Status.Description))

	fmt.Println(utils.PrettyPrint(submissionStatusResponse))
}

func init() {
	rootCmd.AddCommand(statusCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// statusCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// statusCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

/*
Copyright © 2022 Rodrigo Medina rodrigo.medina.neri@gmail.com

*/
package cmd

import (
	"fmt"

	logger "github.com/roeeyn/judge0-uploader/pkg/logger"
	statusFetcher "github.com/roeeyn/judge0-uploader/pkg/status_fetcher"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// statusCmd represents the status command
var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Get the status of a challenge submission",
	Long: `Get the status of a challenge submission by its submission ID.

Based on the submission ID of the submission we can fetch the status from it,
and you can choose to wait until the execution has finished, or just get the current
status as it is.

A execution is considered finished when the status is neither "In Queue" or "Processing".`,
	Run: runStatus,
}

func runStatus(cmd *cobra.Command, args []string) {
	j0AuthToken := viper.GetString("judge0_auth_token")
	j0ServerUrl := viper.GetString("judge0_server_url")
	statusFetcher := statusFetcher.NewStatusFetcher(j0AuthToken, j0ServerUrl)

	logger.LogInfo("Status command called")
	logger.LogInfo("Requesting status for Submission Id: PENDING")

	submissionStatus, err := statusFetcher.Run()
	if err != nil {
		logger.LogFatal(err)
	}

	logger.LogInfo(fmt.Sprintf("Submission status: %s", submissionStatus))
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

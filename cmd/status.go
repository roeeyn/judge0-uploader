/*
Copyright © 2022 Rodrigo Medina rodrigo.medina.neri@gmail.com

*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
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
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("status called")
		fmt.Println("verbose:", verbose)
	},
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

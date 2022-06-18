/*
Copyright © 2022 Rodrigo Medina rodrigo.medina.neri@gmail.com

*/
package cmd

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/roeeyn/judge0-uploader/pkg/j0_submitter"

	"github.com/spf13/cobra"
)

var expectedChallengeFiles = [4]string{"index", "run", "test", "testframework"}

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

func isExpectedFile(fullFileName string) (isExpected bool, fileName string) {
	// Remove extension from the file
	fileName = strings.TrimSuffix(fullFileName, filepath.Ext(fullFileName))

	for _, expectedFile := range expectedChallengeFiles {
		if fileName == expectedFile {
			return true, expectedFile
		}
	}

	return false, ""
}

func getChallengeFiles(absBasePath string) (j0SubmitterFiles *j0_submitter.J0SubmitterFiles, err error) {
	// Verify that the basePath contains expected files
	files, err := ioutil.ReadDir(absBasePath)
	if err != nil {
		err = fmt.Errorf(fmt.Sprintf("Error reading files inside folder: %s", err.Error()))
		return
	}

	submitterFiles := j0_submitter.NewJ0SubmitterFiles()

	for _, file := range files {
		// We do not support nested folders
		if file.IsDir() {
			continue
		}

		absFilePath := path.Join(absBasePath, file.Name())
		InfoLogger.Println("Found file: ", absFilePath)

		if isExpected, fileNameKey := isExpectedFile(absFilePath); isExpected {
			error := submitterFiles.AddFile(fileNameKey, absFilePath)
			if error != nil {
				err = fmt.Errorf(fmt.Sprintf("Error adding file property: %s", error.Error()))
				return
			}
		}
	}

	if submitterFiles.ContainsEmptyFileProperties() {
		err = fmt.Errorf(fmt.Sprintf("Not all needed files are present. Expected files are: %s", expectedChallengeFiles))
		ErrorLogger.Println("Current Files:", submitterFiles)
		return
	}

	InfoLogger.Println("Challenge Files:", submitterFiles)

	return
}

func getAbsolutePath(basePath string) (absPath string, err error) {
	// Validate if the basePath exists
	_, err = os.Stat(basePath)
	if err != nil {
		if os.IsNotExist(err) {
			err = fmt.Errorf(fmt.Sprintf("Base folder: '%s' does not exist", basePath))
		}
		return
	}

	// Get the absolute path
	absPath, err = filepath.Abs(basePath)
	if err != nil {
		err = fmt.Errorf(fmt.Sprintf("Error getting absolute basePath: %s", err.Error()))
		return
	}

	InfoLogger.Println("Challenge Absolute Path:", absPath)
	return
}

func run(cmd *cobra.Command, args []string) {
	InfoLogger.Println("submit command called")
	InfoLogger.Println("Challenge Relative Path:", args[0])

	challengePath := args[0]
	absPath, err := getAbsolutePath(challengePath)
	if err != nil {
		ErrorLogger.Fatal(err)
	}

	challengeFiles, err := getChallengeFiles(absPath)
	if err != nil {
		ErrorLogger.Fatal(err)
	}

	InfoLogger.Println("Challenge Files [REMOVE THIS]:", challengeFiles)

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

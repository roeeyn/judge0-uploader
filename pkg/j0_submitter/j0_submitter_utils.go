package j0_submitter

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

var expectedChallengeFiles = [4]string{"index", "run", "test", "testframework"}

func IsExpectedFile(fullFileName string) (isExpected bool, fileName string) {
	// Remove extension from the file
	fileName = strings.TrimSuffix(fullFileName, filepath.Ext(fullFileName))

	for _, expectedFile := range expectedChallengeFiles {
		if fileName == expectedFile {
			return true, expectedFile
		}
	}

	return false, ""
}

func GetAbsolutePath(basePath string) (absPath string, err error) {
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

	return
}

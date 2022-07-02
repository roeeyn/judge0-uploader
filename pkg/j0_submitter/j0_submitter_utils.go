package j0_submitter

import (
	"bufio"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

var ExpectedChallengeFiles = [4]string{"index", "run", "test", "testframework"}

func IsExpectedFile(fullFileName string) (isExpected bool, fileName string) {
	// Remove extension from the file
	fileName = strings.TrimSuffix(fullFileName, filepath.Ext(fullFileName))

	for _, expectedFile := range ExpectedChallengeFiles {
		if fileName == expectedFile {
			return true, expectedFile
		}
	}

	return false, ""
}

func B64ZipFile() (encodedFile string, err error) {
	// Open file on disk.
	f, err := os.Open(ZIP_FILE_NAME)
	if err != nil {
		err = fmt.Errorf("Error opening zip file %s: %s", ZIP_FILE_NAME, err.Error())
		return
	}

	// Read entire JPG into byte slice.
	reader := bufio.NewReader(f)
	content, err := ioutil.ReadAll(reader)
	if err != nil {
		err = fmt.Errorf("Error reading zip file %s: %s", ZIP_FILE_NAME, err.Error())
		return
	}

	// Encode as base64.
	encodedFile = base64.StdEncoding.EncodeToString(content)
	return
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

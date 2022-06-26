package j0_submitter

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
)

type J0SubmitterFiles struct {
	Run   string
	Index string
	Test  string
	// We're leaving it without camelcase to maintain consistency
	// with the expected challenge file.
	Testframework string
}

type J0Submitter struct {
	Files            *J0SubmitterFiles
	ChallengePath    string
	AbsChallengePath string
}

func (j0SubmitterFiles *J0SubmitterFiles) AddFile(fileKey string, fileAbsPath string) error {
	switch fileKey {
	case "run":
		j0SubmitterFiles.Run = fileAbsPath
	case "index":
		j0SubmitterFiles.Index = fileAbsPath
	case "test":
		j0SubmitterFiles.Test = fileAbsPath
	case "testframework":
		j0SubmitterFiles.Testframework = fileAbsPath
	default:
		return fmt.Errorf("Unknown file key: %s", fileKey)
	}
	return nil
}

func (j0SubmitterFiles *J0SubmitterFiles) ContainsEmptyFileProperties() bool {
	// Iterate over the struct fields and check if any of them is empty
	// because at this point every file is needed.
	if j0SubmitterFiles.Run == "" || j0SubmitterFiles.Index == "" || j0SubmitterFiles.Test == "" || j0SubmitterFiles.Testframework == "" {
		return true
	}

	return false
}

func (j0Submitter *J0Submitter) Run() (err error) {
	challengePath := j0Submitter.ChallengePath
	absPath, err := GetAbsolutePath(challengePath)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Println("Absolute path:", absPath)
	j0Submitter.AbsChallengePath = absPath

	err = j0Submitter.GetChallengeFiles()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	return
}

func NewJ0SubmitterFiles() (j0SubmitterFiles *J0SubmitterFiles) {
	return &J0SubmitterFiles{}
}

func NewJ0Submitter(challengePath string) (j0Submitter *J0Submitter) {
	j0Submitter = &J0Submitter{
		Files:         NewJ0SubmitterFiles(),
		ChallengePath: challengePath,
	}
	return
}

func (j0Submitter *J0Submitter) GetChallengeFiles() (err error) {
	// Verify that the basePath contains expected files
	absBasePath := j0Submitter.AbsChallengePath
	files, err := ioutil.ReadDir(absBasePath)
	if err != nil {
		err = fmt.Errorf(fmt.Sprintf("Error reading files inside folder: %s", err.Error()))
		return
	}

	for _, file := range files {
		// We do not support nested folders
		if file.IsDir() {
			continue
		}

		absFilePath := path.Join(absBasePath, file.Name())
		fmt.Println("Found file: ", absFilePath)

		if isExpected, fileNameKey := IsExpectedFile(file.Name()); isExpected {
			error := j0Submitter.Files.AddFile(fileNameKey, absFilePath)
			fmt.Println("Added file: ", absFilePath)
			if error != nil {
				err = fmt.Errorf(fmt.Sprintf("Error adding file property: %s", error.Error()))
				return
			}
		}
	}

	if j0Submitter.Files.ContainsEmptyFileProperties() {
		err = fmt.Errorf(fmt.Sprintf("Not all needed files are present. Expected files are: %s", ExpectedChallengeFiles))
		fmt.Println("Current Files:", j0Submitter.Files)
		return
	}

	fmt.Println("Challenge Files:", j0Submitter.Files)
	return
}

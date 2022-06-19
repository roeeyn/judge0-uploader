package j0_submitter

import (
	"fmt"
	"io/ioutil"
	"log"
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
	DebugLogger      *log.Logger
	ErrorLogger      *log.Logger
	InfoLogger       *log.Logger
	WarningLogger    *log.Logger
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
		j0Submitter.ErrorLogger.Fatal(err)
	}

	j0Submitter.InfoLogger.Println("Absolute path:", absPath)
	j0Submitter.AbsChallengePath = absPath

	err = j0Submitter.getChallengeFiles()
	if err != nil {
		j0Submitter.ErrorLogger.Fatal(err)
	}

	return
}

func NewJ0SubmitterFiles() (j0SubmitterFiles *J0SubmitterFiles) {
	return &J0SubmitterFiles{}
}

func NewJ0Submitter(challengePath string, debugLogger *log.Logger, errorLogger *log.Logger, infoLogger *log.Logger, warningLogger *log.Logger) (j0Submitter *J0Submitter) {
	j0Submitter = &J0Submitter{
		Files:         NewJ0SubmitterFiles(),
		ChallengePath: challengePath,
		DebugLogger:   debugLogger,
		ErrorLogger:   errorLogger,
		InfoLogger:    infoLogger,
		WarningLogger: warningLogger,
	}
	return
}

func (j0Submitter *J0Submitter) getChallengeFiles() (err error) {
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
		j0Submitter.InfoLogger.Println("Found file: ", absFilePath)

		if isExpected, fileNameKey := IsExpectedFile(file.Name()); isExpected {
			error := j0Submitter.Files.AddFile(fileNameKey, absFilePath)
			j0Submitter.InfoLogger.Println("Added file: ", absFilePath)
			if error != nil {
				err = fmt.Errorf(fmt.Sprintf("Error adding file property: %s", error.Error()))
				return
			}
		}
	}

	if j0Submitter.Files.ContainsEmptyFileProperties() {
		err = fmt.Errorf(fmt.Sprintf("Not all needed files are present. Expected files are: %s", expectedChallengeFiles))
		j0Submitter.ErrorLogger.Println("Current Files:", j0Submitter.Files)
		return
	}

	j0Submitter.InfoLogger.Println("Challenge Files:", j0Submitter.Files)

	return
}

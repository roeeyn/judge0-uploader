package j0_submitter

import (
	"fmt"
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
	Files *J0SubmitterFiles
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

func NewJ0SubmitterFiles() (j0SubmitterFiles *J0SubmitterFiles) {
	return &J0SubmitterFiles{}
}

func NewJ0Submitter(files *J0SubmitterFiles) (j0Submitter *J0Submitter) {
	return &J0Submitter{Files: files}
}

package j0_submitter

import (
	"fmt"
)

type J0SubmitterFiles struct {
	run           string
	index         string
	test          string
	testframework string
}

type J0Submitter struct {
	files J0SubmitterFiles
}

func (j0SubmitterFiles *J0SubmitterFiles) AddFile(fileKey string, fileAbsPath string) error {
	switch fileKey {
	case "run":
		j0SubmitterFiles.run = fileAbsPath
	case "index":
		j0SubmitterFiles.index = fileAbsPath
	case "test":
		j0SubmitterFiles.test = fileAbsPath
	case "testframework":
		j0SubmitterFiles.testframework = fileAbsPath
	default:
		return fmt.Errorf("Unknown file key: %s", fileKey)
	}
	return nil
}

func (j0SubmitterFiles *J0SubmitterFiles) ContainsEmptyFileProperties() bool {
	// Iterate over the struct fields and check if any of them is empty
	// because at this point every file is needed.
	if j0SubmitterFiles.run == "" || j0SubmitterFiles.index == "" || j0SubmitterFiles.test == "" || j0SubmitterFiles.testframework == "" {
		return true
	}

	return false
}

func NewJ0SubmitterFiles() (j0SubmitterFiles *J0SubmitterFiles) {
	return &J0SubmitterFiles{}
}

func NewJ0Submitter(files J0SubmitterFiles) (j0Submitter *J0Submitter) {
	return &J0Submitter{files: files}
}

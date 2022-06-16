package j0_submitter

import (
	"fmt"
)

type J0Submitter struct {
	fileNames []string
}

func NewJ0Submitter(baseDir string) (j0Submitter *J0Submitter, err error) {
	fmt.Println("Uploading:", baseDir)
	return &J0Submitter{fileNames: []string{"holi"}}, nil
}

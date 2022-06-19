package j0_submitter_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"path/filepath"

	"github.com/roeeyn/judge0-uploader/pkg/j0_submitter"
)

var _ = Describe("J0 Submitter Utils Test", func() {
	Context("When having found a file", func() {
		It("Should now if it's a valid file", func() {
			// Arrange
			files := []string{"index.py", "run", "test.js", "testframework.esotericextension"}
			for _, file := range files {
				// Act
				isExpected, _ := j0_submitter.IsExpectedFile(file)

				// Assert
				Expect(isExpected).To(BeTrue())
			}

		})

		It("Should now if it's NOT a valid file", func() {
			// Arrange
			file := "not_expected_file.ex"

			// Act
			isExpected, _ := j0_submitter.IsExpectedFile(file)

			// Assert
			Expect(isExpected).To(BeFalse())
		})
	})

	Context("When getting the absolute path from the challenge", func() {
		It("Should fail if it's a non existant path", func() {
			// Arrange
			path := "./non_existant_path"

			// Act
			_, err := j0_submitter.GetAbsolutePath(path)

			// Assert
			Expect(err.Error()).To(Equal("Base folder: './non_existant_path' does not exist"))

		})

		It("Should get the absolute path of the file", func() {
			// Arrange
			basePath := "."
			absPath, _ := filepath.Abs(basePath)

			// Act
			calculatedPath, err := j0_submitter.GetAbsolutePath(basePath)

			// Assert
			Expect(err).To(BeNil())
			Expect(calculatedPath).To(Equal(absPath))
		})
	})
})

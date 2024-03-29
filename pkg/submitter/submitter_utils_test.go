package submitter_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"path/filepath"

	"github.com/roeeyn/judge0-uploader/pkg/submitter"
)

var _ = Describe("J0 Submitter Utils Test", func() {
	Context("When having an URL", func() {
		It("Should trim the trailing dash", func() {
			// Arrange
			url := "https://example.com/"

			// Act
			cleanedUrl := submitter.CleanUrl(url)

			// Assert
			Expect(cleanedUrl).To(Equal("https://example.com"))
		})

		It("Should use the URL if doesn't have trailing dash", func() {
			// Arrange
			url := "https://example.com"

			// Act
			cleanedUrl := submitter.CleanUrl(url)

			// Assert
			Expect(cleanedUrl).To(Equal("https://example.com"))
		})
	})

	Context("When having found a file", func() {
		It("Should now if it's a valid file", func() {
			// Arrange
			files := []string{"index.py", "run", "test.js", "testframework.esotericextension"}
			for _, file := range files {
				// Act
				isExpected, _ := submitter.IsExpectedFile(file)

				// Assert
				Expect(isExpected).To(BeTrue())
			}

		})

		It("Should now if it's NOT a valid file", func() {
			// Arrange
			file := "not_expected_file.ex"

			// Act
			isExpected, _ := submitter.IsExpectedFile(file)

			// Assert
			Expect(isExpected).To(BeFalse())
		})
	})

	Context("When getting the absolute path from the challenge", func() {
		It("Should fail if it's a non existant path", func() {
			// Arrange
			path := "./non_existant_path"

			// Act
			_, err := submitter.GetAbsolutePath(path)

			// Assert
			Expect(err.Error()).To(Equal("Base folder: './non_existant_path' does not exist"))

		})

		It("Should get the absolute path of the file", func() {
			// Arrange
			basePath := "."
			absPath, _ := filepath.Abs(basePath)

			// Act
			calculatedPath, err := submitter.GetAbsolutePath(basePath)

			// Assert
			Expect(err).To(BeNil())
			Expect(calculatedPath).To(Equal(absPath))
		})
	})
})

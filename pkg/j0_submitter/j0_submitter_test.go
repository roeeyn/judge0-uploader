package j0_submitter_test

import (
	"log"
	"os"

	"github.com/bxcodec/faker/v3"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	submitter "github.com/roeeyn/judge0-uploader/pkg/j0_submitter"
)

var _ = Describe("Judge0 Submitter Struct Tests", func() {

	Context("When creating a new submitter", func() {

		It("should return a new submitter", func() {
			// Arrange
			DebugLogger := log.New(os.Stdout, "DEBUG: ", log.Ldate|log.Ltime|log.Lshortfile)
			ErrorLogger := log.New(os.Stdout, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)
			InfoLogger := log.New(os.Stdout, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
			WarningLogger := log.New(os.Stdout, "WARNING: ", log.Ldate|log.Ltime|log.Lshortfile)
			ChallengeName := faker.FirstName()

			// Act
			submitter := submitter.NewJ0Submitter(ChallengeName, DebugLogger, ErrorLogger, InfoLogger, WarningLogger)

			// Assert
			Expect(submitter).ToNot(BeNil())
			Expect(submitter.ChallengePath).To(Equal(ChallengeName))
			Expect(submitter.DebugLogger).To(Equal(DebugLogger))
			Expect(submitter.ErrorLogger).To(Equal(ErrorLogger))
			Expect(submitter.InfoLogger).To(Equal(InfoLogger))
			Expect(submitter.WarningLogger).To(Equal(WarningLogger))
		})

	})

})

var _ = Describe("Judge0 Submitter File Struct Tests", func() {

	Context("When creating a new submitter file", func() {
		It("should return a new submitter file", func() {
			// Arrange
			// Act
			file := submitter.NewJ0SubmitterFiles()

			// Assert
			Expect(file).ToNot(BeNil())
		})

		It("should return be able to add files", func() {
			// Arrange
			runName := faker.FirstName()
			indexName := faker.FirstName()
			testName := faker.FirstName()
			testframeworkName := faker.FirstName()

			file := submitter.NewJ0SubmitterFiles()

			// Act
			errorList := []error{
				file.AddFile("index", indexName),
				file.AddFile("run", runName),
				file.AddFile("test", testName),
				file.AddFile("testframework", testframeworkName),
			}

			// Assert
			Expect(file.Run).To(Equal(runName))
			Expect(file.Index).To(Equal(indexName))
			Expect(file.Test).To(Equal(testName))
			Expect(file.Testframework).To(Equal(testframeworkName))

			for _, err := range errorList {
				Expect(err).To(BeNil())
			}

		})

		It("should return an error if invalid file added", func() {
			// Arrange
			file := submitter.NewJ0SubmitterFiles()

			// Act
			err := file.AddFile("nonexisting", "any value")

			// Assert
			Expect(err.Error()).To(Equal("Unknown file key: nonexisting"))

		})

		It("should know if empty properties are held", func() {
			// Arrange
			runName := faker.FirstName()
			indexName := faker.FirstName()
			testName := faker.FirstName()
			testframeworkName := faker.FirstName()

			file := submitter.NewJ0SubmitterFiles()

			// Act
			errorList := []error{
				file.AddFile("index", indexName),
				file.AddFile("run", runName),
				file.AddFile("test", testName),
			}

			shouldBeTrue := file.ContainsEmptyFileProperties()

			extraError := file.AddFile("testframework", testframeworkName)
			shouldBeFalse := file.ContainsEmptyFileProperties()

			// Assert
			Expect(shouldBeTrue).To(BeTrue())
			Expect(shouldBeFalse).To(BeFalse())

			for _, err := range errorList {
				Expect(err).To(BeNil())
			}

			Expect(extraError).To(BeNil())
		})
	})

})

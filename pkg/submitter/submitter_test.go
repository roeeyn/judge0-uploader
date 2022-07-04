package submitter_test

import (
	"io/ioutil"
	"os"
	"path"

	"github.com/bxcodec/faker/v3"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	submitter "github.com/roeeyn/judge0-uploader/pkg/submitter"
)

var baseTestChallengeDir = "test-challenge/"
var fileExtension = ".py"

func CreateTestChallengeDir() (err error) {
	// Create test challenge dir
	err = os.MkdirAll(baseTestChallengeDir, 0755)
	if err != nil {
		return
	}

	// Create the files
	for _, file := range submitter.ExpectedChallengeFiles {
		data := []byte("data for" + file)
		err = ioutil.WriteFile(path.Join(baseTestChallengeDir, file+fileExtension), data, 0777)
		if err != nil {
			return
		}
	}

	return
}

var _ = Describe("Judge0 Submitter Struct Tests", func() {

	Context("When creating a new submitter", func() {

		It("should return a new submitter", func() {
			// Arrange
			ChallengeName := faker.FirstName()
			fakeUrl := faker.URL()
			fakeToken := faker.Word()

			// Act
			submitter := submitter.NewSubmitter(ChallengeName, fakeToken, fakeUrl)

			// Assert
			Expect(submitter).ToNot(BeNil())
			Expect(submitter.ChallengePath).To(Equal(ChallengeName))
			Expect(submitter.AuthToken).To(Equal(fakeToken))
			Expect(submitter.ServerUrl).To(Equal(fakeUrl))
		})

	})

})

var _ = Describe("Judge0 Submitter File Struct Tests", func() {

	Context("When creating a new submitter file", func() {
		It("should return a new submitter file", func() {
			// Arrange
			// Act
			file := submitter.NewSubmitterFiles()

			// Assert
			Expect(file).ToNot(BeNil())
		})

		It("should return be able to add files", func() {
			// Arrange
			runName := faker.FirstName()
			indexName := faker.FirstName()
			testName := faker.FirstName()
			testframeworkName := faker.FirstName()

			file := submitter.NewSubmitterFiles()

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
			file := submitter.NewSubmitterFiles()

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

			file := submitter.NewSubmitterFiles()

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

var _ = Describe("Judge0 Submitter Tests", func() {

	BeforeEach(func() {
		// Initial Cleanup in case of previous test failure
		err := os.RemoveAll(baseTestChallengeDir)
		Expect(err).To(BeNil())

		// Copy the sample challenge directory
		err = CreateTestChallengeDir()
		Expect(err).To(BeNil())
	})

	AfterEach(func() {
		// Cleanup
		err := os.RemoveAll(baseTestChallengeDir)
		Expect(err).To(BeNil())
	})

	Context("When getting the files of the submission", func() {
		It("should get the correct absolute path of the files", func() {
			// Arrange
			j0submitter := submitter.NewSubmitter(baseTestChallengeDir, "", "")
			absFilePath, err := submitter.GetAbsolutePath(baseTestChallengeDir)
			Expect(err).To(BeNil())

			// Act
			j0submitter.AbsChallengePath = absFilePath
			expectedError := j0submitter.GetChallengeFiles()

			// Assert
			Expect(expectedError).To(BeNil())
			Expect(j0submitter.Files.Run).To(Equal(path.Join(absFilePath, "run"+fileExtension)))
			Expect(j0submitter.Files.Index).To(Equal(path.Join(absFilePath, "index"+fileExtension)))
			Expect(j0submitter.Files.Test).To(Equal(path.Join(absFilePath, "test"+fileExtension)))
			Expect(j0submitter.Files.Testframework).To(Equal(path.Join(absFilePath, "testframework"+fileExtension)))
		})

		It("should return an error if an expected file is missing", func() {
			// Arrange
			j0submitter := submitter.NewSubmitter(baseTestChallengeDir, "", "")
			absFilePath, err := submitter.GetAbsolutePath(baseTestChallengeDir)
			Expect(err).To(BeNil())

			// Remove an expected file to trigger the error
			removeErr := os.Remove(path.Join(absFilePath, "index"+fileExtension))
			Expect(removeErr).To(BeNil())

			// Act
			j0submitter.AbsChallengePath = absFilePath
			expectedError := j0submitter.GetChallengeFiles()

			// Assert
			Expect(expectedError).ToNot(BeNil())
			Expect(expectedError.Error()).To(Equal("Not all needed files are present. Expected files are: [index run test testframework]"))

		})

		It("should return an error if the path is not valid", func() {

			// Arrange
			invalid_path := "/invalid/path"
			// Act
			j0submitter := submitter.NewSubmitter(invalid_path, "", "")
			j0submitter.AbsChallengePath = invalid_path
			expectedError := j0submitter.GetChallengeFiles()

			// Assert
			Expect(expectedError).ToNot(BeNil())
			Expect(expectedError.Error()).To(Equal("Error reading files inside folder: open /invalid/path: no such file or directory"))

		})
	})

})

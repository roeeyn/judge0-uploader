package j0_submitter_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gstruct"

	"github.com/bxcodec/faker/v3"
	submitter "github.com/roeeyn/judge0-uploader/pkg/j0_submitter"
)

var _ = Describe("Judge0 Submitter Tests", func() {

	Context("When creating a new submitter", func() {

		It("should return a new submitter", func() {
			// Arrange
			runName := faker.FirstName()
			indexName := faker.FirstName()
			testName := faker.FirstName()
			testframeworkName := faker.FirstName()

			files := submitter.J0SubmitterFiles{
				Run:           runName,
				Index:         indexName,
				Test:          testName,
				Testframework: testframeworkName,
			}

			// Act
			submitter := submitter.NewJ0Submitter(&files)

			// Assert
			Expect(submitter).ToNot(BeNil())
			Expect(submitter.Files).To(PointTo(Equal(files)))
		})

	})

})

var _ = Describe("Judge0 Submitter File Tests", func() {

})

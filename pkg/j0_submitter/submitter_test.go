package j0_submitter_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gstruct"

	submitter "github.com/roeeyn/judge0-uploader/pkg/j0_submitter"
)

var _ = Describe("Judge0 Submitter Tests", func() {

	Context("When creating a new submitter", func() {

		It("should return a new submitter", func() {
			// Arrange
			files := submitter.J0SubmitterFiles{}

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

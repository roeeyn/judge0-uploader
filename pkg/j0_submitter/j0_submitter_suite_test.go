package j0_submitter_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestJ0Submitter(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "J0Submitter Suite")
}

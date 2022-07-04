package submitter_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestSubmitter(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Submitter Suite")
}

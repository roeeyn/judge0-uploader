.SILENT:
.PHONY: test-submitter coverage-submitter run-submitter

test-submitter:
	ginkgo pkg/j0_submitter/

coverage-submitter:
	ginkgo -coverprofile=coverage.out pkg/j0_submitter/
	go tool cover -html coverage.out

run-submitter:
	go run main.go submit sample_challenge

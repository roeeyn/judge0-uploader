.SILENT:
.PHONY: test-submitter coverage-submitter run-submitter

submitter-test:
	ginkgo pkg/j0_submitter/

submitter-coverage:
	ginkgo -coverprofile=coverage.out pkg/j0_submitter/
	go tool cover -html coverage.out

submitter-run:
	go run main.go submit sample_challenge

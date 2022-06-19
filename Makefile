.SILENT:

test-submitter:
	ginkgo pkg/j0_submitter/

run-submitter:
	go run main.go submit sample_challenge

package status

import (
	logger "github.com/roeeyn/judge0-uploader/pkg/logger"
)

type StatusFetcher struct {
	AuthToken string
	ServerUrl string
}

type SubmissionStatus struct {
	// status string
}

func NewStatusFetcher(authToken string, serverUrl string) *StatusFetcher {
	return &StatusFetcher{
		AuthToken: authToken,
		ServerUrl: serverUrl,
	}
}

func (statusFetcher *StatusFetcher) Run() (status SubmissionStatus, err error) {
	logger.LogInfo("Running")
	return
}

package status

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httputil"

	logger "github.com/roeeyn/judge0-uploader/pkg/logger"
	utils "github.com/roeeyn/judge0-uploader/pkg/utils"
)

type StatusFetcher struct {
	AuthToken    string
	ServerUrl    string
	SubmissionId string
}

type SubmissionStatus struct {
	Id          int    `json:"id"`
	Description string `json:"description"`
}

type SubmissionStatusResponse struct {
	Stdout        string           `json:"stdout"`
	Time          string           `json:"time"`
	Memory        int              `json:"memory"`
	Stderr        string           `json:"stderr"`
	Token         string           `json:"token"`
	CompileOutput string           `json:"compile_output"`
	Message       string           `json:"message"`
	Status        SubmissionStatus `json:"status"`
}

func NewStatusFetcher(authToken string, serverUrl string, submissionId string) StatusFetcher {
	return StatusFetcher{
		AuthToken:    authToken,
		ServerUrl:    serverUrl,
		SubmissionId: submissionId,
	}
}

func (statusFetcher StatusFetcher) GetSubmissionStatus() (submissionStatusResponse SubmissionStatusResponse, err error) {

	url := utils.CleanUrl(statusFetcher.ServerUrl) + "/submissions/" + statusFetcher.SubmissionId

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		err = fmt.Errorf("Error creating the GET request: %s", err.Error())
		return
	}

	// req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Auth-Token", statusFetcher.AuthToken)

	reqDump, err := httputil.DumpRequestOut(req, true)
	if err != nil {
		err = fmt.Errorf("Error dumping request: %s", err.Error())
		return
	}

	logger.LogInfo(fmt.Sprintf("REQUEST:\n%s", string(reqDump)))

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		err = fmt.Errorf("Error sending the POST request: %s", err.Error())
		return
	}

	logger.LogInfo(fmt.Sprintf("RESPONSE status: %s", resp.Status))

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		err = fmt.Errorf("Error reading response body: %s", err.Error())
		return
	}

	logger.LogInfo(fmt.Sprintf("Request sent successfully:\n %s", string(body)))

	err = json.Unmarshal(body, &submissionStatusResponse)
	if err != nil {
		err = fmt.Errorf("Error unmarshalling response: %s", err.Error())
		return
	}

	logger.LogInfo(fmt.Sprintf("Submission status: %s", submissionStatusResponse.Status.Description))
	return submissionStatusResponse, nil
}

func (statusFetcher StatusFetcher) Run() (response SubmissionStatusResponse, err error) {
	logger.LogInfo("Running")
	response, err = statusFetcher.GetSubmissionStatus()
	return
}

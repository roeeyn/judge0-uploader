package submitter

import (
	"archive/zip"
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httputil"
	"os"
	"path"
	"path/filepath"

	logger "github.com/roeeyn/judge0-uploader/pkg/logger"
	utils "github.com/roeeyn/judge0-uploader/pkg/utils"
)

const ZIP_FILE_NAME = "upload.judge0.zip"
const MAIN_FILE_NAME = "upload.judge0"

type SubmitterFiles struct {
	Index string
	Run   string
	Test  string
	// We're leaving it without camelcase to maintain consistency
	// with the expected challenge file.
	Testframework string
}

type Submitter struct {
	AbsChallengePath string
	AuthToken        string
	ChallengePath    string
	EncodedZipFile   string
	Files            *SubmitterFiles
	ServerUrl        string
}

type Submission struct {
	SubmissionId string `json:"token"`
}

func (submitterFiles *SubmitterFiles) AddFile(fileKey string, fileAbsPath string) error {
	switch fileKey {
	case "run":
		submitterFiles.Run = fileAbsPath
	case "index":
		submitterFiles.Index = fileAbsPath
	case "test":
		submitterFiles.Test = fileAbsPath
	case "testframework":
		submitterFiles.Testframework = fileAbsPath
	default:
		return fmt.Errorf("Unknown file key: %s", fileKey)
	}
	return nil
}

func (submitterFiles *SubmitterFiles) ContainsEmptyFileProperties() bool {
	// Iterate over the struct fields and check if any of them is empty
	// because at this point every file is needed.
	if submitterFiles.Run == "" || submitterFiles.Index == "" || submitterFiles.Test == "" || submitterFiles.Testframework == "" {
		return true
	}

	return false
}

func (submitter *Submitter) Run() (submissionId string, err error) {
	challengePath := submitter.ChallengePath
	absPath, err := GetAbsolutePath(challengePath)
	if err != nil {
		return
	}

	logger.LogInfo(fmt.Sprintf("Absolute path: %s", absPath))
	submitter.AbsChallengePath = absPath

	err = submitter.GetChallengeFiles()
	if err != nil {
		return
	}

	err = submitter.ZipFiles()
	if err != nil {
		return
	}

	encodedFile, err := B64ZipFile()
	if err != nil {
		return
	}

	logger.LogInfo("Encoded file successfully")
	submitter.EncodedZipFile = encodedFile

	submissionId, err = submitter.SubmitEncodedFile()
	if err != nil {
		return
	}

	err = submitter.Cleanup()
	if err != nil {
		return
	}

	logger.LogInfo("Cleanup executed successfully")
	return
}

func NewSubmitterFiles() (submitterFiles *SubmitterFiles) {
	return &SubmitterFiles{}
}

func NewSubmitter(challengePath string, authToken string, serverUrl string) (submitter *Submitter) {
	submitter = &Submitter{
		Files:         NewSubmitterFiles(),
		ChallengePath: challengePath,
		AuthToken:     authToken,
		ServerUrl:     serverUrl,
	}
	return
}

func (submitter *Submitter) Cleanup() (err error) {
	err = os.Remove(ZIP_FILE_NAME)
	return
}

func (submitter *Submitter) GetChallengeFiles() (err error) {
	// Verify that the basePath contains expected files
	absBasePath := submitter.AbsChallengePath
	files, err := ioutil.ReadDir(absBasePath)
	if err != nil {
		err = fmt.Errorf("Error reading files inside folder: %s", err.Error())
		return
	}

	for _, file := range files {
		// We do not support nested folders
		if file.IsDir() {
			continue
		}

		absFilePath := path.Join(absBasePath, file.Name())
		logger.LogInfo(fmt.Sprintf("Found file: %s", absFilePath))

		if isExpected, fileNameKey := IsExpectedFile(file.Name()); isExpected {
			error := submitter.Files.AddFile(fileNameKey, absFilePath)
			logger.LogInfo(fmt.Sprintf("ADDED file: %s", absFilePath))
			if error != nil {
				err = fmt.Errorf("Error adding file property: %s", error.Error())
				return
			}
		} else {
			logger.LogInfo(fmt.Sprintf("IGNORED file: %s", absFilePath))
		}
	}

	if submitter.Files.ContainsEmptyFileProperties() {
		err = fmt.Errorf("Not all needed files are present. Expected files are: %s", ExpectedChallengeFiles)
		return
	}

	logger.LogInfo(fmt.Sprintf("Challenge Files: %s", submitter.Files))
	return
}

func (submitter *Submitter) ZipFiles() (err error) {
	logger.LogInfo("Creating zip archive...")
	archive, err := os.Create(ZIP_FILE_NAME)
	if err != nil {
		return
	}
	defer archive.Close()

	zipWriter := zip.NewWriter(archive)

	filesMap := map[string]string{
		"index":         submitter.Files.Index,
		"testframework": submitter.Files.Testframework,
		"test":          submitter.Files.Test,
	}

	logger.LogInfo(fmt.Sprintf("Files map: %s", filesMap))

	mainContent := ""

	for fileName, filePath := range filesMap {
		fileExtension := filepath.Ext(filePath)
		fullFileName := fileName + fileExtension

		logger.LogInfo(fmt.Sprintf("Opening: %s", fullFileName))

		content, openErr := ioutil.ReadFile(filePath)
		if openErr != nil {
			return openErr
		}

		logger.LogInfo(fmt.Sprintf("Appending '%s' to archive...", fullFileName))

		mainContent = mainContent + string(content) + "\n"

	}

	mainWriter, err := zipWriter.Create(MAIN_FILE_NAME)
	if err != nil {
		return
	}

	_, err = mainWriter.Write([]byte(mainContent))
	if err != nil {
		return
	}

	runContent, err := ioutil.ReadFile(submitter.Files.Run)
	if err != nil {
		return
	}

	runWriter, err := zipWriter.Create("run")
	if err != nil {
		return
	}

	_, err = runWriter.Write([]byte(runContent))
	if err != nil {
		return
	}

	logger.LogInfo("Closing zip archive...")
	zipWriter.Close()

	return
}

func (submitter Submitter) SubmitEncodedFile() (submissionId string, err error) {

	values := map[string]string{"language_id": "89", "additional_files": submitter.EncodedZipFile}
	json_data, err := json.Marshal(values)

	if err != nil {
		err = fmt.Errorf("Error marshalling request data json: %s", err.Error())
		return
	}

	url := utils.CleanUrl(submitter.ServerUrl) + "/submissions?wait=false"

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(json_data))
	if err != nil {
		err = fmt.Errorf("Error creating the POST request: %s", err.Error())
		return
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Auth-Token", submitter.AuthToken)

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

	var submission Submission
	err = json.Unmarshal(body, &submission)
	if err != nil {
		err = fmt.Errorf("Error unmarshalling response body: %s", err.Error())
		return
	}

	logger.LogInfo(fmt.Sprintf("Obtained Submission ID: %s", submission.SubmissionId))

	return submission.SubmissionId, nil
}

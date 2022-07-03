package j0_submitter

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
)

const ZIP_FILE_NAME = "upload.judge0.zip"
const MAIN_FILE_NAME = "upload.judge0"

type J0SubmitterFiles struct {
	Run   string
	Index string
	Test  string
	// We're leaving it without camelcase to maintain consistency
	// with the expected challenge file.
	Testframework string
}

type J0Submitter struct {
	Files            *J0SubmitterFiles
	AuthToken        string
	ChallengePath    string
	AbsChallengePath string
	IsVerbose        bool
	EncodedZipFile   string
}

type Submission struct {
	SubmissionId string `json:"token"`
}

func (j0SubmitterFiles *J0SubmitterFiles) AddFile(fileKey string, fileAbsPath string) error {
	switch fileKey {
	case "run":
		j0SubmitterFiles.Run = fileAbsPath
	case "index":
		j0SubmitterFiles.Index = fileAbsPath
	case "test":
		j0SubmitterFiles.Test = fileAbsPath
	case "testframework":
		j0SubmitterFiles.Testframework = fileAbsPath
	default:
		return fmt.Errorf("Unknown file key: %s", fileKey)
	}
	return nil
}

func (j0SubmitterFiles *J0SubmitterFiles) ContainsEmptyFileProperties() bool {
	// Iterate over the struct fields and check if any of them is empty
	// because at this point every file is needed.
	if j0SubmitterFiles.Run == "" || j0SubmitterFiles.Index == "" || j0SubmitterFiles.Test == "" || j0SubmitterFiles.Testframework == "" {
		return true
	}

	return false
}

func (j0Submitter *J0Submitter) Run() (submissionId string, err error) {
	challengePath := j0Submitter.ChallengePath
	absPath, err := GetAbsolutePath(challengePath)
	if err != nil {
		return
	}

	fmt.Println("Absolute path:", absPath)
	j0Submitter.AbsChallengePath = absPath

	err = j0Submitter.GetChallengeFiles()
	if err != nil {
		return
	}

	err = j0Submitter.ZipFiles()
	if err != nil {
		return
	}

	encodedFile, err := B64ZipFile()
	if err != nil {
		return
	}

	fmt.Println("Encoded file successfully")
	j0Submitter.EncodedZipFile = encodedFile

	submissionId, err = j0Submitter.SubmitEncodedFile()
	return
}

func NewJ0SubmitterFiles() (j0SubmitterFiles *J0SubmitterFiles) {
	return &J0SubmitterFiles{}
}

func NewJ0Submitter(challengePath string, authToken string, isVerbose bool) (j0Submitter *J0Submitter) {
	j0Submitter = &J0Submitter{
		Files:         NewJ0SubmitterFiles(),
		ChallengePath: challengePath,
		AuthToken:     authToken,
		IsVerbose:     isVerbose,
	}
	return
}

func (j0Submitter *J0Submitter) GetChallengeFiles() (err error) {
	// Verify that the basePath contains expected files
	absBasePath := j0Submitter.AbsChallengePath
	files, err := ioutil.ReadDir(absBasePath)
	if err != nil {
		err = fmt.Errorf(fmt.Sprintf("Error reading files inside folder: %s", err.Error()))
		return
	}

	for _, file := range files {
		// We do not support nested folders
		if file.IsDir() {
			continue
		}

		absFilePath := path.Join(absBasePath, file.Name())
		fmt.Println("Found file: ", absFilePath)

		if isExpected, fileNameKey := IsExpectedFile(file.Name()); isExpected {
			error := j0Submitter.Files.AddFile(fileNameKey, absFilePath)
			fmt.Println("Added file: ", absFilePath)
			if error != nil {
				err = fmt.Errorf(fmt.Sprintf("Error adding file property: %s", error.Error()))
				return
			}
		} else {
			fmt.Println("Ignored file: ", absFilePath)
		}
	}

	if j0Submitter.Files.ContainsEmptyFileProperties() {
		err = fmt.Errorf(fmt.Sprintf("Not all needed files are present. Expected files are: %s", ExpectedChallengeFiles))
		fmt.Println("Current Files:", j0Submitter.Files)
		return
	}

	fmt.Println("Challenge Files:", j0Submitter.Files)
	return
}

func (j0Submitter *J0Submitter) ZipFiles() (err error) {
	fmt.Println("creating zip archive...")
	archive, err := os.Create(ZIP_FILE_NAME)
	if err != nil {
		return
	}
	defer archive.Close()

	zipWriter := zip.NewWriter(archive)

	filesMap := map[string]string{
		"index":         j0Submitter.Files.Index,
		"testframework": j0Submitter.Files.Testframework,
		"test":          j0Submitter.Files.Test,
	}

	fmt.Println(filesMap)

	mainContent := ""

	for fileName, filePath := range filesMap {
		fileExtension := filepath.Ext(filePath)
		fullFileName := fileName + fileExtension

		fmt.Println("Opening: ", fullFileName)

		content, openErr := ioutil.ReadFile(filePath)
		if openErr != nil {
			return openErr
		}

		fmt.Printf("Appending '%s' to archive...\n", fullFileName)

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

	runContent, err := ioutil.ReadFile(j0Submitter.Files.Run)
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

	fmt.Println("closing zip archive...")
	zipWriter.Close()

	return
}

func (j0Submitter J0Submitter) SubmitEncodedFile() (submissionId string, err error) {

	values := map[string]string{"language_id": "89", "additional_files": j0Submitter.EncodedZipFile}
	json_data, err := json.Marshal(values)

	if err != nil {
		err = fmt.Errorf(fmt.Sprintf("Error marshalling request data json: %s", err.Error()))
		return
	}

	url := "https://judge.hackademy.mx/submissions?wait=false"

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(json_data))
	if err != nil {
		err = fmt.Errorf(fmt.Sprintf("Error creating the POST request: %s", err.Error()))
		return
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Auth-Token", j0Submitter.AuthToken)

	reqDump, err := httputil.DumpRequestOut(req, true)
	if err != nil {
		err = fmt.Errorf(fmt.Sprintf("Error dumping request: %s", err.Error()))
		return
	}

	fmt.Printf("REQUEST:\n%s\n", string(reqDump))

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		err = fmt.Errorf(fmt.Sprintf("Error sending the POST request: %s", err.Error()))
		return
	}

	fmt.Printf("RESPONSE status: %s\n", resp.Status)

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		err = fmt.Errorf(fmt.Sprintf("Error reading response body: %s", err.Error()))
		return
	}

	fmt.Printf("Request sent successfully: %s\n", string(body))

	var submission Submission
	err = json.Unmarshal(body, &submission)
	if err != nil {
		err = fmt.Errorf(fmt.Sprintf("Error unmarshalling response body: %s", err.Error()))
		return
	}

	fmt.Printf("Obtained Submission ID: %s\n", submission.SubmissionId)

	return submission.SubmissionId, nil
}

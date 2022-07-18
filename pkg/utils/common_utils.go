package common_utils

import (
	"encoding/json"
	"strings"
)

func PrettyPrint(i interface{}) string {
	s, _ := json.MarshalIndent(i, "", "\t")
	return string(s)
}

func CleanUrl(url string) (cleanedUrl string) {
	cleanedUrl = strings.TrimSuffix(url, "/")
	return
}

package utils

import (
	"archive/zip"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"runtime"
	"strings"
)

const (
	DESMUME_ORG  = "TASEmulators"
	DESMUME_REPO = "desmume"
)

func downloadFile(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func DownloadAndExtract(destinationPath string) error {
	var platform string
	goos := runtime.GOOS

	if goos == "windows" {
		platform = "win64"
	} else if goos == "darwin" {
		platform = "macOS"
	} else if goos == "linux" {
		platform = "OpenEmu"
	} else {
		return fmt.Errorf("unsupported OS: %s", goos)
	}

	respData, err := downloadFile(fmt.Sprintf("https://api.github.com/repos/%s/%s/releases/latest", DESMUME_ORG, DESMUME_REPO))
	if err != nil {
		return fmt.Errorf("error getting latest release: %v", err)
	}

	var releaseData map[string]interface{}
	if err := json.Unmarshal(respData, &releaseData); err != nil {
		return fmt.Errorf("error parsing release data: %v", err)
	}

	assets, ok := releaseData["assets"].([]interface{})
	if !ok || len(assets) == 0 {
		return fmt.Errorf("no assets found in the latest release")
	}

	var assetURL, assetName string
	for _, asset := range assets {
		assetInfo := asset.(map[string]interface{})
		name := assetInfo["name"].(string)
		if strings.Contains(name, platform) {
			assetURL = assetInfo["browser_download_url"].(string)
			assetName = name
			break
		}
	}

	if assetURL == "" || assetName == "" {
		return fmt.Errorf("no matching assets found in the latest release")
	}

	zipData, err := downloadFile(assetURL)
	if err != nil {
		return fmt.Errorf("error downloading asset: %v", err)
	}

	zipReader, err := zip.NewReader(strings.NewReader(string(zipData)), int64(len(zipData)))
	if err != nil {
		return fmt.Errorf("error creating zip reader: %v", err)
	}

	for _, file := range zipReader.File {
		fileReader, err := file.Open()
		if err != nil {
			return fmt.Errorf("error opening file inside zip: %v", err)
		}
		defer fileReader.Close()

		destFilePath := destinationPath + "/" + file.Name
		destFile, err := os.Create(destFilePath)
		if err != nil {
			return fmt.Errorf("error creating destination file: %v", err)
		}
		defer destFile.Close()

		_, err = io.Copy(destFile, fileReader)
		if err != nil {
			return fmt.Errorf("error copying file contents: %v", err)
		}
	}

	return nil
}

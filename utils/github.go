package utils

import (
	"archive/zip"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
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

	fmt.Println("Checking for goos=", goos)

	if goos == "windows" {
		platform = "win64"
	} else if goos == "darwin" {
		platform = "macOS"
	} else {
		fmt.Println("Unsupported platform")
		return fmt.Errorf("unsupported platform")
	}

	fmt.Println("Downloading latest release for platform ", platform)

	respData, err := downloadFile(fmt.Sprintf("https://api.github.com/repos/%s/%s/releases/latest", DESMUME_ORG, DESMUME_REPO))
	if err != nil {
		fmt.Println("Error getting latest release")
		return fmt.Errorf("error getting latest release: %v", err)
	}

	var releaseData map[string]interface{}
	if err := json.Unmarshal(respData, &releaseData); err != nil {
		fmt.Println("Error parsing release data")
		return fmt.Errorf("error parsing release data: %v", err)
	}

	assets, ok := releaseData["assets"].([]interface{})
	if !ok || len(assets) == 0 {
		fmt.Println("No assets found in the latest release")
		return fmt.Errorf("no assets found in the latest release")
	}

	fmt.Println(assets)

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
		fmt.Println("No matching assets found in the latest release")
		return fmt.Errorf("no matching assets found in the latest release")
	}

	zipData, err := downloadFile(assetURL)
	if err != nil {
		fmt.Println("Error downloading asset")
		return fmt.Errorf("error downloading asset: %v", err)
	}

	zipReader, err := zip.NewReader(strings.NewReader(string(zipData)), int64(len(zipData)))
	if err != nil {
		fmt.Println("Error creating zip reader")
		return fmt.Errorf("error creating zip reader: %v", err)
	}

	for _, file := range zipReader.File {
		fileReader, err := file.Open()
		if err != nil {
			fmt.Println("Error opening file inside zip")
			return fmt.Errorf("error opening file inside zip: %v", err)
		}

		defer fileReader.Close()

		destFilePath := filepath.Join(destinationPath, file.Name)
		destFile, err := os.Create(destFilePath)
		if err != nil {
			fmt.Println("Error creating destination file")
			return fmt.Errorf("error creating destination file: %v", err)
		}
		defer destFile.Close()

		_, err = io.Copy(destFile, fileReader)
		if err != nil {
			fmt.Println("Error copying file contents")
			return fmt.Errorf("error copying file contents: %v", err)
		}
	}

	return nil
}

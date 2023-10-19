package utils

import (
	"encoding/binary"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"unsafe"

	nds "github.com/szerookii/desmume-launcher/nds"
)

func LoadGameFiles() ([]*nds.NDSFile, error) {
	dir, err := os.Open(filepath.Join(".", "games"))
	if err != nil {
		return nil, err
	}

	defer dir.Close()

	fileInfos, err := dir.Readdir(-1)
	if err != nil {
		return nil, err
	}

	var ndsFiles []*nds.NDSFile

	for _, fileInfo := range fileInfos {
		if !fileInfo.IsDir() && strings.HasSuffix(fileInfo.Name(), ".nds") {
			fullPath := filepath.Join(".", "games", fileInfo.Name())
			file, err := os.Open(fullPath)
			if err != nil {
				continue
			}

			var header nds.NDSHeader
			err = binary.Read(file, binary.LittleEndian, &header)
			if err != nil {
				continue
			}

			var banner nds.NDSBanner
			bannerBytes := (*[unsafe.Sizeof(banner)]byte)(unsafe.Pointer(&banner))[:]
			_, err = file.ReadAt(bannerBytes, int64(header.IconBannerOffset))
			if err != nil {
				continue
			}

			icon, err := banner.IconPNG(128)
			if err != nil {
				continue
			}

			file.Close()

			ndsFiles = append(ndsFiles, &nds.NDSFile{
				Path:              fileInfo.Name(),
				Size:              FormatBytes(float64(fileInfo.Size())),
				Name:              banner.FrenchTitleString(),
				Developer:         banner.Author(),
				Base64EncodedIcon: string(EncodeBytesToBase64(icon)),
			})

			fmt.Printf("Loaded %s\n", fileInfo.Name())
		}
	}

	for _, ndsFile := range ndsFiles {
		fmt.Printf("%s\n", ndsFile.Name)
	}

	return ndsFiles, nil
}

func FileExists(path string) bool {
	var err error

	if _, err := os.Stat(path); err == nil {
		return true
	}

	if os.IsNotExist(err) {
		return false
	}

	return false
}

func FormatBytes(bytes float64) string {
	units := []string{"B", "KB", "MB", "GB", "TB", "PB", "EB", "ZB", "YB"}
	size := bytes
	var i int

	for i = 0; size >= 1024 && i < len(units)-1; i++ {
		size /= 1024
	}
	return fmt.Sprintf("%.2f%s", size, units[i])
}

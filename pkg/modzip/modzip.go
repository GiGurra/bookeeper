package modzip

import (
	"archive/zip"
	"encoding/json"
	"fmt"
	"github.com/GiGurra/bookeeper/pkg/config"
	"io"
	"log"
	"log/slog"
	"os"
	"path/filepath"
)

type ModData struct {
	Mods []ModDataEntry `json:"Mods"`
	MD5  string         `json:"MD5"`
}

type ModDataEntry struct {
	Author       string `json:"Author"`
	Name         string `json:"Name"`
	Folder       string `json:"Folder"`
	Version      string `json:"Version"`
	Description  string `json:"Description"`
	UUID         string `json:"UUID"`
	Created      string `json:"Created"`
	Dependencies []any  `json:"Dependencies"`
	Group        string `json:"Group"`
}

func InspectModZip(
	zipPath string,
) ModData {
	if !config.ExistsFile(zipPath) {
		panic(fmt.Errorf("mod zip file not found: %s", zipPath))
	}

	tmpDir := DeflateZipFileToTempFolder(zipPath)
	defer func() {
		err := os.RemoveAll(tmpDir)
		if err != nil {
			panic(fmt.Errorf("failed to remove temporary directory: %w", err))
		}
	}()

	infoJsonPath := filepath.Join(tmpDir, "info.json")
	if !config.ExistsFile(infoJsonPath) {
		panic(fmt.Errorf("info.json not found in zip: %s", zipPath))
	}

	// parse info.json into ModData
	var modData ModData
	bs, err := os.ReadFile(infoJsonPath)
	if err != nil {
		panic(fmt.Errorf("failed to read file: %w", err))
	}

	err = json.Unmarshal(bs, &modData)
	if err != nil {
		panic(fmt.Errorf("failed to unmarshal json: %w", err))
	}

	return modData
}

func DeflateZipFileToTempFolder(zipPath string) string {
	// Create a temporary directory
	tempDir, err := os.MkdirTemp("", "zip_extract_*")
	if err != nil {
		log.Fatalf("Failed to create temporary directory: %v", err)
		return ""
	}

	// Open the zip file
	reader, err := zip.OpenReader(zipPath)
	if err != nil {
		log.Fatalf("Failed to open zip file: %v", err)
		return ""
	}
	defer func(reader *zip.ReadCloser) {
		err := reader.Close()
		if err != nil {
			slog.Error(fmt.Sprintf("Failed to close zip file: %v", err))
		}
	}(reader)

	// Extract each file from the zip archive
	for _, file := range reader.File {
		// Construct the full path for the extracted file
		path := filepath.Join(tempDir, file.Name)

		// If the entry is a directory, create it
		if file.FileInfo().IsDir() {
			err := os.MkdirAll(path, file.Mode())
			if err != nil {
				panic(fmt.Errorf("failed to create directory: %w", err))
			}
			continue
		}

		// Ensure the parent directory exists
		if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
			panic(fmt.Errorf("failed to create parent directory: %w", err))
		}

		func() {
			// Create the file
			dstFile, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, file.Mode())
			if err != nil {
				panic(fmt.Errorf("failed to create file: %w", err))
			}
			defer func() {
				err := dstFile.Close()
				if err != nil {
					panic(fmt.Errorf("failed to close file: %w", err))
				}
			}()

			// Open the file in the zip
			srcFile, err := file.Open()
			if err != nil {
				panic(fmt.Errorf("failed to open file in zip: %w", err))
			}
			defer func() {
				err := srcFile.Close()
				if err != nil {
					panic(fmt.Errorf("failed to close file in zip: %w", err))
				}
			}()

			// Copy the contents
			_, err = io.Copy(dstFile, srcFile)
			if err != nil {
				panic(fmt.Errorf("failed to copy file: %w", err))
			}
		}()
	}

	return tempDir
}

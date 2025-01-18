package modzip

import (
	"archive/zip"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"github.com/GiGurra/bookeeper/pkg/config"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
)

type ModData struct {
	Mods []ModDataEntry `json:"Mods"`
	MD5  string         `json:"MD5"`
}

func ReadInfoJson(path string) ModData {
	bs, err := os.ReadFile(path)
	if err != nil {
		panic(fmt.Errorf("failed to read file: %w", err))
	}

	var modData ModData
	err = json.Unmarshal(bs, &modData)
	if err != nil {
		panic(fmt.Errorf("failed to unmarshal json: %w", err))
	}

	return modData
}

func (m ModData) Entry() ModDataEntry {

	if len(m.Mods) == 0 {
		panic("No mods found in zip")
	}
	if len(m.Mods) > 2 {
		panic("Multiple mods found in zip, not supported")
	}
	if len(m.Mods) == 2 {
		// only OK if one is GustavDev (which is the game itself)
		first := m.Mods[0]
		second := m.Mods[1]
		if first.Name == "GustavDev" {
			return second
		} else if second.Name == "GustavDev" {
			return first
		} else {
			panic("Multiple mods found in zip, not supported")
		}
	}

	return m.Mods[0]
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

func InspectModZipUsingTempFolderExtract(
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

func ExtractSpecificFilesFromZip(
	zipPath string,
	fileNames []string,
	targetDir string,
) {
	if !config.ExistsFile(zipPath) {
		panic(fmt.Errorf("mod zip file not found: %s", zipPath))
	}

	// Open the zip file
	reader, err := zip.OpenReader(zipPath)
	if err != nil {
		panic(fmt.Errorf("failed to open zip file: %w", err))
	}
	defer func() {
		err := reader.Close()
		if err != nil {
			panic(fmt.Errorf("failed to close zip file: %w", err))
		}
	}()

	for _, fileName := range fileNames {
		file, err := reader.Open(fileName)
		if err != nil {
			panic(fmt.Errorf("failed to open file in zip: %w", err))
		}

		// calculate the target path
		targetPath := filepath.Join(targetDir, fileName)

		// check if the target file exists
		// only continue if the checksums are equal, otherwise panic
		if config.ExistsFile(targetPath) {
			targetFileReader, err := os.Open(targetPath)
			if err != nil {
				panic(fmt.Errorf("failed to open target file: %w", err))
			}
			targetFileChecksum := func() string {
				defer func() {
					err := targetFileReader.Close()
					if err != nil {
						panic(fmt.Errorf("failed to close target file: %w", err))
					}
				}()
				return checksumFile(targetFileReader)
			}()
			if checksumFile(file) != targetFileChecksum {
				_ = targetFileReader.Close()
				panic(fmt.Errorf("checksum mismatch for file: %s. reconciliation not yet implemented", targetPath))
			} else {
				continue // go to next file :S
			}
		}

		// Now copy the file
		func() {
			dstFile, err := os.OpenFile(targetPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
			if err != nil {
				panic(fmt.Errorf("failed to create file: %w", err))
			}
			defer func() {
				err := dstFile.Close()
				if err != nil {
					panic(fmt.Errorf("failed to close file: %w", err))
				}
			}()

			_, err = io.Copy(dstFile, file)
			if err != nil {
				panic(fmt.Errorf("failed to copy file: %w", err))
			}
		}()
	}
}

func checksumFile(fileReader io.Reader) string {
	hash := sha256.New()
	_, err := io.Copy(hash, fileReader)
	if err != nil {
		panic(fmt.Errorf("failed to copy file: %w", err))
	}

	return fmt.Sprintf("%x", hash.Sum(nil))
}

func InspectModZip(
	zipPath string,
) (ModData ModData, PakFiles []string) {
	if !config.ExistsFile(zipPath) {
		panic(fmt.Errorf("mod zip file not found: %s", zipPath))
	}

	// Open the zip file
	reader, err := zip.OpenReader(zipPath)
	if err != nil {
		panic(fmt.Errorf("failed to open zip file: %w", err))
	}
	defer func() {
		err := reader.Close()
		if err != nil {
			panic(fmt.Errorf("failed to close zip file: %w", err))
		}
	}()

	// Extract info.json from the zip archive
	var foundInfoJson *zip.File
	for _, file := range reader.File {
		if file.Name == "info.json" {
			foundInfoJson = file
		}
		if strings.HasSuffix(strings.ToLower(file.Name), ".pak") {
			PakFiles = append(PakFiles, file.Name)
		}
	}
	if foundInfoJson == nil {
		panic(fmt.Errorf("info.json not found in zip: %s", zipPath))
	}
	if len(PakFiles) == 0 {
		panic(fmt.Errorf("no .pak file found in zip: %s", zipPath))
	}

	// Open the file in the zip
	srcFile, err := foundInfoJson.Open()
	if err != nil {
		panic(fmt.Errorf("failed to open file in zip: %w", err))
	}
	defer func() {
		err := srcFile.Close()
		if err != nil {
			panic(fmt.Errorf("failed to close file in zip: %w", err))
		}
	}()

	// Read the contents
	bs, err := io.ReadAll(srcFile)
	if err != nil {
		panic(fmt.Errorf("failed to read file: %w", err))
	}

	// parse info.json into ModData
	err = json.Unmarshal(bs, &ModData)
	if err != nil {
		panic(fmt.Errorf("failed to unmarshal json: %w", err))
	}

	return ModData, PakFiles
}

func DeflateZipFileToTempFolder(zipPath string) string {
	// Create a temporary directory
	tempDir, err := os.MkdirTemp("", "zip_extract_*")
	if err != nil {
		panic(fmt.Errorf("failed to create temporary directory: %w", err))
	}

	// Open the zip file
	reader, err := zip.OpenReader(zipPath)
	if err != nil {
		panic(fmt.Errorf("failed to open zip file: %w", err))
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

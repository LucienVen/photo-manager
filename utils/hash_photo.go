package utils

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/LucienVen/photo-manager/objects"
)

func GetFileSHA256(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	hasher := sha256.New()
	if _, err := io.Copy(hasher, file); err != nil {
		return "", err
	}

	hash := hasher.Sum(nil)
	return fmt.Sprintf("%x", hash), nil
}

func HashExists(targetHash, recordsDir string) (bool, string, error) {
	files, err := os.ReadDir(recordsDir)
	if err != nil {
		return false, "", err
	}

	for _, file := range files {
		if file.IsDir() || !strings.HasSuffix(file.Name(), ".json") {
			continue
		}

		path := filepath.Join(recordsDir, file.Name())
		content, err := os.ReadFile(path)
		if err != nil {
			return false, "", err
		}

		var records []objects.PhotoRecord
		if err := json.Unmarshal(content, &records); err != nil {
			return false, "", fmt.Errorf("解析 %s 出错: %w", path, err)
		}

		for _, r := range records {
			if r.Hash == targetHash {
				return true, path, nil
			}
		}
	}

	return false, "", nil
}

package utils

import (
	"crypto/sha256"
	"fmt"
	"io"
	"os"
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

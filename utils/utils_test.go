package utils

import (
	"github.com/LucienVen/photo-manager/config"
	"testing"
)

func TestCleanTestFile(t *testing.T) {
	config.InitConfig()

	err := CleanTestFile(config.RootPath)
	if err != nil {
		t.Error(err)
	}
}

func TestCheckPicgoExecutable(t *testing.T) {
	config.InitConfig()
	CheckPicgoExecutable()
}

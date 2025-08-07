package utils

import (
	"github.com/LucienVen/photo-manager/config"
	"os"
	"path/filepath"
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
	picgo := config.GetConfig().PicgoPath
	err := CheckPicgoExecutable(picgo)
	if err != nil {
		t.Log(err)
	}
}

func TestUploadImages(t *testing.T) {
	config.InitConfig()
	rootPath := config.RootPath

	dir := filepath.Join(rootPath, "../", "test_images")

	files, err := os.ReadDir(dir)
	if err != nil {
		t.Log(err)
		return
	}

	images := make([]string, 0)
	for _, file := range files {
		name := file.Name()
		fullPath := filepath.Join(dir, name)
		//fmt.Println(fullPath)
		images = append(images, fullPath)
	}

	res, err := UploadImages(config.GetConfig().PicgoPath, images)
	if err != nil {
		t.Log(err)
		return
	}

	PrettyPrint(res)

}

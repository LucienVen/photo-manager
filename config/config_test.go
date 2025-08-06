package config

import (
	"github.com/LucienVen/photo-manager/utils"
	"testing"
)

func TestGetConfig(t *testing.T) {
	InitConfig()
	config := GetConfig()
	utils.PrettyPrint(config)

	t.Log(config.PicgoPath)
}

package controller

import (
	"github.com/LucienVen/photo-manager/config"
	"github.com/LucienVen/photo-manager/utils"
)

func CheckConfig() {
	cfg := config.GetConfig()
	utils.PrettyPrint(cfg)
}

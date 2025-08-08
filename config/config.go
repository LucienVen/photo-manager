package config

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"

	"github.com/joho/godotenv"
)

var (
	RootPath  string
	EnvConfig Config
)

type Config struct {
	GithubName string `env:"GITHUB_NAME"`
	GithubRepo string `env:"GITHUB_REPO"`
	CDNPrefix  string `env:"CDN_PREFIX"`
	ThumbWidth int    `env:"THUMB_WIDTH"`
	Debug      bool   `env:"DEBUG"`
	RecordDir  string `env:"RECORD_DIR"`
	PicgoPath  string `env:"PICGO_PATH"`
}

func InitConfig() {
	// 自动计算项目根路径（假设 config/loader.go 调用处）
	// 返回：repo/config/config.go
	_, currentFile, _, _ := runtime.Caller(0)

	RootPath = filepath.Join(filepath.Dir(currentFile), "../cmd")

	// 加载 .env
	envPath := filepath.Join(RootPath, ".env")
	if err := godotenv.Load(envPath); err != nil {
		log.Fatalf("Failed to load .env file from %s: %v", envPath, err)
	}

	picgoPath := os.Getenv("PICGO_PATH")
	if picgoPath == "" {
		picgoPath = "picgo" // 回退为 PATH 中的命令
	}
	// 绑定环境变量到结构体
	EnvConfig = Config{
		GithubName: os.Getenv("GITHUB_NAME"),
		GithubRepo: os.Getenv("GITHUB_REPO"),
		CDNPrefix:  os.Getenv("CDN_PREFIX"),
		ThumbWidth: getInt("THUMB_WIDTH", 320),
		Debug:      getBool("DEBUG", false),
		RecordDir:  os.Getenv("RECORD_DIR"),
		PicgoPath:  picgoPath,
	}

}

func getInt(key string, fallback int) int {
	if val := os.Getenv(key); val != "" {
		var i int
		fmt.Sscanf(val, "%d", &i)
		return i
	}
	return fallback
}

func getBool(key string, fallback bool) bool {
	if val := os.Getenv(key); val == "true" || val == "1" {
		return true
	}
	if val := os.Getenv(key); val == "false" || val == "0" {
		return false
	}
	return fallback
}

func GetConfig() Config {
	return EnvConfig
}

package utils

import (
	"encoding/json"
	"fmt"
	"github.com/LucienVen/photo-manager/config"
	"github.com/disintegration/imaging"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// 格式化打印结构体数据
func PrettyPrint(v interface{}) {
	data, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		fmt.Println("Failed to marshal config:", err)
		return
	}
	fmt.Println("当前数据:")
	fmt.Println(string(data))
}

// 获取图片基础信息
func GetImageInfo(path string) (width, height int, sizeKB int64, format string, err error) {
	file, err := os.Open(path)
	if err != nil {
		return
	}
	defer file.Close()

	// 获取文件大小
	stat, err := file.Stat()
	if err != nil {
		return
	}
	sizeKB = stat.Size() / 1024

	// 解码获取尺寸
	img, imgFormat, err := image.DecodeConfig(file)
	if err != nil {
		return
	}

	width = img.Width
	height = img.Height
	format = imgFormat
	return
}

// 重命名图片
func RenamePhoto(originalPath, hash string) (string, error) {

	// 获取原文件所在目录
	dir := filepath.Dir(originalPath)
	base := filepath.Base(originalPath)                  // e.g. "sunset.jpg"
	name := strings.TrimSuffix(base, filepath.Ext(base)) // "sunset"
	ext := filepath.Ext(base)                            // ".jpg"
	newFileName := name + "." + hash[:8] + ext           // "sunset.e4d909c2.jpg"

	newPath := filepath.Join(dir, newFileName)

	// 检查新文件名是否已存在
	if _, err := os.Stat(newPath); err == nil {
		return "", fmt.Errorf("文件已存在: %s", newPath)
	}

	// 重命名文件
	err := os.Rename(originalPath, newPath)
	if err != nil {
		return "", fmt.Errorf("重命名文件失败: %v", err)
	}

	return newPath, nil
}

// 生成略缩图，并存放与原图文件夹，并根据原图命名
func ResizePhoto(originalPath string, width int) (string, error) {
	newPath := GenNewFilePath(originalPath)
	// 读取原图
	src, err := imaging.Open(originalPath)
	if err != nil {
		return "", fmt.Errorf("打开图片失败: %v", err)
	}

	// 生成缩略图（例如宽度 300，按比例缩放）
	thumb := imaging.Resize(src, width, 0, imaging.Lanczos)

	// 保存到 thumbs 文件夹
	err = imaging.Save(thumb, newPath)
	if err != nil {
		return "", fmt.Errorf("保存缩略图失败: %v", err)
	}

	return newPath, nil
}

func GenNewFilePath(originalPath string) string {
	thumbTag := "thumb"
	// 获取原文件所在目录
	dir := filepath.Dir(originalPath)
	base := filepath.Base(originalPath)                  // e.g. "sunset.e4d909c2.jpg"
	name := strings.TrimSuffix(base, filepath.Ext(base)) // "sunset.e4d909c2"
	ext := filepath.Ext(base)                            // ".jpg"
	newFileName := name + "." + thumbTag + "." + ext     // "sunset.e4d909c2.thumb.jpg"

	newPath := filepath.Join(dir, newFileName)
	return newPath
}

func CleanTestFile(rootPath string) error {
	fmt.Println(rootPath)
	dir := filepath.Join(rootPath, "../", "test_images")

	files, err := os.ReadDir(dir)
	if err != nil {
		return fmt.Errorf("读取目录失败: %v", err)
	}

	for _, file := range files {
		name := file.Name()
		fullPath := filepath.Join(dir, name)

		// 删除所有 .thumb. 文件
		if strings.Contains(name, ".thumb.") {
			err := os.Remove(fullPath)
			if err != nil {
				return fmt.Errorf("删除缩略图失败 %s: %v", name, err)
			}
			continue
		}

		// 恢复 test1.<hash>.png -> test1.png
		if strings.HasPrefix(name, "test1.") &&
			strings.HasSuffix(name, ".png") &&
			!strings.Contains(name, ".thumb.") &&
			name != "test1.png" {
			original := filepath.Join(dir, "test1.png")

			// 删除原来的 test1.png（如果存在）
			_ = os.Remove(original)

			// 重命名回 test1.png
			err := os.Rename(fullPath, original)
			if err != nil {
				return fmt.Errorf("重命名回 test1.png 失败: %v", err)
			}
		}
	}

	return nil
}

// 验证 PicGo 是否可以执行
func CheckPicgoExecutable() {

	picgo := config.GetConfig().PicgoPath
	fmt.Println("picgo path: ", picgo)

	cmd1 := exec.Command(picgo, "--version")
	out1, err1 := cmd1.CombinedOutput()
	if err1 == nil {
		fmt.Println("PicGo (from $PATH) is executable:")
		fmt.Println(string(out1))
		return
	} else {
		fmt.Println("PicGo (from $PATH) failed:", err1)
	}
}

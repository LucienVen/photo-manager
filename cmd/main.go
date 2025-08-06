package main

import (
	"encoding/json"
	"fmt"
	"github.com/LucienVen/photo-manager/objects"
	"log"
	"os/exec"
	"path/filepath"

	"github.com/LucienVen/photo-manager/args"
	"github.com/LucienVen/photo-manager/config"
	"github.com/LucienVen/photo-manager/utils"
)

func main() {
	// 加载配置
	config.InitConfig()

	cfg := config.GetConfig()
	utils.PrettyPrint(cfg)

	// 解析命令行参数
	args, err := args.ParseArgs()
	if err != nil {
		log.Fatalf("参数解析失败: %v", err)
	}

	// 打印解析结果
	fmt.Printf("解析结果: %s\n", args.String())

	// TODO 环境检查

	photoHash, err := utils.GetFileSHA256(args.ImagePath)
	if err != nil {
		panic(fmt.Sprintf("生成 hash 失败: err:%v, path:%s", err, args.ImagePath))
	}

	fmt.Println("hash : ", photoHash)

	recordPath := filepath.Join(config.RootPath, "../", cfg.RecordDir)
	fmt.Println("recordPath: ", recordPath)

	// exist, _, err := utils.HashExists(photoHash, recordPath)
	// if err != nil {
	// 	panic(err.Error())
	// }

	// fmt.Println("exist: ", exist)

	width, height, sizeKB, format, err := utils.GetImageInfo(args.ImagePath)
	if err != nil {
		log.Fatalf("获取图片信息失败: %v", err)
	}
	fmt.Printf("图片信息：%dx%d, %d KB, 格式: %s\n", width, height, sizeKB, format)

	// 重命名图片
	photoPath, err := utils.RenamePhoto(args.ImagePath, photoHash)
	if err != nil {
		log.Fatalf("重命名图片失败: %v", err)
	}
	fmt.Printf("图片已重命名为: %s\n", photoPath)

	// 生成略缩图
	thumbPath, err := utils.ResizePhoto(photoPath, cfg.ThumbWidth)
	if err != nil {
		log.Fatalf("生成略缩图失败: %v", err)
	}

	fmt.Println("thumbPath: ", thumbPath)

	// 上传
	cmd := exec.Command("/Users/liangliangtoo/.nvm/versions/node/v24.5.0/bin/picgo", "upload", "--json", photoPath, thumbPath)

	output, err := cmd.Output()
	if err != nil {
		fmt.Printf("上传失败: %v\n", err)
		return
	}

	// 解析 JSON 返回值
	var results []objects.PicGoResult
	err = json.Unmarshal(output, &results)
	if err != nil {
		fmt.Printf("JSON解析失败: %v\n", err)
		fmt.Println("原始输出：", string(output))
		return
	}

	// 遍历每张上传结果
	for _, r := range results {
		fmt.Printf("文件：%s\n图片地址：%s\n\n", r.FileName, r.URL)
	}
}

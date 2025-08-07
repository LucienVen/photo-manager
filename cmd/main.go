package main

import (
	"fmt"
	"github.com/LucienVen/photo-manager/args"
	"github.com/LucienVen/photo-manager/config"
	"github.com/LucienVen/photo-manager/objects"
	"github.com/LucienVen/photo-manager/utils"
	"github.com/spf13/cobra"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type PhotoRecord struct {
	Filename   string   `json:"filename"`
	URL        string   `json:"url"`
	ThumbURL   string   `json:"thumb_url"`
	Path       string   `json:"path"`
	UploadedAt int64    `json:"uploaded_at"`
	Tags       []string `json:"tags"`
	Desc       string   `json:"desc"`
	SizeKB     int      `json:"size_kb"`
	Width      int      `json:"width"`
	Height     int      `json:"height"`
	Hash       string   `json:"hash"`
}

func main() {
	var desc string
	var tags string
	var imagePath string

	var rootCmd = &cobra.Command{
		Use:   "main [image path]",
		Short: "Image processing CLI",
		Args:  cobra.ExactArgs(1), // 要求必须提供一个 image path
		Run: func(cmd *cobra.Command, args []string) {
			imagePath = args[0]
			if _, err := os.Stat(imagePath); os.IsNotExist(err) {
				fmt.Printf("Error: file does not exist: %s\n", imagePath)
				os.Exit(1)
			}
		},
	}

	rootCmd.Flags().StringVarP(&desc, "desc", "d", "", "Image description")
	rootCmd.Flags().StringVarP(&tags, "tags", "t", "", "Comma-separated tags")

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Println(imagePath)

	if !args.IsImageFile(imagePath) {
		panic(fmt.Errorf("image file %s is not valid", imagePath))
	}

	argsInstance := args.NewArgs(imagePath, strings.Split(tags, ","), desc)
	utils.PrettyPrint(argsInstance)

	//return

	// 加载配置
	config.InitConfig()

	cfg := config.GetConfig()
	utils.PrettyPrint(cfg)

	// TODO 环境检查

	photoHash, err := utils.GetFileSHA256(argsInstance.ImagePath)
	if err != nil {
		panic(fmt.Sprintf("生成 hash 失败: err:%v, path:%s", err, argsInstance.ImagePath))
	}

	fmt.Println("hash : ", photoHash)

	recordPath := filepath.Join(config.RootPath, "../", cfg.RecordDir)
	fmt.Println("recordPath: ", recordPath)

	//exist, _, err := utils.HashExists(photoHash, recordPath)
	//if err != nil {
	//	panic(err.Error())
	//}
	//
	//fmt.Println("exist: ", exist)

	width, height, sizeKB, format, err := utils.GetImageInfo(argsInstance.ImagePath)
	if err != nil {
		log.Fatalf("获取图片信息失败: %v", err)
	}
	fmt.Printf("图片信息：%dx%d, %d KB, 格式: %s\n", width, height, sizeKB, format)

	// 重命名图片
	photoPath, err := utils.RenamePhoto(argsInstance.ImagePath, photoHash)
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
	picgo := config.GetConfig().PicgoPath

	err = utils.CheckPicgoExecutable(picgo)
	if err != nil {
		panic(err)
	}

	uploadRes, err := utils.UploadImages(picgo, []string{photoPath, thumbPath})
	if err != nil {
		panic(err)
	}

	// 生成记录，写入 records
	timenow := time.Now().Unix()
	for _, item := range uploadRes {

		record := PhotoRecord{
			Filename:   "",
			URL:        "",
			ThumbURL:   "",
			Path:       "",
			UploadedAt: timenow,
			Tags:       nil,
			Desc:       "",
			SizeKB:     0,
			Width:      0,
			Height:     0,
			Hash:       "",
		}

		// 追加到 record
	}

}

func (pr *PhotoRecord) GetHash() {

}
func (pr *PhotoRecord) SetHash() {

}

func (pr *PhotoRecord) SetImageInfo() {

}

func (pr *PhotoRecord) GetPath() {

}

func (pr *PhotoRecord) SetPath() {

}

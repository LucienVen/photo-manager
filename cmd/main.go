package main

import (
	"encoding/json"
	"fmt"
	"github.com/LucienVen/photo-manager/config"
	"github.com/LucienVen/photo-manager/logger"
	"github.com/LucienVen/photo-manager/utils"
	"github.com/spf13/cobra"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type PhotoRecord struct {
	Filename   string   `json:"filename"`
	URL        string   `json:"url"`        // 图片 URL
	Path       string   `json:"path"`       // 原始图片路径
	ThumbURL   string   `json:"thumb_url"`  // 略缩图 URL
	ThumbName  string   `json:"thumb_name"` // 略缩图文件名
	ThumbPath  string   `json:"thumb_path"` // 略缩图路径
	CreatedAt  int64    `json:"created_at"`
	Tags       []string `json:"tags"`
	Desc       string   `json:"desc"`
	SizeKB     int      `json:"size_kb"`
	Width      int      `json:"width"`
	Height     int      `json:"height"`
	Hash       string   `json:"hash"`
	RenamePath string   `json:"rename_path"` // 重命名后的路径
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

	logger.Debug("imagePath:%s", imagePath)

	// 加载配置
	config.InitConfig()
	cfg := config.GetConfig()
	if cfg.Debug {
		utils.PrettyPrint(cfg)
	}

	// 检查 PicGo 可执行文件
	picgo := config.GetConfig().PicgoPath
	err := utils.CheckPicgoExecutable(picgo)
	if err != nil {
		panic(err)
	}

	if !utils.IsImageFile(imagePath) {
		panic(fmt.Errorf("image file %s is not valid", imagePath))
	}

	recordPath := filepath.Join(config.RootPath, "../", cfg.RecordDir)
	logger.Debug("recordPath: %s", recordPath)

	imageInstance := PhotoRecord{
		Path:      imagePath,
		Tags:      strings.Split(tags, ","),
		Desc:      desc,
		CreatedAt: time.Now().Unix(),
	}

	err = imageInstance.SetHash()
	if err != nil {
		panic(err)
	}

	exist, existRecord, err := IsHashExistsInRecords(recordPath, imageInstance.GetHash())
	if err != nil {
		panic(err)
	}

	if exist {
		fmt.Printf("图片 hash [%s] 已经存在\n", imageInstance.GetHash())
		utils.PrettyPrint(existRecord)
		return
	}

	err = imageInstance.SetImageInfo()
	if err != nil {
		panic(err)
	}

	// 重命名图片
	err = imageInstance.Rename()
	if err != nil {
		panic(err)
	}

	// 生成略缩图
	err = imageInstance.ResizePhoto(cfg.ThumbWidth)
	if err != nil {
		panic(err)
	}

	uploadRes, err := utils.UploadImages(picgo, []string{imageInstance.RenamePath, imageInstance.ThumbPath})
	if err != nil {
		panic(err)
	}

	// 识别，并写入记录
	err = imageInstance.SetUploadUrl(uploadRes)
	if err != nil {
		panic(err)
	}

	// 生成记录，写入 records
	toRecordFileName := imageInstance.GetRecordFileName(recordPath)
	err = AppendRecordToFile(toRecordFileName, imageInstance)
	if err != nil {
		panic(err)
	}

	fmt.Printf(">>> 图片处理完成: %s\n", imageInstance.Filename)
	fmt.Printf("    - 记录已保存到: %s\n", toRecordFileName)
	fmt.Printf("    - 原图地址     : %s\n", imageInstance.URL)
	fmt.Printf("    - 缩略图地址   : %s\n", imageInstance.ThumbURL)
	return
}

func (pr *PhotoRecord) GetHash() string {
	return pr.Hash
}

func (pr *PhotoRecord) SetHash() error {
	photoHash, err := utils.GetFileSHA256(pr.Path)
	if err != nil {
		return fmt.Errorf("生成 hash 失败: err:%v, path:%s", err, pr.Path)
	}
	pr.Hash = photoHash
	return nil
}

func (pr *PhotoRecord) SetImageInfo() error {
	width, height, sizeKB, _, err := utils.GetImageInfo(pr.Path)
	if err != nil {
		return fmt.Errorf("获取图片信息失败: %v", err)
	}

	pr.Width = width
	pr.Height = height
	pr.SizeKB = int(sizeKB)

	return nil
}

func (pr *PhotoRecord) Rename() error {
	newPath, newFileName, err := utils.RenamePhoto(pr.Path, pr.Hash)
	if err != nil {
		return fmt.Errorf("图片重命名失败: %v", err)
	}
	pr.Filename = newFileName
	pr.RenamePath = newPath
	return nil
}

func (pr *PhotoRecord) ResizePhoto(width int) error {
	thumbPath, err := utils.ResizePhoto(pr.RenamePath, width)
	if err != nil {
		return fmt.Errorf("生成略缩图失败: %v", err)
	}

	pr.ThumbPath = thumbPath

	// 获取生成略缩图名称
	thumbName := filepath.Base(thumbPath)
	pr.ThumbName = thumbName
	return nil
}

func (pr *PhotoRecord) SetUploadUrl(res []string) error {
	if len(res) == 0 {
		return fmt.Errorf("上传结果为空")
	}

	for _, item := range res {
		// 判断是否略缩图
		if strings.Contains(item, "thumb") {
			pr.ThumbURL = item
		} else {
			pr.URL = item
		}
	}

	return nil
}

// 返回
func (pr *PhotoRecord) GetRecordFileName(recordDir string) string {
	t := time.Unix(pr.CreatedAt, 0)
	filename := fmt.Sprintf("%d-%02d.json", t.Year(), int(t.Month()))
	return filepath.Join(recordDir, filename)
}

// **************************************************

// AppendRecordToFile 读取 json 文件 -> 追加新记录 -> 写回
func AppendRecordToFile(filePath string, newRecord PhotoRecord) error {
	var records []PhotoRecord

	// 若文件存在，尝试读取原有内容
	if _, err := os.Stat(filePath); err == nil {
		data, err := os.ReadFile(filePath)
		if err != nil {
			return fmt.Errorf("读取文件失败: %w", err)
		}
		if len(data) > 0 {
			if err := json.Unmarshal(data, &records); err != nil {
				return fmt.Errorf("解析 JSON 失败: %w", err)
			}
		}
	}

	// 追加新记录
	records = append(records, newRecord)

	// 写回文件（带缩进）
	jsonData, err := json.MarshalIndent(records, "", "  ")
	if err != nil {
		return fmt.Errorf("JSON 序列化失败: %w", err)
	}

	// 确保目录存在
	if err := os.MkdirAll(filepath.Dir(filePath), os.ModePerm); err != nil {
		return fmt.Errorf("创建目录失败: %w", err)
	}

	if err := os.WriteFile(filePath, jsonData, 0644); err != nil {
		return fmt.Errorf("写入 JSON 文件失败: %w", err)
	}

	return nil
}

// ReadRecordsFromFile 读取所有记录
func ReadRecordsFromFile(filePath string) ([]PhotoRecord, error) {
	var records []PhotoRecord

	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("读取文件失败: %w", err)
	}
	if len(data) == 0 {
		return records, nil
	}
	if err := json.Unmarshal(data, &records); err != nil {
		return nil, fmt.Errorf("解析 JSON 失败: %w", err)
	}
	return records, nil
}

func ReadRecordsFromDir(dirPath string) ([]PhotoRecord, error) {
	var allRecords []PhotoRecord

	// 遍历目录下所有 json 文件
	files, err := os.ReadDir(dirPath)
	if err != nil {
		return nil, fmt.Errorf("读取目录失败: %w", err)
	}

	for _, file := range files {
		if file.IsDir() {
			continue
		}
		if filepath.Ext(file.Name()) != ".json" {
			continue
		}

		fullPath := filepath.Join(dirPath, file.Name())
		data, err := os.ReadFile(fullPath)
		if err != nil {
			return nil, fmt.Errorf("读取文件失败: %s: %w", fullPath, err)
		}

		var records []PhotoRecord
		if len(data) > 0 {
			if err := json.Unmarshal(data, &records); err != nil {
				return nil, fmt.Errorf("解析 JSON 失败: %s: %w", fullPath, err)
			}
			allRecords = append(allRecords, records...)
		}
	}

	return allRecords, nil
}

func IsHashExistsInRecords(dirPath string, targetHash string) (bool, PhotoRecord, error) {
	files, err := os.ReadDir(dirPath)
	if err != nil {
		return false, PhotoRecord{}, err
	}

	for _, file := range files {
		if file.IsDir() || !strings.HasSuffix(file.Name(), ".json") {
			continue
		}

		filePath := filepath.Join(dirPath, file.Name())
		f, err := os.Open(filePath)
		if err != nil {
			continue // 跳过无法读取的文件
		}

		decoder := json.NewDecoder(f)

		// 读取 [
		t, err := decoder.Token()
		if err != nil || t != json.Delim('[') {
			f.Close()
			continue
		}

		// 遍历数组中的每个元素
		for decoder.More() {
			var record PhotoRecord
			if err := decoder.Decode(&record); err != nil {
				continue
			}

			if record.Hash == targetHash {
				f.Close()
				return true, record, nil
			}
		}

		f.Close()
	}

	return false, PhotoRecord{}, nil
}

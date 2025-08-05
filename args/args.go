package args

import (
	"fmt"
	"os"
	"strings"
)

// Args 命令行参数结构体
type Args struct {
	ImagePath string   // 图片路径
	Tags      []string // 标签列表
	Desc      string   // 描述
}

// ParseArgs 解析命令行参数
func ParseArgs() (*Args, error) {
	args := os.Args[1:] // 跳过程序名

	if len(args) == 0 {
		return nil, fmt.Errorf("请提供图片路径")
	}

	// 第一个参数必须是图片路径
	imagePath := args[0]
	if !isImageFile(imagePath) {
		return nil, fmt.Errorf("第一个参数必须是图片文件路径")
	}

	result := &Args{
		ImagePath: imagePath,
		Tags:      []string{},
		Desc:      "",
	}

	// 解析剩余参数
	remainingArgs := args[1:]
	if len(remainingArgs) == 0 {
		return result, nil
	}

	// 智能解析参数
	return parseRemainingArgs(result, remainingArgs)
}

// parseRemainingArgs 智能解析剩余参数
func parseRemainingArgs(result *Args, args []string) (*Args, error) {
	for _, arg := range args {
		// 检查是否有明确的标签前缀
		if strings.HasPrefix(arg, "tags:") {
			tags := strings.TrimPrefix(arg, "tags:")
			if tags != "" {
				result.Tags = append(result.Tags, parseTags(tags)...)
			}
			continue
		}

		// 检查是否有明确的描述前缀
		if strings.HasPrefix(arg, "desc:") {
			desc := strings.TrimPrefix(arg, "desc:")
			if desc != "" {
				result.Desc = desc
			}
			continue
		}

		// 智能识别：如果参数包含逗号，认为是标签
		if strings.Contains(arg, ",") {
			result.Tags = append(result.Tags, parseTags(arg)...)
			continue
		}

		// 如果已经有描述，则认为是标签
		if result.Desc != "" {
			result.Tags = append(result.Tags, arg)
			continue
		}

		// 如果这是最后一个参数且没有描述，则认为是描述
		if len(args) == 1 || isLastArg(arg, args) {
			result.Desc = arg
		} else {
			// 否则认为是标签
			result.Tags = append(result.Tags, arg)
		}
	}

	return result, nil
}

// parseTags 解析标签字符串，支持逗号分隔
func parseTags(tagsStr string) []string {
	if tagsStr == "" {
		return []string{}
	}

	// 按逗号分割并清理空白
	tags := strings.Split(tagsStr, ",")
	var result []string

	for _, tag := range tags {
		tag = strings.TrimSpace(tag)
		if tag != "" {
			result = append(result, tag)
		}
	}

	return result
}

// isImageFile 检查是否为图片文件
func isImageFile(filename string) bool {
	ext := strings.ToLower(getFileExtension(filename))
	imageExts := []string{".jpg", ".jpeg", ".png", ".gif", ".bmp", ".webp", ".svg"}

	for _, ext2 := range imageExts {
		if ext == ext2 {
			return true
		}
	}

	return false
}

// getFileExtension 获取文件扩展名
func getFileExtension(filename string) string {
	lastDot := strings.LastIndex(filename, ".")
	if lastDot == -1 {
		return ""
	}
	return filename[lastDot:]
}

// isLastArg 检查是否为最后一个参数
func isLastArg(arg string, args []string) bool {
	return args[len(args)-1] == arg
}

// String 返回参数的字符串表示
func (a *Args) String() string {
	return fmt.Sprintf("图片路径: %s, 标签: %v, 描述: %s",
		a.ImagePath, a.Tags, a.Desc)
}

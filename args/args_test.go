package args

import (
	"os"
	"reflect"
	"testing"
)

func TestParseArgs(t *testing.T) {
	tests := []struct {
		name    string
		args    []string
		want    *Args
		wantErr bool
	}{
		{
			name: "基本图片路径",
			args: []string{"test.jpg"},
			want: &Args{
				ImagePath: "test.jpg",
				Tags:      []string{},
				Desc:      "",
			},
			wantErr: false,
		},
		{
			name: "带标签和描述",
			args: []string{"test.jpg", "tags:风景,广州塔,金融城", "desc:拍摄于广州"},
			want: &Args{
				ImagePath: "test.jpg",
				Tags:      []string{"风景", "广州塔", "金融城"},
				Desc:      "拍摄于广州",
			},
			wantErr: false,
		},
		{
			name: "智能识别标签和描述",
			args: []string{"test.jpg", "风景,广州塔,金融城", "拍摄于广州"},
			want: &Args{
				ImagePath: "test.jpg",
				Tags:      []string{"风景", "广州塔", "金融城"},
				Desc:      "拍摄于广州",
			},
			wantErr: false,
		},
		{
			name: "只有标签",
			args: []string{"test.jpg", "风景,广州塔,金融城"},
			want: &Args{
				ImagePath: "test.jpg",
				Tags:      []string{"风景", "广州塔", "金融城"},
				Desc:      "",
			},
			wantErr: false,
		},
		{
			name: "单个标签",
			args: []string{"test.jpg", "风景"},
			want: &Args{
				ImagePath: "test.jpg",
				Tags:      []string{"风景"},
				Desc:      "",
			},
			wantErr: false,
		},
		{
			name: "只有描述",
			args: []string{"test.jpg", "desc:拍摄于广州"},
			want: &Args{
				ImagePath: "test.jpg",
				Tags:      []string{},
				Desc:      "拍摄于广州",
			},
			wantErr: false,
		},
		{
			name:    "无参数",
			args:    []string{},
			want:    nil,
			wantErr: true,
		},
		{
			name:    "非图片文件",
			args:    []string{"test.txt"},
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 保存原始参数
			originalArgs := os.Args
			defer func() { os.Args = originalArgs }()

			// 设置测试参数
			os.Args = append([]string{"test"}, tt.args...)

			got, err := ParseArgs()
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseArgs() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ParseArgs() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParseTags(t *testing.T) {
	tests := []struct {
		name string
		tags string
		want []string
	}{
		{
			name: "空字符串",
			tags: "",
			want: []string{},
		},
		{
			name: "单个标签",
			tags: "风景",
			want: []string{"风景"},
		},
		{
			name: "多个标签",
			tags: "风景,广州塔,金融城",
			want: []string{"风景", "广州塔", "金融城"},
		},
		{
			name: "带空格的标签",
			tags: "风景, 广州塔 , 金融城",
			want: []string{"风景", "广州塔", "金融城"},
		},
		{
			name: "空标签",
			tags: "风景,,金融城",
			want: []string{"风景", "金融城"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := parseTags(tt.tags)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parseTags() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsImageFile(t *testing.T) {
	tests := []struct {
		name     string
		filename string
		want     bool
	}{
		{"JPG文件", "test.jpg", true},
		{"JPEG文件", "test.jpeg", true},
		{"PNG文件", "test.png", true},
		{"GIF文件", "test.gif", true},
		{"BMP文件", "test.bmp", true},
		{"WEBP文件", "test.webp", true},
		{"SVG文件", "test.svg", true},
		{"大写扩展名", "test.JPG", true},
		{"混合大小写", "test.PnG", true},
		{"非图片文件", "test.txt", false},
		{"无扩展名", "test", false},
		{"空字符串", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := isImageFile(tt.filename)
			if got != tt.want {
				t.Errorf("isImageFile() = %v, want %v", got, tt.want)
			}
		})
	}
}

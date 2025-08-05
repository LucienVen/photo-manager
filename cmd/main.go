package main

import (
	"fmt"
	"log"

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

}

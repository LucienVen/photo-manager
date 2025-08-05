package utils

import (
	"encoding/json"
	"fmt"
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

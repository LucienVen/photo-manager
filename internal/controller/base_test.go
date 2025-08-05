package controller

import (
	"testing"

	"github.com/LucienVen/photo-manager/config"
)

func TestCheckConfig(t *testing.T) {

	// remark: go test 默认不显示标准输出（stdout），除非测试失败或者你显式地使用了 -v 参数（verbose 模式）。

	config.InitConfig()
	CheckConfig()
}

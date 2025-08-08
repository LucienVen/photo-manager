package logger

import (
	"fmt"
	"github.com/LucienVen/photo-manager/config"
)

// Debug 打印调试信息，只有在 DEBUG=true 时才打印
func Debug(format string, a ...any) {
	if config.GetConfig().Debug {
		fmt.Printf("[DEBUG] "+format+"\n", a...)
	}
}

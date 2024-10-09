package utils

import (
	"os"
	"path/filepath"
)

// 检查权限
func IsAdmin() bool {
	fd, err := os.Open("\\\\.\\PHYSICALDRIVE0")
	if err != nil {
		return false
	}
	defer fd.Close()
	return true
}

// 设置工作目录
func InitWorkPath() {
	args := os.Args
	if len(args) > 1 && args[1] == "--enable-workspace" {
		exePath, err := os.Executable()
		if err != nil {
			return
		}

		exeDir := filepath.Dir(exePath)

		err = os.Chdir(exeDir)
		if err != nil {
			return
		}
	}
}

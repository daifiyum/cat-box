package utils

import (
	"os"
	"path/filepath"
)

func AppInit() {
	InitWorkPath()
	InitLog()
	InitDPI()
	InitToast()
}

// 设置工作目录
func InitWorkPath() {
	args := os.Args
	if len(args) >= 2 {
		if args[1] == "--enable-workspace" {
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
}

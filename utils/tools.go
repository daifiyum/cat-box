package utils

import (
	"os"
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

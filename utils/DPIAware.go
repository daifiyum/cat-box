package utils

import (
	"fmt"
	"syscall"
)

var (
	moduser32              = syscall.NewLazyDLL("user32.dll")
	procSetProcessDPIAware = moduser32.NewProc("SetProcessDPIAware")
)

func setProcessDPIAware() error {
	status, r, err := procSetProcessDPIAware.Call()
	if status == 0 {
		return fmt.Errorf("SetProcessDPIAware failed %d: %v %v", status, r, err)
	}
	return nil
}

func init() {
	err := setProcessDPIAware()
	if err != nil {
		LogError(err.Error())
	}
}

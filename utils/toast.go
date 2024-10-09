package utils

import (
	"fmt"
	"path/filepath"
	"syscall"
	"unsafe"
)

var (
	toast        = syscall.NewLazyDLL("./resources/libs/toast.dll")
	toastProc    = toast.NewProc("toast")
	registerProc = toast.NewProc("register_toast_notification")
)

func ShowToast(aumid, title, message string) error {
	cAumid, err := syscall.BytePtrFromString(aumid)
	if err != nil {
		return fmt.Errorf("error converting title to C string: %v", err)
	}
	cTitle, err := syscall.BytePtrFromString(title)
	if err != nil {
		return fmt.Errorf("error converting title to C string: %v", err)
	}
	cMessage, err := syscall.BytePtrFromString(message)
	if err != nil {
		return fmt.Errorf("error converting message to C string: %v", err)
	}

	ret, _, _ := toastProc.Call(uintptr(unsafe.Pointer(cAumid)), uintptr(unsafe.Pointer(cTitle)), uintptr(unsafe.Pointer(cMessage)))
	if ret == 0 {
		return fmt.Errorf("failed to show toast")
	}

	return nil
}

// HKEY_CURRENT_USER\Software\Classes\AppUserModelId\<id>
func RegisterToast(id, name, path string) error {
	cId, err := syscall.BytePtrFromString(id)
	if err != nil {
		return fmt.Errorf("error converting title to C string: %v", err)
	}
	cName, err := syscall.BytePtrFromString(name)
	if err != nil {
		return fmt.Errorf("error converting message to C string: %v", err)
	}
	cPath, err := syscall.BytePtrFromString(path)
	if err != nil {
		return fmt.Errorf("error converting message to C string: %v", err)
	}

	ret, _, _ := registerProc.Call(uintptr(unsafe.Pointer(cId)), uintptr(unsafe.Pointer(cName)), uintptr(unsafe.Pointer(cPath)))
	if ret == 0 {
		return fmt.Errorf("failed to register toast")
	}

	return nil
}

// APPID 注册, 以显示通知图标和名称
func InitToast() {
	toastIcon, _ := filepath.Abs("./resources/icons/box.ico")
	RegisterToast("cat-box", "cat-box", toastIcon)
}

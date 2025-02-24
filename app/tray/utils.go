package tray

import (
	"golang.org/x/sys/windows"
)

// 打开默认浏览器
func OpenBrowser(url string) {
	windows.ShellExecute(0, nil, windows.StringToUTF16Ptr(url), nil, nil, windows.SW_SHOWNORMAL)
}

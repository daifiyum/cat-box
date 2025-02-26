package app

import (
	"fmt"
	"path/filepath"
	"runtime"
	"syscall"
	"unsafe"

	W "github.com/daifiyum/cat-box/app/windows"
)

func init() {
	// windows GUI 编程需要在os线程上
	runtime.LockOSThread()
}

type App struct {
	windowHandle    syscall.Handle
	instanceHandle  syscall.Handle
	iconResource    string
	tooltipText     string
	windowClassName *uint16
	windowTitle     *uint16
	nid             W.NOTIFYICONDATAW
	trayMenu        *W.Menu
	onLeftClick     func()
	onInitialize    func()
}

// 初始化
// 参数：t 提示文字，i 图标路径
func New(t, i string) *App {
	return &App{
		tooltipText:  t,
		iconResource: i,
		onLeftClick:  func() {},
		onInitialize: func() {},
		trayMenu: &W.Menu{
			Handle:    0,
			Callbacks: make(map[uint32]func()),
		},
	}
}

// Run
func (t *App) Run() error {
	if err := t.SetAumid(); err != nil {
		return err
	}

	if err := t.setProcessDPIAware(); err != nil {
		return err
	}

	if err := t.registerWindowClass(); err != nil {
		return err
	}

	if err := t.createWindow(); err != nil {
		return err
	}

	if err := t.addTrayIcon(); err != nil {
		return err
	}

	if err := t.messageLoop(); err != nil {
		return err
	}

	return nil
}

// 设置 AUMID
func (t *App) SetAumid() error {
	iconURL, err := filepath.Abs(t.iconResource)
	if err != nil {
		return fmt.Errorf("failed to get absolute path for icon: %w", err)
	}

	err = W.RegisterAUMID(t.tooltipText, t.tooltipText, iconURL)
	if err != nil {
		return fmt.Errorf("failed to register AUMID: %w", err)
	}

	err = W.SetAUMID(t.tooltipText)
	if err != nil {
		return fmt.Errorf("failed to set AUMID: %w", err)
	}

	return nil
}

// 消息循环
// https://learn.microsoft.com/zh-cn/windows/win32/api/winuser/nf-winuser-getmessage
func (t *App) messageLoop() error {
	var msg W.MSG
	for {
		ret, _, err := W.GetMessage.Call(uintptr(unsafe.Pointer(&msg)), 0, 0, 0)

		switch int32(ret) {
		case -1:
			return fmt.Errorf("GetMessage failed with error: %d", err)
		case 0:
			return nil
		default:
			W.TranslateMessage.Call(uintptr(unsafe.Pointer(&msg)))
			W.DispatchMessage.Call(uintptr(unsafe.Pointer(&msg)))
		}
	}
}

// 设置进程DPI感知
func (t *App) setProcessDPIAware() error {
	status, _, err := W.SetProcessDPIAware.Call()
	if status == 0 {
		return fmt.Errorf("SetProcessDPIAware failed: %v", err)
	}
	return nil
}

// 注册窗口类
func (t *App) registerWindowClass() error {
	hin, _, _ := W.GetModuleHandle.Call(0)

	hIcon, _ := W.LoadIconFromFile(t.iconResource)
	cursor, _, _ := W.LoadCursor.Call(0, uintptr(W.IDC_ARROW))

	wcn, _ := syscall.UTF16PtrFromString("gowintray")
	t.windowClassName = wcn
	wn, _ := syscall.UTF16PtrFromString("gowintray")
	t.windowTitle = wn

	var wcex W.WNDCLASSEX
	wcex.Size = uint32(unsafe.Sizeof(wcex))
	wcex.Style = W.CS_HREDRAW | W.CS_VREDRAW
	wcex.WndProc = syscall.NewCallback(t.windowProc)
	wcex.ClsExtra = 0
	wcex.WndExtra = 0
	wcex.Instance = syscall.Handle(hin)
	wcex.Icon = syscall.Handle(hIcon)
	wcex.Cursor = syscall.Handle(cursor)
	wcex.Background = syscall.Handle(W.COLOR_WINDOW + 1)
	wcex.MenuName = nil
	wcex.ClassName = t.windowClassName
	wcex.IconSm = 0

	ret, _, err := W.RegisterClassEx.Call(uintptr(unsafe.Pointer(&wcex)))
	if ret == 0 {
		return fmt.Errorf("RegisterClassEx failed: %w", err)
	}

	t.instanceHandle = syscall.Handle(hin)
	return nil
}

// 创建窗口
func (t *App) createWindow() error {
	hw, _, err := W.CreateWindowEx.Call(
		0,
		uintptr(unsafe.Pointer(t.windowClassName)),
		uintptr(unsafe.Pointer(t.windowTitle)),
		uintptr(W.WS_OVERLAPPEDWINDOW),
		uintptr(W.CW_USEDEFAULT),
		0,
		uintptr(W.CW_USEDEFAULT),
		0,
		0,
		0,
		uintptr(t.instanceHandle),
		0,
	)
	if hw == 0 {
		return fmt.Errorf("CreateWindowEx failed: %w", err)
	}

	t.windowHandle = syscall.Handle(hw)

	W.ShowWindow.Call(
		uintptr(hw),
		uintptr(0), // 0 隐藏窗口，1 显示
	)

	W.UpdateWindow.Call(
		uintptr(hw),
	)

	return nil
}

// 添加化托盘
func (t *App) addTrayIcon() error {
	hIcon, _ := W.LoadIconFromFile(t.iconResource)
	var nid W.NOTIFYICONDATAW
	nid.CbSize = uint32(unsafe.Sizeof(nid))
	nid.HWnd = t.windowHandle
	nid.UID = 1
	nid.UFlags = W.NIF_ICON | W.NIF_TIP | W.NIF_MESSAGE
	nid.HIcon = hIcon
	nid.UCallbackMessage = W.WM_TRAY_NOTIFYICON
	nid.SzTip = W.TipFromStr(t.tooltipText)

	t.nid = nid

	ret, _, err := W.ShellNotifyIcon.Call(W.NIM_ADD, uintptr(unsafe.Pointer(&t.nid)))
	if ret == 0 {
		return fmt.Errorf("failed to add tray icon: %w", err)
	}

	return nil
}

// 设置托盘提示
func (t *App) SetToolTip(s string) error {
	t.nid.SzTip = W.TipFromStr(s)
	ret, _, err := W.ShellNotifyIcon.Call(W.NIM_MODIFY, uintptr(unsafe.Pointer(&t.nid)))
	if ret == 0 {
		return fmt.Errorf("failed to update tray icon: %w", err)
	}
	return nil
}

// 更新托盘图标
func (t *App) SetIcon(iconPath string) error {
	hIcon, err := W.LoadIconFromFile(iconPath)
	if err != nil {
		return fmt.Errorf("failed to load icon: %w", err)
	}

	t.nid.HIcon = hIcon
	ret, _, err := W.ShellNotifyIcon.Call(W.NIM_MODIFY, uintptr(unsafe.Pointer(&t.nid)))
	if ret == 0 {
		return fmt.Errorf("failed to update tray icon: %w", err)
	}

	return nil
}

// 弹出气泡通知
func (t *App) ShowTrayNotification(title, msg string) error {
	var notify W.NOTIFYICONDATAW
	notify.CbSize = uint32(unsafe.Sizeof(notify))
	notify.UID = 1
	notify.HWnd = t.windowHandle
	notify.UFlags = W.NIF_INFO

	W.SetUTF16String(&notify.SzInfoTitle, title)
	W.SetUTF16String(&notify.SzInfo, msg)

	ret, _, err := W.ShellNotifyIcon.Call(W.NIM_MODIFY, uintptr(unsafe.Pointer(&notify)))
	if ret == 0 {
		return fmt.Errorf("Shell_NotifyIcon failed: %w", err)
	}

	return nil
}

// 弹出托盘菜单
func (t *App) showContextMenu() {
	pt := W.POINT{}
	W.GetCursorPos.Call(uintptr(unsafe.Pointer(&pt)))
	x, y := int(pt.X), int(pt.Y)

	W.SetForegroundWindow.Call(uintptr(t.windowHandle))

	W.TrackPopupMenu.Call(
		uintptr(t.trayMenu.Handle),
		uintptr(W.TPM_LEFTBUTTON),
		uintptr(x),
		uintptr(y),
		0,
		uintptr(t.windowHandle),
		0,
	)
}

func (t *App) SetMenu(m *W.Menu) {
	t.trayMenu = m
}

func (t *App) Quit() {
	W.DestroyWindow.Call(uintptr(t.windowHandle))
}

func (t *App) SetOnLeftClick(fn func()) {
	t.onLeftClick = fn
}

func (t *App) SetOnInitialize(fn func()) {
	t.onInitialize = fn
}

func (t *App) cleanup() {
	W.ShellNotifyIcon.Call(W.NIM_DELETE, uintptr(unsafe.Pointer(&t.nid)))

	W.UnregisterClass.Call(
		uintptr(unsafe.Pointer(t.windowClassName)),
		uintptr(t.instanceHandle),
	)
}

// 消息处理函数
func (t *App) windowProc(hwnd syscall.Handle, msg uint32, wparam, lparam uintptr) uintptr {
	switch msg {
	case W.WM_CREATE:
		t.onInitialize()
		return 0
	case W.WM_TRAY_NOTIFYICON:
		// 托盘左右键
		switch lparam {
		case W.WM_LBUTTONUP:
			t.onLeftClick()
		case W.WM_RBUTTONUP:
			t.showContextMenu()
		}
		return 0
	case W.WM_COMMAND:
		// 菜单项回调
		if W.HIWORD(uint64(wparam)) == 0 {
			if callback, exists := t.trayMenu.Callbacks[W.LOWORD(uint64(wparam))]; exists {
				callback()
			}
		}
		return 0
	case W.WM_CLOSE:
		W.DestroyWindow.Call(uintptr(t.windowHandle))
		return 0
	case W.WM_DESTROY:
		t.cleanup()
		W.PostQuitMessage.Call(0)
		return 0
	default:
		ret, _, _ := W.DefWindowProc.Call(uintptr(hwnd), uintptr(msg), wparam, lparam)
		return ret
	}
}

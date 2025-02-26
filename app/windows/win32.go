package windows

import (
	"syscall"
)

var (
	// kernel32
	kernel32        = syscall.MustLoadDLL("kernel32")
	GetModuleHandle = kernel32.MustFindProc("GetModuleHandleW")

	// user32
	user32              = syscall.MustLoadDLL("user32.dll")
	LoadIcon            = user32.MustFindProc("LoadIconW")
	LoadCursor          = user32.MustFindProc("LoadCursorW")
	RegisterClassEx     = user32.MustFindProc("RegisterClassExW")
	CreateWindowEx      = user32.MustFindProc("CreateWindowExW")
	ShowWindow          = user32.MustFindProc("ShowWindow")
	UpdateWindow        = user32.MustFindProc("UpdateWindow")
	GetMessage          = user32.MustFindProc("GetMessageW")
	TranslateMessage    = user32.MustFindProc("TranslateMessage")
	DispatchMessage     = user32.MustFindProc("DispatchMessageW")
	PostQuitMessage     = user32.MustFindProc("PostQuitMessage")
	DestroyWindow       = user32.MustFindProc("DestroyWindow")
	DefWindowProc       = user32.MustFindProc("DefWindowProcW")
	LoadImage           = user32.MustFindProc("LoadImageW")
	CreatePopupMenu     = user32.MustFindProc("CreatePopupMenu")
	AppendMenu          = user32.MustFindProc("AppendMenuW")
	GetCursorPos        = user32.MustFindProc("GetCursorPos")
	TrackPopupMenu      = user32.MustFindProc("TrackPopupMenu")
	SetForegroundWindow = user32.MustFindProc("SetForegroundWindow")
	UnregisterClass     = user32.MustFindProc("UnregisterClassW")
	CheckMenuItem       = user32.MustFindProc("CheckMenuItem")
	SetProcessDPIAware  = user32.MustFindProc("SetProcessDPIAware")

	// shell32
	shell32                                 = syscall.MustLoadDLL("shell32.dll")
	ShellNotifyIcon                         = shell32.MustFindProc("Shell_NotifyIconW")
	SetCurrentProcessExplicitAppUserModelID = shell32.MustFindProc("SetCurrentProcessExplicitAppUserModelID")

	// advapi32
	advapi32                 = syscall.MustLoadDLL("advapi32.dll")
	CheckTokenMembership     = advapi32.MustFindProc("CheckTokenMembership")
	AllocateAndInitializeSid = advapi32.MustFindProc("AllocateAndInitializeSid")
	FreeSid                  = advapi32.MustFindProc("FreeSid")
)

// 常量标识
const (
	// 窗口
	IDC_ARROW           = 32512 // 正常选择光标
	COLOR_WINDOW        = 5
	CW_USEDEFAULT       = 0x80000000 // 窗口x、y、w、h
	WS_OVERLAPPED       = 0x00000000
	WS_CAPTION          = 0x00C00000 // WS_BORDER | WS_DLGFRAME
	WS_SYSMENU          = 0x00080000
	WS_THICKFRAME       = 0x00040000
	WS_MINIMIZEBOX      = 0x00020000
	WS_MAXIMIZEBOX      = 0x00010000
	WS_OVERLAPPEDWINDOW = (WS_OVERLAPPED | WS_CAPTION | WS_SYSMENU | WS_THICKFRAME | WS_MINIMIZEBOX | WS_MAXIMIZEBOX)
	CS_HREDRAW          = 0x0002
	CS_VREDRAW          = 0x0001

	// 从本地加载图片
	IMAGE_ICON      = 1
	LR_LOADFROMFILE = 0x00000010
	LR_DEFAULTSIZE  = 0x00000040

	// 托盘
	NIM_ADD     = 0x00000000
	NIM_DELETE  = 0x00000002
	NIM_MODIFY  = 0x00000001
	NIF_ICON    = 0x00000002
	NIF_MESSAGE = 0x00000001
	NIF_TIP     = 0x00000004
	NIF_INFO    = 0x00000010
	NIIF_INFO   = 0x00000001
	NIF_STATE   = 0x00000008

	// 消息
	WM_COMMAND         = 0x0111
	WM_LBUTTONUP       = 0x0202
	WM_RBUTTONUP       = 0x0205
	WM_USER            = 0x0400
	WM_TRAY_NOTIFYICON = WM_USER + 1 // 自定义托盘消息
	WM_CLOSE           = 0x0010
	WM_CREATE          = 0x0001
	WM_SYSCOMMAND      = 0x0112 // 当用户选择最大化按钮、最小化按钮、还原按钮或关闭按钮时，窗口会收到此消息
	WM_DESTROY         = 0x0002 // 在窗口被销毁时发送

	// 菜单
	TPM_RETURNCMD  = 0x0100
	TPM_NONOTIFY   = 0x0080
	TPM_LEFTBUTTON = 0x0000
	MF_SEPARATOR   = 0x00000800
	MF_STRING      = 0x00000000
	MF_CHECKED     = 0x00000008
	MF_UNCHECKED   = 0x00000000
	MF_BYCOMMAND   = 0x000
	MF_POPUP       = 0x00000010

	// 安全标识SID
	SECURITY_BUILTIN_DOMAIN_RID = 0x00000020
	DOMAIN_ALIAS_RID_ADMINS     = 0x00000220
)

// WNDCLASSEX 结构体
type WNDCLASSEX struct {
	Size       uint32
	Style      uint32
	WndProc    uintptr
	ClsExtra   int32
	WndExtra   int32
	Instance   syscall.Handle
	Icon       syscall.Handle
	Cursor     syscall.Handle
	Background syscall.Handle
	MenuName   *uint16
	ClassName  *uint16
	IconSm     syscall.Handle
}

// 鼠标坐标
type POINT struct {
	X, Y int32
}

// MSG 结构体
type MSG struct {
	Hwnd    syscall.Handle
	Message uint32
	WParam  uintptr
	LParam  uintptr
	Time    uint32
	Pt      POINT
}

// NOTIFYICONDATAW 结构体
// UTimeout 仅在 Windows 2000 和 Windows XP 中有效
type NOTIFYICONDATAW struct {
	CbSize           uint32
	HWnd             syscall.Handle
	UID              uint32
	UFlags           uint32
	UCallbackMessage uint32
	HIcon            syscall.Handle
	SzTip            [128]uint16
	DwState          uint32
	DwStateMask      uint32
	SzInfo           [256]uint16
	UVersion         uint32
	SzInfoTitle      [64]uint16
	DwInfoFlags      uint32
	GuidItem         syscall.GUID
	HBalloonIcon     syscall.Handle
}

// SID_IDENTIFIER_AUTHORITY 结构体
type SID_IDENTIFIER_AUTHORITY struct {
	Value [6]byte
}

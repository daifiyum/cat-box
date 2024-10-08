package tray

import (
	"os/exec"
	"syscall"

	"github.com/daifiyum/cat-box/singbox"
	"github.com/daifiyum/cat-box/utils"

	"github.com/energye/systray"
)

var (
	moduser32              = syscall.NewLazyDLL("user32.dll")
	procSetProcessDPIAware = moduser32.NewProc("SetProcessDPIAware")
)

// 高分辨率显示
func SetProcessDPIAware() {
	procSetProcessDPIAware.Call()
}

func RunTray() {
	SetProcessDPIAware()
	InitIcons()
	CreateTray()
	CreateMenu()
}

func OpenBrowser(url string) {
	cmd := exec.Command("cmd", "/c", "start", url)
	cmd.SysProcAttr = &syscall.SysProcAttr{
		HideWindow: true,
	}
	cmd.Run()
}

func CreateTray() {
	systray.SetIcon(AppIcon)
	systray.SetTitle("cat-box")
	systray.SetTooltip("cat-box")
	systray.SetOnClick(func(menu systray.IMenu) {
		if utils.IsProxy {
			err := singbox.Stop()
			if err != nil {
				return
			}
			systray.SetIcon(AppIcon)
		} else {
			if utils.IsTun {
				if !utils.IsAdmin() {
					utils.ShowToast("cat-box", "警告", "开启Tun模式需以管理员模式运行！")
					return
				}
			}
			err := singbox.Start()
			if err != nil {
				return
			}
			systray.SetIcon(ProxyIcon)

		}
		singbox.CheckCoreStatus()
	})
}

func CreateMenu() {
	// menu
	mHome := systray.AddMenuItem("面板", "打开代理面板")
	mHome.SetIcon(HomeIcon)

	mSub := systray.AddMenuItem("订阅", "打开订阅面板")
	mSub.SetIcon(SubIcon)

	systray.AddSeparator()

	mSysProxy := systray.AddMenuItemCheckbox("系统代理", "System Proxy", true)
	mTunMode := systray.AddMenuItemCheckbox("TUN模式", "TUN Mode", false)

	systray.AddSeparator()

	mQuit := systray.AddMenuItem("退出", "Quit the whole app")
	mQuit.SetIcon(CloseIcon)

	// click
	mHome.Click(func() {
		OpenBrowser("http://localhost:9090/ui")
	})

	mSub.Click(func() {
		OpenBrowser("http://localhost:3000")
	})

	mSysProxy.Click(func() {
		if mTunMode.Checked() {
			mTunMode.Uncheck()
			mSysProxy.Check()
			utils.IsTun = false
		}
	})

	mTunMode.Click(func() {
		if mSysProxy.Checked() {
			mSysProxy.Uncheck()
			mTunMode.Check()
			utils.IsTun = true
		}
	})

	mQuit.Click(func() {
		systray.Quit()
	})
}

package tray

import (
	"log"

	"github.com/daifiyum/cat-box/app"
	U "github.com/daifiyum/cat-box/config"
	S "github.com/daifiyum/cat-box/sing-box"
)

// 监控托盘
func Watcher(app *app.App) {
	U.IsCoreRunning.Watch(func(newValue bool) {
		if newValue {
			if U.IsTun.Get() {
				app.SetIcon("./resources/icons/tun.ico")
				app.SetToolTip("核心运行中: Tun")
			} else {
				app.SetIcon("./resources/icons/proxy.ico")
				app.SetToolTip("核心运行中: 系统代理")
			}
			return
		}
		app.SetIcon("./resources/icons/box.ico")
		app.SetToolTip("cat-box")
	})

	U.IsTun.Watch(func(newValue bool) {
		err := S.SwitchProxyMode(newValue)
		if err != nil {
			log.Println(err)
		}
	})
}

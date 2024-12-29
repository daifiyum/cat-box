package tray

import (
	"log"

	"github.com/daifiyum/cat-box/app"
	U "github.com/daifiyum/cat-box/config"
	"github.com/daifiyum/cat-box/singbox"
)

// 监控托盘
func Watcher(app *app.App) {
	U.IsCoreRunning.Watch(func(newValue bool) {
		if newValue {
			app.SetIcon("./resources/icons/proxy.ico")
			app.SetToolTip("核心运行中")
			return
		}
		app.SetIcon("./resources/icons/box.ico")
		app.SetToolTip("cat-box")
	})

	U.IsTun.Watch(func(newValue bool) {
		err := singbox.SwitchProxyMode(newValue)
		if err != nil {
			log.Println(err)
		}
	})
}

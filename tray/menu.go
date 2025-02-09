package tray

import (
	"fmt"

	"github.com/daifiyum/cat-box/app"
	W "github.com/daifiyum/cat-box/app/windows"
	U "github.com/daifiyum/cat-box/config"
	S "github.com/daifiyum/cat-box/sing-box"
)

func Menu(app *app.App) *W.Menu {
	menu := W.NewMenu()
	menu.AddItem(1, "面板", func() {
		if !U.IsCoreRunning.Get() {
			app.ShowTrayNotification("警告", "核心未启动，无法访问")
			return
		}
		OpenBrowser(fmt.Sprintf("http://localhost:%s/ui", U.Box.ClashAPIPort))
	})
	menu.AddItem(2, "订阅", func() {
		OpenBrowser(fmt.Sprintf("http://localhost:%s", U.Port))
	})
	menu.AddCheckMenu(3, "TUN模式", false, func() {
		isAdmin, _ := W.IsUserAdmin()
		if !isAdmin {
			app.ShowTrayNotification("警告", "开启Tun模式需以管理员模式运行！")
			return
		}
		if menu.ToggleCheck(3) {
			U.IsTun.Set(true)
		} else {
			U.IsTun.Set(false)
		}
	})
	menu.AddSeparator()
	menu.AddItem(4, "退出", func() {
		S.Stop()
		app.Quit()
	})

	return menu
}

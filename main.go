package main

import (
	"github.com/daifiyum/cat-box/singbox"
	"github.com/daifiyum/cat-box/subservice"
	"github.com/daifiyum/cat-box/tray"
	"github.com/daifiyum/cat-box/utils"
	"github.com/energye/systray"
)

func main() {
	utils.AppInit()
	systray.Run(onReady, onExit)
}

func onReady() {
	tray.InitTray()
	subservice.SubService()
}

func onExit() {
	singbox.Stop()
}

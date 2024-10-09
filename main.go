package main

import (
	"github.com/daifiyum/cat-box/singbox"
	"github.com/daifiyum/cat-box/subservice"
	"github.com/daifiyum/cat-box/tray"
	"github.com/daifiyum/cat-box/utils"
	"github.com/energye/systray"
)

func main() {
	systray.Run(onReady, onExit)
}

func onReady() {
	utils.AppInit()
	tray.RunTray()
	subservice.SubService()
}

func onExit() {
	singbox.Stop()
}

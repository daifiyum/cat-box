package main

import (
	"log"

	"github.com/daifiyum/cat-box/api"
	"github.com/daifiyum/cat-box/app"
	"github.com/daifiyum/cat-box/app/tray"
	"github.com/daifiyum/cat-box/database"
	I "github.com/daifiyum/cat-box/initializer"
	S "github.com/daifiyum/cat-box/singbox"
	"github.com/daifiyum/cat-box/tasks"
)

func init() {
	I.Initialize()
}

func main() {
	app := app.New("cat-box", "./resources/icons/box.ico")

	app.Click(func() {
		S.SwitchCore()
	})

	menu := tray.Menu(app)
	app.SetMenu(menu)

	app.Ready(func() {
		database.Init()
		tasks.Run()
		tray.Watcher(app)
		go api.Run()
	})

	if err := app.Run(); err != nil {
		log.Fatal(err)
	}
}

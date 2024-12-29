package watcher

import (
	"fmt"

	"github.com/daifiyum/cat-box/database/models"
	"github.com/daifiyum/cat-box/tasks"
)

func Setting(settings []models.Setting) {
	for _, setting := range settings {
		switch setting.Key {
		case "update_delay":
			tasks.UpdateDelay(setting.Value)
		default:
			fmt.Printf("Handling update for key: %s, value: %s\n", setting.Key, setting.Value)
		}
	}
}

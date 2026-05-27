package watcher

import (
	"fmt"

	U "github.com/daifiyum/cat-box/common"
	"github.com/daifiyum/cat-box/database/models"
	"github.com/daifiyum/cat-box/tasks"
)

func Setting(setting models.Setting) {
	switch setting.Label {
	case "update_interval":
		tasks.UpdateDelay(setting.Value)
	case "custom_user_agent":
		U.UserAgent = setting.Value
	default:
		fmt.Printf("Handling update for key: %s, value: %s\n", setting.Label, setting.Value)
	}
}

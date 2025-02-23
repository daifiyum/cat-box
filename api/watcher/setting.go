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
	case "user_agent_type":
		U.DefaultUserAgent = setting.Value
	default:
		fmt.Printf("Handling update for key: %s, value: %s\n", setting.Label, setting.Value)
	}
}

package task

import (
	"github.com/daifiyum/cat-box/converter"
	"github.com/daifiyum/cat-box/singbox"
	"github.com/daifiyum/cat-box/subservice/database"
	"github.com/daifiyum/cat-box/subservice/models"
	"github.com/daifiyum/cat-box/utils"

	"github.com/robfig/cron/v3"
)

var scheduler *cron.Cron

func Scheduler() {
	if scheduler != nil {
		scheduler.Stop()
	}
	db := database.DB
	setting := new(models.Setting)
	db.Where("key = ?", "update_delay").First(setting)
	scheduler = cron.New()
	scheduler.AddFunc("@every "+setting.Value, func() {
		handleUpdate()
	})
	scheduler.Start()
}

func handleUpdate() {
	db := database.DB
	var subscriptions []models.Subscriptions
	db.Find(&subscriptions)
	for _, subscription := range subscriptions {
		if subscription.AutoUpdate {
			config, err := converter.Handler(subscription.Link)
			if err != nil {
				utils.LogError("Failed to generate configuration")
				continue
			}
			db.Model(&subscription).Where(subscription.ID).Update("data", config)
			if subscription.Active {
				if utils.IsProxy {
					err = singbox.Start()
					if err != nil {
						utils.LogError("Failed to reload configuration")
						continue
					}
				}
			}
			utils.LogInfo("Automatic update successful")
		}
	}
}

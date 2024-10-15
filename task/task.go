package task

import (
	"time"

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
		SubUpdate()
	})
	scheduler.Start()
}

func SubUpdate() {
	db := database.DB
	var subscriptions []models.Subscriptions
	db.Find(&subscriptions)
	var updates []models.Subscriptions
	for _, subscription := range subscriptions {
		if subscription.AutoUpdate {
			config, err := converter.Handler(subscription.Link)
			if err != nil {
				continue
			}
			subscription.Data = string(config)
			subscription.UpdatedTime = time.Now()
			updates = append(updates, subscription)
			if subscription.Active {
				if utils.IsProxy {
					err = singbox.Start()
					if err != nil {
						continue
					}
					time.Sleep(5 * time.Second)
				}
			}
		}
	}
	if len(updates) > 0 {
		db.Save(&updates)
	}
}

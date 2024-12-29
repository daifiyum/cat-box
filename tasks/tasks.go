package tasks

import (
	"log"
	"time"

	U "github.com/daifiyum/cat-box/config"
	"github.com/daifiyum/cat-box/database"
	"github.com/daifiyum/cat-box/database/models"
	"github.com/daifiyum/cat-box/parser"
	"github.com/daifiyum/cat-box/singbox"
	"github.com/daifiyum/cat-box/tasks/every"
)

var task *every.Task

func Run() error {
	var err error
	interval := getUpdateDelay()
	task, err = every.NewTask(interval, updateSubscriptions)
	task.Start()
	return err
}

func UpdateDelay(interval string) {
	task.UpdateInterval(interval)
}

func getUpdateDelay() string {
	db := database.DBConn
	setting := new(models.Setting)
	db.Where("key = ?", "update_delay").First(setting)
	return setting.Value
}

func updateSubscriptions() {
	db := database.DBConn
	var subscriptions []models.Subscriptions
	db.Find(&subscriptions)
	var updates []models.Subscriptions
	var isActive bool
	for _, subscription := range subscriptions {
		if subscription.AutoUpdate {
			r, err := parser.Parser(subscription.Link, subscription.UserAgent)
			if err != nil {
				continue
			}
			subscription.Data = r
			subscription.UpdatedTime = time.Now()
			updates = append(updates, subscription)
			if subscription.Active {
				isActive = true
			}
		}
	}
	if len(updates) > 0 {
		db.Save(&updates)
	}
	if isActive {
		err := singbox.SwitchProxyMode(U.IsTun.Get())
		if err != nil {
			log.Println(err)
		}
		if U.IsCoreRunning.Get() {
			singbox.Start()
		}
	}
}

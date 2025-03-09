package database

import (
	"fmt"

	"github.com/daifiyum/cat-box/database/models"

	U "github.com/daifiyum/cat-box/common"
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

var DBConn *gorm.DB

func Init() error {
	var err error
	DBConn, err = gorm.Open(sqlite.Open("./resources/db/app.db"), &gorm.Config{})
	if err != nil {
		return fmt.Errorf("failed to connect database")
	}

	// migrate database
	DBConn.AutoMigrate(&models.Subscriptions{}, &models.Setting{})

	// create default settings if not exists
	settings := &[]models.Setting{}
	DBConn.Find(settings)
	if len(*settings) == 0 {
		settings = &[]models.Setting{
			{
				Label:       "update_interval",
				Type:        "text",
				Value:       "30m",
				Description: "更新间隔",
			},
			{
				Label:       "user_agent_type",
				Type:        "select",
				Value:       "sing-box",
				Options:     `["sing-box", "clash"]`,
				Description: "默认User-Agent类型",
			},
		}
		DBConn.Create(settings)
	}

	// apply settings
	if len(*settings) != 0 {
		for _, setting := range *settings {
			switch setting.Label {
			case "user_agent_type":
				U.DefaultUserAgent = setting.Value
			}
		}
	}

	return nil
}

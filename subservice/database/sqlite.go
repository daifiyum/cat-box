package database

import (
	"github.com/daifiyum/cat-box/subservice/models"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDatabase() {
	var err error
	DB, err = gorm.Open(sqlite.Open("./resources/db/app.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	// migrate database
	DB.AutoMigrate(&models.Subscriptions{}, &models.Setting{})

	// create default options if not exists
	DB.FirstOrCreate(&models.Setting{}, models.Setting{Key: "update_delay"})
	setting := new(models.Setting)
	DB.Where("key = ?", "update_delay").First(setting)
	if setting.Value == "none" {
		DB.Model(setting).Where("key = ?", "update_delay").Update("value", "30m")
	}
}

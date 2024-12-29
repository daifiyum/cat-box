package database

import (
	"fmt"

	"github.com/daifiyum/cat-box/database/models"

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
		settings = &[]models.Setting{{Key: "update_delay", Value: "30m"}, {Key: "default_user_agent", Value: "sing-box"}}
		DBConn.Create(settings)
	}

	return nil
}

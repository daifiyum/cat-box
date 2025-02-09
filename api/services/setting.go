package services

import (
	"errors"

	"github.com/daifiyum/cat-box/api/watcher"
	"github.com/daifiyum/cat-box/database"
	"github.com/daifiyum/cat-box/database/models"
)

func GetSetting() ([]models.Setting, error) {
	db := database.DBConn
	var setting []models.Setting

	if err := db.Find(&setting).Error; err != nil {
		return nil, errors.New("setting not found")
	}

	return setting, nil
}

func UpdateSetting(setting models.Setting) error {
	db := database.DBConn
	db.Save(&setting)

	watcher.Setting(setting)

	return nil
}

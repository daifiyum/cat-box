package services

import (
	"errors"
	"fmt"

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

func UpdateSetting(settings []models.Setting) error {
	db := database.DBConn

	// 更新数据库
	tx := db.Begin()
	for _, setting := range settings {
		var existingSetting models.Setting
		if err := tx.Where("key = ?", setting.Key).First(&existingSetting).Error; err != nil {
			tx.Rollback()
			return fmt.Errorf("setting %s not found", setting.Key)
		}

		existingSetting.Value = setting.Value
		if err := tx.Save(&existingSetting).Error; err != nil {
			tx.Rollback()
			return fmt.Errorf("failed to update setting %s", setting.Key)
		}
	}
	tx.Commit()

	// 更新监听器
	watcher.Setting(settings)

	return nil
}

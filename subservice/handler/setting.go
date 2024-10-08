package handler

import (
	"github.com/daifiyum/cat-box/subservice/database"
	"github.com/daifiyum/cat-box/subservice/models"
	"github.com/daifiyum/cat-box/task"
	"github.com/gofiber/fiber/v2"
)

// 查询指定的设置项
func GetSetting(c *fiber.Ctx) error {
	key := c.Params("key")
	db := database.DB
	setting := new(models.Setting)

	if err := db.Where("key = ?", key).First(setting).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"status": "error", "message": "Setting not found", "data": nil})
	}

	return c.JSON(fiber.Map{"status": "success", "message": "Successfully retrieved setting", "data": setting})
}

// 更新指定的设置项
func UpdateSetting(c *fiber.Ctx) error {
	key := c.Params("key")
	setting := new(models.Setting)

	if err := c.BodyParser(setting); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "error", "message": "Cannot parse body input", "data": nil})
	}

	db := database.DB
	if err := db.Model(&models.Setting{Key: key}).Where("key = ?", key).Updates(setting).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "error", "message": "Failed to update setting", "data": nil})
	}
	// 指定设置项：触发更新
	if key == "update_delay" {
		task.Scheduler()
	}
	return c.JSON(fiber.Map{"status": "success", "message": "Successfully updated setting", "data": nil})
}

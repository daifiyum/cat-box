package handler

import (
	"github.com/daifiyum/cat-box/api/services"
	"github.com/daifiyum/cat-box/database/models"
	"github.com/gofiber/fiber/v2"
)

// 获取指定的设置项
func GetSetting(c *fiber.Ctx) error {
	setting, err := services.GetSetting()
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status":  "error",
			"message": "Setting not found",
		})
	}

	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "Successfully got setting",
		"data":    setting,
	})
}

// 更新指定的设置项
func UpdateSetting(c *fiber.Ctx) error {
	var setting models.Setting
	if err := c.BodyParser(&setting); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid request body",
		})
	}

	if err := services.UpdateSetting(setting); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "Successfully updated settings",
	})
}

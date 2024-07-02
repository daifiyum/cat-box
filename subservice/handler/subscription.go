package handler

import (
	"strconv"

	"github.com/daifiyum/cat-box/subservice/database"
	"github.com/daifiyum/cat-box/subservice/models"
	"github.com/daifiyum/cat-box/subservice/parser"

	"github.com/daifiyum/cat-box/singbox"
	"github.com/daifiyum/cat-box/utils"
	"github.com/gofiber/fiber/v2"
)

// GetAllSubscribe query all subscribe
func GetAllSubscribe(c *fiber.Ctx) error {
	db := database.DB
	var subscribe []models.Subscriptions
	db.Find(&subscribe)
	return c.JSON(fiber.Map{"status": "success", "message": "All subscribe", "data": subscribe})
}

// CreateSubscribe new subscribe
func CreateSubscribe(c *fiber.Ctx) error {
	db := database.DB
	subscribe := new(models.Subscriptions)
	if err := c.BodyParser(subscribe); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "error", "message": "Couldn't create subscribe", "data": err})
	}
	res, err := parser.Handler(subscribe.Link)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "error", "message": "Couldn't create subscribe", "data": err})
	}
	subscribe.Data = string(res)
	db.Create(&subscribe)
	return c.JSON(fiber.Map{"status": "success", "message": "Created subscribe", "data": subscribe})
}

// DeleteSubscribe delete subscribe
func DeleteSubscribe(c *fiber.Ctx) error {
	id := c.Params("id")
	db := database.DB
	subscribe := new(models.Subscriptions)
	db.First(subscribe, id)
	if subscribe.Name == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "error", "message": "No subscribe found with ID", "data": nil})

	}
	db.Delete(subscribe)
	return c.JSON(fiber.Map{"status": "success", "message": "Subscribe successfully deleted", "data": nil})
}

// EditSubscribe edit subscribe
func EditSubscribe(c *fiber.Ctx) error {
	id := c.Params("id")
	db := database.DB
	subscribe := new(models.Subscriptions)
	if err := c.BodyParser(subscribe); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "error", "message": "Couldn't edit subscribe", "data": err})
	}
	db.Model(subscribe).Where("id = ?", id).Updates(subscribe)
	return c.JSON(fiber.Map{"status": "success", "message": "Subscribe successfully edited", "data": subscribe})
}

// OperateSubscribe operate subscribe to active or update
func OperateSubscribe(c *fiber.Ctx) error {
	id := c.Params("id")
	update := c.Query("update", "0")
	sw_update := c.Query("sw_update", "0")
	auto_update_str := c.FormValue("auto_update")
	auto_update, _ := strconv.ParseBool(auto_update_str)
	db := database.DB
	subscribe := new(models.Subscriptions)
	db.First(subscribe, id)

	// 手动更新
	if update == "1" {
		res, err := parser.Handler(subscribe.Link)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "error", "message": "Couldn't update subscribe", "data": err})
		}
		subscribe.Data = string(res)
		db.Save(subscribe)
		if subscribe.Active {
			if utils.IsProxy {
				err := singbox.Reload()
				if err != nil {
					utils.LogError("Failed to reload sing-box")
				}
			}
		}
		return c.JSON(fiber.Map{"status": "success", "message": "Subscribe successfully updated", "data": nil})
	}

	// 自动更新
	if sw_update == "1" {
		subscribe.AutoUpdate = auto_update
		db.Save(subscribe)
		return c.JSON(fiber.Map{"status": "success", "message": "Enable auto update", "data": nil})
	}

	// 激活订阅（选中订阅）
	db.Exec("UPDATE subscriptions SET active = CASE WHEN id = ? THEN 1 ELSE 0 END", id)
	if utils.IsProxy {
		err := singbox.Reload()
		if err != nil {
			utils.LogError("Failed to reload sing-box")
		}
	}
	return c.JSON(fiber.Map{"status": "success", "message": "Subscribe successfully active", "data": nil})
}

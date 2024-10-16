package handler

import (
	"time"

	"github.com/daifiyum/cat-box/converter"
	"github.com/daifiyum/cat-box/subservice/database"
	"github.com/daifiyum/cat-box/subservice/models"

	"github.com/daifiyum/cat-box/singbox"
	"github.com/daifiyum/cat-box/utils"
	"github.com/gofiber/fiber/v2"
)

// GetAllSubscribe query all subscribe
func GetAllSubscribe(c *fiber.Ctx) error {
	db := database.DB
	var subscribe []models.Subscriptions
	db.Order("sort_order asc").Find(&subscribe)
	return c.JSON(fiber.Map{"status": "success", "message": "All subscribe", "data": subscribe})
}

// CreateSubscribe new subscribe
func CreateSubscribe(c *fiber.Ctx) error {
	db := database.DB
	subscribe := new(models.Subscriptions)
	if err := c.BodyParser(subscribe); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "error", "message": "Couldn't create subscribe", "data": err})
	}
	res, err := converter.Handler(subscribe.Link)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "error", "message": "Couldn't create subscribe", "data": err})
	}
	var maxSortOrder int
	db.Model(subscribe).Select("MAX(sort_order)").Scan(&maxSortOrder)
	subscribe.UpdatedTime = time.Now()
	subscribe.Data = string(res)
	subscribe.SortOrder = maxSortOrder + 1
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
	db.First(subscribe, id)
	if err := c.BodyParser(subscribe); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "error", "message": "Couldn't edit subscribe", "data": err})
	}
	db.Save(subscribe)
	return c.JSON(fiber.Map{"status": "success", "message": "Subscribe successfully edited", "data": subscribe})
}

// 激活订阅
func ActiveSubscribe(c *fiber.Ctx) error {
	id := c.Params("id")
	db := database.DB
	// 激活订阅（选中订阅）
	db.Exec("UPDATE subscriptions SET active = CASE WHEN id = ? THEN 1 ELSE 0 END", id)
	if utils.IsProxy {
		err := singbox.Start()
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "error", "message": "Failed to reload sing-box"})
		}
	}
	return c.JSON(fiber.Map{"status": "success", "message": "Subscribe successfully active", "data": nil})
}

// 更新订阅
func UpdateSubscribe(c *fiber.Ctx) error {
	id := c.Params("id")
	db := database.DB
	subscribe := new(models.Subscriptions)
	db.First(subscribe, id)
	res, err := converter.Handler(subscribe.Link)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "error", "message": "Couldn't update subscribe", "data": err})
	}
	subscribe.Data = string(res)
	subscribe.UpdatedTime = time.Now()
	db.Save(subscribe)
	if subscribe.Active {
		if utils.IsProxy {
			err := singbox.Start()
			if err != nil {
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "error", "message": "Failed to reload sing-box"})
			}
		}
	}
	return c.JSON(fiber.Map{"status": "success", "message": "Subscribe successfully updated", "data": nil})
}

// 排序
func OrderSubscribe(c *fiber.Ctx) error {
	db := database.DB
	subscribe := new([]models.Subscriptions)
	if err := c.BodyParser(subscribe); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "error", "message": "Couldn't parse data", "data": err})
	}
	for index, item := range *subscribe {
		db.Model(subscribe).Where("id = ?", item.ID).Update("sort_order", index+1)
	}
	return c.JSON(fiber.Map{"status": "success", "message": "Subscribe successfully ordered", "data": nil})
}

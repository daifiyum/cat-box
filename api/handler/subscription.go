package handler

import (
	"github.com/daifiyum/cat-box/api/services"
	"github.com/daifiyum/cat-box/database/models"

	"github.com/gofiber/fiber/v2"
)

// 获取所有订阅
func GetAllSubscribe(c *fiber.Ctx) error {
	subscriptions, err := services.GetAllSubscriptions()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to get subscriptions",
		})
	}
	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "All subscriptions got",
		"data":    subscriptions,
	})
}

// 创建新订阅
func CreateSubscribe(c *fiber.Ctx) error {
	subscribe := new(models.Subscriptions)
	if err := c.BodyParser(subscribe); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid request data",
		})
	}
	if err := services.CreateSubscription(subscribe); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to create subscription",
			"error":   err.Error(),
		})
	}
	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "Subscription created successfully",
		"data":    subscribe,
	})
}

// 删除订阅
func DeleteSubscribe(c *fiber.Ctx) error {
	id := c.Params("id")
	if err := services.DeleteSubscription(id); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to delete subscription",
		})
	}
	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "Subscription deleted successfully",
	})
}

// 编辑订阅
func EditSubscribe(c *fiber.Ctx) error {
	id := c.Params("id")
	updatedSubscribe := new(models.Subscriptions)
	if err := c.BodyParser(updatedSubscribe); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid request data",
		})
	}
	if err := services.EditSubscription(id, updatedSubscribe); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to edit subscription",
		})
	}
	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "Subscription updated successfully",
		"data":    updatedSubscribe,
	})
}

// 激活订阅
func ActiveSubscribe(c *fiber.Ctx) error {
	id := c.Params("id")
	if err := services.ActivateSubscription(id); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to activate subscription",
		})
	}
	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "Subscription activated successfully",
	})
}

// 更新订阅
func UpdateSubscribe(c *fiber.Ctx) error {
	id := c.Params("id")
	if err := services.UpdateSubscription(id); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to update subscription",
		})
	}
	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "Subscription updated successfully",
	})
}

// 排序订阅
func OrderSubscribe(c *fiber.Ctx) error {
	subscriptions := new([]models.Subscriptions)
	if err := c.BodyParser(subscriptions); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid request data",
		})
	}
	if err := services.OrderSubscriptions(*subscriptions); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to order subscriptions",
		})
	}
	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "Subscriptions reordered successfully",
	})
}

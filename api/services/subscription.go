package services

import (
	"time"

	U "github.com/daifiyum/cat-box/common"
	"github.com/daifiyum/cat-box/database"
	"github.com/daifiyum/cat-box/database/models"
	S "github.com/daifiyum/cat-box/singbox"
	P "github.com/daifiyum/cat-box/subscription"
)

func GetAllSubscriptions() ([]models.Subscriptions, error) {
	db := database.DBConn
	var subscriptions []models.Subscriptions
	if err := db.Order("sort_order asc").Find(&subscriptions).Error; err != nil {
		return nil, err
	}
	return subscriptions, nil
}

func CreateSubscription(subscribe *models.Subscriptions) error {
	db := database.DBConn

	res, err := P.Subscription(subscribe.Link, U.DefaultUserAgent)
	if err != nil {
		return err
	}
	var maxSortOrder int
	db.Model(subscribe).Select("MAX(sort_order)").Scan(&maxSortOrder)
	subscribe.UpdatedTime = time.Now()
	subscribe.Data = res
	subscribe.UserAgent = U.DefaultUserAgent
	subscribe.SortOrder = maxSortOrder + 1
	return db.Create(&subscribe).Error
}

func DeleteSubscription(id string) error {
	db := database.DBConn
	subscribe := new(models.Subscriptions)
	if err := db.First(subscribe, id).Error; err != nil {
		return err
	}
	return db.Delete(subscribe).Error
}

func EditSubscription(id string, updatedSubscribe *models.Subscriptions) error {
	db := database.DBConn
	subscribe := new(models.Subscriptions)
	if err := db.First(subscribe, id).Error; err != nil {
		return err
	}
	subscribe.Name = updatedSubscribe.Name
	subscribe.Link = updatedSubscribe.Link
	subscribe.UserAgent = updatedSubscribe.UserAgent
	subscribe.AutoUpdate = updatedSubscribe.AutoUpdate
	return db.Save(subscribe).Error
}

func ActivateSubscription(id string) error {
	db := database.DBConn

	if err := db.Exec("UPDATE subscriptions SET active = CASE WHEN id = ? THEN 1 ELSE 0 END", id).Error; err != nil {
		return err
	}

	if err := S.SwitchProxyMode(U.IsTun.Get()); err != nil {
		return err
	}

	if U.IsCoreRunning.Get() {
		return S.Start()
	}

	return nil
}

func UpdateSubscription(id string) error {
	db := database.DBConn
	subscribe := new(models.Subscriptions)

	if err := db.First(subscribe, id).Error; err != nil {
		return err
	}

	res, err := P.Subscription(subscribe.Link, subscribe.UserAgent)
	if err != nil {
		return err
	}

	subscribe.Data = res
	subscribe.UpdatedTime = time.Now()
	if err := db.Save(subscribe).Error; err != nil {
		return err
	}

	if subscribe.Active {
		err := S.SwitchProxyMode(U.IsTun.Get())
		if err != nil {
			return err
		}
		if U.IsCoreRunning.Get() {
			return S.Start()
		}
	}

	return nil
}

func OrderSubscriptions(subscriptions []models.Subscriptions) error {
	db := database.DBConn
	for index, item := range subscriptions {
		if err := db.Model(&item).Where("id = ?", item.ID).Update("sort_order", index).Error; err != nil {
			return err
		}
	}
	return nil
}

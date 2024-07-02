package singbox

import (
	"encoding/json"
	"errors"
	"os"

	"github.com/daifiyum/cat-box/subservice/database"
	"github.com/daifiyum/cat-box/subservice/models"
	"github.com/daifiyum/cat-box/utils"
)

func getConfig() string {
	db := database.DB
	var subscriptions []models.Subscriptions
	db.Find(&subscriptions)

	for _, subscription := range subscriptions {
		if subscription.Active {
			return subscription.Data
		}
	}
	return ""
}

func GenerateConfig() error {
	config := getConfig()
	if config == "" {
		errmsg := "singbox config is empty"
		utils.LogError(errmsg)
		return errors.New(errmsg)
	}
	var newInbounds []interface{}
	mapConfig := make(map[string]interface{})
	json.Unmarshal([]byte(config), &mapConfig)

	isTun := utils.IsTun
	if !isTun {
		for _, inbound := range mapConfig["inbounds"].([]interface{}) {
			inbound, _ := inbound.(map[string]interface{})
			if inbound["type"] == "mixed" {
				inbound["set_system_proxy"] = true
			}
			if inbound["type"] != "tun" {
				newInbounds = append(newInbounds, inbound)
			}
		}
		mapConfig["inbounds"] = newInbounds
	} else {
		for _, inbound := range mapConfig["inbounds"].([]interface{}) {
			inbound, _ := inbound.(map[string]interface{})
			if inbound["type"] == "mixed" {
				inbound["set_system_proxy"] = false
			}
		}
	}

	singConf, _ := json.MarshalIndent(mapConfig, "", "  ")
	err := os.WriteFile("./resources/core/config.json", []byte(singConf), 0666)
	if err != nil {
		return err
	}
	return nil
}

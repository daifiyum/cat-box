package singbox

import (
	"encoding/json"
	"errors"
	"os"
	"syscall"

	"github.com/daifiyum/cat-box/subservice/database"
	"github.com/daifiyum/cat-box/subservice/models"
	"github.com/daifiyum/cat-box/utils"
	"golang.org/x/sys/windows"
)

func terminateProc(pid int) error {
	dll, err := windows.LoadDLL("kernel32.dll")
	if err != nil {
		return err
	}
	defer dll.Release()

	f, err := dll.FindProc("AttachConsole")
	if err != nil {
		return err
	}
	r1, _, err := f.Call(uintptr(pid))
	if r1 == 0 && err != syscall.ERROR_ACCESS_DENIED {
		return err
	}

	f, err = dll.FindProc("SetConsoleCtrlHandler")
	if err != nil {
		return err
	}
	r1, _, err = f.Call(0, 1)
	if r1 == 0 {
		return err
	}
	f, err = dll.FindProc("GenerateConsoleCtrlEvent")
	if err != nil {
		return err
	}
	r1, _, err = f.Call(windows.CTRL_BREAK_EVENT, uintptr(pid))
	if r1 == 0 {
		return err
	}
	return nil
}

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

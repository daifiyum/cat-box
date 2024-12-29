package singbox

import (
	"encoding/json"
	"fmt"
	"hash/crc32"
	"os"
	"strings"

	U "github.com/daifiyum/cat-box/config"
	"github.com/daifiyum/cat-box/database"
	"github.com/daifiyum/cat-box/database/models"
)

const (
	templatePath   = "./resources/template/template.json"
	coreConfigPath = "./resources/core/config.json"
)

// 初始化配置函数，生成默认配置模板，并提取必要配置传给公共变量BoxConfig
func InitConfig() (map[string]any, error) {
	template, err := ReadTemplate()
	if err != nil {
		return nil, err
	}

	var templateMap map[string]any
	err = json.Unmarshal(template, &templateMap)
	if err != nil {
		return nil, err
	}

	/*
		experimental默认处理
	*/
	experiments := templateMap["experimental"].(map[string]any)
	// 读取clash api端口
	clashAPI := experiments["clash_api"].(map[string]any)
	externalController, ok := clashAPI["external_controller"].(string)
	if !ok {
		return nil, fmt.Errorf("external_controller 配置不存在或格式错误")
	}
	clashAPIPort := strings.Split(externalController, ":")[1]
	U.Box.ClashAPIPort = clashAPIPort

	/*
		inbounds默认处理
	*/
	inbounds, _ := templateMap["inbounds"].([]any)
	inboundsIndex := []int{}
	for i := range inbounds {
		inbound := inbounds[i].(map[string]any)
		if inbound["type"].(string) == "mixed" {
			inbound["set_system_proxy"] = true
			U.Box.MixedListenPort = fmt.Sprintf("%.0f", inbound["listen_port"].(float64))
		}

		// 删除tun入站
		if inbound["type"].(string) == "tun" {
			inboundsIndex = append(inboundsIndex, i)
			U.Box.TunInbound = inbound
		}
	}

	// 根据收集的索引删除不需要的入站
	for i := len(inboundsIndex) - 1; i >= 0; i-- {
		indexToRemove := inboundsIndex[i]
		inbounds = append(inbounds[:indexToRemove], inbounds[indexToRemove+1:]...)
	}

	// 将修改后的入站覆写给模板入站
	templateMap["inbounds"] = inbounds

	return templateMap, nil
}

// 将初始化后的配置模板与传入的部分配置合并
func OutboundsConfig(templateMap map[string]any, outbounds []byte) (map[string]any, error) {
	var outboundsSlice []any
	err := json.Unmarshal(outbounds, &outboundsSlice)
	if err != nil {
		return nil, err
	}

	outboundTags := []any{}
	for _, outboundAny := range outboundsSlice {
		outboundMap := outboundAny.(map[string]any)
		outboundTags = append(outboundTags, outboundMap["tag"])
	}

	templateOutbounds, _ := templateMap["outbounds"].([]any)
	for i := range templateOutbounds {
		templateOutbound := templateOutbounds[i].(map[string]any)
		t := templateOutbound["type"].(string)
		if t == "selector" || t == "urltest" {
			templateOutbound["outbounds"] = append(templateOutbound["outbounds"].([]any), outboundTags...)
		}
	}
	templateOutbounds = append(templateOutbounds, outboundsSlice...)
	templateMap["outbounds"] = templateOutbounds

	return templateMap, nil
}

// 生成默认配置并写入本地文件
func DefaultConfig() error {
	configMap, err := MergeConfig()
	if err != nil {
		return err
	}
	err = WriteCoreFile(configMap)
	if err != nil {
		return err
	}
	return nil
}

// 生成tun配置并写入本地文件
func TunConfig() error {
	configMap, err := MergeConfig()
	if err != nil {
		return err
	}

	inbounds, _ := configMap["inbounds"].([]any)
	for i := range inbounds {
		inbound := inbounds[i].(map[string]any)
		if inbound["type"].(string) == "mixed" {
			inbound["set_system_proxy"] = false
		}
	}
	configMap["inbounds"] = inbounds
	configMap["inbounds"] = append(configMap["inbounds"].([]any), U.Box.TunInbound)

	err = WriteCoreFile(configMap)
	if err != nil {
		return err
	}

	return nil
}

// 检测模板是否改动
func CompareTemplate() error {
	data, err := ReadTemplate()
	if err != nil {
		return err
	}

	curr := crc32.ChecksumIEEE(data)
	if U.PrevCrc32 == 0 {
		U.PrevCrc32 = curr
		return SwitchProxyMode(U.IsTun.Get())
	}

	if curr != U.PrevCrc32 {
		U.PrevCrc32 = curr
		return SwitchProxyMode(U.IsTun.Get())
	}
	return nil
}

// sing-box核心配置写入到本地
func WriteCoreFile(data map[string]any) error {
	config, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return err
	}
	err = os.WriteFile(coreConfigPath, config, 0666)
	if err != nil {
		return err
	}
	return nil
}

// 读取未初始化的模板配置
func ReadTemplate() ([]byte, error) {
	data, err := os.ReadFile(templatePath)
	if err != nil {
		return nil, err
	}
	return data, nil
}

// 获取激活订阅
func GetActiveSub() (string, error) {
	db := database.DBConn
	var subscription models.Subscriptions
	if err := db.Where("active = ?", true).First(&subscription).Error; err != nil {
		return "", err
	}
	return subscription.Data, nil
}

// 合并初始化后的配置和出站配置（默认配置）
func MergeConfig() (map[string]any, error) {
	data, err := GetActiveSub()
	if err != nil {
		return nil, err
	}

	defaultConfig, err := InitConfig()
	if err != nil {
		return nil, err
	}
	config, err := OutboundsConfig(defaultConfig, []byte(data))
	if err != nil {
		return nil, err
	}
	return config, nil
}

// 切换代理模式
func SwitchProxyMode(s bool) error {
	if s {
		return TunConfig()
	}
	return DefaultConfig()
}

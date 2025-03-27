package singbox

import (
	"encoding/json"
	"fmt"
	"hash/crc32"
	"os"
	"strings"

	U "github.com/daifiyum/cat-box/common"
	"github.com/daifiyum/cat-box/database"
	"github.com/daifiyum/cat-box/database/models"
)

type ConfigurationError struct {
	Context string
	Err     error
}

func (e *ConfigurationError) Error() string {
	return fmt.Sprintf("%s: %v", e.Context, e.Err)
}

func (e *ConfigurationError) Unwrap() error {
	return e.Err
}

func InitConfig() (map[string]any, error) {
	template, err := ReadTemplate()
	if err != nil {
		return nil, &ConfigurationError{"读取模板失败", err}
	}

	var templateMap map[string]any
	if err = json.Unmarshal(template, &templateMap); err != nil {
		return nil, &ConfigurationError{"解析模板失败", err}
	}

	if err = processExperimentalConfig(templateMap); err != nil {
		return nil, err
	}

	if err = processInboundsConfig(templateMap); err != nil {
		return nil, err
	}

	return templateMap, nil
}

func processExperimentalConfig(templateMap map[string]any) error {
	experiments, ok := templateMap["experimental"].(map[string]any)
	if !ok {
		return &ConfigurationError{"配置错误", fmt.Errorf("experimental 字段缺失或类型错误")}
	}

	clashAPI, ok := experiments["clash_api"].(map[string]any)
	if !ok {
		return &ConfigurationError{"配置错误", fmt.Errorf("clash_api 字段缺失或类型错误")}
	}

	controller, ok := clashAPI["external_controller"].(string)
	if !ok {
		return &ConfigurationError{"配置错误", fmt.Errorf("external_controller 字段缺失或类型错误")}
	}

	parts := strings.Split(controller, ":")
	if len(parts) < 2 {
		return &ConfigurationError{"配置错误", fmt.Errorf("invalid external_controller format")}
	}
	U.Box.ClashAPIPort = parts[1]
	return nil
}

func processInboundsConfig(templateMap map[string]any) error {
	inbounds, ok := templateMap["inbounds"].([]any)
	if !ok {
		return &ConfigurationError{"配置错误", fmt.Errorf("inbounds 字段缺失或类型错误")}
	}

	var validInbounds []any
	var tunInbound map[string]any

	for _, item := range inbounds {
		inbound, ok := item.(map[string]any)
		if !ok {
			continue
		}

		switch inbound["type"] {
		case "mixed":
			if err := processMixedInbound(inbound); err != nil {
				return err
			}
			validInbounds = append(validInbounds, inbound)
		case "tun":
			tunInbound = inbound
		default:
			validInbounds = append(validInbounds, inbound)
		}
	}

	U.Box.TunInbound = tunInbound
	templateMap["inbounds"] = validInbounds
	return nil
}

func processMixedInbound(inbound map[string]any) error {
	port, ok := inbound["listen_port"].(float64)
	if !ok {
		return &ConfigurationError{"配置错误", fmt.Errorf("mixed 入站缺少 listen_port")}
	}

	inbound["set_system_proxy"] = true
	U.Box.MixedListenPort = fmt.Sprintf("%.0f", port)
	return nil
}

func OutboundsConfig(templateMap map[string]any, outbounds []byte) (map[string]any, error) {
	var outboundsSlice []any
	if err := json.Unmarshal(outbounds, &outboundsSlice); err != nil {
		return nil, &ConfigurationError{"解析出站配置失败", err}
	}

	tags, err := extractOutboundTags(outboundsSlice)
	if err != nil {
		return nil, err
	}

	if err := updateTemplateOutbounds(templateMap, tags); err != nil {
		return nil, err
	}

	templateOutbounds, ok := templateMap["outbounds"].([]any)
	if !ok {
		return nil, &ConfigurationError{"配置错误", fmt.Errorf("outbounds 字段缺失或类型错误")}
	}

	templateMap["outbounds"] = append(templateOutbounds, outboundsSlice...)
	return templateMap, nil
}

func extractOutboundTags(outbounds []any) ([]any, error) {
	tags := make([]any, 0, len(outbounds))
	for _, ob := range outbounds {
		outbound, ok := ob.(map[string]any)
		if !ok {
			return nil, &ConfigurationError{"配置错误", fmt.Errorf("无效的出站配置格式")}
		}

		tag, exists := outbound["tag"]
		if !exists {
			return nil, &ConfigurationError{"配置错误", fmt.Errorf("出站配置缺少 tag 字段")}
		}
		tags = append(tags, tag)
	}
	return tags, nil
}

func updateTemplateOutbounds(templateMap map[string]any, tags []any) error {
	templateOutbounds, ok := templateMap["outbounds"].([]any)
	if !ok {
		return &ConfigurationError{"配置错误", fmt.Errorf("outbounds 字段缺失或类型错误")}
	}

	for i, item := range templateOutbounds {
		outbound, ok := item.(map[string]any)
		if !ok {
			continue
		}

		t, ok := outbound["type"].(string)
		if !ok {
			continue
		}

		if t == "selector" || t == "urltest" {
			existing, _ := outbound["outbounds"].([]any)
			outbound["outbounds"] = append(existing, tags...)
			templateOutbounds[i] = outbound
		}
	}
	return nil
}

func DefaultConfig() error {
	return handleCoreConfig(false)
}

func TunConfig() error {
	return handleCoreConfig(true)
}

func handleCoreConfig(tunMode bool) error {
	configMap, err := MergeConfig()
	if err != nil {
		return err
	}

	if tunMode {
		if err := configureTunInbounds(configMap); err != nil {
			return err
		}
	}

	return WriteCoreFile(configMap)
}

func configureTunInbounds(configMap map[string]any) error {
	inbounds, ok := configMap["inbounds"].([]any)
	if !ok {
		return &ConfigurationError{"配置错误", fmt.Errorf("inbounds 字段缺失或类型错误")}
	}

	for _, item := range inbounds {
		inbound, ok := item.(map[string]any)
		if !ok {
			continue
		}
		if inbound["type"] == "mixed" {
			inbound["set_system_proxy"] = false
		}
	}

	if U.Box.TunInbound != nil {
		configMap["inbounds"] = append(inbounds, U.Box.TunInbound)
	}
	return nil
}

func CompareTemplate() error {
	data, err := ReadTemplate()
	if err != nil {
		return err
	}

	currentCRC := crc32.ChecksumIEEE(data)
	if U.PrevCrc32 == 0 {
		U.PrevCrc32 = currentCRC
		return SwitchProxyMode(U.IsTun.Get())
	}

	if currentCRC != U.PrevCrc32 {
		U.PrevCrc32 = currentCRC
		return SwitchProxyMode(U.IsTun.Get())
	}
	return nil
}

func WriteCoreFile(data map[string]any) error {
	config, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return &ConfigurationError{"序列化配置失败", err}
	}

	if err := os.WriteFile("./resources/core/config.json", config, 0644); err != nil {
		return &ConfigurationError{"写入配置文件失败", err}
	}
	return nil
}

func ReadTemplate() ([]byte, error) {
	data, err := os.ReadFile("./resources/template/template.json")
	if err != nil {
		return nil, &ConfigurationError{"读取模板文件失败", err}
	}
	return data, nil
}

func GetActiveSub() (string, error) {
	db := database.DBConn
	var subscription models.Subscriptions
	if err := db.Where("active = ?", true).First(&subscription).Error; err != nil {
		return "", &ConfigurationError{"获取激活订阅失败", err}
	}
	return subscription.Data, nil
}

func MergeConfig() (map[string]any, error) {
	data, err := GetActiveSub()
	if err != nil {
		return nil, err
	}

	baseConfig, err := InitConfig()
	if err != nil {
		return nil, err
	}

	return OutboundsConfig(baseConfig, []byte(data))
}

func SwitchProxyMode(tunMode bool) error {
	if tunMode {
		return TunConfig()
	}
	return DefaultConfig()
}

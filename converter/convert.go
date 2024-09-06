package converter

import (
	"encoding/json"
	"errors"
	"os"
	"strings"

	C "github.com/daifiyum/cat-box/converter/common"
	"github.com/daifiyum/cat-box/converter/model"
	"github.com/daifiyum/cat-box/converter/parser"
	"github.com/daifiyum/cat-box/utils"
)

func ConvertCProxyToSProxy(proxy string) (model.Outbound, error) {
	for prefix, parseFunc := range parser.ParserMap {
		if strings.HasPrefix(proxy, prefix) {
			proxy, err := parseFunc(proxy)
			if err != nil {
				return model.Outbound{}, err
			}
			return proxy, nil
		}
	}
	return model.Outbound{}, errors.New("unknown proxy format")
}

func ConvertCProxyToJson(proxy string) (string, error) {
	sProxy, err := ConvertCProxyToSProxy(proxy)
	if err != nil {
		return "", err
	}
	data, err := json.Marshal(&sProxy)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func convertSubscribe(data []byte) ([]map[string]interface{}, error) {
	content := string(data)
	arr := strings.Split(content, "\n")
	var outbounds []map[string]interface{}
	for _, line := range arr {
		line = strings.TrimRight(line, " \r")
		if line == "" {
			continue
		}

		node, err := ConvertCProxyToJson(line)
		if err != nil {
			continue
		}
		mapNode, _ := C.StringToMap(node)
		outbounds = append(outbounds, mapNode)

	}
	if len(outbounds) <= 0 {
		return nil, errors.New("no content")
	}
	return outbounds, nil
}

func Handler(url string) ([]byte, error) {

	req, err := C.FetchSubscribe(url)
	if err != nil {
		utils.LogError("fetchSubscribe errors")
		return nil, err
	}

	outbounds, err := convertSubscribe(req)
	if err != nil {
		utils.LogError("convertSubscribe error")
		return nil, err
	}

	template, err := os.ReadFile("./resources/template/template.json")
	if err != nil {
		return nil, err
	}
	mapTemplate := make(map[string]interface{})
	json.Unmarshal(template, &mapTemplate)

	// add tag
	for _, i := range mapTemplate["outbounds"].([]interface{}) {
		m, _ := i.(map[string]interface{})
		if m["type"] == "selector" || m["type"] == "urltest" {
			for _, f := range outbounds {
				m["outbounds"] = append(m["outbounds"].([]interface{}), f["tag"])
			}
		}
	}

	// add outbound
	for _, i := range outbounds {
		mapTemplate["outbounds"] = append(mapTemplate["outbounds"].([]interface{}), i)
	}

	config, err := json.MarshalIndent(mapTemplate, "", "  ")
	if err != nil {
		return nil, err
	}

	return config, nil
}

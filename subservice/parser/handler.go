package parser

import (
	"encoding/json"
	"errors"
	"os"
	"strings"

	C "github.com/daifiyum/cat-box/converter"
	"github.com/daifiyum/cat-box/utils"
)

func convertSubscribe(data []byte) ([]map[string]interface{}, error) {
	content := string(data)
	arr := strings.Split(content, "\n")
	var outbounds []map[string]interface{}
	for _, line := range arr {
		line = strings.TrimRight(line, " \r")
		if line == "" {
			continue
		}

		node, err := C.ConvertCProxyToJson(line)
		if err != nil {
			continue
		}
		mapNode, _ := StringToMap(node)
		outbounds = append(outbounds, mapNode)

	}
	if len(outbounds) <= 0 {
		return nil, errors.New("no content")
	}
	return outbounds, nil
}

func Handler(url string) ([]byte, error) {

	req, err := FetchSubscribe(url)
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

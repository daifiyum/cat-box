package parser

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/daifiyum/cat-box/utils"
	"github.com/hiddify/ray2sing/ray2sing"
)

func fetchSubscribe(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	decodedBody, err := base64.URLEncoding.DecodeString(string(body))
	if err != nil {
		return nil, err
	}
	return decodedBody, nil
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
		scheme, _, found := strings.Cut(line, "://")
		if !found {
			continue
		}
		scheme = strings.ToLower(scheme)
		switch scheme {
		case "ss":
			node, err := ray2sing.ShadowsocksSingbox(line)
			if err != nil {
				continue
			}
			mapNode, _ := StructToMap(node)
			outbounds = append(outbounds, mapNode)
		case "vmess":
			node, err := ray2sing.VmessSingbox(line)
			if err != nil {
				continue
			}
			mapNode, _ := StructToMap(node)
			outbounds = append(outbounds, mapNode)
		case "trojan":
			node, err := ray2sing.TrojanSingbox(line)
			if err != nil {
				continue
			}
			mapNode, _ := StructToMap(node)
			outbounds = append(outbounds, mapNode)
		case "vless":
			node, err := ray2sing.VlessSingbox(line)
			if err != nil {
				continue
			}
			mapNode, _ := StructToMap(node)
			outbounds = append(outbounds, mapNode)
		case "hysteria":
			node, err := ray2sing.HysteriaSingbox(line)
			if err != nil {
				continue
			}
			mapNode, _ := StructToMap(node)
			outbounds = append(outbounds, mapNode)
		case "hysteria2":
			node, err := ray2sing.Hysteria2Singbox(line)
			if err != nil {
				continue
			}
			mapNode, _ := StructToMap(node)
			outbounds = append(outbounds, mapNode)
		}
	}
	if len(outbounds) <= 0 {
		return nil, errors.New("no content")
	}
	return outbounds, nil
}

func StructToMap(obj interface{}) (map[string]interface{}, error) {
	jsonBytes, err := json.Marshal(obj)
	if err != nil {
		return nil, err
	}

	var result map[string]interface{}
	err = json.Unmarshal(jsonBytes, &result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func Handler(url string) ([]byte, error) {
	req, err := fetchSubscribe(url)
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

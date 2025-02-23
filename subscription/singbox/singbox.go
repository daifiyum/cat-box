package singbox

import (
	"encoding/json"
)

func CutOutbounds(data string) ([]any, error) {
	var subscribeMap map[string]any
	err := json.Unmarshal([]byte(data), &subscribeMap)
	if err != nil {
		return nil, err
	}
	subscribeOutbounds, _ := subscribeMap["outbounds"].([]any)
	outbounds := []any{}
	for i := range subscribeOutbounds {
		subscribeOutbound := subscribeOutbounds[i].(map[string]any)
		t := subscribeOutbound["type"].(string)
		switch t {
		case "shadowsocks":
			outbounds = append(outbounds, subscribeOutbound)
		case "vmess":
			outbounds = append(outbounds, subscribeOutbound)
		case "vless":
			outbounds = append(outbounds, subscribeOutbound)
		case "trojan":
			outbounds = append(outbounds, subscribeOutbound)
		case "hysteria":
			outbounds = append(outbounds, subscribeOutbound)
		case "hysteria2":
			outbounds = append(outbounds, subscribeOutbound)
		case "shadowtls":
			outbounds = append(outbounds, subscribeOutbound)
		case "tuic":
			outbounds = append(outbounds, subscribeOutbound)
		}
	}

	return outbounds, nil
}

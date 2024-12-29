package singbox

import (
	"encoding/json"
)

var typeMap = map[string]struct{}{
	"shadowsocks": {},
	"vmess":       {},
	"vless":       {},
	"trojan":      {},
	"hysteria":    {},
	"hysteria2":   {},
	"shadowtls":   {},
	"tuic":        {},
}

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
		if _, exists := typeMap[t]; exists {
			outbounds = append(outbounds, subscribeOutbound)
		}
	}

	return outbounds, nil
}

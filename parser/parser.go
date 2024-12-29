package parser

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/daifiyum/cat-box/parser/native"
	"github.com/daifiyum/cat-box/parser/singbox"
)

// 订阅解析，输出为json格式出站
func Parser(url, user_agent string) (string, error) {
	resp, err := httpGet(url, user_agent, 10*time.Second)
	if err != nil {
		return "", fmt.Errorf("failed to fetch URL: %w", err)
	}

	var outbounds any

	switch user_agent {
	case "sing-box":
		outbounds, err = singbox.CutOutbounds(resp)
		if err != nil {
			return "", fmt.Errorf("failed to parse sing-box outbounds: %w", err)
		}
	case "native":
		safeData, err := base64URLSafe(resp)
		if err != nil {
			return "", fmt.Errorf("failed to decode base64: %w", err)
		}
		outboundsRaw, err := native.NewNativeURIParser(safeData)
		if err != nil {
			log.Println("native parser error:", err)
			return "", fmt.Errorf("failed to parse native outbounds: %w", err)
		}
		outbounds = badOutbounds(outboundsRaw)
	default:
		return "", fmt.Errorf("unsupported source: %s", user_agent)
	}

	outboundsJson, err := json.Marshal(outbounds)
	if err != nil {
		return "", err
	}

	return string(outboundsJson), nil
}

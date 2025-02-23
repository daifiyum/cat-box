package subscription

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	U "github.com/daifiyum/cat-box/common"
	C "github.com/daifiyum/cat-box/subscription/singbox"
	"github.com/daifiyum/cat-box/subscription/v2ray"
)

// 订阅解析，输出为json格式出站
func Subscription(url, ua string) (string, error) {
	resp, err := httpGet(url, ua, 10*time.Second)
	if err != nil {
		return "", fmt.Errorf("failed to fetch URL: %w", err)
	}

	var outbounds any

	switch ua {
	case "sing-box":
		outbounds, err = C.CutOutbounds(resp)
		if err != nil {
			return "", fmt.Errorf("failed to parse sing-box outbounds: %w", err)
		}
	case "v2ray":
		safeData, err := base64URLSafe(resp)
		if err != nil {
			return "", fmt.Errorf("failed to decode base64: %w", err)
		}
		outboundsRaw, err := v2ray.NewNativeURIParser(safeData)
		if err != nil {
			log.Println("native parser error:", err)
			return "", fmt.Errorf("failed to parse native outbounds: %w", err)
		}
		outbounds = badOutbounds(outboundsRaw)
	default:
		return "", fmt.Errorf("unsupported source: %s", U.DefaultUserAgent)
	}

	outboundsJson, err := json.Marshal(outbounds)
	if err != nil {
		return "", err
	}

	return string(outboundsJson), nil
}

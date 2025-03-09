package subscription

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	U "github.com/daifiyum/cat-box/common"
	"github.com/daifiyum/cat-box/subscription/clash"
	C "github.com/daifiyum/cat-box/subscription/singbox"
	"github.com/sagernet/sing/common/json/badjson"
)

// 订阅解析，输出为json格式出站
func Subscription(url, ua string) (string, error) {
	content, err := httpGet(url, ua, 10*time.Second)
	if err != nil {
		return "", fmt.Errorf("failed to fetch URL: %w", err)
	}

	switch ua {
	case "sing-box":
		outbounds, err := C.CutOutbounds(content)
		if err != nil {
			return "", fmt.Errorf("failed to parse sing-box outbounds: %w", err)
		}
		outboundsJson, err := json.Marshal(outbounds)
		if err != nil {
			return "", err
		}
		return string(outboundsJson), nil
	case "clash":
		outbounds, err := clash.ParseClashSubscription(content)
		if err != nil {
			return "", fmt.Errorf("failed to parse sing-box outbounds: %w", err)
		}
		outboundsJson, err := badjson.MarshallObjectsContext(context.Background(), outbounds)
		if err != nil {
			return "", err
		}
		return string(outboundsJson), nil
	default:
		return "", fmt.Errorf("unsupported ua: %s", U.DefaultUserAgent)
	}
}

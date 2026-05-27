package subscription

import (
	"context"
	"fmt"
	"time"

	"github.com/daifiyum/cat-box/subscription/clash"
	"github.com/sagernet/sing/common/json/badjson"
)

// 订阅解析，输出为json格式出站
func Subscription(url, ua string) (string, error) {
	content, err := httpGet(url, ua, 10*time.Second)
	if err != nil {
		return "", fmt.Errorf("failed to fetch URL: %w", err)
	}

	outbounds, err := clash.ParseClashSubscription(context.Background(), content)
	if err != nil {
		return "", fmt.Errorf("failed to parse clash outbounds: %w", err)
	}
	outboundsJson, err := badjson.MarshallObjectsContext(context.Background(), outbounds)
	if err != nil {
		return "", err
	}
	return string(outboundsJson), nil
}

package subscription

import (
	"context"
	"fmt"
	"time"

	"github.com/daifiyum/cparser"
	"github.com/sagernet/sing-box/option"
	"github.com/sagernet/sing/common/json/badjson"
)

// Subscription 解析订阅并输出 JSON 出站。insecure 为 true 时跳过节点 TLS
// 证书验证；省略或为 false 时保留源订阅设置，不影响下载订阅时的 HTTPS 验证。
func Subscription(url, ua string, insecure ...bool) (string, error) {
	content, err := httpGet(url, ua, 10*time.Second)
	if err != nil {
		return "", fmt.Errorf("failed to fetch URL: %w", err)
	}

	outbounds, err := cparser.ParseClashSubscription(context.Background(), content)
	if err != nil {
		return "", fmt.Errorf("failed to parse clash outbounds: %w", err)
	}
	if len(insecure) > 0 && insecure[0] {
		for _, outbound := range outbounds {
			if wrapper, ok := outbound.Options.(option.OutboundTLSOptionsWrapper); ok {
				if tls := wrapper.TakeOutboundTLSOptions(); tls != nil && tls.Enabled {
					tls.Insecure = true
				}
			}
		}
	}
	outboundsJson, err := badjson.MarshallObjectsContext(context.Background(), outbounds)
	if err != nil {
		return "", err
	}
	return string(outboundsJson), nil
}

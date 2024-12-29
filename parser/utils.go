package parser

import (
	"encoding/base64"
	"fmt"
	"strings"
	"time"

	U "github.com/daifiyum/cat-box/config"
	"github.com/sagernet/sing-box/option"
	"github.com/valyala/fasthttp"
	"github.com/valyala/fasthttp/fasthttpproxy"
)

func userAgent(ua string) string {
	switch ua {
	case "sing-box":
		return "sing-box"
	case "native":
		return "native"
	}
	return ""
}

func httpGet(url, user_agent string, timeout time.Duration) (string, error) {
	var client fasthttp.Client

	if U.IsCoreRunning.Get() {
		proxy := fmt.Sprintf("socks5://127.0.0.1:%s", U.Box.MixedListenPort)
		client.Dial = fasthttpproxy.FasthttpSocksDialer(proxy)
	}

	req := fasthttp.AcquireRequest()
	resp := fasthttp.AcquireResponse()
	defer func() {
		fasthttp.ReleaseRequest(req)
		fasthttp.ReleaseResponse(resp)
	}()

	req.SetRequestURI(url)
	req.Header.SetMethod("GET")

	ua := userAgent(user_agent)
	if ua == "" {
		return "", fmt.Errorf("user-agent not set")
	}
	req.Header.Set("User-Agent", ua)

	if err := client.DoTimeout(req, resp, timeout); err != nil {
		return "", fmt.Errorf("request failed: %w", err)
	}

	return string(resp.Body()), nil
}

// 安全的base64
func base64URLSafe(data string) (string, error) {
	urlSafe := strings.NewReplacer("+", "-", "/", "_", "=", "").Replace(data)
	r, err := base64.RawURLEncoding.DecodeString(urlSafe)
	if err != nil {
		return "", err
	}
	return string(r), nil
}

// 糟糕节点处理(重复节点名称)
func badOutbounds(outbounds []option.Outbound) []option.Outbound {
	nameCount := make(map[string]int)
	r := make([]option.Outbound, len(outbounds))

	for i, outbound := range outbounds {
		nameCount[outbound.Tag]++
		if nameCount[outbound.Tag] == 1 {
			r[i] = outbound
		} else {
			outbound.Tag = fmt.Sprintf("%s%d", outbound.Tag, nameCount[outbound.Tag])
			r[i] = outbound
		}
	}
	return r
}

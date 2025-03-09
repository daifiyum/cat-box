package subscription

import (
	"fmt"
	"time"

	U "github.com/daifiyum/cat-box/common"
	"github.com/valyala/fasthttp"
	"github.com/valyala/fasthttp/fasthttpproxy"
)

func httpGet(url, ua string, timeout time.Duration) (string, error) {
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
	req.Header.Set("User-Agent", ua)

	if err := client.DoTimeout(req, resp, timeout); err != nil {
		return "", fmt.Errorf("request failed: %w", err)
	}

	return string(resp.Body()), nil
}

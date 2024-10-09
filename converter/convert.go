package converter

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/daifiyum/cat-box/converter/native"
	"github.com/daifiyum/cat-box/utils"
	"github.com/sagernet/sing-box/option"
	"golang.org/x/net/proxy"
)

// 将base64转换成url安全的格式进行解码
func DecodeBase64URLSafe(content string) ([]byte, error) {
	urlSafe := strings.NewReplacer("+", "-", "/", "_", "=", "").Replace(content)
	result, err := base64.RawURLEncoding.DecodeString(urlSafe)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// 获取代理端口
func GetMixedPort() (string, error) {
	template, err := os.ReadFile("./resources/template/template.json")
	if err != nil {
		return "", err
	}

	var options option.Options

	options.UnmarshalJSON(template)
	var portStr string
	for _, i := range options.Inbounds {
		if i.Type == "mixed" {
			portStr = strconv.Itoa(int(i.MixedOptions.ListenPort))
			break
		}
	}
	return portStr, nil
}

// 获取订阅数据，开启代理：通过代理请求，未开启代理：通过直连请求
func FetchSubscribe(suburl string) ([]byte, error) {
	var client *http.Client
	if utils.IsProxy {
		mixedPort, err := GetMixedPort()
		if err != nil {
			return nil, err
		}
		proxyURI, err := url.Parse("socks5://127.0.0.1:" + mixedPort)
		if err != nil {
			return nil, err
		}

		dialer, err := proxy.FromURL(proxyURI, proxy.Direct)
		if err != nil {
			return nil, err
		}

		transport := &http.Transport{
			Dial: dialer.Dial,
		}

		client = &http.Client{
			Transport: transport,
			Timeout:   30 * time.Second,
		}
	} else {
		client = &http.Client{
			Timeout: 30 * time.Second,
		}
	}

	resp, err := client.Get(suburl)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	decodedBody, err := DecodeBase64URLSafe(string(body))
	if err != nil {
		return nil, err
	}

	return decodedBody, nil
}

// 订阅处理
func Handler(url string) ([]byte, error) {

	req, err := FetchSubscribe(url)
	if err != nil {
		return nil, err
	}

	outbounds, err := native.NewNativeURIParser(string(req))
	if err != nil {
		return nil, err
	}

	// 重复节点处理
	nameCount := make(map[string]int)
	result := make([]option.Outbound, len(outbounds))

	for i, outbound := range outbounds {
		nameCount[outbound.Tag]++
		if nameCount[outbound.Tag] == 1 {
			result[i] = outbound
		} else {
			outbound.Tag = fmt.Sprintf("%s%d", outbound.Tag, nameCount[outbound.Tag])
			result[i] = outbound
		}
	}

	var tags []string
	for _, t := range result {
		tags = append(tags, t.Tag)
	}

	template, err := os.ReadFile("./resources/template/template.json")
	if err != nil {
		return nil, err
	}

	var options option.Options

	options.UnmarshalJSON(template)

	// 添加出站
	options.Outbounds = append(options.Outbounds, result...)

	// 添加出站标签
	for index := range options.Outbounds {
		i := &options.Outbounds[index]
		switch i.Type {
		case "selector":
			i.SelectorOptions.Outbounds = append(i.SelectorOptions.Outbounds, tags...)
		case "urltest":
			i.URLTestOptions.Outbounds = append(i.URLTestOptions.Outbounds, tags...)
		}
	}

	config, err := json.MarshalIndent(options, "", "  ")
	if err != nil {
		return nil, err
	}

	return config, nil
}

package common

import (
	"encoding/base64"
	"io"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

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

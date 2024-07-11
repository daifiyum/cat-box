package parser

import (
	"io"
	"net/http"
	"net/url"
	"time"

	"golang.org/x/net/proxy"
)

func ProxyHttp(link string) ([]byte, error) {
	proxyURI, _ := url.Parse("socks5://127.0.0.1:8888")

	dialer, _ := proxy.FromURL(proxyURI, proxy.Direct)

	transport := &http.Transport{
		Dial: dialer.Dial,
	}

	client := &http.Client{
		Transport: transport,
		Timeout:   30 * time.Second,
	}

	resp, err := client.Get(link)
	if err != nil {
		return nil, err
	}

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

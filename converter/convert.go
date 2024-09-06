package converter

import (
	"encoding/json"
	"errors"
	"strings"

	"github.com/daifiyum/cat-box/converter/model"
	"github.com/daifiyum/cat-box/converter/parser"
)

func ConvertCProxyToSProxy(proxy string) (model.Outbound, error) {
	for prefix, parseFunc := range parser.ParserMap {
		if strings.HasPrefix(proxy, prefix) {
			proxy, err := parseFunc(proxy)
			if err != nil {
				return model.Outbound{}, err
			}
			return proxy, nil
		}
	}
	return model.Outbound{}, errors.New("unknown proxy format")
}

func ConvertCProxyToJson(proxy string) (string, error) {
	sProxy, err := ConvertCProxyToSProxy(proxy)
	if err != nil {
		return "", err
	}
	data, err := json.Marshal(&sProxy)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

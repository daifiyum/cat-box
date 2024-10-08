package converter

import (
	"encoding/json"
	"os"

	C "github.com/daifiyum/cat-box/converter/common"
	"github.com/daifiyum/cat-box/converter/native"
	"github.com/daifiyum/cat-box/utils"
	"github.com/sagernet/sing-box/option"
)

func Handler(url string) ([]byte, error) {

	req, err := C.FetchSubscribe(url)
	if err != nil {
		utils.LogError("fetchSubscribe errors")
		return nil, err
	}

	outbounds, err := native.NewNativeURIParser(string(req))
	if err != nil {
		utils.LogError("convertSubscribe error")
		return nil, err
	}

	var tags []string
	for _, t := range outbounds {
		tags = append(tags, t.Tag)
	}

	template, err := os.ReadFile("./resources/template/template.json")
	if err != nil {
		return nil, err
	}

	var options option.Options

	options.UnmarshalJSON(template)

	// 添加出站
	options.Outbounds = append(options.Outbounds, outbounds...)

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

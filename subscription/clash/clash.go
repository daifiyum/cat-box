package clash

import (
	"context"

	P "github.com/sagernet/serenity/subscription/parser"
	"github.com/sagernet/sing-box/option"
)

func ParseClashSubscription(content string) ([]option.Outbound, error) {
	outbounds, err := P.ParseClashSubscription(context.Background(), content)
	if err != nil {
		return nil, err
	}

	return outbounds, nil
}

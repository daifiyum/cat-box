package main

import (
	"os"

	S "github.com/daifiyum/cat-box/subscription"
)

func main() {
	data, err := S.Subscription("", "clash-verge/v2.4.0")
	if err != nil {
		panic(err)
	}

	err = os.WriteFile("subscription.json", []byte(data), 0644)
	if err != nil {
		panic(err)
	}
}

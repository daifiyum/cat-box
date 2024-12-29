package setup

import (
	"log"
	"os"
)

func Log() error {
	f, err := os.OpenFile("./app.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
	if err != nil {
		return err
	}
	log.SetOutput(f)
	return nil
}

package setup

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	U "github.com/daifiyum/cat-box/config"
)

// cat-box.exe -workspace=true
// cat-box.exe -workspace=false
// cat-box.exe -port=3000
// cat-box.exe -workspace=true -port=3000
func Commands() error {
	workspace := flag.Bool("workspace", false, "Enable workspace mode")
	port := flag.Int("port", 0, "Set the port number")
	flag.Parse()

	if *workspace {
		e, _ := os.Executable()
		err := os.Chdir(filepath.Dir(e))
		if err != nil {
			return err
		}
	}

	if *port != 0 {
		U.Port = fmt.Sprint(*port)
	}
	return nil
}

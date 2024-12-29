package setup

import (
	"fmt"
	"path/filepath"

	W "github.com/daifiyum/cat-box/app/windows"
)

func Aumid() error {
	iconURL, err := filepath.Abs("./resources/icons/box.ico")
	if err != nil {
		return fmt.Errorf("failed to get absolute path for icon: %w", err)
	}

	err = W.RegisterAUMID("cat-box", "cat-box", iconURL)
	if err != nil {
		return fmt.Errorf("failed to register AUMID: %w", err)
	}

	err = W.SetAUMID("cat-box")
	if err != nil {
		return fmt.Errorf("failed to set AUMID: %w", err)
	}

	return nil
}

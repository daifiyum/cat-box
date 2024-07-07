package singbox

import (
	"os/exec"
	"syscall"

	"github.com/daifiyum/cat-box/utils"
)

var (
	cmd *exec.Cmd
)

func CheckConfig() error {
	cmd := exec.Command("./resources/core/sing-box.exe", "check", "-c", "./resources/core/config.json")
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	output, err := cmd.CombinedOutput()
	if err != nil {
		output = output[21:]
		utils.ShowToast("cat-box", "错误", string(output))
		return err
	}
	return nil
}

func Start() error {
	if cmd != nil {
		return nil
	}
	err := GenerateConfig()
	if err != nil {
		return err
	}
	err = CheckConfig()
	if err != nil {
		return err
	}
	cmd = exec.Command("./resources/core/sing-box.exe", "run", "-c", "./resources/core/config.json")
	// cmd.Stdout = os.Stdout
	// cmd.Stderr = os.Stderr
	cmd.SysProcAttr = &syscall.SysProcAttr{
		HideWindow: true,
	}
	err = cmd.Start()
	if err != nil {
		return err
	}
	return nil
}

func Stop() error {
	if cmd == nil {
		return nil
	}

	err := cmd.Process.Kill()
	if err != nil {
		return err
	}

	cmd.Wait()

	cmd = nil
	return nil
}

func Reload() error {
	err := Stop()
	if err != nil {
		return err
	}
	err = Start()
	if err != nil {
		return err
	}
	return nil
}

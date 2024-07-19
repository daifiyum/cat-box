package singbox

import (
	"os/exec"
	"sync"
	"syscall"

	"github.com/daifiyum/cat-box/utils"
	"golang.org/x/sys/windows"
)

var (
	cmd          *exec.Cmd
	closeCoreMux sync.Mutex
)

func CheckConfig() error {
	cmdcheck := exec.Command("./resources/core/sing-box.exe", "check", "-c", "./resources/core/config.json")
	cmdcheck.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	output, err := cmdcheck.CombinedOutput()
	if err != nil {
		output = output[21:]
		utils.ShowToast("cat-box", "错误", string(output))
		return err
	}
	return nil
}

func Start() error {
	err := Stop()
	if err != nil {
		return err
	}

	err = GenerateConfig()
	if err != nil {
		return err
	}
	err = CheckConfig()
	if err != nil {
		return err
	}

	cmd = exec.Command("./resources/core/sing-box.exe", "run", "-c", "./resources/core/config.json")
	cmd.SysProcAttr = &syscall.SysProcAttr{
		CreationFlags: windows.CREATE_UNICODE_ENVIRONMENT | windows.CREATE_NEW_PROCESS_GROUP,
		HideWindow:    true,
	}
	err = cmd.Start()
	if err != nil {
		return err
	}

	return nil
}

func Stop() error {
	closeCoreMux.Lock()
	defer closeCoreMux.Unlock()

	if cmd == nil || cmd.ProcessState != nil && cmd.ProcessState.Exited() {
		return nil
	}

	err := terminateProc(cmd.Process.Pid)
	if err != nil {
		return err
	}

	err = cmd.Wait()
	if err != nil {
		return err
	}

	cmd = nil
	return nil
}

func CheckCoreStatus() {
	if cmd == nil || cmd.ProcessState != nil && cmd.ProcessState.Exited() {
		utils.IsProxy = false
	} else {
		utils.IsProxy = true
	}
}

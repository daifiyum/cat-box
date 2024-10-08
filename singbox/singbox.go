package singbox

import (
	"encoding/json"
	"os"
	"os/exec"
	"sync"
	"syscall"

	"github.com/daifiyum/cat-box/subservice/database"
	"github.com/daifiyum/cat-box/subservice/models"
	"github.com/daifiyum/cat-box/utils"
	"github.com/sagernet/sing-box/option"
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

	err = modeSwitch()
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

func terminateProc(pid int) error {
	dll, err := windows.LoadDLL("kernel32.dll")
	if err != nil {
		return err
	}
	defer dll.Release()

	f, err := dll.FindProc("AttachConsole")
	if err != nil {
		return err
	}
	r1, _, err := f.Call(uintptr(pid))
	if r1 == 0 && err != syscall.ERROR_ACCESS_DENIED {
		return err
	}

	f, err = dll.FindProc("SetConsoleCtrlHandler")
	if err != nil {
		return err
	}
	r1, _, err = f.Call(0, 1)
	if r1 == 0 {
		return err
	}
	f, err = dll.FindProc("GenerateConsoleCtrlEvent")
	if err != nil {
		return err
	}
	r1, _, err = f.Call(windows.CTRL_BREAK_EVENT, uintptr(pid))
	if r1 == 0 {
		return err
	}
	return nil
}

func modeSwitch() error {
	isTun := utils.IsTun
	db := database.DB
	var subscriptions models.Subscriptions
	db.Where("active = ?", true).First(&subscriptions)

	var options option.Options

	options.UnmarshalJSON([]byte(subscriptions.Data))
	if !isTun {
		for index := range options.Inbounds {
			i := &options.Inbounds[index]
			if i.Type == "mixed" {
				i.MixedOptions.SetSystemProxy = true
			}
			if i.Type == "tun" {
				options.Inbounds = append(options.Inbounds[:index], options.Inbounds[index+1:]...)
				index--
			}
		}
	} else {
		for index := range options.Inbounds {
			i := &options.Inbounds[index]
			if i.Type == "mixed" {
				i.MixedOptions.SetSystemProxy = false
			}
		}
	}

	config, _ := json.MarshalIndent(options, "", "  ")
	err := os.WriteFile("./resources/core/config.json", []byte(config), 0666)
	if err != nil {
		return err
	}
	return nil
}

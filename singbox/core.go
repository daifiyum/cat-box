package singbox

import (
	"bufio"
	"log"
	"os/exec"
	"sync"
	"syscall"

	U "github.com/daifiyum/cat-box/common"
	"golang.org/x/sys/windows"
)

var (
	cmd *exec.Cmd
	wg  sync.WaitGroup
)

func Start() error {
	err := Stop()
	if err != nil {
		return err
	}

	err = CompareTemplate()
	if err != nil {
		return err
	}

	cmd = exec.Command("./resources/core/sing-box.exe", "run", "-c", "./resources/core/config.json")
	cmd.SysProcAttr = &syscall.SysProcAttr{
		CreationFlags: windows.CREATE_UNICODE_ENVIRONMENT | windows.CREATE_NEW_PROCESS_GROUP,
		HideWindow:    true,
	}

	stderr, err := cmd.StderrPipe()
	if err != nil {
		return err
	}

	err = cmd.Start()
	if err != nil {
		log.Println(err)
		return err
	}

	U.IsCoreRunning.Set(true)
	wg.Add(2)
	errScanner := bufio.NewScanner(stderr)
	go func() {
		defer wg.Done()
		for errScanner.Scan() {
			text := errScanner.Text()
			U.Broadcaster.Broadcast(text)
		}
	}()
	go func() {
		defer wg.Done()
		cmd.Wait()
		U.IsCoreRunning.Set(false)
	}()

	return nil
}

func Stop() error {
	// cmd未初始化或cmd.Start()出错，此时cmd和cmd.Process为空，则一定没有子进程
	if cmd == nil || cmd.Process == nil {
		return nil
	}
	if U.IsCoreRunning.Get() {
		err := terminateProc(cmd.Process.Pid)
		if err != nil {
			log.Println(err)
			return err
		}
	}
	wg.Wait()
	return nil
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

// 开关核心
func SwitchCore() {
	if U.IsCoreRunning.Get() {
		err := Stop()
		if err != nil {
			return
		}
		return
	}

	err := Start()
	if err != nil {
		return
	}
}

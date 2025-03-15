package singbox

import (
	"bufio"
	"log"
	"os/exec"
	"sync"
	"syscall"

	W "github.com/daifiyum/cat-box/app/windows"
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
	// cmd未初始化或cmd.Start()出错，cmd和cmd.Process为空，则没有子进程
	if cmd == nil || cmd.Process == nil {
		return nil
	}
	if U.IsCoreRunning.Get() {
		err := W.TerminateProc(cmd.Process.Pid)
		if err != nil {
			log.Println(err)
			return err
		}
	}
	wg.Wait()
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

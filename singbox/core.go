package singbox

import (
	"bufio"
	"io"
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
	if err := Stop(); err != nil {
		return err
	}

	if err := CompareTemplate(); err != nil {
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

	if err = cmd.Start(); err != nil {
		log.Println("Core start failed:", err)
		return err
	}

	U.IsCoreRunning.Set(true)

	wg.Add(2)
	go scanStderr(stderr)
	go waitProcess()

	return nil
}

func scanStderr(stderrPipe io.ReadCloser) {
	defer wg.Done()
	scanner := bufio.NewScanner(stderrPipe)
	for scanner.Scan() {
		U.Broadcaster.Broadcast(scanner.Text())
	}
}

func waitProcess() {
	defer wg.Done()
	if err := cmd.Wait(); err != nil {
		log.Println("Error waiting for the process to exit:", err)
	}
	U.IsCoreRunning.Set(false)
}

func Stop() error {
	if cmd == nil || cmd.Process == nil {
		return nil
	}
	if U.IsCoreRunning.Get() {
		if err := W.TerminateProc(cmd.Process.Pid); err != nil {
			log.Println("Error terminating process:", err)
			return err
		}
	}
	wg.Wait()
	return nil
}

func SwitchCore() {
	if U.IsCoreRunning.Get() {
		if err := Stop(); err != nil {
			log.Printf("Core stop failed: %v", err)
		}
	} else {
		if err := Start(); err != nil {
			log.Printf("Core start failed: %v", err)
		}
	}
}

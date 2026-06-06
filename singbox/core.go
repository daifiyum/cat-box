package singbox

import (
	"bufio"
	"io"
	"log"
	"os/exec"
	"sync"
	"syscall"
	"time"

	W "github.com/daifiyum/cat-box/app/windows"
	U "github.com/daifiyum/cat-box/common"
	"golang.org/x/sys/windows"
)

var (
	cmd         *exec.Cmd
	wg          sync.WaitGroup
	lifecycleMu sync.Mutex
)

func Start() error {
	lifecycleMu.Lock()
	defer lifecycleMu.Unlock()

	if err := stopLocked(); err != nil {
		return err
	}

	if err := CompareTemplate(); err != nil {
		return err
	}

	newCmd := exec.Command("./resources/core/sing-box.exe", "run", "-c", "./resources/core/config.json")
	newCmd.SysProcAttr = &syscall.SysProcAttr{
		CreationFlags: windows.CREATE_UNICODE_ENVIRONMENT | windows.CREATE_NEW_PROCESS_GROUP,
		HideWindow:    true,
	}

	stderr, err := newCmd.StderrPipe()
	if err != nil {
		return err
	}

	if err = newCmd.Start(); err != nil {
		log.Println("Core start failed:", err)
		return err
	}

	cmd = newCmd
	U.IsCoreRunning.Set(true)

	wg.Add(2)
	go scanStderr(stderr)
	go waitProcess(newCmd)

	return nil
}

func scanStderr(stderrPipe io.ReadCloser) {
	defer wg.Done()
	scanner := bufio.NewScanner(stderrPipe)
	for scanner.Scan() {
		U.Broadcaster.Broadcast(scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		log.Println("Error scanning core stderr:", err)
	}
}

func waitProcess(processCmd *exec.Cmd) {
	if err := processCmd.Wait(); err != nil {
		log.Println("Error waiting for the process to exit:", err)
	}
	U.IsCoreRunning.Set(false)
	wg.Done()

	lifecycleMu.Lock()
	if cmd == processCmd {
		cmd = nil
	}
	lifecycleMu.Unlock()
}

func Stop() error {
	lifecycleMu.Lock()
	defer lifecycleMu.Unlock()

	return stopLocked()
}

func stopLocked() error {
	if cmd == nil || cmd.Process == nil {
		return nil
	}

	processCmd := cmd
	if U.IsCoreRunning.Get() {
		if err := W.TerminateProc(processCmd.Process.Pid); err != nil {
			log.Println("Error terminating process:", err)
			if killErr := processCmd.Process.Kill(); killErr != nil {
				return killErr
			}
		}
	}

	done := make(chan struct{})
	go func() {
		wg.Wait()
		close(done)
	}()

	select {
	case <-done:
	case <-time.After(5 * time.Second):
		log.Println("Core stop timeout, killing process")
		if err := processCmd.Process.Kill(); err != nil {
			return err
		}
		<-done
	}

	if cmd == processCmd {
		cmd = nil
	}
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

package singbox

import (
	"bufio"
	"fmt"
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
	mu  sync.Mutex
)

func Start() error {
	if err := Stop(); err != nil {
		return fmt.Errorf("failed to stop existing process: %w", err)
	}

	if err := CompareTemplate(); err != nil {
		return fmt.Errorf("config template validation failed: %w", err)
	}

	mu.Lock()
	defer mu.Unlock()

	cmd = exec.Command("./resources/core/sing-box.exe", "run", "-c", "./resources/core/config.json")
	cmd.SysProcAttr = &syscall.SysProcAttr{
		CreationFlags: windows.CREATE_UNICODE_ENVIRONMENT | windows.CREATE_NEW_PROCESS_GROUP,
		HideWindow:    true,
	}

	stderr, err := cmd.StderrPipe()
	if err != nil {
		return fmt.Errorf("failed to get stderr pipe: %w", err)
	}

	if err := cmd.Start(); err != nil {
		stderr.Close()
		return fmt.Errorf("process start failed: %w", err)
	}

	U.IsCoreRunning.Set(true)

	wg = sync.WaitGroup{}
	wg.Add(2)

	go processStderr(stderr)

	go monitorProcess()

	return nil
}

func Stop() error {
	mu.Lock()
	defer mu.Unlock()
	return stopLocked()
}

func stopLocked() error {
	if cmd == nil || cmd.Process == nil {
		return nil
	}

	if U.IsCoreRunning.Get() {
		if err := W.TerminateProc(cmd.Process.Pid); err != nil {
			return fmt.Errorf("process termination failed: %w", err)
		}
	}

	wg.Wait()
	cmd = nil
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

func processStderr(stderrPipe io.ReadCloser) {
	defer wg.Done()
	defer stderrPipe.Close()

	scanner := bufio.NewScanner(stderrPipe)
	for scanner.Scan() {
		U.Broadcaster.Broadcast(scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		U.Broadcaster.Broadcast(fmt.Sprintf("[ERROR] Stderr read failed: %v", err))
	}
}

func monitorProcess() {
	defer wg.Done()

	err := cmd.Wait()
	U.IsCoreRunning.Set(false)

	exitCode := 0
	if exitErr, ok := err.(*exec.ExitError); ok {
		if status, ok := exitErr.Sys().(syscall.WaitStatus); ok {
			exitCode = status.ExitStatus()
		}
	}

	U.Broadcaster.Broadcast(fmt.Sprintf("[CORE] Process exited (Code: %d)", exitCode))
}

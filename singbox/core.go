package singbox

import (
	"bufio"
	"errors"
	"io"
	"log"
	"os"
	"os/exec"
	"sync"
	"syscall"
	"time"

	W "github.com/daifiyum/cat-box/app/windows"
	U "github.com/daifiyum/cat-box/common"
	"golang.org/x/sys/windows"
)

type coreProcess struct {
	cmd  *exec.Cmd
	done chan struct{}
}

var (
	core        *coreProcess
	lifecycleMu sync.Mutex
)

func Start() error {
	for {
		if err := Stop(); err != nil {
			return err
		}
		lifecycleMu.Lock()
		if core == nil {
			break
		}
		// Another Start acquired the lock first; stop its process before retrying.
		lifecycleMu.Unlock()
	}
	defer lifecycleMu.Unlock()
	if err := CompareTemplate(); err != nil {
		return err
	}

	cmd := exec.Command("./resources/core/sing-box.exe", "run", "-c", "./resources/core/config.json")
	cmd.SysProcAttr = &syscall.SysProcAttr{
		CreationFlags: windows.CREATE_UNICODE_ENVIRONMENT | windows.CREATE_NEW_PROCESS_GROUP,
		HideWindow:    true,
	}
	stderr, err := cmd.StderrPipe()
	if err != nil {
		return err
	}
	if err := cmd.Start(); err != nil {
		return err
	}

	process := &coreProcess{cmd: cmd, done: make(chan struct{})}
	core = process
	U.IsCoreRunning.Set(true)
	go process.wait(stderr)
	return nil
}

func (p *coreProcess) wait(stderr io.Reader) {
	defer close(p.done)

	// Finish reading before Wait closes the stderr pipe.
	scanner := bufio.NewScanner(stderr)
	for scanner.Scan() {
		U.Broadcaster.Broadcast(scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		log.Println("Error scanning core stderr:", err)
		// Keep draining after a scan error so the child cannot block on stderr.
		_, _ = io.Copy(io.Discard, stderr)
	}
	if err := p.cmd.Wait(); err != nil {
		log.Println("Error waiting for the process to exit:", err)
	}
	U.IsCoreRunning.Set(false)
}

func Stop() error {
	lifecycleMu.Lock()
	defer lifecycleMu.Unlock()
	if core == nil {
		return nil
	}
	select {
	case <-core.done:
		core = nil
		return nil
	default:
	}

	if err := W.TerminateProc(core.cmd.Process.Pid); err != nil {
		log.Println("Error terminating process:", err)
		if err := core.cmd.Process.Kill(); err != nil && !errors.Is(err, os.ErrProcessDone) {
			return err
		}
	}

	select {
	case <-core.done:
	case <-time.After(5 * time.Second):
		log.Println("Core stop timeout, killing process")
		if err := core.cmd.Process.Kill(); err != nil && !errors.Is(err, os.ErrProcessDone) {
			return err
		}
		<-core.done
	}
	core = nil
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

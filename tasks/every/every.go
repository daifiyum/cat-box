package every

import (
	"fmt"
	"strconv"
	"sync"
	"time"
)

type Task struct {
	duration   time.Duration
	taskFunc   func()
	timer      *time.Timer
	stopChan   chan struct{}
	updateChan chan time.Duration
	wg         sync.WaitGroup
}

func NewTask(interval string, task func()) (*Task, error) {
	duration, err := parseDuration(interval)
	if err != nil {
		return nil, err
	}

	return &Task{
		duration:   duration,
		taskFunc:   task,
		stopChan:   make(chan struct{}),
		updateChan: make(chan time.Duration),
	}, nil
}

func parseDuration(interval string) (time.Duration, error) {
	unitMap := map[byte]time.Duration{'s': time.Second, 'm': time.Minute, 'h': time.Hour, 'd': 24 * time.Hour}
	n := len(interval)
	if n < 2 {
		return 0, fmt.Errorf("invalid format: %s", interval)
	}

	value, err := strconv.Atoi(interval[:n-1])
	if err != nil || unitMap[interval[n-1]] == 0 {
		return 0, fmt.Errorf("invalid duration: %s", interval)
	}

	return time.Duration(value) * unitMap[interval[n-1]], nil
}

func (t *Task) Start() {
	t.wg.Add(1)
	go func() {
		defer t.wg.Done()

		t.timer = time.NewTimer(t.duration)

		for {
			select {
			case <-t.stopChan:
				t.timer.Stop()
				return
			case newDuration := <-t.updateChan:
				t.timer.Stop()
				t.duration = newDuration
				t.timer.Reset(t.duration)
			case <-t.timer.C:
				t.taskFunc()
				t.timer.Reset(t.duration)
			}
		}
	}()
}

func (t *Task) Stop() {
	close(t.stopChan)
	t.wg.Wait()
}

func (t *Task) UpdateInterval(interval string) error {
	newDuration, err := parseDuration(interval)
	if err != nil {
		return err
	}

	t.updateChan <- newDuration
	return nil
}

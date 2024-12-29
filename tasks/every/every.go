package every

import (
	"fmt"
	"strconv"
	"sync"
	"time"
)

type Task struct {
	duration   time.Duration      // 当前任务的执行间隔
	taskFunc   func()             // 定时任务的逻辑
	timer      *time.Timer        // 用于控制定时任务的 Timer
	stopChan   chan struct{}      // 停止信号
	updateChan chan time.Duration // 更新间隔信号
	wg         sync.WaitGroup     // 用于等待任务退出
}

// 创建一个新的定时任务
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

// 解析时间间隔
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

// 启动定时任务
func (t *Task) Start() {
	t.wg.Add(1)
	go func() {
		defer t.wg.Done()

		// 初始化定时器
		t.timer = time.NewTimer(t.duration)

		for {
			select {
			case <-t.stopChan:
				// 停止任务
				t.timer.Stop()
				return
			case newDuration := <-t.updateChan:
				// 更新间隔时间并重置定时器
				t.timer.Stop()
				t.duration = newDuration
				t.timer.Reset(t.duration)
			case <-t.timer.C:
				// 执行任务
				t.taskFunc()

				// 重置定时器
				t.timer.Reset(t.duration)
			}
		}
	}()
}

// 停止任务
func (t *Task) Stop() {
	close(t.stopChan)
	t.wg.Wait()
}

// 更新任务间隔时间
func (t *Task) UpdateInterval(interval string) error {
	newDuration, err := parseDuration(interval)
	if err != nil {
		return err
	}

	// 发送更新信号
	t.updateChan <- newDuration
	return nil
}

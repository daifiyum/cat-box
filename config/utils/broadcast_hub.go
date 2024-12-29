package utils

import "sync"

type BroadcastHub struct {
	subscribers map[chan string]struct{}
	mu          sync.RWMutex
}

func NewBroadcaster() *BroadcastHub {
	return &BroadcastHub{
		subscribers: make(map[chan string]struct{}),
	}
}

func (b *BroadcastHub) Subscribe(bufferSize int) chan string {
	logChan := make(chan string, bufferSize)
	b.mu.Lock()
	b.subscribers[logChan] = struct{}{}
	b.mu.Unlock()
	return logChan
}

func (b *BroadcastHub) Unsubscribe(logChan chan string) {
	b.mu.Lock()
	if _, ok := b.subscribers[logChan]; ok {
		delete(b.subscribers, logChan)
		close(logChan)
	}
	b.mu.Unlock()
}

func (b *BroadcastHub) Broadcast(logLine string) {
	b.mu.RLock()
	defer b.mu.RUnlock()
	for logChan := range b.subscribers {
		select {
		case logChan <- logLine:
		default: // 缓冲区满了丢弃
		}
	}
}

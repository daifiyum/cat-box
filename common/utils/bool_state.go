package utils

import "sync"

type BoolState struct {
	mu        sync.RWMutex
	value     bool
	listeners []func(bool)
}

func (v *BoolState) Set(newValue bool) {
	v.mu.Lock()
	v.value = newValue
	listeners := make([]func(bool), len(v.listeners))
	copy(listeners, v.listeners)
	v.mu.Unlock()

	for _, listener := range listeners {
		listener(newValue)
	}
}

func (v *BoolState) Get() bool {
	v.mu.RLock()
	defer v.mu.RUnlock()
	return v.value
}

func (v *BoolState) Watch(listener func(bool)) {
	v.mu.Lock()
	defer v.mu.Unlock()
	v.listeners = append(v.listeners, listener)
}

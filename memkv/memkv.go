package memkv

import (
	"errors"
	"reflect"
	"strings"
	"sync"
)

type MemKV struct {
	data map[string]any
	lock *sync.RWMutex
}

func New() *MemKV {
	return &MemKV{
		data: make(map[string]any),
		lock: &sync.RWMutex{},
	}
}

func (m *MemKV) Set(key string, value any) {
	m.lock.Lock()
	defer m.lock.Unlock()

	m.data[key] = value
}

func (m *MemKV) Equal(key string, value any) bool {
	m.lock.RLock()
	defer m.lock.RUnlock()

	v, ok := m.data[key]
	if !ok {
		return false
	}
	return reflect.DeepEqual(value, v)
}

func (m *MemKV) Exists(key string) bool {
	m.lock.RLock()
	defer m.lock.Unlock()

	v, ok := m.data[key]

	return ok && v != nil
}

var ErrKeyNotExist = errors.New("key is not exist")

func (m *MemKV) Get(key string) (value any, err error) {
	m.lock.RLock()
	defer m.lock.RUnlock()

	v, ok := m.data[key]
	if !ok {
		return nil, ErrKeyNotExist
	}
	return v, nil
}

type KV struct {
	Key   string
	Value any
}

func (m *MemKV) Gets(keys ...string) (items []KV, err error) {
	m.lock.RLock()
	defer m.lock.RUnlock()

	items = make([]KV, 0)
	for _, key := range keys {
		value, ok := m.data[key]
		if !ok {
			return nil, ErrKeyNotExist
		}
		items = append(items, KV{Key: key, Value: value})
	}
	return items, nil
}

func (m *MemKV) GetWithPrefix(prefix string) (items []KV) {
	m.lock.RLock()
	defer m.lock.RUnlock()

	items = make([]KV, 0)
	for key, value := range m.data {
		if strings.HasPrefix(key, prefix) {
			items = append(items, KV{Key: key, Value: value})
		}
	}
	return items
}

func (m *MemKV) FuncMaps() map[string]any {
	return map[string]any{
		"get":   m.Get,
		"gets":  m.Gets,
		"getp":  m.GetWithPrefix,
		"exist": m.Exists,
	}
}

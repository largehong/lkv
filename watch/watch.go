package watch

import (
	"errors"
	"sync"
)

type Client interface {
	Get() ([]KV, error)
}

type KV struct {
	Key   string
	Value string
}

type Factory func(config any, prefixes []string, callback func(...KV)) (Client, error)

var (
	factories = make(map[string]Factory)
	lock      = &sync.Mutex{}
)

func Register(name string, fn Factory) {
	lock.Lock()
	defer lock.Unlock()

	if _, ok := factories[name]; ok {
		panic("dup client: " + name)
	}

	factories[name] = fn
}

func New(name string, config any, prefixes []string, callback func(...KV)) (Client, error) {
	fn, ok := factories[name]
	if !ok {
		return nil, errors.New("not found client factory: " + name)
	}
	return fn(config, prefixes, callback)
}

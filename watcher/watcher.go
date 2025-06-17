package watcher

import "sync"

type Watcher interface {
	Add(keys ...string)
	Watch() <-chan []KV
}

type KV struct {
	Key   string
	Value string
}

type Factory func(config map[string]any) (Watcher, error)

var (
	factories = make(map[string]Factory)
	lock      = &sync.Mutex{}
)

func Register(name string, fn Factory) {
	lock.Lock()
	defer lock.Unlock()

	_, ok := factories[name]
	if ok {
		panic("dup watcher: " + name)
	}

	factories[name] = fn
}

func Get(name string) (Factory, bool) {
	fn, ok := factories[name]
	return fn, ok
}

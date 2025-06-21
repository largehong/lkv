package engine

import (
	"encoding/json"
	"slices"
	"strings"
	"sync"
	"time"

	"github.com/largehong/lkv/memkv"
	"github.com/largehong/lkv/processor"
	"github.com/largehong/lkv/watch"
)

type Engine struct {
	processors map[string][]*processor.Processor
	lock       *sync.Mutex
	ch         chan []watch.KV
	kv         *memkv.MemKV
	interval   int
	max        int
}

func New(kv *memkv.MemKV, max int, interval int) *Engine {
	return &Engine{
		processors: make(map[string][]*processor.Processor),
		lock:       &sync.Mutex{},
		ch:         make(chan []watch.KV, max),
		kv:         kv,
		interval:   interval,
		max:        max,
	}
}

func (engine *Engine) Callback(kvs ...watch.KV) {
	engine.ch <- kvs
}

func (engine *Engine) Register(prefix string, p *processor.Processor) {
	engine.lock.Lock()
	defer engine.lock.Unlock()

	processors, ok := engine.processors[prefix]
	if !ok {
		processors = make([]*processor.Processor, 0)
	}
	if !slices.Contains(processors, p) {
		processors = append(processors, p)
		engine.processors[prefix] = processors
	}
}

func (engine *Engine) Run() {
	ticker := time.NewTicker(time.Duration(engine.interval) * time.Second)
	defer ticker.Stop()

	buf := make([]watch.KV, 0)
	for {
		select {
		case kvs := <-engine.ch:
			buf = append(buf, kvs...)
			if len(buf) > engine.max {
				engine.handle(buf...)
				buf = make([]watch.KV, 0)
			}
		case <-ticker.C:
			if len(buf) > 0 {
				engine.handle(buf...)
				buf = make([]watch.KV, 0)
			}
		}
	}
}

func (engine *Engine) Once(kvs ...watch.KV) {
	engine.handle(kvs...)
}

func (engine *Engine) handle(kvs ...watch.KV) {
	need := make([]*processor.Processor, 0)

	for _, kv := range kvs {
		//更新本地存储
		if kv.Value == "" {
			engine.kv.Del(kv.Key)
		} else {
			var v any
			if err := json.Unmarshal([]byte(kv.Value), &v); err == nil {
				engine.kv.Set(kv.Key, v)
			} else {
				engine.kv.Set(kv.Key, kv.Value)
			}
		}

		//获取订阅变更的processor，同时进行去重，避免多次通知
		for prefix, processors := range engine.processors {
			if strings.HasPrefix(kv.Key, prefix) {
				for _, p := range processors {
					if !slices.Contains(need, p) {
						need = append(need, p)
					}
				}
			}
		}
	}

	for _, p := range need {
		p.Redenering()
	}
}

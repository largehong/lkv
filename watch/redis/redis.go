package redis

import "github.com/largehong/lkv/watch"

type Redis struct {
	callback func(...watch.KV)
	prefixes []string
}

func init() {
	watch.Register("redis", New)
}

func New(config any, prefixes []string, callback func(...watch.KV)) (watch.Client, error) {
	return nil, nil
}

func (r *Redis) Get() ([]watch.KV, error) { return nil, nil }

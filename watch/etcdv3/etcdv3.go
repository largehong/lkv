package etcdv3

import (
	"context"

	"github.com/largehong/lkv/watch"
	"github.com/mitchellh/mapstructure"
	etcdv3 "go.etcd.io/etcd/client/v3"
)

type Etcdv3 struct {
	prefixes []string
	callback func(...watch.KV)
	client   *etcdv3.Client
}

type Config struct {
	Endpoints []string `yaml:"endpoints" mapstructure:"endpoints"`
}

func init() {
	watch.Register("etcdv3", New)
}

func New(config any, prefixes []string, callback func(...watch.KV)) (watch.Client, error) {
	var c Config
	if err := mapstructure.Decode(config, &c); err != nil {
		return nil, err
	}
	client, err := etcdv3.New(etcdv3.Config{
		Endpoints: c.Endpoints,
	})
	if err != nil {
		return nil, err
	}
	cli := &Etcdv3{
		client:   client,
		prefixes: prefixes,
		callback: callback,
	}
	for _, prefix := range prefixes {
		go cli.watch(prefix)
	}
	return cli, nil
}

func (etcd *Etcdv3) Get() ([]watch.KV, error) {
	kvs := make([]watch.KV, 0)
	for _, prefix := range etcd.prefixes {
		resp, err := etcd.client.Get(context.TODO(), prefix, etcdv3.WithPrefix())
		if err != nil {
			return nil, err
		}
		for _, kv := range resp.Kvs {
			kvs = append(kvs, watch.KV{
				Key:   string(kv.Key),
				Value: string(kv.Value),
			})
		}
	}
	return kvs, nil
}

func (etcd *Etcdv3) watch(prefix string) {
	ch := etcd.client.Watch(context.TODO(), prefix, etcdv3.WithPrefix())
	for resp := range ch {
		for _, event := range resp.Events {
			etcd.callback(watch.KV{
				Key:   string(event.Kv.Key),
				Value: string(event.Kv.Value),
			})
		}
	}
}

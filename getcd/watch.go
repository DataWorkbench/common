package getcd

import (
	"context"
	"time"

	"github.com/DataWorkbench/glog"
	"go.etcd.io/etcd/api/v3/mvccpb"
	etcdv3 "go.etcd.io/etcd/client/v3"
)

// RetryWatch do watch the specified key(prefix) and auto retry when etcd error.
func RetryWatch(ctx context.Context, cli *Client, key string, handler func(ctx context.Context, eventType EventType, kv *KeyValue)) {
	// new logger.
	nl := glog.FromContext(ctx)

	var sleep bool

	revision := int64(0)
	more := true
	for more {
		if sleep {
			// Sleep to prevents died loop.
			time.Sleep(time.Millisecond * 100)
		}
		sleep = true

		var err error
		var resp *etcdv3.GetResponse
		resp, err = cli.Get(ctx, key, etcdv3.WithPrefix(), etcdv3.WithLimit(100), etcdv3.WithRev(revision))
		if err != nil {
			if err == context.Canceled {
				nl.Warn().Msg("etcd: context canceled when get exists key, return now").Fire()
				return
			}
			nl.Error().Msg("etcd: get exists keys failed, retry after 10s").String("key", key).Error("error", err).Fire()
			time.Sleep(time.Second * 10)
			continue
		}

		for _, kv := range resp.Kvs {
			handler(ctx, EventPUT, kv)
		}

		revision = resp.Header.Revision + 1
		more = resp.More
	}

	watchChan := cli.Watch(ctx, key, etcdv3.WithPrefix(), etcdv3.WithRev(revision))
LOOP:
	for {
		select {
		case resp := <-watchChan:
			for _, event := range resp.Events {
				switch event.Type {
				case mvccpb.PUT:
					handler(ctx, EventPUT, event.Kv)
				case mvccpb.DELETE:
					handler(ctx, EventDELETE, event.Kv)
				}
			}
		case <-ctx.Done():
			// done
			break LOOP
		}
	}
}

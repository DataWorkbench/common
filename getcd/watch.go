package getcd

import (
	"context"
	"time"

	"github.com/DataWorkbench/glog"
	"go.etcd.io/etcd/api/v3/mvccpb"
	etcdv3 "go.etcd.io/etcd/client/v3"
)

// RetryWatch putCallback func(k, v []byte), delCallback func(k, v []byte)
func RetryWatch(ctx context.Context, cli *Client, key string, handler func(eventType EventType, kv *KeyValue)) {
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
			nl.Error().Error("etcd: get exists keys error", err).String("key", key).Fire()
			continue
		}

		for _, kv := range resp.Kvs {
			handler(EventPUT, kv)
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
					handler(EventPUT, event.Kv)
				case mvccpb.DELETE:
					handler(EventDELETE, event.Kv)
				}
			}
		case <-ctx.Done():
			// done
			break LOOP
		}
	}
}

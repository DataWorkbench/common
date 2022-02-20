package getcd

import (
	"context"

	"github.com/DataWorkbench/glog"
	"go.etcd.io/etcd/client/v3/concurrency"
)

// NewMutex is a wrapper for create an etcd mutex by giving key prefix.
func NewMutex(ctx context.Context, cli *Client, key string) (mutex *Mutex, err error) {
	nl := glog.FromContext(ctx)

	var session *concurrency.Session

	nl.Debug().Msg("etcd: creating a session for mutex").String("key", key).Fire()
	session, err = concurrency.NewSession(cli)
	if err != nil {
		nl.Debug().Error("etcd: create session error", err).Fire()
		return
	}
	mutex = concurrency.NewMutex(session, key)
	return
}

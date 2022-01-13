package getcd

import (
	"go.etcd.io/etcd/api/v3/mvccpb"
	etcdv3 "go.etcd.io/etcd/client/v3"
	"go.etcd.io/etcd/client/v3/concurrency"
)

type Client = etcdv3.Client

type EventType = mvccpb.Event_EventType

type KeyValue = mvccpb.KeyValue

type Mutex = concurrency.Mutex

const (
	EventPUT    = mvccpb.PUT
	EventDELETE = mvccpb.DELETE
)

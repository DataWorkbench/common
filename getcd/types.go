package getcd

import (
	"go.etcd.io/etcd/api/v3/mvccpb"
	etcdv3 "go.etcd.io/etcd/client/v3"
)

type Client = etcdv3.Client

type EventType = mvccpb.Event_EventType

type KeyValue = mvccpb.KeyValue

const (
	EventPUT    = mvccpb.PUT
	EventDELETE = mvccpb.DELETE
)

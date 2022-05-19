package iaas

import (
	"errors"
)

var (
	ErrUserNotExists      = errors.New("user_not_exists")
	ErrAccessKeyNotExists = errors.New("access_key_not_exists")
	ErrVXNetNotExists     = errors.New("vxnet_not_exists")
	ErrRouterNotExists    = errors.New("router_not_exists")
)

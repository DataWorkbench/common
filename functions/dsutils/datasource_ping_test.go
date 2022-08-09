package dsutils

import (
	"testing"

	"github.com/DataWorkbench/gproto/xgo/types/pbmodel/pbdatasource"
	"github.com/stretchr/testify/require"
)

func Test_pingClickHouse(t *testing.T) {
	u1 := &pbdatasource.ClickHouseURL{
		Host:     "139.198.41.197",
		Port:     8123,
		User:     "root",
		Password: "zhu88jie",
		Database: "demo",
	}
	err := pingClickHouse(u1)
	require.Nil(t, err)

	u1.Port = 18123
	err = pingClickHouse(u1)
	require.NotNil(t, err)
}

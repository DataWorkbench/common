package idgeneratorv2

import (
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNew(t *testing.T) {
	prefix := "wks-"
	g := New(prefix)
	require.Equal(t, g.prefix, prefix)
	require.NotNil(t, g.worker)
}

func TestIDGeneratorV2_Take(t *testing.T) {
	prefix := "wks-"
	g := New(prefix)

	id, err := g.Take()
	require.Nil(t, err, "%+v", err)
	require.True(t, strings.HasPrefix(id, prefix), id)
	require.Equal(t, len(id), 20)
}

func TestIDGeneratorV2_TakeMany(t *testing.T) {
	g := New("wks-")
	for i := 0; i < 10; i++ {
		id, err := g.Take()
		require.Nil(t, err, "%+v", err)
		require.Equal(t, len(id), 20)
		fmt.Println(id)
	}
}

func TestIDGeneratorV2_TakeUnique(t *testing.T) {
	g := New("wks-")

	n := 10000000
	idMap := make(map[string]struct{})

	var lastId string
	_ = lastId
	for i := 0; i < n; i++ {
		id, err := g.Take()
		require.Nil(t, err, "%+v", err)
		require.Equal(t, len(id), 20)
		idMap[id] = struct{}{}

		//if lastId != "" {
		//	require.True(t, id > lastId)
		//}
		//lastId = id
	}

	require.Equal(t, len(idMap), n)
}

func BenchmarkGeneratorV2_Take(b *testing.B) {
	generator := New("wks-")

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			id, err := generator.Take()
			_ = id
			_ = err
		}
	})
}

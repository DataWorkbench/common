package idgenerator

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

func TestRandomInstanceID(t *testing.T) {
	for i := 0; i < 10; i++ {
		id := randomInstanceID()
		fmt.Println(id)
	}
}

func TestIDGenerator_Take(t *testing.T) {
	prefix := "wks-"
	g := New(prefix)

	t.Run("Take", func(t *testing.T) {
		id, err := g.Take()
		require.Nil(t, err, "%+v", err)
		require.True(t, strings.HasPrefix(id, prefix), id)
		require.Equal(t, len(id), 20)
	})
}

func TestIDGenerator_TakeMany(t *testing.T) {
	g := New("wks-")
	t.Run("Take2", func(t *testing.T) {
		for i := 0; i < 10; i++ {
			id, err := g.Take()
			require.Nil(t, err, "%+v", err)
			require.Equal(t, len(id), 20)
			fmt.Println(id)
		}
	})
}

func TestIDGenerator_UniqueTake(t *testing.T) {
	g := New("wks-")

	n := 1000000
	t.Run("Take", func(t *testing.T) {
		idMap := make(map[string]struct{})
		for i := 0; i < n; i++ {
			id, err := g.Take()
			require.Nil(t, err, "%+v", err)
			require.Equal(t, len(id), 20)
			idMap[id] = struct{}{}
		}
		require.Equal(t, len(idMap), n)
	})
}

func BenchmarkGenerator_Take(b *testing.B) {
	generator := New("wks-")

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			id, err := generator.Take()
			_ = id
			_ = err
		}
	})
}

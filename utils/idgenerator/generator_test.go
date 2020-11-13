package idgenerator

import (
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

func TestIDGenerator_Take(t *testing.T) {
	prefix := "wks-"
	g := New(prefix)

	id, err := g.Take()
	require.Nil(t, err, "%+v", err)
	require.True(t, strings.HasPrefix(id, prefix), id)
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

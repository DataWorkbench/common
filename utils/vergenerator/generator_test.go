package vergenerator

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNew(t *testing.T) {
	g := New()
	require.NotNil(t, g.worker)
}

func TestIDGenerator_Take(t *testing.T) {
	g := New()

	id, err := g.Take()
	require.Nil(t, err, "%+v", err)
	require.Greater(t, id, int64(0))
}

func BenchmarkGenerator_Take(b *testing.B) {
	generator := New()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			id, err := generator.Take()
			_ = id
			_ = err
		}
	})
}

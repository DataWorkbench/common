package qerror

import (
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"
)

func TestTypeEqual(t *testing.T) {
	a := ResourceNotExists.Format("aaaaaa")
	b := ResourceNotExists.Format("bbbbb")
	c := ResourceAlreadyExists

	require.True(t, errors.Is(a, b))
	require.False(t, errors.Is(a, c))
}

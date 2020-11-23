package qerror

import (
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"
)

func TestTypeEqual(t *testing.T) {
	a := WithDesc(NotExists, "aaaaaa")
	b := WithDesc(NotExists, "bbbbbb")
	c := AlreadyExists

	require.True(t, errors.Is(a, b))
	require.False(t, errors.Is(a, c))

}

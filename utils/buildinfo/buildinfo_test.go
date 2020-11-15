package buildinfo

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestJSONString(t *testing.T) {
	s := JSONString
	m := make(map[string]string)

	err := json.Unmarshal([]byte(s), &m)
	require.Nil(t, err, "%+v", err)
	require.Equal(t, m, MapValue)
}

func TestSingleString(t *testing.T) {
	s := SingleString
	m := make(map[string]string)

	l1 := strings.Split(s, " ")

	for i := range l1 {
		l2 := strings.Split(l1[i], "=")
		m[l2[0]] = l2[1]
	}

	require.Equal(t, m, MapValue)
}

func TestMultiString(t *testing.T) {
	s := MultiString
	require.NotEqual(t, s, "")
}

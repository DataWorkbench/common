package password

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func Test_Password(t *testing.T) {
	password := "zhu88jie"
	encodePassword, err := Encode(password)
	require.Nil(t, err, "%+v", err)
	require.NotEmpty(t, encodePassword)
	check := Check(password, encodePassword)
	require.True(t, check)
}

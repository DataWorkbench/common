package schargs

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_newDecoder(t *testing.T) {
	data := "xxx"
	d := newDecoder(data)
	require.NotNil(t, d)
	require.Equal(t, d.off, 0)
	require.Equal(t, d.data, data)
}

func Test_decoder_read(t *testing.T) {
	var b byte
	var ok bool

	data := "xyz"
	d := newDecoder(data)

	b, ok = d.read()
	require.True(t, ok)
	require.Equal(t, b, uint8('x'))

	b, ok = d.read()
	require.True(t, ok)
	require.Equal(t, b, uint8('y'))

	b, ok = d.read()
	require.True(t, ok)
	require.Equal(t, b, uint8('z'))

	b, ok = d.read()
	require.False(t, ok)
	require.Equal(t, b, uint8(0))
}

func Test_decoder_peek(t *testing.T) {
	var b byte
	var ok bool

	data := "xyz"
	d := newDecoder(data)

	// ------
	b, ok = d.peek()
	require.True(t, ok)
	require.Equal(t, b, uint8('x'))
	b, ok = d.read()
	require.True(t, ok)
	require.Equal(t, b, uint8('x'))

	// ------
	b, ok = d.peek()
	require.True(t, ok)
	require.Equal(t, b, uint8('y'))
	b, ok = d.read()
	require.True(t, ok)
	require.Equal(t, b, uint8('y'))

	// ------
	b, ok = d.peek()
	require.True(t, ok)
	require.Equal(t, b, uint8('z'))
	b, ok = d.read()
	require.True(t, ok)
	require.Equal(t, b, uint8('z'))

	// ------
	b, ok = d.peek()
	require.False(t, ok)
	require.Equal(t, b, uint8(0))
}

func Test_decoder_extractVariable1(t *testing.T) {
	data := "Test for ${var1} ${var2} $var3} ${var4}"
	d := newDecoder(data)
	keys, err := d.extractVariables()
	require.Nil(t, err)
	fmt.Println(keys)
}

func Test_decoder_extractVariable2(t *testing.T) {
	t.Run("variable name is empty", func(t *testing.T) {
		data := "Test for ${}"
		d := newDecoder(data)
		_, err := d.extractVariables()
		require.NotNil(t, err)
		require.Equal(t, err, ErrVariableIsEmpty)
	})

	t.Run("variable no end keyword", func(t *testing.T) {
		data := "Test for ${var1"
		d := newDecoder(data)
		_, err := d.extractVariables()
		require.NotNil(t, err)
		require.Equal(t, err, ErrInvalidVariableFormat)
	})

	t.Run("variable nested", func(t *testing.T) {
		data := "Test for ${var1${var2}}"
		d := newDecoder(data)
		_, err := d.extractVariables()
		require.NotNil(t, err)
		require.Equal(t, err, ErrUnsupportedVariableNested)
	})
}

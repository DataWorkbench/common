package schargs

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDecoder_ExtractVariables(t *testing.T) {
	data := `
		SELECT * FROM '${var1}';
		SELECT * FROM '${var2}';
		SELECT * FROM '${var3}';
		SELECT * FROM '${var4}';
		SELECT * FROM '$var5}';
		SELECT * FROM '{var6}';
		SELECT * FROM 'var7$';
		SELECT * FROM '$$var8';
	`
	keys, err := ExtractVariables(data)
	require.Nil(t, err)
	require.Equal(t, keys, []string{"var1", "var2", "var3", "var4"})
}

func TestDecoder_Encode(t *testing.T) {
	data := `
		SELECT * FROM '${var1}';
		SELECT * FROM '${var2}';
		SELECT * FROM '${var3}';
		SELECT * FROM '${var4}';
	`
	valueMap := map[string]string{
		"var1": "table1",
		"var2": "table2",
		"var3": "table3",
		"var4": "table4",
	}

	expectedData := `
		SELECT * FROM 'table1';
		SELECT * FROM 'table2';
		SELECT * FROM 'table3';
		SELECT * FROM 'table4';
	`

	nd, err := Encode(data, valueMap)
	require.Nil(t, err)
	require.Equal(t, nd, expectedData)
}

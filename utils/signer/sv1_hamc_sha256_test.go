package signer

import (
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
	"net/http"
	"net/url"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_BuildSV1Authorization(t *testing.T) {
	ak := "access-key-id"
	sk := "secret-access-key"

	body := []byte(`{"x2": 3, "x1": 2}`)

	h := md5.New()
	h.Write(body)
	b := h.Sum(nil)
	md5Hex := make([]byte, hex.EncodedLen(len(b)))
	hex.Encode(md5Hex, b)
	contentMd5 := base64.StdEncoding.EncodeToString(md5Hex)

	method := http.MethodPost
	urlPath := "/v1/workspace"

	query := make(url.Values)
	query.Add("limit", "100")
	query.Add("offset", "20")
	query.Add("order_by", "updated")
	query.Add("search", "name1")
	query.Add("search", "name2")
	query.Add("name", "")

	headers := make(http.Header)
	//date := time.Now().Format("Mon, 02 Jan 2006 15:04:05 GMT")
	date := "Sun, 07 Nov 2021 14:15:03 GMT"
	headers.Add("Date", date)
	headers.Add("Content-Type", "application/json;charset=UTF-8")
	headers.Add("Content-MD5", contentMd5)

	stringToSign := BuildSV1StringToSignature(method, urlPath, headers, query)
	signature := CalculateSV1Hmac256Signature(sk, stringToSign)

	require.Equal(t, signature, "gngQQkCADywByBbHGyleOI3MkogTDmRXGn4HKtxPt9o=")

	authorization := BuildSV1Authorization(ak, signature)
	require.Equal(t, authorization, "SV1-HMAC-SHA256 access-key-id:gngQQkCADywByBbHGyleOI3MkogTDmRXGn4HKtxPt9o=")
}

func Test_ParseSV1Authorization(t *testing.T) {
	authorization := "SV1-HMAC-SHA256 access-key-id:gngQQkCADywByBbHGyleOI3MkogTDmRXGn4HKtxPt9o="
	accessKeyId, signature, err := ParseSV1Authorization(authorization)
	require.Nil(t, err)
	require.Equal(t, accessKeyId, "access-key-id")
	require.Equal(t, signature, "gngQQkCADywByBbHGyleOI3MkogTDmRXGn4HKtxPt9o=")
}

package signer

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"net/http"
	"net/url"
	"sort"
	"strings"

	"github.com/DataWorkbench/common/qerror"
)

const (
	VersionSV1HmacSha256 = "SV1-HMAC-SHA256"
)

// BuildSV1StringToSignature build the string that will used by signature;
//
// example:
//    GET\n
//    Wed, 22 Jul 2019 17:20:31 GMT\n
//    application/json;charset=UTF-8\n
//    4gJE4saaMU4BqNR0kLY+lw==
//    /v1/workspace?limit=10&offset=20&search=name1&search=name2
func BuildSV1StringToSignature(method string, signPath string, headers http.Header, query url.Values) (stringToSign string) {
	// Build canonicalize query parameter
	var qKeys []string
	for key := range query {
		qKeys = append(qKeys, key)
	}
	sort.Strings(qKeys)
	var queryParts []string
	for _, key := range qKeys {
		values := query[key]
		key = url.QueryEscape(key)
		if len(values) == 0 {
			queryParts = append(queryParts, key+"=")
			continue
		}
		sort.Strings(values)
		for _, v := range values {
			queryParts = append(queryParts, key+"="+url.QueryEscape(v))
		}
	}
	// Build canonicalize resource
	var canonicalizeResource string
	if len(queryParts) != 0 {
		canonicalizeResource = signPath + "?" + strings.Join(queryParts, "&")
	} else {
		canonicalizeResource = signPath
	}

	date := headers.Get("X-Date")
	if date == "" {
		date = headers.Get("Date")
	}
	contentType := headers.Get("Content-Type")
	var mimeType string
	if i := strings.Index(contentType, ";"); i > 0 {
		mimeType = contentType[:i]
	} else {
		mimeType = contentType
	}
	contentMD5 := headers.Get("Content-MD5")

	// Build stringToSign
	signParts := []string{
		method,
		canonicalizeResource,
		date,
		mimeType,
		contentMD5,
	}
	return strings.Join(signParts, "\n")
}

// CalculateSV1Hmac256Signature calculate signature by sha256 algorithm.
func CalculateSV1Hmac256Signature(secretAccessKey string, stringToSign string) (signature string) {
	h := hmac.New(sha256.New, []byte(secretAccessKey))
	h.Write([]byte(stringToSign))
	b := h.Sum(nil)
	signature = base64.StdEncoding.EncodeToString(b)
	return
}

func BuildSV1Authorization(accessKeyId string, signature string) (authorization string) {
	return VersionSV1HmacSha256 + " " + accessKeyId + ":" + signature
}

func ParseSV1Authorization(authorization string) (accessKeyId string, signature string, err error) {
	if authorization == "" {
		err = qerror.MissingAuthorizationHeader
		return
	}
	authParts := strings.Split(authorization, " ")
	if len(authParts) != 2 {
		err = qerror.InvalidAuthorizationHeader
		return
	}
	if authParts[0] != VersionSV1HmacSha256 {
		// return invalid sign version.
		err = qerror.UnsupportedSignatureVersion
		return
	}

	signParts := strings.Split(authParts[1], ":")
	if len(signParts) != 2 {
		err = qerror.InvalidAuthorizationHeader
		return
	}

	accessKeyId = signParts[0]
	signature = signParts[1]
	return
}

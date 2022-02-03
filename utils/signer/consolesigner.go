package signer

import (
	"bytes"
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"io/ioutil"
	"net/http"
	"net/url"
	"sort"
	"strings"

	"github.com/DataWorkbench/gproto/xgo/types/pbrequest"
)

type ConsoleSigner struct {
	AccessKeyId     string
	SecretAccessKey string
	Zone            string
}

func (s *ConsoleSigner) Init(accessKeyID string, secretAccessKey string, zone string) {
	s.AccessKeyId = accessKeyID
	s.SecretAccessKey = secretAccessKey
	s.Zone = zone
}

func (s *ConsoleSigner) CalculateSignature(req *pbrequest.ValidateRequestSignature) string {
	m := md5.New()
	m.Write([]byte(req.ReqBody))
	stringToSign := strings.ToUpper(req.ReqMethod) + "\n" + req.ReqPath + "\n" + req.ReqQueryString + "\n" + hex.EncodeToString(m.Sum(nil))
	h := hmac.New(sha256.New, []byte(s.SecretAccessKey))
	h.Write([]byte(stringToSign))
	signature := strings.TrimSpace(base64.StdEncoding.EncodeToString(h.Sum(nil)))
	signature = strings.Replace(signature, " ", "+", -1)
	signature = url.QueryEscape(signature)
	return signature
}

func (s *ConsoleSigner) buildReqBody(request *http.Request) string {
	var data string
	if request.Method == "GET" || request.Method == "HEAD" || request.Method == "DELETE" {
		data = "null"
	} else {
		var bodyBytes []byte
		if request.Body != nil {
			bodyBytes, _ = ioutil.ReadAll(request.Body)
		}
		request.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))
		data = string(bodyBytes)
	}
	return data
}
func (s *ConsoleSigner) buildReqUrl(request *http.Request) string {
	if s.Zone != "" {
		return "/" + s.Zone + request.URL.Path + "/"
	}
	return request.URL.Path + "/"
}

func (s *ConsoleSigner) BuildValidateSignatureRequest(request *http.Request) *pbrequest.ValidateRequestSignature {
	reqBody := s.buildReqBody(request)
	params := request.URL.Query()
	keys := []string{}
	for key := range params {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	parts := []string{}
	var signature string
	var accessKeyID string
	for _, key := range keys {
		values := params[key]
		if len(values) > 0 {

			if key == "signature" {
				signature = values[0]
				continue
			}
			if key == "access_key_id" {
				accessKeyID = values[0]
			}
			if values[0] != "" {
				value := strings.TrimSpace(strings.Join(values, ""))
				value = url.QueryEscape(value)
				value = strings.Replace(value, "+", "%20", -1)
				parts = append(parts, key+"="+value)
			} else {
				parts = append(parts, key+"=")
			}
		} else {
			parts = append(parts, key+"=")
		}
	}

	urlParams := strings.Join(parts, "&")

	return &pbrequest.ValidateRequestSignature{
		ReqMethod:      request.Method,
		ReqPath:        s.buildReqUrl(request),
		ReqSignature:   signature,
		ReqQueryString: urlParams,
		ReqAccessKeyId: accessKeyID,
		ReqBody:        reqBody,
		ReqUserAgent:   ConsoleUserAgent,
	}
}

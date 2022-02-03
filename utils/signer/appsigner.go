package signer

import (
	"net/http"

	"github.com/DataWorkbench/gproto/xgo/types/pbrequest"
)

type AppSigner struct {
	AccessKeyId     string
	SecretAccessKey string
	Zone            string
}

func (s *AppSigner) Init(accessKeyID string, secretAccessKey string, zone string) {
	s.AccessKeyId = accessKeyID
	s.SecretAccessKey = secretAccessKey
	s.Zone = zone
}

func (s *AppSigner) CalculateSignature(req *pbrequest.ValidateRequestSignature) string {
	return ""
}

func (s *AppSigner) BuildValidateSignatureRequest(request *http.Request) *pbrequest.ValidateRequestSignature {
	return &pbrequest.ValidateRequestSignature{}
}

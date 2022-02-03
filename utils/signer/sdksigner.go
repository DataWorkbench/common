package signer

import (
	"net/http"

	"github.com/DataWorkbench/gproto/xgo/types/pbrequest"
)

type SdkSigner struct {
	AccessKeyId     string
	SecretAccessKey string
	Zone            string
}

func (s *SdkSigner) Init(accessKeyID string, secretAccessKey string, zone string) {
	s.AccessKeyId = accessKeyID
	s.SecretAccessKey = secretAccessKey
	s.Zone = zone
}

func (s *SdkSigner) CalculateSignature(req *pbrequest.ValidateRequestSignature) string {
	return ""
}

func (s *SdkSigner) BuildValidateSignatureRequest(request *http.Request) *pbrequest.ValidateRequestSignature {
	return &pbrequest.ValidateRequestSignature{}
}

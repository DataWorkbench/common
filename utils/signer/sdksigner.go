package signer

import (
	"net/http"

	"github.com/DataWorkbench/gproto/pkg/accountpb"
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

func (s *SdkSigner) CalculateSignature(req *accountpb.ValidateRequestSignatureRequest) string {
	return ""
}

func (s *SdkSigner) BuildValidateSignatureRequest(request *http.Request) *accountpb.ValidateRequestSignatureRequest {
	return &accountpb.ValidateRequestSignatureRequest{}
}

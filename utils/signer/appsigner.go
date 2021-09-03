package signer

import (
	"net/http"

	"github.com/DataWorkbench/gproto/pkg/accountpb"
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

func (s *AppSigner) CalculateSignature(req *accountpb.ValidateRequestSignatureRequest) string {
	return ""
}

func (s *AppSigner) BuildValidateSignatureRequest(request *http.Request) *accountpb.ValidateRequestSignatureRequest {
	return &accountpb.ValidateRequestSignatureRequest{}
}

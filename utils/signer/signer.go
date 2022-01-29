package signer

import (
	"net/http"

	"github.com/DataWorkbench/gproto/pkg/types/pbrequest"
)

const (
	ConsoleUserAgent = "QingCloud-Web-Console"
	AppUserAgent     = "App"
	SdkUserAgent     = "SDK"
)

type Signer interface {
	CalculateSignature(req *pbrequest.ValidateRequestSignature) string
	BuildValidateSignatureRequest(request *http.Request) *pbrequest.ValidateRequestSignature
	Init(accessKeyID string, secretAccessKey string, zone string)
}

func CreateSigner(userAgent string) Signer {
	switch userAgent {
	case ConsoleUserAgent:
		return &ConsoleSigner{}
	case AppUserAgent:
		return &AppSigner{}
	case SdkUserAgent:
		return &SdkSigner{}
	default:
		return &ConsoleSigner{}
	}

}

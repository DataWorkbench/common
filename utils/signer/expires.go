package signer

import (
	"net/http"
	"time"

	"github.com/DataWorkbench/common/qerror"
)

// CheckRequestExpired check the header "Date" or "X-Date" is valid.
// And check the signature whether expired;
// expires is the expire time seconds
func CheckRequestExpired(headers http.Header, expires int) (err error) {
	date := headers.Get("X-Date")
	if date == "" {
		date = headers.Get("Date")
	}
	if date == "" {
		err = qerror.MissingDateHeader
		return
	}
	var reqTime time.Time
	reqTime, err = time.ParseInLocation("Mon, 02 Jan 2006 15:04:05 GMT", date, time.UTC)
	if err != nil {
		err = qerror.InvalidDateHeader
		return
	}
	// To prevent replay attack, each signature is just valid within expires
	if time.Now().UTC().Sub(reqTime).Seconds() > float64(expires) {
		err = qerror.ExpiredSignature
		return
	}
	return
}

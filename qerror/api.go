package qerror

// Response used to return http error response body
type Response struct {
	Code      string `json:"code"`
	Desc      string `json:"desc"`
	Status    int    `json:"status"`
	RequestID string `json:"request_id"`
}

func NewResponse(err *Error, reqId string) *Response {
	return &Response{
		Code:      err.code,
		Desc:      err.desc,
		Status:    err.status,
		RequestID: reqId,
	}
}

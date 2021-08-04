package qerror

// Response used to return http error response body
type Response struct {
	// The error summary.
	Code string `json:"code"`
	// The error description format with en_us.
	EnUS string `json:"en_us"`
	// The error description format with zh_cn.
	ZhCN string `json:"zh_cn"`
	// The http status code.
	Status int `json:"status"`
	// The request id that same as Header "X-Request-Id".
	RequestID string `json:"request_id"`
}

func NewResponse(err *Error, reqId string) *Response {
	return &Response{
		Code:      err.code,
		EnUS:      err.enUS,
		ZhCN:      err.zhCN,
		Status:    err.status,
		RequestID: reqId,
	}
}

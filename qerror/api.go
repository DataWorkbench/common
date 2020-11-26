package qerror

// Response used to return http error response body
type Response struct {
	Code      string `json:"code"`
	EnUS      string `json:"en_us"`
	ZhCN      string `json:"zh_cn"`
	Status    int    `json:"status"`
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

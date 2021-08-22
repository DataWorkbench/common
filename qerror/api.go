package qerror

type Detail struct {
	// The error description format with en_us.
	EnUs string `json:"en_us"`
	// The error description format with zh_cn.
	ZhCn string `json:"zh_cn"`
}

// Response used to http error response body
type Response struct {
	// The is summary information of error.
	Code string `json:"code"`
	// The http status code.
	Status int `json:"status"`
	// The request id that same as Header "X-Request-Id".
	RequestID string `json:"request_id"`
	// detail is detail information of error.
	Detail Detail `json:"detail"`
}

func NewResponse(err *Error, reqId string) *Response {
	return &Response{
		Code:      err.code,
		Status:    err.status,
		RequestID: reqId,
		Detail: Detail{
			EnUs: err.enUS,
			ZhCn: err.zhCN,
		},
	}
}

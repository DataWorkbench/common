package qerror

import (
	"fmt"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/DataWorkbench/common/qerror/internal/qerror"
)

type pbError = qerror.Error

type Error struct {
	code   string
	status int
	enUS   string
	zhCN   string
}

func (e Error) Error() string {
	return e.code
}

func (e *Error) Format(a ...interface{}) *Error {
	e.enUS = fmt.Sprintf(e.enUS, a...)
	e.zhCN = fmt.Sprintf(e.zhCN, a...)
	return e
}

func (e *Error) String() string {
	return fmt.Sprintf("code: %s, desc: %s, status: %d", e.code, e.enUS, e.status)
}

func (e *Error) Code() string {
	return e.code
}

func (e *Error) Status() int {
	return e.status
}

// GRPCStatus returns the Status represented by grpc.Status.
func (e *Error) GRPCStatus() *status.Status {
	s := status.New(codes.Unknown, e.Error())
	sd, err := s.WithDetails(&pbError{
		Code:   e.code,
		Status: int32(e.status),
		EnUS:   e.enUS,
		ZhCN:   e.zhCN,
	})
	if err == nil {
		return sd
	}
	return s
}

// FromGRPC check the err's type if returned by gRPC
func FromGRPC(err error) *Error {
	s, ok := status.FromError(err)
	if !ok {
		return nil
	}

	var pe *pbError
	d := s.Details()
	for i := range d {
		pe, ok = d[i].(*pbError)
		if ok {
			break
		}
	}
	if pe == nil {
		return nil
	}
	return &Error{
		code:   pe.Code,
		status: int(pe.Status),
		enUS:   pe.EnUS,
		zhCN:   pe.ZhCN,
	}
}

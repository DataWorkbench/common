package qerror

import (
	"fmt"

	"github.com/DataWorkbench/gproto/pkg/model"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type pbError = model.Error

type Error struct {
	code   string // error code
	status int    // http status code.
	enUS   string
	zhCN   string
}

func (e Error) Error() string {
	return e.code
}

func (e *Error) Format(a ...interface{}) *Error {
	return &Error{
		code:   e.code,
		status: e.status,
		enUS:   fmt.Sprintf(e.enUS, a...),
		zhCN:   fmt.Sprintf(e.zhCN, a...),
	}
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
		EnUs:   e.enUS,
		ZhCn:   e.zhCN,
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
		enUS:   pe.EnUs,
		zhCN:   pe.ZhCn,
	}
}

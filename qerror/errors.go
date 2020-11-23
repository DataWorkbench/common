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
	desc   string
	status int
}

func (e Error) Error() string {
	return e.desc
}

func (e Error) String() string {
	return fmt.Sprintf("code: %s, desc: %s, status: %d", e.code, e.desc, e.status)
}

func (e Error) Code() string {
	return e.code
}

func (e Error) Desc() string {
	return e.desc
}

func (e Error) Status() int {
	return e.status
}

// GRPCStatus returns the Status represented by grpc.Status.
func (e *Error) GRPCStatus() *status.Status {
	s := status.New(codes.Unknown, e.Error())
	sd, err := s.WithDetails(&pbError{
		Code:   e.code,
		Desc:   e.desc,
		Status: int32(e.status),
	})
	if err == nil {
		return sd
	}
	return s
}

// WithDesc set desc message to specified error
func WithDesc(err *Error, desc string) *Error {
	err.desc = desc
	return err
}

// FromGRPC check the err's type if returned by gRPC
func FromGRPC(err error) *Error {
	s, ok := status.FromError(err)
	if !ok {
		return nil
	}

	switch s.Code() {
	case codes.InvalidArgument:
		return WithDesc(InvalidRequest, s.Message())
	default:
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
			desc:   pe.Desc,
			status: int(pe.Status),
		}
	}
}

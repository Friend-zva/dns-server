package apperror

import (
	"fmt"

	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

type AppError struct {
	Code    int
	Message string
	Err     error
}

var (
	ErrInvArgs = &AppError{
		Code:    int(codes.InvalidArgument),
		Message: "invalid arguments",
	}
	ErrInternal = &AppError{
		Code:    int(codes.Internal),
		Message: "internal server error",
	}
	ErrExternal = &AppError{
		Code:    int(codes.Unavailable),
		Message: "external service error",
	}
)

func (e *AppError) Wrap(err error) *AppError {
	return &AppError{
		Code:    e.Code,
		Message: e.Message,
		Err:     err,
	}
}

func (e *AppError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.Err)
	}
	return e.Message
}

func (e *AppError) Is(target error) bool {
	t, ok := target.(*AppError)
	if !ok {
		return false
	}
	return e.Code == t.Code
}

func (e *AppError) Unwrap() error {
	return e.Err
}

func (e *AppError) GRPCStatus() *status.Status {
	return status.New(codes.Code(e.Code), e.Error())
}

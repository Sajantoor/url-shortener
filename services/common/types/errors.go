package types

import (
	"fmt"

	codes "google.golang.org/grpc/codes"
)

type RequestError struct {
	Code codes.Code
	Err  error
}

var ReqError *RequestError

func (r *RequestError) Error() string {
	return fmt.Sprintf("status %d: err %v", r.Code, r.Err)
}

func New(statusCode codes.Code, err error) *RequestError {
	return &RequestError{
		Code: statusCode,
		Err:  err,
	}
}

func NotFoundError(err error) *RequestError {
	return &RequestError{
		Code: codes.NotFound,
		Err:  err,
	}
}

func InternalServerError(err error) *RequestError {
	return &RequestError{
		Code: codes.Internal,
		Err:  err,
	}
}

func AlreadyExistsError(err error) *RequestError {
	return &RequestError{
		Code: codes.AlreadyExists,
		Err:  err,
	}
}

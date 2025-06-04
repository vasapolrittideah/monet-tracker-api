package errorcode

import (
	"net/http"

	"google.golang.org/grpc/codes"
)

type Code string

const (
	// gRPC-like codes
	Unknown            Code = "ERR_UNKNOWN"
	Canceled           Code = "ERR_CANCELED"
	InvalidArgument    Code = "ERR_INVALID_ARGUMENT"
	DeadlineExceeded   Code = "ERR_DEADLINE_EXCEEDED"
	NotFound           Code = "ERR_NOT_FOUND"
	AlreadyExists      Code = "ERR_ALREADY_EXISTS"
	PermissionDenied   Code = "ERR_PERMISSION_DENIED"
	ResourceExhausted  Code = "ERR_RESOURCE_EXHAUSTED"
	FailedPrecondition Code = "ERR_FAILED_PRECONDITION"
	Aborted            Code = "ERR_ABORTED"
	OutOfRange         Code = "ERR_OUT_OF_RANGE"
	Unimplemented      Code = "ERR_UNIMPLEMENTED"
	Internal           Code = "ERR_INTERNAL"
	Unavailable        Code = "ERR_UNAVAILABLE"
	DataLoss           Code = "ERR_DATA_LOSS"
	Unauthenticated    Code = "ERR_UNAUTHENTICATED"

	// Custom app-specific codes
	NotRegistered Code = "ERR_NOT_REGISTERED"
)

var codeToHttpStatusMap = map[Code]int{
	Canceled:           http.StatusRequestTimeout,
	InvalidArgument:    http.StatusBadRequest,
	DeadlineExceeded:   http.StatusGatewayTimeout,
	NotFound:           http.StatusNotFound,
	AlreadyExists:      http.StatusConflict,
	PermissionDenied:   http.StatusForbidden,
	ResourceExhausted:  http.StatusTooManyRequests,
	FailedPrecondition: http.StatusPreconditionFailed,
	Aborted:            http.StatusConflict,
	OutOfRange:         http.StatusBadRequest,
	Unimplemented:      http.StatusNotImplemented,
	Internal:           http.StatusInternalServerError,
	Unavailable:        http.StatusServiceUnavailable,
	DataLoss:           http.StatusInternalServerError,
	Unauthenticated:    http.StatusUnauthorized,
}

var grpcCodeToCodeMap = map[codes.Code]Code{
	codes.Canceled:           Canceled,
	codes.InvalidArgument:    InvalidArgument,
	codes.DeadlineExceeded:   DeadlineExceeded,
	codes.NotFound:           NotFound,
	codes.AlreadyExists:      AlreadyExists,
	codes.PermissionDenied:   PermissionDenied,
	codes.ResourceExhausted:  ResourceExhausted,
	codes.FailedPrecondition: FailedPrecondition,
	codes.Aborted:            Aborted,
	codes.OutOfRange:         OutOfRange,
	codes.Unimplemented:      Unimplemented,
	codes.Internal:           Internal,
	codes.Unavailable:        Unavailable,
	codes.DataLoss:           DataLoss,
	codes.Unauthenticated:    Unauthenticated,
}

var codeToGrpcCodeMap = map[Code]codes.Code{
	Canceled:           codes.Canceled,
	InvalidArgument:    codes.InvalidArgument,
	DeadlineExceeded:   codes.DeadlineExceeded,
	NotFound:           codes.NotFound,
	AlreadyExists:      codes.AlreadyExists,
	PermissionDenied:   codes.PermissionDenied,
	ResourceExhausted:  codes.ResourceExhausted,
	FailedPrecondition: codes.FailedPrecondition,
	Aborted:            codes.Aborted,
	OutOfRange:         codes.OutOfRange,
	Unimplemented:      codes.Unimplemented,
	Internal:           codes.Internal,
	Unavailable:        codes.Unavailable,
	DataLoss:           codes.DataLoss,
	Unauthenticated:    codes.Unauthenticated,
}

func FromGrpcCode(code codes.Code) Code {
	if c, ok := grpcCodeToCodeMap[code]; ok {
		return c
	}
	return Unknown
}

func (c Code) ToGrpcCode() codes.Code {
	if grpcCode, ok := codeToGrpcCodeMap[c]; ok {
		return grpcCode
	}
	return codes.Unknown
}

func (c Code) ToHttpStatus() int {
	if status, ok := codeToHttpStatusMap[c]; ok {
		return status
	}
	return http.StatusInternalServerError
}

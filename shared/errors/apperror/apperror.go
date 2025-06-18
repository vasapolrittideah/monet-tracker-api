package apperror

import (
	"net/http"

	"google.golang.org/grpc/codes"
)

const (
	ErrUnknown            = "ERR_UNKNOWN"
	ErrCanceled           = "ERR_CANCELED"
	ErrInvalidArgument    = "ERR_INVALID_ARGUMENT"
	ErrDeadlineExceeded   = "ERR_DEADLINE_EXCEEDED"
	ErrNotFound           = "ERR_NOT_FOUND"
	ErrAlreadyExists      = "ERR_ALREADY_EXISTS"
	ErrPermissionDenied   = "ERR_PERMISSION_DENIED"
	ErrResourceExhausted  = "ERR_RESOURCE_EXHAUSTED"
	ErrFailedPrecondition = "ERR_FAILED_PRECONDITION"
	ErrAborted            = "ERR_ABORTED"
	ErrOutOfRange         = "ERR_OUT_OF_RANGE"
	ErrUnimplemented      = "ERR_UNIMPLEMENTED"
	ErrInternal           = "ERR_INTERNAL"
	ErrUnavailable        = "ERR_UNAVAILABLE"
	ErrDataLoss           = "ERR_DATA_LOSS"
	ErrUnauthenticated    = "ERR_UNAUTHENTICATED"
)

var AppCodeToGRPCCodeMap = map[string]codes.Code{
	ErrUnknown:            codes.Unknown,
	ErrCanceled:           codes.Canceled,
	ErrInvalidArgument:    codes.InvalidArgument,
	ErrDeadlineExceeded:   codes.DeadlineExceeded,
	ErrNotFound:           codes.NotFound,
	ErrAlreadyExists:      codes.AlreadyExists,
	ErrPermissionDenied:   codes.PermissionDenied,
	ErrResourceExhausted:  codes.ResourceExhausted,
	ErrFailedPrecondition: codes.FailedPrecondition,
	ErrAborted:            codes.Aborted,
	ErrOutOfRange:         codes.OutOfRange,
	ErrUnimplemented:      codes.Unimplemented,
	ErrInternal:           codes.Internal,
	ErrUnavailable:        codes.Unavailable,
	ErrDataLoss:           codes.DataLoss,
	ErrUnauthenticated:    codes.Unauthenticated,
}

var GRPCCodeToAppCodeMap = map[codes.Code]string{
	codes.Unknown:            ErrUnknown,
	codes.Canceled:           ErrCanceled,
	codes.InvalidArgument:    ErrInvalidArgument,
	codes.DeadlineExceeded:   ErrDeadlineExceeded,
	codes.NotFound:           ErrNotFound,
	codes.AlreadyExists:      ErrAlreadyExists,
	codes.PermissionDenied:   ErrPermissionDenied,
	codes.ResourceExhausted:  ErrResourceExhausted,
	codes.FailedPrecondition: ErrFailedPrecondition,
	codes.Aborted:            ErrAborted,
	codes.OutOfRange:         ErrOutOfRange,
	codes.Unimplemented:      ErrUnimplemented,
	codes.Internal:           ErrInternal,
	codes.Unavailable:        ErrUnavailable,
	codes.DataLoss:           ErrDataLoss,
	codes.Unauthenticated:    ErrUnauthenticated,
}

const StatusClientClosedRequest = 499

var AppCodeToHTTPStatusMap = map[string]int{
	ErrUnknown:            http.StatusInternalServerError,
	ErrCanceled:           StatusClientClosedRequest,
	ErrInvalidArgument:    http.StatusBadRequest,
	ErrDeadlineExceeded:   http.StatusGatewayTimeout,
	ErrNotFound:           http.StatusNotFound,
	ErrAlreadyExists:      http.StatusConflict,
	ErrPermissionDenied:   http.StatusForbidden,
	ErrResourceExhausted:  http.StatusTooManyRequests,
	ErrFailedPrecondition: http.StatusBadRequest,
	ErrAborted:            http.StatusConflict,
	ErrOutOfRange:         http.StatusBadRequest,
	ErrUnimplemented:      http.StatusNotImplemented,
	ErrInternal:           http.StatusInternalServerError,
	ErrUnavailable:        http.StatusServiceUnavailable,
	ErrDataLoss:           http.StatusInternalServerError,
	ErrUnauthenticated:    http.StatusUnauthorized,
}

type AppError struct {
	Code    string
	Message string
}

func (e *AppError) Error() string {
	return e.Message
}

func NewError(code string, message ...string) *AppError {
	err := &AppError{
		Code:    code,
		Message: "somethin went wrong",
	}
	if len(message) > 0 {
		err.Message = message[0]
	}

	return err
}

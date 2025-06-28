package apperror

import (
	"net/http"

	"google.golang.org/grpc/codes"
)

const (
	CodeUnknown            = "UNKNOWN"
	CodeCanceled           = "CANCELED"
	CodeInvalidArgument    = "INVALID_ARGUMENT"
	CodeDeadlineExceeded   = "DEADLINE_EXCEEDED"
	CodeNotFound           = "NOT_FOUND"
	CodeAlreadyExists      = "ALREADY_EXISTS"
	CodePermissionDenied   = "PERMISSION_DENIED"
	CodeResourceExhausted  = "RESOURCE_EXHAUSTED"
	CodeFailedPrecondition = "FAILED_PRECONDITION"
	CodeAborted            = "ABORTED"
	CodeOutOfRange         = "OUT_OF_RANGE"
	CodeUnimplemented      = "UNIMPLEMENTED"
	CodeInternal           = "INTERNAL"
	CodeUnavailable        = "UNAVAILABLE"
	CodeDataLoss           = "DATA_LOSS"
	CodeUnauthenticated    = "UNAUTHENTICATED"
)

var AppCodeToGRPCCodeMap = map[string]codes.Code{
	CodeUnknown:            codes.Unknown,
	CodeCanceled:           codes.Canceled,
	CodeInvalidArgument:    codes.InvalidArgument,
	CodeDeadlineExceeded:   codes.DeadlineExceeded,
	CodeNotFound:           codes.NotFound,
	CodeAlreadyExists:      codes.AlreadyExists,
	CodePermissionDenied:   codes.PermissionDenied,
	CodeResourceExhausted:  codes.ResourceExhausted,
	CodeFailedPrecondition: codes.FailedPrecondition,
	CodeAborted:            codes.Aborted,
	CodeOutOfRange:         codes.OutOfRange,
	CodeUnimplemented:      codes.Unimplemented,
	CodeInternal:           codes.Internal,
	CodeUnavailable:        codes.Unavailable,
	CodeDataLoss:           codes.DataLoss,
	CodeUnauthenticated:    codes.Unauthenticated,
}

var GRPCCodeToAppCodeMap = map[codes.Code]string{
	codes.Unknown:            CodeUnknown,
	codes.Canceled:           CodeCanceled,
	codes.InvalidArgument:    CodeInvalidArgument,
	codes.DeadlineExceeded:   CodeDeadlineExceeded,
	codes.NotFound:           CodeNotFound,
	codes.AlreadyExists:      CodeAlreadyExists,
	codes.PermissionDenied:   CodePermissionDenied,
	codes.ResourceExhausted:  CodeResourceExhausted,
	codes.FailedPrecondition: CodeFailedPrecondition,
	codes.Aborted:            CodeAborted,
	codes.OutOfRange:         CodeOutOfRange,
	codes.Unimplemented:      CodeUnimplemented,
	codes.Internal:           CodeInternal,
	codes.Unavailable:        CodeUnavailable,
	codes.DataLoss:           CodeDataLoss,
	codes.Unauthenticated:    CodeUnauthenticated,
}

const StatusClientClosedRequest = 499

var AppCodeToHTTPStatusMap = map[string]int{
	CodeUnknown:            http.StatusInternalServerError,
	CodeCanceled:           StatusClientClosedRequest,
	CodeInvalidArgument:    http.StatusBadRequest,
	CodeDeadlineExceeded:   http.StatusGatewayTimeout,
	CodeNotFound:           http.StatusNotFound,
	CodeAlreadyExists:      http.StatusConflict,
	CodePermissionDenied:   http.StatusForbidden,
	CodeResourceExhausted:  http.StatusTooManyRequests,
	CodeFailedPrecondition: http.StatusBadRequest,
	CodeAborted:            http.StatusConflict,
	CodeOutOfRange:         http.StatusBadRequest,
	CodeUnimplemented:      http.StatusNotImplemented,
	CodeInternal:           http.StatusInternalServerError,
	CodeUnavailable:        http.StatusServiceUnavailable,
	CodeDataLoss:           http.StatusInternalServerError,
	CodeUnauthenticated:    http.StatusUnauthorized,
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
		Message: "something went wrong",
	}
	if len(message) > 0 {
		err.Message = message[0]
	}

	return err
}

package httperror

import (
	"net/http"

	"google.golang.org/grpc/codes"
)

type HTTPError struct {
	Code    codes.Code `json:"code"    swaggertype:"integer"`
	Message string     `json:"message"`
}

type HTTPValidationError struct {
	Code    codes.Code        `json:"code"    swaggertype:"integer"`
	Message string            `json:"message"`
	Details []ValidationError `json:"details"`
}

type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

const StatusClientClosedRequest = 499

var grpcToHTTPStatusMap = map[codes.Code]int{
	codes.OK:                 http.StatusOK,
	codes.Canceled:           StatusClientClosedRequest,
	codes.Unknown:            http.StatusInternalServerError,
	codes.InvalidArgument:    http.StatusBadRequest,
	codes.DeadlineExceeded:   http.StatusGatewayTimeout,
	codes.NotFound:           http.StatusNotFound,
	codes.AlreadyExists:      http.StatusConflict,
	codes.PermissionDenied:   http.StatusForbidden,
	codes.ResourceExhausted:  http.StatusTooManyRequests,
	codes.FailedPrecondition: http.StatusBadRequest,
	codes.Aborted:            http.StatusConflict,
	codes.OutOfRange:         http.StatusBadRequest,
	codes.Unimplemented:      http.StatusNotImplemented,
	codes.Internal:           http.StatusInternalServerError,
	codes.Unavailable:        http.StatusServiceUnavailable,
	codes.DataLoss:           http.StatusInternalServerError,
	codes.Unauthenticated:    http.StatusUnauthorized,
}

func HTTPStatusFromCode(code codes.Code) int {
	if status, ok := grpcToHTTPStatusMap[code]; ok {
		return status
	}

	return http.StatusInternalServerError
}

func NewHTTPError(code codes.Code, message string) HTTPError {
	return HTTPError{
		Code:    code,
		Message: message,
	}
}

func NewValidationError(details []ValidationError) HTTPValidationError {
	return HTTPValidationError{
		Code:    http.StatusBadRequest,
		Message: "validation failed",
		Details: details,
	}
}

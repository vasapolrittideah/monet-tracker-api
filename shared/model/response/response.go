package response

import (
	"github.com/vasapolrittideah/money-tracker-api/shared/constants/errorcode"
)

type AppStatus string

const (
	StatusSuccess AppStatus = "SUCCESS"
	StatusFailure AppStatus = "FAILURE"
	StatusError   AppStatus = "ERROR"
)

type response struct {
	Status       AppStatus      `json:"status"`
	ErrorCode    errorcode.Code `json:"error_code,omitempty"`
	ErrorMessage string         `json:"error_message,omitempty"`
	ErrorDetails any            `json:"error_details,omitempty"`
	Data         any            `json:"data"`
}

type InvalidField struct {
	Field  string `json:"field"`
	Reason string `json:"reason"`
}

func Success(data any) response {
	return response{
		Status: StatusSuccess,
		Data:   data,
	}
}

func Error(errorCode errorcode.Code, message string) response {
	return response{
		Status:       StatusError,
		ErrorCode:    errorCode,
		ErrorMessage: message,
	}
}

func ValidationFailed(invalidFields []InvalidField) response {
	return response{
		Status:       StatusFailure,
		ErrorCode:    errorcode.InvalidArgument,
		ErrorMessage: "Some fields are invalid",
		ErrorDetails: invalidFields,
	}
}

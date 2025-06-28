package httperror

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/vasapolrittideah/money-tracker-api/shared/errors/apperror"
)

type HTTPError struct {
	Code    string `json:"code"    swaggertype:"integer"`
	Message string `json:"message"`
}

type HTTPValidationError struct {
	Code    string            `json:"code"    swaggertype:"integer"`
	Message string            `json:"message"`
	Details []ValidationError `json:"details"`
}

type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

func toHTTPStatus(code string) int {
	if status, ok := apperror.AppCodeToHTTPStatusMap[code]; ok {
		return status
	}
	return http.StatusInternalServerError
}

func FromAppError(c *fiber.Ctx, err *apperror.AppError) error {
	return c.Status(toHTTPStatus(err.Code)).JSON(HTTPError{
		Code:    err.Code,
		Message: err.Message,
	})
}

func NewBadRequestError(c *fiber.Ctx, message string) error {
	return c.Status(http.StatusBadRequest).JSON(HTTPError{
		Code:    apperror.CodeInvalidArgument,
		Message: message,
	})
}

func NewUnauthorizedError(c *fiber.Ctx, message string) error {
	return c.Status(http.StatusUnauthorized).JSON(HTTPError{
		Code:    apperror.CodeUnauthenticated,
		Message: message,
	})
}

func NewValidationError(c *fiber.Ctx, details []ValidationError) error {
	return c.Status(http.StatusBadRequest).JSON(HTTPValidationError{
		Code:    apperror.CodeInvalidArgument,
		Message: "validation error",
		Details: details,
	})
}

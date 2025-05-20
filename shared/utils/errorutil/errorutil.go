package errorutil

import (
	"errors"
	"fmt"
	"strings"

	"github.com/vasapolrittideah/money-tracker-api/shared/domain/apperror"
	"google.golang.org/grpc/codes"
	"gorm.io/gorm"
)

func HandleUnknownDatabaseError(err error) *apperror.Error {
	return apperror.New(codes.Unknown, fmt.Errorf("unknown database error: %s", err))
}

func HandleRecordNotFoundError(err error) *apperror.Error {
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return apperror.New(codes.NotFound, fmt.Errorf("record not found: %s", err))
	}

	return HandleUnknownDatabaseError(err)
}

func HandleUnqiueConstraintError(err error) *apperror.Error {
	if strings.Contains(err.Error(), "duplicate key") {
		return apperror.New(codes.AlreadyExists, fmt.Errorf("duplicate key: %s", err))
	}

	return HandleUnknownDatabaseError(err)
}

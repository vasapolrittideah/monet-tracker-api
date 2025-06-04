package errorutil

import (
	"errors"
	"fmt"
	"strings"

	"github.com/vasapolrittideah/money-tracker-api/shared/constants/errorcode"
	"github.com/vasapolrittideah/money-tracker-api/shared/model/apperror"
	"gorm.io/gorm"
)

func HandleUnknownDatabaseError(err error) *apperror.Error {
	return apperror.New(errorcode.Unknown, fmt.Errorf("unknown database error: %s", err))
}

func HandleRecordNotFoundError(err error) *apperror.Error {
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return apperror.New(errorcode.NotFound, fmt.Errorf("record not found: %s", err))
	}

	return HandleUnknownDatabaseError(err)
}

func HandleUnqiueConstraintError(err error) *apperror.Error {
	if strings.Contains(err.Error(), "duplicate key") {
		return apperror.New(errorcode.AlreadyExists, fmt.Errorf("duplicate key: %s", err))
	}

	return HandleUnknownDatabaseError(err)
}

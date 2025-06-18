package grpcerror

import (
	"github.com/vasapolrittideah/money-tracker-api/shared/errors/apperror"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func toGRPCCode(code string) codes.Code {
	if grpcCode, ok := apperror.AppCodeToGRPCCodeMap[code]; ok {
		return grpcCode
	}
	return codes.Unknown
}

func FromAppError(err apperror.AppError) error {
	return status.Error(toGRPCCode(err.Code), err.Message)
}

func ToAppError(err error) error {
	st := status.Convert(err)
	return apperror.NewError(apperror.GRPCCodeToAppCodeMap[st.Code()], st.Message())
}

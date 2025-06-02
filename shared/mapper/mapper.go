package mapper

import (
	"github.com/google/uuid"
	userpb "github.com/vasapolrittideah/money-tracker-api/generated/protobuf/user"
	"github.com/vasapolrittideah/money-tracker-api/shared/model/domain"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func MapUserEntityToProto(user *domain.User) *userpb.User {
	return &userpb.User{
		Id:                 uuid.UUID(user.Id).String(),
		FullName:           user.FullName,
		Email:              user.Email,
		HashedPassword:     user.HashedPassword,
		HashedRefreshToken: user.HashedRefreshToken,
		Verified:           user.Verified,
		CreatedAt:          timestamppb.New(user.CreatedAt),
		UpdatedAt:          timestamppb.New(user.UpdatedAt),
		LastSignInAt:       timestamppb.New(user.LastSignInAt),
	}
}

func MapUserProtoToEntity(user *userpb.User) *domain.User {
	return &domain.User{
		Id:                 uuid.MustParse(user.Id),
		FullName:           user.FullName,
		Email:              user.Email,
		HashedPassword:     user.HashedPassword,
		HashedRefreshToken: user.HashedRefreshToken,
		Verified:           user.Verified,
		CreatedAt:          user.CreatedAt.AsTime(),
		UpdatedAt:          user.UpdatedAt.AsTime(),
		LastSignInAt:       user.LastSignInAt.AsTime(),
	}
}

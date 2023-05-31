package gapi

import (
	db "github.com/quynhtruong/backend-master-class/db/sqlc"
	"github.com/quynhtruong/backend-master-class/pb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func ConvertUser(user db.User) *pb.User {
	return &pb.User{
		Username:          user.Username,
		FullName:          user.FullName,
		Email:             user.Email,
		PasswordChangedAt: timestamppb.New(user.PasswordChanged),
		CreatedAt:         timestamppb.New(user.CreatedAt),
	}
}

package gapi

import (
	"context"
	"log"

	"github.com/lib/pq"
	db "github.com/quynhtruong/backend-master-class/db/sqlc"
	"github.com/quynhtruong/backend-master-class/pb"
	"github.com/quynhtruong/backend-master-class/util"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (server *Server) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {

	hashedPassword, err := util.HashPassword(req.GetPassword())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "method CreateUser is not implemented: %s", err)
	}
	arg := db.CreateUserParams{
		Username:       req.GetUsername(),
		HashedPassword: hashedPassword,
		FullName:       req.GetFullName(),
		Email:          req.GetEmail(),
	}
	user, err := server.store.CreateUser(ctx, arg)
	if err != nil {
		if pgErr, ok := err.(*pq.Error); ok {
			switch pgErr.Code.Name() {
			case "unique_violation":
				return nil, status.Errorf(codes.AlreadyExists, "usename already exist: %s", err)
			}
			log.Println(pgErr.Code.Name())
		}
		return nil, status.Errorf(codes.Internal, "failed to create user: %s", err)
	}
	response := &pb.CreateUserResponse{
		User: ConvertUser(user),
	}

	return response, nil

}

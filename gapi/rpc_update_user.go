package gapi

import (
	"context"
	"database/sql"
	"time"

	db "github.com/quynhtruong/backend-master-class/db/sqlc"
	"github.com/quynhtruong/backend-master-class/pb"
	"github.com/quynhtruong/backend-master-class/util"
	"github.com/quynhtruong/backend-master-class/val"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (server *Server) UpdateUser(ctx context.Context, req *pb.UpdateUserRequest) (*pb.UpdateUserResponse, error) {
	violations := validateUpdateUserRequest(req)
	if violations != nil {
		return nil, invalidArgumentError(violations)
	}

	arg := db.UpdateUserParams{
		Username: req.GetUsername(),
		FullName: sql.NullString{
			String: req.GetFullName(),
			Valid:  req.FullName != nil,
		},
		Email: sql.NullString{
			String: req.GetEmail(),
			Valid:  req.Email != nil,
		},
	}
	if req.Password != nil {
		hashedPassword, err := util.HashPassword(req.GetPassword())
		if err != nil {
			return nil, status.Errorf(codes.Internal, "method CreateUser is not implemented: %s", err)
		}
		arg.HashedPassword = sql.NullString{
			String: hashedPassword,
			Valid:  true,
		}
		arg.PasswordChanged = sql.NullTime{
			Time:  time.Now(),
			Valid: true,
		}
	}
	user, err := server.store.UpdateUser(ctx, arg)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, status.Errorf(codes.NotFound, "user does not exist")
		}
		return nil, status.Errorf(codes.Internal, "failed to create user: %s", err)
	}
	response := &pb.UpdateUserResponse{
		User: ConvertUser(user),
	}

	return response, nil
}

func validateUpdateUserRequest(req *pb.UpdateUserRequest) (violations [](*errdetails.BadRequest_FieldViolation)) {
	if err := val.ValidateUsername(req.GetUsername()); err != nil {
		violations = append(violations, fieldViolation("username", err))
	}
	if req.Password != nil {
		if err := val.ValidPassword(req.GetPassword()); err != nil {
			violations = append(violations, fieldViolation("password", err))
		}
	}
	if req.Email != nil {
		if err := val.ValidateEmail(req.GetEmail()); err != nil {
			violations = append(violations, fieldViolation("email", err))
		}
	}
	if req.FullName != nil {
		if err := val.ValidateUserFullName(req.GetFullName()); err != nil {
			violations = append(violations, fieldViolation("full_name", err))
		}
	}
	return violations
}

package gapi

import (
	"fmt"

	db "github.com/quynhtruong/backend-master-class/db/sqlc"
	"github.com/quynhtruong/backend-master-class/pb"
	"github.com/quynhtruong/backend-master-class/token"
	"github.com/quynhtruong/backend-master-class/util"
)

type Server struct {
	pb.UnimplementedSimpleBankServer
	config     util.Config
	store      db.Store
	tokenMaker token.Maker
}

func NewServer(config util.Config, store db.Store) (*Server, error) {
	tokenMaker, err := token.NewPasetoMaker(config.TokenSymetricKey)
	if err != nil {
		return nil, fmt.Errorf("can not create token maker: %w", err)
	}
	server := &Server{store: store, tokenMaker: tokenMaker, config: config}

	return server, nil
}

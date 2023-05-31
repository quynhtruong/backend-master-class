package main

import (
	"database/sql"
	"log"
	"net"

	_ "github.com/lib/pq"
	"github.com/quynhtruong/backend-master-class/api"
	db "github.com/quynhtruong/backend-master-class/db/sqlc"
	"github.com/quynhtruong/backend-master-class/gapi"
	"github.com/quynhtruong/backend-master-class/pb"
	"github.com/quynhtruong/backend-master-class/util"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal("Can not load config ", err)
	}
	conn, err := sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal("Cannot connect to DB: ", err)
	}

	store := db.NewStore(conn)
	runGrpcServer(config, store)
}

func runGinServer(config util.Config, store db.Store) {
	server, err := api.NewServer(config, store)
	if err != nil {
		log.Fatal("can not creat server", err)
	}
	err = server.Start(config.HttpServerAddress)
	if err != nil {
		log.Fatal("Cannot start server: ", err)
	}
}

func runGrpcServer(config util.Config, store db.Store) {
	server, err := gapi.NewServer(config, store)
	if err != nil {
		log.Fatal("can not creat server", err)
	}
	grpcServer := grpc.NewServer()
	pb.RegisterSimpleBankServer(grpcServer, server)
	reflection.Register(grpcServer)

	listener, err := net.Listen("tcp", config.GrpcServerAddress)
	if err != nil {
		log.Fatal("can not create listener")
	}
	log.Printf("start gRPC server at %s ", listener.Addr().String())
	err = grpcServer.Serve(listener)
	if err != nil {
		log.Fatal("can not start gRPC server")
	}
}

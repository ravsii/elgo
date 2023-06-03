package main

import (
	"log"
	"net"

	elgo_grpc "github.com/ravsii/elgo/grpc"
	elgo_schema "github.com/ravsii/elgo/grpc/schema"
	"google.golang.org/grpc"
)

func main() {
	srv, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatalf("net listen: %s", err)
	}

	grpcSrv := grpc.NewServer()
	elgoGrpc := elgo_grpc.NewGrpcServer()
	elgo_schema.RegisterPoolServer(grpcSrv, elgoGrpc)
	if err := grpcSrv.Serve(srv); err != nil {
		log.Fatalln(err)
	}
}

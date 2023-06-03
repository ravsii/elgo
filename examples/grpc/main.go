package main

import (
	"log"
	"net"

	elgo_grpc "github.com/ravsii/elgo/grpc"
	pb "github.com/ravsii/elgo/grpc/pb"
	"google.golang.org/grpc"
)

func main() {
	srv, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatalf("net listen: %s", err)
	}

	grpcSrv := grpc.NewServer()
	pb.RegisterPoolServer(grpcSrv, elgo_grpc.NewGrpcServer())
	if err := grpcSrv.Serve(srv); err != nil {
		log.Fatalln(err)
	}
}

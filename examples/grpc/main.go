package main

import (
	"context"
	"log"
	"net"

	"github.com/ravsii/elgo"
	elgo_grpc "github.com/ravsii/elgo/grpc"
	pb "github.com/ravsii/elgo/grpc/pb"
	"google.golang.org/grpc"
)

func main() {
	srv, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatalf("net listen: %s", err)
	}
	defer srv.Close()

	grpcSrv := grpc.NewServer()
	pb.RegisterPoolServer(grpcSrv, elgo_grpc.NewServer())
	go func() {
		if err := grpcSrv.Serve(srv); err != nil {
			log.Fatalln(err)
		}
	}()

	client, err := elgo_grpc.NewClient(":8080")
	if err != nil {
		log.Fatalln(err)
	}
	defer client.Close()

	player := elgo.BaseRatingPlayer{
		ID:  "1",
		ELO: 1,
	}

	err = client.Add(context.Background(), &player, &player)
	if err != nil {
		log.Fatalln(err)
	}
}

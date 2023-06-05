package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"
	"time"

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
	elgoSrv := elgo_grpc.NewPoolServer()
	defer elgoSrv.Close()
	pb.RegisterPoolServer(grpcSrv, elgoSrv)
	go func() {
		if err := grpcSrv.Serve(srv); err != nil {
			log.Fatalln("grpc listener:", err)
		}
	}()
	defer grpcSrv.GracefulStop()

	client, err := elgo_grpc.NewClient(":8080")
	if err != nil {
		log.Fatalln(err)
	}
	defer client.Close()

	client2, err := elgo_grpc.NewClient(":8080")
	if err != nil {
		log.Fatalln(err)
	}
	client2.Close()

	p1 := elgo.BaseRatingPlayer{
		ID:  "1",
		ELO: 1,
	}

	err = client.Add(context.Background(), &p1)
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Println("Let's wait 5 secs for nothing")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	match, err := client.RecieveMatch(ctx)
	if match != nil || !errors.Is(err, elgo.ErrNoMatchFound) {
		log.Fatalln("match or error is not empty, but expected to be", match, err)
	}

	fmt.Println("Let's wait for a match")

	p2 := elgo.BaseRatingPlayer{
		ID:  "2",
		ELO: 2,
	}

	err = client.Add(context.Background(), &p2)
	if err != nil {
		log.Fatalln(err)
	}

	match, err = client.RecieveMatch(context.Background())
	fmt.Println("match found!", match.Player1, match.Player2)

	client.Close()
	fmt.Println("closing conenction")
	time.Sleep(time.Second)
	fmt.Println("closed conenction")
}

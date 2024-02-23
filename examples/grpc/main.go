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
	log.Println("Creating tcp listener at 8080")
	srv, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatalf("net listen: %s", err)
	}
	defer srv.Close()

	grpcSrv := grpc.NewServer()
	pool := elgo.NewPool(
		elgo.WithIncreasePlayerBorderBy(0.1),
		elgo.WithPlayerRetryInterval(time.Second),
	)
	go pool.Run()
	elgoSrv := elgo_grpc.NewPoolServer(pool)
	defer elgoSrv.Close()
	pb.RegisterPoolServer(grpcSrv, elgoSrv)
	go func() {
		if err := grpcSrv.Serve(srv); err != nil {
			log.Fatalln("grpc listener:", err)
		}
	}()
	defer grpcSrv.GracefulStop()

	client1, err := elgo_grpc.NewClient(":8080")
	if err != nil {
		log.Println("can't create elgo_grpc client1:", err)
		return
	}
	defer client1.Close()

	client2, err := elgo_grpc.NewClient(":8080")
	if err != nil {
		log.Println("can't create elgo_grpc client2:", err)
		return
	}
	defer client2.Close()

	p1 := elgo.BaseRatingPlayer{
		ID:  "1",
		ELO: 1,
	}

	log.Println("Adding p1...")

	if err := client1.Add(context.Background(), &p1); err != nil {
		log.Println("can't add player1 to the pool using client1:", err)
		return
	}

	log.Println("Let's wait 5 secs for nothing")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	match, err := client1.RecieveMatch(ctx)
	if match != nil || !errors.Is(err, elgo.ErrNoMatchFound) {
		log.Println("match or error is not empty, but expected to be", match, err)
		return
	}

	log.Println("Let's wait for a match")

	p2 := elgo.BaseRatingPlayer{
		ID:  "2",
		ELO: 2,
	}

	err = client1.Add(context.Background(), &p2)
	if err != nil {
		log.Fatalln(err)
	}

	match, err = client1.RecieveMatch(context.Background())
	fmt.Println("match found!", match.Player1, match.Player2)

	client1.Close()
	fmt.Println("closing conenction")
	time.Sleep(time.Second)
	fmt.Println("closed conenction")
}

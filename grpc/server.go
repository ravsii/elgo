package grpc

import (
	"context"
	"errors"
	"fmt"
	"net"

	"github.com/ravsii/elgo"
	"github.com/ravsii/elgo/grpc/pb"
	"google.golang.org/grpc"
)

// ensuring we've implemented our server correctly.
var _ pb.PoolServer = (*PoolServer)(nil)

type PoolServer struct {
	pool *elgo.Pool
	pb.UnimplementedPoolServer
}

type ListenServer struct {
	conn    net.Listener
	poolSrv *PoolServer
	grpcSrv *grpc.Server
}

// NewListener will create a new grpc listener.
// gRPC server implementation should be created beforehand.
func NewListener(network, address string, pool *PoolServer, opts ...grpc.ServerOption) (*ListenServer, error) {
	var (
		srv ListenServer
		err error
	)

	srv.conn, err = net.Listen(network, address)
	if err != nil {
		return nil, fmt.Errorf("net listen: %s", err)
	}

	srv.grpcSrv = grpc.NewServer(opts...)
	srv.poolSrv = pool
	pb.RegisterPoolServer(srv.grpcSrv, srv.poolSrv)

	return &srv, nil
}

// Listen will start listening for incoming commands.
//
// This is a blocking operation.
func (s *ListenServer) Listen() error {
	if err := s.grpcSrv.Serve(s.conn); err != nil {
		return fmt.Errorf("grpc listen: %w", err)
	}

	return nil
}

// Close returns a map of players left in the queue at the time of closing.
func (s *ListenServer) Close() (elgo.Players, error) {
	if err := s.conn.Close(); err != nil {
		return nil, fmt.Errorf("net listener close: %w", err)
	}

	players := s.poolSrv.Close()
	s.grpcSrv.GracefulStop()

	return players, nil
}

// NewPoolServer returns a new pool grpc server, which simply implements
// a grpc interface. If you need a server with listener, use [NewListener]
func NewPoolServer(poolOpts ...elgo.PoolOpt) *PoolServer {
	pool := elgo.NewPool(poolOpts...)

	go pool.Run()

	return &PoolServer{
		pool: pool,
	}
}

// Add implements pb.PoolServer.
func (s *PoolServer) Add(ctx context.Context, player *pb.Player) (*pb.Empty, error) {
	select {
	case <-ctx.Done():
		return &pb.Empty{}, nil
	default:
		if err := s.pool.AddPlayer(player); err != nil {
			if errors.Is(err, &elgo.PlayerAlreadyExistsError{}) {
				return &pb.Empty{}, NewAlreadyExistsErr(player)
			}
			return &pb.Empty{}, NewCantAddErr(player)
		}
		return &pb.Empty{}, nil
	}
}

// Match implements pb.PoolServer.
func (s *PoolServer) Match(_ *pb.Empty, matches pb.Pool_MatchServer) error {
	for {
		select {
		case <-matches.Context().Done():
			return nil
		case m := <-s.pool.Matches():
			grpcMatch := &pb.PlayerMatch{
				P1: &pb.Player{Id: m.Player1.Identify()},
				P2: &pb.Player{Id: m.Player2.Identify()},
			}
			err := matches.Send(grpcMatch)
			if err != nil {
				return errors.Join(ErrCreateMatch, err)
			}

			return nil
		}
	}
}

// Remove implements pb.PoolServer.
func (s *PoolServer) Remove(ctx context.Context, player *pb.Player) (*pb.Empty, error) {
	select {
	case <-ctx.Done():
	default:
		s.pool.Remove(player)
	}
	return &pb.Empty{}, nil
}

// Size implements pb.PoolServer.
func (s *PoolServer) Size(ctx context.Context, _ *pb.Empty) (*pb.SizeResponse, error) {
	select {
	case <-ctx.Done():
		return &pb.SizeResponse{Size: 0}, nil
	default:
		return &pb.SizeResponse{
			Size: int32(s.pool.Size()),
		}, nil
	}
}

func (s *PoolServer) Close() map[string]elgo.Player {
	return s.pool.Close()
}

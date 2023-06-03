package grpc

import (
	context "context"

	"github.com/ravsii/elgo"
	"github.com/ravsii/elgo/grpc/pb"
)

// ensuring we've implemented our server correctly
var _ pb.PoolServer = (*grpcServer)(nil)

type grpcServer struct {
	pool *elgo.Pool
	pb.UnimplementedPoolServer
}

func NewServer(poolOpts ...elgo.PoolOpt) *grpcServer {
	return &grpcServer{
		pool: elgo.NewPool(poolOpts...),
	}
}

// Add implements pb.PoolServer
func (s *grpcServer) Add(ctx context.Context, player *pb.Player) (*pb.Empty, error) {
	select {
	case <-ctx.Done():
		return &pb.Empty{}, nil
	default:
		if err := s.pool.AddPlayer(player); err != nil {
			return &pb.Empty{}, err
		}
		return &pb.Empty{}, nil
	}
}

// Match implements pb.PoolServer
func (s *grpcServer) Match(_ *pb.Empty, matches pb.Pool_MatchServer) error {
	for {
		select {
		case <-matches.Context().Done():
			return nil
		case m := <-s.pool.Matches():
			grpcMatch := &pb.PlayerMatch{
				P1: &pb.Player{Id: m.Player1.Identify()},
				P2: &pb.Player{Id: m.Player2.Identify()},
			}
			return matches.Send(grpcMatch)
		}
	}
}

// Remove implements pb.PoolServer
func (s *grpcServer) Remove(ctx context.Context, player *pb.Player) (*pb.Empty, error) {
	select {
	case <-ctx.Done():
		return nil, nil
	default:
		s.pool.Remove(player)
		return nil, nil
	}
}

// Size implements pb.PoolServer
func (s *grpcServer) Size(ctx context.Context, _ *pb.Empty) (*pb.SizeResponse, error) {
	select {
	case <-ctx.Done():
		return nil, nil
	default:
		return &pb.SizeResponse{
			Size: int32(s.pool.Size()),
		}, nil
	}
}

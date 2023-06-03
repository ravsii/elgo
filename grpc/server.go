package grpc

import (
	context "context"

	"github.com/ravsii/elgo"
	schema "github.com/ravsii/elgo/grpc/schema"
)

// ensuring we've implemented our server correctly
var _ schema.PoolServer = (*grpcServer)(nil)

type grpcServer struct {
	pool *elgo.Pool
	schema.UnimplementedPoolServer
}

func NewGrpcServer(poolOpts ...elgo.PoolOpt) *grpcServer {
	return &grpcServer{
		pool: elgo.NewPool(poolOpts...),
	}
}

// Add implements schema.PoolServer
func (s *grpcServer) Add(ctx context.Context, player *schema.Player) (*schema.Empty, error) {
	select {
	case <-ctx.Done():
		return &schema.Empty{}, nil
	default:
		if err := s.pool.AddPlayer(player); err != nil {
			return &schema.Empty{}, err
		}
		return &schema.Empty{}, nil
	}
}

// Match implements schema.PoolServer
func (s *grpcServer) Match(_ *schema.Empty, matches schema.Pool_MatchServer) error {
	for {
		select {
		case <-matches.Context().Done():
			return nil
		case m := <-s.pool.Matches():
			grpcMatch := &schema.PlayerMatch{
				P1: &schema.Player{Id: m.Player1.Identify()},
				P2: &schema.Player{Id: m.Player2.Identify()},
			}
			return matches.Send(grpcMatch)
		}
	}
}

// Remove implements schema.PoolServer
func (s *grpcServer) Remove(ctx context.Context, player *schema.Player) (*schema.Empty, error) {
	select {
	case <-ctx.Done():
		return nil, nil
	default:
		s.pool.Remove(player)
		return nil, nil
	}
}

// Size implements schema.PoolServer
func (s *grpcServer) Size(ctx context.Context, _ *schema.Empty) (*schema.SizeResponse, error) {
	select {
	case <-ctx.Done():
		return nil, nil
	default:
		return &schema.SizeResponse{
			Size: int32(s.pool.Size()),
		}, nil
	}
}

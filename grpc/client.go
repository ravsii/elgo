package grpc

import (
	context "context"
	"fmt"

	"github.com/ravsii/elgo"
	"github.com/ravsii/elgo/grpc/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
)

type grpcClient struct {
	conn   *grpc.ClientConn
	client pb.PoolClient
}

// NewClient returns a new client, a closeFunc for closing the connection.
func NewClient(addr string, opts ...grpc.DialOption) (*grpcClient, error) {
	credentials := insecure.NewCredentials()
	opts = append(opts, grpc.WithTransportCredentials(credentials))
	conn, err := grpc.Dial(addr, opts...)
	if err != nil {
		return nil, fmt.Errorf("grpc dial: %w", err)
	}

	pbClient := pb.NewPoolClient(conn)

	return &grpcClient{
		conn:   conn,
		client: pbClient,
	}, nil
}

// Add adds a player to the queue.
// If a player already in the queue, elgo.PlayerAlreadyExistsError is returned.
func (c *grpcClient) Add(ctx context.Context, players ...elgo.Player) error {
	for _, player := range players {
		_, err := c.client.Add(ctx, &pb.Player{Id: player.Identify(), Elo: player.Rating()})
		if err != nil {
			errStatus := status.Convert(err)
			if errStatus.Code() == codes.AlreadyExists {
				return elgo.NewAlreadyExistsErr(player)
			}

			return err
		}
	}

	return nil
}

// RecieveMatch waits for a match to be created and returns it. This is a
// blocking operation and ctx.Done() or other contexts could be used to
// about it.
//
// ErrNoMatchFound is returned in case of timeouts.
//
// Note: unfortunately, it's hard to implement it with channels, because
// some matches could be lost on the process, i.e. if server randomly shuts
// down. So use a simple infinite for loop here.
func (c *grpcClient) RecieveMatch(ctx context.Context) (*elgo.Match, error) {
	matchClient, err := c.client.Match(ctx, &pb.Empty{})
	if err != nil {
		return nil, fmt.Errorf("match client: %w", err)
	}
	defer func() {
		if err == nil {
			err = matchClient.CloseSend()
		}
	}()

	match, err := matchClient.Recv()
	if err != nil {
		errStatus := status.Convert(err)
		if errStatus.Code() == codes.DeadlineExceeded {
			return nil, elgo.ErrNoMatchFound
		}

		return nil, fmt.Errorf("match recv: %w", err)
	}

	return &elgo.Match{
		Player1: match.GetP1(),
		Player2: match.GetP2(),
	}, nil
}

// Remove removes a player from the queue.
func (c *grpcClient) Remove(ctx context.Context, player elgo.Identifier) error {
	_, err := c.client.Remove(ctx, &pb.Player{Id: player.Identify()})
	return err
}

// Size returns a size of a queue.
func (c *grpcClient) Size(ctx context.Context) (int, error) {
	size, err := c.client.Size(ctx, &pb.Empty{})
	return int(size.GetSize()), err
}

func (c *grpcClient) Close() error {
	if err := c.conn.Close(); err != nil {
		return fmt.Errorf("client conn close: %w", err)
	}

	return nil
}

package grpcclient

import (
	"context"
	"time"

	pb "github.com/nikitadada/load-tester/internal/proto/gen"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type PingClient struct {
	conn   *grpc.ClientConn
	client pb.PingServiceClient
}

func NewPingClient(target string) (*PingClient, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	conn, err := grpc.DialContext(
		ctx,
		target,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock(),
	)
	if err != nil {
		return nil, err
	}

	client := pb.NewPingServiceClient(conn)

	return &PingClient{
		conn:   conn,
		client: client,
	}, nil
}

func (c *PingClient) Close() error {
	return c.conn.Close()
}

func (c *PingClient) Call(ctx context.Context) error {
	_, err := c.client.Ping(ctx, &pb.PingRequest{
		Message: "load-test",
	})
	return err
}

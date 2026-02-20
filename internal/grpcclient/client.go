package grpcclient

import (
	"context"
	"time"
)

type Client interface {
	Call(ctx context.Context) error
}

type DummyClient struct{}

func NewDummy() *DummyClient {
	return &DummyClient{}
}

func (c *DummyClient) Call(ctx context.Context) error {
	time.Sleep(25 * time.Millisecond) // имитация latency
	return nil
}

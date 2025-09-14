package stream

import (
    "context"
    "time"

    "google.golang.org/grpc"
    "google.golang.org/grpc/credentials/insecure"

    pb "github.com/jbyun0101/yellowstone-metrics-exporter/internal/proto/geyser"
)

type Client struct {
    conn   *grpc.ClientConn
    client pb.GeyserClient
}

// Dial connects to Dragon’s Mouth gRPC server (e.g. localhost:10000).
func Dial(addr string) (*Client, error) {
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    conn, err := grpc.DialContext(
        ctx,
        addr,
        grpc.WithTransportCredentials(insecure.NewCredentials()), // plaintext for localnet
    )
    if err != nil {
        return nil, err
    }

    return &Client{
        conn:   conn,
        client: pb.NewGeyserClient(conn),
    }, nil
}

func (c *Client) Close() error { return c.conn.Close() }

// StreamSlots subscribes to slot updates and calls handler for each new slot.
func (c *Client) StreamSlots(ctx context.Context, handler func(slot uint64)) error {
    stream, err := c.client.Subscribe(ctx) // generic Subscribe()
    if err != nil {
        return err
    }

    // minimal request: ask for slots
    req := &pb.SubscribeRequest{
        Slots: map[string]*pb.SubscribeRequestFilterSlots{
            "all": {}, // all slots, no filter
        },
    }

    if err := stream.Send(req); err != nil {
        return err
    }

    go func() {
        <-ctx.Done()
        stream.CloseSend()
    }()

    for {
        msg, err := stream.Recv()
        if err != nil {
            return err
        }
        if slotMsg := msg.GetSlot(); slotMsg != nil {
            handler(slotMsg.Slot)
        }
    }
}

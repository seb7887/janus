package grpc

import (
	"context"
	// to test
	"log"
	"sync"
	"time"

	"github.com/seb7887/janus/janusrpc"
)

type janusGRPCHandler struct {
	janusrpc.UnimplementedJanusServiceServer
}

func NewJanusGRPCServer() janusrpc.JanusServiceServer {
	return &janusGRPCHandler{}
}

func (h janusGRPCHandler) GetState(ctx context.Context, req *janusrpc.SingleStateRequest) (*janusrpc.StateResponse, error) {
	return &janusrpc.StateResponse{
		DeviceId: req.DeviceId,
		NodeId:   "pepe",
	}, nil
}

func (h janusGRPCHandler) StreamState(req *janusrpc.SingleStateRequest, stream janusrpc.JanusService_StreamStateServer) error {
	var wg sync.WaitGroup
	for i := 0; i < 3; i++ {
		wg.Add(1)
		go func(count int64) {
			defer wg.Done()
			time.Sleep(time.Duration(count) * time.Second)
			resp := janusrpc.StateResponse{
				DeviceId:    req.DeviceId,
				Temperature: count,
			}
			if err := stream.Send(&resp); err != nil {
				log.Printf("send error %s", err.Error())
			}
		}(int64(i))
	}

	wg.Wait()
	return nil
}

func (h janusGRPCHandler) GetNodeStates(ctx context.Context, req *janusrpc.MultipleStateRequest) (*janusrpc.MultipleStateResponse, error) {
	r := &janusrpc.StateResponse{
		DeviceId: "1",
		NodeId:   req.NodeId,
	}
	var s []*janusrpc.StateResponse
	s = append(s, r)
	return &janusrpc.MultipleStateResponse{States: s}, nil
}

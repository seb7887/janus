package grpc

import (
	"context"

	"github.com/seb7887/janus/internal/query"
	"github.com/seb7887/janus/janusrpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type janusGRPCHandler struct {
	stateService query.QueryServiceState
}

func NewJanusGRPCServer(stateService query.QueryServiceState) janusrpc.JanusServiceServer {
	return &janusGRPCHandler{
		stateService: stateService,
	}
}

func (h janusGRPCHandler) GetState(ctx context.Context, req *janusrpc.SingleStateRequest) (*janusrpc.StateResponse, error) {
	resp, err := h.stateService.GetDeviceState(req)
	if err != nil {
		return nil, status.Error(codes.NotFound, err.Error())
	}

	return resp, nil
}

func (h janusGRPCHandler) StreamState(req *janusrpc.SingleStateRequest, stream janusrpc.JanusService_StreamStateServer) error {
	return h.stateService.StateSubscription(req, stream)
}

func (h janusGRPCHandler) GetNodeStates(ctx context.Context, req *janusrpc.MultipleStateRequest) (*janusrpc.MultipleStateResponse, error) {
	resp, err := h.stateService.GetNodeStates(req)
	if err != nil {
		status.Errorf(codes.NotFound, err.Error())
	}

	return resp, nil
}

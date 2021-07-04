package grpc

import (
	"context"

	"github.com/seb7887/janus/internal/query"
	"github.com/seb7887/janus/janusrpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type janusGRPCHandler struct {
	stateService     query.QueryServiceState
	telemetryService query.QueryServiceTelemetry
}

func NewJanusGRPCServer(stateService query.QueryServiceState, telemetryService query.QueryServiceTelemetry) janusrpc.JanusServiceServer {
	return &janusGRPCHandler{
		stateService:     stateService,
		telemetryService: telemetryService,
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

func (h janusGRPCHandler) GetTelemetryTimeline(ctx context.Context, req *janusrpc.TimelineQuery) (*janusrpc.TimelineQueryResponse, error) {
	err := h.telemetryService.GetTimeline(req)
	if err != nil {
		status.Errorf(codes.Internal, err.Error())
	}

	item := &janusrpc.TimelineItem{
		Name:  "current",
		Count: 3,
	}
	var items []*janusrpc.TimelineItem
	items = append(items, item)

	return &janusrpc.TimelineQueryResponse{
		Items: items,
		Total: 1,
	}, nil
}

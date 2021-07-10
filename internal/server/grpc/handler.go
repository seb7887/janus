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
	res, err := h.telemetryService.GetTimeline(req)
	if err != nil {
		status.Errorf(codes.Internal, err.Error())
	}

	total, err := h.telemetryService.GetTotalSamples(req.Interval, req.Filters)
	if err != nil {
		status.Errorf(codes.Internal, err.Error())
	}

	return &janusrpc.TimelineQueryResponse{
		Result: res,
		Total:  int64(total),
	}, nil
}

func (h janusGRPCHandler) GetSegmentedTimeline(ctx context.Context, req *janusrpc.SegmentedTimelineQuery) (*janusrpc.TimelineQueryResponse, error) {
	var res []*janusrpc.TimelineResponse
	return &janusrpc.TimelineQueryResponse{
		Result: res,
		Total:  0,
	}, nil
}

func (h janusGRPCHandler) GetSegmentQuery(ctx context.Context, req *janusrpc.SegmentQuery) (*janusrpc.SegmentedQueryResponse, error) {
	segmentItems, err := h.telemetryService.GetSegments(req)
	if err != nil {
		status.Errorf(codes.Internal, err.Error())
	}

	total, err := h.telemetryService.GetTotalSamples(req.Interval, req.Filters)
	if err != nil {
		status.Errorf(codes.Internal, err.Error())
	}

	return &janusrpc.SegmentedQueryResponse{
		Segments: segmentItems,
		Total:    int64(total),
	}, nil
}

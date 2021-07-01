package query

import (
	"context"
	"fmt"

	mg "github.com/seb7887/janus/internal/storage/mongodb"
	"github.com/seb7887/janus/janusrpc"
	"go.mongodb.org/mongo-driver/mongo"
)

const (
	meter         = "meter"
	generator     = "generator"
	NOT_FOUND_MSG = "Device not found"
)

type StreamChMsg struct {
	MeterState     mg.Meter
	GeneratorState mg.Generator
}

var qschan = make(chan StreamChMsg, 10)

type QueryServiceState interface {
	GetDeviceState(r *janusrpc.SingleStateRequest) (*janusrpc.StateResponse, error)
	GetNodeStates(r *janusrpc.MultipleStateRequest) (*janusrpc.MultipleStateResponse, error)
	StateSubscription(req *janusrpc.SingleStateRequest, stream janusrpc.JanusService_StreamStateServer) error
}

type queryServiceState struct{}

func NewQueryStateService() QueryServiceState {
	return &queryServiceState{}
}

func formatMeterStateResponse(m *mg.Meter) *janusrpc.StateResponse {
	return &janusrpc.StateResponse{
		DeviceId:       m.DeviceId,
		NodeId:         m.NodeId,
		Temperature:    m.Temperature,
		Consumption:    m.Consumption,
		EnergyConsumed: m.EnergyConsumed,
		LastReport:     m.LastReport,
		Conected:       m.Connected,
	}
}

func formatGeneratorStateResponse(g *mg.Generator) *janusrpc.StateResponse {
	return &janusrpc.StateResponse{
		DeviceId:        g.DeviceId,
		NodeId:          g.NodeId,
		Temperature:     g.Temperature,
		EnergyGenerated: g.EnergyGenerated,
		Enabled:         g.Enabled,
		NeedManteinance: g.NeedManteinance,
		LastManteinance: g.LastManteinance,
	}
}

func (qs *queryServiceState) GetDeviceState(r *janusrpc.SingleStateRequest) (*janusrpc.StateResponse, error) {
	var resp *janusrpc.StateResponse
	var err error

	if r.DeviceType == meter {
		m, err := mg.GetMeter(r.DeviceId)
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf(NOT_FOUND_MSG)
		}
		resp = formatMeterStateResponse(m)
	} else if r.DeviceType == generator {
		g, err := mg.GetGenerator(r.DeviceId)
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf(NOT_FOUND_MSG)
		}
		resp = formatGeneratorStateResponse(g)
	} else {
		return nil, fmt.Errorf("Unknown device type")
	}

	return resp, err
}

func (qs *queryServiceState) GetNodeStates(r *janusrpc.MultipleStateRequest) (*janusrpc.MultipleStateResponse, error) {
	var states []*janusrpc.StateResponse

	m, err := mg.GetNodeMeters(r.NodeId)
	if err != nil {
		return nil, err
	}
	for i := range m {
		states = append(states, formatMeterStateResponse(m[i]))
	}

	g, err := mg.GetNodeGenerators(r.NodeId)
	if err != nil {
		return nil, err
	}
	for j := range g {
		states = append(states, formatGeneratorStateResponse(g[j]))
	}

	return &janusrpc.MultipleStateResponse{States: states}, nil
}

func (qs *queryServiceState) StateSubscription(req *janusrpc.SingleStateRequest, stream janusrpc.JanusService_StreamStateServer) error {
	ctx := context.Background()
	for {
		select {
		case msg := <-qschan:
			if req.DeviceType == meter {
				meterState := msg.MeterState
				if meterState.DeviceId == req.DeviceId {
					stream.Send(formatMeterStateResponse(&msg.MeterState))
				}
			} else if req.DeviceType == generator {
				genState := msg.GeneratorState
				if genState.DeviceId == req.DeviceId {
					stream.Send(formatGeneratorStateResponse(&msg.GeneratorState))
				}
			}
		case <-ctx.Done():
			return nil
		}
	}
}

// Concurrency in Golang:
// If you try to read data from an empty channel, then the goroutine will be blocked
// If you try to write to a channel that already has some data, then the goroutine will be blocked
func StreamState(s *StreamChMsg) {
	qschan <- *s
}

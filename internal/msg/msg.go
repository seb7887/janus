package msg

import (
	"github.com/seb7887/janus/internal/storage/mongodb"
)

type Msg struct {
	ClientId  string
	Topic     string
	Payload   string
	Timestamp int64
}

type StateMsg struct {
	DeviceType      string
	NodeId          string
	Temperature     int64
	Consumption     int64
	EnergyConsumed  int64
	LastReport      int64
	Connected       bool
	EnergyGenerated int64
	Enabled         bool
	NeedManteinance bool
	LastManteinance int64
}

func GetMeterState(id string, state *StateMsg) mongodb.Meter {
	return mongodb.Meter{
		DeviceId:       id,
		NodeId:         state.NodeId,
		Temperature:    state.Temperature,
		Consumption:    state.Consumption,
		EnergyConsumed: state.EnergyConsumed,
		LastReport:     state.LastReport,
		Connected:      state.Connected,
	}
}

func GetGeneratorState(id string, state *StateMsg) mongodb.Generator {
	return mongodb.Generator{
		DeviceId:        id,
		NodeId:          state.NodeId,
		Temperature:     state.Temperature,
		EnergyGenerated: state.EnergyGenerated,
		Enabled:         state.Enabled,
		NeedManteinance: state.NeedManteinance,
		LastManteinance: state.LastManteinance,
	}
}

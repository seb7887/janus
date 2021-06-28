package msg

import (
	"time"

	"github.com/seb7887/janus/internal/storage/mongodb"
	ts "github.com/seb7887/janus/internal/storage/timescaledb"
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

type TelemetryMsg struct {
	DeviceType  string
	MsgType     string
	Power       int64
	Voltage     int64
	Current     int64
	Temperature int64
	Subject     string
	Message     string
}

func GetTelemetryMsg(id string, msg *TelemetryMsg, millis int64) ts.Telemetry {
	return ts.Telemetry{
		DeviceId:    id,
		DeviceType:  msg.DeviceType,
		Power:       msg.Power,
		Voltage:     msg.Voltage,
		Current:     msg.Current,
		Temperature: msg.Temperature,
		Timestamp:   time.Unix(0, millis*int64(time.Millisecond)),
	}
}

func GetLogMsg(id string, msg *TelemetryMsg, millis int64) ts.Log {
	return ts.Log{
		DeviceId:  id,
		Subject:   msg.Subject,
		Message:   msg.Message,
		Timestamp: time.Unix(0, millis*int64(time.Millisecond)),
	}
}
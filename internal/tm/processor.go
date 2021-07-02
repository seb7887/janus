package tm

import (
	"encoding/json"

	m "github.com/seb7887/janus/internal/msg"
	ts "github.com/seb7887/janus/internal/storage/timescaledb"
	log "github.com/sirupsen/logrus"
)

const (
	telemetryMsg = "telemetry"
	logMsg       = "log"
)

func ProcessTelemetryMsg(msg *m.Msg) {
	log.Debugf("telemetry msg %s", msg)
	var payload m.TelemetryMsg
	err := json.Unmarshal([]byte(msg.Payload), &payload)
	if err != nil {
		log.Errorf("error parsing telemetry payload %s", err.Error())
	}

	if payload.MsgType == telemetryMsg {
		tmRow := m.GetTelemetryMsg(msg.ClientId, &payload)
		err = ts.InsertTelemetryEntry(&tmRow)
	} else if payload.MsgType == logMsg {
		logRow := m.GetLogMsg(msg.ClientId, &payload)
		err = ts.InsertLogEntry(&logRow)
	}
}

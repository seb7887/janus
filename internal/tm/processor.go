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

// channel to receive telemetry messages
var tchan = make(chan m.Msg, 10)

func StartTelemetryListener() error {
	for {
		select {
		case msg := <-tchan:
			log.Debugf("telemetry msg %s", msg)
			var payload m.TelemetryMsg
			err := json.Unmarshal([]byte(msg.Payload), &payload)
			if err != nil {
				log.Errorf("error parsing telemetry payload %s", err.Error())
			}

			if payload.MsgType == telemetryMsg {
				tmRow := m.GetTelemetryMsg(msg.ClientId, &payload, msg.Timestamp)
				err = ts.InsertTelemetryEntry(&tmRow)
			} else if payload.MsgType == logMsg {
				logRow := m.GetLogMsg(msg.ClientId, &payload, msg.Timestamp)
				err = ts.InsertLogEntry(&logRow)
			}

			return err
		}
	}
}

func ProcessTelemetryMsg(msg *m.Msg) {
	tchan <- *msg
}

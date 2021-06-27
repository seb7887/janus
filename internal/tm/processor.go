package tm

import (
	m "github.com/seb7887/janus/internal/msg"
	log "github.com/sirupsen/logrus"
)

// channel to receive telemetry messages
var tchan = make(chan m.Msg, 10)

func StartTelemetryListener() error {
	log.Info("Start Telemetry listener")
	for {
		select {
		case msg := <-tchan:
			log.Infof("telemetry msg %s", msg)
		}
	}
}

func ProcessTelemetryMsg(msg *m.Msg) {
	tchan <- *msg
}

package st

import (
	"encoding/json"

	m "github.com/seb7887/janus/internal/msg"
	"github.com/seb7887/janus/internal/storage/mongodb"
	log "github.com/sirupsen/logrus"
)

const (
	meter     = "meter"
	generator = "generator"
)

// channel to receive state messages
var schan = make(chan m.Msg, 10)

func StartStateListener() error {
	for {
		select {
		case msg := <-schan:
			log.Debugf("state msg %s", msg)
			var payload m.StateMsg
			err := json.Unmarshal([]byte(msg.Payload), &payload)
			if err != nil {
				log.Errorf("error parsing state payload %s")
			}

			if payload.DeviceType == meter {
				doc := m.GetMeterState(msg.ClientId, &payload)
				err = mongodb.UpsertMeter(doc)
			} else if payload.DeviceType == generator {
				doc := m.GetGeneratorState(msg.ClientId, &payload)
				err = mongodb.UpsertGenerator(doc)
			}

			return err
		}
	}
}

func ProcessStateMsg(msg *m.Msg) {
	schan <- *msg
}

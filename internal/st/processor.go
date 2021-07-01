package st

import (
	"encoding/json"

	m "github.com/seb7887/janus/internal/msg"
	"github.com/seb7887/janus/internal/query"
	"github.com/seb7887/janus/internal/storage/mongodb"
	log "github.com/sirupsen/logrus"
)

const (
	meter     = "meter"
	generator = "generator"
)

func ProcessStateMsg(msg *m.Msg) {
	log.Debugf("state msg %s", msg)
	var payload m.StateMsg
	err := json.Unmarshal([]byte(msg.Payload), &payload)
	if err != nil {
		log.Errorf("error parsing state payload %s", err.Error())
	}

	if payload.DeviceType == meter {
		doc := m.GetMeterState(msg.ClientId, &payload)
		err = mongodb.UpsertMeter(doc)
		if err == nil {
			// Stream new state to subscribers
			query.StreamState(&query.StreamChMsg{MeterState: doc})
		}
	} else if payload.DeviceType == generator {
		doc := m.GetGeneratorState(msg.ClientId, &payload)
		err = mongodb.UpsertGenerator(doc)
		if err == nil {
			// Stream new state to subscribers
			query.StreamState(&query.StreamChMsg{GeneratorState: doc})
		}
	}

	if err != nil {
		log.Error(err.Error())
	}
}

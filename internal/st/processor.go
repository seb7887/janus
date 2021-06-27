package st

import (
	m "github.com/seb7887/janus/internal/msg"
	log "github.com/sirupsen/logrus"
)

// channel to receive state messages
var schan = make(chan m.Msg, 10)

func StartStateListener() error {
	log.Info("Start State listener")
	for {
		select {
		case msg := <-schan:
			log.Infof("state msg %s", msg)
		}
	}
}

func ProcessStateMsg(msg *m.Msg) {
	schan <- *msg
}

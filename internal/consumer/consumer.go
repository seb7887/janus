package consumer

import (
	"encoding/json"
	"strings"

	"github.com/seb7887/janus/internal/config"
	m "github.com/seb7887/janus/internal/msg"
	"github.com/seb7887/janus/internal/st"
	"github.com/seb7887/janus/internal/tm"
	log "github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
)

const (
	consumer  = "consumer"
	state     = "state"
	telemetry = "telemetry"
)

func InitConsumer() error {
	// conn
	conn, err := amqp.Dial(config.GetConfig().AMQPUrl)
	if err != nil {
		log.Fatalf("ERROR: fail init RabbitMQ consumer %s", err.Error())
		return err
	}

	// create channel
	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("ERROR: fail to create a channel %s", err.Error())
		return err
	}

	// create queue
	queue, err := ch.QueueDeclare(
		consumer, // channel name
		true,     // durable
		false,    // delete when unused
		false,    // exclusive
		false,    // no-wait
		nil,      // arguments
	)
	if err != nil {
		log.Fatalf("ERROR: fail to create a queue %s", err.Error())
		return err
	}

	// channel
	msgChannel, err := ch.Consume(
		queue.Name, // queue
		"",         // consumer
		false,      // auto-ack
		false,      // exclusive
		false,      // no-local
		false,      // no-wait
		nil,        // arguments
	)
	if err != nil {
		log.Fatalf("ERROR: fail to create a message channel %s", err.Error())
		return err
	}

	// consume
	for {
		select {
		case msg := <-msgChannel:
			log.Debugf("received msg: %s", msg.Body)

			// parse message
			var parsedMsg m.Msg
			err = json.Unmarshal(msg.Body, &parsedMsg)
			if err != nil {
				log.Errorf("fail to parse message %s", err.Error())
			}

			// ack for message
			err = msg.Ack(true)
			if err != nil {
				log.Errorf("fail to ack: %s", err.Error())
				return err
			}

			// handle message
			handleMsg(&parsedMsg)
		}
	}
}

func handleMsg(msg *m.Msg) {
	log.Infof("TOPIC %s", msg.Topic)

	if strings.Contains(msg.Topic, state) {
		st.ProcessStateMsg(msg)
	} else if strings.Contains(msg.Topic, telemetry) {
		tm.ProcessTelemetryMsg(msg)
	}
}

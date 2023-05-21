package queue

import (
	"log"
	"os"

	"github.com/nats-io/nats.go"
)

var nc *nats.Conn
var js nats.JetStreamContext

func LoadNats() {
	url := os.Getenv("NATS_URL")
	client, err := nats.Connect(url)
	if err != nil {
		log.Fatal(err)
	}
	nc = client

	jet, err := nc.JetStream()
	if err != nil {
		log.Fatal(err)
	}
	js = jet

	cfg := &nats.StreamConfig{
		Name:      "WARP",
		Retention: nats.WorkQueuePolicy,
		Subjects:  []string{"store"},
	}
	js.AddStream(cfg)
}

func GetStream() nats.JetStreamContext {
	return js
}

func Drain() {
	nc.Drain()
}

func Publish(subject string, data []byte) {
	js.Publish(subject, data)
}

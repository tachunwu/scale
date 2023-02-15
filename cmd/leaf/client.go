package main

import (
	"log"

	"github.com/nats-io/nats.go"
)

func Client() {

	nc, err := nats.Connect(nats.DefaultURL)
	defer nc.Close()
	if err != nil {
		log.Fatal(err)
	}
	js, err := nc.JetStream(nats.PublishAsyncMaxPending(256))
	if err != nil {
		log.Fatal(err)
	}

	js.AddStream(&nats.StreamConfig{
		Name:     "PEBBLE",
		Subjects: []string{"pebble"},
		Storage:  nats.FileStorage,
	})

	for i := 0; i < 10; i++ {
		js.PublishAsync("pebble", []byte("Hello JS Async!"))
	}

}

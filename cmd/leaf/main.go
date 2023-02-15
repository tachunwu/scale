package main

import (
	"log"

	natsserver "github.com/nats-io/nats-server/v2/server"
	"github.com/nats-io/nats.go"
	"github.com/tachunwu/scale/pkg/database"
	"github.com/tachunwu/scale/pkg/jetstream"
	"go.uber.org/zap"
)

func main() {
	app := NewLeafApp()
	Client()
	app.Run()
}

type LeafApp struct {
	db     database.DB
	js     *natsserver.Server
	logger *zap.Logger
}

func NewLeafApp() *LeafApp {
	logger, _ := zap.NewProduction()
	defer logger.Sync()
	return &LeafApp{
		db:     database.NewPebbleDB("./data"),
		js:     jetstream.NewJetStream(),
		logger: logger,
	}
}

func (a *LeafApp) Run() {

	nc, err := nats.Connect(nats.DefaultURL)
	if err != nil {
		log.Fatal(err)
	}
	defer nc.Close()
	js, err := nc.JetStream(nats.PublishAsyncMaxPending(256))

	sub, err := js.PullSubscribe("pebble", "PEBBLE")
	for {
		msgs, err := sub.Fetch(1)
		if err != nil {
			a.logger.Warn("consumer error", zap.Any("error", err))
		}
		for i := range msgs {
			a.logger.Info("Consumer", zap.Binary("data", msgs[i].Data))
		}
	}

}

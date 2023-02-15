package main

// 45.059151905 sec process 40,0000 (TPS: 8877.21990071)
// 74.768570822 sec process 80,0000 (TPS: 10699.6829176)

import (
	"log"

	"github.com/cockroachdb/pebble"
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
	if err != nil {
		a.logger.Error("JetStream error", zap.Any("error", err))
		return
	}

	kv, err := js.KeyValue(jetstream.BucketName)
	if err != nil {
		a.logger.Error("KV error", zap.Any("error", err))
		return
	}

	w, _ := kv.WatchAll()
	kveCh := make(chan nats.KeyValueEntry, 1)
	// Producer
	go func() {
		for {
			kve := <-w.Updates()
			kveCh <- kve
		}

	}()
	// Consumer
	for {
		kve := <-kveCh
		if kve != nil {
			b := a.db.NewBatch()
			b.Set([]byte(kve.Key()), kve.Value(), pebble.Sync)
			b.Commit(pebble.Sync)
			// a.logger.Info(
			// 	"Pebble store",
			// 	zap.String("Key", kve.Key()),
			// 	zap.Binary("Value", kve.Value()),
			// )
		}

	}
}

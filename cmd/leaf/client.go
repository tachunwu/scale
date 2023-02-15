package main

import (
	"log"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/nats-io/nats.go"
	"github.com/tachunwu/scale/pkg/jetstream"
)

var wg sync.WaitGroup

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
	// Test

	wg.Add(8)
	start := time.Now()
	go MicroBench(js)
	go MicroBench(js)
	go MicroBench(js)
	go MicroBench(js)
	go MicroBench(js)
	go MicroBench(js)
	go MicroBench(js)
	go MicroBench(js)
	wg.Wait()
	log.Println(time.Since(start).Seconds())

}

func MicroBench(js nats.JetStreamContext) {
	for i := 0; i < 100000; i++ {
		kv, _ := js.KeyValue(jetstream.BucketName)
		k := uuid.New().String()
		kv.Put(k, []byte("value"))
	}
	defer wg.Done()
}

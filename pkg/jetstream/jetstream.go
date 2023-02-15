package jetstream

import (
	"github.com/nats-io/nats-server/v2/server"
	natsserver "github.com/nats-io/nats-server/v2/server"
	"github.com/nats-io/nats.go"
	"go.uber.org/zap"
)

var BucketName string = "materialize"

func NewJetStream() *natsserver.Server {
	// Enable logger
	logger, _ := zap.NewProduction()
	defer logger.Sync()

	// Setting server
	opts := &server.Options{
		JetStream: true,
		NoLog:     true,
		Trace:     true,
		NoSigs:    true,
		StoreDir:  "./js",
	}

	// Initialize new server with options
	ns, err := server.NewServer(opts)

	if err != nil {
		panic(err)
	}

	// Start the server via goroutine
	ns.ConfigureLogger()
	ns.Start()
	logger.Info(
		"Embedded JetStream started",
		zap.String("ClientURL", ns.ClientURL()),
	)
	// Create KV Bucket
	nc, _ := nats.Connect(nats.DefaultURL)
	js, _ := nc.JetStream()
	js.CreateKeyValue(&nats.KeyValueConfig{
		Bucket: BucketName,
	})

	logger.Info(
		"JetStream KV Bucket created",
		zap.String("bucket_name", BucketName),
	)
	return ns
}

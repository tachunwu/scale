package jetstream

import (
	"github.com/nats-io/nats-server/v2/server"
	natsserver "github.com/nats-io/nats-server/v2/server"
	"go.uber.org/zap"
)

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
		"Embedded info",
		zap.String("ClientURL", ns.ClientURL()),
	)
	return ns
}

package runner

import (
	"github.com/nats-io/nats.go"
)

func NewTestRunner(nats *nats.Conn, jetstream nats.JetStreamContext) *Runner {
	initializeConfiguration()

	return &Runner{
		logger:    getLogger(),
		nats:      nats,
		jetstream: jetstream,
	}
}

package messaging

import (
	"github.com/go-logr/logr"
	kai "github.com/konstellation-io/kre-runners/go-sdk/v1/protos"
	"github.com/nats-io/nats.go"
)

func NewTestMessaging(logger logr.Logger, nats *nats.Conn, jetstream nats.JetStreamContext,
	requestMessage *kai.KaiNatsMessage, messageUtils messageUtils) *Messaging {
	return &Messaging{
		logger.WithName("[MESSAGING]"),
		nats,
		jetstream,
		requestMessage,
		messageUtils,
	}
}

func (ms Messaging) SendError(requestID, errorMessage string) {
	ms.publishError(requestID, errorMessage)
}

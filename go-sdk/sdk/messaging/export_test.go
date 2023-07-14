package messaging

import (
	"github.com/go-logr/logr"
	kai "github.com/konstellation-io/kai-sdk/go-sdk/protos"
	"github.com/nats-io/nats.go"
)

func NewTestMessaging(logger logr.Logger, ns *nats.Conn, js nats.JetStreamContext,
	requestMessage *kai.KaiNatsMessage, messageUtils messageUtils,
) *Messaging {
	//nolint: gocritic
	return &Messaging{
		logger.WithName("[MESSAGING]"),
		ns,
		js,
		requestMessage,
		messageUtils,
	}
}

func (ms Messaging) SendError(requestID, errorMessage string) {
	ms.publishError(requestID, errorMessage)
}

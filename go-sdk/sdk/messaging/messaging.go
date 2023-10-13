package messaging

import (
	"github.com/go-logr/logr"
	kai "github.com/konstellation-io/kai-sdk/go-sdk/protos"
	"github.com/nats-io/nats.go"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"
)

type Scope string

const (
	GlobalScope   Scope = "global"
	ProductScope  Scope = "product"
	WorkflowScope Scope = "workflow"
	ProcessScope  Scope = "process"
)

type Messaging struct {
	logger         logr.Logger
	nats           *nats.Conn
	jetstream      nats.JetStreamContext
	requestMessage *kai.KaiNatsMessage
	messagingUtils messagingUtils
}

func NewMessaging(logger logr.Logger, ns *nats.Conn, js nats.JetStreamContext,
	requestMessage *kai.KaiNatsMessage,
) *Messaging {
	return &Messaging{
		logger.WithName("[MESSAGING]"),
		ns,
		js,
		requestMessage,
		NewMessagingUtils(ns, js),
	}
}

func (ms Messaging) SendOutput(response proto.Message, channelOpt ...string) error {
	return ms.publishMsg(response, ms.requestMessage.GetRequestId(), kai.MessageType_OK, ms.getOptionalString(channelOpt))
}

func (ms Messaging) SendOutputWithRequestID(response proto.Message, requestID string, channelOpt ...string) error {
	return ms.publishMsg(response, requestID, kai.MessageType_OK, ms.getOptionalString(channelOpt))
}

func (ms Messaging) SendAny(response *anypb.Any, channelOpt ...string) {
	ms.publishAny(response, ms.requestMessage.GetRequestId(), kai.MessageType_OK, ms.getOptionalString(channelOpt))
}

func (ms Messaging) SendAnyWithRequestID(response *anypb.Any, requestID string, channelOpt ...string) {
	ms.publishAny(response, requestID, kai.MessageType_OK, ms.getOptionalString(channelOpt))
}

// TODO remove this.
//
//nolint:godox // Task to be done.
func (ms Messaging) SendEarlyReply(response proto.Message, channelOpt ...string) error {
	return ms.publishMsg(response, ms.requestMessage.GetRequestId(), kai.MessageType_EARLY_REPLY, ms.getOptionalString(channelOpt))
}

// TODO remove this.
//
//nolint:godox // Task to be done.
func (ms Messaging) SendEarlyExit(response proto.Message, channelOpt ...string) error {
	return ms.publishMsg(response, ms.requestMessage.GetRequestId(), kai.MessageType_EARLY_EXIT, ms.getOptionalString(channelOpt))
}

func (ms Messaging) GetErrorMessage() string {
	if ms.IsMessageError() {
		return ms.requestMessage.GetError()
	}

	return ""
}

func (ms Messaging) IsMessageOK() bool {
	return ms.requestMessage.MessageType == kai.MessageType_OK
}

func (ms Messaging) IsMessageError() bool {
	return ms.requestMessage.MessageType == kai.MessageType_ERROR
}

// TODO remove this.
//
//nolint:godox // Task to be done.
func (ms Messaging) IsMessageEarlyReply() bool {
	return ms.requestMessage.MessageType == kai.MessageType_EARLY_REPLY
}

// TODO remove this.
//
//nolint:godox // Task to be done.
func (ms Messaging) IsMessageEarlyExit() bool {
	return ms.requestMessage.MessageType == kai.MessageType_EARLY_EXIT
}

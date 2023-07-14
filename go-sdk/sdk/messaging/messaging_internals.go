package messaging

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/spf13/viper"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"

	"github.com/konstellation-io/kai-sdk/go-sdk/internal/common"
	"github.com/konstellation-io/kai-sdk/go-sdk/internal/errors"
	kai "github.com/konstellation-io/kai-sdk/go-sdk/protos"
)

const (
	defaultValue = ""
)

func (ms Messaging) getOptionalString(values []string) string {
	if len(values) > 0 {
		return values[0]
	}
	return defaultValue
}

func (ms Messaging) publishMsg(msg proto.Message, requestID string, msgType kai.MessageType, channel string) error {
	payload, err := anypb.New(msg)
	if err != nil {
		return fmt.Errorf("the handler result is not a valid protobuf: %s", err)
	}
	if requestID == "" {
		requestID = uuid.New().String()
	}
	responseMsg := ms.newResponseMsg(payload, requestID, msgType)

	ms.publishResponse(responseMsg, channel)

	return nil
}

func (ms Messaging) publishAny(payload *anypb.Any, requestID string, msgType kai.MessageType, channel string) {
	if requestID == "" {
		requestID = uuid.New().String()
	}
	responseMsg := ms.newResponseMsg(payload, requestID, msgType)
	ms.publishResponse(responseMsg, channel)
}

func (ms Messaging) publishError(requestID, errMsg string) {
	responseMsg := &kai.KaiNatsMessage{
		RequestId:   requestID,
		Error:       errMsg,
		FromNode:    viper.GetString("metadata.process_id"),
		MessageType: kai.MessageType_ERROR,
	}
	ms.publishResponse(responseMsg, "")
}

func (ms Messaging) newResponseMsg(payload *anypb.Any, requestID string,
	msgType kai.MessageType,
) *kai.KaiNatsMessage {
	ms.logger.V(1).Info("Preparing response message",
		"requestID", requestID, "msgType", msgType)

	return &kai.KaiNatsMessage{
		RequestId:   requestID,
		Payload:     payload,
		FromNode:    viper.GetString("metadata.process_id"),
		MessageType: msgType,
	}
}

func (ms Messaging) publishResponse(responseMsg *kai.KaiNatsMessage, channel string) {
	outputSubject := ms.getOutputSubject(channel)

	outputMsg, err := proto.Marshal(responseMsg)
	if err != nil {
		ms.logger.Error(err, "Error generating output result because handler result is not "+
			"a serializable Protobuf")
		return
	}

	outputMsg, err = ms.prepareOutputMessage(outputMsg)
	if err != nil {
		ms.logger.Error(err, "Error preparing output msg")
		return
	}

	ms.logger.Info("Publishing response", "subject", outputSubject)

	_, err = ms.jetstream.Publish(outputSubject, outputMsg)
	if err != nil {
		ms.logger.Error(err, "Error publishing output")
	}
}

func (ms Messaging) getOutputSubject(channel string) string {
	outputSubject := viper.GetString("nats.output")
	if channel != "" {
		return fmt.Sprintf("%s.%s", outputSubject, channel)
	}
	return outputSubject
}

func (ms Messaging) prepareOutputMessage(msg []byte) ([]byte, error) {
	maxSize, err := ms.messageUtils.GetMaxMessageSize()
	if err != nil {
		return nil, fmt.Errorf("error getting max message size: %s", err)
	}

	lenMsg := int64(len(msg))
	if lenMsg <= maxSize {
		return msg, nil
	}

	ms.logger.V(1).Info("message exceeds maximum size allowed, compressing data")
	outMsg, err := common.CompressData(msg)
	if err != nil {
		return nil, err
	}

	lenOutMsg := int64(len(outMsg))
	if lenOutMsg > maxSize {
		ms.logger.V(1).Info("compressed message exceeds maximum size allowed",
			"Current Message Size", sizeInMB(lenOutMsg), "Max. Allowed Size", sizeInMB(maxSize))
		return nil, errors.ErrMessageToBig
	}

	ms.logger.Info("Message compressed", "Original size", sizeInMB(lenMsg),
		"Compressed size", sizeInMB(lenOutMsg))

	return outMsg, nil
}

package messaging

import (
	"fmt"

	"github.com/konstellation-io/kai-sdk/go-sdk/internal/errors"
	"github.com/nats-io/nats.go"

	"github.com/google/uuid"
	"github.com/spf13/viper"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"

	"github.com/konstellation-io/kai-sdk/go-sdk/internal/common"
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
		return fmt.Errorf("the handler result is not a valid protobuf: %s", err) //nolint:goerr113 // error is wrapped
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
		common.LoggerRequestID, requestID, "msgType", msgType)

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
			"a serializable Protobuf", common.LoggerRequestID, responseMsg.RequestId)
		return
	}

	outputMsg, err = ms.prepareOutputMessage(outputMsg)
	if err != nil {
		ms.logger.Error(err, "Error preparing output msg", common.LoggerRequestID, responseMsg.RequestId)
		return
	}

	ms.logger.Info("Publishing response", "subject", outputSubject,
		common.LoggerRequestID, responseMsg.RequestId)

	_, err = ms.jetstream.Publish(outputSubject, outputMsg)
	if err != nil {
		ms.logger.Error(err, "Error publishing output", common.LoggerRequestID, responseMsg.RequestId)
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
	maxSize, err := ms.messagingUtils.GetMaxMessageSize()
	if err != nil {
		return nil, fmt.Errorf("error getting max message size: %s", err) //nolint:goerr113 // error is wrapped
	}

	lenMsg := int64(len(msg))
	if lenMsg <= maxSize {
		return msg, nil
	}

	ms.logger.V(1).Info("Message exceeds maximum size allowed, compressing data")

	outMsg, err := common.CompressData(msg)
	if err != nil {
		return nil, err
	}

	lenOutMsg := int64(len(outMsg))
	if lenOutMsg > maxSize {
		ms.logger.V(1).Info("Compressed message exceeds maximum size allowed",
			"Current Message Size", sizeInMB(lenOutMsg), "Max. Allowed Size", sizeInMB(maxSize))
		return nil, errors.ErrMessageToBig
	}

	ms.logger.Info("Message compressed", "Original size", sizeInMB(lenMsg),
		"Compressed size", sizeInMB(lenOutMsg))

	return outMsg, nil
}

func (ms Messaging) GetRequestID(msg *nats.Msg) (string, error) {
	requestMsg := &kai.KaiNatsMessage{}

	data := msg.Data

	var err error
	if common.IsCompressed(data) {
		data, err = common.UncompressData(data)
		if err != nil {
			ms.logger.Error(err, "Error reading compressed message")
			return "", err
		}
	}

	err = proto.Unmarshal(data, requestMsg)
	if err != nil {
		ms.logger.Error(err, "Error unmarshalling message")
		return "", err
	}

	return requestMsg.GetRequestId(), err
}

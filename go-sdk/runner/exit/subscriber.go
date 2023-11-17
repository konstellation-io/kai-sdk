package exit

import (
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/go-logr/logr"

	"github.com/nats-io/nats.go"
	"github.com/spf13/viper"
	"google.golang.org/protobuf/proto"

	"github.com/konstellation-io/kai-sdk/go-sdk/internal/common"
	"github.com/konstellation-io/kai-sdk/go-sdk/internal/errors"
	kai "github.com/konstellation-io/kai-sdk/go-sdk/protos"
	"github.com/konstellation-io/kai-sdk/go-sdk/sdk"
)

const _subscriberLoggerName = "[SUBSCRIBER]"

func (er *Runner) getLoggerWithName() logr.Logger {
	return er.sdk.Logger.WithName(_subscriberLoggerName)
}

func (er *Runner) startSubscriber() {
	inputSubjects := viper.GetStringSlice(common.ConfigNatsInputsKey)

	if len(inputSubjects) == 0 {
		er.getLoggerWithName().Info("Undefined input subjects")
		os.Exit(1)
	}

	subscriptions := make([]*nats.Subscription, 0, len(inputSubjects))

	for _, subject := range inputSubjects {
		consumerName := fmt.Sprintf("%s-%s", strings.ReplaceAll(subject, ".", "-"),
			strings.ReplaceAll(strings.ReplaceAll(er.sdk.Metadata.GetProcess(), ".", "-"), " ", "-"))

		er.getLoggerWithName().V(1).Info(fmt.Sprintf("Subscribing to subject %s with queue group %s", subject, consumerName))

		s, err := er.jetstream.QueueSubscribe(
			subject,
			consumerName,
			er.processMessage,
			nats.DeliverNew(),
			nats.Durable(consumerName),
			nats.ManualAck(),
			nats.AckWait(viper.GetDuration(common.ConfigRunnerSubscriberAckWaitTimeKey)),
		)
		if err != nil {
			er.getLoggerWithName().Error(err, fmt.Sprintf("Error subscribing to subject %s", subject))
			os.Exit(1)
		}

		subscriptions = append(subscriptions, s)

		er.getLoggerWithName().V(1).Info(fmt.Sprintf("Listening to subject %s with queue group %s", subject, consumerName))
	}

	er.sdk.Logger.Info("Subscribed to all subjects successfully")

	// Handle sigterm and await termChan signal
	termChan := make(chan os.Signal, 1)
	signal.Notify(termChan, syscall.SIGINT, syscall.SIGTERM)
	<-termChan

	// Handle shutdown
	er.getLoggerWithName().Info("Shutdown signal received")

	er.sdk.Logger.Info("Unsubscribing from all subjects")

	for _, s := range subscriptions {
		er.getLoggerWithName().V(1).Info(fmt.Sprintf("Unsubscribing from subject %s", s.Subject))

		err := s.Unsubscribe()
		if err != nil {
			er.getLoggerWithName().Error(err, fmt.Sprintf("Error unsubscribing from the subject %s", s.Subject))
			os.Exit(1)
		}
	}

	er.getLoggerWithName().Info("Unsubscribed from all subjects")
}

func (er *Runner) processMessage(msg *nats.Msg) {
	requestMsg, err := er.newRequestMessage(msg.Data)
	if err != nil {
		errMsg := fmt.Sprintf("Error parsing msg.data coming from subject %s because is not a valid protobuf: %s",
			msg.Subject, err)
		er.processRunnerError(msg, errMsg, requestMsg.RequestId)

		return
	}

	er.getLoggerWithName().Info(fmt.Sprintf("New message received with subject %s",
		msg.Subject), requestMsg.RequestId)

	handler := er.getResponseHandler(strings.ToLower(requestMsg.FromNode))
	if handler == nil {
		errMsg := fmt.Sprintf("Error missing handler for node %q", requestMsg.FromNode)
		er.processRunnerError(msg, errMsg, requestMsg.RequestId)

		return
	}

	// Make a shallow copy of the sdk object to set inside the request msg.
	hSdk := sdk.ShallowCopyWithRequest(&er.sdk, requestMsg)

	if er.preprocessor != nil {
		err := er.preprocessor(hSdk, requestMsg.Payload)
		if err != nil {
			errMsg := fmt.Sprintf("Error in node %q executing handler preprocessor for node %q: %s",
				er.sdk.Metadata.GetProcess(), requestMsg.FromNode, err)
			er.processRunnerError(msg, errMsg, requestMsg.RequestId)

			return
		}
	}

	err = handler(hSdk, requestMsg.Payload)
	if err != nil {
		errMsg := fmt.Sprintf("Error in node %q executing handler for node %q: %s",
			er.sdk.Metadata.GetProcess(), requestMsg.FromNode, err)
		er.processRunnerError(msg, errMsg, requestMsg.RequestId)

		return
	}

	if er.postprocessor != nil {
		err := er.postprocessor(hSdk, requestMsg.Payload)
		if err != nil {
			errMsg := fmt.Sprintf("Error in node %q executing handler postprocessor for node %q: %s",
				er.sdk.Metadata.GetProcess(), requestMsg.FromNode, err)
			er.processRunnerError(msg, errMsg, requestMsg.RequestId)

			return
		}
	}

	// Tell NATS we don't need to receive the message anymore, and we are done processing it.
	ackErr := msg.Ack()
	if ackErr != nil {
		er.getLoggerWithName().Error(ackErr, errors.ErrMsgAck)
	}
}

func (er *Runner) processRunnerError(msg *nats.Msg, errMsg, requestID string) {
	ackErr := msg.Ack()
	if ackErr != nil {
		er.getLoggerWithName().Error(ackErr, errors.ErrMsgAck)
	}

	er.getLoggerWithName().V(1).Info(errMsg)
	er.publishError(requestID, errMsg)
}

func (er *Runner) newRequestMessage(data []byte) (*kai.KaiNatsMessage, error) {
	requestMsg := &kai.KaiNatsMessage{}

	var err error
	if common.IsCompressed(data) {
		data, err = common.UncompressData(data)
		if err != nil {
			er.getLoggerWithName().Error(err, "Error reading compressed message")
			return nil, err
		}
	}

	err = proto.Unmarshal(data, requestMsg)

	return requestMsg, err
}

func (er *Runner) publishError(requestID, errMsg string) {
	responseMsg := &kai.KaiNatsMessage{
		RequestId:   requestID,
		Error:       errMsg,
		FromNode:    viper.GetString(common.ConfigMetadataProcessIDKey),
		MessageType: kai.MessageType_ERROR,
	}
	er.publishResponse(responseMsg, "")
}

func (er *Runner) publishResponse(responseMsg *kai.KaiNatsMessage, channel string) {
	outputSubject := er.getOutputSubject(channel)

	outputMsg, err := proto.Marshal(responseMsg)
	if err != nil {
		er.getLoggerWithName().
			Error(err, "Error generating output result because handler result is not a serializable Protobuf")
		return
	}

	outputMsg, err = er.prepareOutputMessage(outputMsg)
	if err != nil {
		er.getLoggerWithName().Error(err, "Error preparing output message")
		return
	}

	er.getLoggerWithName().V(1).Info(fmt.Sprintf("Publishing response with subject %s", outputSubject))

	_, err = er.jetstream.Publish(outputSubject, outputMsg)
	if err != nil {
		er.getLoggerWithName().Error(err, "Error publishing output")
	}
}

func (er *Runner) getOutputSubject(channel string) string {
	outputSubject := viper.GetString(common.ConfigNatsOutputKey)
	if channel != "" {
		return fmt.Sprintf("%s.%s", outputSubject, channel)
	}

	return outputSubject
}

// prepareOutputMessage will check the length of the message and compress it if necessary.
// Fails on compressed messages bigger than the threshold.
func (er *Runner) prepareOutputMessage(msg []byte) ([]byte, error) {
	maxSize, err := er.getMaxMessageSize()
	if err != nil {
		return nil, fmt.Errorf("error getting max message size: %s", err) //nolint:goerr113 // error is wrapped
	}

	lenMsg := int64(len(msg))
	if lenMsg <= maxSize {
		return msg, nil
	}

	outMsg, err := common.CompressData(msg)
	if err != nil {
		return nil, err
	}

	lenOutMsg := int64(len(outMsg))
	if lenOutMsg > maxSize {
		er.getLoggerWithName().V(1).Info("Compressed message size %s "+
			"exceeds maximum size allowed %s", sizeInMB(lenOutMsg), sizeInMB(maxSize))
		return nil, errors.ErrMessageToBig
	}

	er.getLoggerWithName().Info(fmt.Sprintf("Message prepared with original size %s "+
		"and compressed size %s", sizeInMB(lenMsg), sizeInMB(lenOutMsg)))

	return outMsg, nil
}

func (er *Runner) getResponseHandler(subject string) Handler {
	if responseHandler, ok := er.responseHandlers[subject]; ok {
		return responseHandler
	}

	// returns the default response handler, or nil if it doesn't exist
	return er.responseHandlers["default"]
}

func (er *Runner) getMaxMessageSize() (int64, error) {
	streamInfo, err := er.jetstream.StreamInfo(viper.GetString(common.ConfigNatsStreamKey))
	if err != nil {
		return 0, fmt.Errorf("error getting stream's max message size: %w", err)
	}

	streamMaxSize := int64(streamInfo.Config.MaxMsgSize)
	serverMaxSize := er.nats.MaxPayload()

	if streamMaxSize == -1 {
		return serverMaxSize, nil
	}

	if streamMaxSize < serverMaxSize {
		return streamMaxSize, nil
	}

	return serverMaxSize, nil
}

func sizeInMB(size int64) string {
	mbSize := float32(size) / 1024 / 1024
	return fmt.Sprintf("%.1f MB", mbSize)
}

package task

import (
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/nats-io/nats.go"
	"github.com/spf13/viper"
	"google.golang.org/protobuf/proto"

	"github.com/konstellation-io/kai-sdk/go-sdk/internal/common"
	"github.com/konstellation-io/kai-sdk/go-sdk/internal/errors"
	kai "github.com/konstellation-io/kai-sdk/go-sdk/protos"
	"github.com/konstellation-io/kai-sdk/go-sdk/sdk"
)

func (tr *Runner) startSubscriber() {
	inputSubjects := viper.GetStringSlice("nats.inputs")
	subscriptions := make([]*nats.Subscription, len(inputSubjects))

	for _, subject := range inputSubjects {
		consumerName := fmt.Sprintf("%s-%s", strings.ReplaceAll(subject, ".", "-"),
			strings.ReplaceAll(strings.ReplaceAll(tr.sdk.Metadata.GetProcess(), ".", "-"), " ", "-"))

		tr.sdk.Logger.WithName("[SUBSCRIBER]").V(1).Info("Subscribing to subject",
			"Subject", subject, "Queue group", consumerName)

		ackWaitTime := 22 * time.Hour

		s, err := tr.jetstream.QueueSubscribe(
			subject,
			consumerName,
			tr.processMessage,
			nats.DeliverNew(),
			nats.Durable(consumerName),
			nats.ManualAck(),
			nats.AckWait(ackWaitTime),
		)
		if err != nil {
			tr.sdk.Logger.WithName("[SUBSCRIBER]").Error(err, "Error subscribing to NATS subject",
				"Subject", subject)
			os.Exit(1)
		}

		subscriptions = append(subscriptions, s)

		tr.sdk.Logger.WithName("[SUBSCRIBER]").V(1).Info("Listening to subject",
			"Subject", subject, "Queue group", consumerName)
	}

	// Handle sigterm and await termChan signal
	termChan := make(chan os.Signal, 1)
	signal.Notify(termChan, syscall.SIGINT, syscall.SIGTERM)
	<-termChan

	// Handle shutdown
	tr.sdk.Logger.WithName("[SUBSCRIBER]").Info("Shutdown signal received")

	for _, s := range subscriptions {
		err := s.Unsubscribe()
		if err != nil {
			tr.sdk.Logger.WithName("[SUBSCRIBER]").Error(err, "Error unsubscribing from the NATS subject",
				"Subject", s.Subject)
			os.Exit(1)
		}
	}
}

func (tr *Runner) processMessage(msg *nats.Msg) {
	requestMsg, err := tr.newRequestMessage(msg.Data)
	if err != nil {
		errMsg := fmt.Sprintf("Error parsing msg.data coming from subject %s because is not a valid protobuf: %s", msg.Subject, err)
		tr.processRunnerError(msg, errMsg, requestMsg.RequestId)

		return
	}

	tr.sdk.Logger.WithName("[SUBSCRIBER]").Info("New message received",
		"Subject", msg.Subject, "Request ID", requestMsg.RequestId)

	handler := tr.getResponseHandler(strings.ToLower(requestMsg.FromNode))
	if handler == nil {
		errMsg := fmt.Sprintf("Error missing handler for node %q", requestMsg.FromNode)
		tr.processRunnerError(msg, errMsg, requestMsg.RequestId)

		return
	}

	// Make a shallow copy of the sdk object to set inside the request msg.
	hSdk := sdk.ShallowCopyWithRequest(&tr.sdk, requestMsg)

	if tr.preprocessor != nil {
		err := tr.preprocessor(hSdk, requestMsg.Payload)
		if err != nil {
			errMsg := fmt.Sprintf("Error in node %q executing handler preprocessor for node %q: %s",
				tr.sdk.Metadata.GetProcess(), requestMsg.FromNode, err)
			tr.processRunnerError(msg, errMsg, requestMsg.RequestId)

			return
		}
	}

	err = handler(hSdk, requestMsg.Payload)
	if err != nil {
		errMsg := fmt.Sprintf("Error in node %q executing handler for node %q: %s",
			tr.sdk.Metadata.GetProcess(), requestMsg.FromNode, err)
		tr.processRunnerError(msg, errMsg, requestMsg.RequestId)

		return
	}

	if tr.postprocessor != nil {
		err := tr.postprocessor(hSdk, requestMsg.Payload)
		if err != nil {
			errMsg := fmt.Sprintf("Error in node %q executing handler postprocessor for node %q: %s",
				tr.sdk.Metadata.GetProcess(), requestMsg.FromNode, err)
			tr.processRunnerError(msg, errMsg, requestMsg.RequestId)

			return
		}
	}

	// Tell NATS we don't need to receive the message anymore, and we are done processing it.
	ackErr := msg.Ack()
	if ackErr != nil {
		tr.sdk.Logger.WithName("[SUBSCRIBER]").Error(ackErr, errors.ErrMsgAck)
	}
}

func (tr *Runner) processRunnerError(msg *nats.Msg, errMsg, requestID string) {
	ackErr := msg.Ack()
	if ackErr != nil {
		tr.sdk.Logger.WithName("[SUBSCRIBER]").Error(ackErr, errors.ErrMsgAck)
	}

	tr.sdk.Logger.WithName("[SUBSCRIBER]").V(1).Info(errMsg)
	tr.publishError(requestID, errMsg)
}

func (tr *Runner) newRequestMessage(data []byte) (*kai.KaiNatsMessage, error) {
	requestMsg := &kai.KaiNatsMessage{}

	var err error
	if common.IsCompressed(data) {
		data, err = common.UncompressData(data)
		if err != nil {
			tr.sdk.Logger.WithName("[SUBSCRIBER]").Error(err, "Error reading compressed message")
			return nil, err
		}
	}

	err = proto.Unmarshal(data, requestMsg)

	return requestMsg, err
}

func (tr *Runner) publishError(requestID, errMsg string) {
	responseMsg := &kai.KaiNatsMessage{
		RequestId:   requestID,
		Error:       errMsg,
		FromNode:    viper.GetString("metadata.process_id"),
		MessageType: kai.MessageType_ERROR,
	}
	tr.publishResponse(responseMsg, "")
}

func (tr *Runner) publishResponse(responseMsg *kai.KaiNatsMessage, channel string) {
	outputSubject := tr.getOutputSubject(channel)

	outputMsg, err := proto.Marshal(responseMsg)
	if err != nil {
		tr.sdk.Logger.WithName("[SUBSCRIBER]").
			Error(err, "Error generating output result because handler result is not a serializable Protobuf")

		return
	}

	outputMsg, err = tr.prepareOutputMessage(outputMsg)
	if err != nil {
		tr.sdk.Logger.WithName("[SUBSCRIBER]").Error(err, "Error preparing output message")
		return
	}

	tr.sdk.Logger.WithName("[SUBSCRIBER]").V(1).Info("Publishing response",
		"Subject", outputSubject)

	_, err = tr.jetstream.Publish(outputSubject, outputMsg)
	if err != nil {
		tr.sdk.Logger.WithName("[SUBSCRIBER]").Error(err, "Error publishing output")
	}
}

func (tr *Runner) getOutputSubject(channel string) string {
	outputSubject := viper.GetString("nats.output")
	if channel != "" {
		return fmt.Sprintf("%s.%s", outputSubject, channel)
	}

	return outputSubject
}

// prepareOutputMessage will check the length of the message and compress it if necessary.
// Fails on compressed messages bigger than the threshold.
func (tr *Runner) prepareOutputMessage(msg []byte) ([]byte, error) {
	maxSize, err := tr.getMaxMessageSize()
	if err != nil {
		return nil, fmt.Errorf("error getting max message size: %s", err)
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
		tr.sdk.Logger.WithName("[SUBSCRIBER]").V(1).
			Info("Compressed message exceeds maximum size allowed",
				"Current message size", sizeInMB(lenOutMsg),
				"Compressed message size", sizeInMB(maxSize))

		return nil, errors.ErrMessageToBig
	}

	tr.sdk.Logger.WithName("[SUBSCRIBER]").Info("Message prepared",
		"Current message size", sizeInMB(lenOutMsg),
		"Compressed message size", sizeInMB(maxSize))

	return outMsg, nil
}

func (tr *Runner) getResponseHandler(subject string) Handler {
	if responseHandler, ok := tr.responseHandlers[subject]; ok {
		return responseHandler
	}

	// returns the default response handler, or nil if it doesn't exist
	return tr.responseHandlers["default"]
}

func (tr *Runner) getMaxMessageSize() (int64, error) {
	streamInfo, err := tr.jetstream.StreamInfo(viper.GetString("nats.stream"))
	if err != nil {
		return 0, fmt.Errorf("error getting stream's max message size: %w", err)
	}

	streamMaxSize := int64(streamInfo.Config.MaxMsgSize)
	serverMaxSize := tr.nats.MaxPayload()

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

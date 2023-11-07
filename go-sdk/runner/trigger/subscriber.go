package trigger

import (
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/go-logr/logr"

	"github.com/google/uuid"

	"github.com/nats-io/nats.go"
	"github.com/spf13/viper"
	"google.golang.org/protobuf/proto"

	"github.com/konstellation-io/kai-sdk/go-sdk/internal/common"
	"github.com/konstellation-io/kai-sdk/go-sdk/internal/errors"
	kai "github.com/konstellation-io/kai-sdk/go-sdk/protos"
	"github.com/konstellation-io/kai-sdk/go-sdk/sdk"
)

const _subscriberLoggerName = "[SUBSCRIBER]"

func (tr *Runner) getLoggerWithName() logr.Logger {
	return tr.sdk.Logger.WithName(_subscriberLoggerName)
}

func (tr *Runner) startSubscriber() {
	inputSubjects := viper.GetStringSlice(common.ConfigNatsInputsKey)
	subscriptions := make([]*nats.Subscription, 0, len(inputSubjects))

	for _, subject := range inputSubjects {
		consumerName := fmt.Sprintf("%s-%s", strings.ReplaceAll(subject, ".", "-"),
			strings.ReplaceAll(strings.ReplaceAll(tr.sdk.Metadata.GetProcess(), ".", "-"), " ", "-"))

		tr.getLoggerWithName().V(1).Info("Subscribing to subject",
			"Subject", subject, "Queue group", consumerName)

		s, err := tr.jetstream.Subscribe(
			subject,
			tr.processMessage,
			nats.DeliverNew(),
			nats.Durable(fmt.Sprintf("%s-%s", consumerName, uuid.New().String())),
			nats.ManualAck(),
			nats.AckWait(viper.GetDuration(common.ConfigRunnerSubscriberAckWaitTimeKey)),
		)
		if err != nil {
			tr.getLoggerWithName().Error(err, "Error subscribing to NATS subject",
				"Subject", subject)
			wg.Done()
			os.Exit(1)
		}

		subscriptions = append(subscriptions, s)

		tr.getLoggerWithName().V(1).Info("Listening to subject",
			"Subject", subject, "Queue group", consumerName)
	}

	// Handle sigterm and await termChan signal
	termChan := make(chan os.Signal, 1)
	signal.Notify(termChan, syscall.SIGINT, syscall.SIGTERM)
	<-termChan

	// Handle shutdown
	tr.getLoggerWithName().Info("Shutdown signal received")

	for _, s := range subscriptions {
		tr.getLoggerWithName().V(1).Info("Unsubscribing from subject",
			"Subject", s.Subject)

		err := s.Unsubscribe()
		if err != nil {
			tr.getLoggerWithName().Error(err, "Error unsubscribing from the NATS subject",
				"Subject", s.Subject)
			wg.Done()
			os.Exit(1)
		}
	}

	tr.getLoggerWithName().Info("Unsubscribed from all subjects")
	wg.Done()
}

func (tr *Runner) processMessage(msg *nats.Msg) {
	tr.getLoggerWithName().V(1).Info("New message received")

	requestMsg, err := tr.newRequestMessage(msg.Data)
	if err != nil {
		errMsg := fmt.Sprintf("Error parsing msg.data coming from subject %s because is not a valid protobuf: %s", msg.Subject, err)
		tr.processRunnerError(msg, errMsg, requestMsg.RequestId)

		return
	}

	tr.getLoggerWithName().Info("New message received",
		"Subject", msg.Subject, "Request ID", requestMsg.RequestId)

	if tr.responseHandler == nil {
		errMsg := "Error missing handler"
		tr.processRunnerError(msg, errMsg, requestMsg.RequestId)

		return
	}

	// Make a shallow copy of the sdk object to set inside the request msg.
	hSdk := sdk.ShallowCopyWithRequest(&tr.sdk, requestMsg)

	err = tr.responseHandler(hSdk, requestMsg.Payload)
	if err != nil {
		errMsg := fmt.Sprintf("Error in node %q executing handler for node %q: %s",
			tr.sdk.Metadata.GetProcess(), requestMsg.FromNode, err)
		tr.processRunnerError(msg, errMsg, requestMsg.RequestId)

		return
	}

	// Tell NATS we don't need to receive the message anymore, and we are done processing it.
	ackErr := msg.Ack()
	if ackErr != nil {
		tr.getLoggerWithName().Error(ackErr, errors.ErrMsgAck)
	}
}

func (tr *Runner) processRunnerError(msg *nats.Msg, errMsg, requestID string) {
	ackErr := msg.Ack()
	if ackErr != nil {
		tr.getLoggerWithName().Error(ackErr, errors.ErrMsgAck)
	}

	tr.getLoggerWithName().V(1).Info(errMsg)
	tr.publishError(requestID, errMsg)
}

func (tr *Runner) newRequestMessage(data []byte) (*kai.KaiNatsMessage, error) {
	requestMsg := &kai.KaiNatsMessage{}

	var err error
	if common.IsCompressed(data) {
		data, err = common.UncompressData(data)
		if err != nil {
			tr.getLoggerWithName().Error(err, "Error reading compressed message")
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
		FromNode:    viper.GetString(common.ConfigMetadataProcessIDKey),
		MessageType: kai.MessageType_ERROR,
	}
	tr.publishResponse(responseMsg, "")
}

func (tr *Runner) publishResponse(responseMsg *kai.KaiNatsMessage, channel string) {
	outputSubject := tr.getOutputSubject(channel)

	outputMsg, err := proto.Marshal(responseMsg)
	if err != nil {
		tr.getLoggerWithName().Error(err, "Error generating output result"+
			" because handler result is not a serializable Protobuf")
		return
	}

	outputMsg, err = tr.prepareOutputMessage(outputMsg)
	if err != nil {
		tr.getLoggerWithName().Error(err, "Error preparing output message")
		return
	}

	tr.getLoggerWithName().V(1).Info("Publishing response",
		"Subject", outputSubject)

	_, err = tr.jetstream.Publish(outputSubject, outputMsg)
	if err != nil {
		tr.getLoggerWithName().Error(err, "Error publishing output")
	}
}

func (tr *Runner) getOutputSubject(channel string) string {
	outputSubject := viper.GetString(common.ConfigNatsOutputKey)
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
		tr.getLoggerWithName().V(1).Info("Compressed message exceeds maximum size allowed",
			"Current message size", sizeInMB(lenOutMsg),
			"Compressed message size", sizeInMB(maxSize))

		return nil, errors.ErrMessageToBig
	}

	tr.getLoggerWithName().Info("Message prepared",
		"Current message size", sizeInMB(lenOutMsg),
		"Compressed message size", sizeInMB(maxSize))

	return outMsg, nil
}

func (tr *Runner) getMaxMessageSize() (int64, error) {
	streamInfo, err := tr.jetstream.StreamInfo(viper.GetString(common.ConfigNatsStreamKey))
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

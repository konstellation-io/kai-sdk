package messaging_test

import (
	"errors"
	"math/rand"
	"testing"

	"github.com/go-logr/logr"
	"github.com/go-logr/logr/testr"
	"github.com/nats-io/nats.go"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"

	"github.com/konstellation-io/kai-sdk/go-sdk/mocks"
	kai "github.com/konstellation-io/kai-sdk/go-sdk/protos"
	"github.com/konstellation-io/kai-sdk/go-sdk/sdk/messaging"
)

type SdkMessagingTestSuite struct {
	suite.Suite
	logger         logr.Logger
	jetstream      mocks.JetStreamContextMock
	messagingUtils mocks.MessagingUtilsMock
}

func (s *SdkMessagingTestSuite) SetupSuite() {
	s.logger = testr.NewWithOptions(s.T(), testr.Options{Verbosity: 1})
}

func (s *SdkMessagingTestSuite) SetupTest() {
	// Reset viper values before each test
	viper.Reset()

	s.jetstream = *mocks.NewJetStreamContextMock(s.T())
	s.messagingUtils = *mocks.NewMessagingUtilsMock(s.T())
}

func (s *SdkMessagingTestSuite) TestMessaging_InstantiateNewMessaging_ExpectOk() {
	// When
	messagingInst := messaging.NewMessaging(s.logger, nil, &s.jetstream, nil)

	// Then
	s.NotNil(messagingInst)
}

func (s *SdkMessagingTestSuite) TestMessaging_PublishError_ExpectOk() {
	// Given
	viper.SetDefault("nats.output", "test-parent")
	viper.SetDefault("metadata.process_id", "parent-node")
	s.jetstream.On("Publish", mock.AnythingOfType("string"), mock.AnythingOfType("[]uint8")).
		Return(&nats.PubAck{}, nil)
	s.messagingUtils.On("GetMaxMessageSize").Return(int64(2048), nil)

	messagingInst := messaging.NewTestMessaging(s.logger, nil, &s.jetstream, nil, &s.messagingUtils)

	// When
	messagingInst.SendError("some-request", "some-error")

	// Then
	s.NotNil(messagingInst)
	s.jetstream.AssertCalled(s.T(),
		"Publish", "test-parent",
		getOutputMessage("some-request", nil, "some-error", "parent-node", kai.MessageType_ERROR))
}

func (s *SdkMessagingTestSuite) TestMessaging_GetRequestID_ExpectOk() {
	// Given
	msg := &kai.KaiNatsMessage{
		RequestId:   "some-request",
		FromNode:    "some-node",
		MessageType: kai.MessageType_OK,
	}
	msgBytes, _ := proto.Marshal(msg)

	messagingInst := messaging.NewTestMessaging(s.logger, nil, &s.jetstream, nil, &s.messagingUtils)

	// When
	requestID, err := messagingInst.GetRequestID(msgBytes)

	// Then
	s.Nil(err)
	s.Equal(msg.RequestId, requestID)
}

func (s *SdkMessagingTestSuite) TestMessaging_GetRequestID_ExpectError() {
	// Given
	msgBytes := []byte("some-invalid-message")

	messagingInst := messaging.NewTestMessaging(s.logger, nil, &s.jetstream, nil, &s.messagingUtils)

	// When
	requestID, err := messagingInst.GetRequestID(msgBytes)

	// Then
	s.NotNil(err)
	s.Equal("", requestID)
}

func (s *SdkMessagingTestSuite) TestMessaging_GetRequestID_ErrorDuringUncompress() {
	// Given
	msgBytes := []byte("compressed-data")

	s.messagingUtils.On("IsCompressed", msgBytes).Return(true, nil)
	s.messagingUtils.On("UncompressData", msgBytes).Return(nil, errors.New("Uncompression error"))

	messagingInst := messaging.NewTestMessaging(s.logger, nil, &s.jetstream, nil, &s.messagingUtils)

	// When
	requestID, err := messagingInst.GetRequestID(msgBytes)

	// Then
	s.NotNil(err)
	s.Equal("", requestID)
}

func TestSdkMessagingTestSuite(t *testing.T) {
	suite.Run(t, new(SdkMessagingTestSuite))
}

func generateRandomString(sizeInBytes int) string {
	validChars := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	validCharCount := len(validChars)
	randomBytes := make([]byte, sizeInBytes)

	for i := 0; i < sizeInBytes; i++ {
		randomBytes[i] = validChars[rand.Intn(validCharCount)]
	}

	return string(randomBytes)
}

//nolint:unparam // false positive
func getOutputMessage(requestID string, msg interface{},
	errorMessage, fromNode string, messageType kai.MessageType) []byte {
	var payload *anypb.Any

	if msg != nil {
		value, ok := msg.(*anypb.Any)
		if ok {
			payload = value
		} else {
			payload, _ = anypb.New(msg.(proto.Message))
		}
	}

	responseMsg := &kai.KaiNatsMessage{
		RequestId:   requestID,
		Payload:     payload,
		FromNode:    fromNode,
		Error:       errorMessage,
		MessageType: messageType,
	}
	outputMsg, _ := proto.Marshal(responseMsg)

	return outputMsg
}

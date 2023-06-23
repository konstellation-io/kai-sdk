package messaging_test

import (
	"github.com/go-logr/logr"
	"github.com/go-logr/logr/testr"
	"github.com/konstellation-io/kre-runners/go-sdk/v1/mocks"
	kai "github.com/konstellation-io/kre-runners/go-sdk/v1/protos"
	"github.com/konstellation-io/kre-runners/go-sdk/v1/sdk/messaging"
	"github.com/nats-io/nats.go"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"
	"math/rand"
	"testing"
)

type SdkMessagingTestSuite struct {
	suite.Suite
	logger       logr.Logger
	nats         *nats.Conn
	jetstream    mocks.JetStreamContextMock
	messageUtils mocks.MessageUtilsMock
}

func (suite *SdkMessagingTestSuite) SetupSuite() {
	suite.logger = testr.NewWithOptions(suite.T(), testr.Options{Verbosity: 1})
}

func (suite *SdkMessagingTestSuite) SetupTest() {
	// Reset viper values before each test
	viper.Reset()

	suite.jetstream = *mocks.NewJetStreamContextMock(suite.T())
	suite.messageUtils = *mocks.NewMessageUtilsMock(suite.T())
}

func (suite *SdkMessagingTestSuite) TestMessaging_InstantiateNewMessaging_ExpectOk() {
	// When
	objectStore := messaging.NewMessaging(suite.logger, nil, &suite.jetstream, nil)

	// Then
	suite.NotNil(objectStore)
}

func (suite *SdkMessagingTestSuite) TestMessaging_PublishError_ExpectOk() {
	// Given
	viper.SetDefault("nats.output", "test-parent")
	viper.SetDefault("metadata.process_id", "parent-node")
	suite.jetstream.On("Publish", mock.AnythingOfType("string"), mock.AnythingOfType("[]uint8")).
		Return(&nats.PubAck{}, nil)
	suite.messageUtils.On("GetMaxMessageSize").Return(int64(2048), nil)
	objectStore := messaging.NewTestMessaging(suite.logger, nil, &suite.jetstream, nil, &suite.messageUtils)

	// When
	objectStore.SendError("some-request", "some-error")

	// Then
	suite.NotNil(objectStore)
	suite.jetstream.AssertCalled(suite.T(),
		"Publish", "test-parent",
		getOutputMessage("some-request", nil, "some-error", "parent-node", kai.MessageType_ERROR))
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

func getOutputMessage(requestId string, msg interface{}, error string, fromNode string, messageType kai.MessageType) []byte {
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
		RequestId:   requestId,
		Payload:     payload,
		FromNode:    fromNode,
		Error:       error,
		MessageType: messageType,
	}
	outputMsg, _ := proto.Marshal(responseMsg)
	return outputMsg
}

package messaging_test

import (
	"fmt"
	"github.com/golang/protobuf/ptypes/wrappers"
	kai "github.com/konstellation-io/kre-runners/go-sdk/v1/protos"
	"github.com/konstellation-io/kre-runners/go-sdk/v1/sdk/messaging"
	"github.com/nats-io/nats.go"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/mock"
	"google.golang.org/protobuf/types/known/anypb"
)

func (suite *SdkMessagingTestSuite) TestMessaging_SendAny_ExpectOk() {
	// Given
	viper.SetDefault("nats.output", "test-parent")
	viper.SetDefault("metadata.process_id", "parent-node")
	suite.jetstream.On("Publish", mock.AnythingOfType("string"), mock.AnythingOfType("[]uint8")).
		Return(&nats.PubAck{}, nil)
	suite.messageUtils.On("GetMaxMessageSize").Return(int64(1024*1024*1024), nil)
	objectStore := messaging.NewTestMessaging(suite.logger, nil, &suite.jetstream, &kai.KaiNatsMessage{}, &suite.messageUtils)

	// When
	msg, err := anypb.New(&wrappers.StringValue{
		Value: "Hi there!",
	})
	objectStore.SendAny(msg)

	// Then
	suite.NoError(err)
	suite.NotNil(objectStore)
	suite.messageUtils.AssertNumberOfCalls(suite.T(), "GetMaxMessageSize", 1)
	suite.jetstream.AssertCalled(suite.T(),
		"Publish", "test-parent", mock.AnythingOfType("[]uint8"))
}

func (suite *SdkMessagingTestSuite) TestMessaging_SendAnyWithExistingRequestMessage_ExpectOk() {
	// Given
	viper.SetDefault("nats.output", "test-parent")
	viper.SetDefault("metadata.process_id", "parent-node")
	suite.jetstream.On("Publish", mock.AnythingOfType("string"), mock.AnythingOfType("[]uint8")).
		Return(&nats.PubAck{}, nil)
	suite.messageUtils.On("GetMaxMessageSize").Return(int64(1024*1024*1024), nil)
	request := kai.KaiNatsMessage{RequestId: "123"}
	objectStore := messaging.NewTestMessaging(suite.logger, nil, &suite.jetstream, &request, &suite.messageUtils)

	// When
	msg, err := anypb.New(&wrappers.StringValue{
		Value: "Hi there!",
	})
	objectStore.SendAny(msg)

	// Then
	suite.NoError(err)
	suite.NotNil(objectStore)
	suite.messageUtils.AssertNumberOfCalls(suite.T(), "GetMaxMessageSize", 1)
	suite.jetstream.AssertCalled(suite.T(),
		"Publish", "test-parent",
		getOutputMessage("123", msg, "", "parent-node", kai.MessageType_OK))
}

func (suite *SdkMessagingTestSuite) TestMessaging_SendAnyWithCustomRequestId_ExpectOk() {
	// Given
	viper.SetDefault("nats.output", "test-parent")
	viper.SetDefault("metadata.process_id", "parent-node")
	suite.jetstream.On("Publish", mock.AnythingOfType("string"), mock.AnythingOfType("[]uint8")).
		Return(&nats.PubAck{}, nil)
	suite.messageUtils.On("GetMaxMessageSize").Return(int64(1024*1024*1024), nil)
	objectStore := messaging.NewTestMessaging(suite.logger, nil, &suite.jetstream, &kai.KaiNatsMessage{}, &suite.messageUtils)

	// When
	msg, err := anypb.New(&wrappers.StringValue{
		Value: "Hi there!",
	})
	objectStore.SendAnyWithRequestID(msg, "myRequestId")

	// Then
	suite.NoError(err)
	suite.NotNil(objectStore)
	suite.messageUtils.AssertNumberOfCalls(suite.T(), "GetMaxMessageSize", 1)
	suite.jetstream.AssertCalled(suite.T(),
		"Publish", "test-parent",
		getOutputMessage("myRequestId", msg, "", "parent-node", kai.MessageType_OK))
}

func (suite *SdkMessagingTestSuite) TestMessaging_SendAny_WithCompression_ExpectOk() {
	// Given
	viper.SetDefault("nats.output", "test-parent")
	viper.SetDefault("metadata.process_id", "parent-node")
	suite.jetstream.On("Publish", mock.AnythingOfType("string"), mock.AnythingOfType("[]uint8")).
		Return(&nats.PubAck{}, nil)
	suite.messageUtils.On("GetMaxMessageSize").Return(int64(2048), nil)
	objectStore := messaging.NewTestMessaging(suite.logger, nil, &suite.jetstream,
		&kai.KaiNatsMessage{}, &suite.messageUtils)

	// When
	msg, err := anypb.New(&wrappers.StringValue{
		Value: generateRandomString(2049),
	})
	objectStore.SendAny(msg)

	// Then
	suite.NoError(err)
	suite.NotNil(objectStore)
	suite.messageUtils.AssertNumberOfCalls(suite.T(), "GetMaxMessageSize", 1)
	suite.jetstream.AssertCalled(suite.T(),
		"Publish", "test-parent", mock.AnythingOfType("[]uint8"))
}

func (suite *SdkMessagingTestSuite) TestMessaging_SendAny_WithCompression_MessageToBig_ExpectError() {
	// Given
	viper.SetDefault("nats.output", "test-parent")
	viper.SetDefault("metadata.process_id", "parent-node")
	suite.jetstream.On("Publish", mock.AnythingOfType("string"), mock.AnythingOfType("[]uint8")).
		Return(&nats.PubAck{}, nil)
	suite.messageUtils.On("GetMaxMessageSize").Return(int64(128), nil)
	objectStore := messaging.NewTestMessaging(suite.logger, nil, &suite.jetstream,
		&kai.KaiNatsMessage{}, &suite.messageUtils)

	// When
	msg, err := anypb.New(&wrappers.StringValue{
		Value: generateRandomString(15000),
	})
	objectStore.SendAny(msg)

	// Then
	suite.NoError(err)
	suite.NotNil(objectStore)
	suite.messageUtils.AssertNumberOfCalls(suite.T(), "GetMaxMessageSize", 1)
	suite.jetstream.AssertNotCalled(suite.T(),
		"Publish", "test-parent")
}

func (suite *SdkMessagingTestSuite) TestMessaging_SendAnyToSubtopic_ExpectOk() {
	// Given
	viper.SetDefault("nats.output", "test-parent")
	viper.SetDefault("metadata.process_id", "parent-node")
	suite.jetstream.On("Publish", mock.AnythingOfType("string"), mock.AnythingOfType("[]uint8")).
		Return(&nats.PubAck{}, nil)
	suite.messageUtils.On("GetMaxMessageSize").Return(int64(1024*1024*1024), nil)
	objectStore := messaging.NewTestMessaging(suite.logger, nil, &suite.jetstream, &kai.KaiNatsMessage{}, &suite.messageUtils)

	// When
	msg, err := anypb.New(&wrappers.StringValue{
		Value: "Hi there!",
	})
	objectStore.SendAny(msg, "subtopic")

	// Then
	suite.NoError(err)
	suite.NotNil(objectStore)
	suite.messageUtils.AssertNumberOfCalls(suite.T(), "GetMaxMessageSize", 1)
	suite.jetstream.AssertCalled(suite.T(),
		"Publish", "test-parent.subtopic", mock.AnythingOfType("[]uint8"))
}

func (suite *SdkMessagingTestSuite) TestMessaging_SendAny_ErrorOnMaxMessageSize_ExpectError() {
	// Given
	viper.SetDefault("nats.output", "test-parent")
	viper.SetDefault("metadata.process_id", "parent-node")
	suite.jetstream.On("Publish", mock.AnythingOfType("string"), mock.AnythingOfType("[]uint8")).
		Return(&nats.PubAck{}, nil)
	suite.messageUtils.On("GetMaxMessageSize").Return(int64(0), fmt.Errorf("error getting size"))
	request := kai.KaiNatsMessage{RequestId: "123"}
	objectStore := messaging.NewTestMessaging(suite.logger, nil, &suite.jetstream, &request, &suite.messageUtils)

	// When
	msg, err := anypb.New(&wrappers.StringValue{
		Value: generateRandomString(1024),
	})
	objectStore.SendAny(msg)

	// Then
	suite.NoError(err)
	suite.NotNil(objectStore)
	suite.messageUtils.AssertNumberOfCalls(suite.T(), "GetMaxMessageSize", 1)
	suite.jetstream.AssertNotCalled(suite.T(), "Publish")
}

func (suite *SdkMessagingTestSuite) TestMessaging_SendAny_ErrorOnPublish_ExpectError() {
	// Given
	viper.SetDefault("nats.output", "test-parent")
	viper.SetDefault("metadata.process_id", "parent-node")
	suite.jetstream.On("Publish", mock.AnythingOfType("string"), mock.AnythingOfType("[]uint8")).
		Return(nil, fmt.Errorf("error publishing"))
	suite.messageUtils.On("GetMaxMessageSize").Return(int64(2048), nil)
	request := kai.KaiNatsMessage{RequestId: "123"}
	objectStore := messaging.NewTestMessaging(suite.logger, nil, &suite.jetstream, &request, &suite.messageUtils)

	// When
	msg, err := anypb.New(&wrappers.StringValue{
		Value: generateRandomString(1024),
	})
	objectStore.SendAny(msg)

	// Then
	suite.NoError(err)
	suite.NotNil(objectStore)
	suite.messageUtils.AssertNumberOfCalls(suite.T(), "GetMaxMessageSize", 1)
	suite.jetstream.AssertNotCalled(suite.T(), "Publish")
}

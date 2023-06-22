package messaging_test

import (
	"fmt"
	"github.com/golang/protobuf/ptypes/wrappers"
	kai "github.com/konstellation-io/kre-runners/go-sdk/v1/protos"
	"github.com/konstellation-io/kre-runners/go-sdk/v1/sdk/messaging"
	"github.com/nats-io/nats.go"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/mock"
)

func (suite *SdkMessagingTestSuite) TestMessaging_SendEarlyExit_ExpectOk() {
	// Given
	viper.SetDefault("nats.output", "test-parent")
	suite.jetstream.On("Publish", mock.AnythingOfType("string"), mock.AnythingOfType("[]uint8")).
		Return(&nats.PubAck{}, nil)
	suite.messageUtils.On("GetMaxMessageSize").Return(int64(1024*1024*1024), nil)
	objectStore := messaging.NewTestMessaging(suite.logger, nil, &suite.jetstream, &kai.KaiNatsMessage{}, &suite.messageUtils)

	// When
	msg := wrappers.StringValue{
		Value: "Hi there!",
	}
	err := objectStore.SendEarlyExit(&msg)

	// Then
	suite.NoError(err)
	suite.NotNil(objectStore)
	suite.messageUtils.AssertNumberOfCalls(suite.T(), "GetMaxMessageSize", 1)
	suite.jetstream.AssertCalled(suite.T(),
		"Publish", "test-parent", mock.AnythingOfType("[]uint8"))
}

func (suite *SdkMessagingTestSuite) TestMessaging_SendEarlyExitWithExistingRequestMessage_ExpectOk() {
	// Given
	viper.SetDefault("nats.output", "test-parent")
	viper.SetDefault("metadata.process_id", "parent-node")
	suite.jetstream.On("Publish", mock.AnythingOfType("string"), mock.AnythingOfType("[]uint8")).
		Return(&nats.PubAck{}, nil)
	suite.messageUtils.On("GetMaxMessageSize").Return(int64(1024*1024*1024), nil)
	request := kai.KaiNatsMessage{RequestId: "123"}
	objectStore := messaging.NewTestMessaging(suite.logger, nil, &suite.jetstream, &request, &suite.messageUtils)

	// When
	msg := wrappers.StringValue{
		Value: "Hi there!",
	}
	err := objectStore.SendEarlyExit(&msg)

	// Then
	suite.NoError(err)
	suite.NotNil(objectStore)
	suite.messageUtils.AssertNumberOfCalls(suite.T(), "GetMaxMessageSize", 1)
	suite.jetstream.AssertCalled(suite.T(),
		"Publish", "test-parent",
		getOutputMessage("123", &msg, "", "parent-node", kai.MessageType_EARLY_EXIT))
}

func (suite *SdkMessagingTestSuite) TestMessaging_SendEarlyExit_WithCompression_ExpectOk() {
	// Given
	viper.SetDefault("nats.output", "test-parent")
	viper.SetDefault("metadata.process_id", "parent-node")
	suite.jetstream.On("Publish", mock.AnythingOfType("string"), mock.AnythingOfType("[]uint8")).
		Return(&nats.PubAck{}, nil)
	suite.messageUtils.On("GetMaxMessageSize").Return(int64(2048), nil)
	objectStore := messaging.NewTestMessaging(suite.logger, nil, &suite.jetstream,
		&kai.KaiNatsMessage{}, &suite.messageUtils)

	// When
	msg := wrappers.StringValue{
		Value: generateRandomString(2049),
	}
	err := objectStore.SendEarlyExit(&msg)

	// Then
	suite.NoError(err)
	suite.NotNil(objectStore)
	suite.messageUtils.AssertNumberOfCalls(suite.T(), "GetMaxMessageSize", 1)
	suite.jetstream.AssertCalled(suite.T(),
		"Publish", "test-parent", mock.AnythingOfType("[]uint8"))
}

func (suite *SdkMessagingTestSuite) TestMessaging_SendEarlyExit_WithCompression_MessageToBig_ExpectError() {
	// Given
	viper.SetDefault("nats.output", "test-parent")
	viper.SetDefault("metadata.process_id", "parent-node")
	suite.jetstream.On("Publish", mock.AnythingOfType("string"), mock.AnythingOfType("[]uint8")).
		Return(&nats.PubAck{}, nil)
	suite.messageUtils.On("GetMaxMessageSize").Return(int64(128), nil)
	objectStore := messaging.NewTestMessaging(suite.logger, nil, &suite.jetstream,
		&kai.KaiNatsMessage{}, &suite.messageUtils)

	// When
	msg := wrappers.StringValue{
		Value: generateRandomString(15000),
	}
	err := objectStore.SendEarlyExit(&msg)

	// Then
	suite.NoError(err)
	suite.NotNil(objectStore)
	suite.messageUtils.AssertNumberOfCalls(suite.T(), "GetMaxMessageSize", 1)
	suite.jetstream.AssertNotCalled(suite.T(),
		"Publish", "test-parent")
}

func (suite *SdkMessagingTestSuite) TestMessaging_SendEarlyExitToSubtopic_ExpectOk() {
	// Given
	viper.SetDefault("nats.output", "test-parent")
	viper.SetDefault("metadata.process_id", "parent-node")
	suite.jetstream.On("Publish", mock.AnythingOfType("string"), mock.AnythingOfType("[]uint8")).
		Return(&nats.PubAck{}, nil)
	suite.messageUtils.On("GetMaxMessageSize").Return(int64(1024*1024*1024), nil)
	objectStore := messaging.NewTestMessaging(suite.logger, nil, &suite.jetstream, &kai.KaiNatsMessage{}, &suite.messageUtils)

	// When
	msg := wrappers.StringValue{
		Value: "Hi there!",
	}
	err := objectStore.SendEarlyExit(&msg, "subtopic")

	// Then
	suite.NoError(err)
	suite.NotNil(objectStore)
	suite.messageUtils.AssertNumberOfCalls(suite.T(), "GetMaxMessageSize", 1)
	suite.jetstream.AssertCalled(suite.T(),
		"Publish", "test-parent.subtopic", mock.AnythingOfType("[]uint8"))
}

func (suite *SdkMessagingTestSuite) TestMessaging_SendEarlyExit_ErrorOnMaxMessageSize_ExpectError() {
	// Given
	viper.SetDefault("nats.output", "test-parent")
	viper.SetDefault("metadata.process_id", "parent-node")
	suite.jetstream.On("Publish", mock.AnythingOfType("string"), mock.AnythingOfType("[]uint8")).
		Return(&nats.PubAck{}, nil)
	suite.messageUtils.On("GetMaxMessageSize").Return(int64(0), fmt.Errorf("error getting size"))
	request := kai.KaiNatsMessage{RequestId: "123"}
	objectStore := messaging.NewTestMessaging(suite.logger, nil, &suite.jetstream, &request, &suite.messageUtils)

	// When
	msg := wrappers.StringValue{
		Value: generateRandomString(1024),
	}
	err := objectStore.SendEarlyExit(&msg)

	// Then
	suite.NoError(err)
	suite.NotNil(objectStore)
	suite.messageUtils.AssertNumberOfCalls(suite.T(), "GetMaxMessageSize", 1)
	suite.jetstream.AssertNotCalled(suite.T(), "Publish")
}

func (suite *SdkMessagingTestSuite) TestMessaging_SendEarlyExit_ErrorOnPublish_ExpectError() {
	// Given
	viper.SetDefault("nats.output", "test-parent")
	viper.SetDefault("metadata.process_id", "parent-node")
	suite.jetstream.On("Publish", mock.AnythingOfType("string"), mock.AnythingOfType("[]uint8")).
		Return(nil, fmt.Errorf("error publishing"))
	suite.messageUtils.On("GetMaxMessageSize").Return(int64(2048), nil)
	request := kai.KaiNatsMessage{RequestId: "123"}
	objectStore := messaging.NewTestMessaging(suite.logger, nil, &suite.jetstream, &request, &suite.messageUtils)

	// When
	msg := wrappers.StringValue{
		Value: generateRandomString(1024),
	}
	err := objectStore.SendEarlyExit(&msg)

	// Then
	suite.NoError(err)
	suite.NotNil(objectStore)
	suite.messageUtils.AssertNumberOfCalls(suite.T(), "GetMaxMessageSize", 1)
	suite.jetstream.AssertNotCalled(suite.T(), "Publish")
}

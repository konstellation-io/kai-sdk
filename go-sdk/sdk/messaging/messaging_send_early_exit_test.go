package messaging_test

import (
	"fmt"

	"github.com/golang/protobuf/ptypes/wrappers"
	"github.com/nats-io/nats.go"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/mock"

	kai "github.com/konstellation-io/kai-sdk/go-sdk/protos"
	"github.com/konstellation-io/kai-sdk/go-sdk/sdk/messaging"
)

func (s *SdkMessagingTestSuite) TestMessaging_SendEarlyExit_ExpectOk() {
	// Given
	viper.SetDefault(natsOutputField, natsOutputValue)
	s.jetstream.On("Publish", mock.AnythingOfType("string"), mock.AnythingOfType(unit8Type)).
		Return(&nats.PubAck{}, nil)
	s.messagingUtils.On("GetMaxMessageSize").Return(int64(1024*1024*1024), nil)
	
	objectStore := messaging.NewTestMessaging(s.logger, nil, &s.jetstream, &kai.KaiNatsMessage{}, &s.messagingUtils)

	// When
	msg := wrappers.StringValue{
		Value: stringValueMessage,
	}
	err := objectStore.SendEarlyExit(&msg)

	// Then
	s.NoError(err)
	s.NotNil(objectStore)
	s.messagingUtils.AssertNumberOfCalls(s.T(), "GetMaxMessageSize", 1)
	s.jetstream.AssertCalled(s.T(),
		"Publish", natsOutputValue, mock.AnythingOfType(unit8Type))
}

func (s *SdkMessagingTestSuite) TestMessaging_SendEarlyExitWithExistingRequestMessage_ExpectOk() {
	// Given
	viper.SetDefault(natsOutputField, natsOutputValue)
	viper.SetDefault(metadataProcessIDField, metadataProcessIDValue)
	s.jetstream.On("Publish", mock.AnythingOfType("string"), mock.AnythingOfType(unit8Type)).
		Return(&nats.PubAck{}, nil)
	s.messagingUtils.On("GetMaxMessageSize").Return(int64(1024*1024*1024), nil)
	
	request := kai.KaiNatsMessage{RequestId: "123"}
	objectStore := messaging.NewTestMessaging(s.logger, nil, &s.jetstream, &request, &s.messagingUtils)

	// When
	msg := wrappers.StringValue{
		Value: stringValueMessage,
	}
	err := objectStore.SendEarlyExit(&msg)

	// Then
	s.NoError(err)
	s.NotNil(objectStore)
	s.messagingUtils.AssertNumberOfCalls(s.T(), "GetMaxMessageSize", 1)
	s.jetstream.AssertCalled(s.T(),
		"Publish", natsOutputValue,
		getOutputMessage("123", &msg, "", metadataProcessIDValue, kai.MessageType_EARLY_EXIT))
}

func (s *SdkMessagingTestSuite) TestMessaging_SendEarlyExit_WithCompression_ExpectOk() {
	// Given
	viper.SetDefault(natsOutputField, natsOutputValue)
	viper.SetDefault(metadataProcessIDField, metadataProcessIDValue)
	s.jetstream.On("Publish", mock.AnythingOfType("string"), mock.AnythingOfType(unit8Type)).
		Return(&nats.PubAck{}, nil)
	s.messagingUtils.On("GetMaxMessageSize").Return(int64(2048), nil)
	
	objectStore := messaging.NewTestMessaging(s.logger, nil, &s.jetstream,
		&kai.KaiNatsMessage{}, &s.messagingUtils)

	// When
	msg := wrappers.StringValue{
		Value: generateRandomString(2049),
	}
	err := objectStore.SendEarlyExit(&msg)

	// Then
	s.NoError(err)
	s.NotNil(objectStore)
	s.messagingUtils.AssertNumberOfCalls(s.T(), "GetMaxMessageSize", 1)
	s.jetstream.AssertCalled(s.T(),
		"Publish", natsOutputValue, mock.AnythingOfType(unit8Type))
}

func (s *SdkMessagingTestSuite) TestMessaging_SendEarlyExit_WithCompression_MessageToBig_ExpectError() {
	// Given
	viper.SetDefault(natsOutputField, natsOutputValue)
	viper.SetDefault(metadataProcessIDField, metadataProcessIDValue)
	s.jetstream.On("Publish", mock.AnythingOfType("string"), mock.AnythingOfType(unit8Type)).
		Return(&nats.PubAck{}, nil)
	s.messagingUtils.On("GetMaxMessageSize").Return(int64(128), nil)
	
	objectStore := messaging.NewTestMessaging(s.logger, nil, &s.jetstream,
		&kai.KaiNatsMessage{}, &s.messagingUtils)

	// When
	msg := wrappers.StringValue{
		Value: generateRandomString(15000),
	}
	err := objectStore.SendEarlyExit(&msg)

	// Then
	s.NoError(err)
	s.NotNil(objectStore)
	s.messagingUtils.AssertNumberOfCalls(s.T(), "GetMaxMessageSize", 1)
	s.jetstream.AssertNotCalled(s.T(),
		"Publish", natsOutputValue)
}

func (s *SdkMessagingTestSuite) TestMessaging_SendEarlyExitToSubtopic_ExpectOk() {
	// Given
	viper.SetDefault(natsOutputField, natsOutputValue)
	viper.SetDefault(metadataProcessIDField, metadataProcessIDValue)
	s.jetstream.On("Publish", mock.AnythingOfType("string"), mock.AnythingOfType(unit8Type)).
		Return(&nats.PubAck{}, nil)
	s.messagingUtils.On("GetMaxMessageSize").Return(int64(1024*1024*1024), nil)

	objectStore := messaging.NewTestMessaging(s.logger, nil, &s.jetstream, &kai.KaiNatsMessage{}, &s.messagingUtils)

	// When
	msg := wrappers.StringValue{
		Value: stringValueMessage,
	}
	err := objectStore.SendEarlyExit(&msg, "subtopic")

	// Then
	s.NoError(err)
	s.NotNil(objectStore)
	s.messagingUtils.AssertNumberOfCalls(s.T(), "GetMaxMessageSize", 1)
	s.jetstream.AssertCalled(s.T(),
		"Publish", "test-parent.subtopic", mock.AnythingOfType(unit8Type))
}

func (s *SdkMessagingTestSuite) TestMessaging_SendEarlyExit_ErrorOnMaxMessageSize_ExpectError() {
	// Given
	viper.SetDefault(natsOutputField, natsOutputValue)
	viper.SetDefault(metadataProcessIDField, metadataProcessIDValue)
	s.jetstream.On("Publish", mock.AnythingOfType("string"), mock.AnythingOfType(unit8Type)).
		Return(&nats.PubAck{}, nil)
	s.messagingUtils.On("GetMaxMessageSize").Return(int64(0), fmt.Errorf("error getting size"))
	
	request := kai.KaiNatsMessage{RequestId: "123"}
	objectStore := messaging.NewTestMessaging(s.logger, nil, &s.jetstream, &request, &s.messagingUtils)

	// When
	msg := wrappers.StringValue{
		Value: generateRandomString(1024),
	}
	err := objectStore.SendEarlyExit(&msg)

	// Then
	s.NoError(err)
	s.NotNil(objectStore)
	s.messagingUtils.AssertNumberOfCalls(s.T(), "GetMaxMessageSize", 1)
	s.jetstream.AssertNotCalled(s.T(), "Publish")
}

func (s *SdkMessagingTestSuite) TestMessaging_SendEarlyExit_ErrorOnPublish_ExpectError() {
	// Given
	viper.SetDefault(natsOutputField, natsOutputValue)
	viper.SetDefault(metadataProcessIDField, metadataProcessIDValue)
	s.jetstream.On("Publish", mock.AnythingOfType("string"), mock.AnythingOfType(unit8Type)).
		Return(nil, fmt.Errorf("error publishing"))
	s.messagingUtils.On("GetMaxMessageSize").Return(int64(2048), nil)
	
	request := kai.KaiNatsMessage{RequestId: "123"}
	objectStore := messaging.NewTestMessaging(s.logger, nil, &s.jetstream, &request, &s.messagingUtils)

	// When
	msg := wrappers.StringValue{
		Value: generateRandomString(1024),
	}
	err := objectStore.SendEarlyExit(&msg)

	// Then
	s.NoError(err)
	s.NotNil(objectStore)
	s.messagingUtils.AssertNumberOfCalls(s.T(), "GetMaxMessageSize", 1)
	s.jetstream.AssertNotCalled(s.T(), "Publish")
}

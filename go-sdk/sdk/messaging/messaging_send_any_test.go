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

func (s *SdkMessagingTestSuite) TestMessaging_SendAny_ExpectOk() {
	// Given
	viper.SetDefault("nats.output", "test-parent")
	viper.SetDefault("metadata.process_id", "parent-node")
	s.jetstream.On("Publish", mock.AnythingOfType("string"), mock.AnythingOfType("[]uint8")).
		Return(&nats.PubAck{}, nil)
	s.messageUtils.On("GetMaxMessageSize").Return(int64(1024*1024*1024), nil)
	objectStore := messaging.NewTestMessaging(s.logger, nil, &s.jetstream, &kai.KaiNatsMessage{}, &s.messageUtils)

	// When
	msg, err := anypb.New(&wrappers.StringValue{
		Value: "Hi there!",
	})
	objectStore.SendAny(msg)

	// Then
	s.NoError(err)
	s.NotNil(objectStore)
	s.messageUtils.AssertNumberOfCalls(s.T(), "GetMaxMessageSize", 1)
	s.jetstream.AssertCalled(s.T(),
		"Publish", "test-parent", mock.AnythingOfType("[]uint8"))
}

func (s *SdkMessagingTestSuite) TestMessaging_SendAnyWithExistingRequestMessage_ExpectOk() {
	// Given
	viper.SetDefault("nats.output", "test-parent")
	viper.SetDefault("metadata.process_id", "parent-node")
	s.jetstream.On("Publish", mock.AnythingOfType("string"), mock.AnythingOfType("[]uint8")).
		Return(&nats.PubAck{}, nil)
	s.messageUtils.On("GetMaxMessageSize").Return(int64(1024*1024*1024), nil)
	request := kai.KaiNatsMessage{RequestId: "123"}
	objectStore := messaging.NewTestMessaging(s.logger, nil, &s.jetstream, &request, &s.messageUtils)

	// When
	msg, err := anypb.New(&wrappers.StringValue{
		Value: "Hi there!",
	})
	objectStore.SendAny(msg)

	// Then
	s.NoError(err)
	s.NotNil(objectStore)
	s.messageUtils.AssertNumberOfCalls(s.T(), "GetMaxMessageSize", 1)
	s.jetstream.AssertCalled(s.T(),
		"Publish", "test-parent",
		getOutputMessage("123", msg, "", "parent-node", kai.MessageType_OK))
}

func (s *SdkMessagingTestSuite) TestMessaging_SendAnyWithCustomRequestId_ExpectOk() {
	// Given
	viper.SetDefault("nats.output", "test-parent")
	viper.SetDefault("metadata.process_id", "parent-node")
	s.jetstream.On("Publish", mock.AnythingOfType("string"), mock.AnythingOfType("[]uint8")).
		Return(&nats.PubAck{}, nil)
	s.messageUtils.On("GetMaxMessageSize").Return(int64(1024*1024*1024), nil)
	objectStore := messaging.NewTestMessaging(s.logger, nil, &s.jetstream, &kai.KaiNatsMessage{}, &s.messageUtils)

	// When
	msg, err := anypb.New(&wrappers.StringValue{
		Value: "Hi there!",
	})
	objectStore.SendAnyWithRequestID(msg, "myRequestId")

	// Then
	s.NoError(err)
	s.NotNil(objectStore)
	s.messageUtils.AssertNumberOfCalls(s.T(), "GetMaxMessageSize", 1)
	s.jetstream.AssertCalled(s.T(),
		"Publish", "test-parent",
		getOutputMessage("myRequestId", msg, "", "parent-node", kai.MessageType_OK))
}

func (s *SdkMessagingTestSuite) TestMessaging_SendAny_WithCompression_ExpectOk() {
	// Given
	viper.SetDefault("nats.output", "test-parent")
	viper.SetDefault("metadata.process_id", "parent-node")
	s.jetstream.On("Publish", mock.AnythingOfType("string"), mock.AnythingOfType("[]uint8")).
		Return(&nats.PubAck{}, nil)
	s.messageUtils.On("GetMaxMessageSize").Return(int64(2048), nil)
	objectStore := messaging.NewTestMessaging(s.logger, nil, &s.jetstream,
		&kai.KaiNatsMessage{}, &s.messageUtils)

	// When
	msg, err := anypb.New(&wrappers.StringValue{
		Value: generateRandomString(2049),
	})
	objectStore.SendAny(msg)

	// Then
	s.NoError(err)
	s.NotNil(objectStore)
	s.messageUtils.AssertNumberOfCalls(s.T(), "GetMaxMessageSize", 1)
	s.jetstream.AssertCalled(s.T(),
		"Publish", "test-parent", mock.AnythingOfType("[]uint8"))
}

func (s *SdkMessagingTestSuite) TestMessaging_SendAny_WithCompression_MessageToBig_ExpectError() {
	// Given
	viper.SetDefault("nats.output", "test-parent")
	viper.SetDefault("metadata.process_id", "parent-node")
	s.jetstream.On("Publish", mock.AnythingOfType("string"), mock.AnythingOfType("[]uint8")).
		Return(&nats.PubAck{}, nil)
	s.messageUtils.On("GetMaxMessageSize").Return(int64(128), nil)
	objectStore := messaging.NewTestMessaging(s.logger, nil, &s.jetstream,
		&kai.KaiNatsMessage{}, &s.messageUtils)

	// When
	msg, err := anypb.New(&wrappers.StringValue{
		Value: generateRandomString(15000),
	})
	objectStore.SendAny(msg)

	// Then
	s.NoError(err)
	s.NotNil(objectStore)
	s.messageUtils.AssertNumberOfCalls(s.T(), "GetMaxMessageSize", 1)
	s.jetstream.AssertNotCalled(s.T(),
		"Publish", "test-parent")
}

func (s *SdkMessagingTestSuite) TestMessaging_SendAnyToSubtopic_ExpectOk() {
	// Given
	viper.SetDefault("nats.output", "test-parent")
	viper.SetDefault("metadata.process_id", "parent-node")
	s.jetstream.On("Publish", mock.AnythingOfType("string"), mock.AnythingOfType("[]uint8")).
		Return(&nats.PubAck{}, nil)
	s.messageUtils.On("GetMaxMessageSize").Return(int64(1024*1024*1024), nil)
	objectStore := messaging.NewTestMessaging(s.logger, nil, &s.jetstream, &kai.KaiNatsMessage{}, &s.messageUtils)

	// When
	msg, err := anypb.New(&wrappers.StringValue{
		Value: "Hi there!",
	})
	objectStore.SendAny(msg, "subtopic")

	// Then
	s.NoError(err)
	s.NotNil(objectStore)
	s.messageUtils.AssertNumberOfCalls(s.T(), "GetMaxMessageSize", 1)
	s.jetstream.AssertCalled(s.T(),
		"Publish", "test-parent.subtopic", mock.AnythingOfType("[]uint8"))
}

func (s *SdkMessagingTestSuite) TestMessaging_SendAny_ErrorOnMaxMessageSize_ExpectError() {
	// Given
	viper.SetDefault("nats.output", "test-parent")
	viper.SetDefault("metadata.process_id", "parent-node")
	s.jetstream.On("Publish", mock.AnythingOfType("string"), mock.AnythingOfType("[]uint8")).
		Return(&nats.PubAck{}, nil)
	s.messageUtils.On("GetMaxMessageSize").Return(int64(0), fmt.Errorf("error getting size"))
	request := kai.KaiNatsMessage{RequestId: "123"}
	objectStore := messaging.NewTestMessaging(s.logger, nil, &s.jetstream, &request, &s.messageUtils)

	// When
	msg, err := anypb.New(&wrappers.StringValue{
		Value: generateRandomString(1024),
	})
	objectStore.SendAny(msg)

	// Then
	s.NoError(err)
	s.NotNil(objectStore)
	s.messageUtils.AssertNumberOfCalls(s.T(), "GetMaxMessageSize", 1)
	s.jetstream.AssertNotCalled(s.T(), "Publish")
}

func (s *SdkMessagingTestSuite) TestMessaging_SendAny_ErrorOnPublish_ExpectError() {
	// Given
	viper.SetDefault("nats.output", "test-parent")
	viper.SetDefault("metadata.process_id", "parent-node")
	s.jetstream.On("Publish", mock.AnythingOfType("string"), mock.AnythingOfType("[]uint8")).
		Return(nil, fmt.Errorf("error publishing"))
	s.messageUtils.On("GetMaxMessageSize").Return(int64(2048), nil)
	request := kai.KaiNatsMessage{RequestId: "123"}
	objectStore := messaging.NewTestMessaging(s.logger, nil, &s.jetstream, &request, &s.messageUtils)

	// When
	msg, err := anypb.New(&wrappers.StringValue{
		Value: generateRandomString(1024),
	})
	objectStore.SendAny(msg)

	// Then
	s.NoError(err)
	s.NotNil(objectStore)
	s.messageUtils.AssertNumberOfCalls(s.T(), "GetMaxMessageSize", 1)
	s.jetstream.AssertNotCalled(s.T(), "Publish")
}

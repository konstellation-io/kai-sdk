package messaging_test

import (
	kai "github.com/konstellation-io/kai-sdk/go-sdk/protos"
	"github.com/konstellation-io/kai-sdk/go-sdk/sdk/messaging"
)

const (
	requestIDValue = "some-request-id"
	errorMessage   = "Some error message"
)

func (s *SdkMessagingTestSuite) TestMessaging_GetErrorMessage_ExpectOk() {
	// Given
	kaiMessage := &kai.KaiNatsMessage{
		RequestId:   requestIDValue,
		MessageType: kai.MessageType_ERROR,
		Error:       "Error message",
	}
	objectStore := messaging.NewTestMessaging(s.logger, nil, &s.jetstream, kaiMessage, &s.messagingUtils)

	// When
	errorMessage := objectStore.GetErrorMessage()

	// Then
	s.NotNil(objectStore)
	s.Equal("Error message", errorMessage)
}

func (s *SdkMessagingTestSuite) TestMessaging_GetErrorMessage_NoErrorMessageExistWhenTypeOK_ExpectError() {
	// Given
	kaiMessage := &kai.KaiNatsMessage{
		RequestId:   requestIDValue,
		MessageType: kai.MessageType_OK,
		Error:       errorMessage,
	}
	objectStore := messaging.NewTestMessaging(s.logger, nil, &s.jetstream, kaiMessage, &s.messagingUtils)

	// When
	errorMessage := objectStore.GetErrorMessage()

	// Then
	s.NotNil(objectStore)
	s.Empty(errorMessage)
}

func (s *SdkMessagingTestSuite) TestMessaging_IsMessageOk_MessageOk_ExpectTrue() {
	// Given
	kaiMessage := &kai.KaiNatsMessage{
		RequestId:   requestIDValue,
		MessageType: kai.MessageType_OK,
	}
	objectStore := messaging.NewTestMessaging(s.logger, nil, &s.jetstream, kaiMessage, &s.messagingUtils)

	// When
	ok := objectStore.IsMessageOK()

	// Then
	s.NotNil(objectStore)
	s.True(ok)
}

func (s *SdkMessagingTestSuite) TestMessaging_IsMessageOk_MessageNotOk_ExpectFalse() {
	// Given
	kaiMessage := &kai.KaiNatsMessage{
		RequestId:   requestIDValue,
		MessageType: kai.MessageType_ERROR,
		Error:       errorMessage,
	}
	objectStore := messaging.NewTestMessaging(s.logger, nil, &s.jetstream, kaiMessage, &s.messagingUtils)

	// When
	ok := objectStore.IsMessageOK()

	// Then
	s.NotNil(objectStore)
	s.False(ok)
}

func (s *SdkMessagingTestSuite) TestMessaging_IsMessageError_MessageError_ExpectTrue() {
	// Given
	kaiMessage := &kai.KaiNatsMessage{
		RequestId:   requestIDValue,
		MessageType: kai.MessageType_ERROR,
		Error:       errorMessage,
	}
	objectStore := messaging.NewTestMessaging(s.logger, nil, &s.jetstream, kaiMessage, &s.messagingUtils)

	// When
	isError := objectStore.IsMessageError()

	// Then
	s.NotNil(objectStore)
	s.True(isError)
}

func (s *SdkMessagingTestSuite) TestMessaging_IsMessageError_MessageNotError_ExpectFalse() {
	// Given
	kaiMessage := &kai.KaiNatsMessage{
		RequestId:   requestIDValue,
		MessageType: kai.MessageType_OK,
	}
	objectStore := messaging.NewTestMessaging(s.logger, nil, &s.jetstream, kaiMessage, &s.messagingUtils)

	// When
	isError := objectStore.IsMessageError()

	// Then
	s.NotNil(objectStore)
	s.False(isError)
}

func (s *SdkMessagingTestSuite) TestMessaging_IsMessageEarlyReply_MessageEarlyReply_ExpectTrue() {
	// Given
	kaiMessage := &kai.KaiNatsMessage{
		RequestId:   requestIDValue,
		MessageType: kai.MessageType_EARLY_REPLY,
	}
	objectStore := messaging.NewTestMessaging(s.logger, nil, &s.jetstream, kaiMessage, &s.messagingUtils)

	// When
	isEarlyReply := objectStore.IsMessageEarlyReply()

	// Then
	s.NotNil(objectStore)
	s.True(isEarlyReply)
}

func (s *SdkMessagingTestSuite) TestMessaging_IsMessageEarlyReply_MessageNotEarlyReply_ExpectFalse() {
	// Given
	kaiMessage := &kai.KaiNatsMessage{
		RequestId:   requestIDValue,
		MessageType: kai.MessageType_OK,
	}
	objectStore := messaging.NewTestMessaging(s.logger, nil, &s.jetstream, kaiMessage, &s.messagingUtils)

	// When
	isEarlyReply := objectStore.IsMessageEarlyReply()

	// Then
	s.NotNil(objectStore)
	s.False(isEarlyReply)
}

func (s *SdkMessagingTestSuite) TestMessaging_IsMessageEarlyExit_MessageEarlyExit_ExpectTrue() {
	// Given
	kaiMessage := &kai.KaiNatsMessage{
		RequestId:   requestIDValue,
		MessageType: kai.MessageType_EARLY_EXIT,
	}
	objectStore := messaging.NewTestMessaging(s.logger, nil, &s.jetstream, kaiMessage, &s.messagingUtils)

	// When
	isEarlyExit := objectStore.IsMessageEarlyExit()

	// Then
	s.NotNil(objectStore)
	s.True(isEarlyExit)
}

func (s *SdkMessagingTestSuite) TestMessaging_IsMessageEarlyExit_MessageNotEarlyExit_ExpectFalse() {
	// Given
	kaiMessage := &kai.KaiNatsMessage{
		RequestId:   requestIDValue,
		MessageType: kai.MessageType_OK,
	}
	objectStore := messaging.NewTestMessaging(s.logger, nil, &s.jetstream, kaiMessage, &s.messagingUtils)

	// When
	isEarlyExit := objectStore.IsMessageEarlyExit()

	// Then
	s.NotNil(objectStore)
	s.False(isEarlyExit)
}

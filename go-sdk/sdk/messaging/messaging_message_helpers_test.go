package messaging_test

import (
	kai "github.com/konstellation-io/kre-runners/go-sdk/v1/protos"
	"github.com/konstellation-io/kre-runners/go-sdk/v1/sdk/messaging"
)

func (s *SdkMessagingTestSuite) TestMessaging_GetErrorMessage_ExpectOk() {
	// Given
	kaiMessage := &kai.KaiNatsMessage{
		RequestId:   "some-request-id",
		MessageType: kai.MessageType_ERROR,
		Error:       "Error message",
	}
	objectStore := messaging.NewTestMessaging(s.logger, nil, &s.jetstream, kaiMessage, &s.messageUtils)

	// When
	errorMessage := objectStore.GetErrorMessage()

	// Then
	s.NotNil(objectStore)
	s.Equal("Error message", errorMessage)
}

func (s *SdkMessagingTestSuite) TestMessaging_GetErrorMessage_NoErrorMessageExistWhenTypeOK_ExpectError() {
	// Given
	kaiMessage := &kai.KaiNatsMessage{
		RequestId:   "some-request-id",
		MessageType: kai.MessageType_OK,
		Error:       "Some error message",
	}
	objectStore := messaging.NewTestMessaging(s.logger, nil, &s.jetstream, kaiMessage, &s.messageUtils)

	// When
	errorMessage := objectStore.GetErrorMessage()

	// Then
	s.NotNil(objectStore)
	s.Empty(errorMessage)
}

func (s *SdkMessagingTestSuite) TestMessaging_IsMessageOk_MessageOk_ExpectTrue() {
	// Given
	kaiMessage := &kai.KaiNatsMessage{
		RequestId:   "some-request-id",
		MessageType: kai.MessageType_OK,
	}
	objectStore := messaging.NewTestMessaging(s.logger, nil, &s.jetstream, kaiMessage, &s.messageUtils)

	// When
	ok := objectStore.IsMessageOK()

	// Then
	s.NotNil(objectStore)
	s.True(ok)
}

func (s *SdkMessagingTestSuite) TestMessaging_IsMessageOk_MessageNotOk_ExpectFalse() {
	// Given
	kaiMessage := &kai.KaiNatsMessage{
		RequestId:   "some-request-id",
		MessageType: kai.MessageType_ERROR,
		Error:       "Some error message",
	}
	objectStore := messaging.NewTestMessaging(s.logger, nil, &s.jetstream, kaiMessage, &s.messageUtils)

	// When
	ok := objectStore.IsMessageOK()

	// Then
	s.NotNil(objectStore)
	s.False(ok)
}

func (s *SdkMessagingTestSuite) TestMessaging_IsMessageError_MessageError_ExpectTrue() {
	// Given
	kaiMessage := &kai.KaiNatsMessage{
		RequestId:   "some-request-id",
		MessageType: kai.MessageType_ERROR,
		Error:       "Some error message",
	}
	objectStore := messaging.NewTestMessaging(s.logger, nil, &s.jetstream, kaiMessage, &s.messageUtils)

	// When
	isError := objectStore.IsMessageError()

	// Then
	s.NotNil(objectStore)
	s.True(isError)
}

func (s *SdkMessagingTestSuite) TestMessaging_IsMessageError_MessageNotError_ExpectFalse() {
	// Given
	kaiMessage := &kai.KaiNatsMessage{
		RequestId:   "some-request-id",
		MessageType: kai.MessageType_OK,
	}
	objectStore := messaging.NewTestMessaging(s.logger, nil, &s.jetstream, kaiMessage, &s.messageUtils)

	// When
	isError := objectStore.IsMessageError()

	// Then
	s.NotNil(objectStore)
	s.False(isError)
}

func (s *SdkMessagingTestSuite) TestMessaging_IsMessageEarlyReply_MessageEarlyReply_ExpectTrue() {
	// Given
	kaiMessage := &kai.KaiNatsMessage{
		RequestId:   "some-request-id",
		MessageType: kai.MessageType_EARLY_REPLY,
	}
	objectStore := messaging.NewTestMessaging(s.logger, nil, &s.jetstream, kaiMessage, &s.messageUtils)

	// When
	isEarlyReply := objectStore.IsMessageEarlyReply()

	// Then
	s.NotNil(objectStore)
	s.True(isEarlyReply)
}

func (s *SdkMessagingTestSuite) TestMessaging_IsMessageEarlyReply_MessageNotEarlyReply_ExpectFalse() {
	// Given
	kaiMessage := &kai.KaiNatsMessage{
		RequestId:   "some-request-id",
		MessageType: kai.MessageType_OK,
	}
	objectStore := messaging.NewTestMessaging(s.logger, nil, &s.jetstream, kaiMessage, &s.messageUtils)

	// When
	isEarlyReply := objectStore.IsMessageEarlyReply()

	// Then
	s.NotNil(objectStore)
	s.False(isEarlyReply)
}

func (s *SdkMessagingTestSuite) TestMessaging_IsMessageEarlyExit_MessageEarlyExit_ExpectTrue() {
	// Given
	kaiMessage := &kai.KaiNatsMessage{
		RequestId:   "some-request-id",
		MessageType: kai.MessageType_EARLY_EXIT,
	}
	objectStore := messaging.NewTestMessaging(s.logger, nil, &s.jetstream, kaiMessage, &s.messageUtils)

	// When
	isEarlyExit := objectStore.IsMessageEarlyExit()

	// Then
	s.NotNil(objectStore)
	s.True(isEarlyExit)
}

func (s *SdkMessagingTestSuite) TestMessaging_IsMessageEarlyExit_MessageNotEarlyExit_ExpectFalse() {
	// Given
	kaiMessage := &kai.KaiNatsMessage{
		RequestId:   "some-request-id",
		MessageType: kai.MessageType_OK,
	}
	objectStore := messaging.NewTestMessaging(s.logger, nil, &s.jetstream, kaiMessage, &s.messageUtils)

	// When
	isEarlyExit := objectStore.IsMessageEarlyExit()

	// Then
	s.NotNil(objectStore)
	s.False(isEarlyExit)
}

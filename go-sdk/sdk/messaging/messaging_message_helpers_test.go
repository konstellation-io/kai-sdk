package messaging_test

import (
	kai "github.com/konstellation-io/kre-runners/go-sdk/v1/protos"
	"github.com/konstellation-io/kre-runners/go-sdk/v1/sdk/messaging"
)

func (suite *SdkMessagingTestSuite) TestMessaging_GetErrorMessage_ExpectOk() {
	// Given
	kaiMessage := &kai.KaiNatsMessage{
		RequestId:   "some-request-id",
		MessageType: kai.MessageType_ERROR,
		Error:       "Error message",
	}
	objectStore := messaging.NewTestMessaging(suite.logger, nil, &suite.jetstream, kaiMessage, &suite.messageUtils)

	// When
	errorMessage := objectStore.GetErrorMessage()

	// Then
	suite.NotNil(objectStore)
	suite.Equal("Error message", errorMessage)
}

func (suite *SdkMessagingTestSuite) TestMessaging_GetErrorMessage_NoErrorMessageExistWhenTypeOK_ExpectError() {
	// Given
	kaiMessage := &kai.KaiNatsMessage{
		RequestId:   "some-request-id",
		MessageType: kai.MessageType_OK,
		Error:       "Some error message",
	}
	objectStore := messaging.NewTestMessaging(suite.logger, nil, &suite.jetstream, kaiMessage, &suite.messageUtils)

	// When
	errorMessage := objectStore.GetErrorMessage()

	// Then
	suite.NotNil(objectStore)
	suite.Empty(errorMessage)
}

func (suite *SdkMessagingTestSuite) TestMessaging_IsMessageOk_MessageOk_ExpectTrue() {
	// Given
	kaiMessage := &kai.KaiNatsMessage{
		RequestId:   "some-request-id",
		MessageType: kai.MessageType_OK,
	}
	objectStore := messaging.NewTestMessaging(suite.logger, nil, &suite.jetstream, kaiMessage, &suite.messageUtils)

	// When
	ok := objectStore.IsMessageOK()

	// Then
	suite.NotNil(objectStore)
	suite.True(ok)
}

func (suite *SdkMessagingTestSuite) TestMessaging_IsMessageOk_MessageNotOk_ExpectFalse() {
	// Given
	kaiMessage := &kai.KaiNatsMessage{
		RequestId:   "some-request-id",
		MessageType: kai.MessageType_ERROR,
		Error:       "Some error message",
	}
	objectStore := messaging.NewTestMessaging(suite.logger, nil, &suite.jetstream, kaiMessage, &suite.messageUtils)

	// When
	ok := objectStore.IsMessageOK()

	// Then
	suite.NotNil(objectStore)
	suite.False(ok)
}

func (suite *SdkMessagingTestSuite) TestMessaging_IsMessageError_MessageError_ExpectTrue() {
	// Given
	kaiMessage := &kai.KaiNatsMessage{
		RequestId:   "some-request-id",
		MessageType: kai.MessageType_ERROR,
		Error:       "Some error message",
	}
	objectStore := messaging.NewTestMessaging(suite.logger, nil, &suite.jetstream, kaiMessage, &suite.messageUtils)

	// When
	isError := objectStore.IsMessageError()

	// Then
	suite.NotNil(objectStore)
	suite.True(isError)
}

func (suite *SdkMessagingTestSuite) TestMessaging_IsMessageError_MessageNotError_ExpectFalse() {
	// Given
	kaiMessage := &kai.KaiNatsMessage{
		RequestId:   "some-request-id",
		MessageType: kai.MessageType_OK,
	}
	objectStore := messaging.NewTestMessaging(suite.logger, nil, &suite.jetstream, kaiMessage, &suite.messageUtils)

	// When
	isError := objectStore.IsMessageError()

	// Then
	suite.NotNil(objectStore)
	suite.False(isError)
}

func (suite *SdkMessagingTestSuite) TestMessaging_IsMessageEarlyReply_MessageEarlyReply_ExpectTrue() {
	// Given
	kaiMessage := &kai.KaiNatsMessage{
		RequestId:   "some-request-id",
		MessageType: kai.MessageType_EARLY_REPLY,
	}
	objectStore := messaging.NewTestMessaging(suite.logger, nil, &suite.jetstream, kaiMessage, &suite.messageUtils)

	// When
	isEarlyReply := objectStore.IsMessageEarlyReply()

	// Then
	suite.NotNil(objectStore)
	suite.True(isEarlyReply)
}

func (suite *SdkMessagingTestSuite) TestMessaging_IsMessageEarlyReply_MessageNotEarlyReply_ExpectFalse() {
	// Given
	kaiMessage := &kai.KaiNatsMessage{
		RequestId:   "some-request-id",
		MessageType: kai.MessageType_OK,
	}
	objectStore := messaging.NewTestMessaging(suite.logger, nil, &suite.jetstream, kaiMessage, &suite.messageUtils)

	// When
	isEarlyReply := objectStore.IsMessageEarlyReply()

	// Then
	suite.NotNil(objectStore)
	suite.False(isEarlyReply)
}

func (suite *SdkMessagingTestSuite) TestMessaging_IsMessageEarlyExit_MessageEarlyExit_ExpectTrue() {
	// Given
	kaiMessage := &kai.KaiNatsMessage{
		RequestId:   "some-request-id",
		MessageType: kai.MessageType_EARLY_EXIT,
	}
	objectStore := messaging.NewTestMessaging(suite.logger, nil, &suite.jetstream, kaiMessage, &suite.messageUtils)

	// When
	isEarlyExit := objectStore.IsMessageEarlyExit()

	// Then
	suite.NotNil(objectStore)
	suite.True(isEarlyExit)
}

func (suite *SdkMessagingTestSuite) TestMessaging_IsMessageEarlyExit_MessageNotEarlyExit_ExpectFalse() {
	// Given
	kaiMessage := &kai.KaiNatsMessage{
		RequestId:   "some-request-id",
		MessageType: kai.MessageType_OK,
	}
	objectStore := messaging.NewTestMessaging(suite.logger, nil, &suite.jetstream, kaiMessage, &suite.messageUtils)

	// When
	isEarlyExit := objectStore.IsMessageEarlyExit()

	// Then
	suite.NotNil(objectStore)
	suite.False(isEarlyExit)
}

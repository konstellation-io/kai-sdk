package messaging_test

import (
	"fmt"
	"github.com/go-logr/logr"
	"github.com/go-logr/zapr"
	"github.com/konstellation-io/kre-runners/go-sdk/v1/mocks"
	"github.com/konstellation-io/kre-runners/go-sdk/v1/sdk/messaging"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/suite"
	"go.uber.org/zap"
	"testing"
)

type SdkMessagingTestSuite struct {
	suite.Suite
	logger    logr.Logger
	jetstream mocks.JetStreamContextMock
}

func (suite *SdkMessagingTestSuite) SetupTest() {
	zapLog, err := zap.NewDevelopment()
	if err != nil {
		panic(fmt.Sprintf("who watches the watchmen (%v)?", err))
	}

	// Reset viper values before each test
	viper.Reset()

	suite.logger = zapr.NewLogger(zapLog)
	suite.jetstream = *mocks.NewJetStreamContextMock(suite.T())
}

func (suite *SdkMessagingTestSuite) TestMessaging_NewMessaging_ExpectOk() {
	// When
	objectStore := messaging.NewMessaging(suite.logger, nil, &suite.jetstream, nil)

	// Then
	suite.NotNil(objectStore)
}

func (suite *SdkMessagingTestSuite) TestMessaging_NewMessagingWithJetStream_ExpectOk() {
	// Given
	viper.SetDefault("nats.jetstream", "jetstream")
	suite.jetstream.On("JetStream").Return(&suite.jetstream, nil)
	objectStore := messaging.NewMessaging(suite.logger, nil, nil, nil)

	// When

	// Then
	suite.NotNil(objectStore)
}

func TestSdkMessagingTestSuite(t *testing.T) {
	suite.Run(t, new(SdkMessagingTestSuite))
}

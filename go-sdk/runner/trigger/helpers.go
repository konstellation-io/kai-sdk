package trigger

import (
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"google.golang.org/protobuf/types/known/anypb"

	"github.com/konstellation-io/kai-sdk/go-sdk/runner/common"
	"github.com/konstellation-io/kai-sdk/go-sdk/sdk"
)

const (
	_initializerLoggerName     = "[INITIALIZER]"
	_runnerLoggerName          = "[RUNNER]"
	_responseHandlerLoggerName = "[RESPONSE HANDLER]"
	_finalizerLoggerName       = "[FINALIZER]"
)

func composeInitializer(initializer common.Initializer) common.Initializer {
	return func(sdk sdk.KaiSDK) {
		logger := sdk.Logger.WithName(_initializerLoggerName)

		logger.V(1).Info("Initializing TriggerRunner...")
		common.InitializeProcessConfiguration(sdk)

		if initializer != nil {
			logger.V(3).Info("Executing user initializer...")
			initializer(sdk)
			logger.V(3).Info("User initializer executed")
		}

		logger.V(1).Info("TriggerRunner initialized")
	}
}

func composeRunner(userRunner RunnerFunc) RunnerFunc {
	return func(runner *Runner, kaiSDK sdk.KaiSDK) {
		logger := kaiSDK.Logger.WithName(_runnerLoggerName)

		logger.V(1).Info("Running TriggerRunner...")

		if userRunner != nil {
			logger.V(3).Info("Executing user runner...")

			go userRunner(runner, kaiSDK)
		}

		// Handle sigterm and await termChan signal
		termChan := make(chan os.Signal, 1)
		signal.Notify(termChan, syscall.SIGINT, syscall.SIGTERM)
		<-termChan

		logger.V(3).Info("User runner executed")

		// Handle shutdown
		logger.Info("Shutting down runner...")
		logger.V(1).Info("Closing opened channels...")
		runner.responseChannels.Range(func(key, value interface{}) bool {
			close(value.(chan *anypb.Any))
			logger.V(1).Info(fmt.Sprintf("Channel closed for request id %s", key), sdk.LoggerRequestID, key)
			return true
		})

		logger.Info("RunnerFunc shutdown")
		wg.Done()
	}
}

func getResponseHandler(handlers *sync.Map) ResponseHandler {
	return func(kaiSDK sdk.KaiSDK, response *anypb.Any) error {
		logger := kaiSDK.Logger.WithName(_responseHandlerLoggerName)

		// Unmarshal response to a KaiNatsMessage type
		logger.Info(fmt.Sprintf("Message received with request id %s", kaiSDK.GetRequestID()), sdk.LoggerRequestID, kaiSDK.GetRequestID())

		responseHandler, ok := handlers.LoadAndDelete(kaiSDK.GetRequestID())

		if ok {
			responseHandler.(chan *anypb.Any) <- response
			return nil
		}

		logger.V(1).Info(fmt.Sprintf("Undefined handler for the message with request id %s",
			kaiSDK.GetRequestID()), sdk.LoggerRequestID, kaiSDK.GetRequestID())

		return nil
	}
}

func composeFinalizer(userFinalizer common.Finalizer) common.Finalizer {
	return func(kaiSDK sdk.KaiSDK) {
		logger := kaiSDK.Logger.WithName(_finalizerLoggerName)

		logger.V(1).Info("Finalizing TriggerRunner...")

		if userFinalizer != nil {
			logger.V(3).Info("Executing user finalizer...")
			userFinalizer(kaiSDK)
			logger.V(3).Info("User finalizer executed")
		}

		logger.V(1).Info("TriggerRunner finalized")
	}
}

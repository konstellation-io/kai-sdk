package trigger

import (
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
		sdk.Logger.WithName(_initializerLoggerName).V(1).Info("Initializing TriggerRunner...")
		common.InitializeProcessConfiguration(sdk)

		if initializer != nil {
			sdk.Logger.WithName(_initializerLoggerName).V(3).Info("Executing user initializer...")
			initializer(sdk)
			sdk.Logger.WithName(_initializerLoggerName).V(3).Info("User initializer executed")
		}

		sdk.Logger.WithName(_initializerLoggerName).V(1).Info("TriggerRunner initialized")
	}
}

func composeRunner(userRunner RunnerFunc) RunnerFunc {
	return func(runner *Runner, sdk sdk.KaiSDK) {
		sdk.Logger.WithName(_runnerLoggerName).V(1).Info("Running TriggerRunner...")

		if userRunner != nil {
			sdk.Logger.WithName(_runnerLoggerName).V(3).Info("Executing user runner...")

			go userRunner(runner, sdk)
		}

		// Handle sigterm and await termChan signal
		termChan := make(chan os.Signal, 1)
		signal.Notify(termChan, syscall.SIGINT, syscall.SIGTERM)
		<-termChan

		sdk.Logger.WithName(_runnerLoggerName).V(3).Info("User runner executed")

		// Handle shutdown
		sdk.Logger.WithName(_runnerLoggerName).Info("Shutting down runner...")
		sdk.Logger.WithName(_runnerLoggerName).V(1).Info("Closing opened channels...")
		runner.responseChannels.Range(func(key, value interface{}) bool {
			close(value.(chan *anypb.Any))
			sdk.Logger.V(1).WithName(_runnerLoggerName).Info("Channel closed for requestID",
				"RequestID", key)
			return true
		})

		sdk.Logger.WithName(_runnerLoggerName).Info("RunnerFunc shutdown")
		wg.Done()
	}
}

func getResponseHandler(handlers *sync.Map) ResponseHandler {
	return func(sdk sdk.KaiSDK, response *anypb.Any) error {
		// Unmarshal response to a KaiNatsMessage type
		sdk.Logger.WithName(_responseHandlerLoggerName).Info("Message received", "RequestID", sdk.GetRequestID())

		responseHandler, ok := handlers.LoadAndDelete(sdk.GetRequestID())

		if ok {
			responseHandler.(chan *anypb.Any) <- response
			return nil
		}

		sdk.Logger.V(1).WithName(_responseHandlerLoggerName).Info("Undefined handler for the message",
			"RequestID", sdk.GetRequestID())

		return nil
	}
}

func composeFinalizer(userFinalizer common.Finalizer) common.Finalizer {
	return func(sdk sdk.KaiSDK) {
		sdk.Logger.V(1).WithName(_finalizerLoggerName).Info("Finalizing TriggerRunner...")

		if userFinalizer != nil {
			sdk.Logger.V(3).WithName(_finalizerLoggerName).Info("Executing user finalizer...")
			userFinalizer(sdk)
			sdk.Logger.V(3).WithName(_finalizerLoggerName).Info("User finalizer executed")
		}

		sdk.Logger.V(1).WithName(_finalizerLoggerName).Info("TriggerRunner finalized")
	}
}

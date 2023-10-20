package trigger

import (
	"os"
	"os/signal"
	"sync"
	"syscall"

	kaiCommon "github.com/konstellation-io/kai-sdk/go-sdk/internal/common"

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
		sdk.Logger = sdk.Logger.WithName(_initializerLoggerName)

		sdk.Logger.V(1).Info("Initializing TriggerRunner...")
		common.InitializeProcessConfiguration(sdk)

		if initializer != nil {
			sdk.Logger.V(3).Info("Executing user initializer...")
			initializer(sdk)
			sdk.Logger.V(3).Info("User initializer executed")
		}

		sdk.Logger.V(1).Info("TriggerRunner initialized")
	}
}

func composeRunner(userRunner RunnerFunc) RunnerFunc {
	return func(runner *Runner, kaiSDK sdk.KaiSDK) {
		kaiSDK.Logger = kaiSDK.Logger.
			WithName(_runnerLoggerName).
			WithValues(
				kaiCommon.LoggerProductID, kaiSDK.Metadata.GetProduct(),
				kaiCommon.LoggerVersionID, kaiSDK.Metadata.GetVersion(),
				kaiCommon.LoggerWorkflowID, kaiSDK.Metadata.GetWorkflow(),
				kaiCommon.LoggerProcessID, kaiSDK.Metadata.GetProcess(),
			)

		kaiSDK.Logger.V(1).Info("Running TriggerRunner...")

		if userRunner != nil {
			kaiSDK.Logger.V(3).Info("Executing user runner...")

			go userRunner(runner, kaiSDK)
		}

		// Handle sigterm and await termChan signal
		termChan := make(chan os.Signal, 1)
		signal.Notify(termChan, syscall.SIGINT, syscall.SIGTERM)
		<-termChan

		kaiSDK.Logger.V(3).Info("User runner executed")

		// Handle shutdown
		kaiSDK.Logger.Info("Shutting down runner...")
		kaiSDK.Logger.V(1).Info("Closing opened channels...")
		runner.responseChannels.Range(func(key, value interface{}) bool {
			close(value.(chan *anypb.Any))
			kaiSDK.Logger.V(1).Info("Channel closed for requestID",
				sdk.LoggerRequestID, key)
			return true
		})

		kaiSDK.Logger.Info("RunnerFunc shutdown")
		wg.Done()
	}
}

func getResponseHandler(handlers *sync.Map) ResponseHandler {
	return func(kaiSDK sdk.KaiSDK, response *anypb.Any) error {
		// Unmarshal response to a KaiNatsMessage type
		kaiSDK.Logger.WithName(_responseHandlerLoggerName).
			Info("Message received", sdk.LoggerRequestID, kaiSDK.GetRequestID())

		responseHandler, ok := handlers.LoadAndDelete(kaiSDK.GetRequestID())

		if ok {
			responseHandler.(chan *anypb.Any) <- response
			return nil
		}

		kaiSDK.Logger.WithName(_responseHandlerLoggerName).V(1).Info("Undefined handler for the message",
			sdk.LoggerRequestID, kaiSDK.GetRequestID())

		return nil
	}
}

func composeFinalizer(userFinalizer common.Finalizer) common.Finalizer {
	return func(kaiSDK sdk.KaiSDK) {
		kaiSDK.Logger = kaiSDK.Logger.WithName(_finalizerLoggerName)

		kaiSDK.Logger.V(1).Info("Finalizing TriggerRunner...")

		if userFinalizer != nil {
			kaiSDK.Logger.V(3).Info("Executing user finalizer...")
			userFinalizer(kaiSDK)
			kaiSDK.Logger.V(3).Info("User finalizer executed")
		}

		kaiSDK.Logger.V(1).Info("TriggerRunner finalized")
	}
}

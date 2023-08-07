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

func composeInitializer(initializer common.Initializer) common.Initializer {
	return func(sdk sdk.KaiSDK) {
		sdk.Logger.WithName("[RUNNER]").V(1).Info("Initializing TriggerRunner...")
		common.InitializeProcessConfiguration(sdk)

		if initializer != nil {
			sdk.Logger.WithName("[RUNNER]").V(3).Info("Executing user initializer...")
			initializer(sdk)
			sdk.Logger.WithName("[RUNNER]").V(3).Info("User initializer executed")
		}
		
		sdk.Logger.WithName("[RUNNER]").V(1).Info("TriggerRunner initialized")
	}
}

func composeRunner(triggerRunner *Runner, userRunner RunnerFunc) RunnerFunc {
	return func(runner *Runner, sdk sdk.KaiSDK) {
		sdk.Logger.WithName("[RUNNER]").V(1).Info("Running TriggerRunner...")

		if userRunner != nil {
			sdk.Logger.WithName("[RUNNER]").V(3).Info("Executing user runner...")
			userRunner(triggerRunner, sdk)
			sdk.Logger.WithName("[RUNNER]").V(3).Info("User runner executed")
		}

		// Handle sigterm and await termChan signal
		termChan := make(chan os.Signal, 1)
		signal.Notify(termChan, syscall.SIGINT, syscall.SIGTERM)
		<-termChan

		// Handle shutdown
		sdk.Logger.WithName("[RUNNER]").Info("Shutting down runner...")
		sdk.Logger.WithName("[RUNNER]").V(1).Info("Closing opened channels...")
		runner.responseChannels.Range(func(key, value interface{}) bool {
			close(value.(chan *anypb.Any))
			sdk.Logger.V(1).WithName("[RUNNER]").Info("Channel closed for requestID",
				"RequestID", key)
			return true
		})

		wg.Done()
		sdk.Logger.WithName("[RUNNER]").Info("RunnerFunc shutdown")
	}
}

func getResponseHandler(handlers *sync.Map) ResponseHandler {
	return func(sdk sdk.KaiSDK, response *anypb.Any) error {
		// Unmarshal response to a KaiNatsMessage type
		sdk.Logger.WithName("[RESPONSE HANDLER]").Info("Message received", "RequestID", sdk.GetRequestID())

		responseHandler, ok := handlers.LoadAndDelete(sdk.GetRequestID())

		if ok {
			responseHandler.(chan *anypb.Any) <- response
			return nil
		}

		sdk.Logger.V(1).WithName("[RESPONSE HANDLER]").Info("No handler found for the message",
			"RequestID", sdk.GetRequestID())

		return nil
	}
}

func composeFinalizer(userFinalizer common.Finalizer) common.Finalizer {
	return func(sdk sdk.KaiSDK) {
		sdk.Logger.V(1).WithName("[FINALIZER]").Info("Finalizing TriggerRunner...")

		if userFinalizer != nil {
			sdk.Logger.V(3).WithName("[FINALIZER]").Info("Executing user finalizer...")
			userFinalizer(sdk)
			sdk.Logger.V(3).WithName("[FINALIZER]").Info("User finalizer executed")
		}

		sdk.Logger.V(1).WithName("[FINALIZER]").Info("TriggerRunner finalized")
	}
}

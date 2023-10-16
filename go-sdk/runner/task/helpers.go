package task

import (
	"github.com/konstellation-io/kai-sdk/go-sdk/runner/common"
	"github.com/konstellation-io/kai-sdk/go-sdk/sdk"
	"github.com/konstellation-io/kai-sdk/go-sdk/sdk/kaiconstants"
	"google.golang.org/protobuf/types/known/anypb"
)

const (
	_initializerLoggerName   = "[INITIALIZER]"
	_preprocessorLoggerName  = "[PREPROCESSOR]"
	_handlerLoggerName       = "[HANDLER]"
	_postprocessorLoggerName = "[POSTPROCESSOR]"
	_finalizerLoggerName     = "[FINALIZER]"
)

func composeInitializer(initializer common.Initializer) common.Initializer {
	return func(sdk sdk.KaiSDK) {

		sdk.Logger = sdk.Logger.WithName(_initializerLoggerName)

		sdk.Logger.V(1).Info("Initializing TaskRunner...")
		common.InitializeProcessConfiguration(sdk)

		if initializer != nil {
			sdk.Logger.V(3).Info("Executing user initializer...")
			initializer(sdk)
			sdk.Logger.V(3).Info("User initializer executed")
		}

		sdk.Logger.V(1).Info("TaskRunner initialized")
	}
}

func composePreprocessor(preprocessor Preprocessor) Preprocessor {
	return func(sdk sdk.KaiSDK, response *anypb.Any) error {
		sdk.Logger = sdk.Logger.
			WithName(_preprocessorLoggerName).
			WithValues(
				kaiconstants.LoggerRequestID, sdk.GetRequestID(),
				kaiconstants.LoggerProductID, sdk.Metadata.GetProduct(),
				kaiconstants.LoggerVersionID, sdk.Metadata.GetVersion(),
				kaiconstants.LoggerWorkflowID, sdk.Metadata.GetWorkflow(),
				kaiconstants.LoggerProcessID, sdk.Metadata.GetProcess(),
			)

		sdk.Logger.V(1).Info("Preprocessing TaskRunner...")

		if preprocessor != nil {
			sdk.Logger.V(3).Info("Executing user preprocessor...")
			return preprocessor(sdk, response)
		}

		return nil
	}
}

func composeHandler(handler Handler) Handler {
	return func(sdk sdk.KaiSDK, response *anypb.Any) error {
		sdk.Logger = sdk.Logger.
			WithName(_handlerLoggerName).
			WithValues(
				kaiconstants.LoggerRequestID, sdk.GetRequestID(),
				kaiconstants.LoggerProductID, sdk.Metadata.GetProduct(),
				kaiconstants.LoggerVersionID, sdk.Metadata.GetVersion(),
				kaiconstants.LoggerWorkflowID, sdk.Metadata.GetWorkflow(),
				kaiconstants.LoggerProcessID, sdk.Metadata.GetProcess(),
			)

		sdk.Logger.V(1).Info("Handling TaskRunner...")

		if handler != nil {
			sdk.Logger.V(3).Info("Executing user handler...")
			return handler(sdk, response)
		}

		return nil
	}
}

func composePostprocessor(postprocessor Postprocessor) Postprocessor {
	return func(sdk sdk.KaiSDK, response *anypb.Any) error {
		sdk.Logger = sdk.Logger.
			WithName(_postprocessorLoggerName).
			WithValues(
				kaiconstants.LoggerRequestID, sdk.GetRequestID(),
				kaiconstants.LoggerProductID, sdk.Metadata.GetProduct(),
				kaiconstants.LoggerVersionID, sdk.Metadata.GetVersion(),
				kaiconstants.LoggerWorkflowID, sdk.Metadata.GetWorkflow(),
				kaiconstants.LoggerProcessID, sdk.Metadata.GetProcess(),
			)

		sdk.Logger.V(1).Info("Postprocessing TaskRunner...")

		if postprocessor != nil {
			sdk.Logger.V(3).Info("Executing user postprocessor...")
			return postprocessor(sdk, response)
		}

		return nil
	}
}

func composeFinalizer(finalizer common.Finalizer) common.Finalizer {
	return func(sdk sdk.KaiSDK) {
		sdk.Logger = sdk.Logger.WithName(_finalizerLoggerName)

		sdk.Logger.V(1).Info("Finalizing TaskRunner...")

		if finalizer != nil {
			sdk.Logger.V(3).Info("Executing user finalizer...")
			finalizer(sdk)
			sdk.Logger.V(3).Info("User finalizer executed")
		}
	}
}

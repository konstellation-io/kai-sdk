package exit

import (
	"github.com/konstellation-io/kai-sdk/go-sdk/runner/common"
	"github.com/konstellation-io/kai-sdk/go-sdk/sdk"
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
	return func(kaiSDK sdk.KaiSDK) {
		kaiSDK.Logger = kaiSDK.Logger.WithName(_initializerLoggerName)

		kaiSDK.Logger.V(1).Info("Initializing ExitRunner...")
		common.InitializeProcessConfiguration(kaiSDK)

		if initializer != nil {
			kaiSDK.Logger.V(3).Info("Executing user initializer...")
			initializer(kaiSDK)
			kaiSDK.Logger.V(3).Info("User initializer executed")
		}

		kaiSDK.Logger.V(1).Info("ExitRunner initialized")
	}
}

//nolint:dupl //Needed duplicated code
func composePreprocessor(preprocessor Preprocessor) Preprocessor {
	return func(kaiSDK sdk.KaiSDK, response *anypb.Any) error {
		kaiSDK.Logger = kaiSDK.Logger.
			WithName(_preprocessorLoggerName).
			WithValues(
				sdk.LoggerRequestID, kaiSDK.GetRequestID(),
				sdk.LoggerProductID, kaiSDK.Metadata.GetProduct(),
				sdk.LoggerVersionID, kaiSDK.Metadata.GetVersion(),
				sdk.LoggerWorkflowID, kaiSDK.Metadata.GetWorkflow(),
				sdk.LoggerProcessID, kaiSDK.Metadata.GetProcess(),
			)
		kaiSDK.Logger.V(1).Info("Preprocessing ExitRunner...")

		if preprocessor != nil {
			kaiSDK.Logger.
				V(3).
				Info("Executing user preprocessor...")

			return preprocessor(kaiSDK, response)
		}

		return nil
	}
}

//nolint:dupl //Needed duplicated code
func composeHandler(handler Handler) Handler {
	return func(kaiSDK sdk.KaiSDK, response *anypb.Any) error {
		kaiSDK.Logger = kaiSDK.Logger.
			WithName(_handlerLoggerName).
			WithValues(
				sdk.LoggerRequestID, kaiSDK.GetRequestID(),
				sdk.LoggerProductID, kaiSDK.Metadata.GetProduct(),
				sdk.LoggerVersionID, kaiSDK.Metadata.GetVersion(),
				sdk.LoggerWorkflowID, kaiSDK.Metadata.GetWorkflow(),
				sdk.LoggerProcessID, kaiSDK.Metadata.GetProcess(),
			)

		kaiSDK.Logger.V(1).Info("Handling ExitRunner...")

		if handler != nil {
			kaiSDK.Logger.V(3).Info("Executing user handler...")
			return handler(kaiSDK, response)
		}

		return nil
	}
}

//nolint:dupl //Needed duplicated code
func composePostprocessor(postprocessor Postprocessor) Postprocessor {
	return func(kaiSDK sdk.KaiSDK, response *anypb.Any) error {
		kaiSDK.Logger = kaiSDK.Logger.
			WithName(_postprocessorLoggerName).
			WithValues(
				sdk.LoggerRequestID, kaiSDK.GetRequestID(),
				sdk.LoggerProductID, kaiSDK.Metadata.GetProduct(),
				sdk.LoggerVersionID, kaiSDK.Metadata.GetVersion(),
				sdk.LoggerWorkflowID, kaiSDK.Metadata.GetWorkflow(),
				sdk.LoggerProcessID, kaiSDK.Metadata.GetProcess(),
			)

		kaiSDK.Logger.V(1).Info("Postprocessing ExitRunner...")

		if postprocessor != nil {
			kaiSDK.Logger.V(3).Info("Executing user postprocessor...")
			return postprocessor(kaiSDK, response)
		}

		return nil
	}
}

func composeFinalizer(finalizer common.Finalizer) common.Finalizer {
	return func(kaiSDK sdk.KaiSDK) {
		kaiSDK.Logger = kaiSDK.Logger.WithName(_finalizerLoggerName)

		kaiSDK.Logger.V(1).Info("Finalizing ExitRunner...")

		if finalizer != nil {
			kaiSDK.Logger.V(3).Info("Executing user finalizer...")
			finalizer(kaiSDK)
			kaiSDK.Logger.V(3).Info("User finalizer executed")
		}
	}
}

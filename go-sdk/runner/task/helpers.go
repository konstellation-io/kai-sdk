package task

import (
	kaiCommon "github.com/konstellation-io/kai-sdk/go-sdk/internal/common"
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

		kaiSDK.Logger.V(1).Info("Initializing TaskRunner...")
		common.InitializeProcessConfiguration(kaiSDK)

		if initializer != nil {
			kaiSDK.Logger.V(3).Info("Executing user initializer...")
			initializer(kaiSDK)
			kaiSDK.Logger.V(3).Info("User initializer executed")
		}

		kaiSDK.Logger.V(1).Info("TaskRunner initialized")
	}
}

//nolint:dupl //Needed duplicated code
func composePreprocessor(preprocessor Preprocessor) Preprocessor {
	return func(kaiSDK sdk.KaiSDK, response *anypb.Any) error {
		kaiSDK.Logger = kaiSDK.Logger.
			WithName(_preprocessorLoggerName).
			WithValues(
				kaiCommon.LoggerRequestID, kaiSDK.GetRequestID(),
				kaiCommon.LoggerProductID, kaiSDK.Metadata.GetProduct(),
				kaiCommon.LoggerVersionID, kaiSDK.Metadata.GetVersion(),
				kaiCommon.LoggerWorkflowID, kaiSDK.Metadata.GetWorkflow(),
				kaiCommon.LoggerProcessID, kaiSDK.Metadata.GetProcess(),
			)

		kaiSDK.Logger.V(1).Info("Preprocessing TaskRunner...")

		if preprocessor != nil {
			kaiSDK.Logger.V(3).Info("Executing user preprocessor...")
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
				kaiCommon.LoggerRequestID, kaiSDK.GetRequestID(),
				kaiCommon.LoggerProductID, kaiSDK.Metadata.GetProduct(),
				kaiCommon.LoggerVersionID, kaiSDK.Metadata.GetVersion(),
				kaiCommon.LoggerWorkflowID, kaiSDK.Metadata.GetWorkflow(),
				kaiCommon.LoggerProcessID, kaiSDK.Metadata.GetProcess(),
			)

		kaiSDK.Logger.V(1).Info("Handling TaskRunner...")

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
				kaiCommon.LoggerRequestID, kaiSDK.GetRequestID(),
				kaiCommon.LoggerProductID, kaiSDK.Metadata.GetProduct(),
				kaiCommon.LoggerVersionID, kaiSDK.Metadata.GetVersion(),
				kaiCommon.LoggerWorkflowID, kaiSDK.Metadata.GetWorkflow(),
				kaiCommon.LoggerProcessID, kaiSDK.Metadata.GetProcess(),
			)

		kaiSDK.Logger.V(1).Info("Postprocessing TaskRunner...")

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

		kaiSDK.Logger.V(1).Info("Finalizing TaskRunner...")

		if finalizer != nil {
			kaiSDK.Logger.V(3).Info("Executing user finalizer...")
			finalizer(kaiSDK)
			kaiSDK.Logger.V(3).Info("User finalizer executed")
		}
	}
}

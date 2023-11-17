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
		logger := kaiSDK.Logger.WithName(_initializerLoggerName)

		logger.V(1).Info("Initializing ExitRunner...")
		common.InitializeProcessConfiguration(kaiSDK)

		if initializer != nil {
			kaiSDK.Logger.V(3).Info("Executing user initializer...")
			initializer(kaiSDK)
			kaiSDK.Logger.V(3).Info("User initializer executed")
		}

		kaiSDK.Logger.V(1).Info("ExitRunner initialized")
	}
}

func composePreprocessor(preprocessor Preprocessor) Preprocessor {
	return func(kaiSDK sdk.KaiSDK, response *anypb.Any) error {
		logger := kaiSDK.Logger.WithName(_preprocessorLoggerName)

		logger.V(1).Info("Preprocessing ExitRunner...")

		if preprocessor != nil {
			logger.V(3).Info("Executing user preprocessor...")

			return preprocessor(kaiSDK, response)
		}

		return nil
	}
}

func composeHandler(handler Handler) Handler {
	return func(kaiSDK sdk.KaiSDK, response *anypb.Any) error {
		kaiSDK.Logger = kaiSDK.Logger.WithName(_handlerLoggerName)

		kaiSDK.Logger.V(1).Info("Handling ExitRunner...")

		if handler != nil {
			kaiSDK.Logger.V(3).Info("Executing user handler...")
			return handler(kaiSDK, response)
		}

		return nil
	}
}

func composePostprocessor(postprocessor Postprocessor) Postprocessor {
	return func(kaiSDK sdk.KaiSDK, response *anypb.Any) error {
		kaiSDK.Logger = kaiSDK.Logger.WithName(_postprocessorLoggerName)

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

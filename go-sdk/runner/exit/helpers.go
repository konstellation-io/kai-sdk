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
	return func(sdk sdk.KaiSDK) {
		sdk.Logger.WithName(_initializerLoggerName).V(1).Info("Initializing ExitRunner...")
		common.InitializeProcessConfiguration(sdk)

		if initializer != nil {
			sdk.Logger.WithName(_initializerLoggerName).V(3).Info("Executing user initializer...")
			initializer(sdk)
			sdk.Logger.WithName(_initializerLoggerName).V(3).Info("User initializer executed")
		}

		sdk.Logger.WithName(_initializerLoggerName).V(1).Info("ExitRunner initialized")
	}
}

func composePreprocessor(preprocessor Preprocessor) Preprocessor {
	return func(sdk sdk.KaiSDK, response *anypb.Any) error {
		sdk.Logger.WithName(_preprocessorLoggerName).V(1).Info("Preprocessing ExitRunner...")

		if preprocessor != nil {
			sdk.Logger.WithName(_preprocessorLoggerName).V(3).Info("Executing user preprocessor...")
			return preprocessor(sdk, response)
		}

		return nil
	}
}

func composeHandler(handler Handler) Handler {
	return func(sdk sdk.KaiSDK, response *anypb.Any) error {
		sdk.Logger.WithName(_handlerLoggerName).V(1).Info("Handling ExitRunner...")

		if handler != nil {
			sdk.Logger.WithName(_handlerLoggerName).V(3).Info("Executing user handler...")
			return handler(sdk, response)
		}

		return nil
	}
}

func composePostprocessor(postprocessor Postprocessor) Postprocessor {
	return func(sdk sdk.KaiSDK, response *anypb.Any) error {
		sdk.Logger.WithName(_postprocessorLoggerName).V(1).Info("Postprocessing ExitRunner...")

		if postprocessor != nil {
			sdk.Logger.WithName(_postprocessorLoggerName).V(3).Info("Executing user postprocessor...")
			return postprocessor(sdk, response)
		}

		return nil
	}
}

func composeFinalizer(finalizer common.Finalizer) common.Finalizer {
	return func(sdk sdk.KaiSDK) {
		sdk.Logger.WithName(_finalizerLoggerName).V(1).Info("Finalizing ExitRunner...")

		if finalizer != nil {
			sdk.Logger.WithName(_finalizerLoggerName).V(3).Info("Executing user finalizer...")
			finalizer(sdk)
			sdk.Logger.WithName(_finalizerLoggerName).V(3).Info("User finalizer executed")
		}
	}
}

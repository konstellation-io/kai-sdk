package task

import (
	"github.com/konstellation-io/kai-sdk/go-sdk/runner/common"
	"github.com/konstellation-io/kai-sdk/go-sdk/sdk"
	"google.golang.org/protobuf/types/known/anypb"
)

func composeInitializer(initializer common.Initializer) common.Initializer {
	return func(sdk sdk.KaiSDK) {
		sdk.Logger.WithName("[INITIALIZER]").V(1).Info("Initializing TaskRunner...")
		common.InitializeProcessConfiguration(sdk)

		if initializer != nil {
			sdk.Logger.WithName("[INITIALIZER]").V(3).Info("Executing user initializer...")
			initializer(sdk)
			sdk.Logger.WithName("[INITIALIZER]").V(3).Info("User initializer executed")
		}
		sdk.Logger.WithName("[INITIALIZER]").V(1).Info("TaskRunner initialized")
	}
}

func composePreprocessor(preprocessor Preprocessor) Preprocessor {
	return func(sdk sdk.KaiSDK, response *anypb.Any) error {
		sdk.Logger.WithName("[PREPROCESSOR]").V(1).Info("Preprocessing TaskRunner...")

		if preprocessor != nil {
			sdk.Logger.WithName("[PREPROCESSOR]").V(3).Info("Executing user preprocessor...")
			return preprocessor(sdk, response)
		}

		return nil
	}
}

func composeHandler(handler Handler) Handler {
	return func(sdk sdk.KaiSDK, response *anypb.Any) error {
		sdk.Logger.WithName("[HANDLER]").V(1).Info("Handling TaskRunner...")

		if handler != nil {
			sdk.Logger.WithName("[HANDLER]").V(3).Info("Executing user handler...")
			return handler(sdk, response)
		}

		return nil
	}
}

func composePostprocessor(postprocessor Postprocessor) Postprocessor {
	return func(sdk sdk.KaiSDK, response *anypb.Any) error {
		sdk.Logger.WithName("[POSTPROCESSOR]").V(1).Info("Postprocessing TaskRunner...")

		if postprocessor != nil {
			sdk.Logger.WithName("[POSTPROCESSOR]").V(3).Info("Executing user postprocessor...")
			return postprocessor(sdk, response)
		}

		return nil
	}
}

func composeFinalizer(finalizer common.Finalizer) common.Finalizer {
	return func(sdk sdk.KaiSDK) {
		sdk.Logger.WithName("[FINALIZER]").V(1).Info("Finalizing TaskRunner...")

		if finalizer != nil {
			sdk.Logger.WithName("[FINALIZER]").V(3).Info("Executing user finalizer...")
			finalizer(sdk)
			sdk.Logger.WithName("[FINALIZER]").V(3).Info("User finalizer executed")
		}
	}
}

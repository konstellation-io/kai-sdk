package main

import (
	"github.com/golang/protobuf/ptypes/wrappers"
	"github.com/konstellation-io/kai-sdk/go-sdk/runner"
	"github.com/konstellation-io/kai-sdk/go-sdk/sdk"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"
)

func main() {
	runner.
		NewRunner().
		ExitRunner().
		WithInitializer(initializer).
		WithPreprocessor(func(ctx sdk.KaiSDK, response *anypb.Any) error {
			return nil
		}).
		WithHandler(defaultHandler).
		WithPostprocessor(func(ctx sdk.KaiSDK, response *anypb.Any) error {
			return nil
		}).
		WithFinalizer(func(sdk sdk.KaiSDK) {
			sdk.Logger.Info("Finalizer")
		}).
		Run()
}

func initializer(sdk sdk.KaiSDK) {
	value, err := sdk.CentralizedConfig.GetConfig("test")
	if err != nil {
		sdk.Logger.Error(err, "Error getting config")
		return
	}
	sdk.Logger.Info("Config value retrieved!", "value", value)

	obj, err := sdk.Storage.Ephemeral.Get("test")
	if err != nil {
		sdk.Logger.Error(err, "Error getting Obj Store values")
		return
	}
	sdk.Logger.Info("ObjectStore value retrieved!", "object", string(obj))
}

func defaultHandler(sdk sdk.KaiSDK, response *anypb.Any) error {
	stringValue := &wrappers.StringValue{}

	// Unmarshall response to StringValue
	err := proto.Unmarshal(response.GetValue(), stringValue)
	if err != nil {
		sdk.Logger.Error(err, "Error unmarshalling response")
		return err
	}

	stringValue = &wrappers.StringValue{
		Value: stringValue.GetValue() + ", Processed by the exit process!",
	}

	err = sdk.Messaging.SendOutput(stringValue)
	if err != nil {
		sdk.Logger.Error(err, "Error sending output")
		return err
	}

	return nil
}

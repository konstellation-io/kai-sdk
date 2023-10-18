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
		WithPreprocessor(func(kaiSDK sdk.KaiSDK, response *anypb.Any) error {
			return nil
		}).
		WithHandler(defaultHandler).
		WithPostprocessor(func(kaiSDK sdk.KaiSDK, response *anypb.Any) error {
			return nil
		}).
		WithFinalizer(func(kaiSDK sdk.KaiSDK) {
			kaiSDK.Logger.Info("Finalizer")
		}).
		Run()
}

func initializer(kaiSDK sdk.KaiSDK) {
	value, err := kaiSDK.CentralizedConfig.GetConfig("test")
	if err != nil {
		kaiSDK.Logger.Error(err, "Error getting config")
		return
	}
	kaiSDK.Logger.Info("Config value retrieved!", "value", value)

	obj, err := kaiSDK.ObjectStore.Get("test")
	if err != nil {
		kaiSDK.Logger.Error(err, "Error getting Obj Store values")
		return
	}
	kaiSDK.Logger.Info("ObjectStore value retrieved!", "object", string(obj))
}

func defaultHandler(kaiSDK sdk.KaiSDK, response *anypb.Any) error {
	stringValue := &wrappers.StringValue{}

	// Unmarshall response to StringValue
	err := proto.Unmarshal(response.GetValue(), stringValue)
	if err != nil {
		kaiSDK.Logger.Error(err, "Error unmarshalling response")
		return err
	}

	stringValue = &wrappers.StringValue{
		Value: stringValue.GetValue() + ", Processed by the exit process!",
	}

	err = kaiSDK.Messaging.SendOutput(stringValue)
	if err != nil {
		kaiSDK.Logger.Error(err, "Error sending output")
		return err
	}

	return nil
}

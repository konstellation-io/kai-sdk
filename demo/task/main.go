package main

import (
	"fmt"
	"github.com/golang/protobuf/ptypes/wrappers"
	"github.com/konstellation-io/kai-sdk/go-sdk/runner"
	"github.com/konstellation-io/kai-sdk/go-sdk/sdk"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"
)

func main() {
	runner.
		NewRunner().
		TaskRunner().
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

	obj, err := kaiSDK.Storage.Ephemeral.Get("test")
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

	objVersion, err := kaiSDK.Storage.Persistent.Save("some-object", []byte("Some value"))
	if err != nil {
		kaiSDK.Logger.Error(err, "Error saving value in the persistent storage")
		return err
	}

	stringValue = &wrappers.StringValue{
		Value: fmt.Sprintf("%s, Processed by the task process, "+
			"uploaded object to persistent storage with ID %s!",
			stringValue.GetValue(),
			objVersion.VersionID,
		),
	}

	err = kaiSDK.Messaging.SendOutput(stringValue)
	if err != nil {
		kaiSDK.Logger.Error(err, "Error sending output")
		return err
	}

	return nil
}

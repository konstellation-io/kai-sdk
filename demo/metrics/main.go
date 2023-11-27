package main

import (
	"context"
	"fmt"

	"github.com/golang/protobuf/ptypes/wrappers"
	"github.com/konstellation-io/kai-sdk/go-sdk/runner"
	"github.com/konstellation-io/kai-sdk/go-sdk/sdk"
	"go.opentelemetry.io/otel/metric"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"
)

func main() {
	runner.
		NewRunner().
		TaskRunner().
		WithInitializer(initializer).
		WithHandler(defaultHandler).
		WithFinalizer(func(kaiSDK sdk.KaiSDK) {
			kaiSDK.Logger.Info("Finalizer")
		}).
		Run()
}

var GlobalCounter metric.Int64Counter

func initializer(kaiSDK sdk.KaiSDK) {
	kaiSDK.Logger.Info("Initializer")

	metricsClient := kaiSDK.Measurements.GetMetricsClient().MetricsClient
	counter, err := metricsClient.Int64Counter("messages_received")
	if err != nil {
		kaiSDK.Logger.Error(err, "Error creating counter")
	}
	GlobalCounter = counter
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
		Value: fmt.Sprintf("%s, Processed by the metrics process, "+
			"and increased the counter by 1",
			stringValue.GetValue(),
		),
	}
	GlobalCounter.Add(context.Background(), 1)

	err = kaiSDK.Messaging.SendOutput(stringValue)
	if err != nil {
		kaiSDK.Logger.Error(err, "Error sending output")
		return err
	}

	return nil
}

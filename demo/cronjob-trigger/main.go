package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/golang/protobuf/ptypes/wrappers"
	"github.com/google/uuid"
	"github.com/konstellation-io/kai-sdk/go-sdk/runner"
	"github.com/konstellation-io/kai-sdk/go-sdk/runner/trigger"
	"github.com/konstellation-io/kai-sdk/go-sdk/sdk"
	"github.com/robfig/cron/v3"
)

func main() {
	runner.
		NewRunner().
		TriggerRunner().
		WithInitializer(initializer).
		WithRunner(cronjobRunner).
		WithFinalizer(func(kaiSDK sdk.KaiSDK) {
			kaiSDK.Logger.Info("Finalizer")
		}).
		Run()
}

func initializer(kaiSDK sdk.KaiSDK) {
	kaiSDK.Logger.V(1).Info("Metadata",
		"process", kaiSDK.Metadata.GetProcess(),
		"product", kaiSDK.Metadata.GetProduct(),
		"workflow", kaiSDK.Metadata.GetWorkflow(),
		"version", kaiSDK.Metadata.GetVersion(),
		"kv_product", kaiSDK.Metadata.GetProductCentralizedConfigurationName(),
		"kv_workflow", kaiSDK.Metadata.GetWorkflowCentralizedConfigurationName(),
		"kv_process", kaiSDK.Metadata.GetProcessCentralizedConfigurationName(),
		"ephemeral-storage", kaiSDK.Metadata.GetEphemeralStorageName(),
	)

	value1, _ := kaiSDK.CentralizedConfig.GetConfig("test1")
	value2, _ := kaiSDK.CentralizedConfig.GetConfig("test2")
	kaiSDK.Logger.Info("Config values from config.yaml", "test1", value1, "test2", value2)
}

func cronjobRunner(tr *trigger.Runner, kaiSDK sdk.KaiSDK) {
	kaiSDK.Logger.Info("Starting cronjob runner")

	c := cron.New(
		cron.WithLogger(kaiSDK.Logger.WithName("[CRONJOB]")),
	)

	_, err := c.AddFunc("@every 30s", func() {
		requestID := uuid.New().String()
		kaiSDK.Logger.Info("Cronjob triggered, new message sent", sdk.LoggerRequestID, requestID)

		val := wrappers.StringValue{
			Value: "Hi there, I'm a cronjob!",
		}

		err := kaiSDK.Messaging.SendOutputWithRequestID(&val, requestID)
		if err != nil {
			kaiSDK.Logger.Error(err, "Error sending output", sdk.LoggerRequestID, requestID)
			return
		}
	})
	if err != nil {
		kaiSDK.Logger.Error(err, "Error adding cronjob")
		return
	}

	c.Start()

	// Handle sigterm and await termChan signal
	termChan := make(chan os.Signal, 1)
	signal.Notify(termChan, syscall.SIGINT, syscall.SIGTERM)
	<-termChan

	c.Stop()
}

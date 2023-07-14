package main

import (
	"github.com/golang/protobuf/ptypes/wrappers"
	"github.com/google/uuid"
	"github.com/konstellation-io/kai-sdk/go-sdk/runner"
	"github.com/konstellation-io/kai-sdk/go-sdk/runner/trigger"
	"github.com/konstellation-io/kai-sdk/go-sdk/sdk"
	"github.com/robfig/cron/v3"

	"os"
	"os/signal"
	"syscall"
)

func main() {
	runner.
		NewRunner().
		TriggerRunner().
		WithInitializer(initializer).
		WithRunner(cronjobRunner).
		WithFinalizer(func(ctx sdk.KaiSDK) {
			ctx.Logger.Info("Finalizer")
		}).
		Run()
}

func initializer(sdk sdk.KaiSDK) {
	sdk.Logger.V(1).Info("Metadata",
		"process", sdk.Metadata.GetProcess(),
		"product", sdk.Metadata.GetProduct(),
		"workflow", sdk.Metadata.GetWorkflow(),
		"version", sdk.Metadata.GetVersion(),
		"kv_product", sdk.Metadata.GetKeyValueStoreProductName(),
		"kv_workflow", sdk.Metadata.GetKeyValueStoreWorkflowName(),
		"kv_process", sdk.Metadata.GetKeyValueStoreProcessName(),
		"object-store", sdk.Metadata.GetObjectStoreName(),
	)

	sdk.Logger.V(1).Info("PathUtils",
		"getBasePath", sdk.PathUtils.GetBasePath(),
		"composeBasePath", sdk.PathUtils.ComposePath("test"))

	value1, _ := sdk.CentralizedConfig.GetConfig("test1")
	value2, _ := sdk.CentralizedConfig.GetConfig("test2")
	sdk.Logger.Info("Config values from config.yaml", "test1", value1, "test2", value2)
}

func cronjobRunner(tr *trigger.Runner, sdk sdk.KaiSDK) {
	sdk.Logger.Info("Starting cronjob runner")

	c := cron.New(
		cron.WithLogger(sdk.Logger),
	)

	_, err := c.AddFunc("@every 30s", func() {
		requestID := uuid.New().String()
		sdk.Logger.Info("Cronjob triggered, new message sent", "requestID", requestID)

		val := wrappers.StringValue{
			Value: "Hi there, I'm a cronjob!",
		}

		err := sdk.Messaging.SendOutputWithRequestID(&val, requestID)
		if err != nil {
			sdk.Logger.Error(err, "Error sending output")
			return
		}
	})
	if err != nil {
		sdk.Logger.Error(err, "Error adding cronjob")
		return
	}

	c.Start()

	// Handle sigterm and await termChan signal
	termChan := make(chan os.Signal, 1)
	signal.Notify(termChan, syscall.SIGINT, syscall.SIGTERM)
	<-termChan
}

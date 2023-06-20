package main

import (
	"github.com/golang/protobuf/ptypes/wrappers"
	"github.com/google/uuid"
	"github.com/konstellation-io/kre-runners/go-sdk/v1/runner"
	"github.com/konstellation-io/kre-runners/go-sdk/v1/runner/trigger"
	"github.com/konstellation-io/kre-runners/go-sdk/v1/sdk"
	"github.com/nats-io/nats.go"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	runner.
		NewRunner().
		TriggerRunner().
		WithInitializer(initializer).
		WithRunner(natsSubscriberRunner).
		WithFinalizer(func(sdk sdk.KaiSDK) {
			sdk.Logger.Info("Finalizer")
		}).
		Run()
}

func initializer(sdk sdk.KaiSDK) {
	sdk.Logger.Info("Writing test value to the object store", "value", "testValue")
	err := sdk.ObjectStore.Save("test", []byte("testValue"))
	if err != nil {
		sdk.Logger.Error(err, "Error saving object")
	}

	sdk.Logger.Info("Writing test value to the centralized config",
		"value", "testConfigValue")
	err = sdk.CentralizedConfig.SetConfig("test", "testConfigValue")
	if err != nil {
		sdk.Logger.Error(err, "Error setting config")
	}

	sdk.Logger.V(1).Info("Metadata",
		"process", sdk.Metadata.GetProcess(),
		"product", sdk.Metadata.GetProduct(),
		"workflow", sdk.Metadata.GetWorkflow(),
		"version", sdk.Metadata.GetVersion(),
		"kv_product", sdk.Metadata.GetKeyValueStoreProductName(),
		"kv_workflow", sdk.Metadata.GetKeyValueStoreWorkflowName(),
		"kv_process", sdk.Metadata.GetKeyValueStoreProcessName(),
		"object_store", sdk.Metadata.GetObjectStoreName(),
	)

	sdk.Logger.V(1).Info("PathUtils",
		"getBasePath", sdk.PathUtils.GetBasePath(),
		"composeBasePath", sdk.PathUtils.ComposeBasePath("test"))
}

func natsSubscriberRunner(tr *trigger.TriggerRunner, sdk sdk.KaiSDK) {
	sdk.Logger.Info("Starting nats subscriber")

	nc, _ := nats.Connect("nats://localhost:4222")
	js, err := nc.JetStream()
	if err != nil {
		panic(err)
	}

	s, err := js.QueueSubscribe(
		"nats-trigger",
		"nats-trigger-queue",
		func(msg *nats.Msg) {
			val := &wrappers.StringValue{
				Value: "Hi there, I'm a NATS subscriber!",
			}

			requestID := uuid.New().String()

			responseChannel := tr.GetResponseChannel(requestID)

			err = sdk.Messaging.SendOutputWithRequestID(val, requestID)
			if err != nil {
				sdk.Logger.Error(err, "Error sending output")
				return
			}

			// Wait for the response before ACKing the message
			<-responseChannel

			sdk.Logger.Info("Message received, acking message")

			err = msg.Ack()
			if err != nil {
				sdk.Logger.Error(err, "Error acking message")
				return
			}
		},
		nats.DeliverNew(),
		nats.Durable("nats-trigger"),
		nats.ManualAck(),
		nats.AckWait(22*time.Hour),
	)

	// Handle sigterm and await termChan signal
	termChan := make(chan os.Signal, 1)
	signal.Notify(termChan, syscall.SIGINT, syscall.SIGTERM)
	<-termChan

	err = s.Unsubscribe()
	if err != nil {
		sdk.Logger.Error(err, "Error unsubscribing")
		return
	}
}

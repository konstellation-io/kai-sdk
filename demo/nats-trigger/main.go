package main

import (
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/nats-io/nats.go"

	"github.com/golang/protobuf/ptypes/wrappers"
	"github.com/google/uuid"
	"github.com/konstellation-io/kai-sdk/go-sdk/runner"
	"github.com/konstellation-io/kai-sdk/go-sdk/runner/trigger"
	"github.com/konstellation-io/kai-sdk/go-sdk/sdk"
)

func main() {
	runner.
		NewRunner().
		TriggerRunner().
		WithInitializer(initializer).
		WithRunner(natsSubscriberRunner).
		WithFinalizer(func(kaiSDK sdk.KaiSDK) {
			kaiSDK.Logger.Info("Finalizer")
		}).
		Run()
}

func initializer(kaiSDK sdk.KaiSDK) {
	kaiSDK.Logger.Info("Writing test value to the ephemeral store", "value", "testValue")
	err := kaiSDK.Storage.Ephemeral.Save("test", []byte("testValue"))
	if err != nil {
		kaiSDK.Logger.Error(err, "Error saving object")
	}

	kaiSDK.Logger.Info("Writing test value to the centralized config",
		"value", "testConfigValue")
	err = kaiSDK.CentralizedConfig.SetConfig("test", "testConfigValue")
	if err != nil {
		kaiSDK.Logger.Error(err, "Error setting config")
	}

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
}

func natsSubscriberRunner(tr *trigger.Runner, kaiSDK sdk.KaiSDK) {
	kaiSDK.Logger.Info("Starting nats subscriber")

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

			err = kaiSDK.Messaging.SendOutputWithRequestID(val, requestID)
			if err != nil {
				kaiSDK.Logger.Error(err, "Error sending output")
				return
			}

			// Wait for the response before ACKing the message
			<-responseChannel

			kaiSDK.Logger.Info("Message received, acking message")

			err = msg.Ack()
			if err != nil {
				kaiSDK.Logger.Error(err, "Error acking message")
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
		kaiSDK.Logger.Error(err, "Error unsubscribing")
		return
	}
}

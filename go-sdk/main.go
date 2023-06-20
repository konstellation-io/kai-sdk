package main

import (
	"fmt"
	"github.com/labstack/gommon/log"
	"github.com/nats-io/nats.go"
	"github.com/spf13/viper"
)

func main() {
	fmt.Println("Initializing NATS config...")
	fmt.Println()
	publishMessageToNats()
}

func publishMessageToNats() {
	nc, err := nats.Connect(viper.GetString("nats.url"))
	if err != nil {
		panic(err)
	}

	// Use the JetStream sdk to produce and consumer messages
	// that have been persisted.
	js, err := nc.JetStream()
	if err != nil {
		panic(err)
	}

	_, err = js.AddStream(&nats.StreamConfig{
		Name: "trigger",
		Subjects: []string{
			"trigger-output",
		},
		Retention: nats.InterestPolicy,
	})
	if err != nil {
		log.Error("Error creating stream", err)
	}

	_, err = js.AddStream(&nats.StreamConfig{
		Name: "task",
		Subjects: []string{
			"task-input",
			"task-output",
		},
		Retention: nats.InterestPolicy,
	})
	if err != nil {
		log.Error("Error creating stream", err)
	}

	_, err = js.AddStream(&nats.StreamConfig{
		Name: "exit",
		Subjects: []string{
			"exit-input",
			"exit-output",
		},
		Retention: nats.InterestPolicy,
	})
	if err != nil {
		log.Error("Error creating stream", err)
	}

	js.CreateKeyValue(&nats.KeyValueConfig{
		Bucket:  "product",
		Storage: nats.MemoryStorage,
	})

	js.CreateKeyValue(&nats.KeyValueConfig{
		Bucket:  "workflow",
		Storage: nats.MemoryStorage,
	})

	js.CreateKeyValue(&nats.KeyValueConfig{
		Bucket:  "process",
		Storage: nats.MemoryStorage,
	})

	js.CreateObjectStore(&nats.ObjectStoreConfig{
		Bucket: "object-store",
	})
}

package messaging

import (
	"fmt"

	"github.com/nats-io/nats.go"
	"github.com/spf13/viper"
)

//go:generate mockery --name messageUtils --output ../../mocks --structname MessageUtilsMock --filename message_utils_mock.go
type messageUtils interface {
	GetMaxMessageSize() (int64, error)
}

type MessageUtilsImpl struct {
	jetstream nats.JetStreamContext
	nats      *nats.Conn
}

func NewMessageUtils(ns *nats.Conn, js nats.JetStreamContext) MessageUtilsImpl {
	return MessageUtilsImpl{
		nats:      ns,
		jetstream: js,
	}
}

func (mu MessageUtilsImpl) GetMaxMessageSize() (int64, error) {
	streamInfo, err := mu.jetstream.StreamInfo(viper.GetString("nats.stream"))
	if err != nil {
		return 0, fmt.Errorf("error getting stream's max message size: %w", err)
	}

	streamMaxSize := int64(streamInfo.Config.MaxMsgSize)
	serverMaxSize := mu.nats.MaxPayload()

	if streamMaxSize == -1 {
		return serverMaxSize, nil
	}

	if streamMaxSize < serverMaxSize {
		return streamMaxSize, nil
	}

	return serverMaxSize, nil
}

func sizeInMB(size int64) string {
	return fmt.Sprintf("%.1f MB", float32(size)/1024/1024)
}

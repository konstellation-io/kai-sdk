package main

import (
	"math/rand"
	"time"
	localProto "training-go/proto"

	"github.com/konstellation-io/kai-sdk/go-sdk/runner"
	"github.com/konstellation-io/kai-sdk/go-sdk/sdk"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"
)

func main() {
	runner.
		NewRunner().
		TaskRunner().
		WithHandler(defaultHandler).
		Run()
}

func defaultHandler(sdk sdk.KaiSDK, response *anypb.Any) error {
	splitterResult := &localProto.Splitter{}

	err := proto.Unmarshal(response.GetValue(), splitterResult)
	if err != nil {
		sdk.Logger.Error(err, "Error unmarshalling incoming message")
		return err
	}

	sdk.Logger.Info("Training model", "url", splitterResult.GetRepoUrl())

	result := &localProto.Training{
		TrainingId: splitterResult.GetTrainingId(),
		ModelId:    GenerateRandString(),
	}

	err = sdk.Messaging.SendOutput(result)
	if err != nil {
		sdk.Logger.Error(err, "Error sending output")
		return err
	}

	return nil
}

func GenerateRandString() string {
	const charset = "abcdefghijklmnopqrstuvwxyz" +
		"ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

	var seededRand *rand.Rand = rand.New(
		rand.NewSource(time.Now().UnixNano()),
	)

	b := make([]byte, seededRand.Intn(20))
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}

	return string(b)
}

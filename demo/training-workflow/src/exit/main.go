package main

import (
	localProto "exit/proto"

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
	validationResult := &localProto.Validation{}

	err := proto.Unmarshal(response.GetValue(), validationResult)
	if err != nil {
		sdk.Logger.Error(err, "Error unmarshalling incoming message")
		return err
	}

	var topModel string
	var topScore int32
	for _, modelScore := range validationResult.GetModelsScores() {
		if modelScore.Score > topScore {
			topModel = modelScore.ModelId
			topScore = modelScore.Score
		}
	}

	// todo: save topModel to storage
	sdk.Logger.Info("Top model", "model", topModel, "score", topScore)

	return nil
}

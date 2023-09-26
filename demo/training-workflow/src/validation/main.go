package main

import (
	"math/rand"
	localProto "validation/proto"

	"github.com/konstellation-io/kai-sdk/go-sdk/runner"
	"github.com/konstellation-io/kai-sdk/go-sdk/sdk"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"
)

var (
	modelScores map[string][]*localProto.ModelScore
)

func main() {
	runner.
		NewRunner().
		TaskRunner().
		WithHandler(defaultHandler).
		Run()
}

func defaultHandler(sdk sdk.KaiSDK, response *anypb.Any) error {
	trainingResult := &localProto.Training{}

	err := proto.Unmarshal(response.GetValue(), trainingResult)
	if err != nil {
		sdk.Logger.Error(err, "Error unmarshalling incoming message")
		return err
	}

	result := validateModel(&sdk, trainingResult.ModelId)

	ok, scores := mergeScore(trainingResult.TrainingId, result)

	if ok {
		result := &localProto.Validation{
			TrainingId:   trainingResult.TrainingId,
			ModelsScores: scores,
		}
		err = sdk.Messaging.SendOutput(result)
		if err != nil {
			sdk.Logger.Error(err, "Error sending output")
			return err
		}

	}

	return nil
}

func validateModel(sdk *sdk.KaiSDK, modelID string) *localProto.ModelScore {
	counter := 0
	for _, letter := range modelID {
		counter = counter + int(letter)
	}

	r := rand.New(rand.NewSource(int64(counter)))
	score := r.Int31n(100)

	sdk.Logger.Info("Validating model", "modelID", modelID, "score", score)

	return &localProto.ModelScore{
		ModelId: modelID,
		Score:   score,
	}
}

func mergeScore(trainingID string, score *localProto.ModelScore) (bool, []*localProto.ModelScore) {
	scores, ok := modelScores[trainingID]
	if !ok {
		scores = []*localProto.ModelScore{}
	}

	scores = append(scores, score)

	if len(scores) == 2 {
		delete(modelScores, trainingID)
		return true, scores
	} else {
		modelScores[trainingID] = scores
		return false, scores
	}
}

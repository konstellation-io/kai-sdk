package main

import (
	"errors"
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	localProto "validation/proto"

	"github.com/konstellation-io/kai-sdk/go-sdk/runner"
	"github.com/konstellation-io/kai-sdk/go-sdk/sdk"
	centralizedconfiguration "github.com/konstellation-io/kai-sdk/go-sdk/sdk/centralized-configuration"
	"github.com/konstellation-io/kai-sdk/go-sdk/sdk/messaging"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"
)

func main() {
	runner.
		NewRunner().
		TaskRunner().
		WithInitializer(initializer).
		WithHandler(defaultHandler).
		Run()
}

func initializer(sdk sdk.KaiSDK) {
	sdk.Logger.Info("Go validation is ready!")
}

func defaultHandler(sdk sdk.KaiSDK, response *anypb.Any) error {
	trainingResult := &localProto.Training{}

	err := proto.Unmarshal(response.GetValue(), trainingResult)
	if err != nil {
		sdk.Logger.Error(err, "Error unmarshalling incoming message")
		return err
	}

	sdk.Logger.Info("Received message", "trainingID", trainingResult.TrainingId)

	result := validateModel(&sdk, trainingResult.ModelId)

	scores, err := mergeScore(&sdk, trainingResult.TrainingId, result)
	if err != nil {
		sdk.Logger.Error(err, "Error merging scores")
		return err
	}

	if scores != nil {
		sdk.Logger.Info("Recieved 2 results for same training, sending output", "trainingID", trainingResult.TrainingId)

		result := &localProto.Validation{
			TrainingId:   trainingResult.TrainingId,
			ModelsScores: scores,
		}

		err = sdk.Messaging.SendOutput(result)
		if err != nil {
			sdk.Logger.Error(err, "Error sending output")
			return err
		}
	} else {
		sdk.Logger.Info("Waiting for second result", "trainingID", trainingResult.TrainingId)
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

func mergeScore(sdk *sdk.KaiSDK, trainingID string, score *localProto.ModelScore) ([]*localProto.ModelScore, error) {

	modelScoreString, err := sdk.CentralizedConfig.GetConfig(trainingID, messaging.ProcessScope)
	if err != nil && !errors.Is(err, centralizedconfiguration.ErrKeyNotFound) {
		return nil, err
	}

	if modelScoreString == "" {
		sdk.Logger.Info("Model Score not found", "trainingID", trainingID)

		newConfig := fmt.Sprintf("%s,%d", score.ModelId, score.Score)
		err = sdk.CentralizedConfig.SetConfig(trainingID, newConfig, messaging.ProcessScope)
		if err != nil {
			return nil, err
		}
		return nil, nil
	}

	sdk.Logger.Info("Model Score found", "trainingID", trainingID)

	configArray := strings.Split(modelScoreString, ",")

	parsedScore, err := strconv.Atoi(configArray[1])
	if err != nil {
		return nil, err
	}

	secondModelScore := &localProto.ModelScore{
		ModelId: configArray[0],
		Score:   int32(parsedScore),
	}

	scores := []*localProto.ModelScore{
		score,
		secondModelScore,
	}

	err = sdk.CentralizedConfig.DeleteConfig(trainingID, messaging.ProcessScope)
	if err != nil {
		return nil, err
	}

	return scores, nil
}

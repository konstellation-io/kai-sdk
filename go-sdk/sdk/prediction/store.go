package prediction

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/konstellation-io/kai-sdk/go-sdk/internal/common"
	"github.com/konstellation-io/kai-sdk/go-sdk/sdk/metadata"
	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
)

var (
	ErrPredictionNotFound = errors.New("prediction not found")
	ErrParsingListResult  = errors.New("error parsing prediction store search result")
)

type Filter struct {
	RequestID string
	Workflow  string
	Process   string
	Version   string
	Timestamp *TimestampRange
}

type TimestampRange struct {
	StartDate int64
	EndDate   int64
}

type RedisPredictionStore struct {
	requestID string
	client    *redis.Client
	metadata  *metadata.Metadata
}

func NewRedisPredictionStore(requestID string) *RedisPredictionStore {
	opts := &redis.Options{
		Addr:     strings.ReplaceAll(viper.GetString(common.ConfigRedisEndpointKey), "redis://", ""),
		Username: viper.GetString(common.ConfigRedisUsernameKey),
		Password: viper.GetString(common.ConfigRedisPasswordKey),
	}

	return &RedisPredictionStore{
		client:    redis.NewClient(opts),
		metadata:  metadata.NewMetadata(),
		requestID: requestID,
	}
}

func (r *RedisPredictionStore) Save(ctx context.Context, predictionID string, value map[string]interface{}) error {
	prediction := Prediction{
		Timestamp: time.Now().UnixMilli(),
		Payload:   value,
		Metadata: Metadata{
			Product:   r.metadata.GetProduct(),
			Version:   r.metadata.GetVersion(),
			Workflow:  r.metadata.GetWorkflow(),
			Process:   r.metadata.GetProcess(),
			RequestID: r.requestID,
		},
	}

	keyWithProductPrefix := r.getKeyWithProductPrefix(predictionID)

	return r.client.JSONSet(ctx, keyWithProductPrefix, "$", prediction).Err()
}

func (r *RedisPredictionStore) Get(ctx context.Context, predictionID string) (*Prediction, error) {
	result, err := r.client.JSONGet(ctx, r.getKeyWithProductPrefix(predictionID), "$").Result()
	if err != nil {
		return nil, err
	}

	if result == "" {
		return nil, ErrPredictionNotFound
	}

	var predictions []*Prediction
	if err := json.Unmarshal([]byte(result), &predictions); err != nil {
		return nil, fmt.Errorf("parsing prediction: %w", err)
	}

	return predictions[0], nil
}

func (r *RedisPredictionStore) Find(ctx context.Context, filter *Filter) ([]Prediction, error) {
	var predictions []Prediction

	result, err := r.client.Do(ctx, "FT.SEARCH", "predictionsIdx", r.buildQueryWithFilters(filter)).Result()
	if err != nil {
		return nil, err
	}

	predictions, err = r.parseResultToPredictionsList(result)
	if err != nil {
		return nil, err
	}

	return predictions, nil
}

func (r *RedisPredictionStore) getKeyWithProductPrefix(key string) string {
	return fmt.Sprintf("%s:%s", r.metadata.GetProduct(), key)
}

func (r *RedisPredictionStore) buildQueryWithFilters(filter *Filter) string {
	queryFilters := []string{
		fmt.Sprintf("@product:{%s} @timestamp:[0 inf]", r.metadata.GetProduct()),
	}

	if filter.Workflow != "" {
		queryFilters = append(queryFilters, fmt.Sprintf("@workflow:{%s}", filter.Workflow))
	}

	if filter.Version != "" {
		queryFilters = append(queryFilters, fmt.Sprintf("@version:{%s}", filter.Version))
	}

	if filter.Process != "" {
		queryFilters = append(queryFilters, fmt.Sprintf("@process:{%s}", filter.Process))
	}

	if filter.RequestID != "" {
		queryFilters = append(queryFilters, fmt.Sprintf("@requestID:{%s}", filter.RequestID))
	}

	if filter.Timestamp != nil {
		queryFilters = append(
			queryFilters,
			fmt.Sprintf("@timestamp:[%d %d]", filter.Timestamp.StartDate, filter.Timestamp.EndDate),
		)
	}

	return strings.ReplaceAll(strings.Join(queryFilters, " "), "-", "\\-")
}

func (r *RedisPredictionStore) parseResultToPredictionsList(rawResult interface{}) ([]Prediction, error) {
	result, ok := rawResult.(map[interface{}]interface{})
	if !ok {
		return nil, ErrParsingListResult
	}

	rawResults, ok := result["results"]
	if !ok {
		return nil, ErrParsingListResult
	}

	results, ok := rawResults.([]interface{})
	if !ok {
		return nil, ErrParsingListResult
	}

	predictions := make([]Prediction, 0, len(results))

	for _, res := range results {
		rawResContent, ok := res.(map[interface{}]interface{})
		if !ok {
			return nil, ErrParsingListResult
		}

		rawExtraAttributes, ok := rawResContent["extra_attributes"]
		if !ok {
			return nil, ErrParsingListResult
		}

		extraAttributes, ok := rawExtraAttributes.(map[interface{}]interface{})
		if !ok {
			return nil, ErrParsingListResult
		}

		resJSONRoot, ok := extraAttributes["$"]
		if !ok {
			return nil, ErrParsingListResult
		}

		resJSONRootString, ok := resJSONRoot.(string)
		if !ok {
			return nil, ErrParsingListResult
		}

		var prediction Prediction

		err := json.Unmarshal([]byte(resJSONRootString), &prediction)
		if err != nil {
			return nil, fmt.Errorf("unmarshalling prediction: %w", err)
		}

		predictions = append(predictions, prediction)
	}

	return predictions, nil
}

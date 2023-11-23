package prediction

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/konstellation-io/kai-sdk/go-sdk/internal/common"
	"github.com/konstellation-io/kai-sdk/go-sdk/sdk/metadata"
	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
)

type Filter struct {
	RequestID string
	Workflow  string
	Process   string
	Version   string
}

type RedisPredictionStore struct {
	requestID string
	client    *redis.Client
	metadata  *metadata.Metadata
}

func NewRedisPredictionStore(requestID string) *RedisPredictionStore {
	//opts, err := redis.ParseURL(viper.GetString(common.ConfigRedisEndpointKey))
	//if err != nil {
	//	return nil, fmt.Errorf("parsing Redis URL: %w", err)
	//}
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

func (r *RedisPredictionStore) Save(ctx context.Context, predictionID string, value map[string]string) error {
	//{
	// 	timestamp: time.Now(),
	//	payload: map[string]string,
	//	metadata: {
	//		 product, version, requestid, workflow, process
	//	}
	//}
	prediction := Prediction{
		Timestamp: strconv.FormatInt(time.Now().Unix(), 10),
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

	var prediction *Prediction
	if err := json.Unmarshal([]byte(result), &prediction); err != nil {
		return nil, fmt.Errorf("parsing prediction: %w", err)
	}

	return prediction, nil
}

func (r *RedisPredictionStore) List(ctx context.Context, filter Filter) ([]Prediction, error) {
	//var predictions []Prediction
	//
	//result, err := r.client.(ctx).Result()
	//if err != nil {
	//	return nil, err
	//}
	//
	//var prediction *Prediction
	//if err := json.Unmarshal([]byte(result), &prediction); err != nil {
	//	return nil, fmt.Errorf("parsing prediction: %w", err)
	//}
	return nil, nil
}

func (r *RedisPredictionStore) getKeyWithProductPrefix(key string) string {
	return fmt.Sprintf("%s:%s", r.metadata.GetProduct(), key)
}

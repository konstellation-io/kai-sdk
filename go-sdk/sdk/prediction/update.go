package prediction

import (
	"context"
	"time"
)

type UpdatePayloadFunc func(Payload) Payload

func (r *RedisPredictionStore) Update(ctx context.Context, predictionID string, updatePayload UpdatePayloadFunc) error {
	prediction, err := r.Get(ctx, predictionID)
	if err != nil {
		return err
	}

	prediction.Payload = updatePayload(prediction.Payload)
	prediction.LastModified = time.Now().UnixMilli()

	err = r.client.JSONSet(ctx, r.getKeyWithProductPrefix(predictionID), "$", prediction).Err()
	if err != nil {
		return err
	}

	return nil
}

package prediction

import (
	"context"
	"time"
)

func (r *RedisPredictionStore) Save(ctx context.Context, predictionID string, value Payload) error {
	prediction := Prediction{
		CreationDate: time.Now().UnixMilli(),
		LastModified: time.Now().UnixMilli(),
		Payload:      value,
		Metadata: Metadata{
			Product:      r.metadata.GetProduct(),
			Version:      r.metadata.GetVersion(),
			Workflow:     r.metadata.GetWorkflow(),
			WorkflowType: r.metadata.GetWorkflowType(),
			Process:      r.metadata.GetProcess(),
			RequestID:    r.requestID,
		},
	}

	keyWithProductPrefix := r.getKeyWithProductPrefix(predictionID)

	return r.client.JSONSet(ctx, keyWithProductPrefix, "$", prediction).Err()
}

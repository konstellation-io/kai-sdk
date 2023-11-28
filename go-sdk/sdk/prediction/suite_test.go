//go:build integration

package prediction_test

import (
	"context"
	"encoding/json"
	"fmt"
	"path/filepath"
	"testing"
	"time"

	"github.com/docker/go-connections/nat"
	"github.com/konstellation-io/kai-sdk/go-sdk/internal/common"
	"github.com/konstellation-io/kai-sdk/go-sdk/sdk/prediction"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/suite"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

const (
	_testRequestID    = "test-request-id"
	_testProduct      = "test-product"
	_testVersion      = "test-version"
	_testWorkflow     = "test-workflow"
	_testWorkflowType = "training"
	_testProcess      = "test-process"

	_testRedisUser          = "test-user"
	_testRedisPassword      = "test-password"
	_testRedisAdminPassword = "testpassword"
)

func TestPredictionStoreSuite(t *testing.T) {
	suite.Run(t, new(PredictionStoreSuite))
}

type PredictionStoreSuite struct {
	suite.Suite
	predictionStore *prediction.RedisPredictionStore
	redisContainer  testcontainers.Container
	redisClient     *redis.Client
}

func (s *PredictionStoreSuite) SetupSuite() {
	s.initializeRedisContainer()
	s.initializeMetadataConfiguration()
	s.initializeRedisClient()
	s.createRedisUser()
	s.initializePredictionStore()
}

func (s *PredictionStoreSuite) SetupTest() {
	s.createRedisPredictionIndex()
}

func (s *PredictionStoreSuite) TearDownTest() {
	err := s.redisClient.FlushAll(context.Background()).Err()
	s.Require().NoError(err)
}

func (s *PredictionStoreSuite) TestPredictionStore_Save_ExpectOK() {
	var (
		ctx              = context.Background()
		predictionID     = "test-prediction"
		predictionValues = map[string]interface{}{
			"test-key": "test-value",
		}
		expectedMetadata = s.getPredictionMetadata()
	)
	// WHEN
	err := s.predictionStore.Save(ctx, predictionID, predictionValues)
	s.Require().NoError(err)

	res, err := s.redisClient.JSONGet(ctx, s.getKeyWithProductPrefix(predictionID)).Result()
	s.Require().NoError(err)

	var actualPrediction []prediction.Prediction
	err = json.Unmarshal([]byte(res), &actualPrediction)
	s.Require().NoError(err)

	// THEN
	s.Equal(predictionValues, actualPrediction[0].Payload)
	s.Equal(expectedMetadata, actualPrediction[0].Metadata)
}

func (s *PredictionStoreSuite) TestPredictionStore_Get_ExpectOK() {
	var (
		ctx                = context.Background()
		predictionID       = "test-prediction"
		expectedPrediction = &prediction.Prediction{
			CreationDate: time.Now().UnixMilli(),
			Payload: map[string]interface{}{
				"test-key": "test-value",
			},
			Metadata: prediction.Metadata{},
		}
	)

	err := s.redisClient.
		JSONSet(ctx, s.getKeyWithProductPrefix(predictionID), "$", expectedPrediction).
		Err()
	s.Require().NoError(err)

	// WHEN
	actualPrediction, err := s.predictionStore.Get(ctx, predictionID)
	s.Require().NoError(err)

	// THEN
	s.Equal(expectedPrediction, actualPrediction)
}

func (s *PredictionStoreSuite) TestPredictionStore_Get_ProductWithoutPemissions() {
	var (
		ctx                    = context.Background()
		predictionID           = "test-prediction"
		otherProductPermission = &prediction.Prediction{
			CreationDate: time.Now().UnixMilli(),
			Payload: map[string]interface{}{
				"test-key": "test-value",
			},
			Metadata: prediction.Metadata{},
		}
	)

	err := s.redisClient.
		JSONSet(ctx, "other-product:"+predictionID, "$", otherProductPermission).
		Err()
	s.Require().NoError(err)

	// WHEN
	actualPrediction, err := s.predictionStore.Get(ctx, predictionID)

	// THEN
	s.ErrorIs(err, prediction.ErrPredictionNotFound)
	s.Nil(actualPrediction)
}

func (s *PredictionStoreSuite) TestPredictionStore_Get_FailsInvalidKey() {
	ctx := context.Background()

	// WHEN
	actualPrediction, err := s.predictionStore.Get(ctx, "not-found-key")

	// THEN
	s.ErrorIs(err, prediction.ErrPredictionNotFound)
	s.Nil(actualPrediction)
}

func (s *PredictionStoreSuite) TestPredictionStore_Find_ExpectOK() {
	var (
		ctx             = context.Background()
		firstPrediction = prediction.Prediction{
			CreationDate: time.Now().UnixMilli(),
			Payload: map[string]interface{}{
				"test-key-2": "test-value-2",
			},
			Metadata: s.getPredictionMetadata(),
		}

		secondPrediction = prediction.Prediction{
			CreationDate: time.Now().UnixMilli(),
			Payload: map[string]interface{}{
				"test-key-2": "test-value-2",
			},
			Metadata: s.getPredictionMetadata(),
		}

		expectedPredictions = []prediction.Prediction{firstPrediction, secondPrediction}
	)

	err := s.redisClient.
		JSONSet(ctx, s.getKeyWithProductPrefix("first-prediction"), "$", firstPrediction).
		Err()
	s.Require().NoError(err)

	err = s.redisClient.
		JSONSet(ctx, s.getKeyWithProductPrefix("second-prediction"), "$", secondPrediction).
		Err()
	s.Require().NoError(err)

	actualPredictions, err := s.predictionStore.Find(ctx, &prediction.Filter{
		Version: _testVersion,
		CreationDate: prediction.TimestampRange{
			EndDate: time.Now().UnixMilli(),
		},
	})
	s.Require().NoError(err)

	// THEN
	s.ElementsMatch(expectedPredictions, actualPredictions)
}

func (s *PredictionStoreSuite) TestPredictionStore_Find_FilterByMultipleFields_ExpectOK() {
	var (
		ctx             = context.Background()
		firstPrediction = prediction.Prediction{
			CreationDate: time.Now().UnixMilli(),
			Payload: map[string]interface{}{
				"test-key-2": "test-value-2",
			},
			Metadata: s.getPredictionMetadata(),
		}

		secondPrediction = prediction.Prediction{
			CreationDate: time.Now().UnixMilli(),
			Payload: map[string]interface{}{
				"test-key-2": "test-value-2",
			},
			Metadata: prediction.Metadata{
				Product:   firstPrediction.Metadata.Product,
				Version:   firstPrediction.Metadata.Version,
				Workflow:  "different-workflow",
				Process:   firstPrediction.Metadata.Process,
				RequestID: firstPrediction.Metadata.RequestID,
			},
		}

		expectedPredictions = []prediction.Prediction{firstPrediction}
	)

	err := s.redisClient.
		JSONSet(ctx, s.getKeyWithProductPrefix("first-prediction"), "$", firstPrediction).
		Err()
	s.Require().NoError(err)

	err = s.redisClient.
		JSONSet(ctx, s.getKeyWithProductPrefix("second-prediction"), "$", secondPrediction).
		Err()
	s.Require().NoError(err)

	actualPredictions, err := s.predictionStore.Find(ctx, &prediction.Filter{
		Version:   _testVersion,
		Workflow:  _testWorkflow,
		Process:   _testProcess,
		RequestID: _testRequestID,
		CreationDate: prediction.TimestampRange{
			EndDate: time.Now().UnixMilli(),
		},
	})
	s.Require().NoError(err)

	// THEN
	s.ElementsMatch(expectedPredictions, actualPredictions)
}

func (s *PredictionStoreSuite) TestPredictionStore_Find_FilterByTimestampRange_ExpectOK() {
	var (
		ctx             = context.Background()
		firstPrediction = prediction.Prediction{
			CreationDate: 1000,
			Payload: map[string]interface{}{
				"test-key-2": "test-value-2",
			},
			Metadata: s.getPredictionMetadata(),
		}

		secondPrediction = prediction.Prediction{
			CreationDate: 2000,
			Payload: map[string]interface{}{
				"test-key-2": "test-value-2",
			},
			Metadata: s.getPredictionMetadata(),
		}

		expectedPredictions = []prediction.Prediction{secondPrediction}
	)

	err := s.redisClient.
		JSONSet(ctx, s.getKeyWithProductPrefix("first-prediction"), "$", firstPrediction).
		Err()
	s.Require().NoError(err)

	err = s.redisClient.
		JSONSet(ctx, s.getKeyWithProductPrefix("second-prediction"), "$", secondPrediction).
		Err()
	s.Require().NoError(err)

	actualPredictions, err := s.predictionStore.Find(ctx, &prediction.Filter{
		CreationDate: prediction.TimestampRange{
			StartDate: 1500,
			EndDate:   3000,
		},
	})
	s.Require().NoError(err)

	// THEN
	s.ElementsMatch(expectedPredictions, actualPredictions)
}

func (s *PredictionStoreSuite) TestPredictionStore_Find_FailsIfIndexDoesNotExist() {
	var (
		ctx = context.Background()
	)

	// FlushAll deletes index too
	err := s.redisClient.FlushAll(ctx).Err()
	s.Require().NoError(err)

	actualPredictions, err := s.predictionStore.Find(ctx, &prediction.Filter{
		Version: _testVersion,
	})

	// THEN
	s.Error(err)
	s.Empty(actualPredictions)
}

func (s *PredictionStoreSuite) TestPredictionStore_Find_ProductFilterAppliedByDefault_ExpectOK() {
	var (
		ctx             = context.Background()
		firstPrediction = prediction.Prediction{
			CreationDate: time.Now().UnixMilli(),
			Payload: map[string]interface{}{
				"test-key-2": "test-value-2",
			},
			Metadata: s.getPredictionMetadata(),
		}

		secondPrediction = prediction.Prediction{
			CreationDate: time.Now().UnixMilli(),
			Payload: map[string]interface{}{
				"test-key-2": "test-value-2",
			},
			Metadata: prediction.Metadata{
				Product:   "different-product",
				Version:   firstPrediction.Metadata.Version,
				Workflow:  firstPrediction.Metadata.Version,
				Process:   firstPrediction.Metadata.Process,
				RequestID: firstPrediction.Metadata.RequestID,
			},
		}

		expectedPredictions = []prediction.Prediction{firstPrediction}
	)

	err := s.redisClient.
		JSONSet(ctx, s.getKeyWithProductPrefix("first-prediction"), "$", firstPrediction).
		Err()
	s.Require().NoError(err)

	err = s.redisClient.
		JSONSet(ctx, s.getKeyWithProductPrefix("second-prediction"), "$", secondPrediction).
		Err()
	s.Require().NoError(err)

	actualPredictions, err := s.predictionStore.Find(ctx, &prediction.Filter{
		Version: _testVersion,
		CreationDate: prediction.TimestampRange{
			EndDate: time.Now().UnixMilli(),
		},
	})
	s.Require().NoError(err)

	// THEN
	s.ElementsMatch(expectedPredictions, actualPredictions)
}

func (s *PredictionStoreSuite) initializeRedisContainer() {
	ctx := context.Background()
	exposedPort := "6379/tcp"

	testdataPath, err := filepath.Abs("./testdata")
	s.Require().NoError(err)

	redisConfContainerPath := "/usr/local/etc/redis"

	req := testcontainers.ContainerRequest{
		Image:        "redis/redis-stack-server:7.2.0-v6",
		ExposedPorts: []string{exposedPort},
		WaitingFor:   wait.ForLog("Ready to accept connections tcp"),
		Mounts: []testcontainers.ContainerMount{
			{
				Source: testcontainers.DockerBindMountSource{
					HostPath: testdataPath,
				},
				Target: testcontainers.ContainerMountTarget(redisConfContainerPath),
			},
		},
	}

	s.redisContainer, err = testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	s.Require().NoError(err)

	redisEndpoint, err := s.redisContainer.PortEndpoint(context.Background(), nat.Port(exposedPort), "")
	s.Require().NoError(err)

	viper.Set(common.ConfigRedisEndpointKey, redisEndpoint)
}

func (s *PredictionStoreSuite) initializeMetadataConfiguration() {
	viper.Set(common.ConfigMetadataProductIDKey, _testProduct)
	viper.Set(common.ConfigMetadataVersionIDKey, _testVersion)
	viper.Set(common.ConfigMetadataWorkflowIDKey, _testWorkflow)
	viper.Set(common.ConfigMetadataWorkflowTypeKey, _testWorkflowType)
	viper.Set(common.ConfigMetadataProcessIDKey, _testProcess)
	viper.Set(common.ConfigRedisUsernameKey, _testRedisUser)
	viper.Set(common.ConfigRedisPasswordKey, _testRedisPassword)
}

func (s *PredictionStoreSuite) initializePredictionStore() {
	s.predictionStore = prediction.NewRedisPredictionStore(_testRequestID)
}

func (s *PredictionStoreSuite) getClient(url, keyID, keySecret string) *minio.Client {
	redisClient, err := minio.New(url, &minio.Options{
		Creds:        credentials.NewStaticV4(keyID, keySecret, ""),
		BucketLookup: minio.BucketLookupPath,
	})
	s.Require().NoError(err)

	return redisClient
}

func (s *PredictionStoreSuite) initializeRedisClient() {
	opts := &redis.Options{
		Addr:     viper.GetString(common.ConfigRedisEndpointKey),
		Username: "default",
		Password: _testRedisAdminPassword,
	}
	s.redisClient = redis.NewClient(opts)
}

func (s *PredictionStoreSuite) createRedisUser() {
	var (
		allowedKeys = fmt.Sprintf("~%s:*", _testProduct)
		passwordArg = fmt.Sprintf(">%s", _testRedisPassword)
	)

	command := s.redisClient.Do(
		context.Background(),
		"ACL", "SETUSER", _testRedisUser, allowedKeys, passwordArg, "+@all", "+ft.search", "~predictionsIdx", "-@dangerous", "on",
	)
	s.Require().NoError(command.Err())
}

func (s *PredictionStoreSuite) createRedisPredictionIndex() {
	err := s.redisClient.Do(context.Background(),
		"FT.CREATE",
		"predictionsIdx",
		"ON",
		"JSON",
		"SCHEMA",
		"$.metadata.product", "AS", "product", "TAG",
		"$.metadata.version", "AS", "version", "TAG",
		"$.metadata.workflow", "AS", "workflow", "TAG",
		"$.metadata.workflowType", "AS", "workflowType", "TAG",
		"$.metadata.process", "AS", "process", "TAG",
		"$.metadata.requestID", "AS", "requestID", "TAG",
		"$.creationDate", "AS", "creationDate", "NUMERIC",
	).Err()
	s.Require().NoError(err)
}

func (s *PredictionStoreSuite) getPredictionMetadata() prediction.Metadata {
	return prediction.Metadata{
		Product:      viper.GetString(common.ConfigMetadataProductIDKey),
		Version:      viper.GetString(common.ConfigMetadataVersionIDKey),
		Workflow:     viper.GetString(common.ConfigMetadataWorkflowIDKey),
		WorkflowType: viper.GetString(common.ConfigMetadataWorkflowTypeKey),
		Process:      viper.GetString(common.ConfigMetadataProcessIDKey),
		RequestID:    _testRequestID,
	}
}

func (s *PredictionStoreSuite) getKeyWithProductPrefix(predictionID string) string {
	return fmt.Sprintf("%s:%s", _testProduct, predictionID)
}

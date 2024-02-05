package common

const (
	ConfigAppConfigPathKey                = "APP_CONFIG_PATH"
	ConfigRunnerLoggerLevelKey            = "runner.logger.level"
	ConfigRunnerLoggerOutputPathsKey      = "runner.logger.output_paths"
	ConfigRunnerLoggerErrorOutputPathsKey = "runner.logger.error_output_paths"
	ConfigRunnerLoggerEncodingKey         = "runner.logger.encoding"
	ConfigRunnerSubscriberAckWaitTimeKey  = "runner.subscriber.ack_wait_time"
	ConfigMetadataProductIDKey            = "metadata.product_id"
	ConfigMetadataWorkflowIDKey           = "metadata.workflow_name"
	ConfigMetadataWorkflowTypeKey         = "metadata.workflow_type"
	ConfigMetadataProcessIDKey            = "metadata.process_name"
	ConfigMetadataProcessTypeKey          = "metadata.process_type"
	ConfigMetadataVersionIDKey            = "metadata.version_tag"
	ConfigNatsURLKey                      = "nats.url"
	ConfigNatsStreamKey                   = "nats.stream"
	ConfigNatsOutputKey                   = "nats.output"
	ConfigNatsInputsKey                   = "nats.inputs"
	ConfigNatsEphemeralStorage            = "nats.object_store"
	ConfigCcGlobalBucketKey               = "centralized_configuration.global.bucket"
	ConfigCcProductBucketKey              = "centralized_configuration.product.bucket"
	ConfigCcWorkflowBucketKey             = "centralized_configuration.workflow.bucket"
	ConfigCcProcessBucketKey              = "centralized_configuration.process.bucket"
	ConfigMinioEndpointKey                = "minio.endpoint"
	ConfigMinioClientUserKey              = "minio.client_user"
	ConfigMinioClientPasswordKey          = "minio.client_password" //nolint:gosec // False positive
	ConfigMinioUseSslKey                  = "minio.ssl"
	ConfigMinioBucketKey                  = "minio.bucket"
	ConfigMinioInternalFolderKey          = "minio.internal_folder"
	ConfigAuthEndpointKey                 = "auth.endpoint"
	ConfigAuthClientKey                   = "auth.client"
	ConfigAuthClientSecretKey             = "auth.client_secret" //nolint:gosec // False positive
	ConfigAuthRealmKey                    = "auth.realm"
	ConfigRedisEndpointKey                = "predictions.endpoint"
	ConfigRedisUsernameKey                = "predictions.username"
	ConfigRedisPasswordKey                = "predictions.password"
	ConfigRedisIndexKey                   = "predictions.index"
	ConfigModelFolderNameKey              = "model_registry.folder_name"
	ConfigMeasurementsEndpointKey         = "measurements.endpoint"
	ConfigMeasurementsInsecureKey         = "measurements.insecure"
	ConfigMeasurementsTimeoutKey          = "measurements.timeout"
	ConfigMeasurementsMetricsIntervalKey  = "measurements.metrics_interval"
)

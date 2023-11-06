package runner

import (
	"fmt"
	"time"

	"github.com/go-logr/logr"
	"github.com/go-logr/zapr"
	"github.com/konstellation-io/kai-sdk/go-sdk/runner/exit"
	"github.com/konstellation-io/kai-sdk/go-sdk/runner/task"
	"github.com/konstellation-io/kai-sdk/go-sdk/runner/trigger"
	"github.com/nats-io/nats.go"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"golang.org/x/exp/slices"
)

type Runner struct {
	logger    logr.Logger
	nats      *nats.Conn
	jetstream nats.JetStreamContext
}

func NewRunner() *Runner {
	initializeConfiguration()

	logger := getLogger()

	nc, err := getNatsConnection(logger)
	if err != nil {
		panic(fmt.Errorf("fatal error connecting to NATS: %w", err))
	}

	js, err := getJetStreamConnection(logger, nc)
	if err != nil {
		panic(fmt.Errorf("fatal error connecting to JetStream: %w", err))
	}

	return &Runner{
		logger:    logger,
		nats:      nc,
		jetstream: js,
	}
}

func validateConfig(keys []string) {
	var mandatoryConfigKeys = []string{
		"metadata.product_id",
		"metadata.workflow_name",
		"metadata.process_name",
		"metadata.version_tag",
		"metadata.base_path",
		"nats.url",
		"nats.stream",
		"nats.output",
		"centralized_configuration.global.bucket",
		"centralized_configuration.product.bucket",
		"centralized_configuration.workflow.bucket",
		"centralized_configuration.process.bucket",
		"minio.endpoint",
		"minio.client_user",     // generated user for the bucket
		"minio.client_password", // generated user's password for the bucket
		"minio.ssl",             // Enable or disable SSL
		"minio.bucket",          // Bucket to be used
		"auth.endpoint",         // keycloak endpoint
		"auth.client",           // Client to be used to authenticate
		"auth.client_secret",    // Client's secret to be used
		"auth.realm",            // Realm
	}

	for _, key := range mandatoryConfigKeys {
		if !slices.Contains(keys, key) {
			panic(fmt.Sprintf("missing mandatory configuration key: %s", key))
		}
	}
}

func initializeConfiguration() {
	// Load environment variables
	viper.SetEnvPrefix("KAI")
	viper.AutomaticEnv()

	if viper.IsSet("APP_CONFIG_PATH") {
		viper.AddConfigPath(viper.GetString("APP_CONFIG_PATH"))
	}

	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")

	err := viper.ReadInConfig() //nolint:ineffassign,staticcheck // err is used in the next if

	viper.SetConfigName("app")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")

	if viper.IsSet("APP_CONFIG_PATH") {
		viper.AddConfigPath(viper.GetString("APP_CONFIG_PATH"))
	}

	err = viper.MergeInConfig()

	keys := viper.AllKeys()
	if len(keys) == 0 {
		panic(fmt.Errorf("configuration could not be loaded: %w", err))
	}

	validateConfig(keys)

	// Set viper default values
	viper.SetDefault("metadata.base_path", "/")
	viper.SetDefault("runner.subscriber.ack_wait_time", 22*time.Hour)
	viper.SetDefault("runner.logger.level", "InfoLevel")
	viper.SetDefault("runner.logger.encoding", "console")
	viper.SetDefault("runner.logger.output_paths", []string{"stdout"})
	viper.SetDefault("runner.logger.error_output_paths", []string{"stderr"})
}

func getNatsConnection(logger logr.Logger) (*nats.Conn, error) {
	nc, err := nats.Connect(viper.GetString("nats.url"))
	if err != nil {
		logger.Error(err, "Error connecting to NATS")
		return nil, err
	}

	return nc, nil
}

func getJetStreamConnection(logger logr.Logger, nc *nats.Conn) (nats.JetStreamContext, error) {
	js, err := nc.JetStream()
	if err != nil {
		logger.Error(err, "Error connecting to JetStream")
		return nil, err
	}

	return js, nil
}

func getLogger() logr.Logger {
	var log logr.Logger

	config := zap.NewDevelopmentConfig()

	logLevel, err := zap.ParseAtomicLevel(viper.GetString("runner.logger.level"))
	if err != nil {
		logLevel = zap.NewAtomicLevelAt(zap.InfoLevel)
	}

	config.Level = zap.NewAtomicLevelAt(logLevel.Level())
	config.OutputPaths = viper.GetStringSlice("runner.logger.output_paths")
	config.ErrorOutputPaths = viper.GetStringSlice("runner.logger.error_output_paths")
	config.Encoding = viper.GetString("runner.logger.encoding")

	logger, err := config.Build()
	if err != nil {
		panic("The logger could not be initialized")
	}

	defer logger.Sync() //nolint:errcheck // Ignore error

	log = zapr.NewLogger(logger)

	log.WithName("[RUNNER CONFIG]").V(1).Info("Logger initialized",
		"log_level", logLevel.String())

	return log
}

func (rn Runner) TriggerRunner() *trigger.Runner {
	return trigger.NewTriggerRunner(rn.logger, rn.nats, rn.jetstream)
}

func (rn Runner) TaskRunner() *task.Runner {
	return task.NewTaskRunner(rn.logger, rn.nats, rn.jetstream)
}

func (rn Runner) ExitRunner() *exit.Runner {
	return exit.NewExitRunner(rn.logger, rn.nats, rn.jetstream)
}

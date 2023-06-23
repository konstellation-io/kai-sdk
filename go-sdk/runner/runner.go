package runner

import (
	"fmt"
	"github.com/nats-io/nats.go"

	"github.com/go-logr/logr"
	"github.com/go-logr/zapr"
	"github.com/konstellation-io/kre-runners/go-sdk/v1/runner/exit"
	"github.com/konstellation-io/kre-runners/go-sdk/v1/runner/task"
	"github.com/konstellation-io/kre-runners/go-sdk/v1/runner/trigger"
	"github.com/spf13/viper"
	"go.uber.org/zap"
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

func initializeConfiguration() {
	//Load base configuration
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()

	//Load app configuration
	viper.SetConfigName("app")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	err = viper.MergeInConfig()

	if len(viper.AllKeys()) == 0 {
		panic(fmt.Errorf("configuration could not be loaded: %w", err))
	}

	// Load environment variables
	viper.AutomaticEnv()

	// Set viper default values
	viper.SetDefault("metadata.base_path", "/")
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

	defer logger.Sync()

	log = zapr.NewLogger(logger)

	log.WithName("[RUNNER CONFIG]").V(1).Info("Logger initialized",
		"log_level", logLevel.String())

	return log
}

func (rn Runner) TriggerRunner() *trigger.TriggerRunner {
	return trigger.NewTriggerRunner(rn.logger, rn.nats, rn.jetstream)
}

func (rn Runner) TaskRunner() *task.TaskRunner {
	return task.NewTaskRunner(rn.logger, rn.nats, rn.jetstream)
}

func (rn Runner) ExitRunner() *exit.ExitRunner {
	return exit.NewExitRunner(rn.logger, rn.nats, rn.jetstream)
}

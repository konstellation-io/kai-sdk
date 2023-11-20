package centralizedconfiguration

import (
	"errors"
	"fmt"

	"github.com/konstellation-io/kai-sdk/go-sdk/internal/common"

	"github.com/go-logr/logr"
	"github.com/nats-io/nats.go"
	"github.com/spf13/viper"

	utilErrors "github.com/konstellation-io/kai-sdk/go-sdk/internal/errors"
	"github.com/konstellation-io/kai-sdk/go-sdk/sdk/messaging"
)

var ErrKeyNotFound = errors.New("config not found in any key-value store for key")

const (
	_centralizedConfigurationLoggerName = "[CENTRALIZED CONFIGURATION]"
)

type CentralizedConfiguration struct {
	logger     logr.Logger
	globalKv   nats.KeyValue
	productKv  nats.KeyValue
	workflowKv nats.KeyValue
	processKv  nats.KeyValue
}

func NewCentralizedConfiguration(logger logr.Logger, js nats.JetStreamContext) (*CentralizedConfiguration, error) {
	wrapErr := utilErrors.Wrapper("configuration init: %w")

	globalKv, productKv, workflowKv, processKv, err := initKVStores(logger, js)
	if err != nil {
		return nil, wrapErr(err)
	}

	return &CentralizedConfiguration{
		logger:     logger,
		globalKv:   globalKv,
		productKv:  productKv,
		workflowKv: workflowKv,
		processKv:  processKv,
	}, nil
}

func initKVStores(logger logr.Logger, js nats.JetStreamContext) (
	globalKv, productKv, workflowKv, processKv nats.KeyValue, err error,
) {
	wrapErr := utilErrors.Wrapper("configuration init: %w")

	name := viper.GetString(common.ConfigCcGlobalBucketKey)
	logger.WithName(_centralizedConfigurationLoggerName).V(1).
		Info(fmt.Sprintf("Initializing global key-value store with name %s", name))

	globalKv, err = js.KeyValue(name)
	if err != nil {
		logger.Error(err, "Error initializing global key-value store")
		return nil, nil, nil, nil, wrapErr(err)
	}

	name = viper.GetString(common.ConfigCcProductBucketKey)
	logger.WithName(_centralizedConfigurationLoggerName).V(1).
		Info(fmt.Sprintf("Initializing product key-value store with name %s", name))

	productKv, err = js.KeyValue(name)
	if err != nil {
		logger.WithName(_centralizedConfigurationLoggerName).
			Error(err, "Error initializing product key-value store")
		return nil, nil, nil, nil, wrapErr(err)
	}

	logger.WithName(_centralizedConfigurationLoggerName).V(1).Info("Product key-value store initialized")

	name = viper.GetString(common.ConfigCcWorkflowBucketKey)
	logger.WithName(_centralizedConfigurationLoggerName).V(1).
		Info(fmt.Sprintf("Initializing workflow key-value store with name %s", name))

	workflowKv, err = js.KeyValue(name)
	if err != nil {
		logger.WithName(_centralizedConfigurationLoggerName).
			Error(err, "Error initializing workflow key-value store")
		return nil, nil, nil, nil, wrapErr(err)
	}

	logger.WithName(_centralizedConfigurationLoggerName).V(1).Info("Workflow key-value store initialized")

	name = viper.GetString(common.ConfigCcProcessBucketKey)
	logger.WithName(_centralizedConfigurationLoggerName).V(1).
		Info(fmt.Sprintf("Initializing process key-value store with name %s", name))

	processKv, err = js.KeyValue(name)
	if err != nil {
		logger.WithName(_centralizedConfigurationLoggerName).
			Error(err, "Error initializing process key-value store")
		return nil, nil, nil, nil, wrapErr(err)
	}

	logger.WithName(_centralizedConfigurationLoggerName).V(1).Info("Process key-value store initialized")

	return globalKv, productKv, workflowKv, processKv, nil
}

func (cc *CentralizedConfiguration) GetConfig(key string, scopeOpt ...messaging.Scope) (string, error) {
	wrapErr := utilErrors.Wrapper("configuration get: %w")

	if len(scopeOpt) > 0 {
		config, err := cc.getConfigFromScope(key, scopeOpt[0])
		if errors.Is(err, nats.ErrKeyNotFound) {
			return "", wrapErr(fmt.Errorf("%w: %q", ErrKeyNotFound, key))
		} else if err != nil {
			return "", wrapErr(err)
		}

		return config, nil
	}

	allScopesInOrder := []messaging.Scope{
		messaging.ProcessScope,
		messaging.WorkflowScope,
		messaging.ProductScope,
		messaging.GlobalScope,
	}
	for _, scope := range allScopesInOrder {
		config, err := cc.getConfigFromScope(key, scope)

		if err != nil && !errors.Is(err, nats.ErrKeyNotFound) {
			return "", wrapErr(err)
		}

		if err == nil {
			return config, nil
		}
	}

	return "", wrapErr(fmt.Errorf("%w: %q", ErrKeyNotFound, key))
}

func (cc *CentralizedConfiguration) SetConfig(key, value string, scopeOpt ...messaging.Scope) error {
	wrapErr := utilErrors.Wrapper("configuration set: %w")

	kvStore := cc.getScopedConfig(scopeOpt...)

	_, err := kvStore.PutString(key, value)
	if err != nil {
		return wrapErr(fmt.Errorf("failed to set config for key %q: %w", key, err))
	}

	return nil
}

func (cc *CentralizedConfiguration) DeleteConfig(key string, scope messaging.Scope) error {
	err := cc.getScopedConfig(scope).Delete(key)
	if err != nil {
		return fmt.Errorf("failed to delete config for key %q: %w", key, err)
	}

	return nil
}

func (cc *CentralizedConfiguration) getConfigFromScope(key string, scope messaging.Scope) (string, error) {
	value, err := cc.getScopedConfig(scope).Get(key)
	if err != nil {
		return "", fmt.Errorf("failed to get config for key %q: %w", key, err)
	}

	return string(value.Value()), nil
}

func (cc *CentralizedConfiguration) getScopedConfig(scope ...messaging.Scope) nats.KeyValue {
	if len(scope) == 0 {
		return cc.processKv
	}

	switch scope[0] {
	case messaging.GlobalScope:
		return cc.globalKv
	case messaging.ProductScope:
		return cc.productKv
	case messaging.WorkflowScope:
		return cc.workflowKv
	case messaging.ProcessScope:
		return cc.processKv
	default:
		return cc.processKv
	}
}

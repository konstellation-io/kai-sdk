package centralizedconfiguration

import (
	"errors"
	"fmt"

	"github.com/go-logr/logr"
	"github.com/nats-io/nats.go"
	"github.com/spf13/viper"

	utilErrors "github.com/konstellation-io/kre-runners/go-sdk/v1/internal/errors"
	"github.com/konstellation-io/kre-runners/go-sdk/v1/sdk/messaging"
)

type CentralizedConfiguration struct {
	logger     logr.Logger
	productKv  nats.KeyValue
	workflowKv nats.KeyValue
	processKv  nats.KeyValue
}

func NewCentralizedConfiguration(logger logr.Logger, js nats.JetStreamContext) (*CentralizedConfiguration, error) {
	wrapErr := utilErrors.Wrapper("configuration init: %w")

	logger = logger.WithName("[CENTRALIZED CONFIGURATION]")

	productKv, workflowKv, processKv, err := initKVStores(logger, js)
	if err != nil {
		return nil, wrapErr(err)
	}

	return &CentralizedConfiguration{
		logger:     logger,
		productKv:  productKv,
		workflowKv: workflowKv,
		processKv:  processKv,
	}, nil
}

func initKVStores(logger logr.Logger, jetstream nats.JetStreamContext) (
	productKv, workflowKv, processKv nats.KeyValue, err error,
) {
	wrapErr := utilErrors.Wrapper("configuration init: %w")

	logger.V(1).Info("Initializing product key-value store",
		"name", viper.GetString("centralized-configuration.product.bucket"))
	productKv, err = jetstream.KeyValue(viper.GetString("centralized-configuration.product.bucket"))
	if err != nil {
		logger.Error(err, "Error initializing product key-value store")
		return nil, nil, nil, wrapErr(err)
	}
	logger.V(1).Info("Product key-value store initialized")

	logger.V(1).Info("Initializing workflow key-value store",
		"name", viper.GetString("centralized-configuration.workflow.bucket"))
	workflowKv, err = jetstream.KeyValue(viper.GetString("centralized-configuration.workflow.bucket"))
	if err != nil {
		logger.Error(err, "Error initializing workflow key-value store")
		return nil, nil, nil, wrapErr(err)
	}
	logger.V(1).Info("Workflow key-value store initialized")

	logger.V(1).Info("Initializing process key-value store",
		"name", viper.GetString("centralized-configuration.process.bucket"))
	processKv, err = jetstream.KeyValue(viper.GetString("centralized-configuration.process.bucket"))
	if err != nil {
		logger.Error(err, "Error initializing process key-value store")
		return nil, nil, nil, wrapErr(err)
	}
	logger.V(1).Info("Process key-value store initialized")

	return productKv, workflowKv, processKv, nil
}

func (cc CentralizedConfiguration) GetConfig(key string, scopeOpt ...messaging.Scope) (string, error) {
	wrapErr := utilErrors.Wrapper("configuration get: %w")

	if len(scopeOpt) > 0 {
		config, err := cc.getConfigFromScope(key, scopeOpt[0])
		if err != nil {
			return "", wrapErr(err)
		}
		return config, nil
	}

	allScopesInOrder := []messaging.Scope{messaging.ProcessScope, messaging.WorkflowScope, messaging.ProductScope}
	for _, scope := range allScopesInOrder {
		config, err := cc.getConfigFromScope(key, scope)

		if err != nil && !errors.Is(err, nats.ErrKeyNotFound) {
			return "", wrapErr(err)
		}

		if err == nil {
			return config, nil
		}
	}

	return "", wrapErr(fmt.Errorf("error retrieving config with key %q, not found in any key-value store", key))
}

func (cc CentralizedConfiguration) SetConfig(key, value string, scopeOpt ...messaging.Scope) error {
	wrapErr := utilErrors.Wrapper("configuration set: %w")

	kvStore := cc.getScopedConfig(scopeOpt...)

	_, err := kvStore.PutString(key, value)
	if err != nil {
		return wrapErr(fmt.Errorf("error storing value with key %q to the key-value store: %w", key, err))
	}

	return nil
}

func (cc CentralizedConfiguration) DeleteConfig(key string, scope messaging.Scope) error {
	err := cc.getScopedConfig(scope).Delete(key)
	if err != nil {
		return fmt.Errorf("error retrieving config with key %q from the configuration: %w", key, err)
	}

	return nil
}

func (cc CentralizedConfiguration) getConfigFromScope(key string, scope messaging.Scope) (string, error) {
	value, err := cc.getScopedConfig(scope).Get(key)
	if err != nil {
		return "", fmt.Errorf("error retrieving config with key %q from the configuration: %w", key, err)
	}

	return string(value.Value()), nil
}

func (cc CentralizedConfiguration) getScopedConfig(scope ...messaging.Scope) nats.KeyValue {
	if len(scope) == 0 {
		return cc.processKv
	}

	switch scope[0] {
	case messaging.ProductScope:
		return cc.productKv
	case messaging.WorkflowScope:
		return cc.workflowKv
	default:
		return cc.processKv
	}
}

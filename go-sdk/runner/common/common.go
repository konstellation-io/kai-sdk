package common

import (
	"fmt"

	"github.com/spf13/viper"
	"google.golang.org/protobuf/types/known/anypb"

	kaisdk "github.com/konstellation-io/kai-sdk/go-sdk/sdk"
)

type Task func(sdk kaisdk.KaiSDK)

type Initializer Task

type Finalizer Task

type Handler func(sdk kaisdk.KaiSDK, response *anypb.Any) error

func InitializeProcessConfiguration(sdk kaisdk.KaiSDK) {
	logger := sdk.Logger.WithName("[CONFIG INITIALIZER]")
	values := viper.GetStringMapString("centralized_configuration.process.config")

	logger.V(1).Info("Initializing process configuration")

	for key, value := range values {
		err := sdk.CentralizedConfig.SetConfig(key, value)
		if err != nil {
			logger.Error(err, fmt.Sprintf("Error initializing process configuration with key %s", key))
		}
	}

	logger.V(1).Info("Process configuration initialized")
}

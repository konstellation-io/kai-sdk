package common

import (
	"github.com/spf13/viper"
	"google.golang.org/protobuf/types/known/anypb"

	kaisdk "github.com/konstellation-io/kai-sdk/go-sdk/sdk"
)

type Task func(sdk kaisdk.KaiSDK)

type Initializer Task

type Finalizer Task

type Handler func(sdk kaisdk.KaiSDK, response *anypb.Any) error

func InitializeProcessConfiguration(sdk kaisdk.KaiSDK) {
	values := viper.GetStringMapString("centralized_configuration.process.config")

	sdk.Logger.WithName("[CONFIG INITIALIZER]").V(1).Info("Initializing process configuration")

	for key, value := range values {
		err := sdk.CentralizedConfig.SetConfig(key, value)
		if err != nil {
			sdk.Logger.WithName("[CONFIG INITIALIZER]").
				Error(err, "Error initializing process configuration", "key", key)
		}

		sdk.Logger.WithName("[CONFIG INITIALIZER]").V(3).
			Info("New process configuration added", "key", key, "value", value)
	}
}

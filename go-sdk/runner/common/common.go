package common

import (
	"github.com/konstellation-io/kre-runners/go-sdk/v1/sdk"
	"github.com/spf13/viper"
	"google.golang.org/protobuf/types/known/anypb"
)

type Task func(sdk sdk.KaiSDK)

type Initializer Task

type Finalizer Task

type Handler func(sdk sdk.KaiSDK, response *anypb.Any) error

func InitializeProcessConfiguration(sdk sdk.KaiSDK) {
	values := viper.GetStringMapString("centralized-configuration.process.config")

	sdk.Logger.WithName("[CONFIG INITIALIZER]").V(1).Info("Initializing process configuration")
	if values != nil {
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
}

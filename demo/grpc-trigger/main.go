package main

import (
	"github.com/konstellation-io/kai-sdk/go-sdk/runner"
	"github.com/konstellation-io/kai-sdk/go-sdk/runner/trigger"
	"github.com/konstellation-io/kai-sdk/go-sdk/sdk"
)

func main() {
	runner.
		NewRunner().
		TriggerRunner().
		WithInitializer(func(ctx sdk.KaiSDK) {
			ctx.Logger.Info("Initializer")
		}).
		WithRunner(grpcRunner).
		WithFinalizer(func(ctx sdk.KaiSDK) {
			ctx.Logger.Info("Finalizer")
		}).
		Run()
}

func grpcRunner(tr *trigger.Runner, sdk sdk.KaiSDK) {
	sdk.Logger.Info("Starting grpc runner")

	
}

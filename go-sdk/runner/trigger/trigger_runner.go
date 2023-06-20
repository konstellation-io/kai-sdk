package trigger

import (
	"github.com/go-logr/logr"
	"github.com/konstellation-io/kre-runners/go-sdk/v1/runner/common"
	"github.com/konstellation-io/kre-runners/go-sdk/v1/sdk"
	"github.com/nats-io/nats.go"
	"google.golang.org/protobuf/types/known/anypb"
	"sync"
)

type Runner func(tr *TriggerRunner, sdk sdk.KaiSDK)

type ResponseHandler func(sdk sdk.KaiSDK, response *anypb.Any) error

type TriggerRunner struct {
	sdk              sdk.KaiSDK
	nats             *nats.Conn
	jetstream        nats.JetStreamContext
	responseHandler  ResponseHandler
	responseChannels sync.Map
	initializer      common.Initializer
	runner           Runner
	finalizer        common.Finalizer
}

var wg sync.WaitGroup

func NewTriggerRunner(logger logr.Logger, nats *nats.Conn, jetstream nats.JetStreamContext) *TriggerRunner {
	return &TriggerRunner{
		sdk:              sdk.NewKaiSDK(logger.WithName("[TRIGGER]"), nats, jetstream),
		nats:             nats,
		jetstream:        jetstream,
		responseChannels: sync.Map{},
	}
}

func (tr *TriggerRunner) WithInitializer(initializer common.Initializer) *TriggerRunner {
	tr.initializer = composeInitializer(initializer)
	return tr
}

func (tr *TriggerRunner) WithRunner(runner Runner) *TriggerRunner {
	tr.runner = composeRunner(tr, runner)
	return tr
}

func (tr *TriggerRunner) WithFinalizer(finalizer common.Finalizer) *TriggerRunner {
	tr.finalizer = composeFinalizer(finalizer)
	return tr
}

func (tr *TriggerRunner) GetResponseChannel(requestID string) <-chan *anypb.Any {
	tr.responseChannels.Store(requestID, make(chan *anypb.Any))
	channel, _ := tr.responseChannels.Load(requestID)

	return channel.(chan *anypb.Any)
}

func (tr *TriggerRunner) Run() {
	// Check required fields are initialized
	if tr.runner == nil {
		panic("Runner function is required")
	}
	if tr.initializer == nil {
		tr.initializer = composeInitializer(nil)
	}

	tr.responseHandler = getResponseHandler(&tr.responseChannels)

	if tr.finalizer == nil {
		tr.finalizer = composeFinalizer(nil)
	}

	tr.initializer(tr.sdk)

	go tr.runner(tr, tr.sdk)

	go tr.startSubscriber()

	wg.Add(2)
	wg.Wait()

	tr.finalizer(tr.sdk)
}

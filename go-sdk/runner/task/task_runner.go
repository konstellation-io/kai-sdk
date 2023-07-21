package task

import (
	"github.com/go-logr/logr"
	"github.com/nats-io/nats.go"

	"github.com/konstellation-io/kai-sdk/go-sdk/runner/common"
	"github.com/konstellation-io/kai-sdk/go-sdk/sdk"
)

type Preprocessor common.Handler

type Handler common.Handler

type Postprocessor common.Handler

type Runner struct {
	sdk              sdk.KaiSDK
	nats             *nats.Conn
	jetstream        nats.JetStreamContext
	responseHandlers map[string]Handler
	initializer      common.Initializer
	preprocessor     Preprocessor
	postprocessor    Postprocessor
	finalizer        common.Finalizer
}

func NewTaskRunner(logger logr.Logger, ns *nats.Conn, js nats.JetStreamContext) *Runner {
	return &Runner{
		sdk:              sdk.NewKaiSDK(logger.WithName("[TASK]"), ns, js),
		nats:             ns,
		jetstream:        js,
		responseHandlers: make(map[string]Handler),
	}
}

func (tr *Runner) WithInitializer(initializer common.Initializer) *Runner {
	tr.initializer = composeInitializer(initializer)
	return tr
}

func (tr *Runner) WithPreprocessor(preprocessor Preprocessor) *Runner {
	tr.preprocessor = composePreprocessor(preprocessor)
	return tr
}

func (tr *Runner) WithHandler(handler Handler) *Runner {
	tr.responseHandlers["default"] = composeHandler(handler)
	return tr
}

func (tr *Runner) WithCustomHandler(subject string, handler Handler) *Runner {
	tr.responseHandlers[subject] = composeHandler(handler)
	return tr
}

func (tr *Runner) WithPostprocessor(postprocessor Postprocessor) *Runner {
	tr.postprocessor = composePostprocessor(postprocessor)
	return tr
}

func (tr *Runner) WithFinalizer(finalizer common.Finalizer) *Runner {
	tr.finalizer = composeFinalizer(finalizer)
	return tr
}

func (tr *Runner) Run() {
	if tr.responseHandlers["default"] == nil {
		panic("No default handler defined")
	}
	if tr.initializer == nil {
		tr.initializer = composeInitializer(nil)
	}
	if tr.finalizer == nil {
		tr.finalizer = composeFinalizer(nil)
	}

	tr.initializer(tr.sdk)

	tr.startSubscriber()

	tr.finalizer(tr.sdk)
}

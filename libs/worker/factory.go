package worker

import (
	"github.com/opentracing/opentracing-go"
	"github.com/zillow/howwegoatzillow/libs/kafka"
)

type Factory struct {
	client kafka.Client
	logger Logger
	tracer opentracing.Tracer
}

// NewFactory initializes a new worker factory applying all the provided options.
func NewFactory(client kafka.Client, options ...FactoryOption) Factory {
	wf := &Factory{
		client: client,
		logger: NoopLogger{},
		tracer: opentracing.NoopTracer{},
	}

	for _, option := range options {
		if option != nil {
			option.apply(wf)
		}
	}

	return *wf
}

// Create creates a new Worker which when run will `DO` the provided work.
func (wf Factory) Create(config kafka.Config, options ...WorkerOption) *Worker {
	w := &Worker{
		client: wf.client,
		config: config,
		logger: wf.logger,
		tracer: wf.tracer,
	}

	for _, option := range options {
		if option != nil {
			option.apply(w)
		}
	}
	return w
}

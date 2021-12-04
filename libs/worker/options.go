package worker

import (
	"time"

	"github.com/opentracing/opentracing-go"
)

// RunOption interface to identify functional options that can control `worker.Run` behavior
type RunOption interface {
	apply(s *runSettings)
}

// WithSleepDuration provides option to override default sleep duration. Default is 0.
// Use this if you want some breathing room after each loop.
func WithSleepDuration(sleep time.Duration) RunOption { return sleepDurationOption{sleep: sleep} }

type sleepDurationOption struct{ sleep time.Duration }

func (l sleepDurationOption) apply(s *runSettings) { s.sleepDuration = l.sleep }

// Speedup increases the concurrencyFactor for a worker.
// concurrencyFactor is how many go routines can be running in parallel.
// NOTE: it's strongly recommended to add more worker instances rather than using this option to speed up each worker. Scale horizontally
func Speedup(times byte) RunOption { return speedupOption{times: times} }

type speedupOption struct{ times byte }

func (s speedupOption) apply(w *runSettings) {
	if s.times > 0 {
		w.concurrencyFactor = s.times
	}
}

// CircuitBreakAfter these many consequitive failures
func CircuitBreakAfter(times uint32) RunOption {
	return circuitBreakAfterOption{times: times}
}

type circuitBreakAfterOption struct{ times uint32 }

func (c circuitBreakAfterOption) apply(w *runSettings) {
	if c.times > 0 {
		w.cbAfter = c.times
	}
}

type WorkerOption interface {
	apply(w *Worker)
}

func WithWorkerLogger(l Logger) WorkerOption { return workerLoggerOption{l} }

type workerLoggerOption struct{ l Logger }

func (s workerLoggerOption) apply(wf *Worker) {
	if s.l != nil {
		wf.logger = s.l
	}
}

// FactoryOption interface to identify functional options that can control the factory and the created workers behavior
type FactoryOption interface {
	apply(s *Factory)
}

// WithLogger  provides option to override the logger to use. default is noop
func WithLogger(l Logger) FactoryOption { return loggerOption{l} }

type loggerOption struct{ l Logger }

func (s loggerOption) apply(wf *Factory) {
	if s.l != nil {
		wf.logger = s.l
	}
}

// WithTracer provides option to override the tracer to use. default is noop
func WithTracer(t opentracing.Tracer) FactoryOption {
	return tracerOption{t}
}

type tracerOption struct{ t opentracing.Tracer }

func (s tracerOption) apply(wf *Factory) {
	if s.t != nil {
		wf.tracer = s.t
	}
}

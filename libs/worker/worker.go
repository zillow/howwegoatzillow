package worker

import (
	"context"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/opentracing/opentracing-go"
	"github.com/sony/gobreaker"
	"github.com/zillow/howegoatzillow/libs/kafka"
)

type Worker struct {
	client    kafka.Client
	config    kafka.Config
	logger    Logger
	tracer    opentracing.Tracer
	wrapup    bool
	wrapupMtx sync.RWMutex
	work      *work
}

func (w *Worker) Run(ctx context.Context, processor func(context.Context, *kafka.Message) error, options ...RunOption) {
	defer func() {
		if w != nil && w.logger != nil {
			w.logger.Sync()
		}
	}()

	defer func() {
		if r := recover(); r != nil {
			w.logger.Error(ctx, "worker run failed", "recover", r)
		}
	}()

	if w.logger == nil {
		w.logger = NoopLogger{}
	}

	settings := &runSettings{
		wrapupDuration:    1 * time.Second,
		cbAfter:           5,
		cbFor:             10 * time.Second,
		concurrencyFactor: 1,
	}

	for _, option := range options {
		if option != nil {
			option.apply(settings)
		}
	}

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	stopCh := make(chan os.Signal, 1)
	signal.Notify(stopCh, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-stopCh
		if settings.wrapupDuration > 0 {
			w.setWrappingUp()
			time.Sleep(settings.wrapupDuration)
		}
		cancel()
	}()

	w.work = makeWork(w, settings, processor)

	for {
		select {
		case <-ctx.Done():
			return
		default:
			w.runSingle(ctx)
		}
		if settings.sleepDuration > 0 {
			time.Sleep(settings.sleepDuration)
		}
	}
}

func makeWork(w *Worker, settings *runSettings, processor func(context.Context, *kafka.Message) error) *work {
	cbSetting := gobreaker.Settings{}

	if settings.cbFor > 0 {
		cbSetting.Timeout = settings.cbFor
	}
	if settings.cbAfter > 0 {
		cbSetting.ReadyToTrip = func(c gobreaker.Counts) bool { return c.ConsecutiveFailures >= settings.cbAfter }
	}

	poolSize := 1
	if settings.concurrencyFactor > 0 {
		poolSize = int(settings.concurrencyFactor)
	}

	return &work{
		kconfig:        w.config,
		kclient:        w.client,
		logger:         w.logger,
		tracer:         w.tracer,
		processor:      processor,
		processTimeout: 1 * time.Minute,
		cb:             gobreaker.NewTwoStepCircuitBreaker(cbSetting),
		goroutinePool:  make(chan struct{}, poolSize),
	}
}

func (w *Worker) runSingle(ctx context.Context) {
	defer func() {
		if r := recover(); r != nil {
			w.logger.Error(ctx, "run once failed", "recover", r)
		}
	}()

	if w.isWrappingUp() {
		return
	}

	w.work.Do(ctx)

}

func (w *Worker) setWrappingUp() {
	w.wrapupMtx.Lock()
	defer w.wrapupMtx.Unlock()
	w.wrapup = true
}

func (w *Worker) isWrappingUp() bool {
	w.wrapupMtx.RLock()
	defer w.wrapupMtx.RUnlock()
	return w.wrapup
}

type runSettings struct {
	wrapupDuration time.Duration
	sleepDuration  time.Duration

	cbAfter           uint32
	cbFor             time.Duration
	concurrencyFactor byte
}

package worker

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"
	"github.com/sony/gobreaker"
	"github.com/zillow/howegoatzillow/libs/kafka"
)

type work struct {
	kconfig        kafka.Config
	kclient        kafka.Client
	processor      func(context.Context, *kafka.Message) error
	logger         Logger
	tracer         opentracing.Tracer
	rdrMtx         sync.RWMutex
	reader         kafka.Reader
	goroutinePool  chan struct{}
	cb             *gobreaker.TwoStepCircuitBreaker
	processTimeout time.Duration
}

func (w *work) Do(ctx context.Context) {
	cfg := w.kconfig

	if err := w.ensureReader(ctx); err != nil {
		err := errors.Wrap(err, "failed to get kafka reader")
		w.logger.Error(ctx, "failed to get kafka reader",
			"error", err,
			"cfg", cfg)
		return
	}

	err := w.do(ctx)

	if err != nil {
		err := errors.Wrap(err, "kafka topic processing failed")
		w.logger.Error(ctx, "kafka topic processing failed",
			"error", err,
			"cfg", cfg)
	}
}

func (w *work) do(ctx context.Context) error {
	successFunc, err := w.cb.Allow()

	// If circuit is open, Allow() returns error.
	// If circuit is open, we dont read.
	if err != nil {
		return nil
	}

loop:
	for {
		select {
		case w.goroutinePool <- struct{}{}:

			msg, err := w.reader.Read(ctx)
			if err != nil {
				successFunc(false)
				<-w.goroutinePool
				return errors.Wrap(err, "failed to read from kafka topic")
			}
			go func(i *kafka.Message) {
				err = w.doSingle(ctx, i)
				successFunc(err == nil)
				<-w.goroutinePool
			}(msg)
		default:
			break loop
		}
	}
	return nil
}

func (w *work) doSingle(ctx context.Context, msg *kafka.Message) (err error) {
	defer func() {
		if r := recover(); r != nil {
			//Panic for one message should not bring down the worker. Log and continue
			w.logger.Error(ctx, "kafka topic single message processing panicked",
				"recover", r,
				"msg", msg,
			)
			switch x := r.(type) {
			case error:
				err = x
			default:
				err = errors.New("kafka topic single message processing panicked")
			}
		}
	}()

	// send the done signal. Always do this.
	// otherwise the message won't be committed
	defer func() {
		msg.Done()
	}()

	operationName := fmt.Sprintf("kafka.work.process.message.%s", w.kconfig.Topic)

	sc, _ := w.tracer.Extract(opentracing.TextMap, &ReadAttributeCarrier{Message: msg})
	span, ctxNew := opentracing.StartSpanFromContext(ctx, operationName, opentracing.ChildOf(sc))
	ctxNew, cancel := context.WithTimeout(ctxNew, w.processTimeout)
	defer cancel()

	err = func() error {
		defer span.Finish()
		errorCh := make(chan error, 1)
		go func() {
			defer func() {
				if r := recover(); r != nil {
					switch x := r.(type) {
					case error:
						errorCh <- x
					default:
						errorCh <- errors.New("kafka topic single message processing panicked")
					}
				}
			}()
			errorCh <- w.processor(ctxNew, msg)
		}()
		select {
		case err2 := <-errorCh:
			return err2
		case <-ctx.Done():
			return errors.New("timeout occurred during kafka process")
		}
	}()

	if err == nil {
		return
	}

	// In case of error, we don't mark it done. It will come back or eventually go to the dead letter queue
	w.logger.Error(ctx, "kafka topic single message processing failed",
		"error", err,
		"msg", msg,
	)
	return err
}

func (w *work) ensureReader(ctx context.Context) error {

	w.rdrMtx.RLock()
	if w.reader != nil {
		w.rdrMtx.RUnlock()
		return nil
	}
	w.rdrMtx.RUnlock()

	w.rdrMtx.Lock()
	defer w.rdrMtx.Unlock()

	rdr, err := w.kclient.Reader(ctx, w.kconfig)
	if err != nil {
		return err
	}

	if rdr == nil {
		return errors.New("nil reader received")
	}

	w.reader = rdr

	return nil
}

type ReadAttributeCarrier struct{ Message *kafka.Message }

// ForeachKey conforms to the opentracing TextMapReader interface.
func (c *ReadAttributeCarrier) ForeachKey(handler func(key, val string) error) error {
	for k, v := range c.Message.Headers {
		val := str(v)
		if err := handler(k, val); err != nil {
			return err
		}
	}
	return nil
}

func str(v interface{}) string {
	if vs, ok := v.([]byte); ok {
		return string(vs)
	} else {
		return fmt.Sprintf("%v", v)
	}
}

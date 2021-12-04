package kafka

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
)

// Client ...
type Client interface {
	Reader(ctx context.Context, topicConfig Config) (Reader, error)
	Writer(ctx context.Context, topicConfig Config) (Writer, error)
}

// Reader ...
type Reader interface {
	Read(ctx context.Context) (*Message, error)
}

// Writer ...
type Writer interface {
	Write(ctx context.Context, key string, value []byte) (Response, error)
}

// Config ...
type Config struct {
	Topic            string
	BootstrapServers []string
}

// Message ...
type Message struct {
	Key       string
	Headers   map[string]string
	Offset    int32
	Partition int32
	value     []byte
}

// Response ...
type Response struct {
	Partition int32
	Offset    int64
}

// Done ...
func (m *Message) Done() {
	// Commit offset
}

// Logger ...
type Logger interface {
	Info(ctx context.Context, msg string, keysAndValues ...interface{})
	Error(ctx context.Context, msg string, keysAndValues ...interface{})
}

type client struct {
	config Config
	tracer opentracing.Tracer
	logger Logger
}

func (c *client) Reader(ctx context.Context, topicConfig Config) (Reader, error) { return nil, nil }

func (c *client) Writer(ctx context.Context, topicConfig Config) (Writer, error) {
	return &writer{c.tracer}, nil
}

// NewClient ...
func NewClient(config Config, tracer opentracing.Tracer, logger Logger) Client {
	return &client{config: config, tracer: tracer, logger: logger}
}

type writer struct {
	tracer opentracing.Tracer
}

func (w *writer) Write(ctx context.Context, key string, value []byte) (Response, error) {
	msg := &Message{
		Key:   key,
		value: value,
		Headers: map[string]string{
			"timestamp": time.Now().String(),
			"guid":      uuid.New().String(),
		}}

	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, w.tracer, "kafka_write", ext.SpanKindProducer)
	defer span.Finish()
	_ = w.tracer.Inject(span.Context(), opentracing.TextMap, &writeAttributeCarrier{msg})
	var resp Response
	//resp = //Make your actual kafka call here
	return resp, nil
}

type writeAttributeCarrier struct{ msg *Message }

// Set conforms to the TextMapWriter interface.
func (c *writeAttributeCarrier) Set(key, val string) {
	c.msg.Headers[key] = val
}

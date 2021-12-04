package http

import (
	"net/http"
	"time"

	"github.com/hashicorp/go-retryablehttp"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
)

type Provider struct {
	tracer        opentracing.Tracer
	leveledLogger LeveledLogger
}

func NewClientProvider(tracer opentracing.Tracer, l LeveledLogger) Provider {
	return Provider{tracer, l}
}

// Config ...
type Config struct {
	TimeoutMs      *int
	RetryWaitMinMs *int
	RetryMax       *int
}

func (p *Provider) GetClient(cfg Config) *http.Client {
	rClient := retryablehttp.NewClient()
	rClient.Logger = p.leveledLogger

	rClient.RetryMax = 0
	if cfg.RetryMax != nil {
		rClient.RetryMax = *cfg.RetryMax
	}

	if cfg.RetryWaitMinMs != nil {
		rClient.RetryWaitMin = time.Duration(*cfg.RetryWaitMinMs) * time.Millisecond
	}

	client := rClient.StandardClient()
	client.Timeout = 10 * time.Second
	if cfg.TimeoutMs != nil {
		client.Timeout = time.Duration(*cfg.TimeoutMs) * time.Millisecond
	}
	return client
}

type HttpClientWrapper struct {
	*http.Client
	tracer opentracing.Tracer
}

func (w *HttpClientWrapper) Do(request *http.Request) (*http.Response, error) {
	ctx := request.Context()

	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, w.tracer, "http-request", ext.SpanKindRPCClient)
	defer span.Finish()

	ext.HTTPMethod.Set(span, request.Method)
	ext.HTTPUrl.Set(span, request.URL.String())

	w.tracer.Inject(span.Context(), opentracing.HTTPHeaders, opentracing.HTTPHeadersCarrier(request.Header))

	request = request.WithContext(ctx)

	resp, err := w.Client.Do(request)
	if err != nil {
		span.SetTag("error", true)
		return resp, err
	}
	ext.HTTPStatusCode.Set(span, uint16(resp.StatusCode))

	return resp, err
}

func (p *Provider) GetWrappedClient(cfg Config) *HttpClientWrapper {
	return &HttpClientWrapper{Client: p.GetClient(cfg), tracer: p.tracer}
}

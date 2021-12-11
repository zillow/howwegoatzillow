package server

import (
	"net/http"

	"github.com/opentracing/opentracing-go"
)

// Option interface to identify functional options
type Option interface{ apply(r *Server) }

// WithServerLogger provides option to provide a logger to use while writing.
func WithServerLogger(l Logger) Option { return serverLoggerOption{l} }

// WithServerTracer provides option to provide a tracer to use while writing.
func WithServerTracer(t opentracing.Tracer) Option { return serverTracerOption{t} }

// WithServerConfig provides option to provide a server configuration.
func WithServerConfig(c Config) Option { return serverConfigOption{c} }

// WithServerPort provides option to provide the port on which the server listens. Default is 80
func WithServerPort(p int) Option { return serverPortOption{p} }

// WithServerReadTimeout provides option to provide the maximum duration in milliseconds for reading the entire
// request, including the body.
// defaults to 10 seconds
func WithServerReadTimeout(t int) Option { return serverReadTimeoutOption{t} }

// WithServerWriteTimeout provides option to provide the maximum duration in milliseconds before timing out writes of the response.
// defaults to 10 seconds
func WithServerWriteTimeout(t int) Option { return serverWriteTimeoutOption{t} }

// WithShutdownDelaySeconds provides option to provide the duration by which server shutdown is delayed after receiving an os signal.
// defaults to 5 seconds
func WithShutdownDelaySeconds(d int) Option { return serverShutdownDelaySecondsOption{d} }

// WithHealthCheck provides option to provide additional health checks that are performed on health check probe.
func WithHealthCheck(f func(http.HandlerFunc) http.HandlerFunc) Option {
	return serverHealthCheckOption{f}
}

// WithLivenessCheck provides option to provide additional liveness checks that are performed on liveness probe.
func WithLivenessCheck(f func(http.HandlerFunc) http.HandlerFunc) Option {
	return serverLivenessCheckOption{f}
}

// WithReadinessCheck provides option to provide additional readiness checks that are performed on readiness probe.
func WithReadinessCheck(f func(http.HandlerFunc) http.HandlerFunc) Option {
	return serverReadinessCheckOption{f}
}

// WithSwaggerFile provides option to provide the swagger file location. Default is '/swagger.json'
func WithSwaggerFile(f string) Option { return serverSwaggerFileOption{f} }

// WithServerRouter provides option to provide hooks to use the http request to mutate the request context.
func WithServerRouter(r Handler) Option {
	return serverRouterOption{r: r}
}

type serverLoggerOption struct{ logger Logger }

func (l serverLoggerOption) apply(s *Server) {
	if l.logger != nil {
		s.logger = l.logger
	}
}

type serverTracerOption struct{ tracer opentracing.Tracer }

func (t serverTracerOption) apply(s *Server) {
	if t.tracer != nil {
		s.tracer = t.tracer
	}
}

type serverConfigOption struct{ config Config }

func (m serverConfigOption) apply(s *Server) {
	s.config = m.config
}

type serverPortOption struct{ port int }

func (p serverPortOption) apply(s *Server) {
	s.config.Port = p.port
}

type serverSwaggerFileOption struct{ f string }

func (p serverSwaggerFileOption) apply(s *Server) {
	s.config.SwaggerFile = p.f
}

type serverReadTimeoutOption struct{ t int }

func (p serverReadTimeoutOption) apply(s *Server) {
	s.config.ReadTimeoutMs = p.t
}

type serverWriteTimeoutOption struct{ t int }

func (p serverWriteTimeoutOption) apply(s *Server) {
	s.config.WriteTimeoutMs = p.t
}

type serverShutdownDelaySecondsOption struct{ t int }

func (p serverShutdownDelaySecondsOption) apply(s *Server) {
	s.config.ShutdownDelaySeconds = p.t
}

type serverReadinessCheckOption struct {
	f func(http.HandlerFunc) http.HandlerFunc
}

func (r serverReadinessCheckOption) apply(s *Server) {
	s.readinessCheck = r.f
}

type serverLivenessCheckOption struct {
	f func(http.HandlerFunc) http.HandlerFunc
}

func (l serverLivenessCheckOption) apply(s *Server) {
	s.livenessCheck = l.f
}

type serverHealthCheckOption struct {
	f func(http.HandlerFunc) http.HandlerFunc
}

func (h serverHealthCheckOption) apply(s *Server) {
	s.healthCheck = h.f
}

type serverRouterOption struct {
	r Handler
}

func (sro serverRouterOption) apply(s *Server) {
	if sro.r != nil {
		s.Router = sro.r
	}
}

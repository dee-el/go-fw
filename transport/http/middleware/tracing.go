package middleware

import (
	"fmt"
	"net/http"

	chi "github.com/go-chi/chi/v5"
	chi_middleware "github.com/go-chi/chi/v5/middleware"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
)

type Tracing struct {
	tracer opentracing.Tracer
}

func NewTracing(tracer opentracing.Tracer) *Tracing {
	return &Tracing{
		tracer: tracer,
	}
}

func (t *Tracing) Handler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// since using chi as template mux http, utilize chi context to get regexed pattern path
		chiCtx := chi.RouteContext(r.Context())
		path := chiCtx.RoutePattern()

		// but if URL path called not the registered one, it will be empty
		// just get the path from URL.Path
		// it can be mean DDoS attack happened or else ...
		if path == "" {
			path = r.URL.Path
		}

		operationName := fmt.Sprintf("HTTP %s: %s", r.Method, path)
		spanCtx, _ := t.tracer.Extract(opentracing.HTTPHeaders, opentracing.HTTPHeadersCarrier(r.Header))

		// Start the global span, it'll wrap the request lifecycle.
		var span opentracing.Span
		if spanCtx != nil {
			span = t.tracer.StartSpan(operationName, opentracing.ChildOf(spanCtx))
		} else {
			span = t.tracer.StartSpan(operationName)
		}

		rw := chi_middleware.NewWrapResponseWriter(w, r.ProtoMajor)
		next.ServeHTTP(rw, r.WithContext(opentracing.ContextWithSpan(r.Context(), span)))

		defer func() {

			// Save request metadata for this span. Note that tags are searchable on UI.
			span.SetTag("component", "net/http")
			span.SetTag(string(ext.HTTPMethod), r.Method)
			span.SetTag(string(ext.HTTPUrl), r.URL.Path)
			span.SetTag(string(ext.HTTPStatusCode), rw.Status())
			span.SetTag("http.headers", getHeaders(r.Header))
			span.SetTag("http.request_ip", getIP(r))
			span.Finish()
		}()
	})
}

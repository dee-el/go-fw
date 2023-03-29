package middleware

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	chi "github.com/go-chi/chi/v5"
	chi_middleware "github.com/go-chi/chi/v5/middleware"

	"github.com/dee-el/go-fw/metrics"
)

type Metrics struct {
	// all metrics will be recorded by this.
	recorder metrics.InboundHTTPRecorder
	// GroupedStatus toggle for grouping HTTP status in the form of `\dxx`.
	// Example: 200, 201, and 203 will have the label `code="2xx"`.
	// Why grouped HTTP status?
	// This impacts on the cardinality of the metrics and also improves the performance of queries
	// that are grouped by status code because there are already aggregated in the metric.
	// By default will be false.
	groupedStatus bool
}

func NewMetrics(recorder metrics.InboundHTTPRecorder, opts ...MetricsOption) *Metrics {
	if recorder == nil {
		return nil
	}

	i := &Metrics{
		recorder:      recorder,
		groupedStatus: false,
	}

	for _, opt := range opts {
		opt(i)
	}

	return i
}

type MetricsOption func(*Metrics)

func GroupedStatusMetricsOption(groupedStatus bool) MetricsOption {
	return func(i *Metrics) {
		i.groupedStatus = groupedStatus
	}
}

func (i *Metrics) Handler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rw := chi_middleware.NewWrapResponseWriter(w, r.ProtoMajor)
		now := time.Now()
		next.ServeHTTP(rw, r)
		defer func() {
			// since using chi as template mux http, utilize chi context to get regexed pattern path
			ctx := chi.RouteContext(r.Context())
			path := ctx.RoutePattern()

			// if URL path called not the registered one, it will be empty
			// it can be mean DDoS attack happened or client hit deprecated routes
			if path == "" {
				path = r.URL.Path
			}

			method := r.Method

			var code string
			code = strconv.Itoa(rw.Status())
			if i.groupedStatus {
				code = fmt.Sprintf("%dxx", rw.Status()/100)
			}

			duration := time.Since(now)
			i.recorder.Record(path, method, code, duration)
		}()
	})
}
